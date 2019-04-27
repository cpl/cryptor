package main

import "cpl.li/go/cryptor/p2p/node"

func defaultHelp() {

}

func init() {
	// initiate other misc
	nodeList = make(map[string]*node.Node)

	// create the commands map
	commands = map[string]command{
		"help": command{
			description: "provide this help message or usage of other commands",
			exec:        help,
			helpMessage: defaultHelp,
		},
		"version": command{
			description: "display the aegis and cryptor package versions",
			exec:        version,
			helpMessage: defaultHelp,
		},
		"node": command{
			description: "creation and management of cryptor nodes",
			exec:        commandNode,
			helpMessage: commandNodeHelp,
		},
		"key": command{
			description: "operations for key derivation and creation",
			exec:        commandKey,
			helpMessage: commandKeyHelp,
		},
	}
}

// the command expected execution function
type cmdFunc func(argc int, argv []string) error

type command struct {
	description string
	exec        cmdFunc
	helpMessage func()
}

// the lookup map containing all command structs and their command names as key
var commands map[string]command
