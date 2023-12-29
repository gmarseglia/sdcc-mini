package main

import (
	"flag"
	"mini/server/front"
	"mini/server/master"
	"time"
)

// define flags
var (
	frontPort  = flag.Int("frontPort", 55555, "The server port for front service")
	masterPort = flag.Int("masterPort", 55556, "The server port for master service")
)

func main() {
	// parse the flags for CLI
	flag.Parse()

	go master.ActivateMasterServer(masterPort)
	time.Sleep(time.Millisecond * 10)

	front.ActivateFrontServer(frontPort)

}
