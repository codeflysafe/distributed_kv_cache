/*
 * @Author: sjhuang
 * @Date: 2022-04-22 10:56:54
 * @LastEditTime: 2022-04-24 10:56:50
 * @FilePath: /nosdb/file/mmap_file_handle.go
 */
package file

import (
	"errors"
	"nosdb/file/mmap"
	"os"
)

type MMapFileHandle struct {
	mm        *mmap.MMap
	file      *os.File
	maxLength int
	close     bool // 文件是否已经关闭
}

func newMMapFileHandle(file *os.File, maxLength int) (FileHandle, error) {
	h := &MMapFileHandle{
		file:      file,
		maxLength: maxLength,
	}
	var err error
	h.mm, err = mmap.NewMMap(file, mmap.RDWR, maxLength)
	return h, err
}

func (h *MMapFileHandle) checkFile() error {
	if h.file == nil || h.mm == nil {
		return errors.New(" no such file")
	}
	return nil
}

func (h *MMapFileHandle) ReadAt(offset int64, length int) (data []byte, err error) {
	if err = h.checkFile(); err != nil {
		return
	}
	return h.mm.ReadAt(int(offset), length)
}

// todo 错误处理
func (h *MMapFileHandle) WriteAt(offset int64, data []byte) (newOffset int64, err error) {
	if err = h.checkFile(); err != nil {
		return
	}
	n := h.mm.WriteAt(int(offset), data)
	if n == -1 {
		err = errors.New(" xxxxx")
		return
	}
	newOffset = int64(n) + offset
	return
}

func (h *MMapFileHandle) Sync() (err error) {
	if err = h.checkFile(); err != nil {
		return
	}
	return h.mm.MSync()
}

func (h *MMapFileHandle) Close() (err error) {
	err = h.file.Close()
	if err != nil {
		return
	}
	h.close = true
	err = h.mm.MunMap()
	return
}

func (h *MMapFileHandle) Delete() (err error) {
	if err = h.checkFile(); err != nil {
		return
	}
	if !h.close {
		err = h.Close()
		if err != nil {
			return
		}
		h.close = true
	}
	return os.Remove(h.file.Name())
}
func (h *MMapFileHandle) IsClose() bool {
	return h.close
}

func (h *MMapFileHandle) Offset() (int64, error) {
	f, err := h.file.Stat()
	if err != nil {
		return 0, err
	}
	return f.Size(), nil
}
