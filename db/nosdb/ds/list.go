package ds

import "nosdb/ds/linkedlist"

// list 接口
type List interface {
	LLen() int
	LPush(value []byte)
	LPop() (value []byte)
	RPush(value []byte)
	RPop() (value []byte)
	ListSeek(idx int) (value []byte, err error)
	ListSet(idx int, newVal []byte) (err error)
	ListDelIndex(idx int)
	LPeek() (value []byte)
	RPeek() (value []byte)
	Empty() bool
}

// 默认采用 linkedlist
func NewList() List {
	return linkedlist.NewLinkedList()
}
