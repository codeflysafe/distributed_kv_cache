// lru 缓存管理模块
package lru

import "container/list"

type Cache struct {
	// cache 的最大容量
	// 0 代表无限
	capactity int
	// 存储 key/value
	l *list.List
	// 存储 cache
	cache map[string]*list.Element

	// 回掉函数，当节点被删除时后的回掉函数，可以为空
	OnEvicted func(key string, value interface{})
}

// type Key interface{}

type entry struct {
	key   string
	value interface{}
}

// 创建一个新的 cahce
func New(capactity int) *Cache {
	return &Cache{
		capactity: capactity,
		l:         list.New(),
		cache:     make(map[string]*list.Element),
	}
}

// 从缓存中加载 key
func (c *Cache) Get(key string) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		// 移到对头
		c.l.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

// 淘汰最久未使用的一个
func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.l.Back()
	if ele != nil {
		c.Remove(ele.Value.(*entry).key)
	}
}

// 新增一个新的缓存
func (c *Cache) Add(key string, value interface{}) {
	if c.cache == nil {
		c.cache = make(map[string]*list.Element)
		c.l = list.New()
	}

	if ele, hit := c.cache[key]; hit {
		c.l.MoveToFront(ele)
		ele.Value.(*entry).value = value
		return
	}

	ele := c.l.PushFront(&entry{key, value})
	c.cache[key] = ele
	if c.capactity != 0 && c.l.Len() > c.capactity {
		c.RemoveOldest()
	}
}

// func (c (Cache)) removeEelemet(key string) *list.Element{
// 	if(c.cache == nil){
// 		return
// 	}
// }

// 删除一个缓存
func (c *Cache) Remove(key string) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		delete(c.cache, key)
		c.l.Remove(ele)
		if c.OnEvicted != nil {
			kv := ele.Value.(*entry)
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.l.Len()
}

func (c *Cache) Clear() {
	if c.OnEvicted != nil {
		for _, e := range c.cache {
			kv := e.Value.(*entry)
			c.OnEvicted(kv.key, kv.value)
		}
	}
	c.l = nil
	c.cache = nil
}
