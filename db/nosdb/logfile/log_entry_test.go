/*
 * @Author: sjhuang
 * @Date: 2022-04-24 09:11:53
 * @LastEditTime: 2022-04-25 10:24:05
 * @FilePath: /nosdb/logfile/log_entry_test.go
 */
package logfile

import (
	"fmt"
	"testing"
)

func TestEntry_Encode(t *testing.T) {
	key := "key1"
	value := "value1"
	var ttl uint32 = 1000
	entry := NewLogEntry([]byte(key), []byte(value), PUT, ttl, B_STRING, STRING)
	b, err := entry.Encode()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(b)
}
