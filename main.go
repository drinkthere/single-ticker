package main

import (
	"fmt"
	_ "net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"single-ticker/client"
	"single-ticker/config"
	"single-ticker/context"
	"single-ticker/utils"
	"single-ticker/utils/logger"
	"single-ticker/watchdog"
	"time"
)

var globalConfig config.Config
var globalBinanceClient client.BinanceClient
var globalContext context.GlobalContext

func main() {
	runtime.GOMAXPROCS(1)
	// 参数判断
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s config_file\n", os.Args[0])
		os.Exit(1)
	}
	utils.RegisterExitSignal(func() {
		os.Exit(1)
	})

	// 加载配置文件
	globalConfig = *config.LoadConfig(os.Args[1])

	// 设置日志级别, 并初始化日志
	logger.InitLogger(globalConfig.LogPath, globalConfig.LogLevel)

	// 初始化okx客户端
	globalBinanceClient.Init(&globalConfig)

	// 解析config，加载杠杆和合约交易对，初始化context，账户初始化设置，拉取仓位、余额等
	globalContext.Init(&globalConfig, &globalBinanceClient)

	time.Sleep(1 * time.Second)

	// 启动zmq
	StartZmq()

	// 开始监听ticker消息
	startTickerMessage()

	watchdog.StartPprofNet(&globalConfig)

	// 阻塞主进程
	for {
		time.Sleep(24 * time.Hour)
	}
}
