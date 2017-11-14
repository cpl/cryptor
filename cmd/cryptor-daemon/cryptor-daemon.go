package main

import (
	"log"
	"path"

	"github.com/sevlyar/go-daemon"
	"github.com/thee-engineer/cryptor/cachedb"
)

const (
	port = 9990
	name = "cryptord"
)

func main() {
	cntxt := &daemon.Context{
		PidFileName: path.Join(cachedb.GetCryptorDir(), "cryptord.pid"),
		PidFilePerm: 0644,
		LogFileName: path.Join(cachedb.GetCryptorDir(), "cryptord.log"),
		LogFilePerm: 0640,
		WorkDir:     cachedb.GetCryptorDir(),
		Umask:       027,
		Args:        []string{name},
	}

	child, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if child != nil {
		return
	}
	defer cntxt.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")

	listen()
}

func listen() {
	for {

	}
}

func parseCommand(command []byte) error {
	return nil
}
