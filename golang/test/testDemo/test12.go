package testDemo

import "fmt"

// 实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。

// 生产者协程向通道中发送100个整数
func producer(ch chan<- int) {
	for i := 1; i <= 100; i++ {
		ch <- i
	}
	close(ch) // 关闭通道，表示写入完成
}

// 消费者协程从通道中接收这些整数并打印
func consumer(ch <-chan int) {
	for num := range ch {
		fmt.Println(num)
	}
}

func Test12() {
	h := make(chan int, 100) // 创建一个容量为100的缓冲通道
	go producer(h)
	consumer(h)
}
