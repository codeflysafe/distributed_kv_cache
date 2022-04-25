package logfile

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"time"
)

// 4 + 8 + 4 + 4 + 4 + 2 + 1 + 1
const EntryMetaSize = 28

type CMD uint16

const (
	Persistent uint32 = 0
)

type ENCODING uint8

const (
	// list
	LINKED_LIST ENCODING = iota
	SLICE_LIST
	// zset
	SKIP_LIST
	// hash
	M_HASH
	B_STRING
)

type TYPE uint8

const (
	//===========  type =========
	LIST TYPE = iota
	HASH
	SET
	ZSET
	STRING
)

// entry 的一些`meta`信息
type Meta struct {
	CRC       uint32   // crc 校验信息
	Timestamp uint64   // 时间戳
	TTL       uint32   // 定时时间, 0 代表用不过期
	KeySize   uint32   // key 的长度
	ValueSize uint32   // value 长度
	cmd       CMD      // 标记，PUT or DEL
	Ty        TYPE     // 存储类型， List, Set, ZSet, Hash, String, 高 8 位
	Encoding  ENCODING // 底层数据结构，编码类型， 低 8 位
}

//  the Entry stored format:
//           |------------------------------------META------------------------------------|
//  |-------------------------------------------------------------------------------------------------------|
//  |   crc  | timestamp |   TTL   |  keySize  |  valueSize  |  cmd   |   Ty   | Encoding |  key  |  value  |
//  |-------------------------------------------------------------------------------------------------------|
//  | uint32 |  uint64   | uint32  |   uint32  |   uint32    | uint16 | uint16 |  uint16  | []byte | []byte |
//  |-------------------------------------------------------------------------------------------------------|
//  |---crc--|----------------------------------------------------------------------------------------------|
//
type LogEntry struct {
	Meta
	Key   []byte
	Value []byte
}

// NewEntry 新建一条记录
func NewLogEntry(key, value []byte, cmd CMD, TTL uint32, encoding ENCODING, ty TYPE) *LogEntry {
	return &LogEntry{
		Key:   key,
		Value: value,
		Meta: Meta{
			Timestamp: uint64(time.Now().Unix()),
			TTL:       TTL,
			KeySize:   uint32(len(key)),
			ValueSize: uint32(len(value)),
			cmd:       cmd,
			Encoding:  encoding,
			Ty:        ty,
		},
	}
}

// GetCRC 计算 crc32
// buf 要检验的字节流
// 返回 校验结果
func (e *LogEntry) GetCRC(buf []byte) uint32 {
	return crc32.ChecksumIEEE(buf)
}

// GetSize 返回长度
func (e *LogEntry) GetSize() int64 {
	return (int64)(EntryMetaSize + e.KeySize + e.ValueSize)
}

// Encode 编码，返回字节数组
func (e *LogEntry) Encode() ([]byte, error) {
	buf := make([]byte, e.GetSize())
	binary.BigEndian.PutUint64(buf[4:12], e.Timestamp)
	binary.BigEndian.PutUint32(buf[12:16], e.TTL)
	binary.BigEndian.PutUint32(buf[16:20], e.KeySize)
	binary.BigEndian.PutUint32(buf[20:24], e.ValueSize)
	binary.BigEndian.PutUint16(buf[24:26], uint16(e.cmd))
	binary.BigEndian.PutUint16(buf[26:28], (uint16(e.Ty)<<8)|uint16(e.Encoding))
	copy(buf[EntryMetaSize:EntryMetaSize+e.KeySize], e.Key)
	copy(buf[EntryMetaSize+e.KeySize:], e.Value)
	binary.BigEndian.PutUint32(buf[0:4], e.GetCRC(buf[4:]))
	return buf, nil
}

// DecodeMeta 将字节流解码为 entry meta 实体
func DecodeMeta(buf []byte) (entry *LogEntry, err error) {
	if len(buf) != EntryMetaSize {
		err = fmt.Errorf(" len is not match ")
		return
	}
	entry = &LogEntry{}
	entry.CRC = binary.BigEndian.Uint32(buf[0:4])
	entry.Timestamp = binary.BigEndian.Uint64(buf[4:12])
	entry.TTL = binary.BigEndian.Uint32(buf[12:16])
	entry.KeySize = binary.BigEndian.Uint32(buf[16:20])
	entry.ValueSize = binary.BigEndian.Uint32(buf[20:24])
	entry.cmd = (CMD)(binary.BigEndian.Uint16(buf[24:26]))
	entry.Ty = (TYPE)(binary.BigEndian.Uint16(buf[26:28]) >> 8)      //取高8位直接截断
	entry.Encoding = (ENCODING)(binary.BigEndian.Uint16(buf[26:28])) // 直接截断
	return
}

// CheckCRC 核对 crc
func (e *LogEntry) CheckCRC(buf []byte) bool {
	// 更新新的crc
	crc := e.GetCRC(buf)
	return crc == e.CRC
}
