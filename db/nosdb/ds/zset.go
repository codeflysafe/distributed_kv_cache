package ds

import "nosdb/ds/skiplist"

// 有序集合数据结构
type ZSet struct {
	// key value
	Items map[string][]byte
	// 跳表实现的有序集合
	list *skiplist.SkipList
}
