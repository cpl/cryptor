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
	viper.AddConfigPath("/etc/cryptor/")
	viper.AddConfigPath("$GOPATH/src/github.com/thee-engineer/cryptor")

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	log.Println("config init finished")
}
