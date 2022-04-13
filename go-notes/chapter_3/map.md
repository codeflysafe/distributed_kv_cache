# map
![](https://raw.githubusercontent.com/codeflysafe/gitalk/main/img/20220412170350.png)

```go
// A header for a Go map.
type hmap struct {
    // 元素个数，调用 len(map) 时，直接返回此值
    count     int
    flags     uint8
    // buckets 的对数 log_2
    B         uint8
    // overflow 的 bucket 近似数
    noverflow uint16
    // 计算 key 的哈希的时候会传入哈希函数
    hash0     uint32
    // 指向 buckets 数组，大小为 2^B
    // 如果元素个数为0，就为 nil
    buckets    unsafe.Pointer
    // 扩容的时候，buckets 长度会是 oldbuckets 的两倍
    oldbuckets unsafe.Pointer
    // 指示扩容进度，小于此地址的 buckets 迁移完成
    nevacuate  uintptr
    extra *mapextra // optional fields
}
```
![](https://raw.githubusercontent.com/codeflysafe/gitalk/main/img/20220412172129.png)
B 是 buckets 数组的长度的对数，也就是说 buckets 数组的长度就是 2^B。bucket 里面存储了 key 和 value.

mapextra 结构体
```go
// mapextra holds fields that are not present on all maps.
type mapextra struct {
    // 如果 key 和 value 都不包含指针，并且可以被 inline(<=128 字节)
    // 就使用 hmap的extra字段 来存储 overflow buckets，这样可以避免 GC 扫描整个 map
    // 然而 bmap.overflow 也是个指针。这时候我们只能把这些 overflow 的指针
    // 都放在 hmap.extra.overflow 和 hmap.extra.oldoverflow 中了
    // overflow 包含的是 hmap.buckets 的 overflow 的 buckets
    // oldoverflow 包含扩容时的 hmap.oldbuckets 的 overflow 的 bucket
    overflow    *[]*bmap
    oldoverflow *[]*bmap

    // 指向空闲的 overflow bucket 的指针
    nextOverflow *bmap
}
```

bmap结构体
```go
// A bucket for a Go map.
type bmap struct {
    // tophash包含此桶中每个键的哈希值最高字节（高8位）信息（也就是前面所述的high-order bits）。
    // 如果tophash[0] < minTopHash，tophash[0]则代表桶的搬迁（evacuation）状态。
    tophash [bucketCnt]uint8
}
```

bmap结构体在编译时，变成了
```go
bmap := MapBucketType(t)
// MapBucketType

type bmap struct {
topbits  [8]uint8
keys     [8]keytype
elems   [8]valuetype
pad      uintptr
overflow uintptr
}
```

![](https://raw.githubusercontent.com/codeflysafe/gitalk/main/img/20220412172103.png)

### bucket 
每个 bucket 设计成最多只能放 8 个 key-value 对，如果有第 9 个 key-value 落入当前的 bucket，那就需要再构建一个 bucket ，
通过 overflow 指针连接起来
![](https://raw.githubusercontent.com/codeflysafe/gitalk/main/img/20220412172809.png)

### key 定位过程
key 经过哈希计算后得到哈希值，共 64 个 bit 位（64位机，32位机就不讨论了，现在主流都是64位机），计算它到底要落在哪个桶时，只会用到最后 B 个 bit 位。还记得前面提到过的 B 吗？如果 B = 5，那么桶的数量，也就是 buckets 数组的长度是 2^5 = 32。
例如，现在有一个 key 经过哈希函数计算后，得到的哈希结果是：
```text
 10010111 | 000011110110110010001111001010100010010110010101010 │ 01010
```
用最后的 5 (B) 个 bit 位，也就是 01010，值为 10，也就是 10 号桶。这个操作实际上就是取余操作，但是取余开销太大，所以代码实现上用的位操作代替。
再用哈希值的高 8 位，找到此 key 在 bucket 中的位置，这是在寻找已有的 key。最开始桶内还没有 key，新加入的 key 会找到第一个空位，放入。
buckets 编号就是桶编号，当两个不同的 key 落在同一个桶中，也就是发生了哈希冲突。冲突的解决手段是用链表法：在 bucket 中，从前往后找到第一个空位。这样，在查找某个 key 时，先找到对应的桶，再去遍历 bucket 中的 key。
这里参考曹大 github 博客里的一张图，原图是 ascii 图，geek 味十足，可以从参考资料找到曹大的博客，推荐大家去看看。
![](https://raw.githubusercontent.com/codeflysafe/gitalk/main/img/20220412172551.png)


### 非并发安全
```go
	// flags
iterator     = 1 // there may be an iterator using buckets
oldIterator  = 2 // there may be an iterator using oldbuckets
hashWriting  = 4 // a goroutine is writing to the map
sameSizeGrow = 8 // the current map growth is to a new map of the same size

if h.flags&hashWriting != 0 {
	throw("concurrent map read and map write")
}
```

如果出现一个goroutine 在修改，而其它在读取的情况，会panic并发错误

### 遍历

```go
// A hash iteration structure.
// If you modify hiter, also change cmd/compile/internal/reflectdata/reflect.go
// and reflect/value.go to match the layout of this structure.
type hiter struct {
	key         unsafe.Pointer // Must be in first position.  Write nil to indicate iteration end (see cmd/compile/internal/walk/range.go).
	elem        unsafe.Pointer // Must be in second position (see cmd/compile/internal/walk/range.go).
	t           *maptype
	h           *hmap
	buckets     unsafe.Pointer // bucket ptr at hash_iter initialization time
	bptr        *bmap          // current bucket
	overflow    *[]*bmap       // keeps overflow buckets of hmap.buckets alive
	oldoverflow *[]*bmap       // keeps overflow buckets of hmap.oldbuckets alive
	startBucket uintptr        // bucket iteration started at
	offset      uint8          // intra-bucket offset to start from during iteration (should be big enough to hold bucketCnt-1)
	wrapped     bool           // already wrapped around from end of bucket array to beginning
	B           uint8
	i           uint8
	bucket      uintptr
	checkBucket uintptr
}

```
由于遍历的时候，随机从一个bucket上开始，所有每次遍历结果都不相同
```go
// decide where to start
	r := uintptr(fastrand())
	if h.B > 31-bucketCntBits {
		r += uintptr(fastrand()) << 31
	}
	it.startBucket = r & bucketMask(h.B)
```

### 扩容
1. loadfactor 超过阈值 6.5
2. too many overflow buckets, overflow buckets 过多 B > 15

## Reference
https://qcrao91.gitbook.io/go/map/map-de-di-ceng-shi-xian-yuan-li-shi-shi-mo