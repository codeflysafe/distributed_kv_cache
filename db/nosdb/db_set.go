/*
 * @Author: sjhuang
 * @Date: 2022-04-20 18:00:34
 * @LastEditTime: 2022-05-09 17:07:28
 * @FilePath: /nosdb/db_set.go
 */
package nosdb

import (
	"nosdb/ds"
	"nosdb/logfile"
)

// ---------------------------------------- set --------------------------------------------------------//

func (db *NosDB) lazySet() {
	db.setOnce.Do(func() {
		db.setIdx = NewSetIndex()
	})
}

// SAdd 命令将一个或多个成员元素加入到集合中，已经存在于集合的成员元素将被忽略。
// 假如集合 key 不存在，则创建一个只包含添加的元素作成员的集合。
func (db *NosDB) SAdd(key string, value []byte) {
	db.lazySet()
	db.setIdx.Lock()
	defer db.setIdx.Unlock()
	db.writeKVLog(key, value, SADD, logfile.SET)
	member := string(value)
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

// 判断是否为容器中的元素
func (db *NosDB) SIsMember(key string, value []byte) bool {
	db.lazySet()
	db.setIdx.RLock()
	defer db.setIdx.RUnlock()
	member := string(value)
	if obj, ok := db.setIdx.kv[key]; ok {
		return obj.SIsMember(member)
	}

	return false
}

// 移除并返回集合中的一个随机元素
// return nil if set is not exists or empty
// todo 日志 如何处理呢？
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
func (db *NosDB) SRem(key string, value []byte) {
	db.lazySet()
	db.setIdx.Lock()
	defer db.setIdx.Unlock()
	member := string(value)
	if obj, ok := db.setIdx.kv[key]; ok {
		db.writeKVLog(key, value, SREM, logfile.SET)
		obj.SRem(member)
	}
}
