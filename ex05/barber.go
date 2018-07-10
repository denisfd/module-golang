package main

// barber works in infinite loop, cntr-C to stop

import (
	"sync"
	"time"
)

type Client struct {
	name   string
	served chan byte
}

type Barber struct {
	name        string
	max_seats   int
	seats_taken int
	clients     chan *Client
	wait_flag   chan byte
	sleep_timer *time.Timer
	mutex       *sync.Mutex
}

func NewClient(name string) *Client {
	c := Client{name: name}
	c.served = make(chan byte)

	return &c
}

func (c *Client) go_to_barber(barber *Barber) {
	println("+", c.name, "attempts to go to", barber.name, "barber")
	barber.mutex.Lock()
	if barber.seats_taken < barber.max_seats {
		stop := barber.sleep_timer.Reset(0)
		if stop {
			println("#", c.name, "woke", barber.name, "up")
		}
		barber.clients <- c
		barber.seats_taken += 1
		println("%", c.name, "takes a seat in a queue")
		barber.mutex.Unlock()
		<-c.served
		println("-", c.name, "is happy and goes home")
	} else {
		barber.mutex.Unlock()
		println("-", c.name, "is upset cause there are no free seats, but he will try to come later")
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		c.go_to_barber(barber)
	}
}

///////
func (b *Barber) start_work() {
	println(b.name, "begins to work, barbershop is now opened")
	println(b.name, "has", b.max_seats, "seat(s) in total")
	b.wait_flag <- 1

	for {
		select {
		case c := <-b.clients:
			println("$", b.name, "serves", c.name)
			time.Sleep(1 * time.Second)
			b.seats_taken -= 1
			c.served <- 1
			println("$", b.name, "served", c.name)
		default:
			println("# No clients,", b.name, "decides to sleep for two seconds")
			b.sleep_timer = time.NewTimer(2 * time.Second)
			<-b.sleep_timer.C
		}
	}

	println(b.name, "finished his work")
	b.wait_flag <- 1
}

func NewBarber(name string, max int) *Barber {
	b := Barber{name: name, max_seats: max, seats_taken: 0}
	b.clients = make(chan *Client, max)
	b.wait_flag = make(chan byte)
	b.mutex = &sync.Mutex{}

	return &b
}

func main() {
	barber := NewBarber("John", 2)
	phil := NewClient("Phil")
	billy := NewClient("Billy")
	steve := NewClient("Steve")
	jack := NewClient("Jack")

	go barber.start_work()
	<-barber.wait_flag

	go billy.go_to_barber(barber)
	go phil.go_to_barber(barber)
	go steve.go_to_barber(barber)
	go jack.go_to_barber(barber)

	<-barber.wait_flag
}
