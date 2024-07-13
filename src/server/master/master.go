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
	WorkerInfoMap     = make(map[string]WorkerInfo)
	WorkerInfoMapLock sync.RWMutex
	workerChannel     = make(chan int, 1000)
	s                 *grpc.Server
)

type MasterServer struct {
	pb.UnimplementedMasterServer
}

type WorkerInfo struct {
	QueueSize int32
	UsageCPU  float32
	Cost      float32
	LastPing  time.Time
}

// SayHello implementation
func (s *MasterServer) NotifyActiveWorker(ctx context.Context, in *pb.NotifyRequest) (*pb.NotifyReply, error) {
	log.Printf("[Master]: Notification from %s", in.GetWorkerAddress())

	// Add worker
	_, exists := WorkerInfoMap[in.GetWorkerAddress()]
	if exists {
		return &pb.NotifyReply{Result: "ALREADY ADDED"}, nil
	} else {
		WorkerInfoMapLock.Lock()
		WorkerInfoMap[in.GetWorkerAddress()] = WorkerInfo{
			Cost:     float32(math.Inf(1)),
			LastPing: time.Now(),
		}
		workerChannel <- len(WorkerInfoMap)
		WorkerInfoMapLock.Unlock()
		log.Printf("[Master]: Added worker %s", in.GetWorkerAddress())
		return &pb.NotifyReply{Result: "NOTIFICATION OK"}, nil
	}
}

// NotifyPing implementation
func (s *MasterServer) NotifyPing(ctx context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	// Update worker info
	cost := in.GetUsageCPU() * (float32(in.QueueSize) + 1)
	WorkerInfoMapLock.Lock()
	WorkerInfoMap[in.GetWorkerAddress()] = WorkerInfo{
		QueueSize: in.GetQueueSize(),
		UsageCPU:  in.GetUsageCPU(),
		Cost:      cost,
		LastPing:  time.Now(),
	}
	WorkerInfoMapLock.Unlock()

	// log.Printf("[Master]: Cost of worker %s: %.2f", in.GetWorkerAddress(), cost)

	// send response
	return &pb.PingReply{Result: "PING OK"}, nil
}

func RemoveWorker(targetWorkerAddr string, cause string) {
	// delete worker address
	WorkerInfoMapLock.Lock()
	_, exists := WorkerInfoMap[targetWorkerAddr]
	if exists {
		delete(WorkerInfoMap, targetWorkerAddr)
		log.Printf("[Master]: Deleted worker %s due to %s", targetWorkerAddr, cause)
	}
	WorkerInfoMapLock.Unlock()
}

func GetWorker() string {
	// Wait for workers to be added
	for len(WorkerInfoMap) == 0 {
		time.Sleep(time.Second * 1)
	}

	WorkerInfoMapLock.Lock()

	// Find the cheapest worker
	var cheapestWorker string
	var cheapestWorkerInfo WorkerInfo
	var lowestCost float32 = float32(math.Inf(1))
	for workerAddr, info := range WorkerInfoMap {
		log.Printf("[Master]: Cost of worker %s: %.2f", workerAddr, info.Cost)
		if info.Cost < lowestCost {
			lowestCost = info.Cost
			cheapestWorker = workerAddr
			cheapestWorkerInfo = info
		}
	}
	// Increase the cost to avoid overuse on burst traffic
	cheapestWorkerInfo.Cost = float32(cheapestWorkerInfo.QueueSize+1) * (cheapestWorkerInfo.UsageCPU * 1.1)
	WorkerInfoMap[cheapestWorker] = cheapestWorkerInfo

	WorkerInfoMapLock.Unlock()

	return cheapestWorker
}

func monitorWorker() {
	go func() {
		for workers := range workerChannel {
			log.Printf("[Master]: Active workers: %d", workers)
		}
	}()

	go func() {
		for {
			// #TODO: set these as env variables
			time.Sleep(time.Second * 3)
			WorkerInfoMapLock.RLock()
			for workerAddr, info := range WorkerInfoMap {
				if time.Since(info.LastPing) > time.Second*3 {
					WorkerInfoMapLock.RUnlock()
					RemoveWorker(workerAddr, "ping timeout")
					WorkerInfoMapLock.RLock()
				}
			}
			WorkerInfoMapLock.RUnlock()
		}
	}()
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
