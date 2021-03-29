package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *pg.DB
}

func (a *App) InitializeDB(hostname, user, password, dbName string) {

	a.DB = pg.Connect(&pg.Options{
		Addr:     hostname,
		User:     user,
		Password: password,
		Database: dbName,
	})
	models := []interface{}{
		(*Exchange)(nil),
	}
	log.Println("creating DB schemas")

	for _, model := range models {
		err := a.DB.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			fmt.Println("error creating table: ", err)
			panic(err)
		}
	}

}

func (a *App) InitializeRoutes() {
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/test", TestHandler).Methods("GET")
	a.Router.HandleFunc("/exchanges", a.GetExchangeHandler).Methods("GET")
	a.Router.HandleFunc("/exchange", a.AddExchangeHandler).Methods("PUT")
	http.Handle("/", a.Router)
	a.Router.Use(loggingMiddleware)
}

func (a *App) Run(port uint) {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Add your routes as needed
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

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		fmt.Printf("Running server @ %s\n", addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
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
