package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/hakobe/grpc-go-client-side-load-balancing-example/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var hostPort string

func init() {
	flag.StringVar(&hostPort, "hostport", "0.0.0.0:5000", "Server listening port")
	flag.Parse()
}

type server struct {
}

func (s *server) Echo(ctx context.Context, in *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{Message: in.GetMessage()}, nil
}

func serve(hostPort string) {
	lis, err := net.Listen("tcp", hostPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	echo.RegisterEchoServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	fmt.Printf("Starting server on %s ...\n", hostPort)
	serve(hostPort)
}
