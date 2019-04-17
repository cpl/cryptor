package main

func init() {
	// create the commands map
	commands = map[string]command{
		"help": command{
			description: "provide this help message",
			exec:        help,
		},
		"version": command{
			description: "display the cryptor-cli and cryptor pkg versions",
			exec:        version,
		},
	}
}

// the command expected execution function
type cmdFunc func(argc int, argv []string) error

type command struct {
	description string
	exec        cmdFunc
}

// the lookup map containing all command structs and their command names as key
var commands map[string]command
