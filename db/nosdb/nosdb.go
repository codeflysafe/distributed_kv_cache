package nosdb

import (
	"sync"
)

type ENCODING uint16

type TYPE uint16

const (
	// list
	LINKED_LIST ENCODING = iota
	SLICE_LIST
	// zset
	SKIPLIST
	// hash
	MHASH
)
const (
	//===========  type =========
	LIST TYPE = iota
	HASH
	SET
	ZSET
	STRING
)

// kv 存储结构
type NosDB struct {
	// 全局的锁
	sync.RWMutex

	// 5 种索引， 每个索引有一个读写锁 和 加载时使用的 sync.Once
	listIdx  *ListIndex
	listOnce sync.Once
	hashIdx  *HashIndex
	hashOnce sync.Once
	setIdx   *SetIndex
	setOnce  sync.Once
	zSetIdx  *ZSetIndex
	zSetOnce sync.Once
	strIdx   *StringIndex
	strOnce  sync.Once

	// todo 缓存淘汰策略
	// lru 策略
}
