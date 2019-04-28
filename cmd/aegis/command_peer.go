package main

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
)

func commandPeer(argc int, argv []string) error {
	// check for selected node
	if nodeSelect == nil {
		return errors.New("peer command needs a selected node, check `help node`")
	}

	// expect arguments
	if argc == 0 {
		return ErrArgumentCount
	}

	switch argv[0] {
	case "add":
		if argc == 3 {
			return ErrArgumentCount
		}
		return commandPeerAdd()
	case "del", "delete", "rm", "remove":
		if argc == 2 {
			return ErrArgumentCount
		}
		return commandPeerDel()
	case "hs":
		if argc == 2 {
			return ErrArgumentCount
		}
		return commandPeerHandshake()
	default:
		return errors.New("unexpected argument " + argv[0])
	}
}

func commandPeerHelp() {
	fmt.Printf("%s can only be used when a node is selected, check %s\n\n",
		color.GreenString("peer"),
		color.YellowString("help node"))

	fmt.Printf("%s %s\n\n",
		color.GreenString("peer"),
		color.YellowString("[ add | del | hs ]"))

	helpPrint("%-13s %-14s %-33s %s\n",
		"peer", "add", "[public key] [host:port]",
		"add a new peer to the selected node, address is optional")
	helpPrint("%-13s %-14s %-33s %s\n",
		"peer", "del", "[public key | id]",
		"delete a peer given its public key or ID (also aliased as: rm, remove, delete)")
	helpPrint("%-13s %-14s %-33s %s\n",
		"peer", "hs", "[public key | id]",
		"manually initialize the handshake protocol with the given peer")
}

func commandPeerAdd() error {
	return nil
}

func commandPeerDel() error {
	return nil
}

func commandPeerHandshake() error {
	return nil
}
