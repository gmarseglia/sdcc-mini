package main

import (
	"log"
	"mini/worker/back"
	"mini/worker/worker"
	"net"
	"time"
)

func listen() net.Listener {
	// listen to request to a free port
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	workerPort := lis.Addr().(*net.TCPAddr).Port
	log.Printf("Worker server listening at %d", workerPort)

	return lis
}

func main() {
	// find free port and listen
	workerListener := listen()

	// Activate the Back Server
	go back.StartServer(workerListener)
	time.Sleep(time.Millisecond * 10)

	// notify master
	worker.NotifyWorkerActive(workerListener.Addr().String())
	select {}

}
