package main

import (
	"MyCache/httpServer"
	single_cache "MyCache/single-cache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	single_cache.NewGroup("scores", single_cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}), 2<<10)

	addr := "localhost:9999"
	peers := httpServer.NewHTTPPool(addr)
	log.Println("MyCache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
