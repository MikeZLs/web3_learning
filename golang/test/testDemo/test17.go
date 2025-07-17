package testDemo

import (
	"fmt"
	"sync"
)

//启动10个协程打印1-100的数字（顺序打印）

func Test17() {
	// 创建一个通道数组，用于协程间的传递
	channels := make([]chan int, 10)
	for i := range channels {
		channels[i] = make(chan int)
	}

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			// 每个协程从自己的通道接收数字，打印，然后传递给下一个协程
			for {
				num, ok := <-channels[workerID]
				if !ok { // 如果通道被关闭，则退出
					return
				}

				fmt.Printf("协程 %d: 打印 %d\n", workerID, num)

				// 准备传递给下一个协程
				nextNum := num + 1
				if nextNum > 100 {
					//如果是最后一个数字，关闭所有通道以终止所有协程
					for j := range channels {
						close(channels[j])
					}
					return
				}

				// 将“接力棒”（下一个数字）交给下一个协程
				nextWorkerID := (workerID + 1) % 10
				channels[nextWorkerID] <- nextNum
			}
		}(i)
	}

	// 将第一个数字“1”交给0号协程，启动接力
	channels[0] <- 1

	wg.Wait() // 等待所有协程完成
}
