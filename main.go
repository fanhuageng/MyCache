package main

import (
	"MyCache/httpServer"
	single_cache "MyCache/single-cache"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"a": "630",
	"b": "589",
	"c": "567",
	"d": "510",
	"e": "633",
	"f": "480",
	"g": "350",
	"h": "430",
}

func createGroup() *single_cache.Group {
	return single_cache.NewGroup("scores", single_cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}), 2<<10)
}

// 启动缓存服务器
func startCacheServer(addr string, addrs []string, Fcache *single_cache.Group) {
	peers := httpServer.NewHTTPPool(addr) // 创建HTTPPool
	peers.Set(addrs...)                   // 添加节点
	Fcache.RegisterPeers(peers)
	log.Printf("MyCache running at: %s", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

// 启动API服务器，与用户交互
func startAPIServer(apiAddr string, Fcache *single_cache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := Fcache.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("MyCache api server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool

	flag.IntVar(&port, "port", 8001, "MyCache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	Fcache := createGroup()
	if api {
		go startAPIServer(apiAddr, Fcache)
	}
	startCacheServer(addrMap[port], []string(addrs), Fcache)
}
