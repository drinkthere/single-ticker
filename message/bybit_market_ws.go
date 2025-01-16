package message

import (
	"context"
	"fmt"
	"github.com/drinkthere/bybit"
	"math/rand"
	"single-ticker/config"
	"single-ticker/container"
	mmContext "single-ticker/context"
	"single-ticker/utils/logger"
	"strconv"
	"sync"
	"time"
)

func StartBybitLinearMarketWs(
	cfg *config.Config, globalContext *mmContext.GlobalContext, tickerChan chan *container.TickerWrapper) {

	bybitMarketWs := newBybitLinearMarketWebsocket(tickerChan)
	source := cfg.Source
	// 循环不同的IP，监听对应的tickers
	bybitMarketWs.Start(cfg, globalContext, source)
	logger.Info("[BiybitTickerWebSocket] Start Listen Bybit Linear Tickers, isPrivate:%t, ip:%s", source.Colo, source.IP)
	time.Sleep(1 * time.Second)
}

type BybitLinearMarketWebSocket struct {
	tickerChan chan *container.TickerWrapper
	rwLock     *sync.RWMutex
}

func newBybitLinearMarketWebsocket(tickerChan chan *container.TickerWrapper) *BybitLinearMarketWebSocket {
	return &BybitLinearMarketWebSocket{
		tickerChan: tickerChan,
		rwLock:     new(sync.RWMutex),
	}
}

func (ws *BybitLinearMarketWebSocket) Start(cfg *config.Config, globalContext *mmContext.GlobalContext, source config.Source) {
	// 初始化合约行情监控
	instIDs := globalContext.InstrumentComposite.InstIDs
	innerLinear := newInnerBybitLinearWebSocket(ws.tickerChan, source)
	innerLinear.startBybitLinearTickers(instIDs, source, cfg.BybitPrivateWsUrl)
	logger.Info("[BybitLTickerWebSocket] Start Listen Bybit Linear Tickers")
}

type innerBybitLinearWebSocket struct {
	source     config.Source
	tickerChan chan *container.TickerWrapper
	isStopped  bool
	stopChan   chan struct{}
	randGen    *rand.Rand
}

func newInnerBybitLinearWebSocket(tickerChan chan *container.TickerWrapper, source config.Source) *innerBybitLinearWebSocket {
	return &innerBybitLinearWebSocket{
		source:     source,
		tickerChan: tickerChan,
		isStopped:  true,
		randGen:    rand.New(rand.NewSource(2)),
	}
}

func (ws *innerBybitLinearWebSocket) handleLinerTickerEvent(event bybit.V5WebsocketPublicTickerResponse) error {
	tickerInfo := event.Data.LinearInverse
	tickerMsg := convertBybitLinearTickerEventToTickerMessage(tickerInfo, event.TimeStamp, ws.source)
	ws.tickerChan <- tickerMsg
	return nil
}

func (ws *innerBybitLinearWebSocket) handleLinerTickerError(isWebsocketClosed bool, err error) {
	if err != nil {
		logger.Error("[BybitMarketWs] Bybit Liner Tickers Occur Error: %+v", err)
	}
	logger.Warn("[BybitMarketWs] Bybit Liner Tickers Ws Will Reconnect In 1 Second")
	time.Sleep(time.Second * 1)
	ws.isStopped = true
}

func (ws *innerBybitLinearWebSocket) startBybitLinearTickers(instIDs []string, source config.Source, privateWsUrl string) {

	go func() {
		defer func() {
			logger.Warn("[BybitMarketWs] Bybit Linear Tickers Listening Exited.")
		}()
		for {
			if !ws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}
			wsUrl := ""
			if source.Colo {
				wsUrl = privateWsUrl
			}
			wsClient := bybit.NewWebsocketClient(wsUrl)
			svc, err := wsClient.V5().Public(bybit.CategoryV5Linear)
			if err != nil {
				logger.Error("[BybitMarketWs] Start Bybit Linear Ws Failed, error: %+v", err)
				logger.Warn("[BybitMarketWs] Bybit Linear Tickers Ws Will Reconnect In 1 Second")
				time.Sleep(time.Second * 1)
				continue
			}

			for _, instID := range instIDs {
				_, err = svc.SubscribeTicker(bybit.V5WebsocketPublicTickerParamKey{
					Symbol: bybit.SymbolV5(instID),
				}, ws.handleLinerTickerEvent)

				if err != nil {
					logger.Error("[BybitMarketWs] Bybit Linear Subscribe %s Tickers Failed, error: %+v", instID, err)
					logger.Warn("[BybitMarketWs] Bybit Linear Tickers Ws Will Reconnect In 10 Second")
					time.Sleep(time.Second * 10)
					continue
				}
				time.Sleep(time.Millisecond * 100)
			}

			go svc.Start(context.Background(), ws.handleLinerTickerError)

			logger.Info("[BybitMarketWs] Bybit Linear Subscribe Tickers Successfully: %+v", instIDs)
			ws.isStopped = false
		}
	}()
}

