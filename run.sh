#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
#./server -port=8001 &
./server -port=8002 -api=0 &
./server -port=8003 -api=0 &

sleep 2
#echo ">>> start test"
#curl "http://localhost:9999/api?key=a" &
#curl "http://localhost:9999/api?key=a" &
#curl "http://localhost:9999/api?key=a" &

wait

