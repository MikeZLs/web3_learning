## 1. 泛型（Generics）

### 什么是泛型？
泛型是 Go 1.18 版本引入的特性，允许编写可以处理不同类型数据的通用代码，而不需要为每种类型重复编写相同的逻辑。

### 泛型的核心原理

#### 1. **类型参数化**
泛型的本质是将类型作为参数传递，在编译时进行类型特化：

```go
// 编译前的泛型代码
func Max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// 编译器会根据调用生成特化版本
// Max[int](1, 2) 会生成：
func Max_int(a, b int) int {
    if a > b {
        return a
    }
    return b
}

// Max[string]("a", "b") 会生成：
func Max_string(a, b string) string {
    if a > b {
        return a
    }
    return b
}
```

#### 2. **类型约束（Type Constraints）**
类型约束定义了类型参数必须满足的条件：

```go
// 类型约束的本质是接口
type Ordered interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64 | ~string
}

// ~ 表示底层类型，允许自定义类型
type MyInt int // MyInt 的底层类型是 int，满足 ~int

// 自定义约束
type Number interface {
    ~int | ~float64
    String() string  // 还可以要求实现特定方法
}
```

#### 3. **类型推断**
Go 编译器可以自动推断类型参数：

```go
func Print[T any](slice []T) {
    for _, v := range slice {
        fmt.Println(v)
    }
}

// 调用时可以省略类型参数
Print([]int{1, 2, 3})      // 自动推断 T = int
Print([]string{"a", "b"})  // 自动推断 T = string
```

#### 4. **泛型的实现原理（Stenciling vs. Dictionary）**
Go 采用了部分单态化（partial monomorphization）的策略：

```go
// GCShape: 相同内存布局的类型共享实现
// 指针类型共享一个实现，值类型根据大小分组

type Container[T any] struct {
    value T
}

// Container[*int] 和 Container[*string] 共享实现
// Container[int] 和 Container[string] 可能有不同实现
```

### 泛型完整示例：

```go
// 泛型接口
type Comparable[T any] interface {
    CompareTo(other T) int
}

// 泛型结构体
type BinaryTree[T constraints.Ordered] struct {
    value T
    left  *BinaryTree[T]
    right *BinaryTree[T]
}

func (t *BinaryTree[T]) Insert(value T) {
    if value < t.value {
        if t.left == nil {
            t.left = &BinaryTree[T]{value: value}
        } else {
            t.left.Insert(value)
        }
    } else {
        if t.right == nil {
            t.right = &BinaryTree[T]{value: value}
        } else {
            t.right.Insert(value)
        }
    }
}

// 泛型方法
func (t *BinaryTree[T]) Contains(value T) bool {
    if t == nil {
        return false
    }
    if value == t.value {
        return true
    }
    if value < t.value {
        return t.left.Contains(value)
    }
    return t.right.Contains(value)
}
```

## 2. 反射（Reflection）

### 什么是反射？
反射是在运行时检查、修改程序结构和行为的能力。Go 的反射主要通过 `reflect` 包实现。

### 反射的核心原理

#### 1. **接口的内部表示**
Go 中每个接口值都包含两个指针：

```go
// 空接口 interface{} 的内部结构
type eface struct {
    _type *_type         // 类型信息指针
    data  unsafe.Pointer // 数据指针
}

// 非空接口的内部结构
type iface struct {
    tab  *itab           // 接口表（包含类型信息和方法表）
    data unsafe.Pointer  // 数据指针
}
```

#### 2. **Type 和 Value 的实现**
```go
// reflect.Type 的核心实现
type Type interface {
    // 基本信息
    Name() string
    Kind() Kind
    Size() uintptr
    
    // 结构体相关
    NumField() int
    Field(i int) StructField
    
    // 方法相关
    NumMethod() int
    Method(int) Method
    
    // ... 更多方法
}

// reflect.Value 的核心结构
type Value struct {
    typ *rtype         // 类型信息
    ptr unsafe.Pointer // 数据指针
    flag               // 标志位（可寻址、可设置等）
}
```

