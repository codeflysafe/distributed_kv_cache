package bitcaskdb

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

const (
	Ratio = 2
)

// 数据库的定义
type BitcaskDB struct {
	io.Closer
	indexes  map[string]int64 // key 与 offset 的索引
	indexes2 map[string]int64 // merge 过程中使用的hash 表
	dbFile   *DBFile          // 数据文件
	dirPath  string           // 数据位置
	mu       sync.RWMutex     // 采用读写锁
	ratio    int64            // indexes2， 新建时的大小
	count    int64
}

// 恢复数据，从磁盘中家在数据
func (db *BitcaskDB) restore() error {
	// 为空无需加载
	if db.dbFile.OffSet == 0 {
		return nil
	}
	var offset int64 = 0
	for offset < db.dbFile.OffSet {
		entry, err := db.dbFile.Read(offset)
		if err != nil {
			return err
		}
		// 遇到删除，就删除掉这个字段
		if entry.Mark == DEL {
			delete(db.indexes, string(entry.Key))
		} else {
			// 更新索引
			db.indexes[string(entry.Key)] = offset
		}
		offset += entry.GetSize()
	}
	db.count = int64(len(db.indexes))
	return nil
}

// Open 开启一个数据库实例
func Open(dirPath string) (db *BitcaskDB, err error) {
	db = &BitcaskDB{
		dirPath: dirPath,
		indexes: make(map[string]int64, 10),
		ratio:   Ratio,
	}
	// 1. 新建db file
	db.dbFile, err = NewDBFile(dirPath)
	if err != nil {
		return
	}
	// 2. 从文件中加载到内存
	err = db.restore()
	return
}

// 放入
func (db *BitcaskDB) Put(key, value []byte) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	offset, err := db.dbFile.Write(NewEntry(key, value, PUT))
	if err != nil {
		return
	}
	db.indexes[string(key)] = offset
	return
}

// 查找
func (db *BitcaskDB) Get(key []byte) (value []byte, err error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var entry *Entry
	if offset, ok := db.indexes[string(key)]; ok {
		entry, err = db.dbFile.Read(offset)
		if err != nil {
			return
		}
		value = entry.Value
		return
	}
	err = fmt.Errorf("no such key %v", string(key))
	return
}

func (db *BitcaskDB) Del(key []byte) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.indexes[string(key)]; ok {
		// 这里无需判断 value 是否相同，因此只需要插入 删除 key 字段即可
		_, err = db.dbFile.Write(NewEntry(key, nil, DEL))
		if err != nil {
			// 删除失败
			return
		}
		// 删除改 key 的元素
		delete(db.indexes, string(key))
	}

	// 如果不存在，无需删除
	return nil
}

func (db *BitcaskDB) scanFiles(ch chan<- *Entry) error {

	// 关闭 ch
	defer close(ch)
	var offset int64 = 0
	for offset < db.dbFile.OffSet {
		entry, err := db.dbFile.Read(offset)
		if err != nil {
			return err
		}
		// 是最新的并且合法的
		if off, ok := db.indexes[string(entry.Key)]; ok && off == offset {
			// fmt.Println("scanFiles", string(entry.Key), string(entry.Value))
			ch <- entry
		}
		// 更新 offset
		offset += entry.GetSize()
	}
	return nil
}

func (db *BitcaskDB) Merge() error {
	//db.mu.Lock()
	//defer db.mu.Unlock()
	// 新设一个chan
	ch := make(chan *Entry, 10)
	go db.scanFiles(ch)
	// 新建一个merge 文件
	mergeDbFile, err := NewMergeDBFile(db.dirPath)
	if err != nil {
		return err
	}
	// 新建 indexes2
	db.indexes2 = make(map[string]int64, db.dbFile.OffSet/db.ratio+1)
	// 从内存中写入到 merge file 中
	// 这里采用一个 channel 来异步更新
	var offset int64
	for entry := range ch {
		offset, err = mergeDbFile.Write(entry)
		fmt.Println("merge", string(entry.Key), string(entry.Value), offset)
		if err != nil {
			log.Println("err: ", err)
		} else {
			// 更新 indexes2
			db.indexes2[string(entry.Key)] = offset
		}
	}

	// 更新 indexes 和 文件
	db.indexes = db.indexes2
	db.indexes2 = nil

	dbFileName := db.dbFile.File.Name()
	db.dbFile.File.Close()
	os.Remove(dbFileName)

	mergeDbFileName := mergeDbFile.File.Name()
	mergeDbFile.File.Close()
	os.Rename(mergeDbFileName, db.dirPath+string(os.PathSeparator)+FileNmae)
	db.dbFile, err = NewDBFile(db.dirPath)
	return err

}

func (db *BitcaskDB) ListKeys() (keys [][]byte, err error) {
	db.mu.RLock()
	db.mu.RUnlock()
	keys = make([][]byte, db.count)
	for key := range db.indexes {
		keys = append(keys, []byte(key))
	}
	return
}

func (db *BitcaskDB) Close() error {
	if db.dbFile.File != nil {
		return db.dbFile.File.Close()
	}
	return nil
}

func (db *BitcaskDB) Sync() error {
	if db.dbFile == nil {
		return fmt.Errorf("BitcaskDB not init")
	}
	return db.dbFile.File.Sync()
}
