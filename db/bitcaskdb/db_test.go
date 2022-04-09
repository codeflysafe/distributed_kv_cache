package bitcaskdb

import (
	"fmt"
	"testing"
)

const (
	path = "./db_test"
)

func TestOpen(t *testing.T) {
	db, err := Open(path)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	t.Log(db)
}

func TestBitcaskDB_Put(t *testing.T) {
	db, err := Open(path)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	for i := 0; i < 10; i++ {
		err = db.Put([]byte(fmt.Sprintf("key%d", i)), []byte(fmt.Sprintf("value%d", i)))
		if err != nil {
			t.Error(err)
		}
	}
}

func TestBitcaskDB_Get(t *testing.T) {
	db, err := Open(path)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	for i := 1; i < 10; i++ {
		key := []byte(fmt.Sprintf("key%d", i))
		value := fmt.Sprintf("value%d", i)
		v, err := db.Get(key)
		if err != nil {
			t.Error(err)
		}
		if value != string(v) {
			t.Error(fmt.Errorf("doesn't match %s, %s", value, string(v)))
		} else {
			fmt.Println(string(key), value)
		}
	}
}

func TestBitcaskDB_Del(t *testing.T) {
	db, err := Open(path)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	key := []byte("key0")
	err = db.Del(key)
	if err != nil {
		t.Error(err)
	}
	v, _ := db.Get(key)
	if v != nil {
		t.Error(fmt.Sprintf("error"))
	}

	keys, _ := db.ListKeys()
	println(keys)
}

func TestBitcaskDB_Merge(t *testing.T) {
	db, err := Open(path)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	err = db.Merge()
	if err != nil {
		t.Error(err)
	}
	var v []byte
	v, err = db.Get([]byte("key1"))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(v))
}
