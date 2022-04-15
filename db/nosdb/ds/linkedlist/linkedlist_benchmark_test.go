package linkedlist

import (
	"strconv"
	"testing"
)

// 测试插入
// BenchmarkLinkedList_RPush-4   	   10000	    522219 ns/op
// BenchmarkLinkedList_RPush-4   	   10000	    519925 ns/op
// BenchmarkLinkedList_RPush-4   	   10000	    526954 ns/op
func BenchmarkLinkedList_RPush(b *testing.B) {
	for n := 0; n < b.N; n++ {
		list := NewLinkedList()
		for i := 0; i < n; i++ {
			list.RPush([]byte(strconv.Itoa(i)))
		}
	}
}

// BenchmarkLinkedList_RPop-4   	1000000000	         0.370 ns/op
func BenchmarkLinkedList_RPop(b *testing.B) {
	for n := 0; n < b.N; n++ {
		list := NewLinkedList()
		list.RPop()
	}
}

// BenchmarkLinkedList_RPeek-4   	1000000000	         0.324 ns/op
func BenchmarkLinkedList_RPeek(b *testing.B) {
	for n := 0; n < b.N; n++ {
		list := NewLinkedList()
		list.RPop()
	}
}

// BenchmarkLinkedList_ListSeek-4   	347112512	         3.22 ns/op
func BenchmarkLinkedList_ListSeek(b *testing.B) {
	for n := 0; n < b.N; n++ {
		list := NewLinkedList()
		list.ListSeek(100)
	}
}
