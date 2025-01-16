package container

import (
	"single-ticker/config"
	"sync"
)

type LastTicker struct {
	BidPrice float64 // 买1价
	BidSize  float64 // 买1量
	AskPrice float64 // 卖1价
	AskSize  float64 // 卖1量
}

type TickerWrapper struct {
	Exchange config.Exchange
	InstType config.InstrumentType
	InstID   string

	BidPrice     float64 // 买1价
	BidSize      float64 // 买1量
	AskPrice     float64 // 卖1价
	AskSize      float64 // 卖1量
	UpdateTimeMs int64   //更新时间（微秒）
	UpdateID     int64   // 更新ID 400900217
	LastTicker   LastTicker
	IP           string // 建立连接时的sourceIP
	Colo         bool   // 是否是内网
}

func (wrapper *TickerWrapper) updateTicker(message *TickerWrapper) bool {
	if message.Exchange == config.BinanceExchange && message.InstType == config.SpotInstrument {
		// 币安现货没有消息更新时间（updateTimeMs），只有updateID
		if wrapper.UpdateID >= message.UpdateID {
			// message is expired
			return false
		}
	} else {
		if wrapper.UpdateTimeMs >= message.UpdateTimeMs {
			// message is expired
			return false
		}
	}

	if message.AskPrice == wrapper.LastTicker.AskPrice &&
		message.BidPrice == wrapper.LastTicker.BidPrice &&
		message.AskSize == wrapper.LastTicker.AskSize &&
		message.BidSize == wrapper.LastTicker.BidSize {
		// ticker值没有变化，返回false为了不推送重复的消息到做市服务，但是updateTime和updateID 还是会更新
		wrapper.UpdateTimeMs = message.UpdateTimeMs
		wrapper.UpdateID = message.UpdateID
		return false
	}

	if message.AskPrice > 0.0 && message.AskSize > 0.0 {
		wrapper.AskPrice = message.AskPrice
		wrapper.AskSize = message.AskSize
		wrapper.LastTicker.AskPrice = message.AskPrice
		wrapper.LastTicker.AskSize = message.AskSize
	}

	if message.BidPrice > 0.0 && message.BidSize > 0.0 {
		wrapper.BidPrice = message.BidPrice
		wrapper.BidSize = message.BidSize
		wrapper.LastTicker.BidPrice = message.BidPrice
		wrapper.LastTicker.BidSize = message.BidSize
	}

	wrapper.UpdateTimeMs = message.UpdateTimeMs
	wrapper.UpdateID = message.UpdateID
	return true
}

type TickerComposite struct {
	Exchange       config.Exchange
	InstType       config.InstrumentType
	tickerWrappers map[string]TickerWrapper
	rwLock         *sync.RWMutex
}

func NewTickerComposite(exchange config.Exchange, instType config.InstrumentType) *TickerComposite {
	return &TickerComposite{
		Exchange:       exchange,
		InstType:       instType,
		tickerWrappers: map[string]TickerWrapper{},
		rwLock:         new(sync.RWMutex),
	}
}

func (composite *TickerComposite) GetTicker(instID string) *TickerWrapper {
	composite.rwLock.RLock()
	wrapper, has := composite.tickerWrappers[instID]
	composite.rwLock.RUnlock()

	if has {
		return &wrapper
	} else {
		return nil
	}
}

func (composite *TickerComposite) UpdateTicker(message *TickerWrapper) bool {
	composite.rwLock.Lock()
	defer composite.rwLock.Unlock()

	updateResult := false
	if composite.Exchange != message.Exchange || composite.InstType != message.InstType {
		return updateResult
	}
	wrapper, has := composite.tickerWrappers[message.InstID]

	if !has {
		updateResult = wrapper.updateTicker(message)
		composite.tickerWrappers[message.InstID] = wrapper
	} else {
		updateResult = wrapper.updateTicker(message)
		composite.tickerWrappers[message.InstID] = wrapper
	}
	return updateResult
}
