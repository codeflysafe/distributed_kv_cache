/*
 * @Author: sjhuang
 * @Date: 2022-04-15 20:41:59
 * @LastEditTime: 2022-05-09 17:25:10
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
	// 和 懒加载时使用的 sync.Once
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

func NewNosDB() (db *NosDB, err error) {
	db = &NosDB{
		RWMutex: sync.RWMutex{},
	}
	// db.logFile = logfile.ReOpenLogFile()
	// return db
	return
}

func (db *NosDB) writeKVLog(key string, value []byte, cmd logfile.CMD, ty logfile.TYPE) {
	db.writeLog([]byte(key), nil, value, -1, cmd, ty)
}

// 写入 log 文件
func (db *NosDB) writeLog(key, member, value []byte, score float64, cmd logfile.CMD, ty logfile.TYPE) {
	db.Lock()
	defer db.Unlock()
	switch ty {
	case logfile.STRING:
		db.logFile.Append(logfile.NewLogEntry([]byte(key),
			nil, value, -1, cmd, logfile.Persistent, logfile.B_STRING, logfile.STRING))
	case logfile.HASH:
		db.logFile.Append(logfile.NewLogEntry([]byte(key),
			member, value, -1, cmd, logfile.Persistent, logfile.M_HASH, logfile.HASH))
	case logfile.LIST:
		db.logFile.Append(logfile.NewLogEntry([]byte(key),
			nil, value, -1, cmd, logfile.Persistent, logfile.LINKED_LIST, logfile.LIST))
	case logfile.SET:
		db.logFile.Append(logfile.NewLogEntry([]byte(key),
			nil, value, -1, cmd, logfile.Persistent, logfile.M_HASH, logfile.HASH))
	default:
		db.logFile.Append(logfile.NewLogEntry([]byte(key),
			nil, value, score, cmd, logfile.Persistent, logfile.SKIP_LIST, logfile.ZSET))

	}

}
