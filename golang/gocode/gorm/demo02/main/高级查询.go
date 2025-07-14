package main

import (
	"fmt"
	"gorm.io/gorm"
)

type Student struct {
	ID     uint   `gorm:"size:3"`
	Name   string `gorm:"size:8" json:"name"`
	Age    int    `gorm:"size:3"`
	Gender bool
	Email  *string `gorm:"size:32"`
}

func PtrString(email string) *string {
	return &email
}

func main() {
	DB = DB.Session(&gorm.Session{
		Logger: mysqlLogger,
	})
	//var studentList []Student
	// // 自动迁移，创建表
	//err := DB.Debug().AutoMigrate(&Student{})
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//DB.Find(&studentList).Delete(&studentList)
	//studentList = []Student{
	//	{ID: 1, Name: "李元芳", Age: 32, Email: PtrString("lyf@yf.com"), Gender: true},
	//	{ID: 2, Name: "张武", Age: 18, Email: PtrString("zhangwu@lly.cn"), Gender: true},
	//	{ID: 3, Name: "枫枫", Age: 23, Email: PtrString("ff@yahoo.com"), Gender: true},
	//	{ID: 4, Name: "刘大", Age: 54, Email: PtrString("liuda@qq.com"), Gender: true},
	//	{ID: 5, Name: "李武", Age: 23, Email: PtrString("liwu@lly.cn"), Gender: true},
	//	{ID: 6, Name: "李琦", Age: 14, Email: PtrString("liqi@lly.cn"), Gender: false},
	//	{ID: 7, Name: "晓梅", Age: 25, Email: PtrString("xiaomeo@sl.com"), Gender: false},
	//	{ID: 8, Name: "如燕", Age: 26, Email: PtrString("ruyan@yf.com"), Gender: false},
	//	{ID: 9, Name: "魔灵", Age: 21, Email: PtrString("moling@sl.com"), Gender: true},
	//}
	//DB.Create(&studentList)

	//var users []Student
	// 查询用户名是枫枫的
	//DB.Where("name = ?", "枫枫").Find(&users)
	//DB.Find(&users, "name = ?", "枫枫")
	//fmt.Println(users)
	//
	//// 查询用户名不是枫枫的（三种写法）
	////DB.Where("not name = ?", "枫枫").Find(&users)
	////DB.Not("name = ?", "枫枫").Find(&users)
	////DB.Where("name <> ?", "枫枫").Find(&users)
	//fmt.Println(users)
	//
	//// 查询用户名包含 如燕，李元芳的
	//DB.Where("name in ?", []string{"如燕", "李元芳"}).Find(&users)
	//fmt.Println(users)
	//
	//// 查询姓李的
	//DB.Where("name like ?", "李%").Find(&users)
	//fmt.Println(users)
	//DB.Where("name like ?", "李_").Find(&users) // 查询姓李并且名字为两个字的人
	//fmt.Println(users)

	// 查询年龄大于23,并且是qq邮箱的
	//DB.Where("age > ? and email like ?", "23", "%@qq.com").Find(&users)
	//DB.Where("age > ?", 23).Where("email like ?", "%@qq.com").Find(&users)
	//fmt.Println(users)

	// 查询是qq邮箱的，或者是女的
	//DB.Where("gender = ? or email like ?", false, "%@qq.com").Find(&users)
	//DB.Where("gender = ?", false).Or("email like ?", "%@qq.com").Find(&users)
	//fmt.Println(users)

	//// 使用结构体查询：结构体中的值只能是等于，不能是其他条件
	//DB.Where(&Student{Age: 32, Name: "李元芳"}).Find(&users)
	//fmt.Println(users)
	//
	//// 使用 map 查询
	//DB.Where(map[string]any{"name": "李元芳", "age": 32}).Find(&users)
	//fmt.Println(users)

	//	使用 Select 筛选字段
	//DB.Select("name", "age").Find(&users)
	//DB.Select([]string{"name", "age"}).Find(&users)
	//fmt.Println(users)

	//// 将数据扫描到一个新的结构体中
	//type User struct {
	//	Username string `gorm:"column:name"` // Scan 是根据 column 列名进行扫描的
	//	Age      int
	//}
	//
	//var userList []User
	//
	////DB.Select([]string{"name", "age"}).Find(&users).Scan(&userList)
	////DB.Select("name", "age").Find(&users).Scan(&userList)
	////DB.Table("students").Select("name", "age").Scan(&userList) // 直接指定表名，省略Find，可以减少查询全表数据步骤，提高效率
	//DB.Model(Student{}).Select("name", "age").Scan(&userList) // 使用 Model 指定结构体，不用知道表名
	//fmt.Println(userList)

	//	排序查询
	//DB.Order("age desc").Find(&studentList)
	//fmt.Println(studentList)

	//	分页查询
	//DB.Limit(2).Find(&studentList) // limit 2
	//DB.Limit(2).Offset(2).Find(&studentList) // limit 2 offset 2
	//DB.Limit(2).Offset(4).Find(&studentList) // limit 2 offset 4
	//DB.Limit(2).Offset(6).Find(&studentList) // limit 2 offset 6
	// 分页查询通用写法
	//limit := 2
	//page := 1
	//DB.Limit(limit).Offset((page - 1) * limit).Find(&studentList)
	//fmt.Println(studentList)

	//	//去重
	//var ageList []int
	//DB.Model(&Student{}).Select("age").Distinct().Scan(&ageList) // 按 age 进行去重
	//DB.Model(&Student{}).Select("distinct age").Scan(&ageList)
	//fmt.Println(ageList)

	////	分组查询
	//type AggeGroup struct {
	//	Gender int
	//	Count  int
	//	Names  string
	//}
	//var groupList []AggeGroup
	// 查询男生的个数和女生的个数
	//DB.Table("students").
	//	Select(
	//		"count(id) as count",
	//		"gender",
	//		"group_concat(name) as names",
	//	).
	//	Group("gender").
	//	Scan(&groupList)

	//DB.Raw("SELECT count(id) as count,gender,group_concat(name) as names FROM `students` GROUP BY gender").Scan(&groupList)
	//DB.Raw("SELECT count(id) as count,gender,group_concat(name) as names FROM `students` where age > ? GROUP BY gender", "18").  // 用 ？ 来占位传参
	//Scan(&groupList)
	//fmt.Println(groupList)

	//	子查询
	//DB.Raw("select * from students where age > (select avg(age) from students)").Find(&studentList)
	//DB.Raw("select * from students where age > (select avg(age) from students)").Scan(&studentList)
	//DB.Where("age > (?)", DB.Model(Student{}).Select("avg(age)")).Find(&studentList)
	//fmt.Println(studentList)

	// 命名参数
	//DB.Where("name = ? and age = ?", "刘大", 54).Find(&studentList)
	//DB.Where("name = @name and age = @age",
	//	sql.Named("name", "刘大"),
	//	sql.Named("age", 54)).Find(&studentList)
	//DB.Where("name = @name and age = @age",
	//	map[string]any{"name": "刘大", "age": 54}).Find(&studentList)
	//fmt.Println(studentList)

	// find 到 map
	var res []map[string]any
	//DB.Model(&Student{}).Find(&res)
	//DB.Table("students").Find(&res)
	//fmt.Println(res)

	//DB.Scopes(age23).Find(&studentList)
	//fmt.Println(studentList)

	DB.Model(&Student{}).Scopes(age23).Scopes(genderMale).Find(&res)
	fmt.Println(res)
}

// 定义一个可以在 Model 层引用的通用方法
func age23(db *gorm.DB) *gorm.DB {
	return db.Where("age > ?", 23)
}

func genderMale(db *gorm.DB) *gorm.DB {
	return db.Where("gender = 1")
}
