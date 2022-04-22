package zskset

import (
	"fmt"
	"math/rand"
	"nosdb/ds"
)

const (
	// 最高层
	SKIPLIST_MAXLEVEL = 32
	// 随机数
	SKIPLIST_P = 2
)

// skipList 跳跃表
type skipListLevel struct {
	// 前进指针
	forward *skipListNode
	span    uint64
}

// https://redisbook.readthedocs.io/en/latest/internal-datastruct/skiplist.html
// 允许重复的 score 值：多个不同的 member 的 score 值可以相同。
// 进行对比操作时，不仅要检查 score 值，还要检查 member
type skipListNode struct {
	// 唯一的member， 比较时使用
	member string
	// the Value
	value []byte
	// score 排序时使用
	score float64
	// 后退指针
	backward *skipListNode
	// 层
	level []*skipListLevel
}

// 新建一个 skipListNode
func NewSkipListNode(level int, score float64, member string, value []byte) *skipListNode {
	node := &skipListNode{
		member:   member,
		value:    value,
		score:    score,
		backward: nil,
		level:    make([]*skipListLevel, level),
	}
	for i := range node.level {
		node.level[i] = new(skipListLevel)
	}
	return node
}

// 跳表
type SkipList struct {
	// 头指针
	head *skipListNode
	// 尾指针
	tail *skipListNode
	// 长度
	length uint64
	// 层高
	level int
}

// 创建 skiplist
func NewSkipList() *SkipList {
	return &SkipList{
		head:   NewSkipListNode(SKIPLIST_MAXLEVEL, 0, "", nil),
		tail:   nil,
		length: 0,
		level:  1,
	}
}

// 随机返回一个层高
func randomLevel() int {
	level := 1
	for (rand.Int() & 0xFFFF) < (0xFFFF >> SKIPLIST_P) {
		level++
	}
	if level > SKIPLIST_MAXLEVEL {
		level = SKIPLIST_MAXLEVEL
	}
	return level
}

// 根据 score 和 member 查找对应的 skipListNode
func (sk *SkipList) skListSearch(score float64, member string) *skipListNode {
	var i int
	var x *skipListNode
	x = sk.head
	for i = sk.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && (x.level[i].forward.score < score || (x.level[i].forward.score == score && x.level[i].forward.member < member)) {
			x = x.level[i].forward
		}
	}
	/* We may have multiple elements with the same score, what we need
	 * is to find the element with both the right score and member. */
	x = x.level[0].forward
	if x != nil && score == x.score && x.member == member {
		return x
	}
	return nil
}

// 插入 节点
// see redis-3.0 zslInsert
// return 插入的节点， 后面共 zset 使用
func (sk *SkipList) SkListInsert(score float64, member string, value []byte) *skipListNode {

	var x *skipListNode
	var update [SKIPLIST_MAXLEVEL]*skipListNode
	var rank [SKIPLIST_MAXLEVEL]uint64
	var i, level int

	x = sk.head
	for i = sk.level - 1; i >= 0; i-- {
		// store rank that is crossed to reach the insert position
		rank[i] = 0
		if i != sk.level-1 {
			rank[i] = rank[i+1]
		}
		for x.level[i].forward != nil && (x.level[i].forward.score < score || (x.level[i].forward.score == score && x.level[i].forward.member < member)) {
			rank[i] += x.level[i].span
			x = x.level[i].forward
		}
		update[i] = x
	}

	/* we assume the key is not already inside, since we allow duplicated
	 * scores, and the re-insertion of score and redis object should never
	 * happen since the caller of zslInsert() should test in the hash table
	 * if the element is already inside or not. */
	level = randomLevel()
	if level > sk.level {
		for i = sk.level; i < level; i++ {
			rank[i] = 0
			update[i] = sk.head
			update[i].level[i].span = sk.length
		}
		sk.level = level
	}

	x = NewSkipListNode(level, score, member, value)
	for i = 0; i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		// update span covered by update[i] as x is insert here
		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	/* increment span for untouched levels */
	for i = level; i < sk.level; i++ {
		update[i].level[i].span++
	}
	if update[0] != sk.head {
		x.backward = update[0]
	}

	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		sk.tail = x
	}

	sk.length++
	return x
}

// 删除节点
func (sk *SkipList) skListDeleteNode(x *skipListNode, update [SKIPLIST_MAXLEVEL]*skipListNode) {
	var i int
	// 删掉 x
	for i = 0; i < sk.level; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].span = x.level[i].span - 1
			update[i].level[i].forward = x.level[i].forward
		} else {
			update[i].level[i].span -= 1
		}
	}
	// 更新 tail 节点 和 backward
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x.backward
	} else {
		sk.tail = x.backward
	}

	for sk.level > 1 && sk.head.level[sk.level-1].forward == nil {
		sk.level--
	}
	sk.length--
}

