package worker

import (
	"context"
	"flag"
	"log"
	pb "mini/proto"
	"mini/worker/back"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	MasterAddr     = flag.String("MasterAddr", "", "The address to connect to.")
	MasterPort     = flag.String("MasterPort", "", "The port of the master service.")
	MasterFullAddr string
	conn           *grpc.ClientConn
	c              pb.MasterClient
)

func dialServerAndSetClient() error {
	// Set up a connection to the gRPC server
	var err error

	if conn == nil {
		conn, err = grpc.Dial(MasterFullAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func NotifyWorkerActive() error {
	// Dial master
	err := dialServerAndSetClient()
	if err != nil {
		return err
	}

	// Save workerAddr
	log.Printf("[Worker]: Address to be advertised to master: %s\n", back.HostFullAddr)

	// create the context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	// contact the server
	r, err := c.NotifyActiveWorker(ctx, &pb.NotifyRequest{WorkerAddress: back.HostFullAddr})
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
	_, err := c.NotifyPing(ctx, &pb.NotifyRequest{WorkerAddress: back.HostFullAddr})

	if err != nil {
		log.Printf("[Worker]: Could not ping Master. More:\n%v", err)
	}

	return err
}
