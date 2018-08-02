package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"sync"
	"time"
)

type Message struct {
	Msg string
}

type Server struct {
	Precounted map[int]*Fibonacci
	*sync.Mutex
}

type Fibonacci struct {
	value *big.Int
	prev  *big.Int
}

func Init() *Server {
	s := &Server{}

	s.Precounted = make(map[int]*Fibonacci)
	s.Precounted[1] = &Fibonacci{value: big.NewInt(1), prev: big.NewInt(0)}
	s.Precounted[2] = &Fibonacci{value: big.NewInt(1), prev: big.NewInt(1)}
	s.Mutex = &sync.Mutex{}

	return s
}

func (s *Server) Run(port string) {
	ln, _ := net.Listen("tcp", port)
	if ln == nil {
		panic("Cannot start listening " + port)
	}
	println("listening on port", port)

	for {
		conn, _ := ln.Accept()
		go s.handleConn(conn)
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (s *Server) countFor(num int) {
	_, ok := s.Precounted[num]
	if ok {
		return
	}
	closest := 1
	for i, _ := range s.Precounted {
		if abs(num-i) < abs(num-closest) {
			closest = i
		}
	}

	cur := new(big.Int).Set(s.Precounted[closest].value)
	prev := new(big.Int).Set(s.Precounted[closest].prev)
	temp := new(big.Int)

	if closest < num {
		for closest < num {
			temp.Set(cur)
			cur.Add(cur, prev)
			prev.Set(temp)
			closest += 1
		}
	} else {
		for closest > num {
			temp.Set(prev)
			prev.Sub(cur, prev)
			cur.Set(temp)
			closest -= 1
		}
	}

	s.Lock()
	s.Precounted[num] = &Fibonacci{value: cur, prev: prev}
	s.Unlock()
}

func (s *Server) handleConn(client net.Conn) {
	for {
		m := Message{}

		d := json.NewDecoder(client)
		e := json.NewEncoder(client)
		err := d.Decode(&m)

		if err != nil {
			continue
		}

		var i int
		_, err = fmt.Sscan(m.Msg, &i)
		if err != nil {
			m.Msg = "Wrong Argument"
			e.Encode(m)
			continue
		}

		start := time.Now()
		s.countFor(i)
		elapsed := time.Since(start)

		//m.Msg = fmt.Sprintf("%s ", elapsed)

		m.Msg = fmt.Sprintf("%s", s.Precounted[i].value)

		e.Encode(m)
	}
}

func main() {
	server := Init()
	go server.Run(":7777")
	println("Press ENTER to stop")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
