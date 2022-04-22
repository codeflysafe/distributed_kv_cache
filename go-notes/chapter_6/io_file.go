package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
)

func main() {
	// 256 kb 采样一次
	runtime.MemProfileRate = 256 * 1024
	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()
	// 加载文件
	file, err := os.OpenFile("a.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0664)
	if err != nil {
		panic(err)
	}
	fmt.Println(file.Name())
	//file.ReadAt()
}
