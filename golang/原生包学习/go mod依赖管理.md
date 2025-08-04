Go modules (go mod) 是 Go 语言的依赖管理系统，从 Go 1.11 开始引入，现在是 Go 项目的标准方式。让我详细介绍如何使用：

## 基本概念

Go modules 是一个包含 Go 包的集合，用于版本控制和分发。每个模块都有一个 `go.mod` 文件来定义模块路径和依赖关系。

## 创建新模块

```bash
# 创建新目录
mkdir myproject
cd myproject

# 初始化模块
go mod init example.com/myproject
```

这会创建一个 `go.mod` 文件：
```
module example.com/myproject

go 1.24
```

## 主要命令

### 1. `go mod init`
初始化新模块：
```bash
go mod init [module-path]
```

### 2. `go mod tidy`
整理依赖，添加缺失的模块，移除不需要的模块：
```bash
go mod tidy
```

### 3. `go mod download`
下载依赖到本地缓存：
```bash
go mod download
```

### 4. `go mod vendor`
将依赖复制到 vendor 目录：
```bash
go mod vendor
```

### 5. `go mod graph`
打印模块依赖关系图：
```bash
go mod graph
```

### 6. `go mod why`
解释为什么需要某个包或模块：
```bash
go mod why golang.org/x/text
```

## 添加依赖

当你在代码中导入新包时，Go 会自动下载：

```go
package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Hello World"})
    })
    r.Run()
}
```

运行 `go run main.go` 或 `go mod tidy` 会自动添加依赖到 `go.mod`：

```
module example.com/myproject

go 1.24

require github.com/gin-gonic/gin v1.9.1
```

## 版本管理

### 升级依赖
```bash
# 升级所有依赖到最新版本
go get -u

# 升级特定模块
go get -u github.com/gin-gonic/gin

# 升级到特定版本
go get github.com/gin-gonic/gin@v1.9.0
```

### 降级依赖
```bash
go get github.com/gin-gonic/gin@v1.8.0
```

## go.sum 文件

`go.sum` 文件包含依赖的校验和，确保依赖的完整性。不要手动编辑这个文件。

## 工作区 (Workspaces) - Go 1.18+

对于多模块开发，可以使用工作区：

```bash
# 创建工作区
go work init ./module1 ./module2

# 添加模块到工作区
go work use ./module3
```

## 实际开发流程

1. **创建项目**：
```bash
mkdir myapp && cd myapp
go mod init github.com/username/myapp
```

2. **编写代码并添加依赖**：
```go
// main.go
package main

import (
    "github.com/gorilla/mux"
    "net/http"
)
```

3. **整理依赖**：
```bash
go mod tidy
```

4. **构建和运行**：
```bash
go build
go run main.go
```

## 常见配置

### 设置 GOPROXY
```bash
# 使用国内代理
go env -w GOPROXY=https://goproxy.cn,direct
```

### 私有模块
```bash
go env -w GOPRIVATE=github.com/yourcompany/*
```