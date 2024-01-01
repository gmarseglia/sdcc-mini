package back

import (
	"context"
	"log"
	"math/rand"
	pb "mini/proto"
	"mini/utils"
	"net"
	"sync"

	"google.golang.org/grpc"
)

var (
	active        int
	activeChannel = make(chan int, 1000)
	counter       int
	counterLock   sync.Mutex
	Prng          = rand.New(rand.NewSource(14))
	s             *grpc.Server
)

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

func StartServer(workerListener net.Listener) {
	// create a new server
	s = grpc.NewServer()

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

func StopServer(wg *sync.WaitGroup) {
	log.Printf("[Back server]: Grafecully stopping...")

	// Graceful stop
	s.GracefulStop()
	log.Printf("[Back server]: Done.")

	// Comunicate on channel so sync
	(*wg).Done()
}
