package logfile

import (
	"fmt"
	"nosdb/file"
	"nosdb/snowflake"
	"nosdb/utils"
	"strconv"
	"strings"
	"sync"
)

const (
	ACTIVE_FILE_PREFIX = "active_"
	WAL_FILE_PREFIX    = "wal_"
	DIR_PATH           = "./log"
	CHECKPOINT_MAX_LEN = 64
	FILE_MAX_LENGTH    = 1 << 20
)

// 文件的一些 配置信息
type Option struct {
	dirPath     string // 文件夹
	maxFileSize int64  // 文件的最大大小
}

// 使用 write ahead log 的方法，来实现数据库的原子性和持久性操作
type LogFile struct {
	sync.RWMutex                     // 读写锁，并发控制
	Option                           // 文件的一些信息
	seq              int64           // 文件的序号
	offset           int64           // 当前偏移量
	activeFileHandle file.FileHandle // 当前的文件管理模块
	activeFileName   string          // 当前的文件名称
	mod              file.MOD        // 文件的读取模式
	Node             *snowflake.Node // 用于生成 seq 的 node
}

// 从上下文恢复
func ReOpenLogFile(dirPath string, maxFileSize int64, mod file.MOD) (log *LogFile, err error) {
	log = &LogFile{
		RWMutex: sync.RWMutex{},
		Option:  Option{dirPath: dirPath, maxFileSize: maxFileSize},
		mod:     mod,
	}
	if log.maxFileSize > FILE_MAX_LENGTH || log.maxFileSize <= 0 {
		log.maxFileSize = FILE_MAX_LENGTH
	}

	// 首先去 该文件下面找到一个 active_ 文件
	if log.activeFileName, err = utils.PrefixPath(dirPath, ACTIVE_FILE_PREFIX); err != nil {
		return
		// 如果不存在这样的文件，就新建一个
		// log.activeFileName = fmt.Sprintf("%s%d.dat", ACTIVE_FILE_PREFIX, log.seq)
	}
	// 去获取 文件的 seq 序号
	var seq int
	if seq, err = strconv.Atoi(strings.Split(strings.Split(log.activeFileName, "_")[1], ".")[0]); err != nil {
		return
	}
	log.seq = int64(seq)
	// 去打开新的 fileHandle， 使用指定的模式
	if log.activeFileHandle, err = file.OpenFile(mod, dirPath, log.activeFileName, int(log.maxFileSize)); err != nil {
		return
	}
	// 去获得文件的当前偏移量
	log.offset, err = log.activeFileHandle.Offset()
	return
}

// 新建一个 logger
func NewLogFile(dirPath string, maxFileSize int64, mod file.MOD) (log *LogFile, err error) {
	log = &LogFile{
		RWMutex:          sync.RWMutex{},
		Option:           Option{dirPath: dirPath, maxFileSize: maxFileSize},
		seq:              0,
		offset:           0,
		activeFileHandle: nil,
		mod:              mod,
	}
	if log.maxFileSize > FILE_MAX_LENGTH || log.maxFileSize <= 0 {
		log.maxFileSize = FILE_MAX_LENGTH
	}
	// activeFileName 新建一个
	log.activeFileName = fmt.Sprintf("%s%d.dat", ACTIVE_FILE_PREFIX, log.seq)
	log.activeFileHandle, err = file.OpenFile(mod, dirPath, log.activeFileName, int(log.maxFileSize))
	if err != nil {
		return
	}
	log.offset, err = log.activeFileHandle.Offset()
	return
}

// Append 日志文件中，追加写
// 如何文件已经满了，则会写入一个新的文件
func (l *LogFile) Append(entry *LogEntry) (err error) {
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
		// 给一个新的seq
		l.seq = l.Node.Generate().Int64()
		// 新的 数据库文件 名称 active_{seq}.dat
		newActiveFilePath := fmt.Sprintf("%s_%d.dat", ACTIVE_FILE_PREFIX, l.seq)
		l.activeFileHandle, err = file.OpenFile(l.mod, l.dirPath, newActiveFilePath, int(l.maxFileSize))
	} else {
		// 还有空间，写入
		l.offset, err = l.activeFileHandle.WriteAt(l.offset, b)
	}
	return
}

// 从 offset 处读出 entry
// 包含两次文件读
func (l *LogFile) ReadAt(offset int64) (e *LogEntry, err error) {
	// 读锁
	l.RLock()
	defer l.RUnlock()
	var b []byte
	// 1. 先从文件中读出对应的 offset
	if b, err = l.activeFileHandle.ReadAt(offset, EntryMetaSize); err != nil {
		return
	}
	// 2.解析 读出的 []byte
	if e, err = DecodeMeta(b); err != nil {
		return
	}
	offset += EntryMetaSize
	// 3. 在从文件中读出， key 和 value
	if e.Key, err = l.activeFileHandle.ReadAt(offset, int(e.KeySize)); err != nil {
		return
	}
	offset += int64(e.KeySize)
	e.Value, err = l.activeFileHandle.ReadAt(offset, int(e.KeySize))
	return
}

// 将缓冲区数据刷入文件中
func (l *LogFile) Flush() error {
	return l.activeFileHandle.Sync()
}

// 关闭logger操作，
func (l *LogFile) Close() error {
	l.Lock()
	defer l.Unlock()
	return l.activeFileHandle.Close()
}
