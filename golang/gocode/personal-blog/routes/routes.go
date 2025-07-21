package routes

// 该包负责定义和组织所有的API路由。

import (
	"personal-blog/controllers" // 导入控制器包
	"personal-blog/middleware"  // 导入中间件包

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SetupRoutes 配置并初始化所有API路由
func SetupRoutes(r *gin.Engine, db *gorm.DB, logger *zap.Logger) {
	// 初始化所有控制器实例
	authController := controllers.NewAuthController(db, logger)
	postController := controllers.NewPostController(db, logger)
	commentController := controllers.NewCommentController(db, logger)

	// --- 认证路由 ---
	// 创建 /api/auth 路由组
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authController.Register) // 注册路由
		auth.POST("/login", authController.Login)       // 登录路由
	}

	// --- 文章路由 ---
	// 创建 /api/posts 路由组
	posts := r.Group("/api/posts")
	{
		// 公开访问的路由
		posts.GET("", postController.GetPosts)    // 获取所有文章
		posts.GET("/:id", postController.GetPost) // 获取单篇文章

		// --- 受保护的路由 ---
		// 对该路由组下的后续路由应用AuthMiddleware中间件
		posts.Use(middleware.AuthMiddleware(db, logger))
		posts.POST("", postController.CreatePost)       // 创建文章
		posts.PUT("/:id", postController.UpdatePost)    // 更新文章
		posts.DELETE("/:id", postController.DeletePost) // 删除文章
	}

	// --- 评论路由 ---
	// 创建 /api/posts/:post_id/comments 路由组
	comments := r.Group("/api")
	{
		// 公开访问的路由
		comments.GET("/:post_id/comments", commentController.GetComments) // 获取文章下的所有评论

		// --- 受保护的路由 ---
		// 应用AuthMiddleware中间件
		comments.Use(middleware.AuthMiddleware(db, logger))
		comments.POST("/:post_id/comments", commentController.CreateComment) // 创建评论
	}

	// --- 健康检查路由 ---
	// 用于检查服务是否正常运行
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
