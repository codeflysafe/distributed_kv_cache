package wal

import (
	"errors"
	"fmt"
	"nosdb/file"
	"nosdb/snowflake"
	"nosdb/utils"
	"sync"
)

const (
	ACTIVE_FILE_PREFIX = "active_"
	WAL_FILE_PREFIX    = "wal_"
	DIR_PATH           = "./log"
	CHECKPOINT_MAX_LEN = 64
	FILE_MAX_LENGTH    = 1 << 20
)

var node, _ = snowflake.NewNode(1)

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

// OpenLogger 从 active_  中恢复文件
// 如果不存在，则新建
// 新建时，要传入指定的文件操作模式 mod
// func OpenLogger(dirPath string, mod file.MOD) (log *Logger, err error) {
// 	// 从
// }

// 新建一个 logger
func NewLogger(activeFileName, dirPath string, maxFileSize int64, mod file.MOD) (log *Logger, err error) {
	log = &Logger{
		Option: Option{
			dirPath:     dirPath,
			maxFileSize: maxFileSize,
		},
		mod:            mod,
		offset:         0,
		activeFileName: activeFileName,
	}
	if log.maxFileSize > FILE_MAX_LENGTH || log.maxFileSize <= 0 {
		log.maxFileSize = FILE_MAX_LENGTH
	}

	if node == nil {
		err = errors.New(" generator node error ")
		return
	}
	log.seq = node.Generate().Int64()
	// activeFileName 不存在
	if len(activeFileName) == 0 {
		// 首先去 该文件下面找到一个 active_ 文件
		if log.activeFileName, err = utils.PrefixPath(dirPath, ACTIVE_FILE_PREFIX); err != nil {
			// 如果不存在这样的文件，就新建一个
			log.activeFileName = fmt.Sprintf("%s%d.dat", ACTIVE_FILE_PREFIX, log.seq)
		}
	}
	log.activeFileHandle, err = file.OpenFile(mod, dirPath, log.activeFileName, int(log.maxFileSize))
	if err != nil {
		return
	}
	log.offset, err = log.activeFileHandle.Offset()
	return
}

// Append 日志文件中，追加写
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
		newActiveFilePath := fmt.Sprintf("%s_%d.dat", ACTIVE_FILE_PREFIX, l.seq)
		l.activeFileHandle, err = file.OpenFile(l.mod, l.dirPath, newActiveFilePath, int(l.maxFileSize))
	} else {
		// 还有空间，写入
		l.offset, err = l.activeFileHandle.WriteAt(l.offset, b)
	}
	return
}

// 将缓冲区数据刷入文件中
func (l *Logger) Flush() error {
	return l.activeFileHandle.Sync()
}

// 关闭logger操作，
func (l *Logger) Close() error {
	l.Lock()
	defer l.Unlock()
	return l.activeFileHandle.Close()
}
