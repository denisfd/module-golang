package main

import (
	"fmt"
	//"time"
)

type Block struct {
	PrevHash string
	CurHash  string
	Data     []string
}

type Node interface {
	Send(*Message)
	Work(chan *Message, chan byte)
	Socket() string
}

type Follower struct {
	port string
}

type Leader struct {
	port string
	buf  chan string
}

type Candidate struct {
	port string
}

//Follower
func (f *Follower) Send(m *Message) {

}

func (f *Follower) Work(m *Message) {
	fmt.Printf(">Follower: %+v\n$: ", *m)
}

func NewFollower(port string) *Follower {
	return &Follower{port: port}
}

func (f *Follower) Socket() string {
	return ip(nil) + ":" + f.port
}

//Leader
func (l *Leader) Send(m *Message) {
	peerPool.Broadcast(m)
}

func (l *Leader) Work(ch chan *Message, stop chan byte) {
	for {
		select {
		case <-stop:
			println("stopping node")
			return
		case message := <-ch:
			fmt.Printf(">Leader: Got Message%+v\n$: ", *message)
		}
	}
}

func (l *Leader) Socket() string {
	return ip(nil) + ":" + l.port
}

func NewLeader(port string) *Leader {
	l := new(Leader)

	l.port = port

	return l
}

//Candidate
func (c *Candidate) Send(m *Message) {

}

func (c *Candidate) Work() {

}

func NewCandidate() *Candidate {
	return &Candidate{}
}
