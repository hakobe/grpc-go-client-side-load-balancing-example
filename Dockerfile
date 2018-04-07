FROM golang:1.10-stretch

# install protoc
RUN apt-get update && apt-get -y install unzip && apt-get clean
ENV PB_VER 3.5.1
ENV PB_URL https://github.com/google/protobuf/releases/download/v${PB_VER}/protoc-${PB_VER}-linux-x86_64.zip
RUN mkdir /protoc && \
    curl -L ${PB_URL} > /protoc/protoc.zip && \
    cd /protoc && \
    unzip protoc.zip
RUN go get -u github.com/golang/protobuf/protoc-gen-go

WORKDIR /go/src/github.com/hakobe/grpc-try-load-balancing
COPY . .

RUN go get -d -v ./...
RUN go build -o server/server github.com/hakobe/grpc-try-load-balancing/server
RUN go build -o client/client github.com/hakobe/grpc-try-load-balancing/client

CMD ["./server/server"]