#### 3. **反射的工作原理**
```go
func demonstrateReflection() {
    var x float64 = 3.14
    
    // 1. interface{} 转换
    var i interface{} = x  // x 被装箱成 eface
    
    // 2. 获取 Type 和 Value
    t := reflect.TypeOf(i)  // 提取 eface._type
    v := reflect.ValueOf(i) // 创建 Value 结构
    
    // 3. 类型检查
    fmt.Println("Type:", t)
    fmt.Println("Kind:", t.Kind())
    
    // 4. 值操作（需要通过指针）
    p := reflect.ValueOf(&x)
    elem := p.Elem()  // 获取指针指向的元素
    if elem.CanSet() {
        elem.SetFloat(2.71)
    }
}
```

#### 4. **反射的性能开销**
```go
// 直接调用
func DirectCall(p *Person) string {
    return p.Name  // 编译时确定偏移量，直接内存访问
}

// 反射调用
func ReflectCall(i interface{}) string {
    v := reflect.ValueOf(i)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    // 运行时查找字段、检查类型、计算偏移量
    return v.FieldByName("Name").String()
}
```

### 反射高级示例：

```go
// 通用的结构体标签验证器
func Validate(obj interface{}) error {
    v := reflect.ValueOf(obj)
    t := reflect.TypeOf(obj)
    
    // 处理指针
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
        t = t.Elem()
    }
    
    if t.Kind() != reflect.Struct {
        return fmt.Errorf("validate: not a struct")
    }
    
    // 遍历所有字段
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        
        // 获取验证标签
        tag := field.Tag.Get("validate")
        if tag == "" {
            continue
        }
        
        // 解析并应用验证规则
        rules := strings.Split(tag, ",")
        for _, rule := range rules {
            if err := applyRule(rule, value); err != nil {
                return fmt.Errorf("field %s: %v", field.Name, err)
            }
        }
    }
    
    return nil
}

// 动态方法调用
func DynamicCall(obj interface{}, methodName string, args ...interface{}) ([]reflect.Value, error) {
    v := reflect.ValueOf(obj)
    method := v.MethodByName(methodName)
    
    if !method.IsValid() {
        return nil, fmt.Errorf("method %s not found", methodName)
    }
    
    // 准备参数
    in := make([]reflect.Value, len(args))
    for i, arg := range args {
        in[i] = reflect.ValueOf(arg)
    }
    
    // 调用方法
    return method.Call(in), nil
}
```

## 3. 接口（Interface）

### 什么是接口？
接口是 Go 语言中定义方法集合的类型，是实现多态和抽象的核心机制。

### 接口的核心原理

#### 1. **接口的内存布局**
```go
// 接口表（itab）结构
type itab struct {
    inter *interfacetype // 接口类型信息
    _type *_type         // 具体类型信息
    hash  uint32         // 类型哈希，用于快速比较
    _     [4]byte        // 内存对齐
    fun   [1]uintptr     // 方法表（可变长度）
}

// 接口赋值的过程
var w io.Writer
var f *os.File = os.Stdout
w = f  // 创建 iface{tab: itab{inter: io.Writer, _type: *os.File}, data: f}
```

#### 2. **动态分派机制**
```go
// 接口方法调用的汇编伪代码
// w.Write([]byte("hello"))
//
// 1. 从接口值加载 itab
// 2. 从 itab.fun 加载对应方法地址
// 3. 调用方法

type Writer interface {
    Write([]byte) (int, error)
}

// 编译器生成的方法表
// itab.fun[0] = (*os.File).Write
// itab.fun[1] = (*os.File).Close  // 如果接口有多个方法
```

#### 3. **类型断言的实现**
```go
// 类型断言的底层实现
func TypeAssert() {
    var i interface{} = "hello"
    
    // s := i.(string) 的底层过程：
    // 1. 检查 i.tab._type 是否等于 string 的类型信息
    // 2. 如果相等，返回 i.data
    // 3. 如果不等，panic
    
    // s, ok := i.(string) 的底层过程：
    // 1. 检查 i.tab._type 是否等于 string 的类型信息
    // 2. 如果相等，返回 i.data 和 true
    // 3. 如果不等，返回零值和 false
}
```

