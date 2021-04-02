package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Test")
}
func (a *App) RequestAirdropHandler(w http.ResponseWriter, r *http.Request) {
	url, err := a.Solana.requestAccountAirdrop(1000000000)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "error getting airdrop: %s", err)
	}
	// w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, url)
}

func (a *App) AddExchangeHandler(w http.ResponseWriter, r *http.Request) {
	var e Exchange
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&e); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := e.addExchange(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, e)
}

func (a *App) GetExchangeHandler(w http.ResponseWriter, r *http.Request) {
	Exchanges, err := getExchanges(a.DB)
	if err != nil {
		fmt.Println("error getting exchanges: ", err)
	}
	log.Println(Exchanges)

	respondWithJSON(w, http.StatusOK, Exchanges)
}
