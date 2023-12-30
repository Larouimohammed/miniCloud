package server

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"net"
	t "time"

	log "github.com/Larouimohammed/miniCloud.git/logger"
	"github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/command"

	consul "github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/consulproxy"
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
	logger       log.Log
	consulClient *consul.ConsulProxy
}

func NewServer() *Server {
	consulClient := consul.DefaultConsulProxy
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	logger := log.Newlogger()
	if err != nil {
		logger.Logger.Sugar().Error("initialisation docker client error %v", err)
		return nil
	}
	// defer client.Close()
	return &Server{
		cli:          client,
		logger:       *logger,
		consulClient: consulClient,
	}

}

var DefaultServer = NewServer()

// provisioning
func (S *Server) Apply(ctx context.Context, config *pb.Req) (*pb.Resp, error) {
	S.logger.Logger.Sugar().Infow("Provisionning infra Starting", "CN", config.Containername, "Image", config.Image, "Numofinstance", config.Nunofinstance)
	defer S.logger.Logger.Sugar().Infow("Provisionning infra complete successfully")
	if err := command.ProvApply(S.cli, config.Containername, config.Image, config.Subnet, config.Nunofinstance, S.logger, S.consulClient); err != nil {
		S.logger.Logger.Sugar().Error(" provisionning error %v", err)
		return &pb.Resp{Resp: "provisionning infra error"}, err
	}
	return &pb.Resp{Resp: t.Now().String()}, nil
}

// droping
func (S *Server) Drop(ctx context.Context, config *pb.DReq) (*pb.Resp, error) {
	S.logger.Logger.Sugar().Infow("Droping infra Starting", "CN", config.Containername, "Numofinstance", config.Nunofinstance)
	defer S.logger.Logger.Sugar().Infow("Droping infra complete successfully")
	if err := command.StopandDropContainer(S.cli, config.Containername, config.Nunofinstance); err != nil {
		S.logger.Logger.Sugar().Error(" droping infra  error %v", err)
	}

	return &pb.Resp{Resp: t.Now().String()}, nil
}

// updating
func (S *Server) Update(ctx context.Context, config *pb.Req) (*pb.Resp, error) {
	S.logger.Logger.Sugar().Infow("Updating infra Starting")
	defer S.logger.Logger.Sugar().Infow("Updating infra complete successfully")
	instance, err := command.Watching(S.cli, config.Containername, S.logger)
	if err != nil {
		S.logger.Logger.Sugar().Error("number of instance is indectectible  %v", err)
	}
	if err := command.StopandDropContainer(S.cli, config.Containername, instance); err != nil {
		S.logger.Logger.Sugar().Error(" droping infra  error  %v", err)
	}
	if err := command.ProvApply(S.cli, config.Containername, config.Image, config.Subnet, config.Nunofinstance, S.logger, S.consulClient); err != nil {
		S.logger.Logger.Sugar().Error(" provisionning error  %v", err)
		return &pb.Resp{Resp: "provisionning infra error"}, err

	}
	return &pb.Resp{Resp: t.Now().String()}, nil

}

// watching
func (S *Server) Watch(ctx context.Context, config *pb.WReq) (*pb.WResp, error) {
	S.logger.Logger.Sugar().Infow("Watching of infra Starting", config.Containername)
	instance, err := command.Watching(S.cli, config.Containername, S.logger)
	if err != nil {
		S.logger.Logger.Sugar().Error("Watchinh error : %v", err)
		return &pb.WResp{Wresp: 01}, nil
	}

	return &pb.WResp{Wresp: instance}, nil
}

func (S *Server) Run() {

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		S.logger.Logger.Sugar().Fatal("failed to listen %v", err)

	}
	s := grpc.NewServer()
	pb.RegisterProvServer(s, S)
	S.logger.Logger.Sugar().Infow("Server Starting", "listing on", lis.Addr())

	sigCh := make(chan os.Signal, 1)
	/*  when sigCh channel gets a signal notify me */
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		sig := <-sigCh
		S.CloseServer(s, sig)
		wg.Done()
	}()

	if err := s.Serve(lis); err != nil {
		S.logger.Logger.Sugar().Fatal("failed to serve: %v", err)

	}
	wg.Wait()
}

func (s *Server) CloseServer(grpcserver *grpc.Server, sig os.Signal) {
	defer s.logger.Logger.Sugar().Infow("shutdow complete", "signal", sig)
	s.logger.Logger.Sugar().Infow("shutdow starting", "signal", sig)

	if err := s.cli.Close(); err != nil {
		s.logger.Logger.Sugar().Error(err)

	}
	grpcserver.GracefulStop()

}