func convertBybitLinearTickerEventToTickerMessage(ticker *bybit.V5WebsocketPublicTickerLinearInverseResult, updateTs int64, source config.Source) *container.TickerWrapper {
	instID := string(ticker.Symbol)

	bidPx, _ := strconv.ParseFloat(ticker.Bid1Price, 64)
	bidSz, _ := strconv.ParseFloat(ticker.Bid1Size, 64)
	askPx, _ := strconv.ParseFloat(ticker.Ask1Price, 64)
	askSz, _ := strconv.ParseFloat(ticker.Ask1Size, 64)
	return &container.TickerWrapper{
		Exchange:     config.BybitExchange,
		InstType:     config.LinearInstrument,
		InstID:       instID,
		AskPrice:     askPx,
		AskSize:      askSz,
		BidPrice:     bidPx,
		BidSize:      bidSz,
		UpdateTimeMs: updateTs,
		IP:           source.IP,
		Colo:         source.Colo,
	}
}

// StartBybitLinearDepthWs Subscribe Bybit Linear Depth
func StartBybitLinearDepthWs(
	cfg *config.Config, globalContext *mmContext.GlobalContext, depthChan chan *bybit.V5WebsocketPublicOrderBookResponse) {

	bybitDepthWs := newBybitLinearDepthWebsocket(depthChan)
	source := cfg.Source
	// 循环不同的IP，监听对应的tickers
	bybitDepthWs.Start(cfg, globalContext, source)
	logger.Info("[BiybitDepthWebSocket] Start Listen Bybit Linear Depth, isPrivate:%t, ip:%s", source.Colo, source.IP)
	time.Sleep(1 * time.Second)
}

type BybitLinearDepthWebSocket struct {
	depthChan chan *bybit.V5WebsocketPublicOrderBookResponse
}

func newBybitLinearDepthWebsocket(depthChan chan *bybit.V5WebsocketPublicOrderBookResponse) *BybitLinearDepthWebSocket {
	return &BybitLinearDepthWebSocket{
		depthChan: depthChan,
	}
}

func (ws *BybitLinearDepthWebSocket) Start(cfg *config.Config, globalContext *mmContext.GlobalContext, source config.Source) {
	// 初始化合约行情监控
	instIDs := globalContext.InstrumentComposite.InstIDs
	innerLinear := newInnerBybitLinearDepthWs(ws.depthChan, source)
	innerLinear.startBybitLinearDepth(instIDs, source, cfg.BybitPrivateWsUrl)

	logger.Info("[BybitLTickerWebSocket] Start Listen Bybit Linear Depth")
}

type innerBybitLinearDepthWs struct {
	source    config.Source
	depthChan chan *bybit.V5WebsocketPublicOrderBookResponse
	isStopped bool
	stopChan  chan struct{}
	randGen   *rand.Rand
}

func newInnerBybitLinearDepthWs(depthChan chan *bybit.V5WebsocketPublicOrderBookResponse, source config.Source) *innerBybitLinearDepthWs {
	return &innerBybitLinearDepthWs{
		source:    source,
		depthChan: depthChan,
		isStopped: true,
		randGen:   rand.New(rand.NewSource(2)),
	}
}

func (ws *innerBybitLinearDepthWs) handleLinerDepthEvent(event bybit.V5WebsocketPublicOrderBookResponse) error {
	ws.depthChan <- &event
	return nil
}

