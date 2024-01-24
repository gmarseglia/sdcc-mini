package main

import (
	"flag"
	"log"
	"mini/server/front"
	"mini/server/master"
	"mini/utils"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// define flags
var (
	wg = sync.WaitGroup{}
)

func main() {
	log.SetOutput(os.Stdout)

	log.Printf("[Main]: Welcome. Main component started. Begin components start.")

	// parse the flags for CLI
	flag.Parse()

	// Save workerAddr
	serverAddr := utils.GetOutboundIP().String()
	log.Printf("!!![Main]: Server Address: %s\n", serverAddr)

	// start master server and add to wait gruop
	wg.Add(1)
	go master.StartServer()
	time.Sleep(time.Millisecond * 10)

	// start front server and add to wait group
	wg.Add(1)
	go front.StartFrontServer()

	// install signal handler
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Printf("[Main]: SIGTERM received. Begin components stop.")
		master.StopServer(&wg)
		front.StopServer(&wg)
	}()

	wg.Wait()
	log.Printf("[Main]: All componentes stopped. Main component stopped. Goodbye.")

}
