# PlexWarp

高性能的 Plex 媒体服务器代理工具，提供智能流媒体优化和多平台支持。

## 功能特性

- 🚀 **高性能代理**: 基于 Gin 框架的高性能 HTTP 代理服务
- 🔧 **灵活配置**: 支持 YAML 配置文件，可自定义各种参数
- 📝 **完善日志**: 多级别日志系统，支持文件和控制台输出
- 🛡️ **安全防护**: 内置安全中间件，防止常见攻击
- 🌐 **跨域支持**: 完整的 CORS 支持
- 📊 **健康检查**: 提供健康检查和版本信息 API
- 🔄 **客户端过滤**: 支持基于 IP 和 User-Agent 的客户端过滤
- 📱 **多平台支持**: 支持 Linux、Windows、macOS 多平台

## 快速开始

### 下载

从 [Releases](https://github.com/NasPilot/PlexWarp/releases) 页面下载适合您系统的预编译二进制文件。

### 配置

1. 复制示例配置文件：
```bash
cp config.yaml.example config/config.yaml
```

2. 编辑配置文件，设置您的 Plex 服务器信息：
```yaml
plex:
  server_url: "http://your-plex-server:32400"
  token: "your-plex-token"
```

### 运行

```bash
# 普通模式
./plexwarp

# 调试模式
./plexwarp -debug

# 查看版本
./plexwarp -version
```

## 配置说明

详细的配置选项请参考 `config.yaml.example` 文件中的注释说明。

## API 接口

- `GET /api/health` - 健康检查
- `GET /api/version` - 版本信息
- `/*` - Plex 代理（所有其他请求）

## 开发

### 构建

```bash
# 克隆仓库
git clone https://github.com/NasPilot/PlexWarp.git
cd PlexWarp

# 安装依赖
go mod tidy

# 构建
go build -o plexwarp .
```

### 测试

```bash
go test ./...
```

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 支持

如果您遇到问题或有建议，请在 [Issues](https://github.com/NasPilot/PlexWarp/issues) 页面提交。