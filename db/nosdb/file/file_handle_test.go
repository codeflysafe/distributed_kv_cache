package file

import (
	"fmt"
	"os"
	"testing"
)

func TestIOFileHandle_WriteAt(t *testing.T) {
	file, err := os.OpenFile("test.txt", os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		t.Error(err)
	}

	handle, err := NewFileHandle(STANDARD_IO, file, -1)
	if err != nil {
		t.Error(err)
	}
	defer handle.Close()
	value := []byte("hsj is a bed man. and more time 杀人诛心!")
	offset, err := handle.WriteAt(1<<20-45, value)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(offset)
}

func TestIOFileHandle_ReadAt(t *testing.T) {
	file, err := os.OpenFile("test.txt", os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		t.Error(err)
	}

	handle, err := NewFileHandle(STANDARD_IO, file, -1)
	if err != nil {
		t.Error(err)
	}
	defer handle.Close()
	b, err := handle.ReadAt(1<<20-45, 45)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}

func TestMMapFileHandle_WriteAt(t *testing.T) {
	file, err := os.OpenFile("test.txt", os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		t.Error(err)
	}

	handle, err := NewFileHandle(M_MAP, file, -1)
	if err != nil {
		t.Error(err)
	}
	defer handle.Close()
	value := []byte("hsj is a bed man. and more time 杀人诛心!")
	offset, err := handle.WriteAt(1024, value)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(offset)
}

func TestMMapFileHandle_ReadAt(t *testing.T) {
	file, err := os.OpenFile("test.txt", os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		t.Error(err)
	}

	handle, err := NewFileHandle(M_MAP, file, -1)
	if err != nil {
		t.Error(err)
	}
	defer handle.Close()
	b, err := handle.ReadAt(1024, 45)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}
