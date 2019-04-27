package main

import (
	"fmt"
	"os"

	"cpl.li/go/cryptor"
	"github.com/fatih/color"
)

const helpMsg = `
aegis is the main Cryptor command line interface that allows for the management
of nodes, peers keys and all other aspects of the Cryptor package.
`

// utility help function
func help(argc int, argv []string) error {
	// check for usage instructions
	if argc == 1 {
		// check if usage exists
		if msg, ok := commands[argv[0]]; ok {
			msg.helpMessage()
		}

		return nil
	}

	// display help message
	fmt.Println(helpMsg)

	// list commands
	fmt.Printf("The commands are:\n\n")
	for name, cmd := range commands {
		color.New(color.FgGreen).Fprintf(os.Stdout, "    %-10s", name)
		color.New(color.FgYellow).Fprintf(os.Stdout, "%s\n", cmd.description)
	}

	// extended help instructions
	fmt.Printf("\nUse `%s` for more information about a command.\n",
		color.YellowString("help <command>"))

	// documentation website
	fmt.Printf("\nFor the full Cryptor documentation visit %s.\n",
		color.BlueString("https://cpl.li/cryptor"))

	return nil
}

func helpPrint(fmtString, cmd, subcmd, args, msg string) {
	fmt.Printf(fmtString, color.GreenString(cmd), color.YellowString(subcmd),
		color.HiYellowString(args), msg)
}

// utility version function
func version(argc int, argv []string) error {

	fmt.Printf("    %-16s %s\n",
		color.GreenString("cryptor"),
		color.YellowString(cryptor.Version))
	fmt.Printf("    %-16s %s\n",
		color.GreenString("aegis"),
		color.YellowString(Version))

	return nil
}
