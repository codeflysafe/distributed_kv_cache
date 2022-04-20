package ds

import "nosdb/ds/set"

type Set interface {
	SAdd(member string, value []byte)
	SCard() int
	SPop() []byte
	SRem(member string)
	Empty() bool
}

func NewSet() Set {
	return set.NewMSet()
}
