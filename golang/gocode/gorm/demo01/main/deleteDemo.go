package main

func DeleteTable() {
	var student Student
	//单条数据删除
	DB.Delete(&student, 20)
	//批量删除
	DB.Delete(&student, []int{18, 20, 21})
	DB.Debug().Delete(&student, "name in (?)", []string{"张三4", "张三8"})

}
