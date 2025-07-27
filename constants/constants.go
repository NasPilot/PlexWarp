package constants

const (
	// 时间格式
	FORMATE_TIME = "2006-01-02 15:04:05"

	// 默认配置
	DEFAULT_PORT = 8080
	DEFAULT_HOST = "0.0.0.0"

	// Plex服务器类型
	PLEX_SERVER = "plex"

	// 日志级别
	LOG_LEVEL_DEBUG = "debug"
	LOG_LEVEL_INFO  = "info"
	LOG_LEVEL_WARN  = "warn"
	LOG_LEVEL_ERROR = "error"

	// 过滤模式
	FILTER_MODE_ALLOW = "allow"
	FILTER_MODE_DENY  = "deny"
)

// PlexServerType Plex服务器类型
type PlexServerType string

const (
	PlexServerTypePlex PlexServerType = "plex"
)

// FilterMode 过滤模式
type FilterMode string

const (
	FilterModeAllow FilterMode = "allow"
	FilterModeDeny  FilterMode = "deny"
)

// LOGO PlexWarp启动Logo
const LOGO = `
 ____  _           __        __
|  _ \| | _____  __\ \      / /_ _ _ __ _ __
| |_) | |/ _ \ \/ / \ \ /\ / / _' | '__| '_ \
|  __/| |  __/>  <   \ V  V / (_| | |  | |_) |
|_|   |_|\___/_/\_\   \_/\_/ \__,_|_|  | .__/
                                      |_|
`