package main

import (
	"errors"
	"os"
	"strings"
)

const (
	cmdMinLen = 3
	cmdMaxLen = 100 // ! DEBUG

	cmdMinArgc = 1
	cmdMaxArgc = 10 // ! DEBUG
)

func parseCommand(cmd string) error {
	// trim any leading and trailing whitespaces
	cmd = strings.TrimSpace(cmd)

	// validate len
	if len(cmd) < cmdMinLen || len(cmd) > cmdMaxLen {
		return errors.New("invalid command length")
	}

	// split command string into individual arguments
	argv := strings.Fields(cmd)
	argc := len(argv)
	if argc < cmdMinArgc || argc > cmdMaxArgc {
		return errors.New("invalid command argument count")
	}

	// builtin simple commands
	switch argv[0] {
	case "exit", "quit":
		os.Exit(0)
	case "clear", "cls":
		if _, err := os.Stdout.WriteString("\033[H\033[2J"); err != nil {
			return err
		}
		return nil
	}

	// lookup command
	c, ok := commands[argv[0]]
	if !ok {
		return errors.New("command not found")
	}

	// execute command
	if argc == 1 {
		return c.exec(0, nil)
	}
	return c.exec(argc-1, argv[1:])
}
