package main

import (
	"log"
	"time"

	cwRest "code.cryptowat.ch/cw-sdk-go/client/rest"
)

type CWApp struct {
	Client *cwRest.CWRESTClient
}

func (cw *CWApp) initializeCW(apiKey string) {
	cw.Client = cwRest.NewCWRESTClient(nil)
	log.Println("CW: initialized")
}
func (cw *CWApp) Run(pollRate uint, done <-chan struct{}) {
	if pollRate == 0 {
		pollRate = 30
	}
	log.Printf("NOM: polling sparkline every %d seconds\n", pollRate)
	go func(t *time.Ticker) {
		for {
			select {
			case <-done:
				log.Println("CW: polling candles stopped")
				t.Stop()
				return
			case <-t.C:
				return
			}
		}
	}(time.NewTicker(time.Duration(pollRate) * time.Second))
}

func (cw *CWApp) getCandles(symbols []string) {

}
