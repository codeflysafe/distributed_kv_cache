package nosdb

import "unsafe"

type DBClient struct {
	// 当前的 bucket 索引
	cur int
	// 所有的 db
	buckets unsafe.Pointer
}
