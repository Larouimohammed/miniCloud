package main

import (
	"log"

	"github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/server"
)

var Server *server.Server

func main() {

	server, err := Server.NewServer()
	if err == nil {
		log.Printf(err.Error())
	}
	if err := server.Run(); err == nil {
		log.Fatal(err)
	}

}
