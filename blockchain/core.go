package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func Mailer(port string, output chan *Message, stop chan byte) {
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		print("> listen: smt went wrong, port is unreachable...\n$: ")
		return
	}

	listenersPool[port] = stop
	print("> listening port ", port, "\n$: ")

	buf := make([]byte, 0, 256)
	tmp := make([]byte, 256)

	for {
		conn, _ := ln.Accept()
		select {
		case <-stop:
			ln.Close()
			println("stopping mailer")
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
		output <- (&message)
	}
}

func Listener(port string) {
	stop := make(chan byte, 2)

	node = NewFollower(port)
	ch := make(chan *Message, 1)
	go Mailer(port, ch, stop)
	for {
		select {
		case <-stop:
			println("Stopping listener")
			return
		default:
			node.Work(ch, stop)
		}
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
		s <- 7
		s <- 0
		s <- 7
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
	if node == nil {
		println("Raft node is not inited, use listen to init")
		return
	}
	m := &Message{
		Mcode:  "SEND",
		Sender: node.Socket(),
		Blocks: []Block{
			Block{
				Data: []string{strings.Join(args[:], " ")},
			},
		},
	}

	node.Send(m)
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

func ip(args []string) string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			saddr := ipv4.String()
			if args == nil {
				return saddr
			}
			fmt.Println("## IPv4:", saddr)
			for port, _ := range listenersPool {
				println(" * ", saddr+":"+port)
			}
			return ""
		}
	}
	return ""
}
