## Once
once是只执行一次动作的对象。
如果once.Do(f)被多次调用，只有第一次调用会执行f，即使f每次调用 Do 提供的f值不同

常见应用场景：
- 初始化全局配置
- 懒加载资源（如数据库连接）
- 注册某些只需注册一次的组件（如 HTTP 路由）

### sync.once源码
```go
package sync

import (
	"sync/atomic"
)

type Once struct {
	_ noCopy
	done atomic.Uint32
	m    Mutex
}
func (o *Once) Do(f func()) {
	// 使用 atomic.Load() 无锁读取 done 状态，如果已是 1（即已执行），直接返回，不加锁，提高性能
	if o.done.Load() == 0 {
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	// 加锁：防止多个 goroutine 同时执行
	o.m.Lock()
	
	defer o.m.Unlock()
	// 再次检查 done 状态，因为多个 goroutine 可能同时进入了快路径判断为未完成，必须在锁保护下再次检查，防止重复执行
	if o.done.Load() == 0 {
		//  保证无论 f() 是否 panic，done 都会被设置，后续不再执行 f，防止重复
		defer o.done.Store(1)
		f()
	}
}
```

---
## 锁
| 类型 | 核心机制 | 主要使用场景 |
| :--- | :--- | :--- |
| `sync.Mutex` | **互斥**。保证同一时间只有一个goroutine访问临界区。 | 对共享资源进行写操作，或读写操作都很频繁的场景。 |
| `sync.RWMutex` | **读写分离**。允许多个读操作并发进行，但写操作独占。 | 读多写少的场景，能显著提升并发读取的性能。 |
| `sync.WaitGroup` | **计数器**。等待一组任务完成。 | 并发任务编排，等待所有子goroutine结束。 |
| `sync.Once` | **单次执行**。保证一个函数只会被执行一次。 | 懒加载、单例模式、全局配置初始化。 |
| `sync.Cond` | **条件等待/通知**。让goroutine等待特定条件满足。 | 复杂的同步场景，如生产者-消费者模型。 |

---
## Pool
sync.Pool：是Go标准库中的一个轻量级对象缓存池，用于 存储和复用临时对象，以减少频繁的内存分配和 GC（垃圾回收）开销，提升性能

##### Pool的底层特点：
- 线程安全：适用于并发环境，内部有每个 P（逻辑处理器）对应的私有对象池，减少竞争。
- 对象可能被 GC 清除：池中的对象生命周期不确定，GC 时可能被清空。
- 不会持久保存数据：适合“临时”对象复用，不适合用作永久缓存。
- 对象可能“不干净”: 从 Pool 中 Get 出来的对象是被复用的，它可能保留了上一次使用时的状态。因此，每次获取对象后，必须手动重置其状态（buffer.Reset()）。

---
## Map
- 提供并发安全的键值存储
- 支持高效读操作，低锁争用
- 适合读多写少或写入分布广泛的场景

### sync.Map 结构：
```go
type Map struct {
    mu     Mutex                         // 锁保护 dirty 表
    read   atomic.Pointer[readOnly]     // 原子只读快表
    dirty  map[any]*entry               // 脏表，存储新写入的数据
    misses int                           // miss 计数器
}

type readOnly struct {
    m        map[any]*entry              // 快表映射
    amended  bool                        // 是否存在脏数据
}

type entry struct {
    p atomic.Pointer[any]               // 值的原子指针
}
```
### 核心原理：双 Map 机制：`read` + `dirty`

| 名称  | 作用                 | 是否并发安全  | 是否锁保护   |
|-------|----------------------|---------|---------|
| `read`  | 快速只读路径          | 是（原子）   | 否（无需加锁） |
| `dirty` | 新增/更新路径        | 否（需要加锁） | 是       |
- `read` 存的是主路径，允许无锁访问
- `dirty` 是辅助路径，写操作或新 key 会进入其中

### sync.Map 内部操作逻辑：

#### Load
```go
m.Load(key)
```
1. 优先查 `read`
2. 如果没命中且 `amended == true`，查 `dirty`
3. 如果从 dirty 查到，增加 `misses` 计数

#### Store
```go
m.Store(key, value)
```
1. 若 key 在 `read` 中 → 直接更新（CAS）
2. 若不在 → 加锁写入 `dirty`
3. 若是第一次写 → 创建 `dirty`，更新 `amended = true`

#### LoadOrStore
```go
m.LoadOrStore(key, value)
```
- 若 key 已存在 → 返回旧值
- 否则写入并返回新值

#### Delete / LoadAndDelete
```go
m.Delete(key)
```
- 设置 `entry.p = nil`
- 延迟清除（惰性删除）

#### miss 触发 read 更新
```go
func (m *Map) missLocked()
```
- 每次 miss +1
- 如果 `misses >= len(dirty)`，将 dirty 提升为新的 read
- 清空 dirty，重置 misses

#### Clear
```go
m.Clear()
```
- 清空 read 和 dirty
- 重置为初始状态

#### Range
```go
m.Range(func(key, value any) bool)
```
- 遍历 read 表
- 若 `amended == true`，提升 dirty 到 read 后再遍历

#### entry 特殊状态
| 状态       | 值              | 含义 |
|------------|------------------|------|
| 正常       | 指向实际值       | 有效数据 |
| 删除       | `nil`            | 删除但仍在表中 |
| 已驱逐     | `expunged` 指针 | 已永久标记为无效，不在 dirty 中存在 |


### 适合场景：
- 缓存、注册表、插件管理器等
- 只写一次、多次读取
- 读多写少

### 不适合场景：
- 高频写入（大量 key 更新）
- 对 key 精细控制或严格顺序的场景

## Singleflight （ 扩展包 "golang.org/x/sync/singleflight" ）
Singleflight 核心方法源码：
```go
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	g.mu.Lock() // 1. 加锁，保护 map 的并发访问
	if g.m == nil {
		g.m = make(map[string]*call) // 2. 懒加载初始化 map
	}
	
	// 3. 检查 key 是否已存在 (即是否已有相同的请求在进行中)
	if c, ok := g.m[key]; ok {
		c.dups++      // 4. 如果存在，说明是重复请求，dups++
		g.mu.Unlock() // 5. 解锁，让其他 goroutine 可以继续
		c.wg.Wait()   // 6.【关键】等待第一个请求完成

		// ... 处理 panic 和 Goexit 的逻辑 ...
		
		return c.val, c.err, true // 7. 返回共享的结果
	}
	
	// 8. 如果 key 不存在，说明是第一个请求
	c := new(call)
	c.wg.Add(1)      // 9. 初始化 WaitGroup，计数器+1
	g.m[key] = c     // 10. 将 call 实例存入 map
	g.mu.Unlock()    // 11. 解锁

	g.doCall(c, key, fn) // 12.【关键】实际执行函数调用
	return c.val, c.err, c.dups > 0 // 13. 返回自己的执行结果
}
```
与 sync.Once 的区别：
- sync.Once  某件事全局只做一次（无参数、无 key）
- singleflight  某一类重复的请求，正在执行时其他人等着共享结果（有 key）

