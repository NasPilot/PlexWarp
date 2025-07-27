package router

import (
	"PlexWarp/internal/handler"
	"PlexWarp/internal/middleware"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.New()

	// 添加中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Security())
	r.Use(middleware.ClientFilter())

	// API路由组
	api := r.Group("/api")
	{
		api.GET("/health", handler.HealthHandler)
		api.GET("/version", handler.VersionHandler)
	}

	// Plex代理路由 - 捕获所有其他请求
	r.NoRoute(handler.ProxyHandler)

	return r
}