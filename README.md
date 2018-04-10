# gRPC client side load balancing example using Go

This is a gRPC client-side load balancing example on top of grpc-go.

The only load balancer bundled to [grpc-go](https://github.com/grpc/grpc-go) is `grpc.RoundRobin`. grpc RoundRoubin requires a `grpc.Resolver` which is intended to implement a DNS resolver or an other resourse resolution mechanism like Consul.

In this example I implemented a `grpc.Resolver` which only returns fixed servers initially passed. It is good for trying gRPC load balancing instantly.

## Prequirements

- [protoc](https://github.com/google/protobuf)
- [protoc-gen-go](https://github.com/golang/protobuf/tree/master/protoc-gen-go)
  - `$ go get -u github.com/golang/protobuf/protoc-gen-go`

## Build

```
$ make
```

## Run

```console
# Run 4 servers
$ ./server/server -hostport 0.0.0.0:5000 &
$ ./server/server -hostport 0.0.0.0:5001 &
$ ./server/server -hostport 0.0.0.0:5002 &
$ ./server/server -hostport 0.0.0.0:5003 &

$ sleep 3 # Wait for server to start up

# Do gRPC method calls 10000 times
$ time ./client/client -n 10000 \
    -server localhost:5000 \
    -server localhost:5001 \
    -server localhost:5002 \
    -server localhost:5003 \
 
# or 
$ ./run.sh
```

# LICENSE
MIT @hakobe