#### 4. **空接口的优化**
```go
// 空接口不需要方法表
type eface struct {
    _type *_type
    data  unsafe.Pointer
}

// 小对象优化（未实现，但是设计考虑）
// 某些小值可以直接存储在 data 字段中，避免堆分配
```

### 接口高级示例：

```go
// 接口组合和嵌入
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    ReadWriter
    Closer
}

// 接口的零值
func NilInterface() {
    var w io.Writer           // w == nil (tab 和 data 都是 nil)
    var f *os.File            // f == nil
    w = f                     // w != nil (tab 不是 nil，但 data 是 nil)
    
    if w != nil {
        fmt.Println("w is not nil")  // 会执行
    }
    
    // 正确的 nil 检查
    if w != nil && reflect.ValueOf(w).IsNil() {
        fmt.Println("w contains nil pointer")
    }
}

// 类型嵌入和方法提升
type MyWriter struct {
    io.Writer  // 嵌入接口
}

// MyWriter 自动实现了 io.Writer 接口

// 接口的私有方法（防止外部实现）
type Token interface {
    token()  // 私有方法，只有本包的类型能实现
}

type myToken struct{}
func (myToken) token() {}  // 实现私有方法
```

## 三者的对比和协同

### 性能对比
```go
// 性能测试示例
func BenchmarkDirect(b *testing.B) {
    s := &MyStruct{Value: 42}
    for i := 0; i < b.N; i++ {
        _ = s.Value  // 直接访问：最快
    }
}

func BenchmarkGeneric(b *testing.B) {
    s := &MyStruct{Value: 42}
    for i := 0; i < b.N; i++ {
        _ = GetValue(s)  // 泛型：编译时特化，接近直接访问
    }
}

func BenchmarkInterface(b *testing.B) {
    var s Valuer = &MyStruct{Value: 42}
    for i := 0; i < b.N; i++ {
        _ = s.GetValue()  // 接口：动态分派，略有开销
    }
}

func BenchmarkReflection(b *testing.B) {
    s := &MyStruct{Value: 42}
    v := reflect.ValueOf(s).Elem()
    f := v.FieldByName("Value")
    for i := 0; i < b.N; i++ {
        _ = f.Int()  // 反射：运行时解析，开销最大
    }
}
```

### 实际应用中的结合使用
```go
// ORM 框架的简化实现，结合三者
type Model interface {
    TableName() string
}

// 泛型 CRUD 操作
type Repository[T Model] struct {
    db *sql.DB
}

func (r *Repository[T]) Create(entity T) error {
    // 使用反射获取字段信息
    v := reflect.ValueOf(entity)
    t := reflect.TypeOf(entity)
    
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
        t = t.Elem()
    }
    
    // 构建 SQL
    tableName := entity.TableName()  // 接口方法
    
    var columns []string
    var values []interface{}
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        if tag := field.Tag.Get("db"); tag != "" && tag != "-" {
            columns = append(columns, tag)
            values = append(values, v.Field(i).Interface())
        }
    }
    
    query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
        tableName,
        strings.Join(columns, ", "),
        strings.Repeat("?, ", len(columns)-1) + "?")
    
    _, err := r.db.Exec(query, values...)
    return err
}

// 使用示例
type User struct {
    ID   int    `db:"id"`
    Name string `db:"name"`
}

func (User) TableName() string { return "users" }

func main() {
    repo := &Repository[User]{db: db}
    user := User{Name: "Alice"}
    repo.Create(user)  // 泛型提供类型安全，接口提供多态，反射处理结构
}
```

这三个特性各有其适用场景：
- **泛型**：编译时的类型抽象，性能好，类型安全
- **接口**：运行时的行为抽象，实现多态
- **反射**：运行时的完全动态性，灵活但有性能开销

理解它们的原理有助于在合适的场景选择合适的工具。