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
	port    = flag.Int("port", 55555, "The server port")
	seed    = flag.Int("seed", 14, "The seed for the PRNG")
	counter = 0
	active  = 0
)

// build the PRNG
var (
	Prng = rand.New(rand.NewSource(int64(*seed)))
)

// server is used to implement the MiniServer interface
type server struct {
	pb.UnimplementedMiniServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) Choice(ctx context.Context, in *pb.ChoiceBiRequest) (*pb.ChoiceReply, error) {
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
	log.Printf("*active: %d", active)

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

func main() {
	// parse the flags for CLI
	flag.Parse()

	// listen to request to specified port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a new server
	s := grpc.NewServer()

	// register the server
	pb.RegisterMiniServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	// start debugging active level
	go debugActive()

	// serve the request
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
