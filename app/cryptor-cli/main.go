package main // import "cpl.li/go/cryptor/app/cryptor-cli"

import (
	"bufio"
	"io"
	"os"

	"github.com/fatih/color"
)

// Version of the Cryptor CLI.
const Version = "v0.0.1"

func main() {
	// set input
	reader := bufio.NewReader(os.Stdin)
	defer reader.Reset(nil)

	// interactive mode
	for {
		// write console symbol
		color.New(color.FgBlue).Fprint(os.Stdout, ">>>")

		// read user input
		text, err := reader.ReadString('\n')
		if err != nil {
			// on EOF, exit
			if err == io.EOF {
				return
			}

			// report error and exit with error
			color.New(color.FgRed).Fprint(os.Stderr, "err: "+err.Error()+"\n")
			os.Exit(1)
		}

		// output
		if err := parseCommand(text); err != nil {
			color.New(color.FgRed).Fprint(os.Stderr, "err: "+err.Error()+"\n")
		}
	}
}
