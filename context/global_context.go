package context

import (
	"single-ticker/client"
	"single-ticker/config"
	"single-ticker/container"
)

type GlobalContext struct {
	InstrumentComposite               *container.InstrumentComposite
	BTCETHSpotTickerUpdateChan        chan *container.TickerWrapper
	ALTSpotTickerUpdateChan           chan *container.TickerWrapper
	BTCETHFuturesTickerUpdateChan     chan *container.TickerWrapper
	ALTFuturesTickerUpdateChan        chan *container.TickerWrapper
	BybitBTCETHSpotTickerUpdateChan   chan *container.TickerWrapper
	BybitALTSpotTickerUpdateChan      chan *container.TickerWrapper
	BybitBTCETHLinearTickerUpdateChan chan *container.TickerWrapper
	BybitALTLinearTickerUpdateChan    chan *container.TickerWrapper
}

func (context *GlobalContext) Init(globalConfig *config.Config, globalBinanceClient *client.BinanceClient) {
	// 初始化交易对数据
	context.initInstrumentComposite(globalConfig)

	context.initChannel()
}

func (context *GlobalContext) initInstrumentComposite(globalConfig *config.Config) {
	instrumentComposite := container.NewInstrumentComposite(globalConfig)
	context.InstrumentComposite = instrumentComposite
}

func (context *GlobalContext) initChannel() {
	context.BTCETHSpotTickerUpdateChan = make(chan *container.TickerWrapper)
	context.ALTSpotTickerUpdateChan = make(chan *container.TickerWrapper)
	context.BTCETHFuturesTickerUpdateChan = make(chan *container.TickerWrapper)
	context.ALTFuturesTickerUpdateChan = make(chan *container.TickerWrapper)

	context.BybitBTCETHSpotTickerUpdateChan = make(chan *container.TickerWrapper)
	context.BybitALTSpotTickerUpdateChan = make(chan *container.TickerWrapper)
	context.BybitBTCETHLinearTickerUpdateChan = make(chan *container.TickerWrapper)
	context.BybitALTLinearTickerUpdateChan = make(chan *container.TickerWrapper)
}
