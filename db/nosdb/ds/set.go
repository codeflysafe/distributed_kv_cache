/*
 * @Author: sjhuang
 * @Date: 2022-04-19 09:30:23
 * @LastEditTime: 2022-04-27 10:44:18
 * @FilePath: /nosdb/ds/set.go
 */
package ds

import "nosdb/ds/set"

type Set interface {
	SAdd(member string, value []byte)
	SCard() int
	SPop() []byte
	SRem(member string)
	Empty() bool
	SIsMember(member string) bool
}

func NewSet() Set {
	return set.NewMSet()
}
