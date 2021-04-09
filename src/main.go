package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/viper"
)

type Config struct {
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
	} `json:"cryptoWatch" mapstructure:"cryptoWatch"`
}
type App struct {
	API    APIApp
	Store  StoreApp
	Solana SolanaApp
	FTX    FTXApp
	Gecko  GeckoApp
	CW     CWApp
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

	config, err := LoadConfig("config.yml")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// database and api routes
	a.InitializeDB(config.DB.Hostname, config.DB.Username, config.DB.Password, config.DB.DBName)
	a.InitializeRoutes()
	a.CW.initializeCW(config.CW.PublicKey)

	a.Gecko.initializeGecko()
	// a.CW.initializeCW(config.CW.ApiKey, config.CW.Secret)

	// solana keys, rpc, and websocket
	a.Solana.InitializeSolana(config.Solana.Network)
	a.Solana.GetSolanaAccount(config.Solana.PrivateKey)
	defer a.Solana.WS.Close()
	go a.Solana.requestAccountAirdrop(1000000000)

	// connect ftx account
	a.FTX.initializeFTX(config.FTX.ApiKey, config.FTX.Secret)

	// poll solana account balances and wait for blocks
	go a.Solana.subscribeAccount(a.ctx.Done())

	// poll FTX funding rates
	go a.FTX.pollFundingRates(45, a.ctx.Done())

	// start api
	a.API.Run(config.API.Port, a.ctx.Done())

	log.Println("shutting down")
	os.Exit(0)

}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path)
	// viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
