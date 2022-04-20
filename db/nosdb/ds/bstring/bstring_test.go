package bstring

import (
	"fmt"
	"testing"
)

func TestBString_IncrByInt(t *testing.T) {
	str := NewBString()
	val := "123"
	str.Set([]byte(val))
	fmt.Println(str.ToString())

	err := str.IncrByInt(-23)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(str.ToString())
	if str.ToString() != "100" {
		t.Error(" err res")
	}
}
