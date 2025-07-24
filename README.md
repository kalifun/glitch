# Glitch

<div align="center">

![Glitch Logo](https://img.shields.io/badge/Glitch-Error%20Management-blue?style=for-the-badge)

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/kalifun/glitch?style=flat-square)](https://goreportcard.com/report/github.com/kalifun/glitch)
[![Release](https://img.shields.io/github/v/release/kalifun/glitch?style=flat-square)](https://github.com/kalifun/glitch/releases)

**A modern, YAML-driven error management tool for Go applications**

[English](README.md) | [中文](README_CN.md)

</div>

## 🚀 Overview

Glitch is a powerful error management tool that revolutionizes how you handle errors in Go applications. It provides a unified, YAML-driven approach to define, generate, and manage errors with multi-language support and flexible output formats.

### ✨ Key Features

- 🎯 **YAML-Driven Configuration**: Define all errors in a centralized YAML file
- 🌍 **Multi-Language Support**: Built-in internationalization with unlimited language support
- 🛠️ **Code Generation**: Automatically generate type-safe Go error definitions
- 🎨 **Flexible Output**: Support for JSON, text, and structured error formats
- 🔌 **Framework Integration**: Support custom processors, extensible to other frameworks
- 📊 **Rich Metadata**: Errors with categories, severity levels, and custom metadata
- 🔧 **Direct Usage**: Generated errors implement Go's standard `error` interface

## 📦 Installation

```bash
go install github.com/kalifun/glitch@latest
```

## 🚀 Quick Start

### 1. Initialize Template

```bash
glitch init
```

This creates an `errors.yaml` template file:

```yaml
error:
  - key: MissingParameterF
    code: MissingParameter
    message:
      cn: "缺少参数: %s"
      en: "Missing Parameter: %s"
```

### 2. Define Your Errors

Edit the `errors.yaml` file with your error definitions:

```yaml
error:
  - key: user_not_found
    code: USER_NOT_FOUND
    category: resource
    severity: error
    description: "User does not exist in the system"
    message:
      en: "User not found: %s"
      cn: "用户未找到: %s"
      fr: "Utilisateur non trouvé: %s"
      es: "Usuario no encontrado: %s"

  - key: invalid_email
    code: INVALID_EMAIL
    category: validation
    severity: error
    message:
      en: "Invalid email address: %s"
      cn: "无效的邮箱地址: %s"
```

### 3. Generate Code

```bash
# Generate
glitch gen -y errors.yaml -p errors
```

### 4. Use Generated Errors

The generated code provides direct, type-safe error usage:

```go
package main

import (
    "fmt"
    "your-project/errors"
)

func findUser(id string) error {
    if id == "" {
        // Direct usage - implements error interface
        return errors.UserNotFoundF.Args(id)
    }
    return nil
}

func main() {
    err := findUser("")
    fmt.Printf("Error: %v\n", err)
    // Output: User not found:
}
```

## 🌐 Framework Integration

### Gin Integration with Custom Handlers

```go
package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kalifun/glitch/repo/gerr"
)

// Create custom error handler for web responses
type WebErrorHandler struct{}

func (h *WebErrorHandler) CanHandle(err *gerr.Error) bool {
    return true // Handle all errors in web context
}

func (h *WebErrorHandler) Handle(ctx context.Context, err *gerr.Error) interface{} {
    // Extract HTTP status from error category
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
    // Setup custom engine for web
    webEngine := gerr.NewProcessorEngine()
    webEngine.AddHandler(&WebErrorHandler{})

    router := gin.Default()

    // Use custom engine in middleware
    router.Use(func(c *gin.Context) {
        c.Next()

        for _, ginErr := range c.Errors {
            if err, ok := ginErr.Err.(*gerr.Error); ok {
                // Add request context
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

## 🌍 Multi-Language Support

### Automatic Language Detection

```go
// Set language in context
ctx := context.WithValue(context.Background(), "language", "cn")
result := gerr.Process(ctx, errors.UserNotFoundF.Args("john"))

// Or use HTTP headers for web applications
// Accept-Language: cn,en;q=0.9
```

### Supported Languages

The system supports unlimited languages. Common examples:
- `en` - English
- `cn` - Chinese (Simplified)
- `fr` - French
- `es` - Spanish
- `de` - German
- `ja` - Japanese
- And any custom language code you define

## 📊 Error Categories and Severity

Organize your errors with categories and severity levels:

```yaml
error:
  - key: validation_failed
    code: VALIDATION_FAILED
    category: validation      # validation, auth, resource, system, business, network
    severity: error          # warning, error, critical
    message:
      en: "Validation failed for field: %s"
```

## Handler

### Custom Handlers

Create specialized handlers for different error types:

```go
// Business logic handler
type BusinessHandler struct{}

func (h *BusinessHandler) CanHandle(err *gerr.Error) bool {
    // Handle business category errors
    if def := err.GetErrWrapper(); def != nil {
        return def.Category == "business"
    }
    return false
}

func (h *BusinessHandler) Handle(ctx context.Context, err *gerr.Error) interface{} {
    return map[string]interface{}{
        "error": map[string]interface{}{
            "type":    "business_error",
            "code":    err.GetCode(),
            "message": err.Error(),
            "suggestion": "Please contact customer service",
        },
    }
}

// Register the handler
engine := gerr.NewProcessorEngine()
engine.AddHandler(&BusinessHandler{})
result := engine.Process(ctx, err)
```

### Middleware Chain

Add middleware for cross-cutting concerns:

```go
// Request tracking middleware
func RequestTrackingMiddleware(ctx context.Context, err *gerr.Error, next func(context.Context, *gerr.Error) interface{}) interface{} {
    // Add request context to error
    if requestID := ctx.Value("request_id"); requestID != nil {
        err = err.With("request_id", requestID)
    }

    // Continue processing
    return next(ctx, err)
}

// Performance monitoring middleware
func PerformanceMiddleware(ctx context.Context, err *gerr.Error, next func(context.Context, *gerr.Error) interface{}) interface{} {
    start := time.Now()
    result := next(ctx, err)
    duration := time.Since(start)

    log.Printf("Error processing took %v for error: %s", duration, err.GetKey())
    return result
}

// Setup engine with middleware
engine := gerr.NewProcessorEngine()
engine.AddMiddleware(RequestTrackingMiddleware)
engine.AddMiddleware(PerformanceMiddleware)
```

### Format Options

Control error output format:

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

## 🔧 Advanced usage

### Custom metadata

```go
err := errors.UserNotFoundF.
    Args("john_doe").
    With("user_id", "12345").
    With("timestamp", time.Now()).
    With("request_id", "req-123")
```

### Error Wrap

```go
originalErr := someFunction()
wrappedErr := errors.InternalError.Wrap(originalErr)
```