// 一致性哈希算法
package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Package consistenthash provides an implementation of a ring hash.
// 计算 hash 值
type Hash func(data []byte) uint32

type Map struct {
	// 对应的hash 函数
	hash Hash
	// 虚拟节点
	replicas int
	// 哈希环， int 的正整数范围刚好在 0 ～ 2^31-1 之间，完美
	keys []int
	// hash值与真实节点的映射
	hashMap map[int]string
}

func New(replicas int, hash Hash) *Map {
	m := &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) IsEmpty() bool {
	return len(m.keys) == 0
}

// 传入一些真实节点的名称，同时创建多个虚拟节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// 获取最近的节点名称
// Get gets the closest item in the hash to the provided key.
// 顺时针查找最近的节点名称
func (m *Map) Get(key string) string {
	// if len(m.keys) == 0 {
	// 	panic(" no live machine or cache node")
	// }
	// 计算key的hash值
	hash := int(m.hash(([]byte)(key)))
	// 顺时针查找离它最近的一个节点
	// 在 hash 环上查找比它大的，采用二分查找，找不到使用第一个
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	if idx == len(m.keys) {
		idx = 0
	}
	return m.hashMap[m.keys[idx]]
}
