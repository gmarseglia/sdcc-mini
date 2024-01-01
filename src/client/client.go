package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "mini/proto"
)

const (
	timeout = 60 * time.Second
)

var (
	frontAddr = flag.String("frontAddr", "localhost", "The address to connect to")
	frontPort = flag.Int("frontPort", 55555, "The port of the master service")
)

func main() {
	// parse the flags
	flag.Parse()

	// Set up a connection to the gRPC server
	serverFullAddr := fmt.Sprintf("%s:%d", *frontAddr, *frontPort)
	conn, err := grpc.Dial(serverFullAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// create the client object
	c := pb.NewFrontClient(conn)

	// create the context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// time the call
	startTime := time.Now()

	// contact the server
	r, err := c.Choice(ctx, &pb.ChoiceBiRequest{Option1: "A", Option2: "B"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// print the result
	log.Printf("Response: (#%d, %s); Time spent: %d ms", r.GetReplyID(), r.GetOption(), time.Since(startTime).Milliseconds())
}
