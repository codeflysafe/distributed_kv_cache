package main

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func add(ptr unsafe.Pointer, offset int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) + uintptr(offset))
}

func pointer() {
	var str string
	str = "123123"
	p := (*reflect.StringHeader)(unsafe.Pointer(&str))
	fmt.Println(p)

	// 访问第一个元素
	ptr := unsafe.Pointer(p.Data)
	fmt.Println(*(*byte)(ptr), '1')

	arr := (*[5]byte)(ptr)
	fmt.Println(arr)

	// 访问数组的第i个元素 地址寻址 按照字
	fmt.Println(arr[3], *(*byte)(add(ptr, 3)))
}

func change() {
	var str, str2, str3, str4 string
	str = "123123"
	p := (*reflect.StringHeader)(unsafe.Pointer(&str))

	str2 = str[0:3]
	p2 := (*reflect.StringHeader)(unsafe.Pointer(&str2))

	str3 = str[1:4]
	p3 := (*reflect.StringHeader)(unsafe.Pointer(&str3))

	// 对比 p1 和 p2 的位置, 位置是相同的, 但是p3的指针向前移动一个字长
	// 说明，与切片一样，字符串切片底层也是共享的一个底层数组
	fmt.Println(p.Data, p2.Data, p3.Data)

	// 对比数组内容, 指针没变，但是长度修改了
	arr1 := (*[5]byte)(unsafe.Pointer(p.Data))
	arr2 := (*[3]byte)(unsafe.Pointer(p2.Data))
	arr3 := (*[3]byte)(unsafe.Pointer(p3.Data))
	fmt.Println(arr1, str, p.Len)
	fmt.Println(arr2, str2, p2.Len)
	fmt.Println(arr3, str3, p3.Len)
	// 修改呢？
	// str[0] = 12 违法的操作
	//1。 采用byte[] 数组修改，然后在变成 string
	fmt.Println(" 采用byte[] 数组修改，然后在变成 string")
	c := []byte(str)
	fmt.Println(c)
	for i, _ := range c {
		// 修改，可以吗？
		c[i] = 'k'
	}
	fmt.Println(c)
	// str4 := string(c) 变回字符串
	fmt.Println(" 采用byte[] 数组修改，------- string")

	// 2. + 号的作用
	fmt.Println("采用 + 修改字符串  ------- ")
	fmt.Println(str4)
	str4 = str2 + str3
	p4 := (*reflect.StringHeader)(unsafe.Pointer(&str4))
	fmt.Println(str4, p4.Data)

	fmt.Println("采用 + 修改字符串  ------- end ")
	// 2。 stringbuild ?
	build := strings.Builder{}
	// 3
	build.WriteString(str2)
	// 3
	build.WriteString(str4)
	// buf := make([]byte, len(b.buf), 2*cap(b.buf)+n) 扩容
	fmt.Println(build.String(), build.Cap(), build.Len())
}

func main() {
	// pointer()
	change()
}
