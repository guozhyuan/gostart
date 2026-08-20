
只要 panic 发生在 Gin 接收请求的同一个主协程（goroutine） 里，gin.Recovery() 中间件100%能保证程序（进程）不会退出，只会返回 500 错误
但是不能处理主协程fatal
不能处理子协程中的panic和fatal，每个子协程都需要recover（）

(可以使用zap替代默认writer)


业务代码中常见的 fatal，90% 都集中在以下 4 种场景：

1. 并发读写 Map（生产环境头号杀手）
这是线上最频发的 fatal，没有之一。Go 的 Map 在并发读写时，运行时检测到数据竞争会直接 throw。

```go
m := make(map[int]int)
go func() { m[1] = 1 }() // 写
go func() { _ = m[1] }() // 同时读（或另一个写），触发 fatal error: concurrent map read and map write
```

2. 栈内存溢出（Stack Overflow）
函数无限递归，导致栈空间耗尽。Go 运行时无法动态扩容无限次，只能崩溃。

```go
func recurse() { recurse() } // 无限递归，触发 fatal error: stack overflow
```

3. 内存耗尽（OOM）
尝试分配超出进程可用物理内存或系统限制的大块内存（注意：不是 Go 堆内存不足触发的 GC，那是 panic，这里是系统级的 runtime: out of memory）。

```go
// 尝试分配 10GB 内存（通常分配失败触发 fatal）
_ = make([]byte, 10<<30) 
```

4. 所有 Goroutine 进入休眠（死锁）
当所有协程都因等待锁、通道等而阻塞，没有协程可调度时，调度器会检测到死锁并 throw。

```go
func main() {
    ch := make(chan int)
    <-ch // 阻塞等待，且没有其他协程往 ch 发数据
}
// 触发 fatal error: all goroutines are asleep - deadlock!
```

---
