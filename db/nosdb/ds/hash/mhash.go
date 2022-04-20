package hash

import (
	"nosdb/utils"
)

// 采用 map 作为hash table

const (
	DEFAULT_MSET_SIZE int = 8
)

// 使用map构成的set
type MHash struct {
	Items map[string][]byte
}

// 新建一个 MSet
func NewMHash() *MHash {
	return &MHash{
		Items: make(map[string][]byte, DEFAULT_MSET_SIZE),
	}
}

// 删除一个哈希表字段
func (mh *MHash) HDel(member string) {
	delete(mh.Items, member)
}

//查看哈希表 member 中，指定的字段是否存在。
func (mh *MHash) HExists(member string) bool {
	v := mh.Items[member]
	return v != nil
}

// 获取存储在哈希表中指定字段的值
func (mh *MHash) HGet(member string) []byte {
	return mh.Items[member]
}

// 获取哈希表中字段的数量
func (mh *MHash) HLen(member string) int {
	return len(mh.Items)
}

// 将哈希表 member 中的字段 field 的值设为 value 。
func (mh *MHash) HSet(member string, value []byte) {
	mh.Items[member] = value
}

// 哈希表中不存在的的字段赋值 。
// 设置成功，返回 1 。 如果给定字段已经存在且没有操作被执行，返回 0 。
func (mh *MHash) HSetNx(member string, value []byte) {
	// 如果存在，则直接返回
	if mh.HExists(member) {
		return
	}
	// 不存在则创建
	mh.HSet(member, value)
}

// Redis Hincrby 命令用于为哈希表中的字段值加上指定增量值。
//增量也可以为负数，相当于对指定字段进行减法操作。
//如果哈希表的 member 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。
//如果指定的字段不存在，那么在执行命令前，字段的值被初始化为 0 。
//对一个储存字符串值的字段执行 HINCRBY 命令将造成一个错误。
//本操作的值被限制在 64 位(bit)有符号数字表示之内。
func (mh *MHash) HIncrBy(member string, offset int) error {
	mh.HSetNx(member, []byte("0"))
	v := mh.Items[member]
	vs, err := utils.BytesIncrBy(v, offset)
	if err != nil {
		return err
	}
	mh.HSet(member, vs)
	return nil
}

//Redis Hincrbyfloat 命令用于为哈希表中的字段值加上指定浮点数增量值。
//如果指定的字段不存在，那么在执行命令前，字段的值被初始化为 0 。
func (mh *MHash) HIncrByFloat(member string, offset float64) (err error) {
	mh.HSetNx(member, []byte("0"))
	v := mh.Items[member]
	var newV []byte
	newV, err = utils.ByteIncrByFloat(v, offset)
	if err != nil {
		return
	}
	mh.HSet(member, newV)
	return
}
