/*
 * @Author: sjhuang
 * @Date: 2022-04-24 09:11:53
 * @LastEditTime: 2022-04-24 14:57:21
 * @FilePath: /nosdb/wal/entry_test.go
 */
package wal

import (
	"fmt"
	"testing"
)

func TestEntry_Encode(t *testing.T) {
	key := "key1"
	value := "value1"
	var ttl uint32 = 1000
	entry := NewWalEntry([]byte(key), []byte(value), PUT, ttl, B_STRING, STRING)
	b, err := entry.Encode()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(b)
}
