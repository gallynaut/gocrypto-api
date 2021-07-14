package store

import (
	"fmt"
	"time"

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

type OHLC struct {
	ExchangeSymbol string `pg:"exchange_symbol"`
	PairSymbol     string `pg:"pair_symbol"`
	Period         int32 `pg:"period"`
	CloseTime      time.Time `pg:"close_time"`
	VolumeBase     float64 `pg:"volume_base"`
	VolumeQuote    float64 `pg:"volume_quote"`
	Open           float64 `pg:"open"`
	High           float64 `pg:"high"`
	Low            float64 `pg:"low"`
	Close          float64 `pg:"close"`
}
func (o OHLC) checkCandleExist(db *pg.DB) (bool, error) {
	// err := db.Model(&o).
  //   Column("exchange_symbol", "pair_symbol", "period", "close_time").
  //   Where("id = ?", 1).
  //   Select()
	var candle OHLC
	_, err := db.QueryOne(&candle, `SELECT * FROM ohlcs WHERE exchange_symbol = ?, pair_symbol = ?, period = ?, close_time = ?`, o.ExchangeSymbol, o.PairSymbol, o.Period, o.CloseTime)
	if err != nil {
		
		return true,err
	}
fmt.Println("CANDLE:", candle)
	// no candle found
	if candle.Close == 0 {
		return false, nil
	}
	return true, err
}

func (o OHLC) WriteToDB(db *pg.DB) error {
	b, err := o.checkCandleExist(db)
	if err != nil {
		
		return nil
	}
	if b {
		fmt.Println("STORE: candle already exist")
		return nil
	}
	
	_, err = db.Model(&o).OnConflict("DO NOTHING").Insert()
	if err != nil {
		fmt.Println("STORE: error writing", err)
		return err
	}
	return err
}

// func getExchanges(db *pg.DB) (fetchedExchanges []Exchange, err error) {
// 	err = db.Model(&fetchedExchanges).Select()
// 	if err != nil {
// 		log.Println("error fetching exchanges", err)
// 	}
// 	return fetchedExchanges, err
// }

// func (e *Exchange) addExchange(db *pg.DB) (err error) {
// 	_, err = db.Model(e).Insert()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return
// }
