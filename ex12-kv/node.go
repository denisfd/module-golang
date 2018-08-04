package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Service struct {
	store *Store
}

func (s *Service) HandleConnection(client net.Conn) {
	println("Connected to client")
	req := &Request{}
	req.Get(client)
	fmt.Printf("%+v\n", *req)
}

func (s *Service) Run(port string) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		panic("Cannot start litening on port " + port)
	}

	println("Listening on port", port)

	for {
		conn, _ := ln.Accept()
		go s.HandleConnection(conn)
	}
}

func NewService() *Service {
	s := &Service{}

	s.store = New(false)

	return s
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	node := NewService()

	for {
		fmt.Print("$: ")
		scanner.Scan()
		text := scanner.Text()
		words := strings.Fields(text)
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "exit":
			return

		case "listen":
			node.Listen(words[1:])

		default:
			println("Unknown command:", words[0])
		}
	}
}

func (s *Service) Listen(args []string) {
	/*if len(args) != 1 {
		println("listen: One argument expected")
		return
	}*/
	go s.Run(":8888")
}
