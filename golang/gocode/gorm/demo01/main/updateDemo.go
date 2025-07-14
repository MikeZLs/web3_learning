package main

func updateTable() {
	////	save更新（单个数据）
	//var student Student
	//DB.Take(&student, 12)
	//student.Name = "广东吴彦祖"
	//student.Age = 1
	//student.Email = nil
	//err := DB.Debug().Save(&student).Error
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("save student success")

	//	update更新（批量更新）
	var studentList []Student
	//// 批量更新单个字段
	//err := DB.Find(&studentList, "age = ?", 18).Update("email", "666@qq.com").Error // 第一种写法
	////err := DB.Model(&studentList).Where("age = ?", 18).Update("email", "888@qq.com").Error // 第二种写法
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("save student success")

	// 批量更新多个字段
	//DB.Debug().Find(&studentList, "age = ?", 18).Updates(
	//	Student{
	//		Age:    66,
	//		Gender: false,
	//	}) // 使用结构体不能更新零值
	//	想要更新零值，使用 map
	DB.Debug().Find(&studentList, "age = ?", 18).Updates(map[string]any{
		"age":    18,
		"gender": true,
	})
}
