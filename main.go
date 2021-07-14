package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/gallynaut/gocrypto-api/store"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type appConfig struct {
	DB struct {
		Hostname string `json:"hostname" mapstructure:"hostname"`
		Username string `json:"username" mapstructure:"username"`
		Password string `json:"password" mapstructure:"password"`
		DBName   string `json:"dbName" mapstructure:"dbName"`
	} `json:"db" mapstructure:"db"`
	API struct {
		Port uint `json:"port" mapstructure:"port"`
	} `json:"api" mapstructure:"api"`
	Solana struct {
		PrivateKey      []byte `json:"privKey" mapstructure:"privKey"`
		Network         string `json:"network" mapstructure:"network"`
		AccountPollRate uint   `json:"accountPollRate" mapstructure:"accountPollRate"`
	} `json:"solana" mapstructure:"solana"`
	FTX struct {
		ApiKey string `json:"apiKey" mapstructure:"apiKey"`
		Secret string `json:"secret" mapstructure:"secret"`
	} `json:"ftx" mapstructure:"ftx"`
	CW struct {
		PublicKey string `json:"pubKey" mapstructure:"pubKey"`
		Secret    string `json:"secret" mapstructure:"secret"`
		URL       string `json:"url" mapstructure:"url"`
	} `json:"cryptoWatch" mapstructure:"cryptoWatch"`
}

type App struct {
	Router *mux.Router
	Store  store.StoreApp
	Solana SolanaApp
	FTX    FTXApp
	Gecko  GeckoApp
	CW     CWApp
	cfg    appConfig
	ctx    context.Context
}

func main() {
	a := App{}

	var cancel context.CancelFunc
	a.ctx, cancel = context.WithCancel(context.Background())

	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt)
		oscall := <-done
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	var err error
	a.cfg, err = loadConfig("config.yml")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	a.initialize(a.cfg)

	// go a.getGoogleTrends("solana")

	// go a.CW.getCandles()

	// start api
	a.RunAPI(a.ctx.Done())

	log.Println("shutting down")
	os.Exit(0)

}

func loadConfig(path string) (config appConfig, err error) {
	viper.SetConfigFile(path)
	// viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func (a *App) initialize(config appConfig) {
	// database and api routes
	// a.Store = store.InitializeDB(config.DB.Hostname, config.DB.Username, config.DB.Password, config.DB.DBName)
	a.initializeRoutes()

	// intitialize data feeds
	// a.initializeCW(config.CW.PublicKey)
	a.initializeGecko()
	// a.initializeFTX(config.FTX.ApiKey, config.FTX.Secret)

	// solana keys, rpc, and websocket
	// a.initializeSolana(config.Solana.Network, config.Solana.PrivateKey)
}
