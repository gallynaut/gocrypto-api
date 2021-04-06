package main

import (
	"fmt"
	"log"
)

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
			_, err := s.Recv()
			if err != nil {
				fmt.Println(" error receiving subscription message: ", err)
			}
			// acctResult, ok := message.(*ws.AccountResult)
			// if !ok {
			// 	log.Printf("error decoding msg: %+v\n", message)
			// }
			// log.Printf("msg received: %+v\n", *acctResult)

			// message unmarshalling not working so using ws as a trigger
			sol.getAccountBalance()
		}
	}
}
