package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/dfuse-io/solana-go/rpc"

	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type Config struct {
	DB     DBConfig     `json:"db" mapstructure:"db"`
	API    APIConfig    `json:"api" mapstructure:"api"`
	Solana SolanaConfig `json:"solana" mapstructure:"solana"`
}
type DBConfig struct {
	Hostname string `json:"hostname" mapstructure:"hostname"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	DBName   string `json:"dbName" mapstructure:"dbName"`
}
type APIConfig struct {
	Port uint `json:"port" mapstructure:"port"`
}
type SolanaConfig struct {
	PrivateKey      []byte `json:"privKey" mapstructure:"privKey"`
	Network         string `json:"network" mapstructure:"network"`
	AccountPollRate uint   `json:"accountPollRate" mapstructure:"accountPollRate"`
}
type App struct {
	Router *mux.Router
	DB     *pg.DB
	Solana SolanaApp
	ctx    context.Context
	done   chan os.Signal
}

func main() {
	a := App{}
	a.done = make(chan os.Signal, 1)
	signal.Notify(a.done, os.Interrupt)

	var cancel context.CancelFunc
	a.ctx, cancel = context.WithCancel(context.Background())

	go func() {
		oscall := <-a.done
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	config, err := LoadConfig("config.yml")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	a.InitializeDB(config.DB.Hostname, config.DB.Username, config.DB.Password, config.DB.DBName)
	a.InitializeRoutes()

	// setup rpc and web sockets
	a.Solana.RPC = rpc.NewClient("https://devnet.solana.com")

	// get public key and request airdrop
	a.Solana.GetSolanaAccount(config.Solana.PrivateKey)
	a.Solana.InitializeSolana()

	// poll account balance
	go a.Solana.pollAccount(config.Solana.AccountPollRate, a.ctx.Done())

	// start api
	a.Run(config.API.Port)

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
