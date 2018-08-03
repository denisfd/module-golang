package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
		panic("Cannot start litening on porti " + port)
	}

	for {
		conn, _ := ln.Accept()
		go s.HandleConnection(conn)
	}
}

func NewService() *Service {
	s := &Service{}

	s.store = NewStore()

	return s
}

func main() {
	node := NewService()

	go node.Run(":8888")

	println("Press ENTER to stop")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
