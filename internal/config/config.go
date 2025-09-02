package config

import (
	"PlexWarp/constants"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
)

var (
	// 基础配置
	Port       int    = constants.DEFAULT_PORT
	Host       string = constants.DEFAULT_HOST
	RootDir    string // 程序根目录
	ConfigDir  string // 配置文件目录
	LogDir     string // 日志文件目录
	StaticDir  string // 静态文件目录
	ConfigFile string // 配置文件路径

	// Plex服务器配置
	PlexServer PlexServerSetting

	// 日志配置
	Logger LoggerSetting

	// 客户端过滤配置
	ClientFilter ClientFilterSetting

	// Plex302重定向配置
	Plex302 Plex302Setting

	// 路径映射配置
	PathMapping PathMappingConfig

	// 软链接配置
	Symlink SymlinkConfig

	// STRM重定向配置
	StrmRedirect StrmRedirectConfig
)

// Init 初始化配置
func Init(configPath string) error {
	// 获取程序根目录
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取程序路径失败: %v", err)
	}
	RootDir = filepath.Dir(execPath)

	// 设置目录路径
	ConfigDir = filepath.Join(RootDir, "config")
	LogDir = filepath.Join(RootDir, "logs")
	StaticDir = filepath.Join(RootDir, "static")

	// 创建必要目录
	if err := createDir(ConfigDir); err != nil {
		return err
	}
	if err := createDir(LogDir); err != nil {
		return err
	}
	if err := createDir(StaticDir); err != nil {
		return err
	}

	// 设置配置文件路径
	if configPath != "" {
		ConfigFile = configPath
	} else {
		ConfigFile = filepath.Join(ConfigDir, "config.yaml")
	}

	// 加载配置
	return loadConfig()
}

// loadConfig 加载配置文件
func loadConfig() error {
	viper.SetConfigFile(ConfigFile)
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在，使用默认配置并创建配置文件
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("创建配置文件失败: %v", err)
			}
		} else {
			return fmt.Errorf("读取配置文件失败: %v", err)
		}
	}

	// 解析配置到结构体
	Port = viper.GetInt("port")
	Host = viper.GetString("host")

	// Plex服务器配置
	PlexServer.ADDR = viper.GetString("plex_server.addr")
	PlexServer.AUTH = viper.GetString("plex_server.auth")

	// 日志配置
	Logger.AccessLogger.Console = viper.GetBool("logger.access_logger.console")
	Logger.AccessLogger.File = viper.GetBool("logger.access_logger.file")
	Logger.ServiceLogger.Console = viper.GetBool("logger.service_logger.console")
	Logger.ServiceLogger.File = viper.GetBool("logger.service_logger.file")

	// 客户端过滤配置
	ClientFilter.Enable = viper.GetBool("client_filter.enable")
	ClientFilter.Mode = viper.GetString("client_filter.mode")
	ClientFilter.Clients = viper.GetStringSlice("client_filter.client_list")

	// Plex302重定向配置
	Plex302.Enable = viper.GetBool("plex302.enable")
	Plex302.MediaMountPaths = viper.GetStringSlice("plex302.media_mount_paths")
	Plex302.TranscodeEnable = viper.GetBool("plex302.transcode_enable")
	Plex302.FallbackOriginal = viper.GetBool("plex302.fallback_original")
	Plex302.CheckLinkValidity = viper.GetBool("plex302.check_link_validity")

	// 加载路径映射规则
	pathMappingData := viper.Get("path_mapping.rules")
	if mappingSlice, ok := pathMappingData.([]interface{}); ok {
		for _, item := range mappingSlice {
			if mapping, ok := item.(map[string]interface{}); ok {
				rule := PathMappingRule{
					From: mapping["from"].(string),
					To:   mapping["to"].(string),
				}
				PathMapping.Rules = append(PathMapping.Rules, rule)
			}
		}
	}

	// 加载软链接规则
	symlinkRulesData := viper.Get("symlink.rules")
	if rulesSlice, ok := symlinkRulesData.([]interface{}); ok {
		for _, item := range rulesSlice {
			if rule, ok := item.(map[string]interface{}); ok {
				symlinkRule := SymlinkRule{
					Path:   rule["path"].(string),
					Target: rule["target"].(string),
				}
				Symlink.Rules = append(Symlink.Rules, symlinkRule)
			}
		}
	}

	// 加载STRM重定向规则
	strmRedirectData := viper.Get("strm_redirect.last_link_rules")
	if rulesSlice, ok := strmRedirectData.([]interface{}); ok {
		for _, item := range rulesSlice {
			if rule, ok := item.(map[string]interface{}); ok {
				strmRule := StrmRedirectRule{
					MatchType: rule["match_type"].(string),
					Patterns:  rule["patterns"].([]string),
					Action:    rule["action"].(string),
				}
				StrmRedirect.LastLinkRules = append(StrmRedirect.LastLinkRules, strmRule)
			}
		}
	}

	return nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	viper.SetDefault("port", constants.DEFAULT_PORT)
	viper.SetDefault("host", constants.DEFAULT_HOST)

	// Plex服务器默认配置
	viper.SetDefault("plex_server.addr", "http://localhost:32400")
	viper.SetDefault("plex_server.auth", "")

	// 日志默认配置
	viper.SetDefault("logger.access_logger.console", true)
	viper.SetDefault("logger.access_logger.file", true)
	viper.SetDefault("logger.service_logger.console", true)
	viper.SetDefault("logger.service_logger.file", true)

	// 客户端过滤默认配置
	viper.SetDefault("client_filter.enable", false)
	viper.SetDefault("client_filter.mode", "allow")
	viper.SetDefault("client_filter.client_list", []string{})

	// Plex302重定向默认配置
	viper.SetDefault("plex302.enable", false)
	viper.SetDefault("plex302.media_mount_paths", []string{"/mnt"})
	viper.SetDefault("plex302.transcode_enable", true)
	viper.SetDefault("plex302.fallback_original", true)
	viper.SetDefault("plex302.check_link_validity", false)

	// 路径映射默认配置
	viper.SetDefault("path_mapping.rules", []map[string]string{})

	// 软链接默认配置
	viper.SetDefault("symlink.rules", []map[string]string{})

	// STRM重定向默认配置
	viper.SetDefault("strm_redirect.rules", []map[string]string{})
}

// createDir 创建目录
func createDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// ListenAddr 返回监听地址
func ListenAddr() string {
	return Host + ":" + strconv.Itoa(Port)
}