// cryptor
package main

import (
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

	app.EnableBashCompletion = true
	app.Version = cryptorVersion
	app.Name = cryptorName
	app.Compiled = time.Now()
	app.Usage = cryptorUsageDescription
	app.UsageText = cryptorName + cryptorUsageString
	app.Copyright = cryptorCopyright

	app.Commands = []cli.Command{
		{
			Name:  "chunk",
			Usage: "chunk a file given a input package and output cache",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "assemble",
			Usage: "assemble a package given the key/pass and input cache",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
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
