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
	counter          int
	counterLock      sync.Mutex
	s                *grpc.Server
)

// FrontServer is used to implement the MiniServer interface
type FrontServer struct {
	pb.UnimplementedFrontServer
}

type ActiveWorker struct {
	Info    master.WorkerInfo
	Active  bool
	Client  *pb.BackClient
	Indexes []int
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

func (s *FrontServer) ConvolutionalLayer(ctx context.Context, in *pb.ConvolutionalLayerFrontRequest) (*pb.ConvolutionalLayerFrontReply, error) {
	// generate ID
	counterLock.Lock()
	counter += 1
	var id = counter
	counterLock.Unlock()

	log.Printf("[Front]: Request #%d started.", id)

	var actualParallelWorkers int
	var actualReplicationGrade int
	if in.GetUseKernels() {
		actualParallelWorkers = min(*ParallelWorkers, master.GetWorkersLen(), len(in.Kernel))
	} else {
		actualParallelWorkers = 1
	}
	actualReplicationGrade = min(*ReplicationGrade, actualParallelWorkers)

	log.Printf("[Front]: Parallel workers: %d, Replication grade: %d", actualParallelWorkers, actualReplicationGrade)

	// Split kernels into blocks
	kernelsMap := make(map[int][]*pb.Matrix)
	var kernelBlockSize int = int(math.Floor(float64(len(in.Kernel)) / float64(actualParallelWorkers)))
	for i := 0; i < actualParallelWorkers-1; i++ {
		kernelsMap[i] = in.Kernel[i*kernelBlockSize : min((i+1)*kernelBlockSize, len(in.Kernel))]
	}
	kernelsMap[actualParallelWorkers-1] = in.Kernel[(actualParallelWorkers-1)*kernelBlockSize : len(in.Kernel)]

	// Get workers and assign
	workers := master.GetWorkers(actualParallelWorkers)
	activeWorker := make([]ActiveWorker, actualParallelWorkers)
	// For each worker
	for i := 0; i < len(workers); i++ {
		// Compute blocks indexes
		workerBlocks := assignIndexes(actualParallelWorkers, actualReplicationGrade, i)

		// Append all kernels to produce the list
		var workerKernels []*pb.Matrix
		for j := 0; j < len(workerBlocks); j++ {
			block := kernelsMap[workerBlocks[j]]
			workerKernels = append(workerKernels, block...)
		}

		// Create the active worker
		activeWorker[i] = ActiveWorker{
			Info:    workers[i],
			Active:  false,
			Client:  nil,
			Indexes: workerBlocks,
			Request: &pb.ConvolutionalLayerBackRequest{
				Target:      in.GetTarget(),
				Kernel:      workerKernels,
				AvgPoolSize: in.GetAvgPoolSize(),
				UseKernels:  in.GetUseKernels(),
				UseSigmoid:  in.GetUseSigmoid(),
				ID:          int32(id)},
		}
	}

	// Test code
	for i := 0; i < len(activeWorker); i++ {
		log.Printf("[Front]: Worker %d: %s", i, activeWorker[i].Info.WorkerAddress)
		log.Printf("[Front]: Indexes: %v", activeWorker[i].Indexes)
		log.Printf("[Front]: Request.Kernel: %v", activeWorker[i].Request.Kernel)
	}

	// Create the back request
	backRequest := pb.ConvolutionalLayerBackRequest{
		Target:      in.GetTarget(),
		Kernel:      in.GetKernel(),
		AvgPoolSize: in.GetAvgPoolSize(),
		UseKernels:  in.GetUseKernels(),
		UseSigmoid:  in.GetUseSigmoid(),
		ID:          int32(id)}

	// Allocate space for the response
	response := pb.ConvolutionalLayerFrontReply{}

	for {
		// get worker
		workers := master.GetWorkers(*ParallelWorkers)
		activeWorker := make([]ActiveWorker, len(workers))
		for i := 0; i < len(workers); i++ {
			activeWorker[i] = ActiveWorker{
				Info:   workers[i],
				Active: false,
				Client: nil,
			}
		}

		// Prepare the channel for sycnhronization
		successChan := make(chan bool)

		// time the call
		startTime := time.Now()

		for i := 0; i < len(workers); i++ {
			go func(i int, actualReply *pb.ConvolutionalLayerFrontReply) {
				// Select the worker
				addr := workers[i].WorkerAddress

				// Set up a connection to the gRPC server
				conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Printf("[Front]: %s did not connect: %v", addr, err)
					master.RemoveWorker(addr, "not connected")
					successChan <- false
					return
				}
				defer conn.Close()

				// create the client object
				c := pb.NewBackClient(conn)
				activeWorker[i].Client = &c

				// create the context
				ctxInternal, cancel := context.WithTimeout(context.Background(), time.Minute*5)
				defer cancel()

				// Set workers as active
				activeWorker[i].Active = true

				// contact the server
				WorkerResponse, err := c.ConvolutionalLayer(ctxInternal, &backRequest)

				// Set workers as inactive
				activeWorker[i].Active = false

				// Check if error
				if err != nil {
					// Check if reply is not avaible due to size too big
					s, ok := status.FromError(err)
					if ok && s.Code() == codes.ResourceExhausted {
						log.Printf("[Front]: Response: #%d: Resource exhausted error occurred: %s by %s.", id, s.Message(), addr)
					} else if ok && s.Code() == codes.Canceled {
						log.Printf("[Front]: Response: #%d: %s by %s.", id, s.Message(), addr)
					} else {
						log.Printf("[Front]: %s is unreachable: %v", addr, err)
						master.RemoveWorker(addr, "unreachable")
					}
					// Notify failure
					successChan <- false
					return
				}

				// Setup actual reply
				response.ID = int32(WorkerResponse.GetID())
				response.Result = WorkerResponse.GetResult()

				// Notify success
				successChan <- true
			}(i, &response)
		}

		for failedWorkes := 0; failedWorkes < len(workers); {
			outcome := <-successChan
			if outcome {
				log.Printf("[Front]: Response: #%d (%d); Time spent: %d ms", id, response.GetID(), time.Since(startTime).Milliseconds())
				// Cancel the request for the active workers
				for i := 0; i < len(workers); i++ {
					if activeWorker[i].Active {
						go func(i int) {
							// create the context
							ctxInternal, cancel := context.WithTimeout(context.Background(), time.Minute*5)
							defer cancel()

							c := *activeWorker[i].Client
							c.Cancel(ctxInternal, &pb.CancelRequest{ID: int32(id)})
						}(i)
					}
				}
				return &response, nil
			} else {
				failedWorkes += 1
			}
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
