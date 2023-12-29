package front

import (
	"context"
	"fmt"
	"log"
	pb "mini/proto"
	"mini/server/master"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	counter       int
	workerCounter int
	counterLock   sync.Mutex
)

// FrontServer is used to implement the MiniServer interface
type FrontServer struct {
	pb.UnimplementedFrontServer
}

// SayHello implements helloworld.GreeterServer
func (s *FrontServer) Choice(ctx context.Context, in *pb.ChoiceBiRequest) (*pb.ChoiceReply, error) {

	// generate ID
	counterLock.Lock()
	counter += 1
	var id = counter
	counterLock.Unlock()

	log.Printf("Request #%d started.", id)

	for {
		var addr string

		// circular workerCounter
		master.WorkerListLock.RLock()
		if len(master.WorkerList) == 0 {
			master.WorkerListLock.RUnlock()
			time.Sleep(time.Second * 1)
			continue
		} else {
			workerCounter = (workerCounter + 1) % len(master.WorkerList)
			addr = master.WorkerList[workerCounter]
			master.WorkerListLock.RUnlock()
		}

		// Set up a connection to the gRPC server
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		// create the client object
		c := pb.NewBackClient(conn)

		// create the context
		ctxInternal, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()

		// time the call
		startTime := time.Now()

		// contact the server
		WorkerResponse, err := c.Choice(ctxInternal, in)
		if err != nil {
			log.Printf("could not greet: %v", err)
			// delete worker address
			master.WorkerListLock.Lock()
			for i, v := range master.WorkerList {
				if v == addr {
					master.WorkerList = append(master.WorkerList[:i], master.WorkerList[i+1:]...)
				}
			}
			master.WorkerListLock.Unlock()
			continue
		}
		log.Printf("Response: (#%d/%d,%s); Time spent: %d ms", WorkerResponse.GetReplyID(), id, WorkerResponse.GetOption(), time.Since(startTime).Milliseconds())

		// send response
		return &pb.ChoiceReply{Option: WorkerResponse.GetOption(), ReplyID: int32(id)}, nil
	}
}

func ActivateFrontServer(port *int) {
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

	// serve the request
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
