package message

import (
	"single-ticker/container"
	"single-ticker/context"
	"single-ticker/utils/logger"
)

func StartGatherBinanceSpotBookTicker(globalContext *context.GlobalContext, tickChan chan *container.TickerWrapper) {
	go func() {
		defer func() {
			logger.Warn("[GatherBSTicker] Binance Spot Ticker Gather Exited.")
		}()
		for tickerMsg := range tickChan {
			if tickerMsg.InstID == "BTCUSDT" || tickerMsg.InstID == "ETHUSDT" {
				globalContext.BTCETHSpotTickerUpdateChan <- tickerMsg
			} else {
				globalContext.ALTSpotTickerUpdateChan <- tickerMsg
			}
		}
	}()
}

func StartGatherBinanceFuturesBookTicker(globalContext *context.GlobalContext, tickChan chan *container.TickerWrapper) {

	go func() {
		defer func() {
			logger.Warn("[GatherBFTicker] Binance Spot Ticker Gather Exited.")
		}()
		for tickerMsg := range tickChan {

			if tickerMsg.InstID == "BTCUSDT" || tickerMsg.InstID == "ETHUSDT" {
				globalContext.BTCETHFuturesTickerUpdateChan <- tickerMsg
			} else {
				globalContext.ALTFuturesTickerUpdateChan <- tickerMsg
			}
		}
	}()
}
