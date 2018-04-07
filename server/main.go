package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	trylb "github.com/hakobe/grpc-try-load-balancing/trylb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
}

func (s *Server) Echo(ctx context.Context, in *trylb.EchoRequest) (*trylb.EchoResponse, error) {

	return &trylb.EchoResponse{Message: in.GetMessage()}, nil
}

func serve(hostPort string) {
	lis, err := net.Listen("tcp", hostPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	trylb.RegisterEchoServiceServer(s, &Server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	fmt.Println("Starting server...")
	hostPort := os.Getenv("HOST_PORT")
	if hostPort == "" {
		hostPort = "0.0.0.0:5000"
	}
	serve(hostPort)
}
