package main

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"

	"log"
	"sort"

	"github.com/go-numb/go-ftx/auth"
	"github.com/go-numb/go-ftx/rest"
	"github.com/go-numb/go-ftx/rest/private/account"
	"github.com/go-numb/go-ftx/rest/public/futures"
	"github.com/go-numb/go-ftx/rest/public/markets"
)

type FTXApp struct {
	Client *rest.Client
}

func (ftx *FTXApp) initializeFTX(apiKey, secret string) {
	// Only main account
	ftx.Client = rest.New(auth.New(apiKey, secret))

	err := ftx.getAccountInformation()
	if err != nil {
		log.Fatal("error getting account information: ", err)
	}

	// go ftx.getCandles("SOL-PERP", 1586143817, 1617679817)
}

func (ftx *FTXApp) pollFundingRates(pollRate uint64, done <-chan struct{}) {
	log.Printf("polling funding rates every %d seconds\n", pollRate)
	go func(t *time.Ticker) {
		for {
			select {
			case <-done:
				log.Printf("funding polling stopped\n")
				t.Stop()
				return
			case <-t.C:
				go ftx.getFundingRates()
				go ftx.getMarket("SOL-PERP")
			}
		}
	}(time.NewTicker(time.Duration(pollRate) * time.Second))
}

func (ftx *FTXApp) getFundingRates() error {
	// FundingRate
	rates, err := ftx.Client.Rates(&futures.RequestForRates{})
	if err != nil {
		return err
	}
	// Sort by FundingRate & Print
	// Custom sort

	sort.Sort(sort.Reverse(rates))
	fmt.Println("===== Funding Rates =====")
	for i, v := range *rates {
		if i < 10 {
			fmt.Printf("%f			%s		%s\n", (v.Rate * 100.0), v.Future, v.Time.String())
		} else {
			break
		}
	}
	return nil
}
func (ftx *FTXApp) getAccountInformation() error {
	info, err := ftx.Client.Information(&account.RequestForInformation{})
	if err != nil {
		// log.Fatal(err)
		return err
	}

	fmt.Printf("account info fetched %s\n", info.Username)
	return nil
}
func (ftx *FTXApp) getMarket(mkt string) error {
	market, err := ftx.Client.Markets(&markets.RequestForMarkets{
		ProductCode: mkt,
	})
	if err != nil {
		return err
	}

	for _, v := range *market {
		if v.Type == "future" {
			fmt.Printf("%s: $%s\n", v.Name, humanize.Commaf(v.VolumeUsd24H))
		}
	}
	return nil
}

func (ftx *FTXApp) getCandles(symbol string, res int, start, end int64) error {
	candle, err := ftx.Client.Candles(&markets.RequestForCandles{
		ProductCode: symbol,
		Resolution:  res,
		Start:       start,
		End:         end,
	})
	if err != nil {
		return err
	}

	for _, v := range *candle {
		fmt.Printf("%s - %f\n", v.StartTime, v.Volume)
	}

	return nil
}
