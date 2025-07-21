package models

// 定义了与用户相关的GORM模型、请求/响应数据结构以及相关方法

import (
	"time"

	"golang.org/x/crypto/bcrypt" // 用于密码的哈希和验证
)

// User 定义了用户数据模型，它精确地映射到数据库中的 `users` 表。
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`            // 用户ID，`gorm:"primaryKey"`指定其为主键
	Username  string    `json:"username" gorm:"unique;not null"` // 用户名，`gorm:"unique;not null"`指定其在数据库中必须唯一且不能为空
	Email     string    `json:"email" gorm:"unique;not null"`    // 电子邮箱，同样唯一且不能为空
	Password  string    `json:"-" gorm:"not null"`               // 用户密码，`json:"-"`标签意味着此字段在序列化为JSON时将被忽略，这是为了防止密码通过API响应意外泄露。
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Posts     []Post    `json:"posts,omitempty" gorm:"foreignKey:UserID"`    // 定义一对多关系：一个用户可以有多篇文章。`omitempty`表示若该字段为空，则在JSON中省略。`gorm:"foreignKey:UserID"`指定了关联的外键。
	Comments  []Comment `json:"comments,omitempty" gorm:"foreignKey:UserID"` // 定义一对多关系：一个用户可以有多条评论。
}

// UserRegisterRequest 定义了用户注册API的请求体结构。
// 使用`binding`标签来让Gin框架自动进行数据验证。
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"` // 用户名，必填，最小长度3，最大长度20
	Email    string `json:"email" binding:"required,email"`           // 电子邮箱，必填，且必须符合email格式
	Password string `json:"password" binding:"required,min=6"`        // 密码，必填，最小长度6
}

// UserLoginRequest 定义了用户登录API的请求体结构。
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 定义了当API成功响应时，返回给客户端的用户信息结构。
// 这是一个数据传输对象（DTO），它的作用是只暴露必要且安全的信息，隐藏如密码等敏感字段。
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// HashPassword 是User模型的一个方法，用于对用户的明文密码进行哈希加密。
// 它应该在创建新用户或用户更新密码时被调用。
func (u *User) HashPassword() error {
	// bcrypt.GenerateFromPassword可以安全地生成密码的哈希值。
	// bcrypt.DefaultCost是推荐的哈希计算成本，成本越高越安全，但计算时间也越长。
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 用于在用户登录时，验证用户输入的密码是否与数据库中存储的哈希密码匹配。
func (u *User) CheckPassword(password string) bool {
	// 直接比较字符串是极不安全的。
	// 必须使用bcrypt.CompareHashAndPassword进行比较，它能有效防止时序攻击。
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil // 如果err为nil，说明密码验证成功。
}

// ToResponse 是一个转换方法，将一个完整的、可能包含敏感信息的User模型对象，转换为一个安全的、只用于API响应的UserResponse对象。
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}
