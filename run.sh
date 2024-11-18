#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
#./server -port=8001 &
#./server -port=8002 -api=0 &
#./server -port=8003 -api=0 &
./server -port=8001 &
./server -port=8002 &
./server -port=8003 -api=1 &

sleep 2
echo ">>> start test"


#curl "http://localhost:9999/api?key=d" &
#curl "http://localhost:9999/api?key=d" &
#curl "http://localhost:9999/api?key=d" &
#curl "http://localhost:9999/api?key=d" &
#curl "http://localhost:9999/api?key=d" &
#curl "http://localhost:9999/api?key=d" &
#curl "http://localhost:9999/api?key=d" &
#curl "http://localhost:9999/api?key=d" &
#curl "http://localhost:9999/api?key=d" &
#curl "http://localhost:9999/api?key=d" &
curl "http://localhost:9999/api?key=a" &
curl "http://localhost:9999/api?key=a" &
curl "http://localhost:9999/api?key=a" &
curl "http://localhost:9999/api?key=a" &
curl "http://localhost:9999/api?key=a" &
curl "http://localhost:9999/api?key=a" &
curl "http://localhost:9999/api?key=a" &
curl "http://localhost:9999/api?key=a" &
curl "http://localhost:9999/api?key=a" &
curl "http://localhost:9999/api?key=a" &
curl "http://localhost:9999/api?key=g" &
curl "http://localhost:9999/api?key=g" &
curl "http://localhost:9999/api?key=g" &
curl "http://localhost:9999/api?key=g" &
curl "http://localhost:9999/api?key=g" &
curl "http://localhost:9999/api?key=g" &
curl "http://localhost:9999/api?key=g" &
curl "http://localhost:9999/api?key=g" &
curl "http://localhost:9999/api?key=g" &
curl "http://localhost:9999/api?key=g" &
curl "http://localhost:9999/api?key=e" &
curl "http://localhost:9999/api?key=e" &
curl "http://localhost:9999/api?key=e" &
curl "http://localhost:9999/api?key=e" &
curl "http://localhost:9999/api?key=e" &
curl "http://localhost:9999/api?key=e" &
curl "http://localhost:9999/api?key=e" &
curl "http://localhost:9999/api?key=e" &
curl "http://localhost:9999/api?key=e" &
curl "http://localhost:9999/api?key=e" &

wait

