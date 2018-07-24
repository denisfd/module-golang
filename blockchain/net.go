package main

import (
	"net"
)

type Message struct {
	id  int
	Str string
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

func (p *Peer) Send(m []byte) {
	conn, err := net.Dial("tcp", p.addr)
	if err != nil {
		println("Smt bad happened, we got an error while sending to ("+p.String()+"):", err)
		return
	}
	conn.Write(m)
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

func (pool *PeerPool) Broadcast(m []byte) {
	for _, peer := range pool.stack {
		peer.Send(m)
	}
}
