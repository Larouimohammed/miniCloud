package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	logger "github.com/Larouimohammed/miniCloud.git/logger"
	"github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/server"
)

func main() {
	logger := logger.Newlogger()
	server := server.DefaultServer
	if server == nil {
		logger.Logger.Error("Server initialisation failled")

	}
	if err := server.Run(); err != nil {
		logger.Logger.Sugar().Error(err)
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sig:
		logger.Logger.Sugar().Infow("shutdow starting", "signal", sig)
		time.Sleep(10 * time.Second)

	}
	defer logger.Logger.Sugar().Infow("shutdow complete", "signal", sig)
}
