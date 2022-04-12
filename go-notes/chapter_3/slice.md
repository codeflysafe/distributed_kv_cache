# slice 切片

```go
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
```

![go slice](https://raw.githubusercontent.com/codeflysafe/gitalk/main/img/20220412142052.png)

## 用法
```go
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
}

```
输出结果
```go
slice 1 ------->
8 24 4 4
8 24 5 8
8 24 6 8
8 24 7 8
8 24 8 8
8 24 9 16
8 24 10 16
8 24 11 16
8 24 12 16
8 24 13 16
8 24 14 16
slice 2 ------->
8 24 0 0
8 24 1 1
8 24 2 2
8 24 3 4
8 24 4 4
8 24 5 8
8 24 6 8
8 24 7 8
8 24 8 8
8 24 9 16
8 24 10 16
```

## 1.make([]T, len, cap int)

```go
// go/src/runtime/slice.go
func makeslice(et *_type, len, cap int) unsafe.Pointer
```
会调用`mallocgc` 去分配内存

```
// Allocate an object of size bytes.
// Small objects are allocated from the per-P cache's free lists. （stack）
// Large objects (> 32 kB) are allocated straight from the heap. 
```

大目标直接分配到堆上(试了试好像并不是)


## 2. 扩容
`growslice` handles slice growth during `append`

它通过slice的元素大小，旧slice和希望的新的最小cap来实现，并且返回一个新的slice，旧的数据通过copy，转移到
新的slice中

```go
// go/src/runtime/slice.go
func growslice(et *_type, old slice, cap int) slice
```

扩容大小从 2x -> 1.25x，随着大小变化，增长率降低

```go
    newcap := old.cap
	doublecap := newcap + newcap
	if cap > doublecap {
		newcap = cap
	} else {
		const threshold = 256
		if old.cap < threshold {
			newcap = doublecap
		} else {
			// Check 0 < newcap to detect overflow
			// and prevent an infinite loop.
			for 0 < newcap && newcap < cap {
				// Transition from growing 2x for small slices
				// to growing 1.25x for large slices. This formula
				// gives a smooth-ish transition between the two.
				newcap += (newcap + 3*threshold) / 4
			}
			// Set newcap to the requested cap when
			// the newcap calculation overflowed.
			if newcap <= 0 {
				newcap = cap
			}
		}
	}
```

### 切片
```go
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

```
#### 分析
![](https://raw.githubusercontent.com/codeflysafe/gitalk/main/img/20220412152544.png)
1. 与old slice 共享底层数组： 分片之后的返回的新的slice，底层数组没有变化，就是len和cap，以及array指向变化。
2. 长度变为切片长度: 如k 和 k2
3. cap变为切片起始开始到数组结尾的长度： 如k2

#### 运行结果
运行结果也证实了这一点
```shell
[1 2] 2 9
[3 4] 2 7
[1 2 3 4 5 6 6 7 8] 9 9
[12 2] 2 9
[3 90] 2 7
[12 2 3 90 5 6 6 7 8] 9 9
[128 2] 2 9
[87 90] 2 7
[128 2 87 90 5 6 6 7 8] 9 9

```

### append 方法

会涉及到扩容

```go

fmt.Println("test append -------------- ")
arr1 := []int{1, 2, 3}
fmt.Println(arr1, len(arr1), cap(arr1))
arr2 := append(arr1, 4, 5, 8, 10)
fmt.Println(arr1, len(arr1), cap(arr1))
fmt.Println(arr2, len(arr2), cap(arr2))
fmt.Println("test append end-------------- ")
```

运行结果

```shell
test append -------------- 
[1 2 3] 3 3
[1 2 3] 3 3
[1 2 3 4 5 8 10] 7 8
test append end-------------- 
```

一个比较奇怪的东西，为何 `arr2` `cap` 会变成 `8`
这是因为 内存对齐
按照前面的计算，`newcap` 应该是`7`
max(7, 2*3)
```go
capmem = roundupsize(uintptr(newcap) * goarch.PtrSize)
newcap = int(capmem / goarch.PtrSize)

//
// Returns size of the memory block that mallocgc will allocate if you ask for the size.
// Returns size of the memory block that mallocgc will allocate if you ask for the size.
func roundupsize(size uintptr) uintptr {
if size < _MaxSmallSize {
	...
	uintptr(class_to_size[size_to_class8[divRoundUp(size, smallSizeDiv)]])
}
_MaxSmallSize   = 32768
smallSizeDiv    = 8
smallSizeMax    = 1024
```

```shell
size = 7*8 = 56
divRoundUp(56, 8) = 7,
size_to_class8[7] = 6
class_to_size[6] = 64
newcap = int(capmem / goarch.PtrSize) = 64/8 = 8
```

当`newcap =6`时
```shell
size = 6*8 = 48
divRoundUp(48, 8) = 6
size_to_class8[6] = 5
class_to_size[5] = 48
newcap = int(capmem / goarch.PtrSize) = 48/8 = 6
```