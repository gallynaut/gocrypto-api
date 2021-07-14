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
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/test", TestHandler).Methods("GET")
	a.Router.HandleFunc("/candles", TestHandler).Methods("GET")
	// a.Router.HandleFunc("/exchanges", a.GetExchangeHandler).Methods("GET")
	// a.Router.HandleFunc("/exchange", a.AddExchangeHandler).Methods("PUT")
	a.Router.HandleFunc("/airdrop", a.RequestAirdropHandler).Methods("GET")
	a.Router.HandleFunc("/candles/ftx/{symbol}/{resolution}", a.GetSymbolCandlesHandler).Methods("GET").Queries("start", "{start}", "end", "{end}")
	a.Router.HandleFunc("/solana/balance", a.GetSolanaAccountBalanceHandler).Methods("GET")
	a.Router.HandleFunc("/google/{keyword}", a.GetGoogleSearchTrends).Methods("GET")
	http.Handle("/", a.Router)
	a.Router.Use(mux.CORSMethodMiddleware(a.Router))
	a.Router.Use(loggingMiddleware)
}

func (a *App) RunAPI(done <-chan struct{}) (err error) {
	if a.cfg.API.Port == 0 {
		a.cfg.API.Port = 8000
	}
	addr := fmt.Sprintf("localhost:%d", a.cfg.API.Port)

	srv := &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.Router, // Pass our instance of gorilla/mux in.
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

	if code == http.StatusOK {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(response)
	} else {
		w.WriteHeader(code)
	}	
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
