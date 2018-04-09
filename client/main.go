package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hakobe/grpc-try-load-balancing/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var errWatcherClose = errors.New("watcher has been closed")

type addrsValue []string

func (as *addrsValue) String() string {
	return "localhost:5000"
}

func (as *addrsValue) Set(addr string) error {
	*as = append(*as, addr)
	return nil
}

var serverAddrs addrsValue
var callTimes int

func init() {
	flag.Var(&serverAddrs, "server", "Server hostports")
	flag.IntVar(&callTimes, "n", 10000, "Times to call")
	flag.Parse()
}

func callEcho(client echo.EchoServiceClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var p peer.Peer
	res, err := client.Echo(
		ctx,
		&echo.EchoRequest{Message: message},
		grpc.FailFast(false), // To wait a resolver returning addrs.
		grpc.Peer(&p),
	)
	if err != nil {
		log.Fatalf("could not call echo: %v", err)
	}
	log.Printf("from: %s, received: %s\n", p.Addr, res.GetMessage())
}

func main() {
	fmt.Println("Starting client...")

	fmt.Println(serverAddrs)

	conn, err := grpc.Dial(
		"dummy",
		grpc.WithInsecure(),
		grpc.WithBalancer(grpc.RoundRobin(NewPseudoResolver(serverAddrs))),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := echo.NewEchoServiceClient(conn)

	var wg sync.WaitGroup
	for i := 0; i < callTimes; i++ {
		wg.Add(1)
		go func(i int) {
			callEcho(c, fmt.Sprintf("hello %d", i))
			wg.Done()
		}(i)
	}
	wg.Wait()
}
