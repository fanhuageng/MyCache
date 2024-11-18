package httpServer

import (
	pb "MyCache/cacheProtoBuf/mycachepb"
	"MyCache/consistentHash"
	"MyCache/distributedNode"
	"MyCache/single-cache"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const defaultBasePath = "/fcache/"
const defaultMultiple = 50

type HTTPPool struct {
	self        string
	basePath    string
	mu          sync.Mutex
	peers       *consistentHash.Map // 一致性哈希
	httpClients map[string]*httpClient
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// 实例化一致性哈希算法，并添加节点
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistentHash.New(defaultMultiple, nil)
	p.peers.Add(peers...)
	p.httpClients = make(map[string]*httpClient, len(peers))
	// 为每一个节点创建HTTP客户端
	for _, peer := range peers {
		p.httpClients[peer] = &httpClient{baseURL: peer + p.basePath}
	}
}

// 根据具体的 key，选择节点，返回节点对应的 HTTP 客户端
func (p *HTTPPool) PickPeer(key string) (distributedNode.PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	peer := p.peers.Get(key) // 通过key找到对应的真实节点
	if peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpClients[peer], true
	}
	return nil, false
}

var _ distributedNode.PeerPicker = (*HTTPPool)(nil)

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool servering unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupname := parts[0]
	key := parts[1]

	group := single_cache.GetGroup(groupname)
	if group == nil {
		http.Error(w, "no such group"+groupname, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

type httpClient struct {
	baseURL string // 要访问的远程节点的地址
}

func (h *httpClient) Get(in *pb.Request, out *pb.Response) error {
	// 下面的代码中%v/%v/%v是错误的，这里调试时发现u会变为http://localhost:8003/fcache//scores导致返回错误404，找不到客户端
	u := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(in.GetGroup()), url.QueryEscape(in.GetKey()))
	response, err := http.Get(u)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", response.Status)
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}

var _ distributedNode.PeerGetter = (*httpClient)(nil)
