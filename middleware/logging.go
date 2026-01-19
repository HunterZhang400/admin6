package middleware

import (
	"admin6/model"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware 创建日志中间件（在认证中间件之后使用）
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()

		// 获取客户端IP
		clientIP := c.ClientIP()

		// 获取用户ID（如果存在）
		userID := "N/A"
		if user, exists := c.Get("user"); exists {
			if userGorm, ok := user.(*model.UserGorm); ok {
				userID = fmt.Sprintf("%d", userGorm.ID)
			}
		}

		// 记录请求信息
		log.Printf("[REQUEST] IP: %s | UserID: %s | Method: %s | Path: %s | User-Agent: %s",
			clientIP,
			userID,
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.UserAgent(),
		)

		// 处理请求
		c.Next()

		// 记录响应信息
		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[RESPONSE] IP: %s | UserID: %s | Method: %s | Path: %s | Status: %d | Latency: %v",
			clientIP,
			userID,
			c.Request.Method,
			c.Request.URL.Path,
			status,
			latency,
		)
	}
}

// SimpleLoggingMiddleware 创建简化版日志中间件（不依赖用户信息）
func SimpleLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()

		// 获取客户端IP
		clientIP := c.ClientIP()

		// 记录请求信息
		log.Printf("[REQUEST] IP: %s | Method: %s | Path: %s | User-Agent: %s",
			clientIP,
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.UserAgent(),
		)

		// 处理请求
		c.Next()

		// 记录响应信息
		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[RESPONSE] IP: %s | Method: %s | Path: %s | Status: %d | Latency: %v",
			clientIP,
			c.Request.Method,
			c.Request.URL.Path,
			status,
			latency,
		)
	}
}
