package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	Id   uint   `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

type User2 struct {
	Sex    string `json:"sex"`
	Height uint   `json:"height"`
	Email  string `json:"email" binding:"required,email"`
	//  binding  告诉框架如何验证将传入的 JSON 或表单数据绑定到 Go 结构体字段时的数据,标签里面可以包含一个或多个验证规则，用空格 分隔
	//	required 强制字段存在：它表示对应的字段在传入的请求数据中必须存在。如果这个字段缺失，验证就会失败
}

// Middleware1 自定义中间件（常用于认证）
func Middleware1() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("mw1 before")
		c.Next()
		fmt.Println("mw1 after")
	}
}

func Middleware2() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("mw2 before")
		c.Next()
		fmt.Println("mw2 after")
	}
}

func main() {
	r := gin.Default()

	// 全局调用中间件
	//r.Use(Middleware1())

	// 接口单独调用中间件
	//r.GET("/user1", Middleware1(), func(c *gin.Context) {})
	//r.GET("/user2", Middleware1(), Middleware2(), func(c *gin.Context) {
	//	fmt.Println("self")
	//}) // 可调用多个中间件

	//grp1 := r.Group("/v1")
	//// 分组调用中间件
	//grp1.Use(Middleware1())

	//r.LoadHTMLGlob("./gin/demo01/templates/*")
	//
	//r.GET("/user", func(c *gin.Context) {
	//	c.HTML(200, "user.impl", gin.H{
	//		"title": "HTML测试",
	//	})
	//})

	//r.LoadHTMLGlob("./gin/demo01/templates/**/*")

	//r.GET("/user", func(c *gin.Context) {
	//	c.HTML(200, "a.impl", gin.H{
	//		"title": "HTML测试",
	//	})
	//})
	//
	//r.GET("/top", func(c *gin.Context) {
	//	c.HTML(200, "a.impl", gin.H{
	//		"title": "HTML测试111",
	//	})
	//})

	//r.GET("/user", func(c *gin.Context) {
	//	//c.String(200, "Hello World")
	//	user := User{
	//		Username: "admin",
	//		Password: "123456",
	//	}
	//	//c.JSON(200, user)
	//
	//	c.XML(200, user)
	//
	//	//c.JSON(200, gin.H{
	//	//	"message": "hello world!",
	//	//	"code":    200,
	//	//})
	//
	//})

	//r.Any("/login", func(c *gin.Context) {
	//	c.String(200, "Hello World")
	//})
	//
	//grp1 := r.Group("/v1")
	//{
	//	grp1.GET("/user", func(c *gin.Context) {
	//		c.String(200, "v1 Hello World")
	//	})
	//}
	//
	//grp2 := r.Group("/v2")
	//{
	//	grp2.GET("/user", func(c *gin.Context) {
	//		c.String(200, "v2 Hello World")
	//	})
	//}

	//r.GET("/user", func(c *gin.Context) {})
	//r.POST("/user", func(c *gin.Context) {})
	//r.PUT("/user", func(c *gin.Context) {})
	//r.DELETE("/user", func(c *gin.Context) {})
	//r.PATCH("/user", func(c *gin.Context) {})
	//
	//// 重定向
	//r.GET("/member", func(c *gin.Context) {
	//	c.Redirect(http.StatusMovedPermanently, "/user")
	//	c.Redirect(http.StatusFound, "www.baidu.com")
	//})

	//dir, _ := os.Getwd()
	//fmt.Println("当前目录：", dir)

	//// 访问静态文件
	//r.Static("/static", "./gin/demo01/static") // 访问目录下的所有文件
	////r.StaticFS("/static", http.Dir("./gin/demo01/static")) // 与上边等价
	//r.StaticFile("/f1", "./gin/demo01/static/1.txt") // 访问单个文件

	//// 传递参数
	//r.GET("/user/:id/:name", func(c *gin.Context) {
	//	id := c.Param("id")
	//	name := c.Param("name")
	//	c.JSON(200, gin.H{
	//		"id":   id,
	//		"name": name,
	//	})
	//})
	//
	//r.GET("/user", func(c *gin.Context) {
	//	id := c.Query("id")
	//	name := c.Query("name")
	//	c.JSON(200, gin.H{
	//		"id":   id,
	//		"name": name,
	//	})
	//})
	//
	//r.POST("/userForm", func(c *gin.Context) {
	//	id := c.PostForm("id")
	//	name := c.PostForm("name")
	//	c.JSON(200, gin.H{
	//		"id":   id,
	//		"name": name,
	//	})
	//})
	//
	//r.POST("/userJson", func(c *gin.Context) {
	//	var user User
	//	//c.MustBindWith(&user, binding.JSON) // 绑定必需参数
	//	//c.ShouldBindBodyWith(&user, binding.JSON) // 绑定多个结构体参数
	//	err := c.ShouldBind(&user)
	//	if err != nil {
	//		c.JSON(400, gin.H{
	//			"err": err.Error(),
	//		})
	//		return
	//	}
	//	c.JSON(200, user)
	//})

	// gin自带用户密码认证
	user := r.Group("/user", gin.BasicAuth(gin.Accounts{
		"user1": "123456",
	}))

	user.GET("/info", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		c.String(http.StatusOK, user)
	})

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}

}
