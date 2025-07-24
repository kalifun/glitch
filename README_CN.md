# Glitch

<div align="center">

![Glitch Logo](https://img.shields.io/badge/Glitch-错误管理-blue?style=for-the-badge)

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/kalifun/glitch?style=flat-square)](https://goreportcard.com/report/github.com/kalifun/glitch)
[![Release](https://img.shields.io/github/v/release/kalifun/glitch?style=flat-square)](https://github.com/kalifun/glitch/releases)

**现代化的 YAML 驱动 Go 应用错误管理工具**

[English](README.md) | [中文](README_CN.md)

</div>

## 🚀 概述

Glitch 是一个强大的错误管理工具，彻底改变了 Go 应用程序中错误处理的方式。它提供了统一的、YAML 驱动的方法来定义、生成和管理错误，支持多语言和灵活的输出格式。

### ✨ 核心特性

- 🎯 **YAML 驱动配置**: 在集中的 YAML 文件中定义所有错误
- 🌍 **多语言支持**: 内置国际化，支持无限制语言
- 🛠️ **代码生成**: 自动生成类型安全的 Go 错误定义
- 🎨 **灵活输出**: 支持 JSON、文本和结构化错误格式
- 🔌 **框架集成**: 支持自定义处理器，可扩展到其他框架
- 📊 **丰富元数据**: 带有分类、严重级别和自定义元数据的错误
- 🔧 **直接使用**: 生成的错误实现 Go 标准 `error` 接口



## 📦 安装

```bash
go install github.com/kalifun/glitch@latest
```

## 🚀 快速开始

### 1. 初始化模板

```bash
glitch init
```

这会创建一个 `errors.yaml` 模板文件：

```yaml
error:
  - key: MissingParameterF
    code: MissingParameter
    message:
      cn: "缺少参数: %s"
      en: "Missing Parameter: %s"
```

### 2. 定义错误

编辑 `errors.yaml` 文件，定义你的错误：

```yaml
error:
  - key: user_not_found
    code: USER_NOT_FOUND
    category: resource
    severity: error
    description: "系统中不存在该用户"
    message:
      cn: "用户未找到: %s"
      en: "User not found: %s"
      fr: "Utilisateur non trouvé: %s"
      es: "Usuario no encontrado: %s"
  
  - key: invalid_email
    code: INVALID_EMAIL
    category: validation
    severity: error
    message:
      cn: "无效的邮箱地址: %s"
      en: "Invalid email address: %s"
```

### 3. 生成代码

```bash
# 使用生成
glitch gen -y errors.yaml -p errors
```

### 4. 使用生成的错误

生成的代码提供直接的、类型安全的错误使用：

```go
package main

import (
    "fmt"
    "your-project/errors"
)

func findUser(id string) error {
    if id == "" {
        // 直接使用 - 实现了 error 接口
        return errors.UserNotFoundF.Args(id)
    }
    return nil
}

func main() {
    err := findUser("")
    fmt.Printf("错误: %v\n", err)
    // 输出: 用户未找到: 
}
```

## 🌐 框架集成

### Gin 集成与自定义处理器

```go
package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kalifun/glitch/repo/gerr"
)

// 为 Web 响应创建自定义错误处理器
type WebErrorHandler struct{}

func (h *WebErrorHandler) CanHandle(err *gerr.Error) bool {
	return true // 在 Web 上下文中处理所有错误
}

func (h *WebErrorHandler) Handle(ctx context.Context, err *gerr.Error) interface{} {
	// 从错误类别提取 HTTP 状态码
	statusCode := http.StatusInternalServerError
	if err.GetErrWrapper() != nil {
		switch err.GetErrWrapper().Category {
		case "validation":
			statusCode = http.StatusBadRequest
		case "auth":
			statusCode = http.StatusUnauthorized
		case "resource":
			statusCode = http.StatusNotFound
		}
	}

	return gin.H{
		"success": false,
		"error": gin.H{
			"code":    err.GetCode(),
			"message": err.Error(),
			"type":    err.GetErrWrapper().Category,
		},
		"status_code": statusCode,
	}
}

func main() {
	// 为 Web 设置自定义引擎
	webEngine := gerr.NewProcessorEngine()
	webEngine.AddHandler(&WebErrorHandler{})

	router := gin.Default()

	// 在中间件中使用自定义引擎
	router.Use(func(c *gin.Context) {
		c.Next()

		for _, ginErr := range c.Errors {
			if err, ok := ginErr.Err.(*gerr.Error); ok {
				// 添加请求上下文
				ctx := context.WithValue(context.Background(), "request_id", c.GetHeader("X-Request-ID"))
				ctx = context.WithValue(ctx, "language", c.GetHeader("Accept-Language"))

				result := webEngine.Process(ctx, err)
				c.JSON(http.StatusOK, result)
				return
			}
		}
	})

	router.Run(":8080")
}
```

## 🌍 多语言支持

### 自动语言检测

```go
// 在上下文中设置语言
ctx := context.WithValue(context.Background(), "language", "cn")
result := gerr.Process(ctx, errors.UserNotFoundF.Args("john"))

// 或在 Web 应用中使用 HTTP 头
// Accept-Language: cn,en;q=0.9
```

### 支持的语言

系统支持无限制语言。常见示例：
- `cn` - 中文（简体）
- `en` - 英语
- `fr` - 法语
- `es` - 西班牙语
- `de` - 德语
- `ja` - 日语
- 以及你定义的任何自定义语言代码

## 📊 错误分类和严重级别

使用分类和严重级别组织你的错误：

```yaml
error:
  - key: validation_failed
    code: VALIDATION_FAILED
    category: validation      # validation, auth, resource, system, business, network
    severity: error          # warning, error, critical
    message:
      cn: "字段验证失败: %s"
```

## Handler

### 自定义处理器

为不同错误类型创建专门的处理器：

```go
// 业务逻辑处理器
type BusinessHandler struct{}

func (h *BusinessHandler) CanHandle(err *gerr.Error) bool {
	// 处理业务类别的错误
	if def := err.GetErrWrapper(); def != nil {
		return def.Category == "business"
	}
	return false
}

func (h *BusinessHandler) Handle(ctx context.Context, err *gerr.Error) interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"type":       "general_error",
			"code":       err.GetCode(),
			"message":    err.Error(),
			"suggestion": "请联系客服",
		},
	}
}

// 注册处理器
engine := gerr.NewProcessorEngine()
engine.AddHandler(&BusinessHandler{})
result := engine.Process(ctx, err)
```

### 中间件链

为横切关注点添加中间件：

```go
// 请求追踪中间件
func RequestTrackingMiddleware(ctx context.Context, err *gerr.Error, next func(context.Context, *gerr.Error) interface{}) interface{} {
	// 向错误添加请求上下文
	if requestID := ctx.Value("request_id"); requestID != nil {
		err = err.With("request_id", requestID)
	}

	// 继续处理
	return next(ctx, err)
}

// 性能监控中间件
func PerformanceMiddleware(ctx context.Context, err *gerr.Error, next func(context.Context, *gerr.Error) interface{}) interface{} {
	start := time.Now()
	result := next(ctx, err)
	duration := time.Since(start)

	log.Printf("错误处理耗时 %v，错误: %s", duration, err.GetKey())
	return result
}

// 设置带中间件的引擎
engine := gerr.NewProcessorEngine()
engine.AddMiddleware(RequestTrackingMiddleware)
engine.AddMiddleware(PerformanceMiddleware)
```

### 格式化选项

控制错误输出格式：

```go
options := gerr.FormatOptions{
    IncludeMetadata: true,
    IncludeCause:    true,
    Language:        "cn",
    Format:          gerr.FormatTypeJSON,
}

formatter := gerr.NewDefaultFormatter()
result := formatter.FormatWithOptions(ctx, err, options)
```

## 🔧 高级用法

### 自定义元数据

```go
err := errors.UserNotFoundF.
    Args("john_doe").
    With("user_id", "12345").
    With("timestamp", time.Now()).
    With("request_id", "req-123")
```

### 错误包装

```go
originalErr := someFunction()
wrappedErr := errors.InternalError.Wrap(originalErr)
```