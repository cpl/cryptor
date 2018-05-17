package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/thee-engineer/cryptor/config"
	"github.com/thee-engineer/cryptor/net/p2p"
)

func main() {
	config.InitViper()

	logFile, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	var nodes [10]*p2p.Node
	var count = 0

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("err", err)
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		args := strings.Split(line, " ")

		switch args[0] {
		case "new":
			switch args[1] {
			case "node":
				if len(args) != 4 {
					fmt.Println("new node <address> <port>")
					break
				}
				nodes[count] = p2p.NewNode(args[2], args[3], nil)
				fmt.Println("created new node, index", count)
				count++
				break
			}
			break
		case "nodes":
			for _, node := range nodes {
				fmt.Println("node:", node)
			}
		case "connect":
			index, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				break
			}
			nodes[index].Connect()
			break
		case "start":
			index, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				break
			}
			nodes[index].Start()
			break
		case "stop":
			index, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				break
			}
			nodes[index].Stop()
			break
		case "dc":
			index, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				break
			}
			nodes[index].Disconnect()
			break
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("not a command")
		}
	}
}
