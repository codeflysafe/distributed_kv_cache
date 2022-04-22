package main

import (
	"fmt"
	"sync"
)

var (
	once sync.Once
)

type A struct {
	index map[string][]byte
	child *A
}

func NewA() *A {
	return new(A)
}

func (a *A) lazyLoadIndex() {
	once.Do(func() {
		fmt.Println("lazyLoad index")
		if a.index == nil {
			a.index = make(map[string][]byte, 10)
		}
	})
}

func (a *A) lazyLoadChild() {
	once.Do(func() {
		fmt.Println("lazyLoad child")
		if a.child == nil {
			a.child = NewA()
		}
	})
}

func main() {
	a := new(A)
	wait := sync.WaitGroup{}
	count := 5
	wait.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			a.lazyLoadIndex()
			wait.Done()
		}()
	}
	wait.Wait()

	wait.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			a.lazyLoadChild()
			wait.Done()
		}()
	}
	wait.Wait()
}
