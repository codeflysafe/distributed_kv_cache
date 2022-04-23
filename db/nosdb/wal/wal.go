package wal

import (
	"fmt"
	"nosdb/file"
	"nosdb/utils"
	"sync"
)

const (
	ACTIVE_FILE_PREFIX = "active_"
	WAL_FILE_PREFIX    = "wal_"
	DIR_PATH           = "log/"
)

// 文件的一些 配置信息
type Option struct {
	dirPath     string // 文件夹
	maxFileSize int64  // 文件的最大大小
}

// 使用 write ahead log 的方法，来实现数据库的原子性和持久性操作
type Logger struct {
	sync.RWMutex                     // 读写锁，并发控制
	Option                           // 文件的一些信息
	seq              int64           // 文件的序号
	offset           int64           // 当前偏移量
	activeFileHandle file.FileHandle // 当前的文件管理模块
	activeFileName   string          // 当前的文件名称
	mod              file.MOD        // 文件的读取模式
}

// 日志文件中，追加写
// 如何文件已经满了，则会写入一个新的文件
func (l *Logger) Append(entry *Entry) (err error) {

	var b []byte
	b, err = entry.Encode()
	if err != nil {
		return
	}
	// 锁住，并发控制
	l.Lock()
	defer l.Unlock()
	// 如果当前文件
	if l.offset+int64(len(b)) > l.maxFileSize {
		// 要新建一个新的wal文件，然后重命名旧的wal文件
		// 命名是 wal_{seq}.dat
		newFilePath := fmt.Sprintf("%s_%d.dat", WAL_FILE_PREFIX, l.seq)
		err = utils.ReNameFile(l.dirPath, l.activeFileName, newFilePath)
		if err != nil {
			return
		}
		// 就文件关闭
		err = l.activeFileHandle.Close()
		if err != nil {
			return
		}
		// 新的 数据库文件 名称 active_{seq}.dat
		newActiveFilePath := fmt.Sprintf("active_%d.dat", ACTIVE_FILE_PREFIX, l.seq)
		l.activeFileHandle, err = file.OpenFile(l.mod, l.dirPath, newActiveFilePath, int(l.maxFileSize))
	} else {
		// 还有空间，写入
		l.offset, err = l.activeFileHandle.WriteAt(l.offset, b)
	}
	return
}
