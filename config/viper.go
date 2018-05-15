package config

import (
	"github.com/spf13/viper"
)

// InitViper ...
func InitViper() {
	viper.SetConfigType("json")
	viper.SetConfigName("config")

	viper.AddConfigPath("$GOPATH/src/github.com/thee-engineer/cryptor")
	viper.AddConfigPath("/etc/cryptor/")
	viper.AddConfigPath("~/.cryptor/")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
