package models

// 定义了文章相关的GORM模型和请求/响应数据结构。

import (
	"time"
)

// Post 定义了博客文章的数据模型，对应数据库中的 `posts` 表。
type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`           // `gorm:"type:text"`确保数据库中使用适合存储长文本的类型
	UserID    uint      `json:"user_id" gorm:"not null"`                     // 外键，指向users表的id
	User      User      `json:"user" gorm:"foreignKey:UserID"`               // 定义多对一关系：一篇文章属于一个用户。GORM会用它来做关联查询（Preload）。
	Comments  []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"` // 定义一对多关系：一篇文章可以有多条评论。
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// PostCreateRequest 定义了创建文章API的请求体结构。
type PostCreateRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

// PostUpdateRequest 定义了更新文章API的请求体结构。
type PostUpdateRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

// PostResponse 定义了返回给客户端的文章信息DTO。
type PostResponse struct {
	ID        uint         `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	User      UserResponse `json:"user"` // 嵌套了安全的UserResponse，而不是完整的User模型，避免泄露作者的敏感信息。
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// ToResponse 将一个Post模型对象转换为安全的PostResponse DTO。
func (p *Post) ToResponse() PostResponse {
	return PostResponse{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		User:      p.User.ToResponse(), // 调用关联的User对象的ToResponse方法
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
