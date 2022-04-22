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
	err = h.mm.MunMap()
	return
}

func (h *MMapFileHandle) Delete() (err error) {
	if err = h.checkFile(); err != nil {
		return
	}
	_ = h.Close()
	return os.Remove(h.file.Name())
}
