package common

import (
	"log"

	"github.com/spf13/viper"
)

// configuration strcu that hold
// application running port and db connection string
type Config struct {
	Port               string `mapstructure:"port"`
	DbConnectionString string `mapstructure:"dbconnectionstring"`
}

// config varibale to access from other files/package
var AppConfig *Config

// load configuration data using viper
func LoadAppConfig() {
	log.Println("Loading Server Configuration")
	viper.AddConfigPath("./src/common")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatal(err)
	}
}
