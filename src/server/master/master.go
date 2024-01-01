package master

import (
	"context"
	"flag"
	"fmt"
	"log"
	pb "mini/proto"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

var (
	masterPort     = flag.Int("masterPort", 55556, "The server port for master service")
	WorkerListLock sync.RWMutex
	WorkerList     []string
	workerCounter  int
	workerChannel  = make(chan int, 1000)
	s              *grpc.Server
)

type MasterServer struct {
	pb.UnimplementedMasterServer
}

// SayHello implementation
func (s *MasterServer) NotifyActiveWorker(ctx context.Context, in *pb.NotifyRequest) (*pb.NotifyReply, error) {
	log.Printf("[Master]: Notification from %s", in.GetWorkerAddress())

	// Add worker
	resultMessage := AddWorker(in.GetWorkerAddress())

	// send response
	return &pb.NotifyReply{Result: resultMessage}, nil
}

// NotifyPing implementation
func (s *MasterServer) NotifyPing(ctx context.Context, in *pb.NotifyRequest) (*pb.NotifyReply, error) {
	// log.Printf("[Master DEV]: ping from %s.", in.GetWorkerAddress())

	// send response
	return &pb.NotifyReply{Result: "PING OK"}, nil
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
	workerChannel <- len(WorkerList)
	WorkerListLock.Unlock()

	return "OK"
}

func RemoveWorker(targetWorkerAddr string) {
	// delete worker address
	WorkerListLock.Lock()
	for i, v := range WorkerList {
		if v == targetWorkerAddr {
			WorkerList = append(WorkerList[:i], WorkerList[i+1:]...)
			workerChannel <- len(WorkerList)
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

func monitorWorker() {
	for workers := range workerChannel {
		log.Printf("[Master]: Active workers: %d", workers)
	}
}

func StartServer() {
	// listen to request to specified port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *masterPort))
	if err != nil {
		log.Fatalf("[Master]: Failed to listen: %v", err)
	}

	// create a new server
	s = grpc.NewServer()

	// register the server
	pb.RegisterMasterServer(s, &MasterServer{})
	log.Printf("[Master]: Listening at %v", lis.Addr())

	workerChannel <- 0
	go monitorWorker()

	// serve the request
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[Master]: Failed to serve: %v", err)
	}
}

func StopServer(wg *sync.WaitGroup) {
	log.Printf("[Master]: Grafecully stopping...")

	// Graceful stop
	s.GracefulStop()
	log.Printf("[Master]: Done.")

	// Comunicate on channel so sync
	(*wg).Done()
}
