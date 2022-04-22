package file

import "io"

// FileHandle 对文件处理的抽象
type FileHandle interface {
	io.Closer
	ReadAt(offset int64, length int) (data []byte, err error)       // 从文件offset位置读出length长度的数据
	WriteAt(offset int64, data []byte) (newOffset int64, err error) // 向文件中写入 data
	Sync() error                                                    // 将数据从缓冲区刷新到文件中
	Delete() error                                                  // 删除文件
}
