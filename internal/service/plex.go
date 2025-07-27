package service

import (
	"PlexWarp/internal/config"
	"PlexWarp/internal/logging"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	PlexClient *http.Client
	PlexBaseURL string
	PlexToken string
)

// InitPlexService 初始化Plex服务
func InitPlexService() {
	PlexClient = &http.Client{
		Timeout: 30 * time.Second,
	}

	PlexBaseURL = strings.TrimSuffix(config.PlexServer.ADDR, "/")
	PlexToken = config.PlexServer.AUTH

	logging.Infof("Plex服务初始化完成，服务器地址: %s", PlexBaseURL)
}

// BuildPlexURL 构建Plex URL
func BuildPlexURL(path string, params map[string]string) string {
	u, err := url.Parse(PlexBaseURL + path)
	if err != nil {
		logging.Errorf("解析URL失败: %v", err)
		return ""
	}

	q := u.Query()
	if PlexToken != "" {
		q.Set("X-Plex-Token", PlexToken)
	}

	for key, value := range params {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()
	return u.String()
}

// ProxyRequest 代理请求到Plex服务器
func ProxyRequest(method, path string, params map[string]string, headers map[string]string) (*http.Response, error) {
	plexURL := BuildPlexURL(path, params)
	if plexURL == "" {
		return nil, fmt.Errorf("构建Plex URL失败")
	}

	req, err := http.NewRequest(method, plexURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 设置默认请求头
	req.Header.Set("User-Agent", "PlexWarp/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := PlexClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}

	return resp, nil
}

// CheckPlexConnection 检查Plex连接
func CheckPlexConnection() error {
	resp, err := ProxyRequest("GET", "/", nil, nil)
	if err != nil {
		return fmt.Errorf("连接Plex服务器失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Plex服务器响应异常: %d", resp.StatusCode)
	}

	logging.Info("Plex服务器连接正常")
	return nil
}