package testDemo

import (
	"fmt"
	"sync"
	"time"
)

// 题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
// 考察点 ：协程原理、并发任务调度。

// 定义任务类型
type Task func()

func Test08(tasks []Task) {
	var wg sync.WaitGroup

	for i, task := range tasks {
		wg.Add(1)

		go func(id int, t Task) {
			defer wg.Done()

			start := time.Now()
			t()
			duration := time.Since(start)

			fmt.Printf("任务 #%d 执行耗时: %v\n", id+1, duration)
		}(i, task) // 显式传值，避免闭包坑
	}

	wg.Wait()
	fmt.Println("所有任务已完成!")
}
