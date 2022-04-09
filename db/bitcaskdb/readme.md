# Bitcask DB


## Bitcask 介绍
[Bitcask](doc/Bitcask.md)

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

### db_file

将一个文件作为持久化存储介质

```go
type DBFile struct {
	File   *os.File // 存储的文件∏
	OffSet int64    // 当前的偏移量
}
```

### DB

```go
// 数据库的定义
type BitcaskDB struct {
	io.Closer
	indexes  map[string]int64 // key 与 offset 的索引
	indexes2 map[string]int64 // merge 过程中使用的hash 表
	dbFile   *DBFile          // 数据文件
	dirPath  string           // 数据位置
	mu       sync.RWMutex     // 采用读写锁
	ratio    int64            // indexes2， 新建时的大小
	count    int64
}
```

#### Get
首先从内存的`indexes`中根据`key`获取`offset`，即指向`value`的地址,然后去文件中读取

#### Write

采用文件的追加写即可

```go
// 顺序写，写入Entry
func (f *DBFile) Write(e *Entry) (offset int64, err error) {
	buf, err := e.Encode()
	if err != nil {
		return
	}
	_, err = f.File.WriteAt(buf, f.OffSet)
	if err != nil {
		return
	}
	offset = f.OffSet
	f.OffSet += e.GetSize()
	return
}
```
#### DEL
`del` 操作是一种特殊的读写，实现采用的是写入一条特殊的删除记录
```go
	_, err = db.dbFile.Write(NewEntry(key, nil, DEL))
```

#### Merge
采用生产者和消费者模式，使用channel作为通信媒介
从旧的文件中生成 `entry`， 写到新的文件中去
![](https://raw.githubusercontent.com/codeflysafe/gitalk/main/img/20220409162109.png)