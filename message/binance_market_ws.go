package message

import (
	binanceSpot "github.com/dictxwang/go-binance"
	binanceFutures "github.com/dictxwang/go-binance/futures"
	"math/rand"
	"single-ticker/container"
	"strconv"
	"sync"

	//"strconv"
	"single-ticker/config"
	"single-ticker/context"
	"single-ticker/utils/logger"
	"time"
)

func StartBinanceSpotMarketWs(globalConfig *config.Config, globalContext *context.GlobalContext, binanceSpotTickerChan chan *container.TickerWrapper) {
	binanceSpotWs := newBinanceSpotMarketWebSocket(binanceSpotTickerChan)
	source := globalConfig.Source
	// 循环不同的IP，监听对应的tickers
	binanceSpotWs.startBinanceMarketWs(globalContext, source.IP)
	logger.Info("[BinanceTickerWebSocket] Start Listen Binance Spot Tickers ip:%s", source.IP)
	time.Sleep(1 * time.Second)
}

type BinanceSpotMarketWebSocket struct {
	tickerChan chan *container.TickerWrapper
	rwLock     *sync.RWMutex
}

func newBinanceSpotMarketWebSocket(tickerChan chan *container.TickerWrapper) *BinanceSpotMarketWebSocket {
	return &BinanceSpotMarketWebSocket{
		tickerChan: tickerChan,
		rwLock:     new(sync.RWMutex),
	}
}

func (ws *BinanceSpotMarketWebSocket) startBinanceMarketWs(globalContext *context.GlobalContext, ip string) {
	// 初始化现货行情监控，这里不分批订阅instID了，每个进程就订阅全部的instIDs
	instIDs := globalContext.InstrumentComposite.InstIDs
	innerSpot := newInnerBinanceSpotWebSocket(ws.tickerChan, ip)
	innerSpot.subscribeBookTickers(instIDs)
	logger.Info("[BSTickerWebSocket] Start Listen Binance Spot Tickers @%s", ip)
}

type innerBinanceSpotWebSocket struct {
	ip         string
	instIDs    []string
	tickerChan chan *container.TickerWrapper
	isStopped  bool
	stopChan   chan struct{}
	randGen    *rand.Rand
}

func newInnerBinanceSpotWebSocket(tickerChan chan *container.TickerWrapper, ip string) *innerBinanceSpotWebSocket {
	return &innerBinanceSpotWebSocket{
		ip:         ip,
		tickerChan: tickerChan,
		isStopped:  true,
		randGen:    rand.New(rand.NewSource(2)),
	}
}

func (iws *innerBinanceSpotWebSocket) handleTickerEvent(event *binanceSpot.WsBookTickerEvent) {
	if iws.randGen.Int31n(10000) < 2 {
		logger.Info("[BSTickerWebSocket] Binance Spot Event: %+v", event)
	}

	iws.tickerChan <- convertSpotEventToBinanceTickerMessage(event, iws.ip)
}

func (iws *innerBinanceSpotWebSocket) handleError(err error) {
	// 出错断开连接，再重连
	logger.Error("[BSTickerWebSocket] Binance Spot Handle Error And Reconnect Ws: %s", err.Error())
	iws.stopChan <- struct{}{}
	iws.isStopped = true
}

func (iws *innerBinanceSpotWebSocket) subscribeBookTickers(instIDs []string) {

	go func() {
		defer func() {
			logger.Warn("[BSTickerWebSocket] Binance Spot BookTicker Listening Exited.")
		}()
		for {
			if !iws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}
			var stopChan chan struct{}
			var err error
			if iws.ip == "" {
				_, stopChan, err = binanceSpot.WsCombinedBookTickerServe(instIDs, iws.handleTickerEvent, iws.handleError)
			} else {
				_, stopChan, err = binanceSpot.WsCombinedBookTickerServeWithIP(iws.ip, instIDs, iws.handleTickerEvent, iws.handleError)
			}
			if err != nil {
				logger.Error("[BSTickerWebSocket] Subscribe Binance Book Tickers Error: %s", err.Error())
				time.Sleep(time.Second * 1)
				continue
			}
			logger.Info("[BSTickerWebSocket] Subscribe Binance Book Tickers: %d", len(instIDs))
			// 重置channel和时间
			iws.stopChan = stopChan
			iws.isStopped = false
		}
	}()
}

func StartBinanceFuturesMarketWs(cfg *config.Config, globalContext *context.GlobalContext, binanceFuturesTickerChan chan *container.TickerWrapper) {
	source := cfg.Source
	binanceFuturesWs := newBinanceFuturesMarketWebSocket(binanceFuturesTickerChan)
	binanceFuturesWs.startBinanceMarketWs(globalContext, source)
	logger.Info("[BinanceTickerWebSocket] Start Listen Binance Futures Tickers, isPrivate:%t, ip:%s", source.Colo, source.IP)
}

type BinanceFuturesMarketWebSocket struct {
	tickerChan chan *container.TickerWrapper
	rwLock     *sync.RWMutex
}

func newBinanceFuturesMarketWebSocket(tickerChan chan *container.TickerWrapper) *BinanceFuturesMarketWebSocket {
	return &BinanceFuturesMarketWebSocket{
		tickerChan: tickerChan,
		rwLock:     new(sync.RWMutex),
	}
}

