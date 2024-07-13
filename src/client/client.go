package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "mini/proto"
	"mini/utils"
)

const (
	timeout = 60 * time.Second
)

var (
	FrontAddr    = flag.String("FrontAddr", "", "The address to connect to.")
	FrontPort    = flag.String("FrontPort", "", "The port of the master service.")
	RequestCount = flag.String("RequestCount", "", "The number of requests to send.")
	counter      int
	counterLock  sync.Mutex
	wg           sync.WaitGroup
	c            pb.FrontClient
)

func setupFields() {
	utils.SetupFieldMandatory(FrontAddr, "FrontAddr", func() {
		log.Printf("[Main]: FrontAddr field is mandatory.")
		exit()
	})
	utils.SetupFieldOptional(FrontPort, "FrontPort", "55555")
	utils.SetupFieldOptional(RequestCount, "RequestCount", "1")
}

func exit() {
	log.Printf("[Main]: All components stopped. Main component stopped. Goodbye.")
	os.Exit(0)
}

func choice() {
	// Internal ID
	counterLock.Lock()
	counter++
	id := counter
	counterLock.Unlock()

	// create the context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// time the call
	startTime := time.Now()

	log.Printf("[Client]: Request #%d sent.", id)

	// contact the server
	r, err := c.Choice(ctx, &pb.ChoiceBiRequest{Option1: "A", Option2: "B"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// print the result
	log.Printf("[Client]: Request #%d -> Response: (#%d , %s) in %d ms", id, r.GetReplyID(), r.GetOption(), time.Since(startTime).Milliseconds())

	wg.Done()
}

func convolutionRun() {
	// Internal ID
	counterLock.Lock()
	counter++
	id := counter
	counterLock.Unlock()

	// Settings
	targetSize := 227
	kernelNum := 96
	kernelSize := 11

	// Serialie the matrices
	frontRequest := &pb.ConvolutionFrontRequest{}
	frontRequest.Target = utils.MatrixToProto(utils.GenerateRandomMatrix(targetSize, targetSize))
	for i := 0; i < kernelNum; i++ {
		frontRequest.Kernel = append(frontRequest.Kernel, utils.MatrixToProto(utils.GenerateRandomMatrix(kernelSize, kernelSize)))
	}
	frontRequest.StrideHorizontal = 1
	frontRequest.StrideVertical = 1

	// create the context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// time the call
	startTime := time.Now()

	log.Printf("[Client]: Request #%d sent.", id)

	// contact the server
	r, err := c.Convulution(ctx, frontRequest)

	if err != nil {
		if s, ok := status.FromError(err); ok {
			log.Printf("[Client]: Request #%d -> Unsuccessful! %s: %v", id, s.Message(), s.Details())
			wg.Done()
			return
		}
	}

	// print the result
	log.Printf("[Client]: Request #%d -> Response: (#%d) in %d ms",
		id,
		r.GetID(),
		time.Since(startTime).Milliseconds())

	// for _, resultProto := range r.GetResult() {
	// 	result := utils.ProtoToMatrix(resultProto)
	// 	utils.PrettyPrint("result", result)
	// }

	wg.Done()
}

func maxPoolRun(powerFactor int) {

	targetSize := int(math.Pow(2, float64(powerFactor)))
	var maxPoolSize int
	if targetSize < 2 {
		maxPoolSize = 1
	} else {
		maxPoolSize = 2
	}

	target := utils.GenerateRandomMatrix(targetSize, targetSize)

	// Serialie the matrices
	targetProto := utils.MatrixToProto(target)

	// Internal ID
	counterLock.Lock()
	counter++
	id := counter
	counterLock.Unlock()

	// create the context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// time the call
	startTime := time.Now()

	log.Printf("[Client]: Request #%d sent.", id)

	// contact the server
	r, err := c.MaxPool(ctx, &pb.MaxPoolFrontRequest{
		Target: targetProto,
		Size:   int32(maxPoolSize)})

	if err != nil {
		if s, ok := status.FromError(err); ok {
			log.Printf("[Client]: Request #%d -> Unsuccessful! %s: %v", id, s.Message(), s.Details())
			wg.Done()
			return
		}
	}

	// print the result
	result := utils.ProtoToMatrix(r.GetResult())
	log.Printf("[Client]: Request #%d -> Response: (#%d, powFactor: %d, size: %dx%d) in %d ms",
		id,
		r.GetID(),
		powerFactor,
		len(result), len(result[0]),
		time.Since(startTime).Milliseconds())

	utils.PrettyPrint("Target Matrix", target)
	utils.PrettyPrint("Result Matrix", result)

	wg.Done()
}

func main() {
	log.SetOutput(os.Stdout)

	// parse the flags
	flag.Parse()
	setupFields()

	// Welcome message
	requestCount, err := strconv.Atoi(*RequestCount)
	if err != nil {
		log.Printf("[Main]: RequestCount given is not a valid integer, reverting to default value: 1.")
		requestCount = 1
	}
	log.Printf("[Main]: Welcome. Client will send %d requests in parallel.", requestCount)

	// Set up a connection to the gRPC server
	serverFullAddr := fmt.Sprintf("%s:%s", *FrontAddr, *FrontPort)
	conn, err := grpc.Dial(serverFullAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("[Main]: Could not not connect. More:\n%v", err)
	}
	defer conn.Close()

	// create the client object
	c = pb.NewFrontClient(conn)

	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		// go choice()
		go convolutionRun()
		// go maxPoolRun(2)
	}

	// wait
	log.Printf("[Main]: All requests sent. Waiting for responses...")
	wg.Wait()

	log.Printf("[Main]: All requests completed. Terminating. Goodbye.")
}
