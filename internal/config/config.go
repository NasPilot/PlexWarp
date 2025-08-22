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

	// Web配置
	Web WebSetting

	// 客户端过滤配置
	ClientFilter ClientFilterSetting

	// HTTPStrm配置
	HTTPStrm HTTPStrmSetting

	// AlistStrm配置
	AlistStrm AlistStrmSetting

	// 字幕配置
	Subtitle SubtitleSetting

	// Strm302重定向配置
	Strm302 Strm302Setting

	// Alist配置
	Alist AlistConfig

	// 重定向配置
	Redirect RedirectConfig
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
	PlexServer.Type = constants.PlexServerTypePlex
	PlexServer.ADDR = viper.GetString("plex_server.addr")
	PlexServer.AUTH = viper.GetString("plex_server.auth")

	// 日志配置
	Logger.AccessLogger.Console = viper.GetBool("logger.access_logger.console")
	Logger.AccessLogger.File = viper.GetBool("logger.access_logger.file")
	Logger.ServiceLogger.Console = viper.GetBool("logger.service_logger.console")
	Logger.ServiceLogger.File = viper.GetBool("logger.service_logger.file")

	// Web配置
	Web.Enable = viper.GetBool("web.enable")
	Web.Custom = viper.GetBool("web.custom")
	Web.Index = viper.GetBool("web.index")
	Web.Head = viper.GetString("web.head")
	Web.ExternalPlayerUrl = viper.GetBool("web.external_player_url")
	Web.Crx = viper.GetBool("web.crx")
	Web.ActorPlus = viper.GetBool("web.actor_plus")
	Web.FanartShow = viper.GetBool("web.fanart_show")
	Web.Danmaku = viper.GetBool("web.danmaku")
	Web.VideoTogether = viper.GetBool("web.video_together")

	// 客户端过滤配置
	ClientFilter.Enable = viper.GetBool("client_filter.enable")
	ClientFilter.Mode = constants.FilterMode(viper.GetString("client_filter.mode"))
	ClientFilter.ClientList = viper.GetStringSlice("client_filter.client_list")

	// HTTPStrm配置
	HTTPStrm.Enable = viper.GetBool("http_strm.enable")
	HTTPStrm.TransCode = viper.GetBool("http_strm.trans_code")
	HTTPStrm.FinalURL = viper.GetBool("http_strm.final_url")
	HTTPStrm.PrefixList = viper.GetStringSlice("http_strm.prefix_list")

	// AlistStrm配置
	AlistStrm.Enable = viper.GetBool("alist_strm.enable")
	AlistStrm.TransCode = viper.GetBool("alist_strm.trans_code")
	AlistStrm.RawURL = viper.GetBool("alist_strm.raw_url")

	// 字幕配置
	Subtitle.Enable = viper.GetBool("subtitle.enable")
	Subtitle.SRT2ASS = viper.GetBool("subtitle.srt2ass")
	Subtitle.ASSStyle = viper.GetStringSlice("subtitle.ass_style")
	Subtitle.SubSet = viper.GetBool("subtitle.sub_set")

	// Strm302重定向配置
	Strm302.Enable = viper.GetBool("strm302.enable")
	Strm302.MediaMountPath = viper.GetStringSlice("strm302.media_mount_path")
	Strm302.TranscodeEnable = viper.GetBool("strm302.transcode_enable")
	Strm302.FallbackOriginal = viper.GetBool("strm302.fallback_original")

	// Alist配置
	Alist.Addr = viper.GetString("alist.addr")
	Alist.Token = viper.GetString("alist.token")
	Alist.SignEnable = viper.GetBool("alist.sign_enable")
	Alist.SignExpireTime = viper.GetInt("alist.sign_expire_time")
	Alist.PublicAddr = viper.GetString("alist.public_addr")
	Alist.RawUrlMapping = viper.GetStringMapString("alist.raw_url_mapping")

	// 重定向配置
	Redirect.Enable = viper.GetBool("redirect.enable")
	Redirect.CheckEnable = viper.GetBool("redirect.check_enable")

	// 加载路径映射规则
	mediaPathMappingData := viper.Get("redirect.media_path_mapping")
	if mappingSlice, ok := mediaPathMappingData.([]interface{}); ok {
		for _, item := range mappingSlice {
			if mapping, ok := item.(map[string]interface{}); ok {
				rule := PathMappingRule{
					From: mapping["from"].(string),
					To:   mapping["to"].(string),
				}
				Redirect.MediaPathMapping = append(Redirect.MediaPathMapping, rule)
			}
		}
	}

	// 加载软链接规则
	symlinkRulesData := viper.Get("redirect.symlink_rules")
	if rulesSlice, ok := symlinkRulesData.([]interface{}); ok {
		for _, item := range rulesSlice {
			if rule, ok := item.(map[string]interface{}); ok {
				symlinkRule := SymlinkRule{
					Path:   rule["path"].(string),
					Target: rule["target"].(string),
				}
				Redirect.SymlinkRules = append(Redirect.SymlinkRules, symlinkRule)
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

	// Web默认配置
	viper.SetDefault("web.enable", true)
	viper.SetDefault("web.custom", false)
	viper.SetDefault("web.index", false)
	viper.SetDefault("web.head", "")
	viper.SetDefault("web.external_player_url", false)
	viper.SetDefault("web.crx", false)
	viper.SetDefault("web.actor_plus", false)
	viper.SetDefault("web.fanart_show", false)
	viper.SetDefault("web.danmaku", false)
	viper.SetDefault("web.video_together", false)

	// 客户端过滤默认配置
	viper.SetDefault("client_filter.enable", false)
	viper.SetDefault("client_filter.mode", "allow")
	viper.SetDefault("client_filter.client_list", []string{})

	// HTTPStrm默认配置
	viper.SetDefault("http_strm.enable", false)
	viper.SetDefault("http_strm.trans_code", true)
	viper.SetDefault("http_strm.final_url", false)
	viper.SetDefault("http_strm.prefix_list", []string{})

	// AlistStrm默认配置
	viper.SetDefault("alist_strm.enable", false)
	viper.SetDefault("alist_strm.trans_code", true)
	viper.SetDefault("alist_strm.raw_url", false)

	// 字幕默认配置
	viper.SetDefault("subtitle.enable", false)
	viper.SetDefault("subtitle.srt2ass", false)
	viper.SetDefault("subtitle.ass_style", []string{})
	viper.SetDefault("subtitle.sub_set", false)

	// Strm302重定向默认配置
	viper.SetDefault("strm302.enable", false)
	viper.SetDefault("strm302.media_mount_path", []string{"/mnt"})
	viper.SetDefault("strm302.transcode_enable", true)
	viper.SetDefault("strm302.fallback_original", true)

	// Alist默认配置
	viper.SetDefault("alist.addr", "")
	viper.SetDefault("alist.token", "")
	viper.SetDefault("alist.sign_enable", false)
	viper.SetDefault("alist.sign_expire_time", 3600)
	viper.SetDefault("alist.public_addr", "")
	viper.SetDefault("alist.raw_url_mapping", map[string]string{})

	// 重定向默认配置
	viper.SetDefault("redirect.enable", false)
	viper.SetDefault("redirect.check_enable", false)
	viper.SetDefault("redirect.media_path_mapping", []map[string]string{})
	viper.SetDefault("redirect.symlink_rules", []map[string]string{})
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