package controllers

// 封装了所有与文章（Post）资源相关的CRUD操作

import (
	"net/http"
	"personal-blog/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PostController 结构体
type PostController struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewPostController 构造函数
func NewPostController(db *gorm.DB, logger *zap.Logger) *PostController {
	return &PostController{db: db, logger: logger}
}

// CreatePost 处理创建新文章的请求 (POST /api/posts)。这是一个受保护的路由。
func (ctrl *PostController) CreatePost(c *gin.Context) {
	var req models.PostCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从Gin上下文中获取由AuthMiddleware设置的当前登录用户的ID。
	userID, exists := c.Get("user_id")
	if !exists {
		// 理论上，如果AuthMiddleware配置正确，这个情况不会发生。但作为防御性编程，检查一下总是好的。
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint), // `c.Get`返回的是interface{}, 需要进行类型断言。
	}

	if err := ctrl.db.Create(&post).Error; err != nil {
		ctrl.logger.Error("创建文章失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	// 在返回响应前，使用Preload("User")来加载新创建文章的作者信息。
	// 这样返回的JSON就会包含完整的作者信息，而不仅仅是user_id。
	ctrl.db.Preload("User").First(&post, post.ID)

	ctrl.logger.Info("文章创建成功", zap.Uint("post_id", post.ID))
	c.JSON(http.StatusCreated, gin.H{
		"message": "文章创建成功",
		"post":    post.ToResponse(),
	})
}

// GetPosts 处理获取所有文章列表的请求 (GET /api/posts)。这是一个公开路由。
func (ctrl *PostController) GetPosts(c *gin.Context) {
	var posts []models.Post
	// Preload("User")告诉GORM在查询文章的同时，把每篇文章关联的User信息也一并查询出来。
	// 这避免了N+1查询问题，提高了效率。
	if err := ctrl.db.Preload("User").Find(&posts).Error; err != nil {
		ctrl.logger.Error("获取文章列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	// 将查询到的Post模型列表，逐个转换为安全的PostResponse DTO列表。
	var responses []models.PostResponse
	for _, post := range posts {
		responses = append(responses, post.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"posts": responses})
}

// GetPost 处理获取单篇文章详情的请求 (GET /api/posts/:id)。这是一个公开路由。
func (ctrl *PostController) GetPost(c *gin.Context) {
	// strconv.ParseUint() 第一个参数是gin框架从URL路径参数中获取文章id，第二个参数表示按十进制解析，第三个参数表示目标大小不超过 32 位
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var post models.Post
	// 根据ID查找文章，并预加载作者信息。
	if err := ctrl.db.Preload("User").First(&post, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果GORM返回ErrRecordNotFound，说明数据库中没有这条记录。
			c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
			return
		}
		ctrl.logger.Error("获取文章失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post.ToResponse()})
}

// UpdatePost 处理更新文章的请求 (PUT /api/posts/:id)。这是一个受保护的路由。
func (ctrl *PostController) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var req models.PostUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var post models.Post
	// 首先，需要从数据库中把要更新的文章找出来。
	if err := ctrl.db.First(&post, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
			return
		}
		ctrl.logger.Error("获取文章失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	// **核心权限检查**: 检查当前登录用户的ID是否与文章的作者ID(post.UserID)匹配。
	if post.UserID != userID.(uint) {
		// 如果不匹配，说明该用户试图修改不属于他的文章，应拒绝该请求。
		c.JSON(http.StatusForbidden, gin.H{"error": "您只能更新自己的文章"})
		return
	}

	// 权限检查通过，更新文章的字段。
	post.Title = req.Title
	post.Content = req.Content

	// `db.Save()`会更新记录的所有字段。
	if err := ctrl.db.Save(&post).Error; err != nil {
		ctrl.logger.Error("更新文章失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	// 重新加载作者信息以返回最新的完整数据。
	ctrl.db.Preload("User").First(&post, post.ID)

	ctrl.logger.Info("文章更新成功", zap.Uint("post_id", post.ID))
	c.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功",
		"post":    post.ToResponse(),
	})
}

// DeletePost 处理删除文章的请求 (DELETE /api/posts/:id)。这是一个受保护的路由。
func (ctrl *PostController) DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var post models.Post
	if err := ctrl.db.First(&post, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
			return
		}
		ctrl.logger.Error("获取文章失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	// **核心权限检查**: 同样需要检查当前用户是否是文章作者。
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "您只能删除自己的文章"})
		return
	}

	// `db.Delete()`会删除指定的记录。模型中定义了`gorm.DeletedAt`字段，GORM会执行软删除，否则会执行物理删除。
	// 在删除文章前，应该先删除其所有关联的评论，以维护数据完整性，但这部分逻辑未在此实现。
	if err := ctrl.db.Delete(&post).Error; err != nil {
		ctrl.logger.Error("删除文章失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}

	ctrl.logger.Info("文章删除成功", zap.Uint("post_id", post.ID))
	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}
