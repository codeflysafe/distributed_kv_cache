package discache

// ☝️不可变类，类似于java中的string
// 类似于 redis 的SDS
type ByteView struct {
	b []byte
}

// return ByteView 的长度
func (v ByteView) Len() int {
	return len(v.b)
}

// 返回一个 copy
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// 返回一个字符串
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
