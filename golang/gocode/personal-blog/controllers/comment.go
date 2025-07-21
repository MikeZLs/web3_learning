package controllers

// 封装了所有与评论（Comment）资源相关的操作。

import (
	"net/http"
	"personal-blog/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CommentController 结构体
type CommentController struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewCommentController 构造函数
func NewCommentController(db *gorm.DB, logger *zap.Logger) *CommentController {
	return &CommentController{db: db, logger: logger}
}

// CreateComment 处理在某篇文章下创建新评论的请求 (POST /api/posts/:post_id/comments)。受保护路由。
func (ctrl *CommentController) CreateComment(c *gin.Context) {
	// 从URL中获取要评论的文章ID
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var req models.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 在创建评论前，先检查目标文章是否存在。
	var post models.Post
	if err := ctrl.db.First(&post, uint(postID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
			return
		}
		ctrl.logger.Error("获取文章失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	comment := models.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  uint(postID),
	}

	if err := ctrl.db.Create(&comment).Error; err != nil {
		ctrl.logger.Error("创建评论失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return
	}

	// 预加载评论的作者信息以便在响应中返回。
	ctrl.db.Preload("User").First(&comment, comment.ID)

	ctrl.logger.Info("评论创建成功", zap.Uint("comment_id", comment.ID))
	c.JSON(http.StatusCreated, gin.H{
		"message": "评论创建成功",
		"comment": comment.ToResponse(),
	})
}

// GetComments 处理获取某篇文章下所有评论的请求 (GET /api/posts/:post_id/comments)。公开路由。
func (ctrl *CommentController) GetComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 同样，先检查文章是否存在。
	var post models.Post
	if err := ctrl.db.First(&post, uint(postID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
			return
		}
		ctrl.logger.Error("获取文章失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	var comments []models.Comment
	// 查询所有`post_id`匹配的评论，并预加载每条评论的作者信息。
	if err := ctrl.db.Preload("User").Where("post_id = ?", uint(postID)).Find(&comments).Error; err != nil {
		ctrl.logger.Error("获取评论失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return
	}

	// 将评论列表转换为安全的响应DTO列表。
	var responses []models.CommentResponse
	for _, comment := range comments {
		responses = append(responses, comment.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"comments": responses})
}
