package back

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	pb "mini/proto"
	"mini/utils"
	"net"
	"sync"

	"google.golang.org/grpc"
)

var (
	BackPort      = flag.String("BackPort", "", "The port of the back service.")
	active        int
	activeChannel = make(chan int, 1000)
	counter       int
	counterLock   sync.Mutex
	Prng          = rand.New(rand.NewSource(14))
	s             *grpc.Server
)

type backServer struct {
	pb.UnimplementedBackServer
}

// SayHello implements helloworld.GreeterServer
func (s *backServer) Choice(ctx context.Context, in *pb.ChoiceBiRequest) (*pb.ChoiceReply, error) {
	// signal as activated
	active += 1
	activeChannel <- active

	// log response
	counterLock.Lock()
	counter += 1
	replyID := counter
	counterLock.Unlock()
	log.Printf("[Back server]: Received request #%d", counter)

	// sleep
	utils.SimulatedCPUIntensiveFunction(1000, &active, 1)

	// randomly choose response
	var response string
	if Prng.Intn(2) == 0 {
		response = in.GetOption1()
	} else {
		response = in.GetOption2()
	}

	// signal as deactivated
	active -= 1
	activeChannel <- active

	// send response
	return &pb.ChoiceReply{Option: response, ReplyID: int32(replyID)}, nil
}

func (s *backServer) Convulution(ctx context.Context, in *pb.ConvolutionBackRequest) (*pb.ConvolutionBackReply, error) {
	// signal as activated
	active += 1
	activeChannel <- active

	// log response
	requestID := in.GetID()
	log.Printf("[Back server]: Received request #%d", requestID)

	// #TODO: implement convolution
	target := utils.ProtoToMatrix(in.GetTarget())
	kernel := utils.ProtoToMatrix(in.GetKernel())

	// Padding
	// resultWidth := len(result[0])
	// resultHeight := len(result)
	kernelWidth := len(kernel[0])
	kernelHeight := len(kernel)

	padTop := int(math.Floor((float64(kernelHeight) - 1.0) / 2.0))
	padBottom := int(math.Ceil((float64(kernelHeight) - 1.0) / 2.0))
	padLeft := int(math.Floor((float64(kernelWidth) - 1.0) / 2.0))
	padRight := int(math.Ceil((float64(kernelWidth) - 1.0) / 2.0))

	// Deep copy + padding
	paddedHeight := padTop + len(target) + padBottom
	paddedWidth := padLeft + len(target[0]) + padRight
	input := make([][]int32, paddedHeight)
	for i := 0; i < paddedHeight; i++ {
		input[i] = make([]int32, paddedWidth)
		for j := 0; j < paddedWidth; j++ {
			if (i < padTop) || (i >= padTop+len(target)) || (j < padLeft) || (j >= padLeft+len(target[0])) {
				input[i][j] = 0
			} else {
				input[i][j] = target[i-padTop][j-padLeft]
			}
		}
	}

	// Apply convolution
	for i := 0; i < len(target); i++ {
		for j := 0; j < len(target[0]); j++ {
			var newTarget int32 = 0
			for kw := 0; kw < kernelHeight; kw++ {
				for kh := 0; kh < kernelWidth; kh++ {
					newTarget += input[i+kh][j+kw] * kernel[kh][kw]
				}
			}
			target[i][j] = newTarget
		}
	}

	// #TODO: remove if working
	// utils.PrettyPrint("input", input)
	// utils.PrettyPrint("target", target)

	// signal as deactivated
	active -= 1
	activeChannel <- active

	// send response
	return &pb.ConvolutionBackReply{Result: utils.MatrixToProto(target), ID: requestID}, nil
}

func (s *backServer) MaxPool(ctx context.Context, in *pb.MaxPoolBackRequest) (*pb.MaxPoolBackReply, error) {
	// signal as activated
	active += 1
	activeChannel <- active

	// log response
	requestID := in.GetID()
	log.Printf("[Back server]: Received request #%d, Type: MaxPool", requestID)

	target := utils.ProtoToMatrix(in.GetTarget())
	targetHeight := len(target)
	targetWidth := len(target[0])

	// Max pool
	poolSize := int(in.GetSize())
	// pooledHeight := int(math.Ceil(float64(len(target)) / float64(poolSize)))
	pooledHeight := targetHeight / poolSize
	if targetHeight%poolSize != 0 {
		pooledHeight += 1
	}
	pooledWidth := targetWidth / poolSize
	if targetWidth%poolSize != 0 {
		pooledWidth += 1
	}

	result := utils.GenerateEmptyMatrix(pooledHeight, pooledWidth)

	for pi := 0; pi < pooledHeight; pi++ {
		for pj := 0; pj < pooledWidth; pj++ {
			var max int32 = 0
			for i := 0; i < poolSize; i++ {
				for j := 0; j < poolSize; j++ {
					ti := pi*poolSize + i
					tj := pj*poolSize + j
					if ti < targetHeight && tj < targetWidth && target[ti][tj] > max {
						max = target[ti][tj]
					}
				}
			}
			result[pi][pj] = max
		}
	}

	utils.PrettyPrint("result", result)

	// signal as deactivated
	active -= 1
	activeChannel <- active

	// send response
	return &pb.MaxPoolBackReply{Result: utils.MatrixToProto(result), ID: requestID}, nil
}

func debugActive() {
	for active := range activeChannel {
		log.Printf("[Back server]: Active rpc: %d", active)
	}
}

func listen() (net.Listener, error) {
	// listen to request to a free port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *BackPort))
	if err != nil {
		return nil, err
	}

	return lis, nil
}

func StartServer() error {
	// Listen
	// find free port and listen
	workerListener, err := listen()
	if err != nil {
		log.Printf("[Back server]: Failed to listen.\nMore: %v", err)
		return err
	}
	log.Printf("[Back server]: Back server listening at port: %s", *BackPort)

	// create a new server
	s = grpc.NewServer()

	// register the server
	pb.RegisterBackServer(s, &backServer{})

	// start debugging active level
	go debugActive()

	// serve the request
	if err := s.Serve(workerListener); err != nil {
		log.Fatalf("[Back server]: failed to serve: %v", err)
	}

	return nil
}

func StopServer(wg *sync.WaitGroup) {
	log.Printf("[Back server]: Grafecully stopping...")

	// Graceful stop
	s.GracefulStop()
	log.Printf("[Back server]: Done.")

	// Comunicate on channel so sync
	(*wg).Done()
}
