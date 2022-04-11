package main

import (
	"fmt"
	"unsafe"
)

func main() {

	var a = 32
	var b int64 = 878
	var p1 *int
	var p2 *int64

	p1, p2 = &a, &b

	fmt.Println(p1, p2, unsafe.Sizeof(p1), unsafe.Sizeof(a), unsafe.Sizeof(p2), unsafe.Sizeof(b), *p1, *p2)

	var c byte = 'a'
	var p3 = &c
	fmt.Println(p3, unsafe.Sizeof(p3), unsafe.Sizeof(c), *p3)

	var d = "1234567890"
	fmt.Println(d, unsafe.Sizeof(d))

	var p4 = &struct {
		name string
		addr string
	}{name: "sss", addr: "上乃事"}

	fmt.Println(p4, unsafe.Sizeof(p4), unsafe.Sizeof(*p4), *p4)

	//	- A pointer value of any type can be converted to a Pointer.
	//	- A Pointer can be converted to a pointer value of any type.
	//	- A uintptr can be converted to a Pointer.
	//	- A Pointer can be converted to a uintptr.
	var p5 = (*uint64)(unsafe.Pointer(p4))
	fmt.Println(p5, *p5)

	// uintptr is an integer type that is large enough to hold the bit pattern of
	// any pointer.
	var x struct {
		a bool
		b int16
		c []int
	}
	pb := (*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
	*pb = 64

	fmt.Println(x.b)

	// bad , 可能会出错
	tmp := uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
	pc := (*int16)(unsafe.Pointer(tmp))
	*pc = 128
	fmt.Println(x.b)

}
