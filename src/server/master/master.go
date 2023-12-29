package master

import (
	"context"
	"fmt"
	"log"
	pb "mini/proto"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

var (
	WorkerListLock sync.RWMutex
	WorkerList     []string
	workerCounter  int
)

type MasterServer struct {
	pb.UnimplementedMasterServer
}

// SayHello implements helloworld.GreeterServer
func (s *MasterServer) NotifyActiveWorker(ctx context.Context, in *pb.NotifyRequest) (*pb.NotifyReply, error) {
	log.Printf("Notification from %s", in.GetWorkerAddress())

	// Add worker
	resultMessage := AddWorker(in.GetWorkerAddress())

	// send response
	return &pb.NotifyReply{Result: resultMessage}, nil
}

func AddWorker(targetWorkerAddr string) string {
	// check if worker already active
	WorkerListLock.RLock()
	for _, v := range WorkerList {
		if v == targetWorkerAddr {
			return "ALREADY ADDED"
		}
	}
	WorkerListLock.RUnlock()

	// append worker to active list
	WorkerListLock.Lock()
	WorkerList = append(WorkerList, targetWorkerAddr)
	WorkerListLock.Unlock()

	return "OK"
}

func RemoveWorker(targetWorkerAddr string) {
	// delete worker address
	WorkerListLock.Lock()
	for i, v := range WorkerList {
		if v == targetWorkerAddr {
			WorkerList = append(WorkerList[:i], WorkerList[i+1:]...)
		}
	}
	WorkerListLock.Unlock()
}

func GetWorker() string {
	for len(WorkerList) == 0 {
		time.Sleep(time.Second * 1)
	}
	WorkerListLock.RLock()
	// circular workerCounter
	workerCounter = (workerCounter + 1) % len(WorkerList)
	addr := WorkerList[workerCounter]
	WorkerListLock.RUnlock()

	return addr
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
