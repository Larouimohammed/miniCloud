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
	dockerclient *client.Client
	pb.UnimplementedProvServer
}

func (S *Server) NewServer() *Server {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf(err.Error())
	}
	defer cli.Close()
	return &Server{
		dockerclient: cli,
	}

}

func (s *Server) Apply(ctx context.Context, config *pb.Req) (*pb.Resp, error) {
	log.Printf("CN: %v  Image:%v Subnet %v Numofinstance %v", config.Containername, config.Image, config.Subnet, config.Nunofinstance)
	//provision
	command.ProvApply(*s.dockerclient, config.Containername, config.Image, config.Subnet, config.Nunofinstance)
	return &pb.Resp{Resp: "your miniCloud is provisioned say :thank you khero"}, nil
}

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
