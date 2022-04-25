package nosdb

import (
	"fmt"
	"nosdb/ds"
)

// --------------------- list操作 ------------------

// lazyList 延迟初始化， listIdx
// 会存在`并发问题`
// 如果多个 goroutine 来同时读取调用这个方法，那么创建多个实例
// 这个时候，要使用 sync.Once
func (db *NosDB) lazyList() {
	db.listOnce.Do(func() {
		db.listIdx = NewListIndex()
	})
}

// LPush 向 key 对应的 list 内 add value
// 将一个值插入到列表头部
func (db *NosDB) LPush(key string, value []byte) {
	db.lazyList()
	db.listIdx.Lock()
	defer db.listIdx.Unlock()
	if obj, ok := db.listIdx.kv[key]; ok {
		// 如果存在
		obj.LPush(value)
	} else {
		// 如果不存在, 则创建
		db.listIdx.kv[key] = ds.NewList()
		db.listIdx.kv[key].LPush(value)
	}
}

// LPushX 将一个值插入到已存在的列表头部，列表不存在时操作无效。
func (db *NosDB) LPushX(key string, value []byte) {
	db.lazyList()
	db.listIdx.Lock()
	defer db.listIdx.Unlock()
	if obj, ok := db.listIdx.kv[key]; ok {
		// 如果存在
		obj.LPush(value)
	}
}

// LPop 移出并获取列表的第一个元素
// return nil if list is empty or not exist
func (db *NosDB) LPop(key string) []byte {
	db.lazyList()
	db.listIdx.Lock()
	defer db.listIdx.Unlock()
	if obj, ok := db.listIdx.kv[key]; ok {
		// 如果存在
		return obj.LPop()
	}
	return nil
}

// LPeek 获取列表的第一个元素并返回
// return nil if list is empty or not exist
func (db *NosDB) LPeek(key string) []byte {
	db.lazyList()
	db.listIdx.RLock()
	defer db.listIdx.RUnlock()
	if obj, ok := db.listIdx.kv[key]; ok {
		// 如果存在
		return obj.LPeek()
	}
	return nil
}

// RPush 向 key 对应的 list 内 add value
// 将一个值插入到列表尾部
func (db *NosDB) RPush(key string, value []byte) {
	db.lazyList()
	db.listIdx.Lock()
	defer db.listIdx.Unlock()
	if obj, ok := db.listIdx.kv[key]; ok {
		// 如果存在
		obj.RPush(value)
	} else {
		// 如果不存在, 则创建
		db.listIdx.kv[key] = ds.NewList()
		db.listIdx.kv[key].RPush(value)
	}
}

// RPushX 命令用于将一个值插入到已存在的列表尾部(最右边)。
// 如果列表不存在，操作无效。
func (db *NosDB) RPushX(key string, value []byte) {
	db.lazyList()
	db.listIdx.Lock()
	defer db.listIdx.Unlock()
	if obj, ok := db.listIdx.kv[key]; ok {
		// 如果存在
		obj.RPush(value)
	}
}

// RPop 移出并获取列表的最后一个元素
// return nil if list is empty or not exist
func (db *NosDB) RPop(key string) []byte {
	db.lazyList()
	db.listIdx.Lock()
	defer db.listIdx.Unlock()
	if obj, ok := db.listIdx.kv[key]; ok {
		// 如果存在
		return obj.RPop()
	}
	return nil
}

// RPeek 获取列表的最后一个元素并返回
// return nil if list is empty or not exist
func (db *NosDB) RPeek(key string) []byte {
	db.lazyList()
	db.listIdx.RLock()
	defer db.listIdx.RUnlock()
	if obj, ok := db.listIdx.kv[key]; ok {
		// 如果存在
		return obj.RPeek()
	}
	return nil
}

// LLen 获取列表长度
// return 0 if list is empty or not exist
func (db *NosDB) LLen(key string) int {
	db.lazyList()
	db.listIdx.RLock()
	defer db.listIdx.RUnlock()
	if obj, ok := db.listIdx.kv[key]; ok {
		// 如果存在
		return obj.LLen()
	}
	return 0
}

// LSet 通过索引来设置元素的值。
// 当索引参数超出范围，或对一个空列表进行 LSet 时，返回一个错误。
func (db *NosDB) LSet(key string, idx int, value []byte) (err error) {
	db.lazyList()
	db.listIdx.Lock()
	defer db.listIdx.Unlock()
	if obj, ok := db.listIdx.kv[key]; ok {
		// 如果存在
		err = obj.ListSet(idx, value)
		return
	}
	err = fmt.Errorf("list not exist")
	return
}

// LRem 根据参数 count 的值，移除列表中与参数 value 相等的元素。
// count 的值可以是以下几种：
// count > 0 : 从表头开始向表尾搜索，移除与 value 相等的元素，数量为 count
// count < 0 : 从表尾开始向表头搜索，移除与 value 相等的元素，数量为 count 的绝对值。
// count = 0 : 移除表中所有与 VALUE 相等的值。
// todo
func (db *NosDB) LRem(key string, count int, value []byte) {
	db.lazyList()
}

// LTrim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
// 下标 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
// todo
func (db *NosDB) LTrim(key string, start, end int) {
	db.lazyList()
}

// LRange 返回列表中指定区间内的元素，区间以偏移量 START 和 END 指定。
// 其中 0 表示列表的第一个元素， 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
// todo
func (db *NosDB) LRange(key string, start, end int) {
	db.lazyList()
	db.listIdx.RLock()
	defer db.listIdx.RUnlock()

}

// LIndex 命令用于通过索引获取列表中的元素。你也可以使用负数下标，
// 以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
// error if idx out of range or list is not exists
func (db *NosDB) LIndex(key string, idx int) ([]byte, error) {
	db.lazyList()
	db.listIdx.RLock()
	defer db.listIdx.RUnlock()
	if obj, ok := db.listIdx.kv[key]; ok {
		// 如果存在
		return obj.ListSeek(idx)
	}
	err := fmt.Errorf(" list is not exists key %s", key)
	return nil, err
}
