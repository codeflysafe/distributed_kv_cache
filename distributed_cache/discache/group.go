package discache

import (
	"discache/singleflight"
	"fmt"
	"log"
	"sync"
)

// 定义接口 Getter 和 回调函数 Get(key string)([]byte, error)，参数是 key，返回值是 []byte
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
	// loader
	loader *singleflight.Group
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// 从 cache 中加载 缓存
// 获取不了就从本地加载后，放入缓存
func (g *Group) Get(key string) (value ByteView, err error) {
	if len(key) == 0 {
		err = fmt.Errorf("key can't be blank")
		return
	}

	if v, ok := g.mainCache.Get(key); ok {
		log.Println("[Cache] hit")
		return v, nil
	}

	// 加载不成功，则从本地加载
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {

	// 保证同时访问同一个缓存key时，只加载一次
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err := g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[Cache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})
	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

//
func (g *Group) getLocally(key string) (value ByteView, err error) {
	bytes, e := g.getter.Get(key)
	if e != nil {
		err = e
		return
	}

	value = ByteView{b: bytes}
	g.put(key, value)
	return value, nil
}

// 向里面增加缓存
func (g *Group) put(key string, value ByteView) {
	g.mainCache.Add(key, value)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, capacity int, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{capactity: capacity},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// return nil if not exists group named name
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}
