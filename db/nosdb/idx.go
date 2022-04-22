package nosdb

import (
	"nosdb/ds"
	"sync"
)

// 采用 分开的 idx 的作用，是为了降低锁的粒度
// 这里定义 list、hash、set、zSet、string 5种 索引
// 但是，也引入了一个问题，那就是针对于 全局 key 的操作， 不太方便
// 现在 key 分散到了 5个 索引中，如何快速的检索key? 可能要遍历这5个索引结构

// ListIndex 采用一个 map
type ListIndex struct {
	*sync.RWMutex
	kv map[string]ds.List
}

func NewListIndex() *ListIndex {
	return &ListIndex{
		RWMutex: new(sync.RWMutex),
		kv:      make(map[string]ds.List, DefaultIdxCap),
	}
}

type HashIndex struct {
	*sync.RWMutex
	kv map[string]ds.Hash
}

func NewHashIndex() *HashIndex {
	return &HashIndex{
		RWMutex: new(sync.RWMutex),
		kv:      make(map[string]ds.Hash, DefaultIdxCap),
	}
}

type SetIndex struct {
	*sync.RWMutex
	kv map[string]ds.Set
}

func NewSetIndex() *SetIndex {
	return &SetIndex{
		RWMutex: new(sync.RWMutex),
		kv:      make(map[string]ds.Set, DefaultIdxCap),
	}
}

type ZSetIndex struct {
	*sync.RWMutex
	kv map[string]ds.ZSet
}

func NewZSetIndex() *ZSetIndex {
	return &ZSetIndex{
		RWMutex: new(sync.RWMutex),
		kv:      make(map[string]ds.ZSet, DefaultIdxCap),
	}
}

type StringIndex struct {
	*sync.RWMutex
	kv map[string]ds.String
}

func NewStringIndex() *StringIndex {
	return &StringIndex{
		RWMutex: new(sync.RWMutex),
		kv:      make(map[string]ds.String, DefaultIdxCap),
	}
}
