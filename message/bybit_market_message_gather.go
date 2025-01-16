package message

import (
	"github.com/drinkthere/bybit"
	"single-ticker/config"
	"single-ticker/container"
	"single-ticker/context"
	"single-ticker/utils"
	"single-ticker/utils/logger"
	"strconv"
)

func StartGatherBybitSpotTicker(globalConfig *config.Config, globalContext *context.GlobalContext, tickerChan chan *container.TickerWrapper) {

	go func() {
		defer func() {
			logger.Error("[GatherBybitMarketMsg] Bybit Spot Tickers Gather Exited.")
		}()

		for tickerMsg := range tickerChan {
			if !utils.InArray(tickerMsg.InstID, globalContext.InstrumentComposite.SpotInstIDs) {
				continue
			}

			if tickerMsg.InstID == "BTCUSDT" || tickerMsg.InstID == "ETHUSDT" {
				globalContext.BybitBTCETHSpotTickerUpdateChan <- tickerMsg
			} else {
				globalContext.BybitALTSpotTickerUpdateChan <- tickerMsg
			}
		}
	}()
}

func StartGatherBybitLinearTicker(globalConfig *config.Config, globalContext *context.GlobalContext, tickerChan chan *container.TickerWrapper) {

	go func() {
		defer func() {
			logger.Error("[GatherBybitMarketMsg] Bybit Linear Tickers Gather Exited.")
		}()
		for tickerMsg := range tickerChan {
			if !utils.InArray(tickerMsg.InstID, globalContext.InstrumentComposite.InstIDs) {
				continue
			}

			fulfillTickerMsg(globalContext, tickerMsg)
			instID := tickerMsg.InstID
			if instID == "BTCUSDT" || instID == "ETHUSDT" {
				globalContext.BybitBTCETHLinearTickerUpdateChan <- tickerMsg
			} else {
				globalContext.BybitALTLinearTickerUpdateChan <- tickerMsg
			}
		}
	}()
}

func fulfillTickerMsg(globalContext *context.GlobalContext, tickerMsg *container.TickerWrapper) {
}

func StartGatherBybitLinearDepth(globalConfig *config.Config, globalContext *context.GlobalContext, depthChan chan *bybit.V5WebsocketPublicOrderBookResponse) {

	go func() {
		defer func() {
			logger.Error("[GatherBybitDepthMsg] Bybit Linear Depth Gather Exited.")
		}()
		for t := range depthChan {
			instID := string(t.Data.Symbol)
			if !utils.InArray(instID, globalContext.InstrumentComposite.InstIDs) {
				continue
			}

			tickerMsg := convertBybitLinearDepthEventToTickerMessage(t)
			fulfillTickerMsg(globalContext, tickerMsg)

			if instID == "BTCUSDT" || instID == "ETHUSDT" {
				globalContext.BybitBTCETHLinearTickerUpdateChan <- tickerMsg
			} else {
				globalContext.BybitALTLinearTickerUpdateChan <- tickerMsg
			}
		}
	}()
}

func convertBybitLinearDepthEventToTickerMessage(event *bybit.V5WebsocketPublicOrderBookResponse) *container.TickerWrapper {
	ticker := event.Data
	instID := string(ticker.Symbol)

	bidPx, bidSz, askPx, askSz := 0.0, 0.0, 0.0, 0.0
	if len(ticker.Bids) > 0 {
		bidPx, _ = strconv.ParseFloat(ticker.Bids[0].Price, 64)
		bidSz, _ = strconv.ParseFloat(ticker.Bids[0].Size, 64)
	}

	if len(ticker.Asks) > 0 {
		askPx, _ = strconv.ParseFloat(ticker.Asks[0].Price, 64)
		askSz, _ = strconv.ParseFloat(ticker.Asks[0].Size, 64)
	}

	return &container.TickerWrapper{
		Exchange:     config.BybitExchange,
		InstType:     config.LinearInstrument,
		InstID:       instID,
		AskPrice:     askPx,
		AskSize:      askSz,
		BidPrice:     bidPx,
		BidSize:      bidSz,
		UpdateTimeMs: event.TimeStamp,
	}
}

func StartGatherBybitLinearTrade(globalConfig *config.Config, globalContext *context.GlobalContext, tradeChan chan *bybit.V5WebsocketPublicTradeResponse) {

	go func() {
		defer func() {
			logger.Error("[GatherBybitTradeMsg] Bybit Linear Trade Gather Exited.")
		}()
		for t := range tradeChan {
			for _, trade := range t.Data {
				instID := string(trade.Symbol)
				if !utils.InArray(instID, globalContext.InstrumentComposite.InstIDs) {
					continue
				}

				tickCfg := globalContext.InstrumentComposite.LinearTickMap[instID]
				tickerMsg := convertBybitLinearTradeEventToTickerMessage(tickCfg, trade, t.TimeStamp)
				fulfillTickerMsg(globalContext, tickerMsg)

				if instID == "BTCUSDT" || instID == "ETHUSDT" {
					globalContext.BybitBTCETHLinearTickerUpdateChan <- tickerMsg
				} else {
					globalContext.BybitALTLinearTickerUpdateChan <- tickerMsg
				}
			}

		}
	}()
}

func convertBybitLinearTradeEventToTickerMessage(tickerInfo *container.TickInfo, trade bybit.V5WebsocketPublicTradeData, updateTs int64) *container.TickerWrapper {
	price, _ := strconv.ParseFloat(trade.Trade, 64)
	size, _ := strconv.ParseFloat(trade.Value, 64)
	bidPx, bidSz, askPx, askSz := 0.0, 0.0, 0.0, 0.0
	if trade.Side == bybit.SideBuy {
		bidPx = price
		askPx = price + tickerInfo.TickSiz
		askSz = size
	} else {
		bidPx = price - tickerInfo.TickSiz
		bidSz = size
		askPx = price
	}

	return &container.TickerWrapper{
		Exchange:     config.BybitExchange,
		InstType:     config.LinearInstrument,
		InstID:       string(trade.Symbol),
		AskPrice:     askPx,
		AskSize:      askSz,
		BidPrice:     bidPx,
		BidSize:      bidSz,
		UpdateTimeMs: updateTs,
	}
}
