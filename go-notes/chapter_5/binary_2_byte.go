package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type T struct {
	B float64
	C int64
}

type A struct {
	Name   string                 // 字符串
	Id     int                    // 整数
	Score  float64                // 浮点数
	Desc   map[string]interface{} // map
	Prices []int                  // slice
}

func (a A) ToString() string {
	return fmt.Sprint(a)
}

func main() {
	a := A{
		Name:   "sjhuang",
		Id:     28,
		Score:  999.0,
		Desc:   make(map[string]interface{}, 10),
		Prices: make([]int, 0, 2),
	}

	a.Desc["city"] = "Henan"
	a.Desc["age"] = 1000
	a.Prices = append(a.Prices, 12)

	t := T{B: 3.1415, C: 0xFFFFFFFF}

	buf := &bytes.Buffer{}
	// 大端优先
	// Write writes the binary representation of data into w.
	// Data must be a fixed-size value or a slice of fixed-size
	// values, or a pointer to such data.
	err := binary.Write(buf, binary.BigEndian, t)
	// error
	// err := binary.Write(buf, binary.BigEndian, a)
	if err != nil {
		panic(err)
	}

	fmt.Println(buf.Bytes())
}
