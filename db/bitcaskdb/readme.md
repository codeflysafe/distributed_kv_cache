# mini DB

https://github.com/flower-corp/minidb

## Bitcask 介绍
[Bitcask](doc/caskdb.md)

## 设计方案
## Entry
![](https://raw.githubusercontent.com/codeflysafe/gitalk/main/img/20220409154644.png)

```go
type Entry struct {
	Key       []byte
	KeySize   uint32 // key 的长度
	Value     []byte
	ValueSize uint32 // value 长度
	Mark      uint16 // 标记，是否删除
}
```

固定大小为 4 + 4 + 2

![](https://raw.githubusercontent.com/codeflysafe/gitalk/main/img/20220409154916.png)

