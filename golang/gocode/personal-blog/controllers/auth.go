package controllers

// 封装了所有与用户认证（注册、登录）相关的HTTP请求处理逻辑

import (
	"net/http"
	"personal-blog/config"
	"personal-blog/middleware"
	"personal-blog/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AuthController 是一个控制器结构体，它聚合了处理认证请求所需的所有依赖（数据库和日志）。
type AuthController struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewAuthController 是一个构造函数，用于创建一个新的 AuthController 实例。
// 这种模式（构造函数注入依赖）有利于代码的解耦和测试。
func NewAuthController(db *gorm.DB, logger *zap.Logger) *AuthController {
	return &AuthController{db: db, logger: logger}
}

// Register 处理用户注册的API请求 (POST /api/auth/register)。
func (ctrl *AuthController) Register(c *gin.Context) {
	var req models.UserRegisterRequest
	// 1. 绑定并验证请求体JSON到`req`结构体。如果失败，Gin会自动返回400 Bad Request。
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 检查用户名或邮箱是否已经被注册。
	var existingUser models.User
	if err := ctrl.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		// `err == nil` 意味着查询成功找到了记录，即用户已存在。
		c.JSON(http.StatusConflict, gin.H{"error": "用户名或邮箱已存在"})
		return
	}

	// 3. 创建一个新的User模型实例。
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // 此时还是明文密码
	}

	// 4. 对用户的明文密码进行哈希加密。
	if err := user.HashPassword(); err != nil {
		ctrl.logger.Error("密码哈希失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	// 5. 将哈希后的用户信息存入数据库。
	if err := ctrl.db.Create(&user).Error; err != nil {
		ctrl.logger.Error("创建用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	ctrl.logger.Info("用户注册成功", zap.String("username", user.Username))
	c.JSON(http.StatusCreated, gin.H{
		"message": "用户注册成功",
		"user":    user.ToResponse(), // 返回安全的、不含密码的用户信息。
	})
}

// Login 处理用户登录的API请求 (POST /api/auth/login)。
func (ctrl *AuthController) Login(c *gin.Context) {
	var req models.UserLoginRequest
	// 1. 绑定并验证请求体。
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 根据请求中的用户名从数据库查找用户。
	var user models.User
	if err := ctrl.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		// 为了安全，不应明确提示是“用户名不存在”还是“密码错误”，统一返回“无效凭证”。
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的凭证"})
		return
	}

	// 3. 验证用户输入的密码是否正确。
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的凭证"})
		return
	}

	// 4. 密码验证成功，为该用户生成JWT。
	token, err := ctrl.generateToken(user.ID, user.Username)
	if err != nil {
		ctrl.logger.Error("生成token失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	ctrl.logger.Info("用户登录成功", zap.String("username", user.Username))
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user":    user.ToResponse(),
	})
}

// generateToken 是一个内部辅助函数，用于生成JWT。
func (ctrl *AuthController) generateToken(userID uint, username string) (string, error) {
	// 创建JWT的Claims部分，包含我们想在token中携带的信息。
	claims := middleware.Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// IssuedAt (iat): token的签发时间。
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// ExpiresAt (exp): token的过期时间。这里设置为1小时后。
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	// 使用指定的签名方法（这里是HS256）和claims创建一个新的token对象。
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用我们配置的密钥对token进行签名，生成最终的token字符串。
	return token.SignedString([]byte(config.GetJWTSecret()))
}
