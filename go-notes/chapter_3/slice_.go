package main

import (
	"fmt"
	"unsafe"
)

// 64 位最大值
const MaxUintptr = ^uintptr(0)

func slice1() {
	nums := make([]int, 4)
	// 1个字长为64位，3个字长， 3*64/8 = 24 个字节，
	println(unsafe.Alignof(nums), unsafe.Sizeof(nums), len(nums), cap(nums))
	for i := 0; i < 10; i++ {
		nums = append(nums, i)
		println(unsafe.Alignof(nums), unsafe.Sizeof(nums), len(nums), cap(nums))
	}
}

func slice2() {
	var nums []int
	// 1个字长为64位，3个字长， 3*64/8 = 24 个字节，
	println(unsafe.Alignof(nums), unsafe.Sizeof(nums), len(nums), cap(nums))
	for i := 0; i < 10; i++ {
		nums = append(nums, i)
		println(unsafe.Alignof(nums), unsafe.Sizeof(nums), len(nums), cap(nums))
	}
}

func main() {
	// 8 个字节
	fmt.Println(MaxUintptr, unsafe.Sizeof(MaxUintptr), ^uint64(0))
	// 4->8->16
	fmt.Println("slice 1 ------->")
	slice1()
	fmt.Println("slice 2 ------->")
	// 0->1->2->4->8->16
	slice2()

	fmt.Println("test append -------------- ")
	arr1 := []int{1, 2, 3}
	fmt.Println(arr1, len(arr1), cap(arr1))
	arr2 := append(arr1, 4, 5, 8, 10)
	fmt.Println(arr1, len(arr1), cap(arr1))
	fmt.Println(arr2, len(arr2), cap(arr2))
	fmt.Println("test append end-------------- ")

	nums := []int{1, 2, 3, 4, 5, 6, 6, 7, 8}
	k := nums[0:2]
	k2 := nums[2:4]
	fmt.Println(k, len(k), cap(k))
	fmt.Println(k2, len(k2), cap(k2))
	fmt.Println(nums, len(nums), cap(nums))
	nums[3] = 90
	nums[0] = 12
	fmt.Println(k, len(k), cap(k))
	fmt.Println(k2, len(k2), cap(k2))
	fmt.Println(nums, len(nums), cap(nums))

	k[0] = 128
	k2[0] = 87
	fmt.Println(k, len(k), cap(k))
	fmt.Println(k2, len(k2), cap(k2))
	fmt.Println(nums, len(nums), cap(nums))

}
