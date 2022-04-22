package file

import (
	"io"
	"nosdb/utils"
	"os"
)

type MOD uint8

const (
	STANDARD_IO MOD = iota
	M_MAP

	FILE_MAX_LENGTH = 1 << 20
)

type FileHandle interface {
	io.Closer
	ReadAt(offset int64, length int) (data []byte, err error)
	WriteAt(offset int64, data []byte) (newOffset int64, err error)
	Sync() error
	Delete() error
	IsClose() bool
}

func NewFileHandle(mod MOD, file *os.File, maxLength int) (FileHandle, error) {
	if maxLength > FILE_MAX_LENGTH || maxLength <= 0 {
		maxLength = FILE_MAX_LENGTH
	}
	switch mod {
	case M_MAP:
		return newMMapFileHandle(file, maxLength)
	default:
		return newIOFileHandle(file, maxLength)
	}
}

func OpenFile(mod MOD, path string, fileName string, maxLength int) (FileHandle, error) {
	file, err := utils.Open(path, fileName)
	if err != nil {
		return nil, err
	}
	return NewFileHandle(mod, file, maxLength)
}
