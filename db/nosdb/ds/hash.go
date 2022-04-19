package ds

// hash 接口
type Hash interface {
	HDel(key string)
	HExists(key string) bool
	HGet(key string) []byte
	HLen(key string) int
	HSet(key string, value []byte)
	HSetNx(key string, value []byte)
	HIncrBy(key string, offset int) error
	HIncrByFloat(key string, offset float64) error
}
