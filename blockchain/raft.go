package main

type Block struct {
	PrevHash string
	CurHash  string
	data     []string
}

type Node interface {
	Send(*Message)
	Work()
	Init()
}

type Follower struct {
}

type Leader struct {
}

type Candidate struct {
}

//Follower
func (f *Follower) Send(m *Message) {

}

func (f *Follower) Work() {

}

func (f *Follower) Init() {
	f = &Follower{}
}

//Leader
func (l *Leader) Send(m *Message) {

}

func (l *Leader) Work() {

}

func (l *Leader) Init() {
	l = &Leader{}
}

//Candidate
func (c *Candidate) Send(m *Message) {

}

func (c *Candidate) Work() {

}

func (c *Candidate) Init() {
	c = &Candidate{}
}
