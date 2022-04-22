package file

import (
	"io"
	"nosdb"
	"os"
)

type IOFileHandle struct {
	file *os.File
}

func (h *IOFileHandle) checkFile() error {
	if h.file == nil {
		return nosdb.FileNotLoadError
	}
	return nil
}

func (h *IOFileHandle) ReadAt(offset int64, length int) (data []byte, err error) {
	err = h.checkFile()
	if err != nil {
		return
	}
	data = make([]byte, length, length)
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
	return h.file.Close()
}

func (h *IOFileHandle) Delete() error {
	_ = h.Close()
	return os.Remove(h.file.Name())
}
