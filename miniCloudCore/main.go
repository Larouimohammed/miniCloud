package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/Larouimohammed/miniCloud.git/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedSendServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) Sendmsg(ctx context.Context, in *pb.Msg) (*pb.Resp, error) {
	log.Printf("CN: %v  Image:%v Subnet %v Numofinstance %v", in.Containername, in.Image, in.Subnet, in.Nunofinstance)
	//provision infra
	
	
	return &pb.Resp{Resp: "thank you khero"}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSendServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
