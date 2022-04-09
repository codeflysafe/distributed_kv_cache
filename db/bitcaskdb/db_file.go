package bitcaskdb

import (
	"fmt"
	"log"
	"os"
)

const (
	FileNmae      = "minidb.data"
	MergeFileName = "minidb.db.merge"
)

type DBFile struct {
	File   *os.File // 存储的文件
	OffSet int64    // 当前的偏移量
}

func newInternal(fileName string) (*DBFile, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Printf(" open file error: %v", err)
		return nil, err
	}
	// 先恢复数据
	stat, err := os.Stat(fileName)
	if err != nil {
		log.Printf(" Stat file error: %v", err)
		return nil, err
	}

	return &DBFile{OffSet: stat.Size(), File: file}, nil
}

func NewDBFile(path string) (*DBFile, error) {
	fileName := path + string(os.PathSeparator) + FileNmae
	return newInternal(fileName)
}

func NewMergeDBFile(path string) (*DBFile, error) {
	mergeFileName := path + string(os.PathSeparator) + MergeFileName
	return newInternal(mergeFileName)
}

// 从 offset 处读取数据
func (f *DBFile) Read(offset int64) (e *Entry, err error) {

	if offset > f.OffSet {
		err = fmt.Errorf("Bad offset, offset can't more than %d, but is %d\n", f.OffSet, offset)
		return
	}
	// 读取固定长度
	buf := make([]byte, entryHeaderSize)
	if _, err = f.File.ReadAt(buf, offset); err != nil {
		return
	}

	if e, err = Decode(buf); err != nil {
		return
	}

	// 读取 key
	offset += entryHeaderSize
	if e.KeySize > 0 {
		key := make([]byte, e.KeySize)
		if _, err = f.File.ReadAt(key, offset); err != nil {
			return
		}
		e.Key = key
	}

	// 读取 value
	// ReadAt always returns a non-nil error when n < len(b).
	offset += int64(e.KeySize)
	if e.ValueSize > 0 {
		value := make([]byte, e.ValueSize)
		if _, err = f.File.ReadAt(value, offset); err != nil {
			return
		}
		e.Value = value
	}
	return
}

// 顺序写，写入Entry
func (f *DBFile) Write(e *Entry) (offset int64, err error) {
	buf, err := e.Encode()
	if err != nil {
		return
	}
	_, err = f.File.WriteAt(buf, f.OffSet)
	if err != nil {
		return
	}
	offset = f.OffSet
	f.OffSet += e.GetSize()
	return
}
