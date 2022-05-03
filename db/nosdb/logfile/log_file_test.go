/*
 * @Author: sjhuang
 * @Date: 2022-04-23 12:10:58
 * @LastEditTime: 2022-04-28 11:04:25
 * @FilePath: /nosdb/logfile/log_file_test.go
 */
package logfile

import (
	"fmt"
	"nosdb/file"
	"testing"
)

var PUT CMD = 0

func TestLogger_Append(t *testing.T) {
	var log *LogFile
	var err error
	log, err = NewLogFile(DIR_PATH, FILE_MAX_LENGTH, file.STANDARD_IO)
	if err != nil {
		t.Error(err)
	}
	defer log.Close()
	fmt.Println(log.offset, log.activeFileName)
	//entry := NewEntry()
	//log.Append()
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		var ttl uint32 = 1000
		entry := NewLogEntry([]byte(key), nil, []byte(value), -1, PUT, ttl, B_STRING, STRING)
		err = log.Append(entry)
		if err != nil {
			t.Error(err)
		}
	}
	fmt.Println(log.offset, log.activeFileName)
}

func TestLogger_ReadAt(t *testing.T) {
	var log *LogFile
	var err error
	log, err = ReOpenLogFile(DIR_PATH, FILE_MAX_LENGTH, file.STANDARD_IO)
	if err != nil {
		t.Error(err)
	}
	defer log.Close()
	fmt.Println(log.offset, log.activeFileName)
	//entry := NewEntry()
	//log.Append()
	var ttl uint32 = 1000
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		entry := NewLogEntry([]byte(key), nil, []byte(value), -1, PUT, ttl, B_STRING, STRING)
		err = log.Append(entry)
		if err != nil {
			t.Error(err)
		}
	}
	fmt.Println(log.offset, log.activeFileName)
}
