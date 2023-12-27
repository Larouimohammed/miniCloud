package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Larouimohammed/miniCloud.git/logger"
	"github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/server"
	"google.golang.org/grpc"
)

func main() {
	logger := logger.Newlogger()
	server := server.DefaultServer
	ServerChannel := make(chan *grpc.Server)
	wg := &sync.WaitGroup{}
	defer close(ServerChannel)
	wg.Add(2)
	if server == nil {
		logger.Logger.Error("Server initialisation failled")

	}
	go func() {
		defer wg.Done()
		
		srv := server.Run()
		ServerChannel <-srv
		if srv == nil {
			logger.Logger.Sugar().Error("grpc server failled to start ")
		}
	}()
	go func() {
		defer wg.Done()
		var serverS *grpc.Server
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		serverS =<-ServerChannel
		server.CloseServer(serverS,<-sig)
        

	}()

	wg.Wait()
}
