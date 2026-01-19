package middleware

import (
	"admin6/handler"
	"admin6/infra/database"
	"admin6/model"
	"admin6/pkg/common"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过管理API路由
		if strings.HasPrefix(c.Request.URL.Path, "/4H1NrV37X7jgqwm4Q12/") {
			c.Next()
			return
		}

		// 跳过静态文件路由
		if strings.HasPrefix(c.Request.URL.Path, "/web/") {
			c.Next()
			return
		}

		// 跳过登录API
		if strings.HasPrefix(c.Request.URL.Path, "/api/user/login") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/validate-session") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/user/logout") {
			c.Next()
			return
		}

		// 尝试获取用户会话，如果没有则使用默认用户
		sessionID := c.GetHeader("X-Session-ID")
		if sessionID == "" {
			// 尝试从Cookie获取会话ID
			if cookie, err := c.Cookie("user_session_id"); err == nil && cookie != "" {
				sessionID = cookie
			}
		}

		var userID uint = 0
		if sessionID != "" {
			// 验证会话
			user, err := ValidateUserSession(sessionID)
			if err == nil {
				userID = user.ID
			}
		}

		// 将用户信息存储到上下文中
		c.Set("userID", userID)
		c.Next()
	}
}

// ValidateUserSession validates a user session and returns user info
func ValidateUserSession(sessionID string) (*model.UserGorm, error) {
	db := database.GetGormDB()
	var session model.UserSession
	err := db.Where("token = ? AND expires_at > ?", sessionID, time.Now()).First(&session).Error
	if err != nil {
		return nil, fmt.Errorf("session not found or expired")
	}

	// Get user from database
	var user model.UserGorm
	err = db.Where("id = ?", session.UserID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先从Header获取会话ID，如果没有则从Cookie获取
		sessionID := c.GetHeader("X-Session-ID")
		if sessionID == "" {
			// 尝试从Cookie获取会话ID
			if cookie, err := c.Cookie("admin_session_id"); err == nil && cookie != "" {
				sessionID = cookie
			}
		}

		if sessionID == "" {
			c.JSON(http.StatusUnauthorized, common.Response{
				Code: -1,
				Msg:  "缺少会话ID，请先登录",
				Data: map[string]interface{}{
					"require_login": true,
					"redirect":      false,
				},
			})
			c.Abort()
			return
		}

		// 验证管理员会话
		admin, err := handler.ValidateAdminSession(sessionID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, common.Response{
				Code: -1,
				Msg:  "会话无效或已过期，请重新登录",
				Data: map[string]interface{}{
					"require_login": true,
					"redirect":      false,
				},
			})
			c.Abort()
			return
		}

		// 将管理员信息存储到上下文中
		c.Set("admin", admin)
		c.Set("adminID", admin.ID)
		c.Next()
	}
}
