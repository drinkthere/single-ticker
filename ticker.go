package main

import (
	"github.com/drinkthere/bybit"
	"single-ticker/config"
	"single-ticker/container"
	"single-ticker/message"
)

func startTickerMessage() {
	// 监听币安数据，并Bind端口，供远程订阅
	if globalConfig.TickerSource == config.BinanceSpotTickerSource {
		binanceSpotTickerChan := make(chan *container.TickerWrapper)
		message.StartBinanceSpotMarketWs(&globalConfig, &globalContext, binanceSpotTickerChan)
		message.StartGatherBinanceSpotBookTicker(&globalContext, binanceSpotTickerChan)
	} else if globalConfig.TickerSource == config.BinanceFuturesTickerSource {
		binanceFuturesTickerChan := make(chan *container.TickerWrapper)
		message.StartBinanceFuturesMarketWs(&globalConfig, &globalContext, binanceFuturesTickerChan)
		message.StartGatherBinanceFuturesBookTicker(&globalContext, binanceFuturesTickerChan)
	} else if globalConfig.TickerSource == config.BybitLinearTickerSource {
		// bybit linear ticker 100ms 一次
		bybitLinearTickerChan := make(chan *container.TickerWrapper)
		message.StartBybitLinearMarketWs(&globalConfig, &globalContext, bybitLinearTickerChan)
		message.StartGatherBybitLinearTicker(&globalConfig, &globalContext, bybitLinearTickerChan)

		// bybit depth level 1 是10ms一次
		bybitLinearDepthChan := make(chan *bybit.V5WebsocketPublicOrderBookResponse)
		message.StartBybitLinearDepthWs(&globalConfig, &globalContext, bybitLinearDepthChan)
		message.StartGatherBybitLinearDepth(&globalConfig, &globalContext, bybitLinearDepthChan)

		// bybit trade 是realtime的
		bybitLinearTradeChan := make(chan *bybit.V5WebsocketPublicTradeResponse)
		message.StartBybitLinearTradeWs(&globalConfig, &globalContext, bybitLinearTradeChan)
		message.StartGatherBybitLinearTrade(&globalConfig, &globalContext, bybitLinearTradeChan)

	} else if globalConfig.TickerSource == config.BybitSpotTickerSource {
		// bybit 现货是realtime的
		bybitSpotTickerChan := make(chan *container.TickerWrapper)
		message.StartBybitSpotMarketWs(&globalConfig, &globalContext, bybitSpotTickerChan)
		message.StartGatherBybitSpotTicker(&globalConfig, &globalContext, bybitSpotTickerChan)
	}
}
