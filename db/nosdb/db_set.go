/*
 * @Author: sjhuang
 * @Date: 2022-04-20 18:00:34
 * @LastEditTime: 2022-04-24 15:09:07
 * @FilePath: /nosdb/db_set.go
 */
package nosdb

import "nosdb/ds"

// ---------------------------------------- set --------------------------------------------------------//

func (db *NosDB) lazySet() {
	db.setOnce.Do(func() {
		if db.setIdx == nil {
			db.setIdx = NewSetIndex()
		}
	})
}

// SAdd 命令将一个或多个成员元素加入到集合中，已经存在于集合的成员元素将被忽略。
// 假如集合 key 不存在，则创建一个只包含添加的元素作成员的集合。
func (db *NosDB) SAdd(key string, member string, value []byte) {
	db.lazySet()
	db.setIdx.Lock()
	defer db.setIdx.Unlock()
	if obj, ok := db.setIdx.kv[key]; ok {
		obj.SAdd(member, value)
	} else {
		set := ds.NewSet()
		set.SAdd(member, value)
		db.setIdx.kv[key] = set
	}
}

// 取集合的成员数
// return 0 if set is not exists or empty
func (db *NosDB) SCard(key string) int {
	db.lazySet()
	db.setIdx.RLock()
	defer db.setIdx.RUnlock()
	if obj, ok := db.setIdx.kv[key]; ok {
		return obj.SCard()
	}
	return 0
}

// todo
func (db *NosDB) SIsMember(key string, member string) bool {
	db.lazySet()
	return false
}

// 移除并返回集合中的一个随机元素
// return nil if set is not exists or empty
func (db *NosDB) SPop(key string) []byte {
	db.lazySet()
	db.setIdx.Lock()
	defer db.setIdx.Unlock()
	if obj, ok := db.setIdx.kv[key]; ok {
		return obj.SPop()
	}
	return nil
}

// 移除集合中一个元素
// todo member 是啥？ 如何设计，这都是一个很大的问题
func (db *NosDB) SRem(key string, member string) {
	db.lazySet()
	db.setIdx.Lock()
	defer db.setIdx.Unlock()
	if obj, ok := db.setIdx.kv[key]; ok {
		obj.SRem(member)
	}
}
