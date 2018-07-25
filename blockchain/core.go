package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func Listener(port string) {
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		print("> listen: smt went wrong, port is unreachable...\n$: ")
		return
	}

	stop := make(chan byte, 1)
	listenersPool[port] = stop
	print("> listening port ", port, "\n$: ")

	buf := make([]byte, 0, 256)
	tmp := make([]byte, 256)

	for {
		conn, _ := ln.Accept()
		select {
		case <-stop:
			ln.Close()
			delete(listenersPool, port)
			return
		default:
		}
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
			print("> error while encoding\n$: ")
			return
		}
		fmt.Printf("> %+v\n$: ", message)
	}
}

func listen(args []string) {
	if len(args) != 1 {
		println("## listen: One argument expected(port)")
		return
	}

	reg, _ := regexp.Compile("^\\d+$")

	if reg.MatchString(args[0]) {
		go Listener(args[0])
	} else {
		println("## listen: Wrong Argument", args[0])
	}
}

func stop(args []string) {
	if len(args) != 1 {
		println("## stop: One argument expected")
		return
	}

	s, presents := listenersPool[args[0]]

	if presents {
		s <- 1
		println("## stop: stopping listenning on port", args[0])
		(&Peer{"", "127.0.0.1:" + args[0]}).Send(&Message{}) //it is a bit tricky, but it allows to bypass net.Accept blocking
	} else {
		println("## stop: not listening to port", args[0])
	}
}

func listeners() {
	for k, _ := range listenersPool {
		println(" * port:", k)
	}
}

func send(args []string) {
	if len(args) == 0 {
		println("## send: Nothing to send")
		return
	}
	m := &Message{Mcode: strings.Join(args[:], " ")}

	peerPool.Broadcast(m)
}

func rmpeer(args []string) {
	if len(args) == 0 {
		println("## rmpeer: no args specified, expected template for domain or ip")
		return
	}
	for _, str := range args {
		reg, err := regexp.Compile(str)
		if err != nil {
			continue
		}
		peerPool.RmPeers(reg)
	}
}

func peer(args []string) {
	if len(args) != 2 {
		println("## peer: Two arguments expected(domain ip:port)")
		return
	}

	reg, _ := regexp.Compile("^\\d+.\\d+.\\d+.\\d+:\\d+$")

	if reg.MatchString(args[1]) {
		peerPool.AddPeer(Peer{domain: args[0], addr: args[1]})
	} else {
		println("## peer: wrong address")
	}
}

func ip(args []string) {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			fmt.Println("## IPv4:", ipv4)
			for port, _ := range listenersPool {
				println(" * ", ipv4.String()+":"+port)
			}
			return
		}
	}
}
