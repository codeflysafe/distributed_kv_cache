package zskset

import "nosdb/ds"

// 有序集合数据结构
type ZSkSet struct {
	// key value
	items map[string]*skipListNode
	// 跳表实现的有序集合
	list *SkipList
}

func NewZSet() *ZSkSet {
	return &ZSkSet{
		items: make(map[string]*skipListNode),
		list:  NewSkipList(),
	}
}

// 返回 zskset 中的元素个数
func (zs *ZSkSet) ZCard() int {
	return len(zs.items)
}

// 像ZSet添加元素
// 如果存在则更新
func (zs *ZSkSet) ZAdd(score float64, member string, value []byte) {
	var node *skipListNode
	// 如果已经存在
	if v, ok := zs.items[member]; ok {
		if v.score == score {
			v.value = value
		} else {
			zs.list.SkListDelete(score, member)
			node = zs.list.SkListInsert(score, member, value)
		}
	} else {
		// 不存在
		node = zs.list.SkListInsert(score, member, value)
	}

	if node != nil {
		zs.items[member] = node
	}
}

// 删除
func (zs *ZSkSet) ZDel(score float64, member string) {
	zs.list.SkListDelete(score, member)
}

// 查询
// Todo
func (zs *ZSkSet) ZRange(rangSpec ds.ZRangeSpec) {
	zs.list.skListRange(rangSpec)
}

func (zs *ZSkSet) ZCount(rangSpec ds.ZRangeSpec) int {
	return len(zs.list.skListRange(rangSpec))
}

func (zs *ZSkSet) ZIncrScore(member string, value []byte, offset float64) {
	if node, ok := zs.items[member]; ok {
		// 删除 再插入
		score := node.score + offset
		value := node.value
		zs.list.SkListDelete(node.score, member)
		delete(zs.items, member)
		zs.items[member] = zs.list.SkListInsert(score, member, value)
	} else {
		zs.items[member] = zs.list.SkListInsert(offset, member, value)
	}
}

func (zs *ZSkSet) ZScore(member string) float64 {
	if node, ok := zs.items[member]; ok {
		return node.score
	}
	return 0
}
