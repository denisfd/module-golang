package main

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strings"
)

func Listener(port string) {
	ln, _ := net.Listen("tcp", ":"+port)

	buf := make([]byte, 0, 256)
	tmp := make([]byte, 256)

	for {
		conn, _ := ln.Accept()
		for {
			n, err := conn.Read(tmp)
			if err != nil {
				break
			}
			buf = append(buf, tmp[:n]...)
		}
		var message Message
		err := json.Unmarshal(buf, &message)
		buf = buf[:0]
		if err != nil {
			println("Listener: error while encoding")
			return
		}
		fmt.Printf("> %+v\n", message)
	}
}

func listen(args []string) {
	if len(args) != 1 {
		println("listen: One argument expected(port)")
		return
	}

	reg, _ := regexp.Compile("^\\d+$")

	if reg.MatchString(args[0]) {
		println("listening port", args[0])
		go Listener(args[0])
	} else {
		println("listen: Wrong Argument", args[0])
	}
}

func send(args []string) {
	if len(args) == 0 {
		println("send: Nothing to send")
		return
	}
	m := Message{Mcode: strings.Join(args[:], " ")}

	peerPool.Broadcast(m)
}

func peer(args []string) {
	if len(args) != 2 {
		println("peer: Two arguments expected(domain ip:port)")
		return
	}

	reg, _ := regexp.Compile("^\\d+.\\d+.\\d+.\\d+:\\d+$")

	if reg.MatchString(args[1]) {
		peerPool.AddPeer(Peer{domain: args[0], addr: args[1]})
	} else {
		println("peer: wrong address")
	}
}
