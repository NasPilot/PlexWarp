package handler

import (
	"PlexWarp/internal/config"
	"PlexWarp/internal/logging"
	"PlexWarp/internal/service"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	// 媒体文件路径正则表达式
	mediaFileRegex = regexp.MustCompile(`/library/parts/(\d+)/(\d+)/file`)
	// 转码相关路径正则表达式
	transcodeRegex = regexp.MustCompile(`/video/:/transcode/`)
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
	// 检查是否需要进行strm重定向
	if shouldHandleStrmRedirect(c.Request) {
		if handleStrmRedirect(c.Writer, c.Request) {
			return // 重定向成功，直接返回
		}
		// 重定向失败，继续正常代理流程
	}

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

// shouldHandleStrmRedirect 检查是否需要处理strm重定向
func shouldHandleStrmRedirect(r *http.Request) bool {
	// 检查是否为媒体文件请求或转码请求
	return mediaFileRegex.MatchString(r.URL.Path) || transcodeRegex.MatchString(r.URL.Path)
}

// handleStrmRedirect 处理strm重定向
func handleStrmRedirect(w http.ResponseWriter, r *http.Request) bool {
	// 检查plex302功能是否启用
	if !config.Plex302.Enable {
		return false
	}

	log.Printf("处理strm重定向请求: %s", r.URL.Path)
	
	// 创建strm服务实例
	strmService := service.NewStrmService()
	
	// 检查是否应该进行重定向
	userAgent := r.Header.Get("User-Agent")
	if !strmService.ShouldRedirect(r.URL.Path, userAgent) {
		return false
	}
	
	// 尝试从请求路径中提取文件路径信息
	filePath := extractFilePathFromRequest(r)
	if filePath == "" {
		log.Printf("无法从请求中提取文件路径: %s", r.URL.Path)
		return false
	}
	
	// 检查是否为strm文件
	if !strmService.IsStrmFile(filePath) {
		return false
	}
	
	// 尝试处理重定向
	err := strmService.HandleRedirect(w, r, filePath)
	if err != nil {
		log.Printf("strm重定向失败: %v", err)
		
		// 如果错误信息包含"fallback"，表示需要回退到原始请求
		if strings.Contains(err.Error(), "fallback") {
			return false
		}
		
		// 其他错误，返回错误响应
		http.Error(w, "Strm redirect failed", http.StatusInternalServerError)
		return true
	}
	
	log.Printf("strm重定向成功: %s", filePath)
	return true
}

// extractFilePathFromRequest 从请求中提取文件路径
func extractFilePathFromRequest(r *http.Request) string {
	// 这里需要根据Plex的API结构来解析文件路径
	// 由于我们无法直接从URL路径获取完整的文件路径，
	// 这里返回一个示例实现，实际使用时需要根据具体情况调整
	
	// 对于媒体文件请求，通常需要查询Plex数据库或API来获取实际文件路径
	// 这里提供一个简化的实现
	path := r.URL.Path
	
	// 如果是媒体文件请求，尝试从查询参数中获取路径信息
	if mediaFileRegex.MatchString(path) {
		// 从查询参数中获取文件路径（如果有的话）
		if filePath := r.URL.Query().Get("path"); filePath != "" {
			return filePath
		}
		
		// 或者从其他参数中获取
		if mediaPath := r.URL.Query().Get("file"); mediaPath != "" {
			return mediaPath
		}
	}
	
	// 如果无法从查询参数获取，返回空字符串
	// 实际实现中，这里应该查询Plex数据库来获取文件路径
	return ""
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