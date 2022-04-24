package wal

import (
	"fmt"
	"testing"
)

func TestEntry_Encode(t *testing.T) {
	key := "key1"
	value := "value1"
	var ttl uint32 = 1000
	entry := NewEntry([]byte(key), []byte(value), PUT, ttl, B_STRING, STRING)
	b, err := entry.Encode()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(b)
}
