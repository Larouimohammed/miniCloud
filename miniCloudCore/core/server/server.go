package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	t "time"

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
	log.Printf("Provisionning infra Starting")
	log.Printf("CN: %v  Image:%v Subnet %v Numofinstance %d", config.Containername, config.Image, config.Subnet, config.Nunofinstance)
	if err := command.ProvApply(S.cli, config.Containername, config.Image, config.Subnet, config.Nunofinstance); err != nil {
		log.Printf(" provisionning error : %v", err)
		return &pb.Resp{Resp: "provisionning infra error"}, err
	}
	return &pb.Resp{Resp: t.Now().String()}, nil
}

// droping
func (S *Server) Drop(ctx context.Context, config *pb.DReq) (*pb.Resp, error) {
	log.Printf("Droping infra Starting")
	log.Printf("Droping CN: %v Numofinstance %d", config.Containername, config.Nunofinstance)

	if err := command.StopandDropContainer(S.cli, config.Containername, config.Nunofinstance); err != nil {
		log.Printf(" droping infra  error : %v", err)
	}

	return &pb.Resp{Resp: t.Now().String()}, nil
}

// updating
func (S *Server) Update(ctx context.Context, config *pb.Req) (*pb.Resp, error) {
	log.Printf("Updating infra Starting")
	instance, err := command.Watching(S.cli, config.Containername)
	if err != nil {
		log.Printf("number of instance is indectectible : %v", err)
	}
	if err := command.StopandDropContainer(S.cli, config.Containername, instance); err != nil {
		log.Printf(" droping infra  error : %v", err)
	}
	if err := command.ProvApply(S.cli, config.Containername, config.Image, config.Subnet, config.Nunofinstance); err != nil {
		log.Printf(" provisionning error : %v", err)
		return &pb.Resp{Resp: "provisionning infra error"}, err

	}
	return &pb.Resp{Resp: t.Now().String()}, nil

}

// watching
func (S *Server) Watch(ctx context.Context, config *pb.WReq) (*pb.WResp, error) {
	log.Printf("Watching %v infra Starting", config.Containername)
	instance, err := command.Watching(S.cli, config.Containername)
	if err != nil {
		log.Printf("Watchinh error : %v", err)
		return &pb.WResp{Wresp: 01}, nil
	}

	return &pb.WResp{Wresp: instance}, nil
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
