### Golang标准库下的 context 核心机制

### Context 接口定义（标准接口）

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool) // 返回取消时间（如有），无则 ok=false
    Done() <-chan struct{}                   // 返回一个 channel，在 context 被取消时关闭
    Err() error                              // 返回取消原因（Canceled 或 DeadlineExceeded）
    Value(key any) any                       // 获取绑定的值，用于跨 API 传递 request 范围内数据
}
```

---

### Context 的构造函数及用途

| 构造函数 | 说明 |
|----------|------|
| `context.Background()` | 根 context，从不取消，适合作为程序入口的基础 context |
| `context.TODO()` | 占位 context，用于暂时不知道用哪个 context 的情况 |
| `context.WithCancel(parent)` | 返回派生子 context 和 cancel 函数，手动触发取消 |
| `context.WithTimeout(parent, dur)` | 设置 context 超时，时间到自动取消 |
| `context.WithDeadline(parent, time)` | 设置绝对时间点作为取消时间 |
| `context.WithValue(parent, key, val)` | 在 context 中携带一对键值，用于跨 API/协程传参 |

---

### Context 的内部结构与派生关系

```go
// 永不取消、无值、无 deadline，用于 Background 和 TODO
type emptyCtx struct{}

// 可取消的 context，支持 Done、Err
type cancelCtx struct {
    Context                     // 指向父 context
    done     atomic.Value       // 懒加载 <-chan struct{}
    children map[canceler]struct{} // 所有子 context
    err      error              // 被取消的错误（Canceled、DeadlineExceeded）
    cause    error              // 取消的详细原因
}

// timerCtx 表示带有截止时间的 cancelCtx
type timerCtx struct {
    cancelCtx
    deadline time.Time
    timer    *time.Timer
}

// 携带键值的 context，用于 WithValue
type valueCtx struct {
    Context
    key, val any
}
```

---

### 取消信号的传播流程（链式）

1. 创建派生 context（例如 WithCancel）
2. `cancelCtx.propagateCancel(parent, child)` 注册子 context 到父 context 的 `children`
3. 父 context 被取消后，会递归调用子 context 的 `cancel()`
4. `cancel()` 会：
    - 设置 `err` 和 `cause`
    - 关闭 `done` channel
    - 调用所有子 context 的 `cancel()`

---

### 常见错误类型

```go
var Canceled = errors.New("context canceled")               // 手动或父 context 取消
var DeadlineExceeded = errors.New("context deadline exceeded") // 超时取消
```

---

### 辅助函数与增强能力

```go
func Cause(ctx Context) error
```
- 获取取消的根本原因（手动设置的 `cause`，或 fallback 到 `Err()`）

```go
func AfterFunc(ctx Context, f func()) func() bool
```
- 注册一个函数，在 context 被取消后调用（在新 goroutine 中）

```go
func WithoutCancel(ctx Context) Context
```
- 派生一个子 context，不会因父 context 取消而被影响（可用于日志上下文传递等）

---

### 典型用法示例

#### 超时控制

```go
func Fetch(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel() // 必须调用，释放 timer 资源

    select {
    case <-ctx.Done():
        return ctx.Err() // 超时或手动取消
    case data := <-getData():
        return process(data)
    }
}
```

#### 跨协程传递取消信号

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("worker canceled:", ctx.Err())
            return
        case item := <-taskChan:
            do(item)
        }
    }
}
```

---

### 使用建议（官方推荐）

#### 应该

- 将 `ctx` 作为第一个参数传递，命名为 `ctx`
- 使用 `defer cancel()` 确保及时释放资源
- 使用 `ctx.Value()` 仅传递请求范围内的数据（如用户信息、追踪 ID）

#### 不应该

- 不应将 `Context` 存入结构体字段中
- 不应使用 `Value()` 储存配置或全局数据
- 不应忽略 `cancel()` 的调用（`go vet` 可自动检测）

---

### 补充说明：Value 查找机制

```go
func (c *valueCtx) Value(key any) any {
    if c.key == key {
        return c.val
    }
    return c.Context.Value(key)
}
```

- Value 是 **链式向上查找** 的
- 为避免 key 冲突，推荐使用非导出类型作为 key，例如：

```go
type ctxKey struct{}
var userKey = ctxKey{}

ctx := context.WithValue(ctx, userKey, user)
```

---

### 总结

- `context` 提供对取消、超时、键值传递的统一抽象
- 通过链式组合和传播机制，适配服务端 request 生命周期管理
- 合理使用 `Done()` 和 `Err()`，是 goroutine 协作的基础

