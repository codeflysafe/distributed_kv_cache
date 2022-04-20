package ds

import "nosdb/ds/bstring"

type String interface {
	Set(value []byte)
	Get() []byte
	GetSet(newVal []byte) []byte
	GetRange(start, end int) []byte
	IncrByInt(offset int) (err error)
	IncrByFloat(offset float64) (err error)
	Append(as []byte)
	ToString() string
	StrLen() int
}

func NewString() String {
	return bstring.NewBString()
}
