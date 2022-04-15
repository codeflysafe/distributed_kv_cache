package ds

type ENCODING uint8

const (
	LINKED_LIST ENCODING = iota
	SLICE_LIST
)

// list 接口
type List interface {
	LLen() int
	LPush(value []byte)
	LPop() (value []byte)
	RPush(value []byte)
	RPop() (value []byte)
	ListSeek(idx int) (value []byte)
	ListDelIndex(idx int)
	LPeek() (value []byte)
	RPeek() (value []byte)
	Empty() bool
}