func (ws *innerBybitLinearDepthWs) handleLinerDepthError(isWebsocketClosed bool, err error) {
	if err != nil {
		logger.Error("[BybitMarketWs] Bybit Liner Depth Occur Error: %+v", err)
	}
	logger.Warn("[BybitMarketWs] Bybit Liner Depth Ws Will Reconnect In 1 Second")
	time.Sleep(time.Second * 1)
	ws.isStopped = true
}

func (ws *innerBybitLinearDepthWs) startBybitLinearDepth(instIDs []string, source config.Source, privateWsUrl string) {

	go func() {
		defer func() {
			logger.Warn("[BybitMarketWs] Bybit Linear Tickers Listening Exited.")
		}()
		for {
			if !ws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}
			wsUrl := ""
			if source.Colo {
				wsUrl = privateWsUrl
			}
			wsClient := bybit.NewWebsocketClient(wsUrl)
			svc, err := wsClient.V5().Public(bybit.CategoryV5Linear)

			if err != nil {
				logger.Error("[BybitMarketWs] Start Bybit Linear Depth Ws Failed, error: %+v", err)
				logger.Warn("[BybitMarketWs] Bybit Liner Depth Ws Will Reconnect In 1 Second")
				time.Sleep(time.Second * 1)
				continue
			}

			for _, instID := range instIDs {
				_, err = svc.SubscribeOrderBook(bybit.V5WebsocketPublicOrderBookParamKey{
					Depth:  1,
					Symbol: bybit.SymbolV5(instID),
				}, ws.handleLinerDepthEvent)

				if err != nil {
					logger.Error("[BybitMarketWs] Bybit Linear Subscribe %s Depth Failed, error: %+v", instID, err)
					logger.Warn("[BybitMarketWs] Bybit Liner Depth Ws Will Reconnect In 10 Second")
					time.Sleep(time.Second * 10)
					continue
				}
				time.Sleep(time.Millisecond * 100)
			}

			go svc.Start(context.Background(), ws.handleLinerDepthError)

			logger.Info("[BybitMarketWs] Bybit Linear Subscribe Depth Successfully: %+v", instIDs)
			ws.isStopped = false
		}
	}()
}

// StartBybitLinearTradeWs Subscribe Bybit Linear Trade
func StartBybitLinearTradeWs(
	cfg *config.Config, globalContext *mmContext.GlobalContext, tradeChan chan *bybit.V5WebsocketPublicTradeResponse) {
	bybitTradeWs := newBybitLinearTradeWebsocket(tradeChan)
	source := cfg.Source
	// 循环不同的IP，监听对应的tickers
	bybitTradeWs.Start(cfg, globalContext, source)
	logger.Info("[BiybitTradeWebSocket] Start Listen Bybit Linear Trade, isPrivate:%t, ip:%s", source.Colo, source.IP)
	time.Sleep(1 * time.Second)
}

type BybitLinearTradeWebSocket struct {
	tradeChan chan *bybit.V5WebsocketPublicTradeResponse
}

func newBybitLinearTradeWebsocket(tradeChan chan *bybit.V5WebsocketPublicTradeResponse) *BybitLinearTradeWebSocket {
	return &BybitLinearTradeWebSocket{
		tradeChan: tradeChan,
	}
}

func (ws *BybitLinearTradeWebSocket) Start(cfg *config.Config, globalContext *mmContext.GlobalContext, source config.Source) {
	// 初始化合约行情监控
	instIDs := globalContext.InstrumentComposite.InstIDs
	innerLinear := newInnerBybitLinearTradeWs(ws.tradeChan, source)
	innerLinear.startBybitLinearTrade(instIDs, source, cfg.BybitPrivateWsUrl)

	logger.Info("[BybitLTickerWebSocket] Start Listen Bybit Linear Trade")
}

type innerBybitLinearTradeWs struct {
	source    config.Source
	tradeChan chan *bybit.V5WebsocketPublicTradeResponse
	isStopped bool
	stopChan  chan struct{}
	randGen   *rand.Rand
}

func newInnerBybitLinearTradeWs(tradeChan chan *bybit.V5WebsocketPublicTradeResponse, source config.Source) *innerBybitLinearTradeWs {
	return &innerBybitLinearTradeWs{
		source:    source,
		tradeChan: tradeChan,
		isStopped: true,
		randGen:   rand.New(rand.NewSource(2)),
	}
}

