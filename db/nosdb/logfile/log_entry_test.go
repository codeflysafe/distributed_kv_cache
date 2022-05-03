/*
 * @Author: sjhuang
 * @Date: 2022-04-24 09:11:53
 * @LastEditTime: 2022-04-28 11:07:58
 * @FilePath: /nosdb/logfile/log_entry_test.go
 */
package logfile

import (
	"fmt"
	"strconv"
	"testing"
)

func TestEntry_Encode(t *testing.T) {
	key := "key1"
	value := "value1"
	var ttl uint32 = 1000
	entry := NewLogEntry([]byte(key), nil, []byte(value), -1, PUT, ttl, B_STRING, STRING)
	b, err := entry.Encode()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(b)
	fmt.Println(b[26:28])
}

func TestEntry_Decode(t *testing.T) {
	key := "key1"
	value := "value1"
	var ttl uint32 = 1000
	entry := NewLogEntry([]byte(key), nil, []byte(value), -1, PUT, ttl, B_STRING, STRING)
	b, err := entry.Encode()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(b, len(b))

	if entry, err = DecodeMeta(b[:EntryMetaSize]); err != nil {
		t.Error(err)
	}
	fmt.Println(entry)
	if entry.Encoding != B_STRING || entry.Ty != STRING {
		t.Errorf("decode error ")
	}
}

func TestFormat(t *testing.T) {
	score := 3.14159
	ss := strconv.FormatFloat(score, 'g', -1, 64)
	fmt.Println(ss, len(ss))
}
