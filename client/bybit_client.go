package client

import (
	"github.com/drinkthere/bybit"
)

type BybitClient struct {
	Client *bybit.Client
}

func (cli *BybitClient) Init() bool {
	restBaseUrl := ""
	cli.Client = bybit.NewClient(restBaseUrl).WithAuth("", "")
	return true
}
