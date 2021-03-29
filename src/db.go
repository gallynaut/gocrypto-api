package main

import (
	"fmt"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func (a *App) InitializeDB(hostname, user, password, dbName string) {

	a.DB = pg.Connect(&pg.Options{
		Addr:     hostname,
		User:     user,
		Password: password,
		Database: dbName,
	})
	models := []interface{}{
		(*Exchange)(nil),
	}
	log.Println("creating DB schemas")

	for _, model := range models {
		err := a.DB.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			fmt.Println("error creating table: ", err)
			panic(err)
		}
	}

}
