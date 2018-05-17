package config

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/viper"
	"github.com/thee-engineer/cryptor/crypt/aes"
)

var configFile = path.Join(
	os.Getenv("GOPATH"),
	"src/github.com/thee-engineer/cryptor/config_test.json")

func TestConfigNotFound(t *testing.T) {
	viper.Reset()

	if err := readConfig("no-such-config"); err != nil {
		if err.Error() != "open : no such file or directory" {
			t.Fatal(err)
		}
		return
	}

	t.Fatalf("found invalid config")
}
func TestBasicConfig(t *testing.T) {
	viper.Reset()
	readConfig("config_test")
	viper.Reset()

	if usedFile := viper.ConfigFileUsed(); usedFile != "" {
		t.Fatalf("read config file should be none")
	}

	if len(viper.AllKeys()) != 0 {
		t.Fatalf("config, not empty")
	}

	if err := aes.EncryptFiles("testpassword", configFile); err != nil {
		t.Fatal(err)
	}
	if err := readConfig("config_test"); err != nil {
		t.Fatal(err)
	}

	if usedFile := viper.ConfigFileUsed(); usedFile != configFile {
		t.Fatalf("expected config file %s, got %s", configFile, usedFile)
	}

	if len(viper.AllKeys()) == 0 {
		t.Fatalf("read config, empty")
	}

	if !viper.IsSet("version") {
		t.Fatalf("failed to find field version")
	}
}
