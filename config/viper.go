package config

import (
	"github.com/spf13/viper"
	"log"
)

func ConfigInit() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}
}
