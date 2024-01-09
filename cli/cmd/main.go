package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/Larouimohammed/miniCloud.git/cli/config"
	pb "github.com/Larouimohammed/miniCloud.git/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProvClient(conn)

	if os.Args[1] == "apply" {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Timeout)*time.Second)
		defer cancel()
		r, err := c.Apply(ctx, &pb.Req{Containername: config.Containername, Image: config.Image, Subnet: config.Subnet, Nunofinstance: config.Replicas, Command: config.Command, AnsiblePlaybookPath: config.AnsiblePlaybookPath})
		if err != nil {
			log.Printf("Server can't provisionning infra: %v ", err)

		}
		log.Printf("Your infra %v was provisioned at : %s", config.Containername, r.Resp)
	}
	if os.Args[1] == "drop" {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Timeout)*time.Second)
		defer cancel()
		d, err := c.Drop(ctx, &pb.DReq{Containername: config.Containername, Nunofinstance: config.Replicas})
		if err != nil {
			log.Printf("Server can't Drop infra  : %v ", err)
		}
		log.Printf("your infra %v was droped at : %s", config.Containername, d.Resp)
	}
	if os.Args[1] == "update" {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Timeout)*time.Second)
		defer cancel()
		u, err := c.Update(ctx, &pb.Req{Containername: config.Containername, Image: config.Image, Subnet: config.Subnet, Nunofinstance: config.Replicas, Command: config.Command, AnsiblePlaybookPath: config.AnsiblePlaybookPath})
		if err != nil {
			log.Printf("Server can't Update infra  : %v ", err)
		}
		log.Printf("Update state: infra for %v updated at %s", config.Containername, u.Resp)
	}
	if os.Args[1] == "watch" {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Timeout)*time.Hour)
		defer cancel()
		// client := Message.NewSendClient(conn)

		stream, err := c.Watch(ctx, &pb.WReq{Containername: config.Containername})
		if err != nil {
			log.Printf("Server can't watch infra  : %v ", err)
		}

		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				log.Printf("End of file error %v", err)
			}
			if err != nil {
				log.Printf("%v.ListFeatures(_) = _, %v", c, err)
			}
			if msg.Werr == "" && msg.Wresp != "" {
				log.Printf("Server Start Watch Streaming With Message : %v ", msg.Wresp)
			}
			if msg.Werr != "" && msg.Wresp == "" {
				log.Printf("Server Start Watch Streaming with Error : %v ", msg.Werr)

			} else {
				log.Printf("Server Start Watch Streaming with Error :%v \n and Message : %v ", msg.Werr, msg.Wresp)
			}

		}
	}

}
