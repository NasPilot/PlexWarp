package middleware

import (
	"PlexWarp/internal/config"
	"PlexWarp/internal/logging"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logging.AccessInfof("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
		return ""
	})
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logging.Errorf("Panic recovered: %v", recovered)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	})
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Plex-Token")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ClientFilter 客户端过滤中间件
func ClientFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.ClientFilter.Enable {
			c.Next()
			return
		}

		userAgent := c.GetHeader("User-Agent")
		if userAgent == "" {
			c.Next()
			return
		}

		isAllowed := false
		for _, client := range config.ClientFilter.Clients {
			if strings.Contains(strings.ToLower(userAgent), strings.ToLower(client)) {
				isAllowed = true
				break
			}
		}

		// 根据过滤模式决定是否允许访问
		if config.ClientFilter.Mode == "allow" && !isAllowed {
			logging.Warnf("客户端被拒绝访问: %s", userAgent)
			c.JSON(http.StatusForbidden, gin.H{"error": "Client not allowed"})
			c.Abort()
			return
		} else if config.ClientFilter.Mode == "deny" && isAllowed {
			logging.Warnf("客户端被拒绝访问: %s", userAgent)
			c.JSON(http.StatusForbidden, gin.H{"error": "Client denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimiter 简单的速率限制中间件
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以实现更复杂的速率限制逻辑
		// 目前只是一个占位符
		c.Next()
	}
}

// Security 安全头中间件
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}