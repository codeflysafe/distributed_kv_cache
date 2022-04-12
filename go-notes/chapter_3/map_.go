package main

import "fmt"

func main() {

	kv := make(map[int]string, 10)
	// 没有初始化，map， 无法写入值，但是可以读取，读取的结果为空
	for i := 0; i < 12; i++ {
		kv[i] = fmt.Sprintf("i is %d ", i)
	}
	if v, ok := kv[1]; ok {
		fmt.Println(v)
	}

}
