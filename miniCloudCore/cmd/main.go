package main

import (
	"log"

	"github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/server"
)

var Server *server.Server

func main() {
	server := Server.NewServer()
	if err := server.Run(); err == nil {
		log.Fatal(err)
	}

}
