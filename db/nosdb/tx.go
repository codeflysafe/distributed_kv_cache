package nosdb

import "nosdb/snowflake"

var node, _ = snowflake.NewNode(2)

// tx 事务，相关
// 如何实现 ACID
type TX struct {
	txId int64 // 事务id
}

// 新建一个事务
func NewTx() *TX {
	tx := &TX{
		txId: node.Generate().Int64(),
	}
	return tx
}
