package main

import (
	"errors"

	"github.com/fatih/color"

	"cpl.li/go/cryptor/p2p/node"
)

// TODO Finish this

var nodeList []*node.Node

func commandNode(argc int, argv []string) error {
	// expect arguments
	if argc == 0 {
		return errors.New("invalid argument count")
	}

	switch argv[0] {
	case "new":
		if argc != 3 {
			return errors.New("invalid argument count")
		}
		return commandNodeNew(argv[1], argv[2])
	case "list":
		commandNodeList()
	default:
		return errors.New("unexpected argument " + argv[0])
	}

	return nil
}

func commandNodeList() {
	if len(nodeList) == 0 {
		color.Blue("no nodes\n")
	}

	for idx, node := range nodeList {
		color.Green("%d ", idx)
		color.Yellow("%s ", node.PublicKey())
		color.Red("%s\n", node.Addr())
	}
}

func commandNodeNew(name, key string) error {
	return nil
}
