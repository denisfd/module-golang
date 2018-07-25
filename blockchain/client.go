package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var peerPool PeerPool
var listenersPool = make(map[string]chan byte)
var node *Node

func Config() {
	listen([]string{"7755"})
	peer([]string{"me", "127.0.0.1:7755"})
}

func help() {
	println("$$ USAGE:")
	println(" * send message -> broadcasts message to all your peers")
	println(" * peer name ip:port -> adds new peer for broadcasting, use any name you like")
	println(" * peers -> display all peers")
	println(" * rmpeer template1 template2 ... -> deletes all peers which satisfy given templates (name OR ip:port)")
	println(" * listen port -> starts listening to given port")
	println(" * listeners -> prints all porte you are listening to")
	println(" * stop port -> stop listening given port")
	println(" * ip -> displays your ip and all ip:port pairs you are listening to")
	println(" * cfg -> adds listener and peer, so you can send messages to yourself")
}

func processInput() {
	defer wg.Done()
	scanner := bufio.NewScanner(os.Stdin)

	println("Type help to get help...")

	for {
		fmt.Print("$: ")
		scanner.Scan()
		text := scanner.Text()
		words := strings.Fields(text)
		if len(words) == 0 {
			continue
		}
		command := words[0]

		switch command {
		case "exit":
			return
		case "cfg":
			Config()
		case "listen":
			listen(words[1:])
		case "listeners":
			listeners()
		case "stop":
			stop(words[1:])
		case "peer":
			peer(words[1:])
		case "peers":
			peerPool.Print()
		case "rmpeer":
			rmpeer(words[1:])
		case "send":
			send(words[1:])
		case "ip":
			ip(words[1:])
		case "help":
			help()
		default:
			println("Unknown command:", command)
		}
	}
}

func main() {
	wg.Add(1)
	go processInput()
	wg.Wait()
}
