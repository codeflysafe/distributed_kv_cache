package nosdb

import (
	"nosdb/ds"
)

// ------------------ hash 操作 -----------------------

// 懒加载 hash 结构
func (db *NosDB) lazyHash() {
	db.hashOnce.Do(func() {
		db.hashIdx = NewHashIndex()
	})
}

//HSet 命令用于为哈希表中的字段赋值 。
//如果哈希表不存在，一个新的哈希表被创建并进行 HSET 操作。
//如果字段已经存在于哈希表中，旧值将被覆盖。
func (db *NosDB) HSet(key string, member string, value []byte) {
	db.lazyHash()
	db.hashIdx.Lock()
	defer db.hashIdx.Unlock()
	if obj, ok := db.hashIdx.kv[key]; ok {
		obj.HSet(member, value)
	} else {
		h := ds.NewHash()
		h.HSet(member, value)
		db.hashIdx.kv[key] = h
	}
}

// HSetNx 命令用于为哈希表中不存在的的字段赋值 。
//如果哈希表不存在，一个新的哈希表被创建并进行 HSet 操作。
//如果字段已经存在于哈希表中，操作无效。
//如果 key 不存在，一个新哈希表被创建并执行 HSetNx 命令。
func (db *NosDB) HSetNx(key string, member string, value []byte) {
	db.lazyHash()
	db.hashIdx.Lock()
	defer db.hashIdx.Unlock()
	if obj, ok := db.hashIdx.kv[key]; ok {
		obj.HSetNx(member, value)
	} else {
		h := ds.NewHash()
		h.HSet(member, value)
		db.hashIdx.kv[key] = h
	}
}

// HDel 删除一个哈希表字段
// 如果不存在该字段或者hash表为空，则 no op
func (db *NosDB) HDel(key string, member string) {
	db.lazyHash()
	db.hashIdx.Lock()
	defer db.hashIdx.Unlock()
	if obj, ok := db.hashIdx.kv[key]; ok {
		obj.HDel(member)
	}
}

// HExists 查看哈希表 key 中，指定的字段是否存在。
func (db *NosDB) HExists(key string, member string) bool {
	db.lazyHash()
	db.hashIdx.RLock()
	defer db.hashIdx.RUnlock()
	if obj, ok := db.hashIdx.kv[key]; ok {
		return obj.HExists(member)
	}
	return false
}

// HGet 获取存储在哈希表中指定字段的值。
// 返回 nil， 如果为空或者不存在
func (db *NosDB) HGet(key string, member string) []byte {
	db.lazyHash()
	db.hashIdx.RLock()
	defer db.hashIdx.RUnlock()
	if obj, ok := db.hashIdx.kv[key]; ok {
		return obj.HGet(member)
	}
	return nil
}

// HIncrBy 为哈希表 key 中的指定字段 member 的整数值加上增量 offset 。
// 返回错误，如果不存在此字段，或者为空，或者不是整数
func (db *NosDB) HIncrBy(key string, member string, offset int) error {
	db.lazyHash()
	db.hashIdx.Lock()
	defer db.hashIdx.Unlock()
	if obj, ok := db.hashIdx.kv[key]; ok {
		return obj.HIncrBy(member, offset)
	} else {
		h := ds.NewHash()
		err := h.HIncrBy(member, offset)
		db.hashIdx.kv[key] = h
		return err
	}
}

// HIncrByFloat 为哈希表 key 中的指定字段的浮点数值加上增量 offset 。
// 返回错误，如果不存在此字段，或者为空，或者不是整数
func (db *NosDB) HIncrByFloat(key string, member string, offset float64) error {
	db.lazyHash()
	db.hashIdx.Lock()
	defer db.hashIdx.Unlock()
	if obj, ok := db.hashIdx.kv[key]; ok {
		return obj.HIncrByFloat(member, offset)
	} else {
		h := ds.NewHash()
		err := h.HIncrByFloat(member, offset)
		db.hashIdx.kv[key] = h
		return err
	}
}

// HLen 获取哈希表中字段的数量
func (db *NosDB) HLen(key string) int {
	db.lazyHash()
	db.hashIdx.RLock()
	defer db.hashIdx.RUnlock()
	return len(db.hashIdx.kv)
}
