package mmap

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

func TestMMap_WriteAt(t *testing.T) {
	file, err := os.Create("test.txt")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	length := 1 << 20
	data := make([]byte, length, length)
	file.Write(data)
	mm, err := NewMMap(file, syscall.PROT_WRITE, length)
	if err != nil {
		t.Error(err)
	}
	value := []byte("hsj good man! says xiao wang!")
	fmt.Println(value)
	mm.WriteAt(0, value)
	err = mm.MunMap()
	if err != nil {
		t.Error(err)
	}
}
