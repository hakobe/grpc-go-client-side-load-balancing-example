#!/bin/bash
trap 'kill $(jobs -p)' EXIT

./server/server -hostport 0.0.0.0:5000 &
./server/server -hostport 0.0.0.0:5001 &
./server/server -hostport 0.0.0.0:5002 &
./server/server -hostport 0.0.0.0:5003 &

sleep 3

time ./client/client -n 10000 \
    -server localhost:5000 \
    -server localhost:5001 \
    -server localhost:5002 \
    -server localhost:5003 \
