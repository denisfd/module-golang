package main

import (
	"encoding/json"
	"net"
	"regexp"
)

type Message struct {
	Len    int
	Sender string
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
	return p.domain + " -> " + p.addr
}

func (p *Peer) Send(m *Message) {
	conn, err := net.Dial("tcp", p.addr)
	if err != nil {
		println("Smt bad happened, we got an error while sending to ("+p.String()+"):", err)
		return
	}
	j, _ := json.Marshal(*m)
	conn.Write(j)
	conn.Close()
}

func (p *Peer) Matches(reg *regexp.Regexp) bool {
	if reg.MatchString(p.domain) || reg.MatchString(p.addr) {
		return true
	}
	return false
}

func (s *PeerPool) AddPeer(p Peer) {
	if s.sp >= cap(s.stack)-1 {
		s.stack = append(s.stack, p)
	} else {
		s.stack[s.sp] = p
	}
	s.sp += 1
}

func (pool *PeerPool) RmPeer(index int) {
	if index >= pool.sp { //pool.sp -> len of peerPool
		return
	}
	pool.stack[index] = pool.stack[pool.sp-1] //we do not care about order, so we can delete
	//elems replacing them with last elements and decrementing stack pointer sp,
	//deleting with append(a[:i], a[i+1:]...) is not very nice
	pool.sp -= 1
}

func (pool *PeerPool) RmPeers(reg *regexp.Regexp) {
	for i := 0; i < pool.sp; i++ { //we cannot use range here... :'(((
		if pool.stack[i].Matches(reg) {
			pool.RmPeer(i)
			i -= 1
		}
	}
}

func (pool *PeerPool) Print() {
	for i := 0; i < pool.sp; i++ {
		println(" * ", pool.stack[i].String())
	}
}

func (pool *PeerPool) Broadcast(m *Message) {
	for i := 0; i < pool.sp; i++ {
		pool.stack[i].Send(m)
	}
}
