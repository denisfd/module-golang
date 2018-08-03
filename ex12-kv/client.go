package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("$: ")
		scanner.Scan()
		text := scanner.Text()
		words := strings.Fields(text)
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "exit":
			return

		case "gen":
			Gen()

		case "get", "set", "inc", "dec":
			Send(words[:])

		default:
			println("Unknown command:", words[0])
		}
	}
}

func Gen() {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		return
	}

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

	println(hex.EncodeToString(u))
}

func Send(args []string) {
	if len(args) != 3 {
		println("error, 2 args expected -> uuid value")
		return
	}

	f, err := strconv.ParseFloat(args[2], 64)

	if len(args[1]) != 32 {
		println("invalid uuid")
		return
	}

	if err != nil {
		println("invalid value")
		return
	}

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		println("error while dialing:", err)
		return
	}

	req := &Request{}

	req.Op = string(args[0][0])
	req.Value = f
	req.Key = args[1]
	req.Send(conn)
	conn.Close()
}
