package main

import (
	"fmt"

	"github.com/dfuse-io/solana-go"
	"github.com/mr-tron/base58"
)

func (a *App) ConnectWallet(privkey []byte) {
	var privkeyStr string = base58.Encode(privkey)
	var err error

	a.Sol.PrivateKey, err = solana.PrivateKeyFromBase58(privkeyStr)
	if err != nil {
		fmt.Println("error generating key: ", err)
	}
	fmt.Println("new solana account generated: ", a.Sol)
}
