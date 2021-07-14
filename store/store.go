package store

import (
	"fmt"
	"log"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type StoreApp struct {
	DB *pg.DB
}

func InitializeDB(hostname, user, password, dbName string) StoreApp {
	var s StoreApp
	s.DB = pg.Connect(&pg.Options{
		Addr:     hostname,
		User:     user,
		Password: password,
		Database: dbName,
	})
	models := []interface{}{
		(*Exchange)(nil),
		(*OHLC)(nil),
	}
	log.Println("STORE: creating DB schemas")

	for _, model := range models {
		err := s.DB.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			fmt.Println("STORE: error creating table: ", err)
			panic(err)
		}
	}
	return s
}

func GetCandle(db *pg.DB, exchangeSymbol string, pairSymbol string, period string, closeTime time.Time ) (*OHLC, error) {
	var candle OHLC
	_, err := db.QueryOne(&candle, `SELECT * FROM ohlcs WHERE exchange_symbol = ?, pair_symbol = ?, period = ?, close_time = ?`, exchangeSymbol, pairSymbol, period, closeTime)
	return &candle, err
}