func (ws *innerBybitLinearTradeWs) handleLinerTradeEvent(event bybit.V5WebsocketPublicTradeResponse) error {
	ws.tradeChan <- &event
	return nil
}

func (ws *innerBybitLinearTradeWs) handleLinerTradeError(isWebsocketClosed bool, err error) {
	if err != nil {
		logger.Error("[BybitMarketWs] Bybit Liner Trade Occur Error: %+v", err)
	}
	logger.Warn("[BybitMarketWs] Bybit Liner Trade Ws Will Reconnect In 1 Second")
	time.Sleep(time.Second * 1)
	ws.isStopped = true
}

func (ws *innerBybitLinearTradeWs) startBybitLinearTrade(instIDs []string, source config.Source, privateWsUrl string) {
	go func() {
		defer func() {
			logger.Warn("[BybitMarketWs] Bybit Linear Trade Listening Exited.")
		}()
		for {
			if !ws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}

			wsUrl := ""
			if source.Colo {
				wsUrl = privateWsUrl
			}
			wsClient := bybit.NewWebsocketClient(wsUrl)
			svc, err := wsClient.V5().Public(bybit.CategoryV5Linear)
			if err != nil {
				logger.Error("[BybitMarketWs] Start Bybit Linear Trade Websocket Failed, error: %+v", err)
				logger.Warn("[BybitMarketWs] Bybit Liner Trade Ws Will Reconnect In 1 Second")
				time.Sleep(time.Second * 1)
				continue
			}

			for _, instID := range instIDs {
				_, err = svc.SubscribeTrade(bybit.V5WebsocketPublicTradeParamKey{
					Symbol: bybit.SymbolV5(instID),
				}, ws.handleLinerTradeEvent)

				if err != nil {
					logger.Error("[BybitMarketWs] Bybit Linear Subscribe %s Trade Failed, error: %+v", instID, err)
					logger.Warn("[BybitMarketWs] Bybit Liner Trade Ws Will Reconnect In 10 Second")
					time.Sleep(time.Second * 10)
					continue
				}
				time.Sleep(time.Millisecond * 100)
			}

			go svc.Start(context.Background(), ws.handleLinerTradeError)

			logger.Info("[BybitMarketWs] Bybit Linear Subscribe Trade: %+v", instIDs)
			ws.isStopped = false
		}
	}()
}

// StartBybitSpotMarketWs Subscribe Bybit Spot Ticker
func StartBybitSpotMarketWs(
	cfg *config.Config, globalContext *mmContext.GlobalContext, tickerChan chan *container.TickerWrapper) {

	bybitMarketWs := newBybitSpotMarketWebsocket(tickerChan)
	source := cfg.Source
	// 循环不同的IP，监听对应的tickers
	bybitMarketWs.Start(cfg, globalContext, source)
	logger.Info("[BybitSpotTickerWebSocket] Start Listen Bybit Spot Tickers, isPrivate:%t, ip:%s", source.Colo, source.IP)
	time.Sleep(1 * time.Second)
}

type BybitSpotMarketWebSocket struct {
	tickerChan chan *container.TickerWrapper
	rwLock     *sync.RWMutex
}

func newBybitSpotMarketWebsocket(tickerChan chan *container.TickerWrapper) *BybitSpotMarketWebSocket {
	return &BybitSpotMarketWebSocket{
		tickerChan: tickerChan,
		rwLock:     new(sync.RWMutex),
	}
}

func (ws *BybitSpotMarketWebSocket) Start(cfg *config.Config, globalContext *mmContext.GlobalContext, source config.Source) {
	// 初始化合约行情监控
	instIDs := globalContext.InstrumentComposite.SpotInstIDs
	innerSpot := newInnerBybitSpotWebSocket(globalContext, ws.tickerChan, source)
	innerSpot.startBybitSpotTickers(instIDs, source, cfg.BybitPrivateWsUrl)
	logger.Info("[BFTickerWebSocket] Start Listen Bybit Spot Tickers")
}

type innerBybitSpotWebSocket struct {
	ctxt       *mmContext.GlobalContext
	source     config.Source
	tickerChan chan *container.TickerWrapper
	isStopped  bool
	stopChan   chan struct{}
	randGen    *rand.Rand
}

