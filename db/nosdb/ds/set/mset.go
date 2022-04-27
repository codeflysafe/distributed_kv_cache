/*
 * @Author: sjhuang
 * @Date: 2022-04-19 09:30:36
 * @LastEditTime: 2022-04-27 10:43:25
 * @FilePath: /nosdb/ds/set/mset.go
 */
package set

const (
	DEFAULT_MSET_SIZE int = 8
)

// 使用map构成的set
type MSet struct {
	Items map[string][]byte
}

// 新建一个 MSet
func NewMSet() *MSet {
	return &MSet{
		Items: make(map[string][]byte, DEFAULT_MSET_SIZE),
	}
}

// 向 MSet 集合中添加元素
// 如果存在就更新，不存在就添加
func (ms *MSet) SAdd(member string, value []byte) {
	ms.Items[member] = value
}

// 返回结合的长度
func (ms *MSet) SCard() int {
	return len(ms.Items)
}

func (ms *MSet) Empty() bool {
	return len(ms.Items) == 0
}

// 	SPOP key
// 移除并返回集合中的一个随机元素
func (ms *MSet) SPop() []byte {
	for k, v := range ms.Items {
		delete(ms.Items, k)
		return v
	}
	return nil
}

// 移除集合中一个或多个成员
func (ms *MSet) SRem(member string) {
	delete(ms.Items, member)
}

func (ms *MSet) SIsMember(member string) bool {
	_, ok := ms.Items[member]
	return ok
}
