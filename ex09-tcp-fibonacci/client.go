package main

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
)

type Message struct {
	Msg string
}

func main() {
	b := bufio.NewReader(os.Stdin)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		println("Smt bad happened, we got an error while sending to:", err)
		return
	}
	defer conn.Close()

	for {
		line, _ := b.ReadBytes('\n')
		if line[0] == 'q' {
			return
		}

		m := Message{string(line[:len(line)-1])}
		e := json.NewEncoder(conn)
		e.Encode(m)

		d := json.NewDecoder(conn)
		err := d.Decode(&m)

		if err != nil {
			continue
		}

		println(m.Msg)
	}
}
