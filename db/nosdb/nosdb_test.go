package nosdb

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNosDB_Set(t *testing.T) {
	db := NewNosDB()
	val := []byte("value1")
	db.Set("key1", val)
}

func TestNosDB_Get(t *testing.T) {
	db := NewNosDB()
	val := []byte("value1")
	db.Set("key1", val)
	v := db.Get("key1")
	fmt.Println(string(v))
	if string(v) != string(val) {
		t.Error("err value")
	}
}

func TestNosDB_GetSet(t *testing.T) {
	db := NewNosDB()
	val := []byte("value1")
	v := db.GetSet("key1", val)
	if v != nil {
		panic("err")
	}
	val2 := []byte("value2")
	v1 := db.GetSet("key1", val2)
	fmt.Println(string(v1))
	if reflect.DeepEqual(v1, val) {
		fmt.Println(string(v1))
	}
}

func TestNosDB_SetNx(t *testing.T) {
	db := NewNosDB()
	val := []byte("value1")
	db.SetNx("key1", val)
	v := db.Get("key1")
	if v == nil || string(v) != "value1" {
		t.Error(" errr setNx")
	}
	fmt.Println(string(v))
}
