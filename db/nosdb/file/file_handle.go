/*
 * @Author: sjhuang
 * @Date: 2022-04-22 15:15:32
 * @LastEditTime: 2022-04-24 10:57:18
 * @FilePath: /nosdb/file/file_handle.go
 */
package file

import (
	"errors"
	"io"
	"nosdb/utils"
	"os"
)

type MOD uint8

const (
	STANDARD_IO MOD = iota
	M_MAP
)

type FileHandle interface {
	io.Closer
	ReadAt(offset int64, length int) (data []byte, err error)
	WriteAt(offset int64, data []byte) (newOffset int64, err error)
	Sync() error
	Delete() error
	IsClose() bool
	Offset() (int64, error)
}

func NewFileHandle(mod MOD, file *os.File, maxLength int) (FileHandle, error) {
	switch mod {
	case M_MAP:
		return newMMapFileHandle(file, maxLength)
	default:
		return newIOFileHandle(file, maxLength)
	}
}

func OpenFile(mod MOD, path string, fileName string, maxLength int) (FileHandle, error) {
	if maxLength <= 0 {
		err := errors.New(" error maxlength ")
		return nil, err
	}
	file, err := utils.Open(path, fileName)
	if err != nil {
		return nil, err
	}
	return NewFileHandle(mod, file, maxLength)
}
