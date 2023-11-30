package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/Larouimohammed/miniCloud.git/proto"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

type Config struct {
	Containername string `json:"containername"`
	Image         string `json:"image"`
	Subnet        string `json:"subnet"`
	Nunofinstance int32 `json:"numofinstance"`
}

func (C *Config) Goyaml(pathfile string) *Config {
	// read the output.yaml file
	data, err := os.ReadFile(pathfile)

	if err != nil {
		panic(err)
	}
	// fmt.Printf("members = %#v\n", string(data))
	var config Config
	databayte := []byte(data)
	// dataconfig := make([]Config, 0)
	err = yaml.Unmarshal([]byte(databayte), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// fmt.Printf("members = %#v\n", config)
	return &config
}

func main() {
	var config *Config
	config = config.Goyaml("config.yaml")
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSendClient(conn)

	// name := defaultName
	// if len(os.Args) > 1 {
	// 	name = os.Args[1]
	// }
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	r, err := c.Sendmsg(ctx, &pb.Msg{Containername: config.Containername, Image: config.Image, Subnet: config.Subnet, Nunofinstance: config.Nunofinstance})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Resp)
}
