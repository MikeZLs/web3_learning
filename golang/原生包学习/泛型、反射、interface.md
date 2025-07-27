## 1. 泛型（Generics）

### 泛型的定义
泛型是 Go 1.18 版本引入的特性，允许编写可以处理不同类型数据的通用代码，而不需要为每种类型重复编写相同的逻辑。

### 泛型的核心概念：
- **类型参数**：用方括号 `[]` 定义
- **类型约束**：限制类型参数可以接受的类型

### 示例代码：
```go
// 泛型函数 - 可以处理任何可比较的类型
func Min[T comparable](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// 泛型切片函数
func Map[T any, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// 泛型结构体
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

// 使用示例
func main() {
    // 使用泛型函数
    fmt.Println(Min(5, 3))        // 输出: 3
    fmt.Println(Min("hello", "world")) // 输出: hello
    
    // 使用泛型 Map
    numbers := []int{1, 2, 3, 4, 5}
    doubled := Map(numbers, func(n int) int { return n * 2 })
    fmt.Println(doubled) // 输出: [2 4 6 8 10]
    
    // 使用泛型 Stack
    intStack := Stack[int]{}
    intStack.Push(10)
    intStack.Push(20)
    fmt.Println(intStack.Pop()) // 输出: 20 true
}
```

## 2. 反射（Reflection）

### 反射的定义
反射是在运行时检查、修改程序结构和行为的能力。Go 的反射主要通过 `reflect` 包实现。

### 反射的核心概念：
- **Type**：表示类型信息
- **Value**：表示值信息
- **Kind**：表示底层类型种类

### 示例代码：
```go
import (
    "fmt"
    "reflect"
)

type Person struct {
    Name string `json:"name" validate:"required"`
    Age  int    `json:"age" validate:"min=0"`
}

func (p Person) Greet() string {
    return fmt.Sprintf("Hello, I'm %s", p.Name)
}

func examineType(i interface{}) {
    // 获取类型信息
    t := reflect.TypeOf(i)
    v := reflect.ValueOf(i)
    
    fmt.Printf("Type: %s, Kind: %s\n", t.Name(), t.Kind())
    
    // 如果是结构体，遍历字段
    if t.Kind() == reflect.Struct {
        for i := 0; i < t.NumField(); i++ {
            field := t.Field(i)
            value := v.Field(i)
            
            fmt.Printf("Field: %s, Type: %s, Value: %v\n", 
                field.Name, field.Type, value.Interface())
            
            // 获取标签
            jsonTag := field.Tag.Get("json")
            validateTag := field.Tag.Get("validate")
            fmt.Printf("  Tags - json: %s, validate: %s\n", jsonTag, validateTag)
        }
    }
    
    // 遍历方法
    for i := 0; i < t.NumMethod(); i++ {
        method := t.Method(i)
        fmt.Printf("Method: %s\n", method.Name)
    }
}

// 使用反射修改值
func modifyValue(x interface{}) {
    v := reflect.ValueOf(x)
    
    // 必须是指针才能修改
    if v.Kind() == reflect.Ptr && !v.IsNil() {
        v = v.Elem() // 获取指针指向的值
        
        if v.Kind() == reflect.String && v.CanSet() {
            v.SetString("Modified!")
        }
    }
}

// 使用反射调用方法
func callMethod(obj interface{}, methodName string) {
    v := reflect.ValueOf(obj)
    method := v.MethodByName(methodName)
    
    if method.IsValid() {
        result := method.Call(nil)
        if len(result) > 0 {
            fmt.Println("Method result:", result[0].Interface())
        }
    }
}

func main() {
    p := Person{Name: "Alice", Age: 30}
    
    // 检查类型
    examineType(p)
    
    // 修改值
    str := "original"
    fmt.Println("Before:", str)
    modifyValue(&str)
    fmt.Println("After:", str)
    
    // 调用方法
    callMethod(p, "Greet")
}
```

## 3. 接口（Interface）

