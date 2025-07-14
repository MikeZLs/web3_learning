package main

import (
	"fmt"
	"gorm.io/gorm"
)

func SelectTable() {
	// 单条记录的查询
	//var student Student
	//DB = DB.Session(&gorm.Session{
	//	Logger: mysqlLogger,
	//})
	//DB.Take(&student)
	//fmt.Println(student)
	//student = Student{}
	//
	//DB.First(&student)
	//fmt.Println(student)
	//student = Student{}
	//
	//DB.Last(&student)
	//fmt.Println(student)
	//student = Student{}
	//
	//err := DB.Take(&student, 13).Error
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(student)

	//DB.Take(&student, "name = ?", "张三3") // 使用 ? 来作为占位符 ，有效防止 sql 注入
	//fmt.Println(student)

	//	查询多条记录
	var studentList []Student
	DB = DB.Session(&gorm.Session{
		Logger: mysqlLogger,
	})

	// 获取查询记录数
	//count := DB.Find(&studentList, "age = ?", "18").RowsAffected
	//count := DB.Find(&studentList).RowsAffected
	//fmt.Println(count)
	//
	//DB.Find(&studentList)
	////for _, student := range studentList {
	////	fmt.Println(student)
	////}
	//
	//data, _ := json.Marshal(studentList) // 将数据转为 json 格式
	//fmt.Println(string(data))

	// 根据主键查询
	//DB.Find(&studentList, []int{1, 11, 16}) // id = 11 没匹配到数据不会报错，返回匹配的数据
	//fmt.Println(studentList)

	// 根据其他条件查询
	//DB.Find(&studentList, "name in (?)", []string{"张三1", "张三5", "张三9"})
	//DB.Where("name in ?", []string{"张三1", "张三5", "张三9"}).Find(&studentList) // GORM 查询是链式构建语法，先构造条件，最后调用 .Find() 来执行查询
	//DB.Find(&studentList, "name in ? and age in ?", []string{"张三1", "张三5", "张三9"}, []int{18, 26})
	DB.Where("name in ? and age in ?", []string{"张三1", "张三5", "张三9"}, []int{18, 26}).Find(&studentList)
	fmt.Println(studentList)

}
