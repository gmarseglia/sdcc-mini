package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "mini/proto"
	"mini/utils"
)

const (
	timeout = 60 * time.Second
)

var (
	FrontAddr    = flag.String("FrontAddr", "", "The address to connect to")
	FrontPort    = flag.String("FrontPort", "", "The port of the master service")
	requestCount = flag.Int("requestCount", 1, "The number of requests to send")
	counter      int
	counterLock  sync.Mutex
	wg           sync.WaitGroup
	c            pb.FrontClient
)

func setupFields() {
	utils.SetupFieldMandatory(FrontAddr, "FrontAddr", func() {
		log.Printf("[Main]: FrontAddr field is mandatory.")
		exit()
	})
	utils.SetupFieldOptional(FrontPort, "FrontPort", "55555")
}

func exit() {
	log.Printf("[Main]: All components stopped. Main component stopped. Goodbye.")
	os.Exit(0)
}

func choice() {
	// Internal ID
	counterLock.Lock()
	counter++
	id := counter
	counterLock.Unlock()

	// create the context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// time the call
	startTime := time.Now()

	log.Printf("[Client]: Request #%d sent.", id)

	// contact the server
	r, err := c.Choice(ctx, &pb.ChoiceBiRequest{Option1: "A", Option2: "B"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// print the result
	log.Printf("[Client]: Request #%d -> Response: (#%d , %s) in %d ms", id, r.GetReplyID(), r.GetOption(), time.Since(startTime).Milliseconds())

	wg.Done()
}

func main() {
	log.SetOutput(os.Stdout)

	// parse the flags
	flag.Parse()
	setupFields()

	// Welcome message
	log.Printf("[Main]: Welcome. Client will send %d requests in parallel.", *requestCount)

	// Set up a connection to the gRPC server
	serverFullAddr := fmt.Sprintf("%s:%s", *FrontAddr, *FrontPort)
	conn, err := grpc.Dial(serverFullAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("[Main]: Could not not connect. More:\n%v", err)
	}
	defer conn.Close()

	// create the client object
	c = pb.NewFrontClient(conn)

	for i := 0; i < *requestCount; i++ {
		wg.Add(1)
		go choice()
	}

	// wait
	log.Printf("[Main]: All requests sent. Waiting for responses...")
	wg.Wait()

	log.Printf("[Main]: All requests completed. Terminating. Goodbye.")
}
