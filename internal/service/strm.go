package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"PlexWarp/internal/config"
)

// StrmService strm文件处理服务
type StrmService struct {
}

// NewStrmService 创建strm服务
func NewStrmService() *StrmService {
	return &StrmService{}
}

// IsStrmFile 判断是否为strm文件
func (s *StrmService) IsStrmFile(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".strm")
}

// ReadStrmContent 读取strm文件内容
func (s *StrmService) ReadStrmContent(filePath string) (string, error) {
	// 应用路径映射
	mappedPath := s.applyPathMapping(filePath)
	
	// 检查文件是否存在
	if _, err := os.Stat(mappedPath); os.IsNotExist(err) {
		return "", fmt.Errorf("strm file not found: %s", mappedPath)
	}

	// 读取文件内容
	data, err := os.ReadFile(mappedPath)
	if err != nil {
		return "", fmt.Errorf("read strm file failed: %v", err)
	}

	content := strings.TrimSpace(string(data))
	if content == "" {
		return "", fmt.Errorf("strm file is empty")
	}

	log.Printf("Read strm content: %s -> %s", filePath, content)
	return content, nil
}

// GetDirectLinkFromStrm 从strm文件获取直链
func (s *StrmService) GetDirectLinkFromStrm(strmPath string) (string, error) {
	// 读取strm文件内容
	content, err := s.ReadStrmContent(strmPath)
	if err != nil {
		return "", err
	}

	// 如果内容已经是HTTP链接，直接返回
	if strings.HasPrefix(content, "http://") || strings.HasPrefix(content, "https://") {
		log.Printf("Strm contains direct HTTP link: %s", content)
		return content, nil
	}

	// 对于本地路径，当前版本不支持，返回错误
	if strings.HasPrefix(content, "/") {
		return "", fmt.Errorf("local path not supported in current version: %s", content)
	}

	return "", fmt.Errorf("unsupported strm content format: %s", content)
}

// applyPathMapping 应用路径映射规则
func (s *StrmService) applyPathMapping(path string) string {
	// 首先检查软链接规则
	for _, rule := range config.Symlink.Rules {
		if strings.HasPrefix(path, rule.Path) {
			mappedPath := strings.Replace(path, rule.Path, rule.Target, 1)
			log.Printf("Applied symlink rule: %s -> %s", path, mappedPath)
			return mappedPath
		}
	}

	// 然后应用媒体路径映射规则
	for _, rule := range config.PathMapping.Rules {
		if strings.HasPrefix(path, rule.From) {
			mappedPath := strings.Replace(path, rule.From, rule.To, 1)
			log.Printf("Applied path mapping: %s -> %s", path, mappedPath)
			return mappedPath
		}
	}

	return path
}

// IsMediaPath 判断路径是否在媒体挂载路径中
func (s *StrmService) IsMediaPath(path string) bool {
	for _, mountPath := range config.Plex302.MediaMountPaths {
		if strings.HasPrefix(path, mountPath) {
			return true
		}
	}
	return false
}

// ShouldRedirect 判断是否应该进行302重定向
func (s *StrmService) ShouldRedirect(path string, userAgent string) bool {
	// 检查功能是否启用
	if !config.Plex302.Enable {
		return false
	}

	// 检查是否为strm文件
	if !s.IsStrmFile(path) {
		return false
	}

	// 检查是否在媒体路径中
	if !s.IsMediaPath(path) {
		return false
	}

	// 检查是否为转码请求（如果禁用转码重定向）
	if !config.Plex302.TranscodeEnable && s.isTranscodeRequest(path) {
		return false
	}

	return true
}

// isTranscodeRequest 判断是否为转码请求
func (s *StrmService) isTranscodeRequest(path string) bool {
	// 检查路径中是否包含转码相关的关键词
	transcodeKeywords := []string{
		"transcode",
		"universal",
		"decision",
		"start",
	}

	pathLower := strings.ToLower(path)
	for _, keyword := range transcodeKeywords {
		if strings.Contains(pathLower, keyword) {
			return true
		}
	}

	return false
}

// HandleRedirect 处理302重定向
func (s *StrmService) HandleRedirect(w http.ResponseWriter, r *http.Request, strmPath string) error {
	// 获取直链
	directLink, err := s.GetDirectLinkFromStrm(strmPath)
	if err != nil {
		log.Printf("Failed to get direct link for strm: %v", err)
		
		// 如果启用了回退到原始链接
		if config.Plex302.FallbackOriginal {
			log.Printf("Falling back to original request")
			return fmt.Errorf("fallback to original")
		}
		
		return err
	}

	// 执行302重定向
	log.Printf("Redirecting to direct link: %s", directLink)
	w.Header().Set("Location", directLink)
	w.WriteHeader(http.StatusFound)
	return nil
}

// CheckStrmHealth 检查strm相关服务健康状态
func (s *StrmService) CheckStrmHealth() error {
	if !config.Plex302.Enable {
		return nil // 功能未启用，跳过检查
	}

	// 检查媒体挂载路径
	for _, mountPath := range config.Plex302.MediaMountPaths {
		if _, err := os.Stat(mountPath); os.IsNotExist(err) {
			log.Printf("Warning: media mount path does not exist: %s", mountPath)
		}
	}

	return nil
}