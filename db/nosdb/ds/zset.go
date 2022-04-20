package ds

import "nosdb/ds/zskset"

/* Struct to hold a inclusive/exclusive range spec by score comparison. */
type ZRangeSpec struct {
	MinScore, MaxScore float64 // min Score -> maxScore
	Minex, Maxex       bool    /* are min or max exclusive? */
}

type ZSet interface {
	ZAdd(score float64, member string, value []byte)
	ZDel(score float64, member string)
	ZRange(rangspec ZRangeSpec)
	ZCount(rangspec ZRangeSpec) int
	ZCard() int
}

func NewZSet() ZSet {
	return zskset.NewZSet()
}