func newInnerBybitSpotWebSocket(globalContext *mmContext.GlobalContext, tickerChan chan *container.TickerWrapper, source config.Source) *innerBybitSpotWebSocket {
	return &innerBybitSpotWebSocket{
		ctxt:       globalContext,
		source:     source,
		tickerChan: tickerChan,
		isStopped:  true,
		randGen:    rand.New(rand.NewSource(2)),
	}
}

func (ws *innerBybitSpotWebSocket) handleSpotTickerEvent(event bybit.V5WebsocketPublicTickerResponse) error {
	tickerInfo := event.Data.Spot
	instID := string(tickerInfo.Symbol)
	SpotTickMap := ws.ctxt.InstrumentComposite.SpotTickMap
	tickCfg := SpotTickMap[instID]
	if tickCfg == nil {
		return fmt.Errorf("[BybitMarketWs] not spot tick config for %s", instID)
	}
	tickerMsg := convertBybitSpotTickerEventToTickerMessage(tickCfg, tickerInfo, event.TimeStamp, ws.source)
	ws.tickerChan <- tickerMsg
	return nil
}

func (ws *innerBybitSpotWebSocket) handleSpotTickerError(isWebsocketClosed bool, err error) {
	if err != nil {
		logger.Error("[BybitMarketWs] Bybit Spot Tickers Occur Error: %+v", err)
	}
	logger.Warn("[BybitMarketWs] Bybit Spot Tickers Ws Will Reconnect In 1 Second")
	time.Sleep(time.Second * 1)
	ws.isStopped = true
}

func (ws *innerBybitSpotWebSocket) startBybitSpotTickers(instIDs []string, source config.Source, privateWsUrl string) {
	go func() {
		defer func() {
			logger.Warn("[BybitMarketWs] Bybit Spot Tickers Listening Exited.")
		}()
		for {
			if !ws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}
			logger.Info("[BybitMarketWs] Start Create Bybit Spot Tickers Ws Connection")
			wsUrl := ""
			if source.Colo {
				wsUrl = privateWsUrl
			}
			wsClient := bybit.NewWebsocketClient(wsUrl)
			svc, err := wsClient.V5().Public(bybit.CategoryV5Spot)
			if err != nil {
				logger.Error("[BybitMarketWs] Start Bybit Spot Ws Failed, error: %+v", err)
				logger.Warn("[BybitMarketWs] Bybit Spot Tickers Ws Will Reconnect In 1 Second")
				time.Sleep(time.Second * 1)
				continue
			}

			for _, instID := range instIDs {
				_, err = svc.SubscribeTicker(bybit.V5WebsocketPublicTickerParamKey{
					Symbol: bybit.SymbolV5(instID),
				}, ws.handleSpotTickerEvent)

				if err != nil {
					logger.Error("[BybitMarketWs] Bybit Spot Subscribe %s Tickers Failed, error: %+v", instID, err)
					logger.Warn("[BybitMarketWs] Bybit Spot Tickers Ws Will Reconnect In 10 Second")
					time.Sleep(time.Second * 10)
					continue
				}
				time.Sleep(time.Millisecond * 100)
			}

			go svc.Start(context.Background(), ws.handleSpotTickerError)

			logger.Info("[BybitMarketWs] Bybit Spot Subscribe Tickers Successfully: %+v", instIDs)
			ws.isStopped = false
		}
	}()
}

func convertBybitSpotTickerEventToTickerMessage(tick *container.TickInfo, ticker *bybit.V5WebsocketPublicTickerSpotResult, updateTs int64, source config.Source) *container.TickerWrapper {
	instID := string(ticker.Symbol)

	// bybit的spot ticker数据中，没有bid price和ask price，所以需要我们自己计算
	lastPrice, _ := strconv.ParseFloat(ticker.LastPrice, 64)
	sz := tick.MinSz
	bidPx := lastPrice - tick.TickSiz
	askPx := lastPrice + tick.TickSiz

	return &container.TickerWrapper{
		Exchange:     config.BybitExchange,
		InstType:     config.SpotInstrument,
		InstID:       instID,
		AskPrice:     askPx,
		AskSize:      sz,
		BidPrice:     bidPx,
		BidSize:      sz,
		UpdateTimeMs: updateTs,
		IP:           source.IP,
		Colo:         source.Colo,
	}
}
