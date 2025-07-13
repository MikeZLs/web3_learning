package testDemo

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。

func Test14() {
	var count int64
	var wg sync.WaitGroup
	wg.Add(10)
	for range 10 {
		go func() {
			defer wg.Done()

			for range 1000 {
				// Go 的标准库 sync/atomic 提供了高性能的 原子操作，可用于实现无锁（lock-free）的并发安全计数器。
				// atomic.AddInt64(&counter, 1)：对 counter 做原子递增，无需加锁
				// &counter：传入计数器变量的地址（必须是指针）
				// 使用 int64 而非 int：因为 atomic 明确要求使用具体类型如 int32 / int64
				atomic.AddInt64(&count, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Println("计数器的值:", count)
}
