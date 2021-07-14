package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"code.cryptowat.ch/cw-sdk-go/common"
	"github.com/gallynaut/gocrypto-api/store"
	"github.com/juju/errors"
)

type CWApp struct {
	// Client *cwRest.CWRESTClient
	URL    string
	ApiKey string
}
type ohlcServer struct {
	Result    map[string][][]json.Number `json:"result"`
	Allowance struct {
		Cost         float64 `json:"cost"`
		Remaining    float64 `json:"remaining"`
		RemaningPaid int64   `json:"remainingPaid"`
		Account      string  `json:"account"`
	} `json:"allowance"`
}

func (a *App) initializeCW(apiKey string) {
	// a.CW.Client = cwRest.NewCWRESTClient(nil)
	// a.CW.URL = "https://api.cryptowat.ch"
	// a.CW.ApiKey = apiKey
	log.Println("CW: initialized")
	a.Router.HandleFunc("/cw/markets/{exchange}/{pair}/ohlc", a.GetOHLCHandler).Methods("GET")
	// a.Router.HandleFunc("/gecko/coin/{symbol}", a.Gecko.GetCoinHandler).Methods("GET")
	// a.Router.HandleFunc("/gecko/{symbol}/price", a.Gecko.GetSymbolPriceHandler).Methods("GET")

	// go func() {
	// 	_, err := a.CW.GetOHLC(&CWOHLCRequest{
	// 		ExchangeSymbol: "binance",
	// 		PairSymbol:     "solusdt",
	// 		Periods:        []string{"60", "3600", "86400"},
	// 		After:          "1577862000", // time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC),
	// 	})
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }()

	// add routes for cryptowatch here
}

// func (cw *CWApp) Run(pollRate uint, done <-chan struct{}) {
// 	if pollRate == 0 {
// 		pollRate = 30
// 	}
// 	log.Printf("NOM: polling sparkline every %d seconds\n", pollRate)
// 	go func(t *time.Ticker) {
// 		for {
// 			select {
// 			case <-done:
// 				log.Println("CW: polling candles stopped")
// 				t.Stop()
// 				return
// 			case <-t.C:
// 				return
// 			}
// 		}
// 	}(time.NewTicker(time.Duration(pollRate) * time.Second))
// }

type CWOHLCRequest struct {
	ExchangeSymbol string   // ftx
	PairSymbol     string   // solusd
	Periods        []string // 60 (1min)
	Before         string
	After          string
}

func (o *CWOHLCRequest) String() string {
	periodStr := strings.Join(o.Periods, ",")
	if o.After == "" {
		return fmt.Sprintf("/markets/%s/%s/ohlc?periods=%s", o.ExchangeSymbol, o.PairSymbol, periodStr)
	} else {
		return fmt.Sprintf("/markets/%s/%s/ohlc?periods=%s&after=%s", o.ExchangeSymbol, o.PairSymbol, periodStr, o.After)
	}

}
func jsonToFloat(num json.Number) float64 {
	f, err := num.Float64()
	if err != nil {
		log.Printf("error converting %v to a float\n", num)
		return 0.0
	}
	return f
}

func (a *App) GetOHLC(req *CWOHLCRequest) (map[common.Period][]store.OHLC, error) {
	// it would be better to set header manually
	url := a.cfg.CW.URL + req.String() + "&apikey=" + a.cfg.CW.PublicKey
	log.Println("CW: ", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Trace(err)
	}

	defer resp.Body.Close()

	srv := ohlcServer{}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&srv); err != nil {
		return nil, errors.Trace(err)
	}
	log.Printf("CW: request cost %f with (%f) remaining credits\n", srv.Allowance.Cost, srv.Allowance.Remaining)

	ret := make(map[common.Period][]store.OHLC, len(srv.Result))
	for srvPeriod, srvCandles := range srv.Result {
		v64, err := strconv.ParseInt(srvPeriod, 10, 64)
		if err != nil {
			log.Println("CW: cant parse ", srvPeriod)
			continue
		}

		period := common.Period(v64)

		candles := make([]store.OHLC, 0, len(srvCandles))
		for _, srvCandle := range srvCandles {
			if len(srvCandle) < 7 {
				return nil, errors.Errorf("unexpected response from the server: wanted 7 elements, got %v", srvCandle)
			}

			ts, err := srvCandle[0].Int64()
			if err != nil {
				return nil, errors.Annotatef(err, "getting timestamp %q", srvCandle[0].String())
			}
			candle := store.OHLC{
				ExchangeSymbol: req.ExchangeSymbol,
				PairSymbol:     req.PairSymbol,
				Period:         int32(period),
				CloseTime:      time.Unix(ts, 0),
				Open:           jsonToFloat(srvCandle[1]),
				High:           jsonToFloat(srvCandle[2]),
				Low:            jsonToFloat(srvCandle[3]),
				Close:          jsonToFloat(srvCandle[4]),
				VolumeBase:     jsonToFloat(srvCandle[5]),
				VolumeQuote:    jsonToFloat(srvCandle[6]),
			}

			candles = append(candles, candle)
		}
		ret[period] = candles
	}
	go a.writeResultsToDB(ret)
	return ret, nil
}

func (a *App) writeResultsToDB(results map[common.Period][]store.OHLC) error {
	minCandle := results[common.Period1M]

	_, err := a.Store.DB.Model(&minCandle).Insert()
	if err != nil {
		panic(err)
	}

	// for _, element := range results {
	// 	_, err := a.Store.DB.Model(&element).Insert()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	// fmt.Println(element)
	// }
	return nil
}