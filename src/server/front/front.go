package front

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	pb "mini/proto"
	"mini/server/master"
	"net"
	"sort"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	FrontPort        = flag.String("FrontPort", "", "The port for front service.")
	ParallelWorkers  = flag.Int("ParallelWorkers", 2, "The number of parallel workers.")
	ReplicationGrade = flag.Int("ReplicationGrade", 2, "The replication grade.")
	counter          int32
	counterLock      sync.Mutex
	s                *grpc.Server
)

// FrontServer is used to implement the MiniServer interface
type FrontServer struct {
	pb.UnimplementedFrontServer
}

type ActiveWorker struct {
	ID      int
	Info    master.WorkerInfo
	Active  bool
	Ready   bool
	Client  *pb.BackClient
	Request *pb.ConvolutionalLayerBackRequest
}

func assignIndexes(parallelWorkers int, replicationGrade int, index int) []int {
	result := make([]int, replicationGrade)
	for i := 0; i < replicationGrade; i++ {
		result[i] = (index + i) % parallelWorkers
	}
	sort.Ints(result)
	return result
}

func checkCompleted(blocks []int, completedBlocks map[int]bool, lock *sync.RWMutex) bool {
	lock.Lock()
	defer lock.Unlock()
	for _, block := range blocks {
		if !completedBlocks[block] {
			return false
		}
	}
	return true
}

func setReady(workersMapsLock *sync.RWMutex, workersMap map[int]ActiveWorker, i int, ready bool) {
	workersMapsLock.RLock()
	entry := workersMap[i]
	workersMapsLock.RUnlock()
	entry.Ready = ready
	workersMapsLock.Lock()
	workersMap[i] = entry
	workersMapsLock.Unlock()
}

