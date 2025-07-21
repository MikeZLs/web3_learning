package middleware

// JWT 认证中间件

import (
	"net/http"
	"personal-blog/config"
	"personal-blog/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Claims 是JWT中存储的自定义数据部分的结构体。
// 嵌入jwt.RegisteredClaims，这样就可以使用像ExpiresAt这样的标准JWT字段。
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthMiddleware 是一个Gin中间件，用于验证JWT并保护需要登录才能访问的路由。
func AuthMiddleware(db *gorm.DB, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 步骤1: 从请求头中获取 "Authorization" 字段的值。
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "需要提供Authorization请求头"})
			c.Abort() // Abort会阻止请求链中后续的处理函数被调用。
			return
		}

		// 步骤2: "Authorization" 头的值通常是 "Bearer <token>" 的格式。我们需要提取出token部分。
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // 如果移除前缀后字符串没变，说明格式不正确。
			c.JSON(http.StatusUnauthorized, gin.H{"error": "需要提供Bearer Token"})
			c.Abort()
			return
		}

		// 步骤3: 解析并验证JWT。
		claims := &Claims{}
		// jwt.ParseWithClaims会解析token字符串，并将解析出的claims填充到我们提供的claims变量中。
		// 它还需要一个回调函数来提供用于验证签名的密钥。
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetJWTSecret()), nil
		})

		// 如果解析过程中发生错误（如格式错误、签名不匹配）或者token本身无效（如已过期）。
		if err != nil || !token.Valid {
			logger.Error("无效的token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 步骤4: [推荐] 验证token中的用户是否仍然真实存在于数据库中。
		// 这是一个额外的安全层，可以防止已删除的用户使用其尚未过期的旧token继续访问系统。
		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			logger.Error("用户未找到", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未找到"})
			c.Abort()
			return
		}

		// 步骤5: 认证成功！将从token中解析出的用户信息存入Gin的上下文中。
		// 这样，后续的处理函数（即真正的业务逻辑控制器）就可以通过 `c.Get("user_id")` 来获取当前登录用户的信息。
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		// `c.Next()` 调用请求处理链中的下一个处理函数。
		c.Next()
	}
}
