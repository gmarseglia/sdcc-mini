package main

import (
	"log"
	"mini/worker/back"
	"mini/worker/worker"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	wg = sync.WaitGroup{}
)

func listen() (net.Listener, error) {
	// listen to request to a free port
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Printf("[Main]: Failed to listen.\nMore: %v", err)
		return nil, err
	}

	workerPort := lis.Addr().(*net.TCPAddr).Port
	log.Printf("[Main]: Worker server listening at %d", workerPort)

	return lis, nil
}

func stopComponentsAndExit(message string) {
	log.Printf("[Main]: %s. Begin components stop.", message)
	back.StopServer(&wg)

	wg.Wait()

	exit()
}

func exit() {
	log.Printf("[Main]: All components stopped. Main component stopped. Goodbye.")

	os.Exit(0)
}

func main() {
	log.Printf("[Main]: Welcome. Main component started. Begin components start.")

	// find free port and listen
	workerListener, err := listen()
	if err != nil {
		exit()
	}

	// Activate the Back Server
	wg.Add(1)
	go back.StartServer(workerListener)
	time.Sleep(time.Millisecond * 10)

	// install signal handler
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		stopComponentsAndExit("SIGTERM received")
	}()

	// notify master
	err = worker.NotifyWorkerActive(workerListener.Addr().String())
	if err != nil {
		stopComponentsAndExit("Master unreachable")
	}

	// infite loop with server pings
	for {
		time.Sleep(time.Second * 10)
		err := worker.PingServer()
		if err != nil {
			stopComponentsAndExit("Ping failed")
		}
	}

}
