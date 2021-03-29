package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHostname string `mapstructure:"APP_DB_HOSTNAME"`
	DBUsername string `mapstructure:"APP_DB_USERNAME"`
	DBPassword string `mapstructure:"APP_DB_PASSWORD"`
	DBName     string `mapstructure:"APP_DB_NAME"`
	APIPort    uint   `mapstructure:"APP_API_PORT"`
}

func main() {
	a := App{}

	config, err := LoadConfig(".env")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	a.InitializeDB(config.DBHostname, config.DBUsername, config.DBPassword, config.DBName)
	a.InitializeRoutes()

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
