package main

import (
	"fmt"
	"log"

	"github.com/go-pg/pg/v10"
)

type Symbol struct {
	ID          string  `json:"id" pg:"id,pk"`
	Symbol      string  `json:"symbol" pg:"symbol"`
	Name        string  `json:"name" pg:"name"`
	TotalSupply float64 `json:"total_supply" pg:"total_supply"`
	LastUpdated string  `json:"last_updated" pg:"last_updated"`
}

type Candles struct {
	ID       string `json:"id" pg:"id,pk"`
	Symbol   string `json:"symbol" pg:"symbol"`
	Exchange string `json:"exchange" pg:"exchange"`
	Period   string `json:"name" pg:"name"`

	TotalSupply float64 `json:"total_supply" pg:"total_supply"`
	LastUpdated string  `json:"last_updated" pg:"last_updated"`
}

type Exchange struct {
	FullName     string `pg:"full_name,pk" json:"fullName"`      //Name + Network (i.e. BYBIT-MAIN)
	ExchangeName string `pg:"exchange_name,notnull" json:"name"` //BYBIT
	EndPoint     string `pg:"end_point,notnull" json:"endPoint"`
}

func (e Exchange) String() string {
	return fmt.Sprintf("%s @ %s", e.FullName, e.EndPoint)
}

func getExchanges(db *pg.DB) (fetchedExchanges []Exchange, err error) {
	err = db.Model(&fetchedExchanges).Select()
	if err != nil {
		log.Println("error fetching exchanges", err)
	}
	return fetchedExchanges, err
}

func (e *Exchange) addExchange(db *pg.DB) (err error) {
	_, err = db.Model(e).Insert()
	if err != nil {
		panic(err)
	}
	return
}
