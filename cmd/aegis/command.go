package main

import "cpl.li/go/cryptor/p2p/node"

func defaultHelp() {

}

func init() {
	// initiate other misc
	nodeList = make(map[string]*node.Node)

	// create the commands map
	commands = map[string]command{
		"help": {
			description: "provide this help message or usage of other commands",
			exec:        help,
			helpMessage: defaultHelp,
		},
		"version": {
			description: "display the aegis and cryptor package versions",
			exec:        version,
			helpMessage: defaultHelp,
		},
		"node": {
			description: "creation and management of cryptor nodes",
			exec:        commandNode,
			helpMessage: commandNodeHelp,
		},
		"key": {
			description: "operations for key derivation and creation",
			exec:        commandKey,
			helpMessage: commandKeyHelp,
		},
		"peer": {
			description: "command for managing node peers and their connections",
			exec:        commandPeer,
			helpMessage: commandPeerHelp,
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
