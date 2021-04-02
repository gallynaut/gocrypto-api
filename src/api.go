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

func (a *App) InitializeRoutes() {
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/test", TestHandler).Methods("GET")
	a.Router.HandleFunc("/exchanges", a.GetExchangeHandler).Methods("GET")
	a.Router.HandleFunc("/exchange", a.AddExchangeHandler).Methods("PUT")
	a.Router.HandleFunc("/airdrop", a.RequestAirdropHandler).Methods("GET")
	http.Handle("/", a.Router)
	a.Router.Use(loggingMiddleware)
}

func (a *App) Run(port uint) (err error) {
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
		Handler:      a.Router, // Pass our instance of gorilla/mux in.
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	log.Printf("server started at %s\n", addr)
	<-a.ctx.Done()
	log.Printf("server stopped\n")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s\n", err)
	}

	log.Printf("server exited properly\n")
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
