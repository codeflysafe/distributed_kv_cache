package nosdb

import "nosdb/ds"

// ======================= string ====================

func (db *NosDB) lazyStr() {
	db.strOnce.Do(func() {
		db.strIdx = NewStringIndex()
	})
}

// 设置指定 key 的值。
// SET 命令用于设置给定 key 的值。
// 如果 key 已经存储其他值， SET 就覆写旧值，且无视类型。
// 如果不存在，则为新建一个
func (db *NosDB) Set(key string, value []byte) {
	db.lazyStr()
	db.strIdx.Lock()
	defer db.strIdx.Unlock()
	if obj, ok := db.strIdx.kv[key]; ok {
		obj.Set(value)
	} else {
		str := ds.NewString()
		str.Set(value)
		db.strIdx.kv[key] = str
	}
}

// 获取指定 key 的值。
// Get 命令用于获取指定 key 的值。如果 key 不存在，返回 nil
func (db *NosDB) Get(key string) []byte {
	db.lazyStr()
	db.strIdx.RLock()
	defer db.strIdx.RUnlock()
	if obj, ok := db.strIdx.kv[key]; ok {
		return obj.Get()
	}
	return nil
}

// 返回 key 中字符串值的子字符
// GetRange 命令用于获取指定子字符的值。如果 key 不存在，返回 nil
func (db *NosDB) GetRange(key string, start, end int) []byte {
	db.lazyStr()
	db.strIdx.RLock()
	defer db.strIdx.RUnlock()
	if obj, ok := db.strIdx.kv[key]; ok {
		return obj.GetRange(start, end)
	}
	return nil
}

// GetSet 命令用于设置指定 key 的值，并返回 key 的旧值。
// 如果不存在，则 set 后返回 nil
func (db *NosDB) GetSet(key string, newVal []byte) []byte {
	db.lazyStr()
	db.strIdx.Lock()
	defer db.strIdx.Unlock()
	if obj, ok := db.strIdx.kv[key]; ok {
		return obj.GetSet(newVal)
	} else {
		//
		str := ds.NewString()
		str.Set(newVal)
		db.strIdx.kv[key] = str
	}
	return nil
}

// （SET if Not eXists） 命令在指定的 key 不存在时，为 key 设置指定的值。
// 若存在， 则不处理
func (db *NosDB) SetNx(key string, value []byte) {
	db.lazyStr()
	db.strIdx.Lock()
	defer db.strIdx.Unlock()
	if _, ok := db.strIdx.kv[key]; !ok {
		str := ds.NewString()
		str.Set(value)
		db.strIdx.kv[key] = str
	}
}

func (db *NosDB) StrLen(key string) int {
	db.lazyStr()
	db.strIdx.RLock()
	defer db.strIdx.RUnlock()
	if obj, ok := db.strIdx.kv[key]; ok {
		return obj.StrLen()
	}
	return 0
}

// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
//如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
//本操作的值限制在 64 位(bit)有符号数字表示之内。
func (db *NosDB) IncrByInt(key string, offset int) error {
	db.lazyStr()
	db.strIdx.Lock()
	defer db.strIdx.Unlock()
	if obj, ok := db.strIdx.kv[key]; ok {
		return obj.IncrByInt(offset)
	} else {
		str := ds.NewString()
		err := str.IncrByInt(offset)
		if err != nil {
			return err
		}
		db.strIdx.kv[key] = str
		return nil
	}
}

func (db *NosDB) IncrByFloat(key string, offset float64) error {
	db.lazyStr()
	db.strIdx.Lock()
	defer db.strIdx.Unlock()
	if obj, ok := db.strIdx.kv[key]; ok {
		return obj.IncrByFloat(offset)
	} else {
		str := ds.NewString()
		err := str.IncrByFloat(offset)
		if err != nil {
			return err
		}
		db.strIdx.kv[key] = str
		return nil
	}
}

//  如果 key 已经存在并且是一个字符串， APPEND
// 命令将指定的 value 追加到该 key 原来值（value）的末尾。
// 如果 key 不存在， APPEND 就简单地将给定 key 设为 value
// 就像执行 SET key value 一样。
func (db *NosDB) Append(key string, value []byte) {
	db.lazyStr()
	db.strIdx.Lock()
	defer db.strIdx.Unlock()
	if obj, ok := db.strIdx.kv[key]; ok {
		obj.Append(value)
	} else {
		str := ds.NewString()
		str.Set(value)
		db.strIdx.kv[key] = str
		return
	}
}

// Del 删除指定 key 的value
func (db *NosDB) Del(key string) {
	db.lazyStr()
	db.strIdx.Lock()
	defer db.strIdx.Unlock()
	if _, ok := db.strIdx.kv[key]; ok {
		delete(db.hashIdx.kv, key)
	}
}
