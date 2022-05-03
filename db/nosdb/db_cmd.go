package nosdb

import "nosdb/logfile"

/*
 * @Author: sjhuang
 * @Date: 2022-04-25 11:49:53
 * @LastEditTime: 2022-04-28 11:12:00
 * @FilePath: /nosdb/db_cmd.go
 */

const (
	// ====== string cmd =======
	SET logfile.CMD = iota
	DEL
	// ======== set cmd ==========
	SADD
	SPOP
	SREM

	// ======== list cmd ==========
	LPUSH
	LPUSHX
	LPOP
	RPUSH
	RPUSHX
	RPOP
)
