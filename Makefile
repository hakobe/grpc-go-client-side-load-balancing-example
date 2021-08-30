all: build
get:
	go get -d -v ./...
build-server:
	go build -o server/server github.com/hakobe/grpc-go-client-side-load-balancing-example/server
build-client:
	go build -o client/client github.com/hakobe/grpc-go-client-side-load-balancing-example/client
gen-pb:
	mkdir -p echo && protoc ./echo.proto --go_out=.
build: gen-pb get build-server build-client

.PHONY: \
	build-server \
	build-client \
	gen-pb \
	build \
	all