// return 1 存在并且已经删除
// return 0 不存在
func (sk *SkipList) SkListDelete(score float64, member string) *skipListNode {
	var update [SKIPLIST_MAXLEVEL]*skipListNode
	var x *skipListNode
	var i int

	x = sk.head

	// 找到对应的 节点
	for i = sk.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && (x.level[i].forward.score < score || (x.level[i].forward.score == score && x.level[i].forward.member < member)) {
			x = x.level[i].forward
		}
		update[i] = x
	}
	/* We may have multiple elements with the same score, what we need
	 * is to find the element with both the right score and member. */
	x = x.level[0].forward
	if x != nil && score == x.score && x.member == member {
		// 删除这一列
		sk.skListDeleteNode(x, update)
	}
	return x
}

// 测试使用
func (sk *SkipList) skListPrintByLevel() {
	var x *skipListNode
	var i int
	// 每层遍历
	for i = sk.level - 1; i >= 0; i-- {
		x = sk.head
		fmt.Println(" level ", i)
		for x != nil {
			fmt.Println(string(x.value))
			x = x.level[i].forward
		}
	}
}

func skListValueGetMin(score float64, rangeSpec ds.ZRangeSpec) bool {
	if rangeSpec.MinEx {
		return score > rangeSpec.MinScore
	} else {
		return score >= rangeSpec.MinScore
	}
}

func skListValueGetMax(score float64, rangeSpec ds.ZRangeSpec) bool {
	if rangeSpec.MaxEx {
		return score < rangeSpec.MaxScore
	} else {
		return score <= rangeSpec.MaxScore
	}
}

func (sk *SkipList) skListIsInRange(rangeSpec ds.ZRangeSpec) bool {
	var x *skipListNode
	if (rangeSpec.MinScore > rangeSpec.MaxScore) || (rangeSpec.MinScore == rangeSpec.MaxScore && (rangeSpec.MinEx || rangeSpec.MaxEx)) {
		return false
	}
	x = sk.tail
	if x == nil || !skListValueGetMin(x.score, rangeSpec) {
		return false
	}
	x = sk.head.level[0].forward
	if x == nil || !skListValueGetMax(x.score, rangeSpec) {
		return false
	}
	return true
}

/* Find the first node that is contained in the specified range.
 * Returns NULL when no element is contained in the rangSpec. */
func (sk *SkipList) skListFirstInRange(rangSpec ds.ZRangeSpec) *skipListNode {
	var x *skipListNode
	var i int

	if !sk.skListIsInRange(rangSpec) {
		return nil
	}

	x = sk.head
	for i = sk.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && !skListValueGetMin(x.level[i].forward.score, rangSpec) {
			x = x.level[i].forward
		}
	}
	//
	x = x.level[0].forward
	if x == nil {
		return nil
	}
	if !skListValueGetMax(x.score, rangSpec) {
		return nil
	}
	return x
}

/* Find the last node that is contained in the specified range.
 * Returns NULL when no element is contained in the ZRangeSpec. */
func (sk *SkipList) skListLastInRange(rangSpec ds.ZRangeSpec) *skipListNode {
	var x *skipListNode
	var i int

	if !sk.skListIsInRange(rangSpec) {
		return nil
	}

	x = sk.head
	for i = sk.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && skListValueGetMax(x.level[i].forward.score, rangSpec) {
			x = x.level[i].forward
		}
	}
	//
	if x == nil {
		return nil
	}
	if !skListValueGetMin(x.score, rangSpec) {
		return nil
	}
	return x
}

// GetByScoreRange returns the nodes whose score within the specific range.
// If options is nil, it searches in interval [start, end] without any limit by default.
//
// Time complexity of this method is : O(log(N) + O(xxx)).
func (sk *SkipList) skListRange(rangSpec ds.ZRangeSpec) []*skipListNode {
	// 默认容量为一半，防止扩容过多
	nodes := make([]*skipListNode, 0, sk.length/2+1)
	first, last := sk.skListFirstInRange(rangSpec), sk.skListLastInRange(rangSpec)
	if first == nil || last == nil {
		return nodes
	}
	for first != nil && first != last {
		nodes = append(nodes, first)
		first = first.level[0].forward
	}
	nodes = append(nodes, last)
	return nodes
}
