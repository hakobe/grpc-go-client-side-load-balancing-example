package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	trylb "github.com/hakobe/grpc-try-load-balancing/trylb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var hostPort string

func init() {
	flag.StringVar(&hostPort, "hostport", "0.0.0.0:5000", "Server listening port")
	flag.Parse()
}

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

	zapLogger, _ := zap.NewProduction()
	opts := []grpc_zap.Option{}
	grpc_zap.ReplaceGrpcLogger(zapLogger)
	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(zapLogger, opts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(zapLogger, opts...),
		),
	)
	trylb.RegisterEchoServiceServer(s, &Server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	fmt.Println("Starting server...")
	serve(hostPort)
}
