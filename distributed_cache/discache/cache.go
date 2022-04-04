package discache

import (
	"discache/lru"
	"sync"
)

// 是对 lru.Cache 的包装，然后增加一些安全并发操作
type cache struct {
	// 读写锁
	mu sync.RWMutex
	// 所有的 keys 和 values
	nbytes    int64
	capactity int
	// lru cache
	lru *lru.Cache
	// 统计 命中率 等
	nhit, nget int64
}

func (c *cache) Add(key string, value ByteView) {
	// 写操作，锁住整个缓存
	c.mu.Lock()
	defer c.mu.Unlock()
	// 延迟加载
	if c.lru == nil {
		c.lru = lru.New(c.capactity)
	}
	c.lru.Add(key, value)
	c.nbytes += int64(len(key)) + int64(value.Len())
}

func (c *cache) Get(key string) (value ByteView, ok bool) {
	// 因为涉及到修改list所以才有 lock
	c.nget++
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		c.nhit++
		return v.(ByteView), ok
	}
	return
}

// 获取缓存大
func (c *cache) bytes() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.nbytes
}

func (c *cache) hits() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.nhit
}

func (c *cache) hgets() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.nget
}

func (c *cache) Len() int64 {
	c.mu.RLock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return 0
	}
	return int64(c.lru.Len())
}
