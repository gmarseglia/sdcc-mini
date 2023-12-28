package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	pb "mini/proto"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	masterAddr        = flag.String("addr", "localhost", "The address to connect to")
	masterPort        = flag.Int("masterPort", 55556, "The port of the master service")
	workerInitialPort = flag.Int("workerPort", 55557, "The port of the worker service")
	workerPort        int
	workerListener    net.Listener
)

func listen() {
	// listen to request to specified port
	var err error
	for port := *workerInitialPort; ; port += 1 {
		workerListener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			workerPort = port
			log.Printf("Worker server listening at %d", workerPort)
			break
		} else {
			log.Fatalf("failed to listen: %v", err)
		}
	}
}

func notifyWorkerActive() {
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
	r, err := c.NotifyActiveWorker(ctx, &pb.NotifyRequest{WorkerAddress: fmt.Sprintf("localhost:%d", workerPort)})
	if err != nil {
		log.Fatalf("could not notify: %v", err)
	}

	log.Printf("r: %s", r.GetResult())
}

func main() {
	listen()
	notifyWorkerActive()
}
