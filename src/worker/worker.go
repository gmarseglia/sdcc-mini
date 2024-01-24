package main

import (
	"flag"
	"fmt"
	"log"
	"mini/utils"
	"mini/worker/back"
	worker "mini/worker/workercomponent"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	wg = sync.WaitGroup{}
)

func setupAddresses() {
	// Setup WorkerPort
	utils.SetupField(back.HostAddr, "HostAddr", utils.GetOutboundIP().String())
	utils.SetupField(back.HostPort, "HostPort", "55557")
	back.HostFullAddr = fmt.Sprintf("%s:%s", *back.HostAddr, *back.HostPort)
	utils.SetupField(worker.MasterAddr, "MasterAddr", "0.0.0.0")
	utils.SetupField(worker.MasterPort, "MasterPort", "55556")
	worker.MasterFullAddr = fmt.Sprintf("%s:%s", *worker.MasterAddr, *worker.MasterPort)
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
	flag.Parse()

	log.Printf("[Main]: Welcome. Main component started. Begin components start.")

	setupAddresses()

	// Activate the Back Server
	wg.Add(1)
	go func() {
		err := back.StartServer()
		if err != nil {
			exit()
		}
	}()
	time.Sleep(time.Millisecond * 10)

	// install signal handler
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		stopComponentsAndExit("SIGTERM received")
	}()

	// notify master
	err := worker.NotifyWorkerActive()
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
