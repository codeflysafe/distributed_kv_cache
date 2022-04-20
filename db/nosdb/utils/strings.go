package utils

import (
	"fmt"
	"strconv"
)

const (
	INT_MAX = int(^uint(0) >> 1)
	INT_MIN = ^INT_MAX
)

// 增加或者减少 一个 offset
// 范围一个 修改后的 string， 或者一个错误 if error !=nil
//本操作的值限制在 64 位(bit)有符号数字表示之内。
func BytesIncrBy(val []byte, offset int) (newVal []byte, err error) {
	if val == nil {
		err = fmt.Errorf(" val is blank  %v", val)
		return
	}
	var str string
	str, err = StrIncrBy(string(val), offset)
	if err != nil {
		return
	}
	newVal = []byte(str)
	return
}

// 增加或者减少 一个 offset
// 范围一个 修改后的 string， 或者一个错误 if error !=nil
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
func StrIncrBy(val string, offset int) (newVal string, err error) {
	var vs int
	vs, err = strconv.Atoi(val)
	if err != nil {
		return
	}
	if offset < 0 {
		if INT_MIN-offset > vs {
			err = fmt.Errorf(" val is overflow minv %d, %d", vs, offset)
			return
		}
	} else {
		if INT_MAX-offset < vs {
			err = fmt.Errorf(" val is overflow maxV %d, %d", vs, offset)
			return
		}
	}
	vs += offset
	newVal = strconv.Itoa(vs)
	return
}

// 字段值加上指定浮点数增量值。
func StrIncrByFloat(val string, offset float64) (newVal string, err error) {
	var vs float64
	vs, err = strconv.ParseFloat(val, 64)
	if err != nil {
		return
	}
	vs += offset
	newVal = strconv.FormatFloat(vs, 'g', -1, 64)
	return
}

// 会存在精度问题
func ByteIncrByFloat(val []byte, offset float64) (newVal []byte, err error) {
	var vs float64
	vs, err = strconv.ParseFloat(string(val), 64)
	if err != nil {
		return
	}
	vs += offset
	newVal = []byte(strconv.FormatFloat(vs, 'g', -1, 64))
	return
}
