package main

import (
	"cpl.li/go/cryptor/p2p/node"
	"github.com/fatih/color"
)

func init() {
	// initiate other misc
	nodeList = make(map[string]*node.Node)

	// create the commands map
	commands = map[string]command{
		"help": command{
			description: "provide this help message",
			exec:        help,
			helpMessage: color.GreenString("help") + " " +
				color.YellowString("[command]") +
				", provides descriptions and usage instructions of commands",
		},
		"version": command{
			description: "display the aegis and cryptor package versions",
			exec:        version,
		},
		"node": command{
			description: "creation and management of cryptor nodes",
			exec:        commandNode,
			// TODO this mess will be fixed by Issue #29
			helpMessage: color.GreenString("node") + " " +
				color.YellowString("[ new | sel | del | list | start | stop | conn | disc ]\n") +
				"\n" + color.GreenString("node new ") + color.YellowString("[name] [key]") + ", create a new node with the given name and private key" +
				"\n" + color.GreenString("node [sel | select] ") + color.YellowString("[name]") + ", marks the given node as the selected node" +
				"\n" + color.GreenString("node [del | rm ] ") + color.YellowString("[name]") + ", delete the given node, if no name is given, delete selected node" +
				"\n" + color.GreenString("node list ") + ", display a list of all nodes, a '->' symbol marks the selected node" +
				"\n" + color.GreenString("node [start | stop | disc] ") + color.YellowString("[name]") + ", as their name suggests, when no name is given, selected node is used" +
				"\n" + color.GreenString("node conn ") + color.YellowString("[name] [addr]") + ", connects the node to the network, if no addr is given, the current one is used",
		},
		"key": command{
			description: "operations for key derivation and creation",
			exec:        commandKey,
			helpMessage: color.GreenString("key") + " " +
				color.YellowString("[ gen | pass | bip39 ]"),
		},
	}
}

// the command expected execution function
type cmdFunc func(argc int, argv []string) error

type command struct {
	description string
	exec        cmdFunc
	helpMessage string
}

// the lookup map containing all command structs and their command names as key
var commands map[string]command
