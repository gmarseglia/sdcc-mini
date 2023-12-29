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
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	masterAddr     = flag.String("addr", "localhost", "The address to connect to")
	masterPort     = flag.Int("masterPort", 55556, "The port of the master service")
	workerListener net.Listener
	active         int
	activeChannel  = make(chan int, 1000)
	counter        int
	counterLock    sync.Mutex
	Prng           = rand.New(rand.NewSource(14))
)

func listen() {
	// listen to request to a free port
	var err error
	workerListener, err = net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	workerPort := workerListener.Addr().(*net.TCPAddr).Port
	log.Printf("Worker server listening at %d", workerPort)

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
	r, err := c.NotifyActiveWorker(ctx, &pb.NotifyRequest{WorkerAddress: workerListener.Addr().String()})
	if err != nil {
		log.Fatalf("could not notify: %v", err)
	}

	log.Printf("r: %s", r.GetResult())
}

type backServer struct {
	pb.UnimplementedBackServer
}

// SayHello implements helloworld.GreeterServer
func (s *backServer) Choice(ctx context.Context, in *pb.ChoiceBiRequest) (*pb.ChoiceReply, error) {
	// signal as activated
	active += 1
	activeChannel <- active

	// log response
	counterLock.Lock()
	counter += 1
	replyID := counter
	counterLock.Unlock()
	log.Printf("Received #%d:(%s, %s)", counter, in.GetOption1(), in.GetOption2())

	// sleep
	utils.SimulatedCPUIntensiveFunction(1000, &active, 1)

	// randomly choose response
	var response string
	if Prng.Intn(2) == 0 {
		response = in.GetOption1()
	} else {
		response = in.GetOption2()
	}

	// signal as deactivated
	active -= 1
	activeChannel <- active

	// send response
	return &pb.ChoiceReply{Option: response, ReplyID: int32(replyID)}, nil
}

func debugActive() {
	for active := range activeChannel {
		log.Printf("*active: %d", active)
	}
}

func activateBackServer() {
	// create a new server
	s := grpc.NewServer()

	// register the server
	pb.RegisterBackServer(s, &backServer{})
	log.Printf("Back server listening at %v", workerListener.Addr())

	// start debugging active level
	go debugActive()

	// serve the request
	if err := s.Serve(workerListener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	// find free port and listen
	listen()

	// Activate the Back Server
	go activateBackServer()
	time.Sleep(time.Millisecond * 10)

	// notify master
	notifyWorkerActive()
	select {}

}
