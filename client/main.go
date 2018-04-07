package main

import (
	"context"
	"fmt"
	"log"
	"time"

	trylb "github.com/hakobe/grpc-try-load-balancing/trylb"
	"google.golang.org/grpc"
)

func getConn(hostPort string) *grpc.ClientConn {
	conn, err := grpc.Dial(hostPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

func main() {
	fmt.Println("Starting client...")
	hostPort := "server1:5000"
	conn := getConn(hostPort)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	c := trylb.NewEchoServiceClient(conn)

	res, err := c.Echo(ctx, &trylb.EchoRequest{Message: "hello"})
	if err != nil {
		log.Fatalf("could not call echo: %v", err)
	}
	log.Printf("received message is %s\n", res.GetMessage())
}
