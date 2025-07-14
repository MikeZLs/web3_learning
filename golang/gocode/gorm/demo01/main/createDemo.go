package main

import "fmt"

type Student struct {
	Id   uint   `gorm:"size:10;primaryKey"` // 默认id为主键，使用primaryKey可以显示指定主键字段
	Name string `gorm:"size:16" json:"name"`
	//Name  string  `gorm:"type:varchar(16)"` // 另一种定义字段类型和大小的写法
	Age    int     `gorm:"size:3"`
	Email  *string `gorm:"size:128"`                      // 使用指针是为了存空值
	Type   string  `gorm:"size:4;column:_type"`           // column:指定字段名
	Date   string  `gorm:"default:2001-01-01;comment:日期"` // default:给字段赋默认值；comment:字段注释
	Gender bool
}

func CreateTable() {
	err := DB.Debug().AutoMigrate(&Student{}) // DB后使用Debug()方法来显示Debug信息
	if err != nil {
		fmt.Println("创建表失败", err)
		return
	}
}
