package config

import (
	"fmt"
	"log"
	"syscall"

	"github.com/spf13/viper"
	"github.com/thee-engineer/cryptor/crypt/aes"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	log.Println("config init ...")
	if err := readConfig("config"); err != nil {
		log.Fatal(err)
	}
	log.Println("config init finished")
}

func readConfig(name string) error {
	viper.SetConfigType("json")
	viper.SetConfigName(name)

	viper.AddConfigPath("./")
	viper.AddConfigPath("~/.cryptor/")
	viper.AddConfigPath("$GOPATH/src/github.com/thee-engineer/cryptor")

	if err := viper.ReadInConfig(); err != nil {
		// log.Println("failed reading", viper.ConfigFileUsed(), err)

		// Read password for decryption
		fmt.Print(">>> ")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Println("failed to read password, using test password")
			password = []byte("testpassword")
		}

		// Attempt config file decryption
		if err := aes.DecryptFiles(string(password),
			viper.ConfigFileUsed()); err != nil {
			log.Println("failed to decrypt config file")
			return err
		}

		// Attempt reading file after decryption
		if err := viper.ReadInConfig(); err != nil {
			log.Println("failed to read config file after decryption")
			return err
		}
	}

	// Set default values where possible
	viper.SetDefault("version", "0.1")
	viper.SetDefault("node.address", "localhost")
	viper.SetDefault("node.port", 2000)

	return nil
}
