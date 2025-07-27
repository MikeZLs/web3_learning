```go
// Mutex 是一个互斥锁的实现
type Mutex struct {
    state int32  // 记录锁的状态（是否加锁、是否有等待者、是否饥饿等）
    sema  uint32 // 信号量，用于阻塞/唤醒等待 goroutine
}
```

### 状态位定义及含义

```go
const (
    mutexLocked      = 1 << iota // 二进制第0位：表示锁是否已被持有
    mutexWoken                   // 第1位：是否有等待者被唤醒
    mutexStarving                // 第2位：是否处于饥饿模式
    mutexWaiterShift = iota      // 从第3位开始是等待者的数量（用于原子管理）
    
    starvationThresholdNs = 1e6  // 等待超 1ms 将触发饥饿模式
)
```

---

## `Lock()` 方法详解

```go
func (m *Mutex) Lock() {
    if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
        if race.Enabled {
            race.Acquire(unsafe.Pointer(m))
        }
        return
    }
    m.lockSlow()
}
```

### 快路径
- `state == 0`：说明没有其他 goroutine 拥有该锁。
- 通过 CAS 将状态更新为 `mutexLocked`，表示成功获取锁。

### 慢路径
- 如果锁已被持有，则进入 `lockSlow()`。

---

## `lockSlow()` 慢路径详解（核心实现）

```go
func (m *Mutex) lockSlow() {
    var waitStartTime int64      // 开始等待的时间
    starving := false            // 当前是否处于饥饿状态
    awoke := false               // 当前 goroutine 是否刚被唤醒
    iter := 0                    // 自旋次数
    old := m.state
```

### 自旋 + 饥饿判断逻辑

```go
    for {
        if old&mutexStarving == 0 && runtime_canSpin(iter) {
            // 自旋条件满足（非饥饿模式且可自旋）
            if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 {
                if atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
                    awoke = true
                }
            }
            runtime_doSpin()
            iter++
            old = m.state
            continue
        }
```

- 如果未进入饥饿模式，且当前允许自旋（由调度器判断），则进行短暂自旋；
- 自旋前优先尝试设置 `mutexWoken` 位，避免唤醒多个 goroutine；

---

### 加入等待队列 & 状态更新

```go
        new := old
        if old&mutexStarving != 0 {
            // 饥饿模式下，不抢锁，只排队
            new |= mutexLocked
        } else {
            new |= mutexLocked
        }
        new += 1 << mutexWaiterShift

        if starving {
            new |= mutexStarving
        }

        if awoke {
            new &^= mutexWoken
        }

        if atomic.CompareAndSwapInt32(&m.state, old, new) {
            break
        }
        old = m.state
    }
```

- 设置 `mutexLocked` 位表示将尝试获得锁；
- `mutexWaiterShift` 是等待计数器，+1 表示新增等待者；
- 饥饿模式下会保留 `mutexStarving`，否则在某些条件下恢复正常模式；
- 成功 CAS 后跳出循环，准备阻塞；

---

### 开始阻塞

```go
    if waitStartTime == 0 {
        waitStartTime = runtime_nanotime()
    }

    runtime_Semacquire(&m.sema)
```

- 如果第一次进入等待，则记录当前时间；
- 阻塞当前 goroutine，直到被 `Unlock()` 唤醒。

---

### 被唤醒后的处理逻辑

```go
    starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
    old = m.state

    for {
        new := old - 1<<mutexWaiterShift
        if starving {
            new |= mutexStarving
        } else {
            new &^= mutexStarving
        }

        if atomic.CompareAndSwapInt32(&m.state, old, new) {
            break
        }
        old = m.state
    }

    if race.Enabled {
        race.Acquire(unsafe.Pointer(m))
    }
}
```

- 判断是否等待超过 1ms：超过则进入饥饿模式；
- 原子减少 waiter 数量；
- 最后，如果启用了 race 检查器，则报告锁的持有；

---

## `Unlock()` 解锁流程（补充）

```go
func (m *Mutex) Unlock() {
    new := atomic.AddInt32(&m.state, -mutexLocked)
    if new&mutexLocked != 0 {
        throw("sync: unlock of unlocked mutex")
    }
    if new>>mutexWaiterShift != 0 {
        m.unlockSlow(new)
    }
}
```

- 原子去除 `mutexLocked`；
- 如果还有等待者，调用 `unlockSlow()` 进行唤醒。

---

## `unlockSlow()` 唤醒等待者

```go
func (m *Mutex) unlockSlow(new int32) {
    for {
        old := m.state
        if old&mutexLocked != 0 {
            return
        }
        if old&mutexStarving == 0 {
            // 正常模式下尝试唤醒一个 waiter
            if old&mutexWoken == 0 {
                new := old | mutexWoken
                if atomic.CompareAndSwapInt32(&m.state, old, new) {
                    runtime_Semrelease(&m.sema)
                    return
                }
            }
        } else {
            // 饥饿模式：直接交出锁，跳过竞争
            runtime_Semrelease(&m.sema)
            return
        }
    }
}
```

- 正常模式：尝试设置 `mutexWoken`，然后唤醒一个等待者；
- 饥饿模式：直接唤醒第一个等待者，由它获得锁；

---

## 总结说明

```text
sync.Mutex 在高并发场景下通过以下机制实现性能与公平性的平衡：

1. 快路径加锁（无竞争）
2. 自旋优化短期竞争
3. 饥饿模式保证等待者不会永远饿死
4. 原子操作+信号量高效阻塞/唤醒 goroutine
5. 内置 race 检查支持并发调试
```