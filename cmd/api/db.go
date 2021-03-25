package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var seedData = []DBSchema{
	DBSchema{
		Name:        "Exchanges",
		Table:       Exchange{},
		InitialData: "./data/exchanges.json",
	},
	// DBSchema{
	// 	Name:        "KlineData",
	// 	Table:       Exchange,
	// 	InitialData: "",
	// },
}

func databaseInit(c *context.Context) {
	pgdb := connect()
	defer pgdb.Close()

	err := createSchema(pgdb)
	if err != nil {
		panic(err)
	}

	// Write exchanges.json to database
	var initialExchanges Exchanges
	initialExchanges, err = getJsonExchanges()
	if err != nil {
		panic(err)
	}
	err = writeExchanges(pgdb, initialExchanges)
	if err != nil {
		panic(err)
	}

	// Check database has exchanges
	var fetchedExchanges Exchanges
	err = pgdb.Model(&fetchedExchanges).Limit(10).Select()
	if err != nil {
		panic(err)
	}

}

func connect() *pg.DB {
	return pg.Connect(&pg.Options{
		Addr:     "gallydb:5432",
		User:     "me",
		Password: "password",
		Database: "go_testing",
	})
}

// createSchema creates database schema for User and Story models.
func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*Exchange)(nil),
		// (*Story)(nil),
	}

	for _, model := range models {
		err := db.Model(model).DropTable(&orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("DROP:   %s\n", reflect.TypeOf(model))

		err = db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
		fmt.Printf("CREATE: %s\n", reflect.TypeOf(model))
	}
	return nil
}

func writeExchanges(db *pg.DB, exchanges Exchanges) error {
	_, err := db.Model(&exchanges).Insert()
	if err != nil {
		panic(err)
	}
	return nil
}
