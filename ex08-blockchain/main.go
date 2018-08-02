package main

import (
	"os"
)

func main() {
	if len(os.Args) < 2 {
		println("too few args specified")
		return
	}

	switch os.Args[1] {
	case "add":
		println("adding block")
		Add(os.Args[2:])
	case "mine":
		println("mining new blocks")
	case "list":
		println("listing all blocks")
		List()
	case "drop":
		println("dropping blocks table")
		Drop()
	default:
		println("unknown command")
	}
}
