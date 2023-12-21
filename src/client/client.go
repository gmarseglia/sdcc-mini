package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "mini/proto"
)

const (
	timeout = 5 * time.Second
)

var (
	addr = flag.String("addr", "localhost:55555", "the address to connect to")
)

func main() {
	// parse the flags
	flag.Parse()

	// Set up a connection to the gRPC server
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// create the client object
	c := pb.NewMiniClient(conn)

	// create the context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// contact the server
	r, err := c.Choice(ctx, &pb.ChoiceBiRequest{Option1: "uno", Option2: "due", Millis: 1000})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// print the result
	log.Printf("Greeting: %s", r.GetOption())
}
