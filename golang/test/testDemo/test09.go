package testDemo

import (
	"fmt"
	"math"
)

// 定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。
// 然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
// 在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。

// Shape 定义一个 Shape 接口
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle 定义一个 Rectangle 结构体，
type Rectangle struct {
	Width, Height float64
}

// Rectangle 实现 Shape 接口
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// 定义 Circle 结构体
type Circle struct {
	Radius float64
}

// Circle 实现 Shape 接口
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func Test09() {
	// 创建 Rectangle 和 Circle 实例
	r := Rectangle{Width: 5, Height: 3}
	c := Circle{Radius: 2.5}

	// // 使用接口变量
	// var s Shape

	// s = r
	fmt.Println("矩形:")
	// fmt.Printf("面积: %.2f\n", s.Area()) // 通过接口调用方法
	fmt.Printf("面积: %.2f\n", r.Area()) // 直接调用结构体方法
	fmt.Printf("周长: %.2f\n", r.Perimeter())

	// s = c
	fmt.Println("\n圆形:")
	fmt.Printf("面积: %.2f\n", c.Area())
	fmt.Printf("周长: %.2f\n", c.Perimeter())
}
