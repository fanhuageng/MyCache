package single_cache

import (
	"MyCache/distributedNode"
	"MyCache/singleFlight"
	"fmt"
	"log"
	"sync"
)

type Group struct {
	name      string
	getter    Getter // 缓存未命中时，获得源数据的回调
	mainCache cache  // 一开始实现的并发缓存
	peers     distributedNode.PeerPicker
	loader    *singleFlight.SFGroup
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, getter Getter, cacheBytes int64) *Group {
	if getter == nil {
		panic("Getter nil")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleFlight.SFGroup{},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

// Get实现Getter接口函数
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 注册选择远程节点。将实现了 PeerPicker 接口的 HTTPPool 注入到 Group 中
func (g *Group) RegisterPeers(peers distributedNode.PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	// 缓存命中，从缓存中查找
	if v, ok := g.mainCache.get(key); ok {
		log.Printf("cache hit")
		return v, nil
	}
	return g.load(key) // 未命中，本地查找或远程查找
}

func (g *Group) load(key string) (value ByteView, err error) {
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err := g.loadRemote(peer, key); err == nil {
					return value, nil
				}
				log.Printf("Failed to load from remote peer %s", peer)
			}
		}
		return g.loadLocally(key)
	})

	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

// 本地查找
func (g *Group) loadLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// 远程查找
func (g *Group) loadRemote(peergetter distributedNode.PeerGetter, key string) (ByteView, error) {
	bytes, error := peergetter.Get(g.name, key) // 从远程节点中查询最终值
	if error != nil {
		return ByteView{}, error
	}
	return ByteView{b: bytes}, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
