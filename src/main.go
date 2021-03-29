package main

import (
	"log"

	"github.com/dfuse-io/solana-go"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type Config struct {
	DBHostname    string `mapstructure:"APP_DB_HOSTNAME"`
	DBUsername    string `mapstructure:"APP_DB_USERNAME"`
	DBPassword    string `mapstructure:"APP_DB_PASSWORD"`
	DBName        string `mapstructure:"APP_DB_NAME"`
	APIPort       uint   `mapstructure:"APP_API_PORT"`
	SolPrivateKey []byte `mapstructure:"SOLANA_PRIV_KEY"`
}

type App struct {
	Router *mux.Router
	DB     *pg.DB
	Sol    solana.Account
	// Exchanges []Exchange
}

func main() {
	a := App{}

	config, err := LoadConfig(".env")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	a.InitializeDB(config.DBHostname, config.DBUsername, config.DBPassword, config.DBName)
	a.InitializeRoutes()

	a.ConnectWallet(config.SolPrivateKey)

	a.Run(config.APIPort)

}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
