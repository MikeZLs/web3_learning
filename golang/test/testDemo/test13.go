package testDemo

import (
	"fmt"
	"sync"
)

// 编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。

func Test13() {
	// 创建一个互斥锁
	var mutex sync.Mutex
	var count int
	var wg sync.WaitGroup
	wg.Add(10) // 10个协程
	for range 10 {
		go func() {
			defer wg.Done()

			for range 1000 {
				// 加锁
				mutex.Lock()
				count++
				// 解锁
				mutex.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Println("计数器的值:", count)

}
