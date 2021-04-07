package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dfuse-io/solana-go"
	"github.com/dfuse-io/solana-go/rpc"
	"github.com/dfuse-io/solana-go/rpc/ws"

	"github.com/mr-tron/base58"
)

type SolanaApp struct {
	Account   *solana.Account
	PublicKey solana.PublicKey
	RPC       *rpc.Client
	Balance   *solBalance
	Network   *solNetwork
	WS        *ws.Client
}
type solNetwork struct {
	URL    string
	Prefix string
}
type solBalance struct {
	Lamports uint64
	Context  uint64
}

func (sol *SolanaApp) InitializeSolana(network string) {
	if network == "mainnet" {
		sol.Network = &solNetwork{
			URL:    "api.mainnet-beta.solana.com",
			Prefix: network,
		}
	} else if network == "testnet" {
		sol.Network = &solNetwork{
			URL:    "testnet.solana.com",
			Prefix: network,
		}
	} else {
		sol.Network = &solNetwork{
			URL:    "devnet.solana.com",
			Prefix: network,
		}
	}
	var err error
	// setup rpc and web sockets
	sol.RPC = rpc.NewClient("https://" + sol.Network.URL)
	sol.WS, err = ws.Dial(context.Background(), "ws://"+sol.Network.URL)
	if err != nil {
		log.Fatal("could not start Solana websocket:", err)
	}

	// _, err = sol.requestAccountAirdrop(1000000000)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func (sol *SolanaApp) GetSolanaAccount(privkey []byte) {
	var privkeyStr string = base58.Encode(privkey)
	var err error

	sol.Account, err = solana.AccountFromPrivateKeyBase58(privkeyStr)
	if err != nil {
		fmt.Println("error generating key: ", err)
	}

	sol.PublicKey = sol.Account.PublicKey()

	log.Printf("pubkey: https://explorer.solana.com/address/%s?cluster=%s\n", sol.PublicKey, sol.Network.Prefix)

}

func (sol *SolanaApp) subscribeAccount(done <-chan struct{}) error {
	s, err := sol.WS.AccountSubscribe(sol.PublicKey, "")
	if err != nil {
		panic(err)
	}
	log.Println("subscribed to account")

	for {
		select {
		case <-done:
			log.Println("unsubscribing account")
			s.Unsubscribe()
			return nil
		default:
			message, err := s.Recv()
			if err != nil {
				fmt.Println(" error receiving subscription message: ", err)
			}
			acctResult, ok := message.(*ws.AccountResult)
			if !ok {
				log.Printf("error decoding msg: %+v\n", message)
			}
			log.Printf("msg received: %+v\n", *acctResult)

			// unmarshalling message not working so using channel as a trigger
			// sol.getAccountBalance()
		}
	}
}

func (sol *SolanaApp) getAccountBalance() (*solBalance, error) {
	acct, err := sol.RPC.GetBalance(context.Background(), fmt.Sprint(sol.PublicKey), "")
	if err != nil {
		log.Println("error getting acct balance: ", err)
	}
	sol.Balance = &solBalance{
		Lamports: (uint64)(acct.Value),
		Context:  (uint64)(acct.Context.Slot),
	}

	log.Println(sol.Balance.String())
	return sol.Balance, nil
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

// func (sol *SolanaApp) pollRPCAccount(pollRate uint, done <-chan struct{}) {
// 	if pollRate == 0 {
// 		pollRate = 30
// 	}
// 	log.Printf("polling account every %d seconds\n", pollRate)
// 	go func(t *time.Ticker) {
// 		for {
// 			select {
// 			case <-done:
// 				log.Println("polling stopped")
// 				t.Stop()
// 				return
// 			case <-t.C:
// 				go sol.getAccountBalance()
// 			}
// 		}
// 	}(time.NewTicker(time.Duration(pollRate) * time.Second))
// }

func (b solBalance) String() string {
	return fmt.Sprintf("block %d: %f", b.Context, float32(b.Lamports)/float32(1000000000))
}
func (b solBalance) FancyString() string {
	return fmt.Sprintf("###################### ACCOUNT BALANCE ######################\nblock %d: %f\n#############################################################", b.Context, float32(b.Lamports)/float32(1000000000))
}
