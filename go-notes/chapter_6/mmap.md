# mmap 函数

https://man7.org/linux/man-pages/man2/mmap.2.html

存储映射io, 讲一个磁盘文件映射到存储空间的一个缓冲区上，于是当从缓冲区取数据时，就相当于读文件中的相应字节。于此类似，将数据存入缓冲区时，相应字节就自动写入文件。这样可以不实用read 和
write 的情况下执行io

## `linux` 下 `mmap` 函数

```c
// linux 中
 #include <sys/mman.h>
  void *
  mmap(void *addr, size_t len, int prot, int flags, int fd, off_t offset);
```
`addr` 参数指定映射存储区的起始位置，通常设置为0，这表示由系统自动选择该映射区的起始位置\
`fd` 参数，代表一个文件描述符。在文件映射到地址空间之前，必须打开该文件。\
`len` 参数是映射的字节\
`offset`是要映射的字节在文件中的偏移量 \
`prot` 是对存储区的保护

| prot       | 说明    |
|------------|-------|
| PROT_READ  | 可读    |
| PROT_WRITE | 可写    |
| PROT_EXEC  | 可执行   |
| PROT_NONE  | 不可访问  |

`flag`字段， 主要是三种

| flag        | 说明                         |
|-------------|----------------------------|
| MAP_FIXED   | 返回值必须等于 addr，移植性差，一般不用     |
| MAP_SHARED  | 指定存储操作修改映射文件               |
| MAP_PRIVATE | 为映射文件创造一个私有副本，存储操作修改的是私有副本 |

![](https://raw.githubusercontent.com/codeflysafe/gitalk/main/img/20220421153043.png)