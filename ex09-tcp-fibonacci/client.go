package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Message struct {
	Msg string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		return
	}
	defer conn.Close()

	for scanner.Scan() {
		text := scanner.Text()
		if text == "q" {
			return
		}

		m := Message{text}
		e := json.NewEncoder(conn)
		e.Encode(m)

		d := json.NewDecoder(conn)
		err = d.Decode(&m)

		if err != nil {
			continue
		}

		fmt.Fprintf(os.Stdout, "%s\n", m.Msg)
	}
}
