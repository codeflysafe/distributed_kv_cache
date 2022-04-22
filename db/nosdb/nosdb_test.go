package nosdb

import (
	"fmt"
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
}
