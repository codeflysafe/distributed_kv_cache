package nosdb

import "unsafe"

type ENCODING uint8

type TYPE uint8

const (
	// list
	LINKED_LIST ENCODING = iota
	SLICE_LIST
	// zset
	SKIPLIST
	// map
	MAP
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

	// todo 缓存淘汰策略
	// lru 策略

}
