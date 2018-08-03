package main

import (
	"github.com/hashicorp/raft"
	//"github.com/tidwall/raft-fastlog"
)

type Store struct {
	raft *raft.Raft
}

func NewStore() *Store {
	s := &Store{}

	return s
}
