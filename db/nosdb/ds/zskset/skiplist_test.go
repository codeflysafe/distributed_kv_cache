package zskset

import (
	"fmt"
	"testing"
)

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

func TestSkLipList_skListFirstInRange(t *testing.T) {
	sk := NewSkipList()
	sk.SkListInsert(10.0, "sjhuang", []byte("1993"))
	sk.SkListInsert(11.0, "sjhuang", []byte("1994"))
	sk.SkListInsert(11.0, "sjhuang2", []byte("1995"))
	sk.SkListInsert(12.0, "sjhuang2", []byte("1996"))
	sk.SkListInsert(12.0, "sjhuang1", []byte("1997"))
	sk.SkListInsert(12.0, "sjhuang", []byte("1998"))
	// 最后一层遍历一下即可
	sk.skListPrintByLevel()
	rangSpec := ZRangeSpec{
		MinScore: 10.0,
		MaxScore: 11.9,
		MinEx:    false,
		MaxEx:    false,
	}
	node := sk.skListFirstInRange(rangSpec)
	if node == nil || node.score != rangSpec.MinScore {
		if node != nil {
			t.Errorf(" ans is error ,except %f get %f", rangSpec.MinScore, node.score)
		} else {
			t.Errorf(" ans is error ")
		}
	}
	fmt.Println(node)

}

func TestSkLipList_skListLastInRange(t *testing.T) {
	sk := NewSkipList()
	sk.SkListInsert(10.0, "sjhuang", []byte("1993"))
	sk.SkListInsert(11.0, "sjhuang", []byte("1994"))
	sk.SkListInsert(11.0, "sjhuang2", []byte("1995"))
	sk.SkListInsert(12.0, "sjhuang2", []byte("1996"))
	sk.SkListInsert(12.0, "sjhuang1", []byte("1997"))
	sk.SkListInsert(12.0, "sjhuang", []byte("1998"))
	// 最后一层遍历一下即可
	sk.skListPrintByLevel()
	rangSpec := ZRangeSpec{
		MinScore: 10.0,
		MaxScore: 11.9,
		MinEx:    false,
		MaxEx:    false,
	}
	node := sk.skListLastInRange(rangSpec)
	if node == nil || node.score != 11.0 {
		if node != nil {
			t.Errorf(" ans is error ,except %f get %f", rangSpec.MinScore, node.score)
		} else {
			t.Errorf(" ans is error ")
		}
	}
	fmt.Println(node)
}

func TestSkLipList_skListRange(t *testing.T) {
	sk := NewSkipList()
	sk.SkListInsert(10.0, "sjhuang", []byte("1993"))
	sk.SkListInsert(11.0, "sjhuang", []byte("1994"))
	sk.SkListInsert(11.0, "sjhuang2", []byte("1995"))
	sk.SkListInsert(12.0, "sjhuang2", []byte("1996"))
	sk.SkListInsert(12.0, "sjhuang1", []byte("1997"))
	sk.SkListInsert(12.0, "sjhuang", []byte("1998"))
	// 最后一层遍历一下即可
	sk.skListPrintByLevel()
	rangSpec := ZRangeSpec{
		MinScore: 29.0,
		MaxScore: 30.0,
		MinEx:    false,
		MaxEx:    false,
	}
	nodes := sk.skListRange(rangSpec)
	for _, node := range nodes {
		fmt.Println(node)
	}
}
