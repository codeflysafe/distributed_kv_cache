package ds

/* Struct to hold a inclusive/exclusive range spec by score comparison. */
type ZRangeSpec struct {
	MinScore, MaxScore float64 // min Score -> maxScore
	Minex, Maxex       bool    /* are min or max exclusive? */
	Limit              int     // limit nums
}

type ZSet interface {
	ZAdd(score float64, member string, value []byte)
	ZDel(score float64, member string)
	ZRange(rangspec ZRangeSpec)
}
