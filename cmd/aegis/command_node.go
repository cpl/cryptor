package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/crypt/mwords"
	"cpl.li/go/cryptor/crypt/ppk"
	"cpl.li/go/cryptor/p2p/node"
)

var nodeList map[string]*node.Node
var nodeSelect *node.Node
var nodeSelectName string

func commandNode(argc int, argv []string) error {
	// expect arguments
	if argc == 0 {
		return ErrArgumentCount
	}

	switch argv[0] {
	case "new":
		switch argc {
		case 1:
			keyPrivate, _ = ppk.NewPrivateKey()
			defer crypt.ZeroBytes(keyPrivate[:])
			return commandNodeNew(
				strings.Join(mwords.RandomWords(3), "-"), keyPrivate.ToHex())
		case 2:
			return commandNodeNew(
				strings.Join(mwords.RandomWords(3), "-"), argv[2])
		case 3:
			return commandNodeNew(argv[2], argv[1])
		default:
			return ErrArgumentCount
		}
	case "list", "ls":
		commandNodeList()
	case "start":
		var name string
		if argc != 2 {
			if nodeSelectName == "" {
				return ErrArgumentCount
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
				return ErrArgumentCount
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
	case "disc", "disconnect", "dc":
		var name string
		if argc != 2 {
			if nodeSelectName == "" {
				return ErrArgumentCount
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
				return ErrArgumentCount
			}
			name = nodeSelectName
		} else {
			name = argv[1]
		}
		return commandNodeDel(name)
	case "sel", "select":
		if argc != 2 {
			return ErrArgumentCount
		}
		return commandNodeSelect(argv[1])
	default:
		return errors.New("unexpected argument " + argv[0])
	}

	return nil
}

func commandNodeHelp() {
	fmt.Printf("%s %s\n\n",
		color.GreenString("%s", "node"),
		color.YellowString("%s", "[ new | del | sel | list | start | stop | conn | disc ]"))

	helpPrint("%-13s %-14s %-30s %s\n",
		"node", "new", "[private key] [name]",
		"create a new node with the given name and private hexadecimal key")
	helpPrint("%-13s %-14s %-30s %s\n",
		"node", "del", "[name]",
		"delete a node given its name, (also aliased as: delete, rm, remove)")
	helpPrint("%-13s %-14s %-30s %s\n",
		"node", "sel", "[name]",
		"select a node for quicker or special operations (marked as -> in list), (also aliased as: select)")
	helpPrint("%-13s %-14s %-30s %s\n",
		"node", "list", "",
		"list all created nodes and some information regarding them, (also aliased as: ls)")
	helpPrint("%-13s %-14s %-30s %s\n",
		"node", "start", "[name]",
		"start the node and set state as running")
	helpPrint("%-13s %-14s %-30s %s\n",
		"node", "stop", "[name]",
		"stop the node from running and disconnect if connected")
	helpPrint("%-13s %-14s %-30s %s\n",
		"node", "conn", "[name] [address]",
		"connect the node to the network, listening on the given address (use current node address if none is given)")
	helpPrint("%-13s %-14s %-30s %s\n",
		"node", "disc", "[name]",
		"disconnect the node from the network (also aliased as: disconnect, dc)")
}

func commandNodeList() {
	// nodelist is empty
	if len(nodeList) == 0 {
		color.Blue("no nodes\n")
		return
	}

	// list nodes
	var padding string
	for name, n := range nodeList {
		// check for selected node
		if n == nodeSelect {
			padding = " -> "
		} else {
			padding = "    "
		}

		// format state string with right color
		var stateString string
		nState := n.State()

		switch nState {
		case node.StateStopped:
			stateString = color.RedString(nState.String())
		case node.StateRunning:
			stateString = color.HiGreenString(nState.String())
		case node.StateConnected:
			stateString = color.CyanString(nState.String())
		}

		fmt.Printf("%s%s %s %s %-18s %s %-40s %s %s\n",
			padding,
			color.GreenString("public key"),
			color.YellowString(n.PublicKey().ToHex()),
			color.GreenString("state"),
			stateString,
			color.GreenString("name"),
			color.YellowString(name),
			color.GreenString("addr"),
			color.YellowString(n.Addr()),
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
