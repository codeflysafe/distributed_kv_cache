package list

import "testing"

func TestSkipList_SkListInsert(t *testing.T) {
	sk := NewSkipList()
	sk.SkListInsert(10.0, "sjhuang", []byte("1993"))
	sk.SkListInsert(11.0, "sjhuang", []byte("1994"))
	sk.SkListInsert(11.0, "sjhuang2", []byte("1995"))
	sk.SkListInsert(12.0, "sjhuang2", []byte("1996"))
	sk.SkListInsert(12.0, "sjhuang1", []byte("1997"))
	sk.SkListInsert(12.0, "sjhuang", []byte("1998"))

	// 最后一层遍历一下即可
	sk.skListPrintByLevel()
}

func TestSkipList_SkListDelete(t *testing.T) {
	sk := NewSkipList()
	sk.SkListInsert(10.0, "sjhuang", []byte("1993"))
	sk.SkListInsert(11.0, "sjhuang", []byte("1994"))
	sk.SkListInsert(11.0, "sjhuang2", []byte("1995"))
	sk.SkListInsert(12.0, "sjhuang2", []byte("1996"))
	sk.SkListInsert(12.0, "sjhuang1", []byte("1997"))
	sk.SkListInsert(12.0, "sjhuang", []byte("1998"))
	// 最后一层遍历一下即可
	sk.skListPrintByLevel()
	sk.SkListDelete(10.0, "sjhuang")

	// 最后一层遍历一下即可
	sk.skListPrintByLevel()
}
