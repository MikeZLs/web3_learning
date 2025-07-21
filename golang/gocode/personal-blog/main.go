// main.go
package main

import (
	"fmt"
	"log"
	"os"
	"personal-blog/config"   // 导入配置包，用于获取环境变量和默认配置
	"personal-blog/database" // 导入数据库包，用于初始化和迁移
	"personal-blog/routes"   // 导入路由包，用于设置API路由

	"github.com/gin-gonic/gin" // 导入Gin框架
	"go.uber.org/zap"          // 导入Zap日志库，用于高性能的结构化日志记录
	"go.uber.org/zap/zapcore"
)

// main 是应用程序的入口函数
func main() {
	// 初始化Zap日志记录器，并配置日志输出到文件
	logger, err := initLogger()
	if err != nil {
		log.Fatal("初始化日志记录器失败:", err)
	}
	defer logger.Sync() // 确保在程序退出时，所有缓冲的日志都被写入

	// 初始化数据库连接
	db, err := database.InitDatabase()
	if err != nil {
		logger.Fatal("连接数据库失败", zap.Error(err))
	}

	// 自动迁移数据库表结构
	err = database.AutoMigrate(db)
	if err != nil {
		logger.Fatal("迁移数据库失败", zap.Error(err))
	}

	// 初始化Gin Web服务器，使用默认中间件(logger, recovery)
	r := gin.Default()

	// 设置所有API路由
	routes.SetupRoutes(r, db, logger)

	// 从配置中获取端口号，并启动服务器
	port := config.GetPort()
	logger.Info("服务启动中", zap.String("port", port))
	if err := r.Run(":" + port); err != nil { // 监听并启动服务
		logger.Fatal("服务启动失败", zap.Error(err))
	}
}

// initLogger 初始化 zap 日志记录器
// 配置日志同时输出到控制台和E盘的文件
func initLogger() (*zap.Logger, error) {
	// 定义日志存储目录
	logDir := "E:/logs"
	// 如果目录不存在，则创建它
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 使用Zap的生产环境预设配置
	cfg := zap.NewProductionConfig()

	// 设置日志输出路径，包括主日志文件和标准输出（控制台）
	cfg.OutputPaths = []string{
		"E:/logs/blog-app.log",
		"stdout",
	}
	// 设置错误日志输出路径，包括错误日志文件和标准错误输出（控制台）
	cfg.ErrorOutputPaths = []string{
		"E:/logs/blog-error.log",
		"stderr",
	}

	// 通过环境变量 "LOG_LEVEL" 设置日志级别，默认为INFO
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "DEBUG":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "WARN":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "ERROR":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// 配置日志编码为JSON格式
	cfg.Encoding = "json"

	// 配置JSON编码器的字段名
	cfg.EncoderConfig.TimeKey = "timestamp"                   // 时间戳字段名
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 时间格式为ISO8601
	cfg.EncoderConfig.LevelKey = "level"                      // 日志级别字段名
	cfg.EncoderConfig.NameKey = "logger"                      // 日志记录器名字段名
	cfg.EncoderConfig.CallerKey = "caller"                    // 调用位置字段名 (文件:行号)
	cfg.EncoderConfig.MessageKey = "message"                  // 日志消息字段名
	cfg.EncoderConfig.StacktraceKey = "stacktrace"            // 堆栈跟踪字段名

	// 构建并返回日志记录器实例
	return cfg.Build()
}
