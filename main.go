package main

import (
	"PlexWarp/constants"
	"PlexWarp/internal/config"
	"PlexWarp/internal/handler"
	"PlexWarp/internal/logging"
	"PlexWarp/internal/router"
	"PlexWarp/internal/service"
	"PlexWarp/utils"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	isDebug     bool   // 开启调试模式
	showVersion bool   // 显示版本信息
	configPath  string // 配置文件路径
)

func init() {
	flag.BoolVar(&showVersion, "version", false, "显示版本信息")
	flag.BoolVar(&isDebug, "debug", false, "是否启用调试模式")
	flag.StringVar(&configPath, "config", "", "指定配置文件路径")
	flag.Parse()

	fmt.Print(constants.LOGO)
	fmt.Println(utils.Center(fmt.Sprintf(" PlexWarp %s ", config.Version().AppVersion), 71, "="))
}

func main() {
	if showVersion {
		versionInfo, _ := json.MarshalIndent(config.Version(), "", "  ")
		fmt.Println(string(versionInfo))
		return
	}

	gin.SetMode(gin.ReleaseMode)

	if isDebug {
		logging.SetLevel(logrus.DebugLevel)
		fmt.Println("已启用调试模式")
	}

	signChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		fmt.Println("PlexWarp 已退出")
	}()

	if err := config.Init(configPath); err != nil { // 初始化配置
		fmt.Println("配置初始化失败：", err)
		return
	}
	logging.Init()                                                                         // 初始化日志
	logging.Infof("Plex服务器地址：%s", config.PlexServer.ADDR)                              // 日志打印
	service.InitPlexService()                                                              // 初始化Plex服务
	if err := handler.Init(); err != nil {                                                 // 初始化处理器
		logging.Error("Plex处理器初始化失败：", err)
		return
	}

	logging.Info("PlexWarp 监听端口：", config.Port)
	ginR := router.InitRouter() // 路由初始化
	logging.Info("PlexWarp 启动成功")
	go func() {
		if err := ginR.Run(config.ListenAddr()); err != nil {
			errChan <- err
		}
	}()

	select {
	case sig := <-signChan:
		logging.Info("PlexWarp 正在退出，信号：", sig)
	case err := <-errChan:
		logging.Error("PlexWarp 运行出错：", err)
	}
}