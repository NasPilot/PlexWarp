package config

// 程序版本信息
type VersionInfo struct {
	AppVersion string // 程序版本号
	CommitHash string // Git Commit Hash
	BuildData  string // 编译时间
	GoVersion  string // 编译 Golang 版本
	OS         string // 操作系统
	Arch       string // 架构
}

// Plex服务器相关设置
type PlexServerSetting struct {
	ADDR string // 地址
	AUTH string // 认证授权TOKEN
}

// 日志设置
type LoggerSetting struct {
	AccessLogger  BaseLoggerSetting // 访问日志相关配置
	ServiceLogger BaseLoggerSetting // 服务日志相关配置
}

// 基础日志配置字段
type BaseLoggerSetting struct {
	Console bool // 是否将日志输出到终端中
	File    bool // 是否将日志输出到文件中
}

// 客户端User-Agent过滤设置
type ClientFilterSetting struct {
	Enable  bool     // 启用过滤
	Mode    string   // 过滤模式：whitelist 或 blacklist
	Clients []string // 客户端列表
}

// Plex302重定向设置
type Plex302Setting struct {
	Enable            bool     // 启用302重定向功能
	MediaMountPaths   []string // 媒体挂载路径列表
	TranscodeEnable   bool     // 是否允许转码
	FallbackOriginal  bool     // 失败时是否回退到原始链接
	CheckLinkValidity bool     // 是否检查链接有效性
}

// 路径映射规则
type PathMappingRule struct {
	From string // 源路径
	To   string // 目标路径
}

// 路径映射配置
type PathMappingConfig struct {
	Enable bool               // 启用路径映射
	Rules  []PathMappingRule  // 映射规则列表
}

// 软链接规则
type SymlinkRule struct {
	Path   string // 路径匹配规则
	Target string // 目标路径
}

// 软链接配置
type SymlinkConfig struct {
	Enable bool          // 启用软链接处理
	Rules  []SymlinkRule // 软链接规则列表
}

// STRM重定向规则
type StrmRedirectRule struct {
	MatchType string   // 匹配类型：startswith, endswith, contains, regex
	Patterns  []string // 匹配模式列表
	Action    string   // 动作：proxy, redirect
}

// STRM重定向配置
type StrmRedirectConfig struct {
	Enable        bool                // 启用STRM重定向
	LastLinkRules []StrmRedirectRule  // 最终链接处理规则
}