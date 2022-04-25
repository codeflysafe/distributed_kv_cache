package nosdb

import "nosdb/logfile"

/*
 * @Author: sjhuang
 * @Date: 2022-04-25 11:49:53
 * @LastEditTime: 2022-04-25 16:51:35
 * @FilePath: /nosdb/db_cmd.go
 */

const (
	// ====== string cmd =======
	SET logfile.CMD = iota
	GetSet
	SetNx
	IncrByInt
	IncrByFloat
	Append

	// ======== string cmd ==========
)
