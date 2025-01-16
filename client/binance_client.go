package client

import (
	"github.com/dictxwang/go-binance"
	"github.com/dictxwang/go-binance/futures"
	"github.com/dictxwang/go-binance/portfolio"
	"single-ticker/config"
)

type BinanceClient struct {
	futuresClient   *futures.Client
	spotClient      *binance.Client
	PortfolioClient *portfolio.Client
}

func (cli *BinanceClient) Init(cfg *config.Config) bool {
	cli.futuresClient = futures.NewClient(cfg.BinanceAPIKey, cfg.BinanceSecretKey)
	return true
}
