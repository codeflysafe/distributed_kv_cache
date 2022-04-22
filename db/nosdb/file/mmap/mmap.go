package mmap

import (
	"errors"
	"os"
	"syscall"
)

// todo mmap 封装好
const (
	M_MAP_READ  = syscall.PROT_READ  // 只允许读，不允许写
	M_MAP_WRITE = syscall.PROT_WRITE // 只写
	M_MAP_EXEC  = syscall.PROT_EXEC  // 可执行
	M_MAP_NONE  = syscall.PROT_NONE  // 不可操作

	M_MAP_SHARED = syscall.MAP_SHARED
)

type MMap struct {
	data   []byte
	prot   int
	length int
}

// NewMMap 是对 syscall.Mmap 一个简单的封装
// 返回值是data 或者 error
func NewMMap(file *os.File, prot int, length int) (*MMap, error) {
	if file == nil {
		err := errors.New(" file is nil ")
		return nil, err
	}
	data, err := mmap(int(file.Fd()), 0, length, prot, M_MAP_SHARED)
	if err != nil {
		return nil, err
	}
	mm := &MMap{
		data:   data,
		prot:   prot,
		length: length,
	}
	return mm, nil
}
