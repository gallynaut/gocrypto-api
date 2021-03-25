package main

import (
	"fmt"
	"strings"
)

type Exchange struct {
	FullName string `sql:"fullName,pk" json:"fullName"` //Name + Network (i.e. BYBIT-MAIN)
	Name     string `sql:"name" json:"name"`            //BYBIT
	EndPoint string `sql:"endPoint" json:"endPoint"`
}
type Exchanges []Exchange

func (e Exchange) String() string {
	return fmt.Sprintf("%s @ %s", e.FullName, e.EndPoint)
}

func (exchanges Exchanges) String() string {
	var eStr []string
	for _, e := range exchanges {
		eStr = append(eStr, e.String())
	}
	return strings.Join(eStr, ", ")
}

type DBSchema struct {
	Name        string
	Table       interface{}
	InitialData string
}

//Clean these up
// type KlineData struct {
// 	exchange   *Exchange `pg:"rel:has-one"`
// 	symbol     string
// 	interval   string
// 	openTime   int64
// 	open       string
// 	closeTime  string
// 	closePrice string
// 	high       string
// 	low        string
// 	volume     string
// 	turnover   string
// }
