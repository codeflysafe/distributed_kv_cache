package main

import (
	"fmt"
	"math/rand"
	"unsafe"
)

type A struct {
	S *string
}

func (f *A) String() string {
	return *f.S
}

type ATrick struct {
	S unsafe.Pointer
}

func (f *ATrick) String() string {
	return *(*string)(f.S)
}

func NewA(s string) A {
	return A{S: &s}
}

func NewATrick(s string) ATrick {
	return ATrick{S: noescape(unsafe.Pointer(&s))}
}

func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func esacpe() *int {
	x1 := 12
	return &x1
}

func Increase() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

func LessThan8192() {
	nums := make([]int, 100) // = 64KB
	for i := 0; i < len(nums); i++ {
		nums[i] = rand.Int()
	}
}

func MoreThan8192() {
	nums := make([]int, 1000000) // = 64KB
	for i := 0; i < len(nums); i++ {
		nums[i] = rand.Int()
	}
}

func F() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func NonConstant() {
	number := 10
	s := make([]int, number)
	for i := 0; i < len(s); i++ {
		s[i] = i
	}
}

func main() {
	s := "hello"

	f1 := NewA(s)
	f2 := NewATrick(s)
	s1 := f1.String()
	s2 := f2.String()

	//函数返回局部指针变量
	esacpe()
	_ = s1 + s2

	//  interface类型逃逸
	str := "wo zai nali"
	fmt.Println(str)

	in := Increase()
	fmt.Println(in()) // 1

	// 变量大小不确定及栈空间不足引发逃逸
	NonConstant()
	MoreThan8192()
	LessThan8192()

	var nums [10]int
	nums[0] = 1

	var slices []int
	slices = append(slices, 0)

	num2 := make([]int, 1, 12)
	num2 = append(num2, 1)

	l := 20
	c := make([]int, 0, l) // 堆 动态分配不定空间 逃逸
	c = append(c, 12)

	f := F()

	for i := 0; i < 10; i++ {
		f()
	}
}
