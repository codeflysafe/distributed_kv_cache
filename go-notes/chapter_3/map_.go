package main

import "fmt"

func main() {

	kv := make(map[int]string, 10)
	// 没有初始化，map， 无法写入值，但是可以读取，读取的结果为空
	for i := 0; i < 12; i++ {
		kv[i] = fmt.Sprintf("i is %d ", i)
	}

	// 访问方式1 mapaccess1_fast64
	va := kv[1]
	fmt.Println(va)

	// 访问方式2 mapaccess2_fast64
	if v, ok := kv[1]; ok {
		fmt.Println(v)
	}

	// 遍历 mapiterinit
	for k, v := range kv {
		fmt.Println(k, v)
	}

}
