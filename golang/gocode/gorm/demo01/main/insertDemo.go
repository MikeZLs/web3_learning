package main

func InsertTable() {

	//email := "12345678@163.com"

	// 单条插入
	//s1 := Student{
	//	Name:   "张三",
	//	Age:    18,
	//	Email:  &email, // 只有指针类型可以传入 nil ，在数据库中为 null
	//	Date:   "2002-02-02",
	//	Gender: true,
	//}

	// Create() 接收的是一个指针，而不是值
	// 由于我们传递的是一个指针，调用完Create() 之后，s1 这个对象上面就有该记录的信息了
	//err := DB.Create(&s1).Error
	//if err != nil {
	//	fmt.Println(err)
	//}

	//fmt.Println("写入成功")

	//	批量插入
	//var studentList []Student
	//
	//for i := 0; i < 10; i++ {
	//	studentList = append(studentList, Student{
	//		Name:   fmt.Sprintf("张三%d", i+1),
	//		Age:    18 + i,
	//		Email:  &email,
	//		Date:   "2002-02-02",
	//		Gender: true,
	//	})
	//}
	//err := DB.Create(&studentList).Error
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("插入成功")

	s2 := Student{
		Name:   "张三66",
		Age:    18,
		Date:   "2002-02-02",
		Gender: true,
	}
	DB.Create(&s2)
}
