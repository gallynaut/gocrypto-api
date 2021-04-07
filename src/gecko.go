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
	log.Println("connected to coingecko")
	go g.getSolanPrice()

}

func (g *GeckoApp) getSolanPrice() {
	vsCurrency := "usd"
	ids := []string{"solana"}
	perPage := 1
	page := 1
	sparkline := true
	pcp := geckoTypes.PriceChangePercentageObject
	priceChangePercentage := []string{pcp.PCP1h, pcp.PCP24h, pcp.PCP7d, pcp.PCP14d, pcp.PCP30d, pcp.PCP200d, pcp.PCP1y}
	order := geckoTypes.OrderTypeObject.MarketCapDesc
	market, err := g.Client.CoinsMarket(vsCurrency, ids, order, perPage, page, sparkline, priceChangePercentage)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Total coins: ", len(*market))
	fmt.Println(*market)
}
