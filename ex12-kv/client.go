package main

import (
	"bufio"
	"fmt"
	"os"
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

		case "get":
			Get(words[1:])
		case "set":
			Set(words[1:])
		case "inc":
			Inc(words[1:])
		case "dec":
			Dec(words[1:])

		default:
			println("Unknown command:", words[0])
		}
	}
}

func Get(args []string) {

}

func Set(args []string) {

}
func Inc(args []string) {

}

func Dec(args []string) {

}