func (s *FrontServer) ConvolutionalLayer(ctx context.Context, in *pb.ConvolutionalLayerFrontRequest) (*pb.ConvolutionalLayerFrontReply, error) {
	// generate ID and InnerID
	counterLock.Lock()
	counter += 1
	var id int32 = counter
	counterLock.Unlock()

	// Compute parallelization parameters
	var actualParallelWorkers int
	var actualReplicationGrade int
	if in.GetUseKernels() {
		actualParallelWorkers = min(*ParallelWorkers, master.GetWorkersLen(), len(in.Kernel))
	} else {
		actualParallelWorkers = 1
	}
	actualReplicationGrade = min(*ReplicationGrade, actualParallelWorkers)

	log.Printf("[Front]: Request #%d started. Parallel workers: %d, Replication grade: %d", id, actualParallelWorkers, actualReplicationGrade)

	// Split kernels into blocks
	kernelsMap := make(map[int][]*pb.Matrix)
	var kernelBlockSize int = int(math.Floor(float64(len(in.Kernel)) / float64(actualParallelWorkers)))
	for i := 0; i < actualParallelWorkers-1; i++ {
		kernelsMap[i] = in.Kernel[i*kernelBlockSize : min((i+1)*kernelBlockSize, len(in.Kernel))]
	}
	kernelsMap[actualParallelWorkers-1] = in.Kernel[(actualParallelWorkers-1)*kernelBlockSize : len(in.Kernel)]

	// Join blocks into splits
	indexesMap := make(map[int][]int)
	backRequestMap := make(map[int]*pb.ConvolutionalLayerBackRequest)
	for i := 0; i < actualParallelWorkers; i++ {
		// Compute blocks indexes
		indexesMap[i] = assignIndexes(actualParallelWorkers, actualReplicationGrade, i)
		splitBlocks := indexesMap[i]

		// Append all kernels to produce the list
		var split []*pb.Matrix
		for j := 0; j < len(splitBlocks); j++ {
			block := kernelsMap[splitBlocks[j]]
			split = append(split, block...)
		}

		// Create back request
		backRequestMap[i] = &pb.ConvolutionalLayerBackRequest{
			Target:      in.GetTarget(),
			Kernel:      split,
			AvgPoolSize: in.GetAvgPoolSize(),
			UseKernels:  in.GetUseKernels(),
			UseSigmoid:  in.GetUseSigmoid(),
			ID:          id,
			InnerID:     int32(i)}
	}

	// Get first workers
	workersMapsLock := sync.RWMutex{}
	workers := master.GetWorkers(actualParallelWorkers, nil)
	workersMap := make(map[int]ActiveWorker)
	for i := 0; i < actualParallelWorkers; i++ {
		workersMap[i] = ActiveWorker{
			ID:     i,
			Info:   workers[i],
			Active: false,
			Ready:  true,
		}
	}

	// time the call
	startTime := time.Now()

	successChan := make(chan int)
	replaceChan := make(chan int)
	fatalChannel := make(chan error)
	readyChan := make([]chan bool, actualParallelWorkers)
	for i := 0; i < actualParallelWorkers; i++ {
		readyChan[i] = make(chan bool)
	}

	// Initialize completedBlocks and completedResult
	completeBlocksLock := sync.RWMutex{}
	completedBlocks := make(map[int]bool, actualParallelWorkers)
	for i := 0; i < actualParallelWorkers; i++ {
		completedBlocks[i] = false
	}
	completedResult := make([]*pb.Matrix, len(in.Kernel))

	// For each part
	for i := 0; i < actualParallelWorkers; i++ {
		go func(i int) {
			// Call RPC to worker
			for !checkCompleted(indexesMap[i], completedBlocks, &completeBlocksLock) {
				// Select the worker
				var activeWorker ActiveWorker

				workersMapsLock.RLock()
				ready := workersMap[i].Ready
				workersMapsLock.RUnlock()

				// Wait for ready, in case of replacing
				if !ready {
					log.Printf("[Front]: Waiting for replace %d.%d", id, i)
					<-readyChan[i]
				}

				workersMapsLock.RLock()
				addr := workersMap[i].Info.WorkerAddress
				workersMapsLock.RUnlock()

				log.Printf("[Front]: Sending %d.%d (%v) to %s", id, i, indexesMap[i], addr)

				// Set up a connection to the gRPC server
				conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Printf("[Front]: %s did not connect: %v", addr, err)
					master.RemoveWorker(addr, "not connected")
					setReady(&workersMapsLock, workersMap, i, false)
					replaceChan <- i
					continue
				}
				defer conn.Close()

				// create the client object
				c := pb.NewBackClient(conn)

				workersMapsLock.RLock()
				activeWorker = workersMap[i]
				workersMapsLock.RUnlock()
				activeWorker.Client = &c
				workersMapsLock.Lock()
				workersMap[i] = activeWorker
				workersMapsLock.Unlock()

				// create the context
				ctxInternal, cancel := context.WithTimeout(context.Background(), time.Minute*5)
				defer cancel()

				// Set workers as active
				workersMapsLock.RLock()
				activeWorker = workersMap[i]
				workersMapsLock.RUnlock()
				activeWorker.Active = true
				workersMapsLock.Lock()
				workersMap[i] = activeWorker
				workersMapsLock.Unlock()

				// contact the server
				WorkerResponse, err := c.ConvolutionalLayer(ctxInternal, backRequestMap[i])

				// Set workers as inactive
				workersMapsLock.RLock()
				activeWorker = workersMap[i]
				workersMapsLock.RUnlock()
				activeWorker.Active = false
				workersMapsLock.Lock()
				workersMap[i] = activeWorker
				workersMapsLock.Unlock()

				// Check if error
				if err != nil {
					// Check if reply is not avaible due to size too big
					s, ok := status.FromError(err)
					if ok && s.Code() == codes.ResourceExhausted {
						log.Printf("[Front]: Response: %d.%d: Resource exhausted error occurred: %s by %s.", id, i, s.Message(), addr)
						setReady(&workersMapsLock, workersMap, i, false)
						fatalChannel <- err
					} else if ok && s.Code() == codes.Canceled {
						log.Printf("[Front]: Cancelled %d.%d (%v) to %s", id, i, indexesMap[i], addr)
					} else {
						log.Printf("[Front]: %s is unreachable: %v", addr, err)
						master.RemoveWorker(addr, "unreachable")
						setReady(&workersMapsLock, workersMap, i, false)
						replaceChan <- i
					}
					continue
				}

				// Set completed
				log.Printf("[Front]: Received %d.%d (%v) to %s", id, i, indexesMap[i], addr)
				completeBlocksLock.Lock()
				newCompletedBlocks := indexesMap[i]
				for j := 0; j < len(newCompletedBlocks); j++ {
					completedBlockIndex := newCompletedBlocks[j]
					completedBlocks[completedBlockIndex] = true

					var begin int = completedBlockIndex * kernelBlockSize
					var end int
					if completedBlockIndex == actualParallelWorkers-1 {
						end = len(in.Kernel)
					} else {
						end = begin + kernelBlockSize
					}
					for k := begin; k < end; k++ {
						completedResult[k] = WorkerResponse.GetResult()[k-begin]
					}
				}
				completeBlocksLock.Unlock()

				// Notify outer
				successChan <- i
			}
		}(i)
	}

	for {
		select {
		// Wait for completion
		case <-successChan:
			// Check if all completed
			allCompleted := true
			completeBlocksLock.RLock()
			for i := 0; i < actualParallelWorkers; i++ {
				if !completedBlocks[i] {
					allCompleted = false
					break
				}
			}
			completeBlocksLock.RUnlock()

			if allCompleted {
				// Cancel pending requests
				for i := 0; i < actualParallelWorkers; i++ {
					if workersMap[i].Active {
						go func(i int) {
							// create the context
							ctxInternal, cancel := context.WithTimeout(context.Background(), time.Minute*5)
							defer cancel()
							c := *workersMap[i].Client
							c.Cancel(ctxInternal, &pb.CancelRequest{ID: id, InnerID: int32(i)})
						}(i)
					}
				}

				log.Printf("[Front]: Response: %d complete in %d ms", id, time.Since(startTime).Milliseconds())
				return &pb.ConvolutionalLayerFrontReply{Result: completedResult}, nil
			}

		case replaceIndex := <-replaceChan:
			log.Printf("[Front]: Replacing %d.%d", id, replaceIndex)
			newWorker := master.GetWorkers(1, workers)

			workersMapsLock.RLock()
			workerMapEntry := workersMap[replaceIndex]
			workersMapsLock.RUnlock()

			workerMapEntry.Info = newWorker[0]
			workerMapEntry.Ready = true

			workersMapsLock.Lock()
			workersMap[replaceIndex] = workerMapEntry
			workersMapsLock.Unlock()

			log.Printf("[Front]: Replaced %d.%d with %s", id, replaceIndex, newWorker[0].WorkerAddress)
			readyChan[replaceIndex] <- true

		case err := <-fatalChannel:
			return nil, err
		}
	}
}

func StartFrontServer() {
	// listen to request to specified port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *FrontPort))
	if err != nil {
		log.Fatalf("[Front]: Failed to listen: %v", err)
	}

	// create a new server
	s = grpc.NewServer()

	// register the server
	pb.RegisterFrontServer(s, &FrontServer{})
	log.Printf("[Front]: Listening at %v", lis.Addr())

	// serve the request
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[Front]: Failed to serve: %v", err)
	}
}

func StopServer(wg *sync.WaitGroup) {
	log.Printf("[Front]: Grafecully stopping...")

	// Graceful stop
	s.GracefulStop()
	log.Printf("[Front]: Done.")

	// Comunicate on channel so sync
	(*wg).Done()
}
