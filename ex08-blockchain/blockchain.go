package main

import (
	"bufio"
	"crypto/sha256"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

type Block struct {
	Data     string
	PrevHash string
	CurHash  string
}

func NewBlock(prevBlock *Block, data string) *Block {
	block := Block{}

	block.PrevHash = prevBlock.CurHash
	block.Data = data
	block.CurHash = fmt.Sprintf("%x", sha256.Sum256([]byte(block.Data+block.PrevHash)))

	return &block
}

func Add(args []string) {
	if len(args) == 0 {
		println("no args specified")
		return
	}

	data := strings.Join(args, " ")

	database, _ := sql.Open("sqlite3", "./blocks.db")
	Init(database)
	b := GetLastBlock(database)

	fmt.Printf("HEAD: %+v\n", *b)

	sqlAdd(NewBlock(b, data).MineHash(""), database)

}

func GetLastBlock(database *sql.DB) *Block {
	last, _ := database.Query("SELECT data, curhash, prevhash FROM blocks ORDER BY id DESC LIMIT 1")

	block := Block{}

	for last.Next() {
		last.Scan(&block.Data, &block.CurHash, &block.PrevHash)
		for last.Next() {
		}
		return &block
	}

	return &block
}

func Init(database *sql.DB) {
	ch, _ := database.Query("SELECT count(*) FROM sqlite_master WHERE type='table' and name='blocks'")
	var i int
	for ch.Next() {
		ch.Scan(&i)
		if i == 0 {
			for ch.Next() {
			}
			println("Empty table")
			statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS blocks (id INTEGER PRIMARY KEY, data TEXT, curhash TEXT, prevhash TEXT)")
			statement.Exec()
			sqlAdd(((&Block{
				Data:     "genesis",
				PrevHash: "",
			}).MineHash("")), database)
			break
		}
	}
}

func sqlAdd(block *Block, database *sql.DB) {
	statement, _ := database.Prepare("INSERT INTO blocks (data, curhash, prevhash) VALUES (?, ?, ?)")
	statement.Exec(block.Data, block.CurHash, block.PrevHash)
}

func List() {
	database, _ := sql.Open("sqlite3", "./blocks.db")

	Init(database)

	rows, _ := database.Query("SELECT id, data, curhash, prevhash FROM blocks")
	var id int
	var data string
	var hash string
	var prev string

	for rows.Next() {
		rows.Scan(&id, &data, &hash, &prev)
		fmt.Printf("ID: #%d, DATA -> %s\n", id, data)
		fmt.Printf("CURHASH:  %s\n", hash)
		fmt.Printf("PREVHASH: %s\n\n", prev)
	}
}

func Drop() {
	database, _ := sql.Open("sqlite3", "./blocks.db")
	statement, _ := database.Prepare("DROP TABLE IF EXISTS blocks")
	statement.Exec()
}

func (block *Block) MineHash(diff string) *Block {
	var counter int = 0
	var eq bool
	for {
		eq = true
		block.CurHash = fmt.Sprintf("%x", sha256.Sum256([]byte(block.PrevHash+block.Data+fmt.Sprintf("%x", counter))))
		for i := 0; i < len(diff); i++ {
			if diff[i] != block.CurHash[i] {
				eq = false
				break
			}
		}
		if eq {
			return block
		}
		counter += 1
	}
}

func MineBlocks(diff string) {
	database, _ := sql.Open("sqlite3", "./blocks.db")
	Init(database)
	for {
		b := GetLastBlock(database)

		new := NewBlock(b, "mined").MineHash(diff)
		sqlAdd(new, database)

		fmt.Printf("Block mined, data -> %s\n", new.Data)
		fmt.Printf("Curhash:  %s\n", new.CurHash)
		fmt.Printf("Prevhash: %s\n\n", new.PrevHash)
	}
}

func Mine(args []string) {
	if len(args) == 0 {
		println("No difficulty specified. Specify a string to start mining (like aad0, 000000, 123412, 11, etc)")
		return
	}
	diff := args[0]
	if len(diff) > 10 {
		println("This difficulty is toooooooo difficult......")
		return
	}

	for i := 0; i < len(diff); i++ {
		if (diff[i] >= 'a' && diff[i] <= 'f') || (diff[i] >= '0' && diff[i] <= '9') {
			continue
		}
		println("Invalid difficulty, expected hexidecimal value")
		return
	}

	println("press Enter to stop mining")
	go MineBlocks(diff)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
