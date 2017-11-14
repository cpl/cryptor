// cryptor
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

const cryptorName = "cryptor"
const cryptorVersion = "0.0.3"
const cryptorCopyright = "Copyright (C) 2017 Alexandru-Paul Copil"
const cryptorUsageDescription = "command line interface for basic tasks"
const cryptorUsageString = " [command] [sub-command] [args]"

func main() {
	app := cli.NewApp()

	// CLI misc settings
	app.EnableBashCompletion = true
	app.Version = cryptorVersion
	app.Name = cryptorName
	app.Compiled = time.Now()
	app.Usage = cryptorUsageDescription
	app.UsageText = cryptorName + cryptorUsageString
	app.Copyright = cryptorCopyright

	// CLI variables
	var key, password, cache, chunk, source string
	var size uint

	app.Commands = []cli.Command{
		{
			// chunk commands
			Name:  "chunk",
			Usage: "chunk a file given a input package and output cache",
			Action: func(c *cli.Context) error {
				if key == "" && password == "" {
					fmt.Printf("provide [--key | -k] or [--password | -p]\n")
					return nil
				}

				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "key, k",
					Usage:       "provide a hex encoded key for encryption",
					Destination: &key},
				cli.StringFlag{
					Name:        "password, p",
					Usage:       "provide a password for encryption",
					Destination: &password},
				cli.StringFlag{
					Name:        "cache, c, out, o",
					Usage:       "provide the path/name of a cache as destination",
					Value:       "default",
					Destination: &cache},
				cli.StringFlag{
					Name:        "source, in, i",
					Usage:       "provide the path for input file(s)/dir(s)",
					Value:       ".",
					Destination: &source},
				cli.UintFlag{
					Name:        "size, s",
					Usage:       "provide the size of each individual chunk",
					Value:       0,
					Destination: &size},
			},
		},
		{
			// assemble commands
			Name:  "assemble",
			Usage: "assemble a package given the key/pass and input cache",
			Action: func(c *cli.Context) error {
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "key, k",
					Usage:       "provide a hex encoded key for decryption",
					Destination: &key},
				cli.StringFlag{
					Name:        "password, p",
					Usage:       "provide a password for decryption",
					Destination: &password},
				cli.StringFlag{
					Name:        "cache, c",
					Usage:       "provide the path/name of a cache as source",
					Value:       "default",
					Destination: &cache},
				cli.StringFlag{
					Name:        "chunk, tail, hash",
					Usage:       "provide the hash of the tail chunk",
					Destination: &chunk},
			},
		},
		{
			// cache commands
			Name:  "cache",
			Usage: "diffrent commands for interacting with caches",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new cache",
					Action: func(c *cli.Context) error {
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing cache",
					Action: func(c *cli.Context) error {
						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
