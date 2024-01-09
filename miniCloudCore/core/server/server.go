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
	"github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/ansible"
	"github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/command"
	"github.com/docker/docker/api/types"
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
	err := command.ProvApply(S.cli, config.Containername, config.Image, config.Subnet, config.AnsiblePlaybookPath, config.Nunofinstance, config.Command, S.logger, S.consulClient)
	if err != nil {
		S.logger.Logger.Sugar().Error(" provisionning infara error %v", err)
		return &pb.Resp{Resp: "provisionning infra error"}, err
	}
	S.logger.Logger.Sugar().Infow("Ansible Path", "ansible Path", config.AnsiblePlaybookPath)

	if config.AnsiblePlaybookPath != "" {
		if err := ansible.RunAnsible(config.AnsiblePlaybookPath, "", S.logger); err != nil {
			S.logger.Logger.Sugar().Error(" provisionning infara error %v", err)

			return &pb.Resp{Resp: "Ansible install error"}, err
		}
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
	instance, err := command.GetInstance(ctx, S.cli, config.Containername, S.logger)
	if err != nil {
		S.logger.Logger.Sugar().Error("number of instance is indectectible  %v", err)
	}
	if err := command.StopandDropContainer(S.cli, config.Containername, instance); err != nil {
		S.logger.Logger.Sugar().Error(" droping infra  error  %v", err)
	}
	if err := command.ProvApply(S.cli, config.Containername, config.Image, config.Subnet, config.AnsiblePlaybookPath, config.Nunofinstance, config.Command, S.logger, S.consulClient); err != nil {
		S.logger.Logger.Sugar().Error(" provisionning error  %v", err)
		return &pb.Resp{Resp: "provisionning infra error"}, err

	}
	S.logger.Logger.Sugar().Infow(config.AnsiblePlaybookPath)
	if config.AnsiblePlaybookPath != "" {
		if err := ansible.RunAnsible(config.AnsiblePlaybookPath, "", S.logger); err != nil {
			S.logger.Logger.Sugar().Error(" provisionning infara error %v", err)

			return &pb.Resp{Resp: "Ansible install error"}, err
		}
	}
	return &pb.Resp{Resp: t.Now().String()}, nil

}

// watching
func (S *Server) Watch(config *pb.WReq, stream pb.Prov_WatchServer) error {
	S.logger.Logger.Sugar().Infow("Watching of infra Starting", config.Containername)
	defer S.logger.Logger.Sugar().Infow("Watching of infra shuting", config.Containername)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {

		defer wg.Done()

		for {
			msgs, serrs := S.cli.Events(context.Background(), types.EventsOptions{})
			select {
			case msg := <-msgs:
				err := stream.Send(&pb.WResp{Wresp: fmt.Sprintf("%+v", msg), Werr: ""})
				if err != nil {
					S.logger.Logger.Sugar().Error(err)

				}
			case errs := <-serrs:
				err := stream.Send(&pb.WResp{Wresp: "", Werr: errs.Error()})
				if err != nil {
					S.logger.Logger.Sugar().Error(err)
				}

			}

		}
	}(&wg)
	wg.Wait()
	return nil
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
	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		sig := <-sigCh
		S.CloseServer(s, sig)

	}(&wg)

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
