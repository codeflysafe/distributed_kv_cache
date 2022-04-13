package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	// 状态值
	flags uint8 = 1
	// 读写锁
	mu sync.RWMutex
)

func read(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.RLock()
	defer mu.RUnlock()
	// 读一会
	time.Sleep(20)
	fmt.Printf(" %d am reading this book, flags is %d \n", id, flags)

}

func write(id int, status uint8, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()
	// 写一会
	time.Sleep(10)
	flags = status
	fmt.Printf(" %d am writing this book, flags : %d \n", id, flags)

}

func main() {
	group := &sync.WaitGroup{}
	x, y := 3, 5
	group.Add(x + y)
	// 三个人在读
	for i := 0; i < 3; i++ {
		go read(i, group)
	}

	for i := 2; i < 7; i++ {
		go write(i, uint8(i), group)
	}
	group.Wait()
}
