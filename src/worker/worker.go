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

	"github.com/shirou/gopsutil/cpu"
)

var (
	PingTimeout  = flag.Int("PingTimeout", 1, "The ping timeout in seconds.")
	FakeCPUScale = flag.Float64("FakeCPUScale", 1, "The scale of the fake CPU usage.")
	wg           = sync.WaitGroup{}
)

func setupFields() {
	// Setup WorkerPort
	utils.SetupFieldOptional(worker.HostAddr, "HostAddr", utils.GetOutboundIP().String())
	utils.SetupFieldOptional(worker.HostPort, "HostPort", "55557")
	worker.HostFullAddr = fmt.Sprintf("%s:%s", *worker.HostAddr, *worker.HostPort)
	utils.SetupFieldMandatory(worker.MasterAddr, "MasterAddr", func() {
		log.Printf("[Main]: MasterAddr is a mandatory field.")
		exit()
	})
	utils.SetupFieldOptional(worker.MasterPort, "MasterPort", "55556")
	worker.MasterFullAddr = fmt.Sprintf("%s:%s", *worker.MasterAddr, *worker.MasterPort)
	back.HostPort = *worker.HostPort
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

	setupFields()

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
		// Prepare the request
		percentages, _ := cpu.Percent(time.Second, false)
		err := worker.PingServer(*PingTimeout*10, float32(percentages[0]*(*FakeCPUScale)))
		if err != nil {
			stopComponentsAndExit("Ping failed")
		}
		time.Sleep(time.Second * time.Duration(*PingTimeout-1))
	}

}
