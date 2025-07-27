package handler

import (
	"PlexWarp/internal/config"
	"PlexWarp/internal/logging"
	"PlexWarp/internal/service"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Init 初始化处理器
func Init() error {
	// 检查Plex连接
	if err := service.CheckPlexConnection(); err != nil {
		return err
	}

	logging.Info("处理器初始化完成")
	return nil
}

// ProxyHandler Plex代理处理器
func ProxyHandler(c *gin.Context) {
	// 获取请求路径
	path := c.Request.URL.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// 获取查询参数
	params := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	// 获取请求头
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 && !isHopByHopHeader(key) {
			headers[key] = values[0]
		}
	}

	// 代理请求到Plex服务器
	resp, err := service.ProxyRequest(c.Request.Method, path, params, headers)
	if err != nil {
		logging.Errorf("代理请求失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "代理请求失败"})
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for key, values := range resp.Header {
		if !isHopByHopHeader(key) {
			for _, value := range values {
				c.Header(key, value)
			}
		}
	}

	// 设置状态码
	c.Status(resp.StatusCode)

	// 复制响应体
	if _, err := io.Copy(c.Writer, resp.Body); err != nil {
		logging.Errorf("复制响应体失败: %v", err)
	}

	// 记录访问日志
	logging.AccessInfof("%s %s %d", c.Request.Method, c.Request.URL.Path, resp.StatusCode)
}

// HealthHandler 健康检查处理器
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "PlexWarp",
		"version": config.Version().AppVersion,
	})
}

// VersionHandler 版本信息处理器
func VersionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, config.Version())
}

// isHopByHopHeader 检查是否为逐跳头部
func isHopByHopHeader(header string) bool {
	hopByHopHeaders := []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}

	header = strings.ToLower(header)
	for _, h := range hopByHopHeaders {
		if strings.ToLower(h) == header {
			return true
		}
	}
	return false
}