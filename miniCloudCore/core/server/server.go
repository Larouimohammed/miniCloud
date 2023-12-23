package server

import (
	"context"
	"flag"
	"fmt"

	"net"
	t "time"

	log "github.com/Larouimohammed/miniCloud.git/logger"

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
	logger log.Log
}

func NewServer() *Server {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	logger := log.Newlogger()
	if err != nil {
		logger.Logger.Sugar().Error(" initialisation docker client error : %v", err)
		return nil
	}
	// defer client.Close()
	return &Server{
		cli: client,
		logger: *logger,
	}

}

var DefaultServer = NewServer()


// provisioning
func (S *Server) Apply(ctx context.Context, config *pb.Req) (*pb.Resp, error) {
	S.logger.Logger.Sugar().Infow("Provisionning infra Starting","CN :",config.Containername,"Image :",config.Image,"Numofinstance :",config.Nunofinstance)
	if err := command.ProvApply(S.cli, config.Containername, config.Image, config.Subnet, config.Nunofinstance); err != nil {
		S.logger.Logger.Sugar().Error(" provisionning error : %v", err)
		return &pb.Resp{Resp: "provisionning infra error"}, err
	}
	return &pb.Resp{Resp: t.Now().String()}, nil
}

// droping
func (S *Server) Drop(ctx context.Context, config *pb.DReq) (*pb.Resp, error) {
	S.logger.Logger.Sugar().Infow("Droping infra Starting","CN :",config.Containername,"Numofinstance :",config.Nunofinstance)
	if err := command.StopandDropContainer(S.cli, config.Containername, config.Nunofinstance); err != nil {
		S.logger.Logger.Sugar().Error(" droping infra  error : %v", err)
	}

	return &pb.Resp{Resp: t.Now().String()}, nil
}

// updating
func (S *Server) Update(ctx context.Context, config *pb.Req) (*pb.Resp, error) {
	S.logger.Logger.Sugar().Infow("Updating infra Starting")
	instance, err := command.Watching(S.cli, config.Containername)
	if err != nil {
		S.logger.Logger.Sugar().Error("number of instance is indectectible : %v", err)
	}
	if err := command.StopandDropContainer(S.cli, config.Containername, instance); err != nil {
		S.logger.Logger.Sugar().Error(" droping infra  error : %v", err)
	}
	if err := command.ProvApply(S.cli, config.Containername, config.Image, config.Subnet, config.Nunofinstance); err != nil {
		S.logger.Logger.Sugar().Error(" provisionning error : %v", err)
		return &pb.Resp{Resp: "provisionning infra error"}, err

	}
	return &pb.Resp{Resp: t.Now().String()}, nil

}

// watching
func (S *Server) Watch(ctx context.Context, config *pb.WReq) (*pb.WResp, error) {
	S.logger.Logger.Sugar().Infow("Watching of infra Starting", config.Containername)
	instance, err := command.Watching(S.cli, config.Containername)
	if err != nil {
		S.logger.Logger.Sugar().Error("Watchinh error : %v", err)
		return &pb.WResp{Wresp: 01}, nil
	}

	return &pb.WResp{Wresp: instance}, nil
}

func (S *Server) Run() error {

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		S.logger.Logger.Sugar().Error("failed to listen: %v", err)
		return err
	}
	s := grpc.NewServer()
	pb.RegisterProvServer(s, S)
	S.logger.Logger.Sugar().Infow("Server Starting",  "listing on", lis.Addr())
	if err := s.Serve(lis); err != nil {
		S.logger.Logger.Sugar().Error("failed to serve: %v", err)
		return err
	}

	return nil
}
