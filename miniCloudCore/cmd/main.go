package main

import (
	"os"
	"os/signal"
	"syscall"

	logger "github.com/Larouimohammed/miniCloud.git/logger"
	"github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/server"
)

func main() {
	logger := logger.Newlogger()
	server := server.DefaultServer
	if server == nil {
		logger.Logger.Error("Server initialisation failled")

	}
	srv := server.Run()
	if srv == nil {
		logger.Logger.Sugar().Error("grpc server failled to start ")
	}
	
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	server.Close(srv,<-sig)
	
	
	

}
