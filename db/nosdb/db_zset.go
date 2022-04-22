package nosdb

import "nosdb/ds"

// ==================== ZSet ==================

func (db *NosDB) lazyZSet() {
	db.zSetOnce.Do(func() {
		if db.zSetIdx == nil {
			db.zSetIdx = NewZSetIndex()
		}
	})
}

// ZAdd 命令用于将一个或多个成员元素及其分数值加入到有序集当中。
//如果某个成员已经是有序集的成员，那么更新这个成员的分数值，并通过重新插入这个成员元素，来保证该成员在正确的位置上。
//分数值可以是整数值或双精度浮点数。
//如果有序集合 key 不存在，则创建一个空的有序集并执行 ZAdd 操作
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

// ZCard 获取有序集合的成员数
func (db *NosDB) ZCard(key string) int {
	db.zSetIdx.RLock()
	defer db.zSetIdx.RUnlock()
	if obj, ok := db.zSetIdx.kv[key]; ok {
		return obj.ZCard()
	}
	return 0
}

// ZCount 命令用于计算有序集合中指定分数区间的成员数量。
func (db *NosDB) ZCount(key string, minScore, maxSCore float64) int {
	db.zSetIdx.RLock()
	defer db.zSetIdx.RUnlock()
	if obj, ok := db.zSetIdx.kv[key]; ok {
		return obj.ZCount(ds.ZRangeSpec{minScore, maxSCore, false, false})
	}
	return 0
}

//ZIncrScore 命令对有序集合中指定成员的分数加上增量 offset
//可以通过传递一个负数值 increment ，让分数减去相应的值，比如 ZIncrScore key -5 member ，就是让 member 的 score 值减去 5 。
//当 key 不存在，或分数不是 key 的成员时， ZIncrScore key increment member 等同于 ZAdd key increment member 。
//分数值可以是整数值或双精度浮点数。
func (db *NosDB) ZIncrScore(key string, member string, value []byte, offset float64) {
	db.zSetIdx.Lock()
	defer db.zSetIdx.Unlock()
	if obj, ok := db.zSetIdx.kv[key]; ok {
		obj.ZIncrScore(member, value, offset)
	} else {
		zs := ds.NewZSet()
		zs.ZAdd(offset, member, value)
		db.zSetIdx.kv[key] = zs
	}
}

// ZScore 命令返回有序集中，成员的分数值。
//如果成员元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
func (db *NosDB) ZScore(key string, member string) float64 {
	db.zSetIdx.RLock()
	defer db.zSetIdx.RUnlock()
	if obj, ok := db.zSetIdx.kv[key]; ok {
		return obj.ZScore(member)
	}
	return 0
}
