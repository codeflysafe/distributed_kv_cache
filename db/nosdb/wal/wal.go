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
func NewLogger(activeFileName, dirPath string, offset, maxFileSize int64, mod file.MOD) (log *Logger, err error) {
	log = &Logger{
		Option: Option{
			dirPath:     dirPath,
			maxFileSize: maxFileSize,
		},
		offset:         offset,
		mod:            mod,
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
	// activeFileName 不存在， 新建一个
	if len(activeFileName) == 0 {
		log.activeFileName = fmt.Sprintf("%s%d.dat", ACTIVE_FILE_PREFIX, log.seq)
	}
	log.activeFileHandle, err = file.OpenFile(mod, dirPath, log.activeFileName, int(log.maxFileSize))
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

// 关闭logger操作，
// 这里会释放一些资源，包括文件上下文备份等等
// todo
func (l *Logger) Close() error {
	return l.activeFileHandle.Close()
}

type FileCheckPoint struct {
	fileName string   //
	offset   int64    // 8
	mod      file.MOD // 2
}