func (ws *BinanceFuturesMarketWebSocket) startBinanceMarketWs(globalContext *context.GlobalContext, source config.Source) {
	// 初始化合约行情监控
	instIDs := globalContext.InstrumentComposite.InstIDs
	innerFutures := newInnerBinanceFuturesWebSocket(ws.tickerChan, source)
	innerFutures.startBookTickers(instIDs)

	logger.Info("[BFTickerWebSocket] Start Listen Binance Futures Tickers")
}

type innerBinanceFuturesWebSocket struct {
	source     config.Source
	tickerChan chan *container.TickerWrapper
	isStopped  bool
	stopChan   chan struct{}
	randGen    *rand.Rand
}

func newInnerBinanceFuturesWebSocket(tickerChan chan *container.TickerWrapper, source config.Source) *innerBinanceFuturesWebSocket {
	return &innerBinanceFuturesWebSocket{
		source:     source,
		tickerChan: tickerChan,
		isStopped:  true,
		randGen:    rand.New(rand.NewSource(2)),
	}
}

func (iws *innerBinanceFuturesWebSocket) handleTickerEvent(event *binanceFutures.WsBookTickerEvent) {
	if iws.randGen.Int31n(10000) < 2 {
		logger.Info("[BFTickerWebSocket] Binance Futures Event: %+v", event)
	}
	iws.tickerChan <- convertFuturesEventToBinanceTickerMessage(event, iws.source)
}

func (iws *innerBinanceFuturesWebSocket) handleError(err error) {
	// 出错断开连接，再重连
	logger.Error("[BFTickerWebSocket] Binance Futures Handle Error And Reconnect Ws: %s", err.Error())
	iws.stopChan <- struct{}{}
	iws.isStopped = true
}

func (iws *innerBinanceFuturesWebSocket) startBookTickers(instIDs []string) {

	go func() {
		defer func() {
			logger.Warn("[BFTickerWebSocket] Binance Futures Flash Ticker Listening Exited.")
		}()
		for {
			if !iws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}
			var stopChan chan struct{}
			var err error
			if iws.source.Colo {
				binanceFutures.UseIntranet = true
			} else {
				binanceFutures.UseIntranet = false
			}
			if iws.source.IP == "" {
				if iws.source.Colo {
					_, stopChan, err = binanceFutures.WsCombinedBookTickerServeIfIntranet(instIDs, iws.handleTickerEvent, iws.handleError)
				} else {
					_, stopChan, err = binanceFutures.WsCombinedBookTickerServe(instIDs, iws.handleTickerEvent, iws.handleError)
				}
			} else {
				if iws.source.Colo {
					_, stopChan, err = binanceFutures.WsCombinedBookTickerServeWithIPIfIntranet(iws.source.IP, instIDs, iws.handleTickerEvent, iws.handleError)
				} else {
					_, stopChan, err = binanceFutures.WsCombinedBookTickerServeWithIP(iws.source.IP, instIDs, iws.handleTickerEvent, iws.handleError)
				}
			}

			if err != nil {
				logger.Error("[BFTickerWebSocket] Subscribe Binance Futures Book Tickers Error: %s", err.Error())
				time.Sleep(time.Second * 1)
				continue
			}
			logger.Info("[BFTickerWebSocket] Subscribe Binance Futures Book Tickers: %d", len(instIDs))
			// 重置channel和时间
			iws.stopChan = stopChan
			iws.isStopped = false
		}
	}()
}

func convertSpotEventToBinanceTickerMessage(ticker *binanceSpot.WsBookTickerEvent, ip string) *container.TickerWrapper {
	bestAskPrice, _ := strconv.ParseFloat(ticker.BestAskPrice, 64)
	bestAskQty, _ := strconv.ParseFloat(ticker.BestAskQty, 64)
	bestBidPrice, _ := strconv.ParseFloat(ticker.BestBidPrice, 64)
	bestBidQty, _ := strconv.ParseFloat(ticker.BestBidQty, 64)
	return &container.TickerWrapper{
		Exchange:     config.BinanceExchange,
		InstType:     config.SpotInstrument,
		InstID:       ticker.Symbol,
		AskPrice:     bestAskPrice,
		AskSize:      bestAskQty,
		BidPrice:     bestBidPrice,
		BidSize:      bestBidQty,
		UpdateTimeMs: time.Now().UnixMilli(),
		UpdateID:     ticker.UpdateID,
		IP:           ip,
		Colo:         false,
	}
}

func convertFuturesEventToBinanceTickerMessage(ticker *binanceFutures.WsBookTickerEvent, source config.Source) *container.TickerWrapper {
	bestAskPrice, _ := strconv.ParseFloat(ticker.BestAskPrice, 64)
	bestAskQty, _ := strconv.ParseFloat(ticker.BestAskQty, 64)
	bestBidPrice, _ := strconv.ParseFloat(ticker.BestBidPrice, 64)
	bestBidQty, _ := strconv.ParseFloat(ticker.BestBidQty, 64)
	return &container.TickerWrapper{
		Exchange:     config.BinanceExchange,
		InstType:     config.FuturesInstrument,
		InstID:       ticker.Symbol,
		AskPrice:     bestAskPrice,
		AskSize:      bestAskQty,
		BidPrice:     bestBidPrice,
		BidSize:      bestBidQty,
		UpdateTimeMs: ticker.Time,
		UpdateID:     ticker.UpdateID,
		IP:           source.IP,
		Colo:         source.Colo,
	}
}
