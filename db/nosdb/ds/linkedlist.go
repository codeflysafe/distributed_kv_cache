package ds

import "unsafe"

// 链表节点
type listNode struct {
	prev, next *listNode
	// 指向 value 的指针，存储各种类型的值
	// todo 这里后面再研究一下
	value unsafe.Pointer
}

func newListNode(value unsafe.Pointer) *listNode {
	return &listNode{
		prev:  nil,
		next:  nil,
		value: value}
}

// 链表的迭代器
type listIter struct {
	next      *listNode
	direction uint8
}

// 双端链表
type list struct {
	// 头节点
	head *listNode
	// 尾节点
	tail *listNode
	// 链表的长度
	length uint64
}

// 新建一个链表
func NewList() *list {
	return &list{
		head:   nil,
		tail:   nil,
		length: 0,
	}
}

// 返回链表的长度
func (l *list) len() uint64 {
	return l.length
}

// 插入到链表头
// value unsafe.Pointer
func (l *list) listAddNodeHead(value unsafe.Pointer) {
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

func (l *list) listAddNodeTail(value unsafe.Pointer) {
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

// 删除节点
// node 待删除的节点
func (l *list) listDelNode(node *listNode) {

}

// 返回头节点
func (l *list) listFirst() *listNode {
	return l.head
}

// 返回尾节点
func (l *list) listLast() *listNode {
	return l.tail
}
