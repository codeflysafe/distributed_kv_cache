package ds

type Set interface {
	SAdd(member string, value []byte)
	SCard() int
	SPop() []byte
	SRem(member string)
	Empty() bool
}
