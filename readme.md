# PlexWarp

PlexWarp 是一个高性能的 Plex 媒体服务器代理工具，提供增强的功能和更好的用户体验。

## 功能特性

- 🚀 高性能 Plex 服务器代理
- 🔒 客户端访问控制和过滤
- 📊 详细的访问日志记录
- 🎨 Web 前端自定义支持
- 🔧 灵活的配置选项
- 📱 跨平台支持 (Linux, Windows, macOS)
- 🎯 多架构支持 (amd64, arm, arm64)

## 快速开始

### 下载

从 [Releases](https://github.com/your-username/PlexWarp/releases) 页面下载适合您系统的版本。

### 配置

1. 复制示例配置文件：
```bash
cp config.yaml.example config.yaml
```

2. 编辑配置文件，设置您的 Plex 服务器地址：
```yaml
plex_server:
  addr: "http://your-plex-server:32400"
  auth: "your-plex-token"  # 可选
```

### 运行

```bash
# 使用默认配置
./plexwarp

# 指定配置文件
./plexwarp -config /path/to/config.yaml

# 启用调试模式
./plexwarp -debug

# 查看版本信息
./plexwarp -version
```

## 配置说明

### 基础配置

- `port`: 服务监听端口 (默认: 8080)
- `host`: 服务监听地址 (默认: 0.0.0.0)

### Plex 服务器配置

- `plex_server.addr`: Plex 服务器地址
- `plex_server.auth`: Plex Token (可选，用于认证)

### 日志配置

- `logger.access_logger`: 访问日志配置
- `logger.service_logger`: 服务日志配置

### 客户端过滤

- `client_filter.enable`: 启用客户端过滤
- `client_filter.mode`: 过滤模式 (allow/deny)
- `client_filter.client_list`: 客户端列表

## API 接口

### 健康检查
```
GET /api/health
```

### 版本信息
```
GET /api/version
```

## 开发

### 构建

```bash
# 安装依赖
go mod tidy

# 构建
go build -o plexwarp .

# 运行
./plexwarp
```

### 测试

```bash
go test ./...
```

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 支持

如果您遇到问题或有建议，请在 [Issues](https://github.com/your-username/PlexWarp/issues) 页面提交。