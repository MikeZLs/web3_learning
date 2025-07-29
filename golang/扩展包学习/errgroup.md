### errgroup是`golang.org/x/sync/errgroup`包提供的一个并发控制工具，它是对标准库`sync.WaitGroup`的增强版本。 ###

## 核心原理

### 1. 结构设计
```go
type Group struct {
    cancel func()        // context取消函数
    wg     sync.WaitGroup // 等待所有goroutine完成
    errOnce sync.Once    // 确保只记录第一个错误
    err     error        // 存储第一个出现的错误
}
```

### 2. 四大核心机制

**错误传播机制**：使用`sync.Once`确保只有第一个出现的错误被记录，避免竞态条件。
```go
func (g *Group) Go(f func() error) {
    g.wg.Add(1)
    go func() {
        defer g.wg.Done()
        if err := f(); err != nil {
            g.errOnce.Do(func() {
                g.err = err          // 记录错误
                if g.cancel != nil {
                    g.cancel()       // 触发context取消
                }
            })
        }
    }()
}
```

**Context取消传播**：当任意goroutine出错时，通过调用`cancel()`函数取消context，实现级联取消。
```go
func WithContext(ctx context.Context) (*Group, context.Context) {
    ctx, cancel := context.WithCancel(ctx)
    return &Group{cancel: cancel}, ctx
}
```

**同步等待机制**：底层使用`sync.WaitGroup`协调所有goroutine的完成。
```go
func (g *Group) Wait() error {
    g.wg.Wait()    // 等待所有goroutine完成
    if g.cancel != nil {
        g.cancel()  // 确保context被取消
    }
    return g.err    // 返回第一个错误（如果有的话）
}
```

**工作流程**：启动阶段增加计数 → 并发执行监听取消信号 → 错误时记录并取消 → 等待清理后返回错误

## 实际使用

### 基本用法
```go
package main

import (
    "context"
    "fmt"
    "time"
    "golang.org/x/sync/errgroup"
)

func main() {
    g, ctx := errgroup.WithContext(context.Background())
    
    // 启动多个并发任务
    for i := 0; i < 3; i++ {
        i := i // 避免闭包问题 (go 1.22及之后可不用)
        g.Go(func() error {
            select {
            case <-time.After(time.Duration(i) * time.Second):
                fmt.Printf("任务 %d 完成\n", i)
                return nil
            case <-ctx.Done():
                fmt.Printf("任务 %d 被取消\n", i)
                return ctx.Err()
            }
        })
    }
    
    // 等待所有任务完成或第一个错误
    if err := g.Wait(); err != nil {
        fmt.Printf("出现错误: %v\n", err)
    }
}
```

### 实际应用场景

**1. 并发网络请求**
```go
func fetchMultipleURLs(urls []string) error {
    g, ctx := errgroup.WithContext(context.Background())
    
    for _, url := range urls {
        url := url
        g.Go(func() error {
            return fetchURL(url, ctx)
        })
    }
    
    return g.Wait()
}
```

**2. 限制并发数量**
```go
func processWithLimit(tasks []Task) error {
    g, ctx := errgroup.WithContext(context.Background())
    
    // 限制最多10个并发任务
    semaphore := make(chan struct{}, 10)
    
    for _, task := range tasks {
        task := task
        g.Go(func() error {
            semaphore <- struct{}{}        // 获取信号量
            defer func() { <-semaphore }() // 释放信号量
            
            return processTask(task, ctx)
        })
    }
    
    return g.Wait()
}
```

**3. 数据管道处理**
```go
func pipelineProcess(data []Data) error {
    g, ctx := errgroup.WithContext(context.Background())
    
    // 阶段1：数据预处理
    g.Go(func() error {
        return preprocessData(data, ctx)
    })
    
    // 阶段2：数据验证
    g.Go(func() error {
        return validateData(data, ctx)
    })
    
    // 阶段3：数据存储
    g.Go(func() error {
        return saveData(data, ctx)
    })
    
    return g.Wait()
}
```

## 设计理念与优势

**快速失败**：一旦有任务失败，立即通知其他任务停止，避免资源浪费

**错误优先级**：只关心第一个错误，简化错误处理逻辑

**优雅关闭**：通过context机制实现协作式取消，而非强制终止

**简化并发代码**：相比手动管理WaitGroup和错误channel，提供更简洁的API

errgroup特别适合需要并发执行多个任务，且希望在任一任务失败时能够优雅地取消其他任务的场景，如并发的网络请求、文件处理、数据库操作等。这种设计使得errgroup既保持了简单的API，又提供了强大的并发控制和错误处理能力。