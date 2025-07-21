package models

// 定义了评论相关的GORM模型和请求/响应数据结构。

import (
	"time"
)

// Comment 定义了评论的数据模型，对应数据库中的 `comments` 表。
type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`       // 评论作者ID
	User      User      `json:"user" gorm:"foreignKey:UserID"` // 关联评论作者
	PostID    uint      `json:"post_id" gorm:"not null"`       // 评论所属文章ID
	Post      Post      `json:"post" gorm:"foreignKey:PostID"` // 关联所属文章 (这个关联在当前逻辑中较少用到，但定义出来是好习惯)
	CreatedAt time.Time `json:"created_at"`
}

// CommentCreateRequest 定义了创建评论API的请求体结构。
type CommentCreateRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

// CommentResponse 定义了返回给客户端的评论信息DTO。
type CommentResponse struct {
	ID        uint         `json:"id"`
	Content   string       `json:"content"`
	User      UserResponse `json:"user"` // 嵌套评论作者的安全信息
	CreatedAt time.Time    `json:"created_at"`
}

// ToResponse 将一个Comment模型对象转换为安全的CommentResponse DTO。
func (c *Comment) ToResponse() CommentResponse {
	return CommentResponse{
		ID:        c.ID,
		Content:   c.Content,
		User:      c.User.ToResponse(),
		CreatedAt: c.CreatedAt,
	}
}
