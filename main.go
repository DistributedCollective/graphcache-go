package main

import (
	"github.com/CrocSwap/graphcache-go/cache"
	"github.com/CrocSwap/graphcache-go/controller"
	"github.com/CrocSwap/graphcache-go/loader"
	"github.com/CrocSwap/graphcache-go/server"
	"github.com/CrocSwap/graphcache-go/views"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	netCfgPath := "./config/networks.json"
	netCfg := loader.LoadNetworkConfig(netCfgPath)
	cache := cache.New()
	onChain := loader.OnChainLoader{Cfg: netCfg}

	goerlChainConfig, _ := netCfg["goerli"]
	cntrl := controller.New(netCfg, cache)
	syncer := controller.NewSubgraphSyncer(cntrl, goerlChainConfig, "goerli")
	syncer.SyncStartup(goerlChainConfig, "goerli")

	mainnetChainConfig, _ := netCfg["mainnet"]
	pricingSyncer := controller.NewSubgraphSyncer(cntrl, mainnetChainConfig, "mainnet")
	pricingSyncer.SyncPricingSwaps()

	views := views.Views{Cache: cache, OnChain: &onChain}
	apiServer := server.APIWebServer{Views: &views}
	apiServer.Serve()
}
