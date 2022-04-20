package nosdb

import (
	"sync"
	"unsafe"
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

// NosObj 用来存储 value
type NosObj struct {
	// 底层数据结构
	encoding ENCODING
	// 存储类型
	ty TYPE
	// 指向底层数据结构的指针
	ptr unsafe.Pointer
}

// k/v 存储结构
type NosDB struct {
	sync.RWMutex
	kv map[string]NosObj

	// todo 缓存淘汰策略
	// lru 策略
}
