package main

import "fmt"

type KV struct {
	Key   int
	Value int
}

func main() {

	cache := make(map[int64]KV, 16)

	cache[1] = KV{Key: 1, Value: 1}
	cache[2] = KV{Key: 2, Value: 2}

	v1 := cache[1]
	fmt.Println(v1)

	if v2, ok := cache[2]; ok {
		fmt.Println(v2)
	}
}
