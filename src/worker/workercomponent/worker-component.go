package worker

import (
	"context"
	"flag"
	"fmt"
	"log"
	pb "mini/proto"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	masterAddr = flag.String("masterAddr", "localhost", "The address to connect to")
	masterPort = flag.Int("masterPort", 55556, "The port of the master service")
	conn       *grpc.ClientConn
	c          pb.MasterClient
	workerAddr string
)

func dialServerAndSetClient() error {
	// Set up a connection to the gRPC server
	var err error

	if conn == nil {
		masterFullAddr := fmt.Sprintf("%s:%d", *masterAddr, *masterPort)
		conn, err = grpc.Dial(masterFullAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("[Worker]: Could not Dial. More:\n%v", err)
			return err
		}
	}

	// create the client object
	if c == nil {
		c = pb.NewMasterClient(conn)
	}

	return nil
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func NotifyWorkerActive(givenWorkerPort int) error {
	// Dial master
	err := dialServerAndSetClient()
	if err != nil {
		return err
	}

	// Save workerAddr
	workerAddr := fmt.Sprintf("%s:%d", GetOutboundIP().String(), givenWorkerPort)
	log.Printf("!!![Worker]: workerAddr: %s\n", workerAddr)

	// create the context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	// contact the server
	r, err := c.NotifyActiveWorker(ctx, &pb.NotifyRequest{WorkerAddress: workerAddr})
	if err != nil {
		log.Printf("[Worker]: Could not notify Master. More:\n%v", err)
		return err
	}

	// Print server response
	log.Printf("[Worker]: Server reply: %s", r.GetResult())

	return nil
}

func PingServer() error {
	dialServerAndSetClient()

	// create the context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// contact the server
	_, err := c.NotifyPing(ctx, &pb.NotifyRequest{WorkerAddress: workerAddr})

	if err != nil {
		log.Printf("[Worker]: Could not ping Master. More:\n%v", err)
	}

	return err
}
