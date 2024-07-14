package front

import (
	"context"
	"flag"
	"fmt"
	"log"
	pb "mini/proto"
	"mini/server/master"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	FrontPort       = flag.String("FrontPort", "", "The port for front service.")
	ParallelWorkers = flag.Int("ParallelWorkers", 2, "The number of parallel workers.")
	counter         int
	counterLock     sync.Mutex
	s               *grpc.Server
)

// FrontServer is used to implement the MiniServer interface
type FrontServer struct {
	pb.UnimplementedFrontServer
}

func (s *FrontServer) ConvolutionalLayer(ctx context.Context, in *pb.ConvolutionalLayerFrontRequest) (*pb.ConvolutionalLayerFrontReply, error) {
	// generate ID
	counterLock.Lock()
	counter += 1
	var id = counter
	counterLock.Unlock()

	log.Printf("[Front]: Request #%d started.", id)

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

				// create the context
				ctxInternal, cancel := context.WithTimeout(context.Background(), time.Minute*5)
				defer cancel()

				// contact the server
				WorkerResponse, err := c.ConvolutionalLayer(ctxInternal, &backRequest)

				// Check if error
				if err != nil {
					// Check if reply is not avaible due to size too big
					s, ok := status.FromError(err)
					if ok && s.Code() == codes.ResourceExhausted {
						log.Printf("[Front]: Response: #%d: Resource exhausted error occurred: %s.", id, s.Message())
					} else {
						log.Printf("[Front]: %s is unreachable: %v", addr, err)
						master.RemoveWorker(addr, "unreachable")
						// #TODO: assign a new worker
					}
					// Notify failure
					successChan <- false
				} else {
					// Setup actual reply
					response.ID = int32(WorkerResponse.GetID())
					response.Result = WorkerResponse.GetResult()

					// Notify success
					successChan <- true
				}
			}(i, &response)
		}

		for failedWorkes := 0; failedWorkes < len(workers); {
			outcome := <-successChan
			if outcome {
				log.Printf("[Front]: Response: #%d (%d); Time spent: %d ms", id, response.GetID(), time.Since(startTime).Milliseconds())
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
