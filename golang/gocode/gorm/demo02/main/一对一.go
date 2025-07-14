package main

type User struct {
	ID       uint
	Name     string
	Age      int
	Gender   bool
	UserInfo UserInfo // 通过UserInfo可以拿到用户详情信息
}

type UserInfo struct {
	UserID uint // 外键（GORM会自动创建外键）
	ID     uint
	Addr   string
	Like   string
}

func main() {
	//err := DB.Debug().AutoMigrate(&User{}, &UserInfo{})
	//if err != nil {
	//	panic(err)
	//}

	//DB.Debug().Create(&User{
	//	Name:   "张三",
	//	Age:    18,
	//	Gender: true,
	//	UserInfo: UserInfo{
	//		Addr: "东莞市",
	//		Like: "洗浴",
	//	},
	//})

	//var userInfo UserInfo
	//DB.Take(&userInfo)
	//fmt.Println(userInfo)

	//var user User
	//DB.Preload("UserInfo").Take(&user)
	//fmt.Println(user)

	// 级联删除
	var user User
	DB.Take(&user)
	DB.Debug().Select("UserInfo").Delete(&user)

}
