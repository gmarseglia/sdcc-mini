package worker

import (
	"context"
	"flag"
	"fmt"
	"log"
	pb "mini/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	masterAddr = flag.String("addr", "localhost", "The address to connect to")
	masterPort = flag.Int("masterPort", 55556, "The port of the master service")
)

func NotifyWorkerActive(workerAddr string) {
	// Set up a connection to the gRPC server
	masterFullAddr := fmt.Sprintf("%s:%d", *masterAddr, *masterPort)
	conn, err := grpc.Dial(masterFullAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// create the client object
	c := pb.NewMasterClient(conn)

	// create the context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	// contact the server
	r, err := c.NotifyActiveWorker(ctx, &pb.NotifyRequest{WorkerAddress: workerAddr})
	if err != nil {
		log.Fatalf("could not notify: %v", err)
	}

	log.Printf("r: %s", r.GetResult())
}
