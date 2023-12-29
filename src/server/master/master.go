package master

import (
	"context"
	"fmt"
	"log"
	pb "mini/proto"
	"net"
	"sync"

	"google.golang.org/grpc"
)

var (
	WorkerListLock sync.RWMutex
	WorkerList     []string
)

type MasterServer struct {
	pb.UnimplementedMasterServer
}

// SayHello implements helloworld.GreeterServer
func (s *MasterServer) NotifyActiveWorker(ctx context.Context, in *pb.NotifyRequest) (*pb.NotifyReply, error) {
	log.Printf("Notification from %s", in.GetWorkerAddress())

	// lock the WorkerList
	// check if worker already active
	WorkerListLock.RLock()
	for _, v := range WorkerList {
		if v == in.GetWorkerAddress() {
			return &pb.NotifyReply{Result: "ALREADY ADDED"}, nil
		}
	}
	WorkerListLock.RUnlock()

	// append worker to active list
	WorkerListLock.Lock()
	WorkerList = append(WorkerList, in.GetWorkerAddress())
	WorkerListLock.Unlock()

	// send response
	return &pb.NotifyReply{Result: "OK"}, nil
}

func ActivateMasterServer(port *int) {
	// listen to request to specified port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a new server
	s := grpc.NewServer()

	// register the server
	pb.RegisterMasterServer(s, &MasterServer{})
	log.Printf("Master server listening at %v", lis.Addr())

	// serve the request
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
