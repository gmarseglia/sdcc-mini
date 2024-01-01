package main

import (
	"flag"
	"log"
	"mini/server/front"
	"mini/server/master"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// define flags
var (
	frontPort  = flag.Int("frontPort", 55555, "The server port for front service")
	masterPort = flag.Int("masterPort", 55556, "The server port for master service")
	wg         = sync.WaitGroup{}
)

func main() {
	log.SetOutput(os.Stdout)

	log.Printf("[Main]: Welcome, begin components start.")

	// parse the flags for CLI
	flag.Parse()

	// start master server and add to wait gruop
	wg.Add(1)
	go master.StartServer(masterPort)
	time.Sleep(time.Millisecond * 10)

	// start front server and add to wait group
	wg.Add(1)
	go front.StartFrontServer(frontPort)

	// install signal handler
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Printf("[Main]: SIGTERM received, begin components stop.")
		master.StopServer(&wg)
		front.StopServer(&wg)
	}()

	wg.Wait()
	log.Printf("[Main]: All componentes are stopped. Stopping server...")
	log.Printf("[Main]: Done")

}
