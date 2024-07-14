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
	WorkerAddress string
	QueueSize     int32
	UsageCPU      float32
	Cost          float32
	LastPing      time.Time
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
			WorkerAddress: in.GetWorkerAddress(),
			Cost:          float32(math.Inf(1)),
			LastPing:      time.Now(),
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
	winfo := WorkerInfoMap[in.GetWorkerAddress()]
	winfo.QueueSize = in.GetQueueSize()
	winfo.UsageCPU = in.GetUsageCPU()
	winfo.Cost = cost
	winfo.LastPing = time.Now()
	WorkerInfoMap[in.GetWorkerAddress()] = winfo
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
		workerChannel <- len(WorkerInfoMap)
		log.Printf("[Master]: Deleted worker %s due to %s", targetWorkerAddr, cause)
	}
	WorkerInfoMapLock.Unlock()
}

func GetWorkers(number int) []WorkerInfo {
	// Wait for workers to be added
	for len(WorkerInfoMap) < number {
		time.Sleep(time.Second * 1)
	}

	var cheapestWorkerInfo = make([]WorkerInfo, number)
	for i := 0; i < number; i++ {
		cheapestWorkerInfo[i] = WorkerInfo{Cost: float32(math.Inf(1))}
	}

	// Find the cheapest worker
	WorkerInfoMapLock.Lock()

	for _, info := range WorkerInfoMap {
		// log.Printf("[Master]: Cost of worker %s: %.2f", info.WorkerAddress, info.Cost)
		for i := 0; i < number; i++ {
			if info.Cost <= cheapestWorkerInfo[i].Cost {
				// Translate the array
				copy(cheapestWorkerInfo[i+1:], cheapestWorkerInfo[i:number-1])
				cheapestWorkerInfo[i] = info
				break
			}
		}
	}

	// Increase the cost to avoid overuse on burst traffic
	for i := 0; i < number; i++ {
		log.Printf("[Master]: Cost of worker %s: %.2f", cheapestWorkerInfo[i].WorkerAddress, cheapestWorkerInfo[i].Cost)
		winfo := cheapestWorkerInfo[i]
		winfo.Cost = float32(winfo.QueueSize+1) * (winfo.UsageCPU * 1.5)
		WorkerInfoMap[cheapestWorkerInfo[i].WorkerAddress] = winfo
		cheapestWorkerInfo[i] = winfo
	}

	WorkerInfoMapLock.Unlock()

	return cheapestWorkerInfo
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
