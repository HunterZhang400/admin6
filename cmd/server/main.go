package main

import (
	"admin6/config"
	"admin6/handler"
	"admin6/infra/database"
	"admin6/job"
	"admin6/middleware"
	"admin6/pkg/httpcros"
	"fmt"
	"time"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	puzzlePath = "4H1NrV37X7jgqwm4Q12"
)

func main() {
	config.Setup()
	database.InitGormDB(config.Cfg.Mysql)

	gin.SetMode(gin.ReleaseMode)
	g := gin.Default()
	g.Use(gin.Recovery())
	g.Use(httpcros.CROS)
	g.Use(middleware.SimpleLoggingMiddleware()) // 添加日志中间件
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	g.Use(middleware.AuthMiddleware())    // 应用英语接口认证中间件
	g.Use(middleware.LoggingMiddleware()) // 添加高级日志中间件（在认证之后）

	// 管理API路由
	adminGroup := g.Group(fmt.Sprintf("/%s/api", puzzlePath))
	{
		// 管理员认证（不需要认证的接口）
		adminGroup.POST("/admin/login", handler.AdminLogin)

		adminGroup.POST("/admin/logout", handler.AdminLogout)

		// 需要管理员认证的业务接口
		adminAuthGroup := adminGroup.Group("")
		adminAuthGroup.Use(middleware.AdminAuthMiddleware())
		{

		}
	}

	g.StaticFS("/%s/web", http.Dir(puzzlePath+"/web"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8010"
	}
	fmt.Printf("http://127.0.0.1:%s/%s/web/\n", port, puzzlePath)
	err := g.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
