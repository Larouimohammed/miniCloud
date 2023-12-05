package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Larouimohammed/miniCloud.git/cli/config"
	pb "github.com/Larouimohammed/miniCloud.git/proto"
	"google.golang.org/grpc"
)

var (
	Address = "localhost:50051"
	Timeout = 100
)

func main() {
	var config *config.Config
	config = config.Goyaml("./config.yaml")
	address := Address
	// if len(os.Args) > 1 {
	// 	address = os.Args[1]
	// }
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProvClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Timeout)*time.Second)
	defer cancel()
	if os.Args[1] == "apply" {
		r, err := c.Apply(ctx, &pb.Req{Containername: config.Containername, Image: config.Image, Subnet: config.Subnet, Nunofinstance: config.Replicas})
		if err != nil {
			log.Printf("Server can't provisionning infra: %v ", err)

		}
		log.Printf("Respending: %s", r.Resp)
	}
	if os.Args[1] == "drop" {

		d, err := c.Drop(ctx, &pb.Req{Containername: config.Containername, Image: config.Image, Subnet: config.Subnet, Nunofinstance: config.Replicas})
		if err != nil {
			log.Printf("Server can't Drop infra  : %v ", err)
		}
		log.Printf("Respending: %s", d.Resp)
	}
}
