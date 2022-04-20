package ds

import "nosdb/ds/hash"

// hash 接口
type Hash interface {
	HDel(member string)
	HExists(member string) bool
	HGet(member string) []byte
	HLen(member string) int
	HSet(member string, value []byte)
	HSetNx(member string, value []byte)
	HIncrBy(member string, offset int) error
	HIncrByFloat(member string, offset float64) error
}

func NewHash() Hash {
	return hash.NewMHash()
}
