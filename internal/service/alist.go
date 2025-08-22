package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"PlexWarp/internal/config"
)

// AlistService Alist服务
type AlistService struct {
	client *http.Client
}

// AlistResponse Alist API响应
type AlistResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// AlistFileInfo Alist文件信息
type AlistFileInfo struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	IsDir    bool   `json:"is_dir"`
	Modified string `json:"modified"`
	Sign     string `json:"sign"`
	Thumb    string `json:"thumb"`
	Type     int    `json:"type"`
	RawURL   string `json:"raw_url"`
	Readme   string `json:"readme"`
	Header   string `json:"header"`
	Provider string `json:"provider"`
	Related  []interface{} `json:"related"`
}

// NewAlistService 创建Alist服务
func NewAlistService() *AlistService {
	return &AlistService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetDirectLink 获取文件直链
func (a *AlistService) GetDirectLink(filePath string) (string, error) {
	if config.Alist.Addr == "" {
		return "", fmt.Errorf("alist address not configured")
	}

	// 构建API URL
	apiURL := strings.TrimRight(config.Alist.Addr, "/") + "/api/fs/get"

	// 构建请求体
	reqBody := map[string]interface{}{
		"path": filePath,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request body failed: %v", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request failed: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if config.Alist.Token != "" {
		req.Header.Set("Authorization", config.Alist.Token)
	}

	// 发送请求
	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response failed: %v", err)
	}

	// 解析响应
	var alistResp AlistResponse
	if err := json.Unmarshal(body, &alistResp); err != nil {
		return "", fmt.Errorf("unmarshal response failed: %v", err)
	}

	if alistResp.Code != 200 {
		return "", fmt.Errorf("alist api error: %s", alistResp.Message)
	}

	// 解析文件信息
	fileInfoData, ok := alistResp.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid file info data")
	}

	rawURL, ok := fileInfoData["raw_url"].(string)
	if !ok || rawURL == "" {
		return "", fmt.Errorf("raw_url not found")
	}

	// 处理URL映射
	if len(config.Alist.RawUrlMapping) > 0 {
		for from, to := range config.Alist.RawUrlMapping {
			if strings.Contains(rawURL, from) {
				rawURL = strings.Replace(rawURL, from, to, 1)
				break
			}
		}
	}

	// 如果启用了签名，添加签名参数
	if config.Alist.SignEnable {
		if sign, ok := fileInfoData["sign"].(string); ok && sign != "" {
			u, err := url.Parse(rawURL)
			if err == nil {
				q := u.Query()
				q.Set("sign", sign)
				u.RawQuery = q.Encode()
				rawURL = u.String()
			}
		}
	}

	log.Printf("Got direct link from Alist: path=%s, url=%s", filePath, rawURL)
	return rawURL, nil
}

// CheckConnection 检查Alist连接
func (a *AlistService) CheckConnection() error {
	if config.Alist.Addr == "" {
		return fmt.Errorf("alist address not configured")
	}

	apiURL := strings.TrimRight(config.Alist.Addr, "/") + "/api/me"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}

	if config.Alist.Token != "" {
		req.Header.Set("Authorization", config.Alist.Token)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("alist connection failed: status %d", resp.StatusCode)
	}

	return nil
}

// ConvertPathToAlist 将本地路径转换为Alist路径
func (a *AlistService) ConvertPathToAlist(localPath string) string {
	// 应用路径映射规则
	for _, rule := range config.Redirect.MediaPathMapping {
		if strings.HasPrefix(localPath, rule.From) {
			return strings.Replace(localPath, rule.From, rule.To, 1)
		}
	}

	// 如果没有匹配的映射规则，返回原路径
	return localPath
}

// IsMediaFile 判断是否为媒体文件
func (a *AlistService) IsMediaFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	mediaExts := []string{".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm", ".m4v", ".ts", ".m2ts"}
	
	for _, mediaExt := range mediaExts {
		if ext == mediaExt {
			return true
		}
	}
	return false
}