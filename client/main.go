package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	trylb "github.com/hakobe/grpc-try-load-balancing/trylb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"
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

func init() {
	flag.Var(&serverAddrs, "server", "specify server hostports")
	flag.Parse()
}

type resolver struct {
	addrs []string
}

func (r *resolver) Resolve(target string) (naming.Watcher, error) {
	w := &watcher{
		updatesChan: make(chan []*naming.Update, 1),
	}
	updates := []*naming.Update{}
	for _, addr := range r.addrs {
		updates = append(updates, &naming.Update{Op: naming.Add, Addr: addr})
	}
	w.updatesChan <- updates
	return w, nil
}

type watcher struct {
	updatesChan chan []*naming.Update
}

func (w *watcher) Next() ([]*naming.Update, error) {
	us, ok := <-w.updatesChan
	if !ok {
		return nil, errWatcherClose
	}
	return us, nil
}

func (w *watcher) Close() {
	close(w.updatesChan)
}

func callEcho(client trylb.EchoServiceClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var p peer.Peer
	res, err := client.Echo(
		ctx,
		&trylb.EchoRequest{Message: message},
		grpc.FailFast(false),
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
		grpc.WithBalancer(grpc.RoundRobin(&resolver{addrs: serverAddrs})),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := trylb.NewEchoServiceClient(conn)

	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			callEcho(c, fmt.Sprintf("hello %d", i))
			wg.Done()
		}(i)
	}
	wg.Wait()
}
