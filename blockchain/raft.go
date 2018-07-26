package main

import (
	"fmt"
	"time"
)

type Node interface {
	Send(*Message)
	Work(chan *Message, chan byte)
	Socket() string
}

type Follower struct {
	port  string
	lport string
}

type Leader struct {
	port string
}

type Candidate struct {
	port string
}

//Follower
func (f *Follower) Send(m *Message) {
	if f.lport == "" {
		println(">Follower: refuse to send, who is a leader?")
	} else {
		println(">Follower: Sending message")
		(&Peer{"", f.lport}).Send(m)
	}
}

func (f *Follower) Work(ch chan *Message, stop chan byte) {
	timer := time.NewTimer(6 * time.Second)
	for {
		select {
		case <-stop:
			println("stopping node")
			return
		case <-timer.C:
			println("FOLLOWER TRANSFORMS TO OMAGELOOL CANDIDATE")
			node = NewCandidate(f.port)
			return
		case message := <-ch:
			fmt.Printf(">Follower: Got Message%+v\n$: ", *message)
			switch message.Mcode {
			case "ELECT":
				(&Peer{"", message.Sender}).Send(&Message{
					Mcode:  "VOTE",
					Sender: f.Socket(),
				})
			case "HB":
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(6 * time.Second)
				println("Ok, follower will wait once more")
				f.lport = message.Sender
			}
		}
	}
}

func NewFollower(port string) *Follower {
	f := new(Follower)

	f.port = port

	return f
}

func (f *Follower) Socket() string {
	return ip(nil) + ":" + f.port
}

//Leader
func (l *Leader) Send(m *Message) {
	(&Peer{"", l.Socket()}).Send(m)
}

func (l *Leader) Work(ch chan *Message, stop chan byte) {
	timer := time.NewTimer(4 * time.Second)
	block := Block{Data: []string{}}
	for {
		select {
		case <-stop:
			println("stopping node")
			return
		case <-timer.C:
			peerPool.Broadcast(&Message{
				Mcode:  "HB",
				Sender: l.Socket(),
				Blocks: []Block{block},
			})
			block.Reset()
			timer.Reset(4 * time.Second)
		case message := <-ch:
			switch message.Mcode {
			case "SEND":
				block.Append(message.Blocks[0].Data[0])
			default:
				fmt.Printf(">Leader: Got Message%+v\n$: ", *message)
			}
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
	println(">Candidate: sending refuse, who is a Leader?")
}

func (c *Candidate) Socket() string {
	return ip(nil) + ":" + c.port
}

func (c *Candidate) Work(ch chan *Message, stop chan byte) {
	timer := time.NewTimer(1 * time.Second)
	peerPool.Broadcast(&Message{
		Mcode:  "ELECT",
		Sender: c.Socket(),
	})
	(&Peer{"", c.Socket()}).Send(&Message{
		Mcode:  "VOTE",
		Sender: c.Socket(),
	})
	voted := 0
	for {
		select {
		case <-stop:
			println("stopping node")
			return
		case <-timer.C:
			if voted > peerPool.sp {
				println("CANDIDATE TRANSFORMS TO OMAGELOOL LEADER")
				node = NewLeader(c.port)
			} else {
				println("CANDIDATE: NOT ENOUGH VOTES")
				node = NewFollower(c.port)
			}
			return
		case message := <-ch:
			fmt.Printf(">Candidate: Got Message%+v\n$: ", *message)
			switch message.Mcode {
			case "VOTE":
				voted += 1
			}
		}
	}
}

func NewCandidate(port string) *Candidate {
	c := new(Candidate)

	c.port = port

	return c
}
