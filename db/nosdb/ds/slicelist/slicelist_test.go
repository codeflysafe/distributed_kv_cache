package slicelist

import (
	"fmt"
	"strconv"
	"testing"
)

func Key(str string) []byte {
	return []byte(str)
}

func TestList_LPush(t *testing.T) {
	list := NewSliceList()
	val := "123"
	list.LPush(Key(val))

	if list.LLen() != 1 {
		t.Error(" 没有数据，测试失败")
	}

	var lens int = 10000
	for i := 0; i < 10000; i++ {
		list.LPush(Key(strconv.Itoa(i)))
	}

	if list.LLen() != lens+1 {
		t.Error(" 没有数据，测试失败")
	}
}

func TestList_LPop(t *testing.T) {
	list := NewSliceList()
	val := "123"
	list.LPush(Key(val))
	vs := list.LPop()
	fmt.Println(string(vs))
	if string(vs) != val {
		t.Errorf(" excepct get %s, bug get %s", val, string(vs))
	}
	if !list.Empty() {
		t.Errorf(" length should is 0")
	}
}

func TestList_RPush(t *testing.T) {
	list := NewSliceList()
	val := "123"
	list.RPush(Key(val))

	if list.LLen() != 1 {
		t.Error(" 没有数据，测试失败")
	}

	var lens int = 10000
	for i := 0; i < 10000; i++ {
		list.RPush(Key(strconv.Itoa(i)))
	}

	if list.LLen() != lens+1 {
		t.Error(" 没有数据，测试失败")
	}
}

func TestList_RPop(t *testing.T) {
	list := NewSliceList()
	val := "123"
	list.RPush(Key(val))
	vs := list.RPop()
	fmt.Println(string(vs))
	if string(vs) != val {
		t.Errorf(" excepct get %s, bug get %s", val, string(vs))
	}
	if !list.Empty() {
		t.Errorf(" length should is 0")
	}
}

func TestList_LPeek(t *testing.T) {
	list := NewSliceList()
	val := "123"
	list.LPush(Key(val))
	vs := list.LPeek()
	fmt.Println(string(vs))
	if string(vs) != val {
		t.Errorf(" excepct get %s, bug get %s", val, string(vs))
	}
	if list.LLen() != 1 {
		t.Errorf(" length should is 1")
	}
}

func TestList_RPeek(t *testing.T) {
	list := NewSliceList()
	val := "123"
	list.RPush(Key(val))
	vs := list.RPeek()
	fmt.Println(string(vs))
	if string(vs) != val {
		t.Errorf(" excepct get %s, bug get %s", val, string(vs))
	}
	if list.LLen() != 1 {
		t.Errorf(" length should is 1")
	}
}

func TestList_ListSeek(t *testing.T) {
	list := NewSliceList()
	for i := 0; i < 100; i++ {
		v := fmt.Sprintf("%d", i)
		list.RPush(Key(v))
	}

	if list.LLen() != 100 {
		t.Errorf(" error lens %d", list.LLen())
	}

	val := list.ListSeek(30)
	fmt.Println("val is", string(val))
	if string(val) != "30" {
		t.Errorf(" error value %s", string(val))
	}

}

func TestList_ListDelIndex(t *testing.T) {
	list := NewSliceList()
	for i := 0; i < 100; i++ {
		v := fmt.Sprintf("%d", i)
		list.RPush(Key(v))
	}

	if list.LLen() != 100 {
		t.Errorf(" error lens %d", list.LLen())
	}

	val := list.ListSeek(30)
	fmt.Println("val is", string(val))
	if string(val) != "30" {
		t.Errorf(" error value %s", string(val))
	}

	list.ListDelIndex(30)
	val = list.ListSeek(30)
	fmt.Println("val is", string(val))
	if string(val) != "31" {
		t.Errorf(" error value %s", string(val))
	}
}
