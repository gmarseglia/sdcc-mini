package back

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	pb "mini/proto"
	"mini/utils"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	HostPort           string = ""
	FakeCPUScale       float64
	Active             int
	ActiveChannel      = make(chan int, 1000)
	activeRequests     = make(map[CompositeID]bool)
	activeRequestsLock sync.RWMutex
	counter            int
	counterLock        sync.Mutex
	Prng               = rand.New(rand.NewSource(14))
	s                  *grpc.Server
)

type backServer struct {
	pb.UnimplementedBackServer
}

type CompositeID struct {
	ID      int32
	InnerID int32
}

// SayHello implements helloworld.GreeterServer
func (s *backServer) Choice(ctx context.Context, in *pb.ChoiceBiRequest) (*pb.ChoiceReply, error) {
	// signal as activated
	Active += 1
	ActiveChannel <- Active

	// log response
	counterLock.Lock()
	counter += 1
	replyID := counter
	counterLock.Unlock()
	log.Printf("[Back server]: Received request #%d", counter)

	// sleep
	utils.SimulatedCPUIntensiveFunction(1000, &Active, 1)

	// randomly choose response
	var response string
	if Prng.Intn(2) == 0 {
		response = in.GetOption1()
	} else {
		response = in.GetOption2()
	}

	// signal as deactivated
	Active -= 1
	ActiveChannel <- Active

	// send response
	return &pb.ChoiceReply{Option: response, ReplyID: int32(replyID)}, nil
}

func (s *backServer) Cancel(ctx context.Context, in *pb.CancelRequest) (*pb.CancelReply, error) {
	log.Printf("[Back server]: Cancel %d.%d", in.GetID(), in.GetInnerID())

	compositeID := CompositeID{ID: in.GetID(), InnerID: in.GetInnerID()}
	if _, exists := activeRequests[compositeID]; exists {
		activeRequestsLock.Lock()
		activeRequests[compositeID] = false
		activeRequestsLock.Unlock()
	}

	return &pb.CancelReply{Result: "OK"}, nil
}

func checkStillActive(compositeID CompositeID) bool {
	activeRequestsLock.RLock()
	defer activeRequestsLock.RUnlock()
	return activeRequests[compositeID]
}

func (s *backServer) ConvolutionalLayer(ctx context.Context, in *pb.ConvolutionalLayerBackRequest) (*pb.ConvolutionalLayerBackReply, error) {
	// signal as activated
	Active += 1
	ActiveChannel <- Active

	// log response
	compositeID := CompositeID{ID: in.GetID(), InnerID: in.GetInnerID()}
	// requestID := in.GetID()

	activeRequestsLock.Lock()
	activeRequests[compositeID] = true
	activeRequestsLock.Unlock()

	log.Printf("[Back server]: Received  %d.%d", compositeID.ID, compositeID.InnerID)

	// Prepare reply
	backReply := &pb.ConvolutionalLayerBackReply{ID: in.GetID(), InnerID: in.GetInnerID()}

	useSigmoid := in.GetUseSigmoid()

	// Padding
	target := utils.ProtoToMatrix(in.GetTarget())
	targetHeight := len(target)
	targetWidth := len(target[0])

	if in.GetUseKernels() {
		// Begin timing
		start := time.Now()

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

			if !checkStillActive(compositeID) {
				// signal as deactivated
				Active -= 1
				ActiveChannel <- Active
				return nil, status.Error(codes.Canceled, "cancelled as requested")
			}
			// Fake delay by CPU
			time.Sleep(time.Duration(time.Since(start).Seconds()*(FakeCPUScale-1)) * time.Second)

		}
	} else {
		start := time.Now()

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

		if !checkStillActive(compositeID) {
			// signal as deactivated
			Active -= 1
			ActiveChannel <- Active
			return nil, status.Error(codes.Canceled, fmt.Sprintf("Request %d.%d cancelled as requested", compositeID.ID, compositeID.InnerID))
		}
		// Fake delay by CPU
		time.Sleep(time.Duration(time.Since(start).Seconds()*(FakeCPUScale-1)) * time.Second)
	}

	activeRequestsLock.Lock()
	delete(activeRequests, compositeID)
	activeRequestsLock.Unlock()

	// signal as deactivated
	Active -= 1
	ActiveChannel <- Active

	log.Printf("[Back server]: Completed %d.%d", compositeID.ID, compositeID.InnerID)

	// send response
	return backReply, nil
}

func debugActive() {
	for Active := range ActiveChannel {
		log.Printf("[Back server]: Active rpc: %d", Active)
	}
}

func listen() (net.Listener, error) {
	// listen to request to a free port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", HostPort))
	if err != nil {
		return nil, err
	}

	return lis, nil
}

func StartServer(portChan chan string) error {
	// Listen
	// find free port and listen
	workerListener, err := listen()
	if err != nil {
		log.Printf("[Back server]: Failed to listen.\nMore: %v", err)
		return err
	}
	HostPort := fmt.Sprintf("%v", workerListener.Addr().(*net.TCPAddr).Port)
	log.Printf("[Back server]: Back server listening at port: %s", HostPort)

	// Comunicate port to worker-component
	portChan <- HostPort

	// create a new server
	s = grpc.NewServer()

	// register the server
	pb.RegisterBackServer(s, &backServer{})

	// start debugging Active level
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
