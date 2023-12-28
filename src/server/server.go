package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	pb "mini/proto"
	"mini/utils"
	"net"
	"time"

	"google.golang.org/grpc"
)

// define flags
var (
	frontPort  = flag.Int("frontPort", 55555, "The server port for front service")
	masterPort = flag.Int("masterPort", 55556, "The server port for master service")
	seed       = flag.Int("seed", 14, "The seed for the PRNG")
	counter    = 0
	active     = 0
)

// build the PRNG
var (
	Prng = rand.New(rand.NewSource(int64(*seed)))
)

// FrontServer is used to implement the MiniServer interface
type FrontServer struct {
	pb.UnimplementedFrontServer
}

// SayHello implements helloworld.GreeterServer
func (s *FrontServer) Choice(ctx context.Context, in *pb.ChoiceBiRequest) (*pb.ChoiceReply, error) {
	// signal as activated
	active += 1

	// log response
	counter += 1
	replyID := counter
	log.Printf("Received #%d:(%s, %s)", counter, in.GetOption1(), in.GetOption2())

	// sleep
	utils.SimulatedCPUIntensiveFunction(1000, &active, 1)

	// signal as deactivated
	active -= 1

	// randomly choose response
	var response string
	if Prng.Intn(2) == 0 {
		response = in.GetOption1()
	} else {
		response = in.GetOption2()
	}

	// send response
	return &pb.ChoiceReply{Option: response, ReplyID: int32(replyID)}, nil
}

type masterServer struct {
	pb.UnimplementedMasterServer
}

// SayHello implements helloworld.GreeterServer
func (s *masterServer) NotifyActiveWorker(ctx context.Context, in *pb.NotifyRequest) (*pb.NotifyReply, error) {
	log.Printf("Notification from %s", in.GetWorkerAddress())

	// send response
	return &pb.NotifyReply{Result: "OK"}, nil
}

func debugActive() {
	lastActive := -1
	for {
		if active != lastActive {
			lastActive = active
			log.Printf("*active: %d", active)
			time.Sleep(utils.Step * 2)
		}
	}
}

func activateFrontServer(port *int) {
	// listen to request to specified port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a new server
	s := grpc.NewServer()

	// register the server
	pb.RegisterFrontServer(s, &FrontServer{})
	log.Printf("Front server listening at %v", lis.Addr())

	// start debugging active level
	go debugActive()

	// serve the request
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func activateMasterServer(port *int) {
	// listen to request to specified port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a new server
	s := grpc.NewServer()

	// register the server
	pb.RegisterMasterServer(s, &masterServer{})
	log.Printf("Master server listening at %v", lis.Addr())

	// serve the request
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	// parse the flags for CLI
	flag.Parse()

	go activateMasterServer(masterPort)
	time.Sleep(time.Millisecond * 10)

	activateFrontServer(frontPort)

}
