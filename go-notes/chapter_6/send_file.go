package main

import (
	"fmt"
	"os"
	"syscall"
)

func fileCopy() {
	var file, newFile *os.File
	var err error
	file, err = os.Open("a.txt")
	if err != nil {
		panic(err)
	}
	newFile, err = os.Create("a_copy.txt")
	if err != nil {
		panic(err)
	}
	var offset int64
	var fileInfo os.FileInfo
	fileInfo, err = file.Stat()
	if err != nil {
		panic(err)
	}
	var w int
	w, err = syscall.Sendfile(int(newFile.Fd()), int(file.Fd()), &offset, int(fileInfo.Size()))
	if err != nil {
		panic(err)
	}
	fmt.Println(w)
}

func main() {
	fileCopy()
}
