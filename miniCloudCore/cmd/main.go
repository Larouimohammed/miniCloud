package main

import (
	"log"

	"github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/server"
)

var Server *server.Server

func main() {

	// server, err := Server.NewServer()
	// if err != nil {
	// 	log.Printf("Server initialisation : %v", err)

	// }
	server := Server.NewServer()
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
