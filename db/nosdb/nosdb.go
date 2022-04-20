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

	// 5 种索引， 每个索引有一个读写锁
	listIdx *ListIndex
	hashIdx *HashIndex
	setIdx  *SetIndex
	zSetIdx *ZSetIndex
	strIdx  *StringIndex

	// todo 缓存淘汰策略
	// lru 策略
}
