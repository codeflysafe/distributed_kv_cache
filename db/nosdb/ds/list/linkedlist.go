package list

import (
	"errors"
)

var (
	// ErrListNotFound return when the list is not found
	ErrListNotFound = errors.New("not found in this list")

	// ErrIndexOutOfRange return when the index out of range
	ErrIndexOutOfRange = errors.New("the index out of range")
)

type CMD uint8
type DIRECTION uint8

const (
	LPush CMD = iota
	LPop
	RPush
	RPop

	LEFT DIRECTION = iota
	RIGHT
)

// 链表节点
type listNode struct {
	prev, next *listNode
	// 指向 value 的指针，存储各种类型的值
	// byte 流
	value []byte
}

// 创建新的节点
func newListNode(value []byte) *listNode {
	return &listNode{
		prev:  nil,
		next:  nil,
		value: value}
}

// 链表的迭代器
type ListIter struct {
	next      *listNode
	direction DIRECTION
}

// 使用 iterator 返回下一个节点
func (it *ListIter) Next() *listNode {
	cur := it.next
	if cur != nil {
		if it.direction == LEFT {
			it.next = cur.next
		} else {
			it.next = cur.prev
		}
	}
	return cur
}

// 双端链表
type List struct {
	// 头节点
	head *listNode
	// 尾节点
	tail *listNode
	// 链表的长度
	length uint64
}

// 新建一个链表
func NewList() *List {
	return &List{
		head:   nil,
		tail:   nil,
		length: 0,
	}
}

// 返回链表的长度
func (l *List) LLen() uint64 {
	return l.length
}

// 链表是否为空
func (l *List) Empty() bool {
	return l.length == 0
}

// 插入到链表头
// value unsafe.Pointer
func (l *List) LPush(value []byte) {
	node := newListNode(value)
	// 如果链表为空，
	if l.length == 0 {
		l.head = node
		l.tail = node
	} else {
		node.next = l.head
		l.head.prev = node
		l.head = node
	}
	l.length++
}

// 从链表left 弹出值
func (l *List) LPop() (value []byte) {
	if l.head != nil {
		value = l.head.value
		l.listDelNode(l.head)
	}
	return
}

// 插入到链表尾
func (l *List) RPush(value []byte) {
	node := newListNode(value)
	// 如果链表为空，
	if l.length == 0 {
		l.head = node
		l.tail = node
	} else {
		l.tail.next = node
		node.prev = l.tail
		l.tail = node
	}
	l.length++
}

// 从链表 right 弹出值
func (l *List) RPop() (value []byte) {
	if l.tail != nil {
		value = l.tail.value
		l.listDelNode(l.tail)
	}
	return
}

/* Return the element at the specified zero-based index
 * where 0 is the head, 1 is the element next to head
 * and so on. Negative integers are used in order to count
 * from the tail, -1 is the last element, -2 the penultimate
 * and so on. If the index is out of range NULL is returned. */
// from redis-3.0
func (l *List) listIndex(idx int64) *listNode {
	var node *listNode
	if idx < 0 {
		idx = (-idx) - 1
		node = l.tail
		for idx > 0 && node != nil {
			node = node.prev
		}
	} else {
		node = l.head
		for idx > 0 && node != nil {
			node = node.next
		}
	}
	return node
}

func (l *List) listDelNode(node *listNode) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}
	node = nil
	l.length--
}

// 删除 idx 处的链表节点
func (l *List) listDelIndex(idx int64) {
	node := l.listIndex(idx)
	if node != nil {
		l.listDelNode(node)
	}
}

// return the Iterator
func (l *List) ListIterator(direction DIRECTION) *ListIter {
	if direction == LEFT {
		return &ListIter{
			next:      l.head,
			direction: direction,
		}
	} else {
		return &ListIter{
			next:      l.tail,
			direction: direction,
		}
	}
}

// 返回头节点，但是不删除
func (l *List) LPeek() (value []byte) {
	if l.head != nil {
		value = l.head.value
	}
	return
}

// 返回尾节点, 但是不删除
func (l *List) RPeek() (value []byte) {
	if l.tail != nil {
		value = l.tail.value
	}
	return
}
