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

func Config() {
	listen([]string{"7755"})
	peer([]string{"me", "127.0.0.1:7755"})
}

func processInput() {
	defer wg.Done()
	scanner := bufio.NewScanner(os.Stdin)

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
		case "listen":
			listen(words[1:])
		case "peer":
			peer(words[1:])
		case "peers":
			peerPool.Print()
		case "send":
			send(words[1:])
		default:
			println("Unknown command:", command)
		}
	}
}

func main() {
	wg.Add(1)
	Config()
	go processInput()
	wg.Wait()
}
