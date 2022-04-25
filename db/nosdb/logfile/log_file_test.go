/*
 * @Author: sjhuang
 * @Date: 2022-04-23 12:10:58
 * @LastEditTime: 2022-04-25 10:25:12
 * @FilePath: /nosdb/logfile/log_file_test.go
 */
package logfile

import (
	"fmt"
	"nosdb/file"
	"testing"
)

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
		entry := NewLogEntry([]byte(key), []byte(value), PUT, ttl, B_STRING, STRING)
		err = log.Append(entry)
		if err != nil {
			t.Error(err)
		}
	}
	fmt.Println(log.offset, log.activeFileName)
}
