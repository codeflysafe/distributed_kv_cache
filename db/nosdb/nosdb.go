/*
 * @Author: sjhuang
 * @Date: 2022-04-15 20:41:59
 * @LastEditTime: 2022-05-09 16:36:01
 * @FilePath: /nosdb/nosdb.go
 */
package nosdb

import (
	"nosdb/file"
	"nosdb/logfile"
	"sync"
)

// kv 存储结构
type NosDB struct {
	// 全局的锁
	sync.RWMutex

	// 5 种索引， 每个索引有一个读写锁
	// 和 加载时使用的 sync.Once
	listIdx  *ListIndex
	listOnce sync.Once
	hashIdx  *HashIndex
	hashOnce sync.Once
	setIdx   *SetIndex
	setOnce  sync.Once
	zSetIdx  *ZSetIndex
	zSetOnce sync.Once
	strIdx   *StringIndex
	strOnce  sync.Once

	// wal 日志文件模块
	logFile *logfile.LogFile
	mod     file.MOD // 操作日志文件的mod
}

func NewNosDB() *NosDB {
	// var err error
	// db := new(NosDB)
	// db.walLogger, err = logfile.NewWalLogger()
	return nil
}

// 写入 log 文件
func (db *NosDB) WriteLog(key, value []byte) {
	db.logFile.Append(logfile.NewLogEntry([]byte(key),
		nil, value, -1, SET, logfile.Persistent, logfile.B_STRING, logfile.STRING))
	// db.logFile.Append()
}
