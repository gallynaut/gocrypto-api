package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func (a *App) GetSymbolCandlesHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	symbol := params["symbol"]
	resolution, err := strconv.Atoi(params["resolution"])
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, "resolution not found")
	}
	start, err := strconv.ParseInt(params["start"], 10, 64)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, "start not found")
	}
	end, err := strconv.ParseInt(params["end"], 10, 64)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, "end not found")
	}
	fmt.Printf("GET params were: %s\nsymbol: %s\nresolution: %d\n", r.URL.Query(), symbol, resolution)
	err = a.FTX.getCandles(symbol, resolution, start, end)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, "candles not found")
	}

	// Exchanges, err := getExchanges(a.DB)
	// if err != nil {
	// 	fmt.Println("error getting exchanges: ", err)
	// }
	// log.Println(Exchanges)

	respondWithJSON(w, http.StatusOK, params)
}

func (a *App) GetSolanaAccountBalance(w http.ResponseWriter, r *http.Request) {
	b, err := a.Solana.getAccountBalance()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting acct balance: %s", err))
	}
	respondWithJSON(w, http.StatusOK, b)
}
