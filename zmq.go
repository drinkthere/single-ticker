package main

import (
	zmq "github.com/pebbe/zmq4"
	"google.golang.org/protobuf/proto"
	"os"
	"single-ticker/config"
	"single-ticker/container"
	"single-ticker/protocol/pb"
	"single-ticker/utils/logger"
	"time"
)

func StartZmq() {
	logger.Info("Start Binance Ticker ZMQ")
	if globalConfig.TickerSource == config.BinanceSpotTickerSource {
		if _, exists := globalConfig.PubBindIPCMap["BTCETH"]; exists {
			startZmq(&globalConfig, "BTCETH", globalContext.BTCETHSpotTickerUpdateChan)
		}
		if _, exists := globalConfig.PubBindIPCMap["ALTCOIN"]; exists {
			startZmq(&globalConfig, "ALTCOIN", globalContext.ALTSpotTickerUpdateChan)
		}
	}

	if globalConfig.TickerSource == config.BinanceFuturesTickerSource {
		if _, exists := globalConfig.PubBindIPCMap["BTCETH"]; exists {
			startZmq(&globalConfig, "BTCETH", globalContext.BTCETHFuturesTickerUpdateChan)
		}
		if _, exists := globalConfig.PubBindIPCMap["ALTCOIN"]; exists {
			startZmq(&globalConfig, "ALTCOIN", globalContext.ALTFuturesTickerUpdateChan)
		}
	}

	if globalConfig.TickerSource == config.BybitSpotTickerSource {
		if _, exists := globalConfig.PubBindIPCMap["BTCETH"]; exists {
			startZmq(&globalConfig, "BTCETH", globalContext.BybitBTCETHSpotTickerUpdateChan)
		}
		if _, exists := globalConfig.PubBindIPCMap["ALTCOIN"]; exists {
			startZmq(&globalConfig, "ALTCOIN", globalContext.BybitALTSpotTickerUpdateChan)
		}
	}

	if globalConfig.TickerSource == config.BybitLinearTickerSource {
		if _, exists := globalConfig.PubBindIPCMap["BTCETH"]; exists {
			startZmq(&globalConfig, "BTCETH", globalContext.BybitBTCETHLinearTickerUpdateChan)
		}
		if _, exists := globalConfig.PubBindIPCMap["ALTCOIN"]; exists {
			startZmq(&globalConfig, "ALTCOIN", globalContext.BybitALTLinearTickerUpdateChan)
		}
	}
}

func startZmq(cfg *config.Config, ipcKey string, tickerChan chan *container.TickerWrapper) {
	go func() {
		defer func() {
			logger.Warn("[StartZmq] %s %s Pub Service Listening Exited.", cfg.Exchange, ipcKey)
		}()

		logger.Info("[StartZmq] %s %s Start Pub Service.", cfg.Exchange, ipcKey)

		var ctx *zmq.Context
		var pub *zmq.Socket
		var err error
		isPubStopped := true
		for {
			if isPubStopped {
				ctx, err = zmq.NewContext()
				if err != nil {
					logger.Error("[StartZmq] %s %s New Context Error: %s", cfg.Exchange, ipcKey, err.Error())
					time.Sleep(time.Second * 1)
					continue
				}

				pub, err = ctx.NewSocket(zmq.PUB)
				if err != nil {
					logger.Error("[StartZmq] %s %s New Socket Error: %s", cfg.Exchange, ipcKey, err.Error())
					time.Sleep(time.Second * 1)
					continue
				}

				ipc, ok := cfg.PubBindIPCMap[ipcKey]
				if !ok {
					logger.Error("[StartZmq] %s %s Pub IPC is missing.", cfg.Exchange, ipcKey)
					os.Exit(1)
				}

				err = pub.Bind(ipc)
				if err != nil {
					logger.Error("[StartZmq] %s %s Bind to Local ZMQ %s Error: %s", cfg.Exchange, ipcKey, ipc, err.Error())
					time.Sleep(time.Second * 1)
					continue
				}
				isPubStopped = false
			}

			select {
			case ticker := <-tickerChan:
				instID := ticker.InstID
				md := &pb.TickerInfo{
					InstID:   instID,
					BestBid:  ticker.BidPrice,
					BestAsk:  ticker.AskPrice,
					BidSz:    ticker.BidSize,
					AskSz:    ticker.AskSize,
					UpdateID: ticker.UpdateID,
					EventTs:  ticker.UpdateTimeMs,
				}

				data, err := proto.Marshal(md)
				if err != nil {
					logger.Warn("[StartZmq]  %s %s  Error marshaling Ticker: %v", cfg.Exchange, ipcKey, err)
				}

				_, err = pub.Send(string(data), 0)
				if err != nil {
					logger.Warn("[StartZmq] %s %s Error %s sending Ticker Data: %v", cfg.Exchange, ipcKey, instID, err)
					isPubStopped = true
					pub.Close()
					time.Sleep(time.Second * 1)
					continue
				}
			}
		}
	}()
}
