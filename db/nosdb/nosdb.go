/*
 * @Author: sjhuang
 * @Date: 2022-04-15 20:41:59
 * @LastEditTime: 2022-04-28 11:15:32
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

	// todo 缓存淘汰策略
	// lru 策略

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


