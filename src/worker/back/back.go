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

func (s *backServer) Convulution(ctx context.Context, in *pb.ConvolutionalLayerBackRequest) (*pb.ConvolutionalLayerBackReply, error) {
	// signal as activated
	active += 1
	activeChannel <- active

	// log response
	requestID := in.GetID()
	log.Printf("[Back server]: Received request #%d", requestID)

	// Prepare reply
	backReply := &pb.ConvolutionalLayerBackReply{ID: requestID}

	useSigmoid := in.GetUseSigmoid()

	// Padding

	target := utils.ProtoToMatrix(in.GetTarget())
	targetHeight := len(target)
	targetWidth := len(target[0])

	if in.GetUseKernels() {
		kernelWidth := len(utils.ProtoToMatrix(in.Kernel[0]))
		kernelHeight := len(utils.ProtoToMatrix(in.Kernel[0])[0])

		padTop := int(math.Floor((float64(kernelHeight) - 1.0) / 2.0))
		padBottom := int(math.Ceil((float64(kernelHeight) - 1.0) / 2.0))
		padLeft := int(math.Floor((float64(kernelWidth) - 1.0) / 2.0))
		padRight := int(math.Ceil((float64(kernelWidth) - 1.0) / 2.0))

		for _, kernelProto := range in.GetKernel() {
			kernel := utils.ProtoToMatrix(kernelProto)

			// Deep copy + padding
			paddedHeight := padTop + targetHeight + padBottom
			paddedWidth := padLeft + targetWidth + padRight
			input := make([][]float32, paddedHeight)
			for i := 0; i < paddedHeight; i++ {
				input[i] = make([]float32, paddedWidth)
				for j := 0; j < paddedWidth; j++ {
					if (i < padTop) || (i >= padTop+targetHeight) || (j < padLeft) || (j >= padLeft+targetWidth) {
						input[i][j] = 0
					} else {
						input[i][j] = target[i-padTop][j-padLeft]
					}
				}
			}

			// Apply convolution + sigmoid
			convoluted := utils.GenerateEmptyMatrix(targetHeight, targetWidth)
			for i := 0; i < targetHeight; i++ {
				for j := 0; j < targetWidth; j++ {
					var newTarget float32 = 0
					for kw := 0; kw < kernelHeight; kw++ {
						for kh := 0; kh < kernelWidth; kh++ {
							newTarget += input[i+kh][j+kw] * kernel[kh][kw]
						}
					}
					if useSigmoid {
						// Apply sigmoid
						convoluted[i][j] = float32(1.0 / (1.0 + math.Exp(float64(-newTarget))))
					} else {
						convoluted[i][j] = newTarget
					}
				}
			}

			// Compute size after pooling
			poolSize := int(in.GetAvgPoolSize())
			if poolSize > 1 {
				pooledHeight := int(targetHeight / poolSize)
				pooledWidth := int(targetWidth / poolSize)

				// Prepare result
				result := utils.GenerateEmptyMatrix(pooledHeight, pooledWidth)

				// Apply AvgPool
				for pi := 0; pi < pooledHeight; pi++ {
					for pj := 0; pj < pooledWidth; pj++ {
						var avgValue float32 = 0
						for i := 0; i < poolSize; i++ {
							for j := 0; j < poolSize; j++ {
								ti := pi*poolSize + i
								tj := pj*poolSize + j
								if ti < targetHeight && tj < targetWidth {
									avgValue += convoluted[ti][tj] / float32(poolSize*poolSize)
								}
							}
						}
						result[pi][pj] = avgValue
					}
				}
				// Append the result
				backReply.Result = append(backReply.Result, utils.MatrixToProto(result))
			} else {
				// Append the result
				backReply.Result = append(backReply.Result, utils.MatrixToProto(convoluted))
			}
		}
	} else {
		poolSize := int(in.GetAvgPoolSize())
		if poolSize > 1 {
			pooledHeight := int(targetHeight / poolSize)
			pooledWidth := int(targetWidth / poolSize)

			// Prepare result
			result := utils.GenerateEmptyMatrix(pooledHeight, pooledWidth)

			// Apply AvgPool
			for pi := 0; pi < pooledHeight; pi++ {
				for pj := 0; pj < pooledWidth; pj++ {
					var avgValue float32 = 0
					for i := 0; i < poolSize; i++ {
						for j := 0; j < poolSize; j++ {
							ti := pi*poolSize + i
							tj := pj*poolSize + j
							if ti < targetHeight && tj < targetWidth {
								avgValue += target[ti][tj] / float32(poolSize*poolSize)
							}
						}
					}
					result[pi][pj] = avgValue
				}
			}
			// Append the result
			backReply.Result = append(backReply.Result, utils.MatrixToProto(result))
		} else {
			backReply.Result = append(backReply.Result, utils.MatrixToProto(target))
		}
	}

	// signal as deactivated
	active -= 1
	activeChannel <- active

	// send response
	return backReply, nil
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
