package main

import (
	"errors"
	"fmt"

	"github.com/fatih/color"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/p2p/node"
)

var nodeList map[string]*node.Node
var nodeSelect *node.Node
var nodeSelectName string

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
	case "start":
		var name string
		if argc != 2 {
			if nodeSelectName == "" {
				return errors.New("invalid argument count")
			}
			name = nodeSelectName
		} else {
			name = argv[1]
		}
		return commandNodeStart(name)
	case "stop":
		var name string
		if argc != 2 {
			if nodeSelectName == "" {
				return errors.New("invalid argument count")
			}
			name = nodeSelectName
		} else {
			name = argv[1]
		}
		return commandNodeStop(name)
	case "conn", "connect":
		var addr, name string
		if argc == 3 {
			addr = argv[2]
			name = argv[1]
		} else {
			addr = ""
			if argc == 2 {
				name = argv[1]
			} else {
				name = nodeSelectName
			}
		}
		return commandNodeConn(name, addr)
	case "disc", "disconnect":
		var name string
		if argc != 2 {
			if nodeSelectName == "" {
				return errors.New("invalid argument count")
			}
			name = nodeSelectName
		} else {
			name = argv[1]
		}
		return commandNodeDisc(name)
	case "del", "delete", "rm", "remove":
		var name string
		if argc != 2 {
			if nodeSelectName == "" {
				return errors.New("invalid argument count")
			}
			name = nodeSelectName
		} else {
			name = argv[1]
		}
		return commandNodeDel(name)
	case "sel", "select":
		if argc != 2 {
			return errors.New("invalid argument count")
		}
		return commandNodeSelect(argv[1])
	default:
		return errors.New("unexpected argument " + argv[0])
	}

	return nil
}

func commandNodeList() {
	// nodelist is empty
	if len(nodeList) == 0 {
		color.Blue("no nodes\n")
		return
	}

	// list nodes
	var padding string
	for name, node := range nodeList {
		// check for selected node
		if node == nodeSelect {
			padding = " -> "
		} else {
			padding = "    "
		}

		fmt.Printf("%s%s %s %s %-16s %s %s\n",
			padding,
			color.GreenString("public key"),
			color.YellowString(node.PublicKey().ToHex()),
			color.GreenString("name"),
			color.YellowString(name),
			color.GreenString("addr"),
			color.YellowString(node.Addr()),
		)
	}
}

func commandNodeStart(name string) error {
	// get node
	n, ok := nodeList[name]
	if !ok {
		return errors.New("no node found with name " + name)
	}

	// start
	return n.Start()
}

func commandNodeStop(name string) error {
	// get node
	n, ok := nodeList[name]
	if !ok {
		return errors.New("no node found with name " + name)
	}

	// stop
	return n.Stop()
}

func commandNodeConn(name, addr string) error {
	// get node
	n, ok := nodeList[name]
	if !ok {
		return errors.New("no node found with name " + name)
	}

	// set address
	if addr != "" {
		if err := n.SetAddr(addr); err != nil {
			return err
		}
	}

	// connect
	return n.Connect()
}

func commandNodeDisc(name string) error {
	// get node
	n, ok := nodeList[name]
	if !ok {
		return errors.New("no node found with name " + name)
	}

	// disconnect
	return n.Disconnect()
}

func commandNodeDel(name string) error {
	// get node
	n, ok := nodeList[name]
	if !ok {
		return errors.New("no node found with name " + name)
	}

	// stop
	if err := n.Stop(); err != nil && err.Error() != "can't stop, node not running" {
		return err
	}

	// remove select node ref
	if n == nodeSelect {
		nodeSelect = nil
		nodeSelectName = ""
	}

	// delete node
	n = nil
	delete(nodeList, name)
	return nil
}

func commandNodeNew(name, key string) error {
	// check for node with the same name
	if _, ok := nodeList[name]; ok {
		return errors.New("node already exists with name " + name)
	}

	// decode key
	if err := keyPrivate.FromHex(key); err != nil {
		return err
	}
	defer crypt.ZeroBytes(keyPrivate[:])

	// create and assign node
	nodeList[name] = node.NewNode(name, keyPrivate)

	return nil
}

func commandNodeSelect(name string) error {
	// get node
	n, ok := nodeList[name]
	if !ok {
		return errors.New("no node found with name " + name)
	}

	nodeSelect = n
	nodeSelectName = name
	return nil
}
