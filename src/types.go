package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dfuse-io/solana-go/rpc"
	"github.com/go-pg/pg/v10"
)

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

type Balance struct {
	Lamports uint64
	Context  uint64
}

func (b Balance) String() string {
	return fmt.Sprintf("block %d: %f", b.Context, float32(b.Lamports)/float32(1000000000))
}
func (b Balance) FancyString() string {
	return fmt.Sprintf("###################### ACCOUNT BALANCE ######################\nblock %d: %f\n#############################################################", b.Context, float32(b.Lamports)/float32(1000000000))
}

type WSRequest struct {
	Version string        `json:"jsonrpc"`
	ID      uint64        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,omitempty"`
}
type SubResponse struct {
	Version string `json:"jsonrpc"`
	Result  uint64 `json:"result"`
	ID      uint64 `json:"id"`
}
type WSResponse struct {
	Version string           `json:"jsonrpc"`
	Params  *params          `json:"params"`
	Error   *json.RawMessage `json:"error"`
}
type params struct {
	Result       *json.RawMessage `json:"result"`
	Subscription int              `json:"subscription"`
}

type AccountResult struct {
	Context struct {
		Slot uint64
	} `json:"context"`
	Value struct {
		Account rpc.Account `json:"account"`
	} `json:"value"`
}

type ReqParams struct {
	Encoding   string `json:"encoding"`
	Commitment string `json:"commitment"`
}
