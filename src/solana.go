package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dfuse-io/solana-go"
	"github.com/dfuse-io/solana-go/rpc"

	"github.com/mr-tron/base58"
)

type SolanaApp struct {
	Account   *solana.Account
	PublicKey solana.PublicKey
	RPC       *rpc.Client
	Balance   *Balance
}

func (sol *SolanaApp) InitializeSolana() {
	_, err := sol.requestAccountAirdrop(1000000000)
	if err != nil {
		fmt.Println(err)
	}

	go sol.getAccountBalance()
}

func (sol *SolanaApp) GetSolanaAccount(privkey []byte) {
	var privkeyStr string = base58.Encode(privkey)
	var err error

	sol.Account, err = solana.AccountFromPrivateKeyBase58(privkeyStr)
	if err != nil {
		fmt.Println("error generating key: ", err)
	}

	sol.PublicKey = sol.Account.PublicKey()

	fmt.Printf("pubkey: https://explorer.solana.com/address/%s?cluster=devnet\n", sol.PublicKey)
}

func (sol *SolanaApp) getAccountBalance() {
	log.Println("fetching account balance")
	acct, err := sol.RPC.GetBalance(context.Background(), fmt.Sprint(sol.PublicKey), "")
	if err != nil {
		log.Println("error getting acct balance: ", err)
	}
	sol.Balance = &Balance{
		Lamports: (uint64)(acct.Value),
		Context:  (uint64)(acct.Context.Slot),
	}
	fmt.Println(sol.Balance.String())
}

func (sol *SolanaApp) requestAccountAirdrop(lamports uint64) (url string, err error) {
	airdrop, err := sol.RPC.RequestAirdrop(context.Background(), &sol.PublicKey, lamports, rpc.CommitmentMax)
	if err != nil {
		return "", fmt.Errorf("error getting airdrop: %e", err)
	}
	url = fmt.Sprintf("https://explorer.solana.com/tx/%s?cluster=devnet", airdrop)
	fmt.Println("requested airdrop: ", url)
	return url, nil
}

func (sol *SolanaApp) pollAccount(pollRate uint, done <-chan struct{}) {
	if pollRate == 0 {
		pollRate = 30
	}
	log.Printf("polling account every %d seconds\n", pollRate)
	go func(t *time.Ticker) {
		for {
			select {
			case <-done:
				log.Println("polling stopped")
				t.Stop()
				return
			case <-t.C:
				go sol.getAccountBalance()
			}
		}
	}(time.NewTicker(time.Duration(pollRate) * time.Second))
}
