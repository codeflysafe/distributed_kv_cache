package nosdb

import "nosdb/logfile"

/*
 * @Author: sjhuang
 * @Date: 2022-04-25 11:49:53
 * @LastEditTime: 2022-05-09 17:18:45
 * @FilePath: /nosdb/db_cmd.go
 */

const (
	// ====== string cmd =======
	SET logfile.CMD = iota
	DEL
	// ======== set cmd ==========
	SADD
	SREM

	// ======== list cmd ==========
	LPUSH
	LPOP
	RPUSH
	RPOP

	// ======= hash cmd ===========
	HSET
	HDEL
)
