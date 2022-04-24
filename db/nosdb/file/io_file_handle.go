/*
 * @Author: sjhuang
 * @Date: 2022-04-22 10:33:13
 * @LastEditTime: 2022-04-24 11:52:37
 * @FilePath: /nosdb/file/io_file_handle.go
 */
package file

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type IOFileHandle struct {
	file      *os.File
	maxLength int
	close     bool // 文件是否已经关闭
}

func newIOFileHandle(file *os.File, maxLength int) (FileHandle, error) {
	return &IOFileHandle{
		file:      file,
		maxLength: maxLength,
	}, nil
}

func (h *IOFileHandle) checkFile() error {
	if h.file == nil {
		return errors.New("no such file")
	}
	return nil
}

func (h *IOFileHandle) ReadAt(offset int64, length int) (data []byte, err error) {
	err = h.checkFile()
	if err != nil {
		return
	}
	if h.maxLength < int(offset)+length {
		err = fmt.Errorf("length out of range %d %d", offset, h.maxLength)
		return
	}
	data = make([]byte, length)
	_, err = h.file.ReadAt(data, offset)
	// ReadAt always returns a non-nil error when n < len(b).
	// At end of file, that error is io.EOF.
	// 不是错误，到达了文件结尾
	if err == io.EOF {
		err = nil
		return
	}
	return
}

func (h *IOFileHandle) WriteAt(offset int64, data []byte) (newOffset int64, err error) {
	if err = h.checkFile(); err != nil {
		return
	}
	if h.maxLength < int(offset)+len(data) {
		err = fmt.Errorf("length out of range %d %d", offset, h.maxLength)
		return
	}
	_, err = h.file.WriteAt(data, offset)
	if err != nil {
		return
	}
	newOffset = offset + int64(len(data))
	return
}

func (h *IOFileHandle) Sync() error {
	return h.file.Sync()
}

func (h *IOFileHandle) Close() error {
	err := h.file.Close()
	if err != nil {
		return err
	}
	h.close = true
	return nil
}

func (h *IOFileHandle) Delete() error {
	if !h.close {
		err := h.Close()
		if err != nil {
			return err
		}
		h.close = true
	}
	return os.Remove(h.file.Name())
}

func (h *IOFileHandle) IsClose() bool {
	return h.close
}

func (h *IOFileHandle) Offset() (int64, error) {
	f, err := h.file.Stat()
	if err != nil {
		return 0, err
	}
	return f.Size(), nil
}
