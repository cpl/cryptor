package config

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	log.Println("config init ...")

	viper.SetConfigType("json")
	viper.SetConfigName("config")

	viper.AddConfigPath("./")
	viper.AddConfigPath("~/.cryptor/")
	viper.AddConfigPath("$GOPATH/src/github.com/thee-engineer/cryptor")

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	viper.SetDefault("version", "0.1")
	viper.SetDefault("node.address", "localhost")
	viper.SetDefault("node.port", 2000)

	log.Println("config init finished")
}
