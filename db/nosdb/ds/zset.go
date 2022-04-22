package ds

import "nosdb/ds/zskset"

/* Struct to hold a inclusive/exclusive range spec by score comparison. */
type ZRangeSpec struct {
	MinScore, MaxScore float64 // min Score -> maxScore
	MinEx, MaxEx       bool    /* are min or max exclusive? */
}

type ZSet interface {
	ZAdd(score float64, member string, value []byte)
	ZDel(score float64, member string)
	ZRange(rangSpec ZRangeSpec)
	ZCount(rangSpec ZRangeSpec) int
	ZCard() int
	ZIncrScore(member string, value []byte, offset float64)
	ZScore(member string) float64
}

func NewZSet() ZSet {
	return zskset.NewZSet()
}
