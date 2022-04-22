package mmap

import (
	"fmt"
	"os"
	"testing"
)

func TestMMap_WriteAt(t *testing.T) {
	file, err := os.OpenFile("test.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0664)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	length := 1 << 8
	data := make([]byte, length, length)
	file.Write(data)
	mm, err := NewMMap(file, RDWR, length)
	if err != nil {
		t.Error(err)
	}
	value := []byte("hsj good man! says xiao wang!")
	fmt.Println(value)
	mm.WriteAt(0, value)
	err = mm.MSync()
	if err != nil {
		t.Error(err)
	}

	// 这里会覆盖
	value = []byte("whh good girl! says xiao huang!")
	fmt.Println(value)
	mm.WriteAt(10, value)
	err = mm.MSync()
	if err != nil {
		t.Error(err)
	}
}

func TestMMap_ReadAt(t *testing.T) {
	file, err := os.OpenFile("test.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0664)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	length := 1 << 8
	mm, err := NewMMap(file, RDWR, length)
	defer mm.MunMap()
	if err != nil {
		t.Error(err)
	}

	b, er := mm.ReadAt(0, 10)
	if er != nil {
		t.Error(er)
	}
	fmt.Println(string(b))
}
