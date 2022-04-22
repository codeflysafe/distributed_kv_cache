package ds

import "nosdb/ds/zskset"

type ZSet interface {
	ZAdd(score float64, member string, value []byte)
	ZDel(score float64, member string)
	ZRange(minScore, maxScore float64, minEx, maxEx bool)
	ZCount(minScore, maxScore float64, minEx, maxEx bool) int
	ZCard() int
	ZIncrScore(member string, value []byte, offset float64)
	ZScore(member string) float64
}

func NewZSet() ZSet {
	return zskset.NewZSet()
}
