package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type APIApp struct {
	Router *mux.Router
	Port   uint
}

// mostly used for debugging right now
func (a *App) initializeRoutes() {
	a.API.Router = mux.NewRouter()
	a.API.Router.HandleFunc("/test", TestHandler).Methods("GET")
	// a.API.Router.HandleFunc("/exchanges", a.GetExchangeHandler).Methods("GET")
	// a.API.Router.HandleFunc("/exchange", a.AddExchangeHandler).Methods("PUT")
	a.API.Router.HandleFunc("/airdrop", a.RequestAirdropHandler).Methods("GET")
	a.API.Router.HandleFunc("/candles/ftx/{symbol}/{resolution}", a.GetSymbolCandlesHandler).Methods("GET").Queries("start", "{start}", "end", "{end}")
	a.API.Router.HandleFunc("/solana/balance", a.GetSolanaAccountBalanceHandler).Methods("GET")
	// a.API.Router.HandleFunc("/gecko/{symbol}", a.GetSymbolHandler).Methods("GET")
	// a.API.Router.HandleFunc("/gecko/coin/{symbol}", a.GetCoinHandler).Methods("GET")
	// a.API.Router.HandleFunc("/gecko/{symbol}/price", a.GetSymbolPriceHandler).Methods("GET")
	http.Handle("/", a.API.Router)
	a.API.Router.Use(loggingMiddleware)
}

func (api *APIApp) Run(port uint, done <-chan struct{}) (err error) {
	if port == 0 {
		port = 8000
	}
	addr := fmt.Sprintf("localhost:%d", port)

	srv := &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      api.Router, // Pass our instance of gorilla/mux in.
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	log.Printf("API: server started at %s\n", addr)
	<-done
	log.Printf("API: server stopped\n")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("API: server Shutdown Failed:%+s\n", err)
	}

	log.Printf("API: server exited properly\n")
	if err == http.ErrServerClosed {
		err = nil
	}

	return
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
