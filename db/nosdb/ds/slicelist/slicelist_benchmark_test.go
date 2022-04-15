package slicelist

import (
	"strconv"
	"testing"
)

// BenchmarkSliceList_RPush-4   	   10000	    385774 ns/op
// BenchmarkSliceList_RPush-4   	   10000	    405426 ns/op
// BenchmarkSliceList_RPush-4   	   10000	    402971 ns/op
func BenchmarkSliceList_RPush(b *testing.B) {
	for n := 0; n < b.N; n++ {
		list := NewSliceList()
		for i := 0; i < n; i++ {
			list.RPush([]byte(strconv.Itoa(i)))
		}
	}
}

// BenchmarkLinkedList_RPop-4   	91373572	        13.5 ns/op
func BenchmarkLinkedList_RPop(b *testing.B) {
	for n := 0; n < b.N; n++ {
		list := NewSliceList()
		list.RPop()
	}
}

// BenchmarkLinkedList_RPeek-4   	83354500	        12.4 ns/op
func BenchmarkLinkedList_RPeek(b *testing.B) {
	for n := 0; n < b.N; n++ {
		list := NewSliceList()
		list.RPop()
	}
}

// BenchmarkLinkedList_ListSeek-4   	90736436	        13.0 ns/op
func BenchmarkLinkedList_ListSeek(b *testing.B) {
	for n := 0; n < b.N; n++ {
		list := NewSliceList()
		list.ListSeek(100)
	}
}
