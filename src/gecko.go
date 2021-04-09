package main

import (
	"fmt"
	"log"

	coingecko "github.com/superoo7/go-gecko/v3"
	geckoTypes "github.com/superoo7/go-gecko/v3/types"
)

type GeckoApp struct {
	Client *coingecko.Client
}

func (g *GeckoApp) initializeGecko() {
	// httpClient := &http.Client{
	// 	Timeout: time.Second * 10,
	// }
	g.Client = coingecko.NewClient(nil)
	log.Println("GEK: connected")

	// fetch and store in DB top 200 markets by name, symbol, marketcap
	// poll every hour and post to DB
	// compare changes

}

// func (g *GeckoApp) getSymbol(symbol string) (*geckoTypes.CoinsMarketItem, error) {

// }

func (g *GeckoApp) getSymbol(symbol string) (*geckoTypes.CoinsMarketItem, error) {
	pcp := geckoTypes.PriceChangePercentageObject
	priceChangePercentage := []string{pcp.PCP1h, pcp.PCP24h, pcp.PCP7d, pcp.PCP14d, pcp.PCP30d, pcp.PCP200d, pcp.PCP1y}
	market, err := g.Client.CoinsMarket("usd", []string{symbol},
		geckoTypes.OrderTypeObject.MarketCapDesc, 1, 1, true, priceChangePercentage)
	if err != nil {
		log.Println("GEK: err fetching symbol prie: ", err)
		return nil, err
	}
	for i, v := range *market {
		if i == 0 {
			// log.Printf("GECKO\t%s: %+v", symbol, v)
			return &v, nil
		}
	}
	return nil, fmt.Errorf("empty list returned")
}

func (g *GeckoApp) getSymbolPrice(symbol string) (float64, error) {
	market, err := g.getSymbol(symbol)
	if err != nil {
		return 0.0, err
	}
	return market.CurrentPrice, nil
}

func (g *GeckoApp) getCoin(symbol string) (*geckoTypes.CoinsID, error) {
	coin, err := g.Client.CoinsID(symbol, false, true, false, true, true, false)
	if err != nil {
		return nil, err
	}
	log.Printf("coinID: %+v", *coin)
	return coin, nil
}
