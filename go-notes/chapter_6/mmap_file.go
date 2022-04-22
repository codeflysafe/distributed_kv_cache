package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
)

const (
	filePath    = "a.txt"
	MaxFileSize = 1 << 25 // file 最大为多少, 32MB
	// 4096
)

// 只执行一次，线程安全
func init() {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	// 这里写入一个 12.6MB 的文件
	l := MaxFileSize
	buf := make([]byte, l, l)
	var n int
	n, err = file.Write(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
}

func main() {
	fmt.Println(os.Getpagesize())
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	fd := int(file.Fd())
	data, err := syscall.Mmap(fd, 0, MaxFileSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(data), MaxFileSize)
	var idx uint32 = 1324
	binary.BigEndian.PutUint32(data[0:4], idx)

	// 刷回磁盘
	syscall.Munmap(data)
}
