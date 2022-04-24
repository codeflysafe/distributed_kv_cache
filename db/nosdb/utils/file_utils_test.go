package utils

import (
	"fmt"
	"testing"
)

func TestCheckDir(t *testing.T) {
	path := "./log"
	err := CheckDir(path)
	fmt.Println(err)
}

func TestOpen(t *testing.T) {
	path := "./log"
	fileName := "activate_123123.dat"
	f, err := Open(path, fileName)
	if err != nil {
		t.Error(err)
	}
	println(f.Name())
}

func TestReNameFile(t *testing.T) {
	path := "./log"
	fileName := "activate_123123.dat"
	err := ReNameFile(path, fileName, "wal_123123.dat")
	if err != nil {
		t.Error(err)
	}
}
