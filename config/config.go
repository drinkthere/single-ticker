package config

import (
	"encoding/json"
	"go.uber.org/zap/zapcore"
	"os"
)

type Source struct {
	IP   string // local IP to connect to OKx websocket
	Colo bool   // is co-location with Okx
}

type Config struct {
	Exchange Exchange
	// 日志配置
	LogLevel zapcore.Level
	LogPath  string

	// 币安配置
	BinanceAPIKey    string
	BinanceSecretKey string

	Source            Source // ip、channel等配置信息
	BybitPrivateWsUrl string // bybit 内网ws地址
	TickerSource      TickerSource
	PubBindIPCMap     map[string]string // 发布ticker数据时绑定的IPC

	InstIDs            []string // 要套利的交易对
	PprofListenAddress string
}

func LoadConfig(filename string) *Config {
	config := new(Config)
	reader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// 加载配置
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}