### 接口的定义
接口是 Go 语言中定义方法集合的类型，是实现多态和抽象的核心机制。

### 接口的特点：
- **隐式实现**：不需要显式声明实现接口
- **鸭子类型**：只要实现了接口的所有方法，就自动满足接口
- **空接口**：`interface{}` 或 `any` 可以接收任何类型

### 示例代码：
```go
// 定义接口
type Writer interface {
    Write([]byte) (int, error)
}

type Reader interface {
    Read([]byte) (int, error)
}

// 组合接口
type ReadWriter interface {
    Reader
    Writer
}

// 带有更多方法的接口
type Animal interface {
    Speak() string
    Move() string
}

// 实现接口的结构体
type Dog struct {
    Name string
}

func (d Dog) Speak() string {
    return fmt.Sprintf("%s says: Woof!", d.Name)
}

func (d Dog) Move() string {
    return fmt.Sprintf("%s runs", d.Name)
}

type Cat struct {
    Name string
}

func (c Cat) Speak() string {
    return fmt.Sprintf("%s says: Meow!", c.Name)
}

func (c Cat) Move() string {
    return fmt.Sprintf("%s sneaks", c.Name)
}

// 使用接口作为参数
func MakeSound(a Animal) {
    fmt.Println(a.Speak())
}

// 类型断言和类型切换
func processInterface(i interface{}) {
    // 类型断言
    if str, ok := i.(string); ok {
        fmt.Printf("String value: %s\n", str)
    }
    
    // 类型切换
    switch v := i.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    case Animal:
        fmt.Printf("Animal: %s\n", v.Speak())
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}

// 空接口的使用
func PrintAnything(values ...interface{}) {
    for _, v := range values {
        fmt.Printf("Value: %v, Type: %T\n", v, v)
    }
}

// 接口嵌套
type Swimmer interface {
    Swim() string
}

type Flyer interface {
    Fly() string
}

type Duck interface {
    Animal
    Swimmer
    Flyer
}

type MallardDuck struct {
    Name string
}

func (m MallardDuck) Speak() string { return "Quack!" }
func (m MallardDuck) Move() string  { return "Waddle" }
func (m MallardDuck) Swim() string  { return "Swimming" }
func (m MallardDuck) Fly() string   { return "Flying" }

func main() {
    // 基本使用
    dog := Dog{Name: "Buddy"}
    cat := Cat{Name: "Whiskers"}
    
    MakeSound(dog)
    MakeSound(cat)
    
    // 接口切片
    animals := []Animal{dog, cat}
    for _, animal := range animals {
        fmt.Println(animal.Move())
    }
    
    // 类型处理
    processInterface("hello")
    processInterface(42)
    processInterface(dog)
    
    // 空接口
    PrintAnything(1, "two", 3.0, true)
    
    // 检查接口实现
    var _ Duck = MallardDuck{} // 编译时检查
}
```

## 三者的关系和区别

### 1. **使用场景对比**
- **泛型**：编译时类型安全的通用代码
- **反射**：运行时类型检查和操作
- **接口**：定义行为契约，实现多态

### 2. **性能对比**
- **泛型**：编译时确定类型，性能最好
- **接口**：有小的运行时开销（动态分派）
- **反射**：性能开销最大，应谨慎使用

### 3. **实际应用示例**
```go
// 结合使用的例子
type Serializer interface {
    Serialize(v interface{}) ([]byte, error)
}

// 泛型序列化函数
func SerializeList[T any](items []T, s Serializer) ([][]byte, error) {
    results := make([][]byte, len(items))
    for i, item := range items {
        data, err := s.Serialize(item)
        if err != nil {
            return nil, err
        }
        results[i] = data
    }
    return results, nil
}

// 使用反射的 JSON 序列化器
type JSONSerializer struct{}

func (j JSONSerializer) Serialize(v interface{}) ([]byte, error) {
    // 内部使用反射实现
    return json.Marshal(v)
}
```