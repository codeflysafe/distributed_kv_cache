# sync 包

## sync.Pool

核心是，对象复用。对于很多需要重复分配、回收内存的地方，sync.Pool 是一个很好的选择。频繁地分配、回收内存会给 GC 带来一定的负担，严重的时候
会引起 `CPU` 的毛刺，而 `sync.Pool` 可以将暂时不用的对象缓存起来，待下次需要的时候直接使用，不用再次经过内存分配，复用对象的内存，减轻 GC` 
的压力，提升系统的性能

## sync.Once
`sync.Once` 是 `Go` 标准库提供的使函数只执行一次的实现，常应用于单例模式，例如初始化配置、保持数据库连接等。作用与 `init` 函数类似，但有区别。
`init` 函数是当所在的 `package` 首次被加载时执行，若迟迟未被使用，则既浪费了内存，又延长了程序加载时间。
`sync.Once` 可以在代码的任意位置初始化和调用，因此可以延迟到使用时再执行，并发场景下是线程安全的。

```go
package main

import (
	"fmt"
	"sync"
)

var (
	once sync.Once
)

type A struct {
	index map[string][]byte
	child *A
}

func NewA() *A {
	return new(A)
}

func (a *A) lazyLoadIndex() {
	once.Do(func() {
		fmt.Println("lazyLoad index")
		if a.index == nil {
			a.index = make(map[string][]byte, 10)
		}
	})
}

func (a *A) lazyLoadChild() {
	once.Do(func() {
		fmt.Println("lazyLoad child")
		if a.child == nil {
			a.child = NewA()
		}
	})
}

func main() {
	a := new(A)
	wait := sync.WaitGroup{}
	count := 5
	wait.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			a.lazyLoadIndex()
			wait.Done()
		}()
	}
	wait.Wait()

	wait.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			a.lazyLoadChild()
			wait.Done()
		}()
	}
	wait.Wait()

}
```

输出
```shell
lazyLoad index
```

它无论是否调用的是一个函数，它只执行一次
### 实现原理

```go
type Once struct {
	// done indicates whether the action has been performed.
	// It is first in the struct because it is used in the hot path.
	// The hot path is inlined at every call site.
	// Placing done first allows more compact instructions on some architectures (amd64/386),
	// and fewer instructions (to calculate offset) on other architectures.
	done uint32
	m    Mutex
}
```

它的实现类似于Java双重检验锁实现单例模式
```java
public class Singleton {  
    private volatile static Singleton singleton;  
    private Singleton (){}  
    public static Singleton getSingleton() {  
    if (singleton == null) {  
        synchronized (Singleton.class) {  
            if (singleton == null) {  
                singleton = new Singleton();  
            }  
        }  
    }  
    return singleton;  
    }  
}
```
不同的是，`go` 里面没有`volatile` 关键字，这里采用原子类操作`atomic.LoadUint32`
其它操作一样，类比即可

```go
func (o *Once) Do(f func()) {
	// Note: Here is an incorrect implementation of Do:
	//
	//	if atomic.CompareAndSwapUint32(&o.done, 0, 1) {
	//		f()
	//	}
	//
	// Do guarantees that when it returns, f has finished.
	// This implementation would not implement that guarantee:
	// given two simultaneous calls, the winner of the cas would
	// call f, and the second would return immediately, without
	// waiting for the first's call to f to complete.
	// This is why the slow path falls back to a mutex, and why
	// the atomic.StoreUint32 must be delayed until after f returns.
	
	if atomic.LoadUint32(&o.done) == 0 {
		// Outlined slow-path to allow inlining of the fast-path.
		o.doSlow(f)
	}
}
func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	// 这里需要判断的原因是， if atomic.LoadUint32(&o.done) == 0 是，
	// 可能会导致多个协程进入临界区
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
```