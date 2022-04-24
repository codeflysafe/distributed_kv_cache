/*
 * @Author: sjhuang
 * @Date: 2022-04-23 12:10:58
 * @LastEditTime: 2022-04-24 10:47:53
 * @FilePath: /nosdb/wal/wal_test.go
 */
package wal

import (
	"fmt"
	"nosdb/file"
	"testing"
)

func TestLogger_Append(t *testing.T) {
	var log *Logger
	var err error
	log, err = NewLogger("", DIR_PATH, 0, FILE_MAX_LENGTH, file.STANDARD_IO)
	if err != nil {
		t.Error(err)
	}
	//entry := NewEntry()
	//log.Append()
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		var ttl uint32 = 1000
		entry := NewEntry([]byte(key), []byte(value), PUT, ttl, B_STRING, STRING)
		err = log.Append(entry)
		if err != nil {
			t.Error(err)
		}
	}
}
