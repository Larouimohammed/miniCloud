package main

import (
	"log"

	"github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/server"
)
func main() {
	server := server.DefaultServer
	if server == nil {
		log.Printf("Server initialisation failled")

	}
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
