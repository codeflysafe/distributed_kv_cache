package string

import (
	"nosdb/utils"
)

const (
	DEFAULT_STRING_SIZE = 8
)

type String struct {
	value []byte
}

func NewString() *String {
	return &String{
		value: make([]byte, 0, DEFAULT_STRING_SIZE),
	}
}

func (s *String) Set(value []byte) {
	s.value = value
}

func (s *String) Get() []byte {
	return s.value
}

// todo support -1
func (s *String) GetRange(start, end int) []byte {
	return s.value[start:end]
}

func (s *String) StrLen() int {
	return len(s.value)
}

//如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 IncrByInt 操作。
//如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
//本操作的值限制在 64 位(bit)有符号数字表示之内。
func (s *String) IncrByInt(offset int) (err error) {
	s.value, err = utils.BytesIncrBy(s.value, offset)
	return
}

//如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 IncrByFloat 操作。
func (s *String) IncrByFloat(offset float64) (err error) {
	s.value, err = utils.ByteIncrByFloat(s.value, offset)
	return
}

// 如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
// 如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
func (s *String) Append(as []byte) {
	s.value = append(s.value, as...)
}

func (s *String) String() string {
	return string(s.value)
}
