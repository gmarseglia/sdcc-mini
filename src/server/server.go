package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	pb "mini/proto"
	"net"
	"time"

	"google.golang.org/grpc"
)

// define flags
var (
	port = flag.Int("port", 55555, "The server port")
	seed = flag.Int("seed", 14, "The seed for the PRNG")
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
	// log response
	log.Printf("Received: options:(%s, %s), sleep: %d", in.GetOption1(), in.GetOption2(), in.GetMillis())

	// sleep
	time.Sleep(time.Duration(in.GetMillis()) * time.Millisecond)

	// randomly choose response
	var response string
	if Prng.Intn(2) == 0 {
		response = in.GetOption1()
	} else {
		response = in.GetOption2()
	}

	// send response
	return &pb.ChoiceReply{Option: response}, nil
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

	// serve the request
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
