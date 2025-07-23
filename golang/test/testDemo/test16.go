package testDemo

import (
	"fmt"
	"sync"
)

// 启动10个协程打印1-100的数字

func Test16() {
	var wg sync.WaitGroup
	// 启动10个协程打印1-100的数字
	totalNum := make(chan int, 100) // 创建通道，容量为100

	// 启动10个协程打印
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(workId int) {
			defer wg.Done()
			for num := range totalNum {
				fmt.Printf("协程%d,打印数字：%d\n", workId, num)
			}
		}(i)
	}

	//////////  确保消费者先启动，再启动生产者，可以避免发生死锁   //////////
	// 通道放入数字
	for i := 1; i <= 100; i++ {
		totalNum <- i
	}

	close(totalNum)
	wg.Wait()
}

func Test1602() {
	wg := sync.WaitGroup{}
	chanNum := make(chan int, 1)

	for i := range 9 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for num := range chanNum {
				if num > 100 {
					chanNum <- num
					return
				}
				fmt.Printf("协程%d,打印数字：%d\n", i, num)
				chanNum <- num + 1
			}
		}()
	}

	chanNum <- 1

	wg.Wait()
	close(chanNum)
}
