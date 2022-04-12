package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Z struct {
	z unsafe.Pointer
}

type Y struct {
}

type X struct {
	name string
	id   int
	p    *Y
}

func getInt() int {
	a := 55
	return a
}

func foo() *int {
	a1 := 55
	return &a1
}

func foo2() *X {
	x2 := X{
		name: "hsj",
		id:   12,
	}
	return &x2
}

func foo3() *X {
	xf2 := &X{
		name: "hsj",
		id:   12,
	}
	return xf2
}

func foo4() *X {
	xf4 := new(X)
	return xf4
}

func foo5() {
	_ = new(X)
}

func reflectv() {
	xv := X{}
	xc := interface{}(xv)
	t := reflect.TypeOf(xc)
	_ = t.Size()
}

// 8192 KB
func lessThanStack() {
	var xs [100]int64 // 800B
	xs[0] = 0
}

// 8192 KB
func more32ThanStack() {
	var xm3 [5000]int64 // 40KB
	//fmt.Println(unsafe.Sizeof(xm))
	xm3[0] = 0
}

// 8192 KB
func moreThanStack() {
	var xm [2000000]int64 // 16MB
	//fmt.Println(unsafe.Sizeof(xm))
	xm[0] = 0
}

// 闭包产生的逃逸
func Increase() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

func unknownum() {
	num := 12
	cn := make([]int, num)
	cn = append(cn, 1)

	cns := make([]int, 121)
	cns = append(cns, 123)
}

func main() {
	b := getInt()
	var c = b
	// 1. 由于fmt.Println导致的内存逃逸
	fmt.Println(c)

	// 2. 函数返回指向栈内对象的指针，或者说是参数泄漏，延长了指针对象的生命周期。
	foo()

	// 3. 涉及到反射使用
	reflectv()
	// 4. 被已经逃逸的变量引用的指针，一定发生逃逸
	x3 := foo2()
	y := Y{}
	x3.p = &y

	// 5.栈空间不足引发逃逸
	lessThanStack()
	moreThanStack()
	// 超过32KB
	more32ThanStack()

	// 6.变量大小不确定
	unknownum()

	_ = Z{unsafe.Pointer(&Y{})}
	p2 := foo3()
	p2.name = "xws"

	foo4()

	foo5()
	// 9. 闭包发生逃逸
	Increase()
}
