package main

import (
	"fmt"
	"sync"
)

var (
	// 对象池
	pool *sync.Pool
)

type Book struct {
	ISBN string
	Name string
}

func initPool() {
	pool = &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating a new Book")
			return new(Book)
		},
	}
}

func main() {
	initPool()
	// 转换为 book
	book := (pool.Get().(*Book))

	book.Name = "first"
	fmt.Printf("设置 book.Name = %s\n", book.Name)

	pool.Put(book)

	fmt.Println("Pool 里已有一个对象：&{first}，调用 Get: ", pool.Get().(*Book))
	fmt.Println("Pool 没有对象了，调用 Get: ", pool.Get().(*Book))
}
