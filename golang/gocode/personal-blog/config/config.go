package config

// 数据库连接、JWT密钥、端口等配置

import (
	"os" // 导入os包以访问环境变量
)

// GetDatabaseURL 从环境变量 "DATABASE_URL" 获取数据库连接字符串。
// 如果环境变量未设置，则返回一个默认的本地MySQL连接字符串。
func GetDatabaseURL() string {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// 默认MySQL连接信息
		return "root:root@tcp(localhost:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return dbURL
}

// GetJWTSecret 从环境变量 "JWT_SECRET" 获取JWT签名密钥。
// 如果环境变量未设置，则返回一个不安全的默认密钥，并提示在生产环境中更改。
func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "your-secret-key-change-this-in-production"
	}
	return secret
}

// GetPort 从环境变量 "PORT" 获取服务监听的端口号。
// 如果环境变量未设置，则返回默认端口 "8080"。
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}
	return port
}

// GetLogPath 从环境变量 "LOG_PATH" 获取日志目录路径。
// 如果环境变量未设置，则返回默认路径 "E:/logs"。
func GetLogPath() string {
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		return "E:/logs"
	}
	return logPath
}
