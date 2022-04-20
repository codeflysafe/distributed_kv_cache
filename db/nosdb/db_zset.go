package nosdb

import "nosdb/ds"

// ==================== zset ==================

// Zadd 命令用于将一个或多个成员元素及其分数值加入到有序集当中。
//如果某个成员已经是有序集的成员，那么更新这个成员的分数值，并通过重新插入这个成员元素，来保证该成员在正确的位置上。
//分数值可以是整数值或双精度浮点数。
//如果有序集合 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
//当 key 存在但不是有序集类型时，返回一个错误。
func (db *NosDB) ZAdd(key string, score float64, member string, value []byte) {
	db.zSetIdx.Lock()
	defer db.zSetIdx.Unlock()
	if obj, ok := db.zSetIdx.kv[key]; ok {
		obj.ZAdd(score, member, value)
	} else {
		zs := ds.NewZSet()
		zs.ZAdd(score, member, value)
		db.zSetIdx.kv[key] = zs
	}
}

// 获取有序集合的成员数
func (db *NosDB) ZCard(key string) int {
	db.zSetIdx.RLock()
	defer db.zSetIdx.RUnlock()
	if obj, ok := db.zSetIdx.kv[key]; ok {
		return obj.ZCard()
	}
	return 0
}

// Zcount 命令用于计算有序集合中指定分数区间的成员数量。
func (db *NosDB) ZCount(key string, minScore, maxSCore float64) int {
	db.zSetIdx.RLock()
	defer db.zSetIdx.RUnlock()
	if obj, ok := db.zSetIdx.kv[key]; ok {
		return obj.ZCount(ds.ZRangeSpec{minScore, maxSCore, false, false})
	}
	return 0
}

// Zincrby 命令对有序集合中指定成员的分数加上增量 increment
//可以通过传递一个负数值 increment ，让分数减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5 。
//当 key 不存在，或分数不是 key 的成员时， ZINCRBY key increment member 等同于 ZADD key increment member 。
//当 key 不是有序集类型时，返回一个错误。
//分数值可以是整数值或双精度浮点数。
func (db *NosDB) ZIncrBy(key string, member string, offset float64) {
	db.zSetIdx.Lock()
	defer db.zSetIdx.Unlock()
	if obj, ok := db.zSetIdx.kv[key]; ok {
		return obj.ZAdd()
	}
	return 0
}
