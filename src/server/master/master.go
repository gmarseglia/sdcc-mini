package master

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	pb "mini/proto"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

var (
	MasterPort        = flag.String("MasterPort", "", "The port for master service.")
	WorkerList        []string
	WorkerListLock    sync.RWMutex
	WorkerCostMap     = make(map[string]float32)
	WorkerCostMapLock sync.RWMutex
	workerChannel     = make(chan int, 1000)
	s                 *grpc.Server
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
func (s *MasterServer) NotifyPing(ctx context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	// log.Printf("[Master]: ping from %s: (QueueSize: %d, UsageCPU: %.2f).", in.GetWorkerAddress(), in.GetQueueSize(), in.GetUsageCPU())

	// Compute cost of the worker
	var cost float32 = in.GetUsageCPU() * (float32(in.QueueSize) + 1)
	WorkerCostMapLock.Lock()
	WorkerCostMap[in.GetWorkerAddress()] = cost
	WorkerCostMapLock.Unlock()
	log.Printf("[Master]: Cost of worker %s: %.2f", in.GetWorkerAddress(), cost)

	// send response
	return &pb.PingReply{Result: "PING OK"}, nil
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
	// Wait for workers to be added
	for len(WorkerCostMap) == 0 {
		time.Sleep(time.Second * 1)
	}

	WorkerCostMapLock.Lock()

	// Find the cheapest worker
	var cheapestWorker string
	var lowestCost float32 = float32(math.Inf(1))
	for workerAddr, cost := range WorkerCostMap {
		log.Printf("[Master]: Cost of worker %s: %.2f", workerAddr, cost)
		if cost < lowestCost {
			lowestCost = cost
			cheapestWorker = workerAddr
		}
	}
	// Increase the cost to avoid overuse on burst traffic
	WorkerCostMap[cheapestWorker] = lowestCost * 1.5

	WorkerCostMapLock.Unlock()

	return cheapestWorker
}

func monitorWorker() {
	for workers := range workerChannel {
		log.Printf("[Master]: Active workers: %d", workers)
	}
}

func StartServer() {
	// listen to request to specified port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *MasterPort))
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
