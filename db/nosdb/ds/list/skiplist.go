package list

import (
	"math/rand"
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

func (sk *SkipList) SearchNode(score float64, member string) *skipListNode {
	return nil
}

// 插入 节点
// see redis-3.0 zslInsert
func (sk *SkipList) InsertNode(score float64, member string, value []byte) *skipListNode {

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
