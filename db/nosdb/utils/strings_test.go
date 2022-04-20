package utils

import (
	"fmt"
	"strconv"
	"testing"
)

func TestStrIncrBy(t *testing.T) {

	vals := []string{"1", "2", "-100", strconv.Itoa(INT_MIN), strconv.Itoa(INT_MAX)}
	offsets := []int{12, -2, 0, INT_MAX, INT_MIN}
	newVals := []string{"13", "0", "-100", "-1", "-1"}

	for i := 0; i < 5; i++ {
		val := vals[i]
		offset := offsets[i]
		newVal := newVals[i]
		newV, err := StrIncrBy(val, offset)
		if err != nil {
			t.Error(err)
		}
		if newV != newVal {
			t.Errorf(" res is error, get %s, real is %s", newV, newVal)
		}
	}

	val := strconv.Itoa(INT_MAX)
	offset := 1
	newV, err := StrIncrBy(val, offset)
	if err == nil {
		t.Errorf("   ------  ")
	} else {
		fmt.Println(err.Error(), newV)
	}

	val = strconv.Itoa(INT_MIN)
	offset = -1
	newV, err = StrIncrBy(val, offset)
	if err == nil {
		t.Errorf("   ------  ")
	} else {
		fmt.Println(err.Error(), newV)
	}

}

func TestBytesIncrBy(t *testing.T) {
	val := "1"
	offset := 12
	newV, err := StrIncrBy(val, offset)
	if err != nil {
		t.Error(err)
	}
	if newV != "13" {
		t.Errorf(" res is err, get %s", newV)
	}
}

// todo 处理精度问题
func TestStrIncrByFloat(t *testing.T) {
	vals := []string{"1", "2", "-100", strconv.Itoa(INT_MIN), strconv.Itoa(INT_MAX)}
	offsets := []float64{12.1, -2.1, 1.0, float64(1.0 * INT_MAX), float64(1.0 * INT_MIN)}
	newVals := []string{"13.1", "-0.1", "-99.0", "-1.0", "-1.0"}

	for i := 0; i < 5; i++ {
		val := vals[i]
		offset := offsets[i]
		newVal := newVals[i]
		newV, err := StrIncrByFloat(val, offset)
		if err != nil {
			t.Error(err)
		}
		if newV != newVal {
			t.Errorf(" res is error, get %s, real is %s", newV, newVal)
		}
	}
}
