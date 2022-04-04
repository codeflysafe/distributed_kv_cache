package discache

import (
	"discache/consistenthash"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_cache"
	// 默认50个虚拟节点
	defaultReplicas = 50
)

type HTTPPool struct {
	//一个是 self，用来记录自己的地址，包括主机名/IP 和端口。
	self string
	// basePath，作为节点间通讯地址的前缀，默认是 /_cache/，
	basePath string
	// 读写锁
	mu sync.RWMutex
	// 一致性hash算法
	peers *consistenthash.Map
	// 每一个 节点对应的 httpGetter
	// baseurl 127.0.0.1:8000等
	httpGetters map[string]*httpGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// 记录日志
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// 实现 handle 接口
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("ServeHTTP: %s, %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath)+1:], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupNmae := parts[0]
	key := parts[1]

	group := GetGroup(groupNmae)
	if group == nil {
		http.Error(w, "no such group "+groupNmae, http.StatusBadRequest)
		return
	}
	value, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(value.ByteSlice())
}

// 通过 缓存的key来获取对应的 PeerGetter
func (p *HTTPPool) PickPeer(key string) (peer PeerGetter, ok bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

// 增加一系列 节点
// peers : http://127.0.0.1
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// ----  httpGetter ----
type httpGetter struct {
	baseURL string
}

// @see PeerGetter
// 实现PeerGetter接口
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	// log.Printf("Get baseurl:%s group:%s, key: %s", h.baseURL, group, key)
	u := fmt.Sprintf(
		"%v/%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	log.Printf("Get: %s", u)
	// 从服务器获取对应的 resp
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", resp.Status)
	}

	// resp 中获取 bytes
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

var _PeerGetter = (*httpGetter)(nil)
