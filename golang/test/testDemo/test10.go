package testDemo

import "fmt"

// 使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，
// 再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。
// 为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。

// 定义一个 Person 结构体
type Person struct {
	Name string
	Age  int
}

// 定义一个 Employee 结构体，组合 Person 结构体,并添加 EmployeeID 字段
type Employee struct {
	Person
	EmployeeID int
}

// 为 Employee 结构体实现一个 PrintInfo() 方法
func (e Employee) PrintInfo() {
	fmt.Printf("员工姓名：%s，年龄：%d，员工ID：%d\n", e.Name, e.Age, e.EmployeeID)
}

func Test10() {
	// 创建一个 Employee 实例
	e := Employee{
		Person: Person{
			Name: "张三",
			Age:  30,
		},
		EmployeeID: 10010,
	}

	e.PrintInfo()
}
