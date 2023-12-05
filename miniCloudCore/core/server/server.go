package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/command"
	pb "github.com/Larouimohammed/miniCloud.git/proto"
	"github.com/docker/docker/client"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type Server struct {
	cli *client.Client
	pb.UnimplementedProvServer
}

func (S *Server) NewServer() (*Server, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf(" initialisation docker client error : %v", err)
		return nil, err
	}
	defer cli.Close()
	return &Server{
		cli: cli,
	}, nil

}

// provisioning
func (S *Server) Apply(ctx context.Context, config *pb.Req) (*pb.Resp, error) {
	log.Printf("CN: %v  Image:%v Subnet %v Numofinstance %d", config.Containername, config.Image, config.Subnet, config.Nunofinstance)
	if err := command.ProvApply(S.cli, config.Containername, config.Image, config.Subnet, config.Nunofinstance); err != nil {
		log.Printf(" provisionning error : %v", err)
		return &pb.Resp{Resp: "provisionning infra error"}, err
	}
	return &pb.Resp{Resp: "your infra was provisioned"}, nil
}

// droping
func (S *Server) Drop(ctx context.Context, config *pb.Req) (*pb.Resp, error) {
	log.Printf("Droping CN: %v  Image:%v Subnet %v Numofinstance %d", config.Containername, config.Image, config.Subnet, config.Nunofinstance)

	return &pb.Resp{Resp: "your infra was droped"}, nil
}

// watching
// updating

func (S *Server) Run() error {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}
	s := grpc.NewServer()
	pb.RegisterProvServer(s, &Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}

	return nil
}
