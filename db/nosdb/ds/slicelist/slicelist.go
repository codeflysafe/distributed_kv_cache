package slicelist

import "errors"

const (
	SLICE_LEN = 8
)

// 采用 slice 包装的 list
type SliceList struct {
	// lpush 放入此中
	left [][]byte
	// rpush 放入此中
	right [][]byte
}

func NewSliceList() *SliceList {
	return &SliceList{
		left:  make([][]byte, 0, SLICE_LEN),
		right: make([][]byte, 0, SLICE_LEN),
	}
}

func (sl *SliceList) LLen() int {
	return len(sl.left) + len(sl.right)
}

// 由于 slice array 的特性， 需要修改 O(N)
// 有时候需要扩容，此时会增加复杂度
func (sl *SliceList) LPush(value []byte) {
	sl.left = append(sl.left, value)
}

func (sl *SliceList) LPop() (value []byte) {
	if len(sl.left) == 0 && len(sl.right) == 0 {
		value = nil
	} else if len(sl.left) > 0 {
		value = sl.left[len(sl.left)-1]
		sl.left = sl.left[:len(sl.left)-1]
	} else {
		value = sl.right[0]
		sl.right = sl.right[1:]
	}
	return
}

// 取出 最左边元素
func (sl *SliceList) LPeek() (value []byte) {
	if len(sl.left) == 0 && len(sl.right) == 0 {
		value = nil
	} else if len(sl.left) > 0 {
		value = sl.left[len(sl.left)-1]
	} else {
		value = sl.right[0]
	}
	return
}

func (sl *SliceList) RPush(value []byte) {
	sl.right = append(sl.right, value)
}

func (sl *SliceList) RPop() (value []byte) {
	// 删除尾部元素
	if len(sl.left) == 0 && len(sl.right) == 0 {
		value = nil
	} else if len(sl.right) > 0 {
		value = sl.right[len(sl.right)-1]
		sl.right = sl.right[:len(sl.right)-1]
	} else {
		value = sl.left[0]
		sl.left = sl.left[1:]
	}
	return
}

func (sl *SliceList) RPeek() (value []byte) {
	if len(sl.left) == 0 && len(sl.right) == 0 {
		value = nil
	} else if len(sl.right) > 0 {
		value = sl.right[len(sl.right)-1]
	} else {
		value = sl.left[0]
	}
	return
}

// idx >= 0
// idx < 0 -1 代表最后一个
func (sl *SliceList) ListSeek(idx int) (value []byte, err error) {
	l := sl.LLen()
	var index = idx
	if idx < 0 {
		index = l + idx
	}
	if l < index {
		err = errors.New(" out of range")
		return
	}
	if len(sl.left) <= index {
		// 在右边找
		value = sl.right[index-len(sl.left)]
	} else {
		value = sl.left[index]
	}
	return
}

// 删除操作非常耗时, 涉及到内存的从新分配
func (sl *SliceList) ListDelIndex(idx int) {
	l := sl.LLen()
	var index = idx
	if idx < 0 {
		index = l + idx
	}
	if l < index {
		return
	}
	if len(sl.left) <= index {
		// 在右边找
		// 在右边删除
		sl.right = append(sl.right[:index-len(sl.left)], sl.right[index-len(sl.left)+1:]...)
	} else {
		sl.right = append(sl.right[:index], sl.right[index+1:]...)
	}
	return
}

func (sl *SliceList) Empty() bool {
	return sl.LLen() == 0
}
