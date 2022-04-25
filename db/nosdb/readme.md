<!--
 * @Author: sjhuang
 * @Date: 2022-04-12 21:29:10
 * @LastEditTime: 2022-04-25 10:29:46
 * @FilePath: /nosdb/readme.md
-->
# notsdb

![](https://raw.githubusercontent.com/codeflysafe/gitalk/main/img/20220412215321.png)

`not simple db`, A high performance NoSQL database based on `bitcask`, supports string, list, hash, set, and zset.



## 文件树

```shell
.
├── cache
│   ├── cache.go
│   └── lru.go
├── config.go
├── db_client.go
├── db_hash.go
├── db_list.go
├── db_set.go
├── db_string.go
├── db_zset.go
├── ds # 实现的数据结构
│   ├── bstring
│   │   ├── bstring.go
│   │   └── bstring_test.go
│   ├── ds.go
│   ├── hash
│   │   ├── mhash.go
│   │   └── mhash_test.go
│   ├── hash.go
│   ├── linkedlist
│   │   ├── cpu.pprof
│   │   ├── linkedlist.go
│   │   ├── linkedlist.test
│   │   ├── linkedlist_benchmark_test.go
│   │   └── linkedlist_test.go
│   ├── list.go
│   ├── set
│   │   └── mset.go
│   ├── set.go
│   ├── slicelist
│   │   ├── slicelist.go
│   │   ├── slicelist_benchmark_test.go
│   │   └── slicelist_test.go
│   ├── string.go
│   ├── zset.go
│   └── zskset
│       ├── skiplist.go
│       ├── skiplist_test.go
│       └── zskset.go
├── file
│   ├── a.txt
│   │   └── bench_mack
│   ├── file_handle.go
│   ├── file_handle_benchmark_test.go
│   ├── file_handle_test.go
│   ├── io_file_handle.go
│   ├── mmap
│   │   ├── mmap.go
│   │   ├── mmap_test.go
│   │   ├── mmap_unix.go
│   │   └── test.txt
│   └── mmap_file_handle.go
├── go.mod
├── go.sum
├── idx.go
├── nosdb.go
├── nosdb_test.go
├── readme.md
├── snowflake
│   ├── snowflake.go
│   └── snowflake_test.go
├── utils
│   ├── file_utils.go
│   ├── file_utils_test.go
│   ├── strings.go
│   └── strings_test.go
└── wal
    ├── entry.go
    ├── entry_test.go
    ├── wal.go
    └── wal_test.go

14 directories, 57 files
```

## Plan
- [ ] bitcask\ B+Tree\ LSM_Tree 了解和学习
- [x] 参考 redis 实现基本数据结构 ziplist\skiplist\linkedlist\string


## 参考
[rosedb](https://github.com/flower-corp/rosedb)
[nutsdb](https://github.com/nutsdb/nutsdb)
