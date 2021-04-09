package main

import (
	"fmt"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type StoreApp struct {
	DB *pg.DB
}

func (a *App) InitializeDB(hostname, user, password, dbName string) {

	a.Store.DB = pg.Connect(&pg.Options{
		Addr:     hostname,
		User:     user,
		Password: password,
		Database: dbName,
	})
	models := []interface{}{
		(*Exchange)(nil),
	}
	log.Println("STORE: creating DB schemas")

	for _, model := range models {
		err := a.Store.DB.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			fmt.Println("STORE: error creating table: ", err)
			panic(err)
		}
	}
}
