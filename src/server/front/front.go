package front

import (
	"context"
	"flag"
	"fmt"
	"log"
	pb "mini/proto"
	"mini/server/master"
	"mini/utils"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	FrontPort   = flag.String("FrontPort", "", "The port for front service.")
	counter     int
	counterLock sync.Mutex
	s           *grpc.Server
)

// FrontServer is used to implement the MiniServer interface
type FrontServer struct {
	pb.UnimplementedFrontServer
}

// SayHello implements helloworld.GreeterServer
func (s *FrontServer) Choice(ctx context.Context, in *pb.ChoiceBiRequest) (*pb.ChoiceReply, error) {

	// generate ID
	counterLock.Lock()
	counter += 1
	var id = counter
	counterLock.Unlock()

	log.Printf("[Front]: Request #%d started.", id)

	for {
		// get worker
		addr := master.GetWorker()

		// Set up a connection to the gRPC server
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("[Front]: Did not connect: %v", err)
		}
		defer conn.Close()

		// create the client object
		c := pb.NewBackClient(conn)

		// create the context
		ctxInternal, cancel := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel()

		// time the call
		startTime := time.Now()

		// contact the server
		WorkerResponse, err := c.Choice(ctxInternal, in)
		if err != nil {
			// log.Printf("could not greet: %v", err)
			log.Printf("[Front]: %s is unreachable.", addr)
			master.RemoveWorker(addr)
			continue
		}

		log.Printf("[Front]: Response: #%d (%d); Time spent: %d ms", id, WorkerResponse.GetReplyID(), time.Since(startTime).Milliseconds())

		// send response
		return &pb.ChoiceReply{Option: WorkerResponse.GetOption(), ReplyID: int32(id)}, nil
	}
}

func (s *FrontServer) Convulution(ctx context.Context, in *pb.ConvolutionFrontRequest) (*pb.ConvolutionFrontReply, error) {
	// generate ID
	counterLock.Lock()
	counter += 1
	var id = counter
	counterLock.Unlock()

	log.Printf("[Front]: Request #%d started.", id)

	// Deserialize
	target := utils.ProtoToMatrix(in.Target)
	kernel := utils.ProtoToMatrix(in.Kernel)

	// #TODO: remove if working
	utils.PrettyPrint("Target Matrix", target)
	utils.PrettyPrint("Kernel Matrix", kernel)

	for {
		// get worker
		addr := master.GetWorker()

		// Set up a connection to the gRPC server
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("[Front]: Did not connect: %v", err)
		}
		defer conn.Close()

		// create the client object
		c := pb.NewBackClient(conn)

		// create the context
		ctxInternal, cancel := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel()

		// time the call
		startTime := time.Now()

		// Create the back request
		backRequest := pb.ConvolutionBackRequest{Target: in.GetTarget(), Kernel: in.GetKernel(), ID: int32(id)}

		// contact the server
		WorkerResponse, err := c.Convulution(ctxInternal, &backRequest)
		if err != nil {
			// log.Printf("could not greet: %v", err)
			log.Printf("[Front]: %s is unreachable.", addr)
			master.RemoveWorker(addr)
			continue
		}

		log.Printf("[Front]: Response: #%d (%d); Time spent: %d ms", id, WorkerResponse.GetID(), time.Since(startTime).Milliseconds())

		// #TODO: remove if working
		result := utils.ProtoToMatrix(WorkerResponse.GetResult())
		utils.PrettyPrint("Result Matrix", result)

		// send response
		return &pb.ConvolutionFrontReply{Result: WorkerResponse.GetResult(), ID: int32(id)}, nil
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
