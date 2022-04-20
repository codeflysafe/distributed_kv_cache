package nosdb

import (
	"nosdb/ds"
	"sync"
)

// 采用 分开的 idx 的作用，是为了降低锁的粒度
// 这里定义 list、hash、set、zset、string 5种 索引
// 但是，也引入了一个问题，那就是针对于 全局 key 的操作， 不太方便
// 现在 key 分散到了 5个 索引中，如何快速的检索key可能要遍历这5个索引结构

// 采用一个 map
type ListIndex struct {
	sync.RWMutex
	kv map[string]ds.List
}

func NewListIndex() *ListIndex {
	return &ListIndex{
		RWMutex: sync.RWMutex{},
		kv:      make(map[string]ds.List),
	}
}

type HashIndex struct {
	sync.RWMutex
	kv map[string]ds.Hash
}
type SetIndex struct {
	sync.RWMutex
	kv map[string]ds.Set
}

type ZSetIndex struct {
	sync.RWMutex
	kv map[string]ds.ZSet
}

type StringIndex struct {
	sync.RWMutex
	kv map[string]ds.String
}
