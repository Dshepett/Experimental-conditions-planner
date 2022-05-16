package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port      string
	Host      string
	APIAddres string
}

func New(path string) *Config {
	viper.SetConfigFile(path)
	viper.ReadInConfig()
	return &Config{
		Port:      fmt.Sprintf("%v", viper.Get("PORT")),
		Host:      fmt.Sprintf("%v", viper.Get("HOST")),
		APIAddres: fmt.Sprintf("%v", viper.Get("APIADDRESS")),
	}
}
