package testDemo

import (
	"fmt"
	"sync"
	"time"
)

/*
题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
考察点 ： go 关键字的使用、协程的并发执行。
*/

// 打印奇数
// 它接受一个 *sync.WaitGroup 作为参数，以便在完成时通知主程序
func printOdd(wg *sync.WaitGroup) {
	// defer 语句确保在函数退出前一定会调用 wg.Done()
	defer wg.Done()

	for i := 1; i <= 10; i += 2 {
		fmt.Println("打印奇数：", i)
		time.Sleep(time.Second * 1)
	}
}

// 打印偶数
func printEven(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 2; i <= 10; i += 2 {
		fmt.Println("打印偶数：", i)
		time.Sleep(time.Second * 1)
	}
}

func Test07() {
	// 创建一个 WaitGroup，它会等待一组协程完成
	var wg sync.WaitGroup
	// 我们需要等待两个协程，所以调用 Add(2)
	wg.Add(2)

	// 使用 go 关键字启动两个协程
	go printOdd(&wg)
	go printEven(&wg)

	// Wait() 会阻塞当前程序，直到 WaitGroup 的计数器变为零
	//（即两个协程都调用了 wg.Done()）
	wg.Wait()

	fmt.Println("所有协程执行完毕")
}
