package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Larouimohammed/miniCloud.git/logger"
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
	wg := &sync.WaitGroup{}
	wg.Add(1)
	// go func(wg *sync.WaitGroup) {
	// 	defer wg.Done()
	// 	print("qfqsfdsgfds\n")

	// }(wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		wg.Add(1)
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		server.Close(srv, <-sig)

	}(wg)
	wg.Wait()
}
