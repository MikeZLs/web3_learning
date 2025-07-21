# 个人博客系统后端

 用到的技术栈：Go、Gin、GORM、MySQL。

## 功能特性

- **用户认证**: 支持用户注册、登录，并使用 JWT (JSON Web Token) 进行认证。
- **博客文章**: 实现了博客文章的增、删、改、查（CRUD）操作。
- **评论功能**: 用户可以对文章发表评论。
- **权限控制**: 只有文章的作者才能编辑或删除自己的文章。
- **错误处理**: 完善的错误处理机制，返回规范的 HTTP 状态码。
- **日志记录**: 使用 Zap 日志库，并将日志文件输出到 E 盘。

## 技术栈

- **Go 1.24+**
- **Gin 框架**: 高性能的 HTTP Web 框架。
- **GORM**: 用于数据库操作的 ORM 库。
- **MySQL**: 8.0.26 版本。
- **JWT**: 用于身份验证的 Token 机制。
- **Zap**: 结构化、分级别的日志库。
- **bcrypt**: 用于密码哈希加密。

## 安装与设置

1.  **克隆仓库**
    ```bash
    git clone <repository-url>
    cd personal-blog
    ```

2.  **安装依赖**
    ```bash
    go mod tidy
    ```

3.  **设置 MySQL 数据库**
   - 创建一个名为 `blog_db` 的 MySQL 数据库。
   - 如果需要，请在 `config/config.go` 文件中更新数据库连接信息。

4.  **环境变量** (可选)
    ```bash
    export DATABASE_URL="root:root@tcp(localhost:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
    export JWT_SECRET="your-secret-key"
    export PORT="8080"
    export LOG_LEVEL="INFO"  # 可选值: DEBUG, INFO, WARN, ERROR
    ```

5.  **创建日志目录**
   - 应用程序会自动创建 `E:/logs/` 目录。
   - 你也可以通过设置 `LOG_PATH` 环境变量来指定自定义路径。

6.  **运行应用**
    ```bash
    go run main.go
    ```

## 日志配置

本应用使用 Zap 日志库，并具有以下配置：

### 日志文件位置
- **主日志**: `E:/logs/blog-app.log` - 记录所有级别的应用日志。
- **错误日志**: `E:/logs/blog-error.log` - 仅记录错误级别的日志。
- **控制台输出**: 日志也会同时在终端中显示。

### 日志级别
通过 `LOG_LEVEL` 环境变量设置：
- `DEBUG`: 最详细的级别，包含所有日志。
- `INFO`: 一般信息（默认级别）。
- `WARN`: 警告信息。
- `ERROR`: 仅错误信息。

### 日志格式
日志采用结构化的 JSON 格式，包含以下字段：
- `timestamp`: ISO8601 格式的时间戳。
- `level`: 日志级别 (DEBUG, INFO, WARN, ERROR)。
- `logger`: 日志记录器名称。
- `caller`: 源代码文件及行号。
- `message`: 日志消息。
- `stacktrace`: 错误堆栈信息（仅在错误日志中出现）。

### 日志条目示例
```json
{
  "timestamp": "2024-01-15T10:30:45.123Z",
  "level": "INFO",
  "logger": "personal-blog",
  "caller": "controllers/auth.go:45",
  "message": "用户注册成功",
  "username": "john_doe"
}
```

## API 接口

### 认证接口
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录

### 文章接口
- `GET /api/posts` - 获取所有文章
- `GET /api/posts/:id` - 获取单篇文章
- `POST /api/posts` - 创建文章 (需要认证)
- `PUT /api/posts/:id` - 更新文章 (需要认证, 仅限作者)
- `DELETE /api/posts/:id` - 删除文章 (需要认证, 仅限作者)

### 评论接口
- `GET /api/posts/:post_id/comments` - 获取文章下的所有评论
- `POST /api/posts/:post_id/comments` - 发表评论 (需要认证)

### 健康检查
- `GET /health` - 服务健康状态

## 数据库结构

### Users (用户) 表
- `id` (主键)
- `username` (唯一)
- `email` (唯一)
- `password` (哈希值)
- `created_at`
- `updated_at`
- `deleted_at`

### Posts (文章) 表
- `id` (主键)
- `title`
- `content`
- `user_id` (外键, 关联 Users 表)
- `created_at`
- `updated_at`
- `deleted_at`

### Comments (评论) 表
- `id` (主键)
- `content`
- `user_id` (外键, 关联 Users 表)
- `post_id` (外键, 关联 Posts 表)
- `created_at`

## 项目结构

```
personal-blog/
├── main.go                 # 应用主入口
├── go.mod                  # Go 模块依赖
├── README.md              # 项目文档
│
├── config/                # 配置管理
│   └── config.go          # 数据库、JWT、端口等配置
│
├── models/                # 数据模型
│   ├── user.go           # 用户模型及相关结构体
│   ├── post.go           # 文章模型及相关结构体
│   └── comment.go        # 评论模型及相关结构体
│
├── database/              # 数据库相关
│   └── database.go       # 数据库初始化与迁移
│
├── middleware/            # 中间件
│   └── auth.go           # JWT 认证中间件
│
├── controllers/           # 控制器 (业务逻辑)
│   ├── auth.go           # 认证控制器
│   ├── post.go           # 文章管理控制器
│   └── comment.go        # 评论管理控制器
│
└── routes/                # 路由配置
    └── routes.go         # API 路由设置

日志文件将存储在:
E:/logs/
├── blog-app.log          # 主应用日志
└── blog-error.log        # 仅错误日志
```

## 认证机制

本 API 使用 JWT (JSON Web Tokens) 进行认证：

1.  通过注册或登录接口获取一个 JWT。
2.  在后续需要认证的请求中，将 JWT 放入 `Authorization` 请求头中，格式为：`Bearer <token>`。
3.  Token 的有效期为 24 小时。

## 错误处理

本 API 返回标准的 HTTP 状态码：
- `200` - 成功 (Success)
- `201` - 已创建 (Created)
- `400` - 错误请求 (Bad Request)
- `401` - 未授权 (Unauthorized)
- `403` - 禁止访问 (Forbidden)
- `404` - 未找到 (Not Found)
- `409` - 冲突 (Conflict)
- `500` - 服务器内部错误 (Internal Server Error)

## 开发

### 在开发模式下运行
```bash
# 将 Gin 设置为调试模式
export GIN_MODE=debug

# 将日志级别设置为调试
export LOG_LEVEL=DEBUG

go run main.go
```

### 测试 API
你可以使用 Postman、curl 或任何 HTTP 客户端工具来测试接口。

注册示例：
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
