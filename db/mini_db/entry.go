package mini_db

import (
	"encoding/binary"
	"fmt"
)

// 4 + 4 + 2
const entryHeaderSize = 10

const (
	PUT uint16 = iota
	DEL
)

// 一个存储条目
//
type Entry struct {
	Key       []byte
	KeySize   uint32 // key 的长度
	Value     []byte
	ValueSize uint32 // value 长度
	Mark      uint16 // 标记，是否修改过
}

func NewEntry(key, value []byte, mark uint16) *Entry {
	return &Entry{
		Key:       key,
		KeySize:   uint32(len(key)),
		Value:     value,
		ValueSize: uint32(len(value)),
		Mark:      mark,
	}
}

func (e *Entry) GetSize() int64 {
	return (int64)(entryHeaderSize + e.KeySize + e.ValueSize)
}

// encode 编码，返回字节数组
func (e *Entry) Encode() ([]byte, error) {
	buf := make([]byte, e.GetSize())
	binary.BigEndian.PutUint32(buf[0:4], e.KeySize)
	binary.BigEndian.PutUint32(buf[4:8], e.ValueSize)
	binary.BigEndian.PutUint16(buf[8:10], e.Mark)
	copy(buf[entryHeaderSize:entryHeaderSize+e.KeySize], e.Key)
	copy(buf[entryHeaderSize+e.KeySize:], e.Value)
	return buf, nil
}

// 将字节流解码为 entry 实体
func Decode(buf []byte) (entry *Entry, err error) {
	if len(buf) != entryHeaderSize {
		err = fmt.Errorf(" len is not match ")
		return
	}
	entry = &Entry{}
	entry.KeySize = binary.BigEndian.Uint32(buf[0:4])
	entry.ValueSize = binary.BigEndian.Uint32(buf[4:8])
	entry.Mark = binary.BigEndian.Uint16(buf[8:10])
	return
}
