package main

import (
	"sync"
)

func fib(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return fib(n-2) + fib(n-1)
}

func lock1() {
	var mu sync.Mutex

	mu.Lock()
	go func() {
		fib(10)
		mu.Unlock()
	}()

	mu.Lock()
	mu.Unlock()
}

func chan1() {
	done := make(chan int)

	go func() {
		fib(10)
		done <- 1
	}()

	<-done
}
