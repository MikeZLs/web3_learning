package main

// User 用户表 一个用户拥有多篇文章
//type User struct {
//	ID       uint      `gorm:"size:4"`
//	Name     string    `gorm:"size:8"`
//	Articles []Article // 用户拥有的文章列表
//}

// Article 文章表 一篇文章只属于一个用户
//type Article struct {
//	ID     uint   `gorm:"size:4"`
//	Title  string `gorm:"size:16"`
//	UserID uint   `gorm:"size:4"` // 属于   这里的类型要和引用的外键类型一致，包括大小
//	User   User   // 属于
//}

func main() {
	//err := DB.Debug().AutoMigrate(&User{}, &Article{})
	//if err != nil {
	//	return
	//}
	//fmt.Println("创建成功！")

	//// 创建用户，带上文章
	//DB.Debug().Create(&User{
	//	Name: "王五",
	//	Articles: []Article{
	//		{
	//			Title: "Python",
	//		}, {
	//			Title: "Php",
	//		},
	//	},
	//})

	////	创建文章，关联已有用户
	//DB.Debug().Create(&Article{
	//	Title:  "开始阅读Golang文档",
	//	UserID: 3,
	//})

	//DB.Debug().Create(&Article{
	//	Title: "开始学习Python",
	//	User: User{
	//		Name: "王五",
	//	},
	//})

	//var user User
	//DB.Take(&user, 1)
	//DB.Create(&Article{
	//	Title: "张三写的C++",
	//	User:  user,
	//})

	//	预加载连表查
	//var articles Article
	//DB.Debug().Preload("User").Take(&articles)
	//fmt.Println(articles)

	//var users User
	//DB.Debug().Preload("Articles").Take(&users)
	//fmt.Println(users)

	//var users User
	//// 预加载函数 Preload() 可以带过滤条件
	//DB.Debug().Preload("Articles", "id > ?", 2).Take(&users)
	//fmt.Println(users)

	//var user User
	//DB.Debug().Preload("Articles").Take(&user, 5)
	//DB.Debug().Model(&user).Association("Articles").Delete(&user.Articles)
	//DB.Delete(&user)

	////	级联删除
	//var user User
	//DB.Take(&user, 6)
	//DB.Debug().Select("Articles").Delete(&user)

}
