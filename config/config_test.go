package config_test

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/viper"

	_ "github.com/thee-engineer/cryptor/config"
)

var configFile = path.Join(
	os.Getenv("GOPATH"), "src/github.com/thee-engineer/cryptor/config.json")

func TestBasicConfig(t *testing.T) {
	t.Parallel()

	if usedFile := viper.ConfigFileUsed(); usedFile != configFile {
		t.Fatalf("expected config file %s, got %s", configFile, usedFile)
	}

	if len(viper.AllKeys()) == 0 {
		t.Fatalf("read config, empty")
	}
}
