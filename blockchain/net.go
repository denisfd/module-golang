package main

import (
	"encoding/json"
	"net"
)

type Message struct {
	Len    int
	Mcode  string
	Blocks []Block
}

type Peer struct {
	domain string
	addr   string
}

type PeerPool struct {
	stack []Peer
	sp    int
}

func (p *Peer) String() string {
	return p.domain + ": " + p.addr
}

func (p *Peer) Send(m Message) {
	conn, err := net.Dial("tcp", p.addr)
	if err != nil {
		println("Smt bad happened, we got an error while sending to ("+p.String()+"):", err)
		return
	}
	j, _ := json.Marshal(m)
	conn.Write(j)
	conn.Close()
}

func (s *PeerPool) AddPeer(p Peer) {
	if s.sp >= cap(s.stack)-1 {
		s.stack = append(s.stack, p)
	} else {
		s.stack[s.sp] = p
	}
	s.sp += 1
}

func (pool *PeerPool) Print() {
	for _, peer := range pool.stack {
		println(peer.String())
	}
}

func (pool *PeerPool) Broadcast(m Message) {
	for _, peer := range pool.stack {
		peer.Send(m)
	}
}
