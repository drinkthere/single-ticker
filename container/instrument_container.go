package container

import (
	"github.com/drinkthere/bybit"
	"single-ticker/client"
	"single-ticker/config"
	"single-ticker/utils"
	"single-ticker/utils/logger"
	"strconv"
)

type TickInfo struct {
	TickSiz float64
	LotSz   float64
	MinSz   float64
}

type InstrumentComposite struct {
	InstIDs       []string             // 永续合约支持多个交易对，如：BTCUSDT
	SpotInstIDs   []string             // 现货，支持多个交易对，如：BTCUSDT
	LinearTickMap map[string]*TickInfo // bybit linear tick info
	SpotTickMap   map[string]*TickInfo // bybit spot tick info
}

func NewInstrumentComposite(globalConfig *config.Config) *InstrumentComposite {
	composite := &InstrumentComposite{
		InstIDs:     make([]string, 0),
		SpotInstIDs: make([]string, 0),
	}

	for _, instID := range globalConfig.InstIDs {
		composite.InstIDs = append(composite.InstIDs, instID)
		if globalConfig.Exchange == config.BybitExchange {
			composite.SpotInstIDs = append(composite.SpotInstIDs, utils.ConvertBybitLinearInstIDToSpotInstID(instID))
		} else {
			composite.SpotInstIDs = append(composite.SpotInstIDs, utils.ConvertBinanceFuturesInstIDToSpotInstID(instID))
		}
	}

	if globalConfig.Exchange == config.BybitExchange {
		var bybitClient = client.BybitClient{}
		bybitClient.Init()
		composite.SpotTickMap = getTickMap(&bybitClient, composite.SpotInstIDs, bybit.CategoryV5Spot)
		composite.LinearTickMap = getTickMap(&bybitClient, composite.InstIDs, bybit.CategoryV5Linear)
	}

	return composite
}

func getTickMap(bybitClient *client.BybitClient, instIDs []string, category bybit.CategoryV5) map[string]*TickInfo {
	tickMap := make(map[string]*TickInfo)

	// 标的现在数量不到1000，所以暂时先不分页了
	limit := 1000
	resp, err := bybitClient.Client.V5().Market().GetInstrumentsInfo(bybit.V5GetInstrumentsInfoParam{
		Category: category,
		Limit:    &limit,
	})

	if err != nil {
		logger.Fatal("[InitInstrumentComposite] Get %s instrument info Failed", category)
	}

	if resp.RetCode != 0 {
		logger.Fatal("[InitInstrumentComposite] Get %s instrument info Failed, Resp is %+v", category, resp)
	}

	if category == bybit.CategoryV5Spot {
		for _, inst := range resp.Result.Spot.List {
			instID := string(inst.Symbol)
			if utils.InArray(instID, instIDs) {
				tickSize, _ := strconv.ParseFloat(inst.PriceFilter.TickSize, 64)
				lotSz, _ := strconv.ParseFloat(inst.LotSizeFilter.BasePrecision, 64)
				minSz, _ := strconv.ParseFloat(inst.LotSizeFilter.MinOrderQty, 64)
				logger.Debug("%s instID is %s, inst is %+v", category, instID, inst)
				tickInfo := &TickInfo{
					TickSiz: tickSize,
					LotSz:   lotSz,
					MinSz:   minSz,
				}
				tickMap[instID] = tickInfo
			}
		}
	} else {
		for _, inst := range resp.Result.LinearInverse.List {
			instID := string(inst.Symbol)
			if utils.InArray(instID, instIDs) {
				tickSize, _ := strconv.ParseFloat(inst.PriceFilter.TickSize, 64)
				lotSz, _ := strconv.ParseFloat(inst.LotSizeFilter.QtyStep, 64)
				minSz, _ := strconv.ParseFloat(inst.LotSizeFilter.MinOrderQty, 64)
				logger.Debug("%s instID is %s, inst is %+v", category, instID, inst)
				tickInfo := &TickInfo{
					TickSiz: tickSize,
					LotSz:   lotSz,
					MinSz:   minSz,
				}
				tickMap[instID] = tickInfo
			}
		}
	}

	return tickMap
}
