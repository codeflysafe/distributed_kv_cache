package mmap

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"unsafe"
)

// todo mmap 封装好
const (
	RDONLY = 0
	RDWR   = 1 << iota
	COPY
	EXEC
)

type MMap struct {
	data   []byte
	inProt int
	length int
}

// NewMMap 是对 syscall.Mmap 一个简单的封装
// 返回值是data 或者 error
func NewMMap(file *os.File, inProt int, length int) (*MMap, error) {
	if file == nil {
		err := errors.New(" file is nil ")
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if info.Size() < int64(length) {
		err = fmt.Errorf(" file size %d less length %d, please check it ", info.Size(), length)
		return nil, err
	}
	data, err := mmap(length, inProt, int(file.Fd()), 0)
	if err != nil {
		return nil, err
	}
	mm := &MMap{
		data:   data,
		inProt: inProt,
		length: length,
	}
	return mm, nil
}

func (m *MMap) header() *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(m))
}

func (m *MMap) addrLen() (uintptr, uintptr) {
	header := m.header()
	return header.Data, uintptr(header.Len)
}

// WriteAt 向 offset 处写入 value
// n > 0 写入成功 -1 写入失败
func (mm *MMap) WriteAt(offset int, value []byte) int {
	// 如果总长度不足 offset+len(value)， 写入失败
	if mm.length < offset+len(value) {
		return -1
	}
	// 注意，有可能会覆盖掉之前写入的内容
	// 将value内的数据 copy 到 mm.data 中
	return copy(mm.data[offset:offset+len(value)], value)
}

func (mm *MMap) ReadAt(offset int, length int) (b []byte, err error) {
	if mm.length < offset+length {
		err = errors.New("offset length out of range ")
		return
	}
	b = make([]byte, length, length)
	copy(b, mm.data[offset:offset+length])
	return b, nil
}

// MSync flushes changes made to the in-core copy of a file that
// was mapped into memory using mmap(2) back to the filesystem.
// Without use of this call, there is no guarantee that changes are
// written back before munmap(2) is called.  To be more precise, the
// part of the file that corresponds to the memory area starting at
// addr and having length length is updated.
func (mm *MMap) MSync() error {
	return msync(mm.data)
}

// The MunMap system call deletes the mappings for the specified
// address range, and causes further references to addresses within
// the range to generate invalid memory references.  The region is
// also automatically unmapped when the process is terminated.  On
// the other hand, closing the file descriptor does not unmap the
// region.
func (mm *MMap) MunMap() error {
	return munmap(mm.data)
}
