package mmap

import (
	"golang.org/x/sys/unix"
)

func mmap(len, prot, flag int, fd uintptr, off int64) ([]byte, error) {
	b, err := unix.Mmap(int(fd), off, len, prot, flag)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// WriteAt 向 offset 处写入 value
// 1 写入成功 0 写入失败
func (mm *MMap) WriteAt(offset int, value []byte) bool {
	// 如果总长度不足 offset+len(value)， 写入失败
	if mm.length < offset+len(value) {
		return false
	}
	// 注意，有可能会覆盖掉之前写入的内容
	// 将value内的数据 copy 到 mm.data 中
	copy(mm.data[offset:offset+len(value)], value)
	return true
}

// MSync flushes changes made to the in-core copy of a file that
// was mapped into memory using mmap(2) back to the filesystem.
// Without use of this call, there is no guarantee that changes are
// written back before munmap(2) is called.  To be more precise, the
// part of the file that corresponds to the memory area starting at
// addr and having length length is updated.
func (mm *MMap) MSync() error {
	return unix.Msync(mm.data, unix.MS_SYNC)
}

// The MunMap system call deletes the mappings for the specified
// address range, and causes further references to addresses within
// the range to generate invalid memory references.  The region is
// also automatically unmapped when the process is terminated.  On
// the other hand, closing the file descriptor does not unmap the
// region.
func (mm *MMap) MunMap() error {
	return unix.Munmap(mm.data)
}
