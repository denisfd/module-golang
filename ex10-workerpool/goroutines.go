package goroutines

//package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	FREE = iota
	WORKING
)

type Job struct {
	time float64
}

type Worker struct {
	id     int
	status int
	job    chan *Job
	stop   chan byte
}

type Scheduler struct {
	size        int
	maxSize     int
	lastID      int
	Jobs        chan *Job
	FreeWorkers chan *Worker // We need channel to pass tests, cause it saves the order,  map does not
	Workers     map[int]*Worker
	wg          sync.WaitGroup
}

/////--------Job
func NewJob(str string) (*Job, error) {
	var f float64
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil, err
	}
	return &Job{time: f}, nil
}

/////--------Worker
func NewWorker(id int) *Worker {
	w := &Worker{}

	w.id = id
	w.status = FREE
	w.job = make(chan *Job, 1)
	w.stop = make(chan byte, 1)

	fmt.Printf("worker:%d spawning\n", w.id)

	return w
}

func (w *Worker) Work(s *Scheduler) {
	for {
		select {
		case <-w.stop:
			fmt.Printf("worker:%d stopping\n", w.id)
			s.size -= 1
			defer s.wg.Done()
			return
		case j := <-w.job:
			fmt.Printf("worker:%d sleep:%.1f\n", w.id, j.time)
			w.status = WORKING
			time.Sleep(time.Millisecond * time.Duration(int(j.time*1000)))
			w.status = FREE
			s.FreeWorkers <- w
		}
	}
}

/////--------Scheduler
func NewScheduler(poolSize int) *Scheduler {
	s := &Scheduler{}

	s.maxSize = poolSize
	s.Jobs = make(chan *Job, (10 + poolSize))
	s.FreeWorkers = make(chan *Worker, poolSize)
	s.Workers = make(map[int]*Worker)

	return s
}

func (s *Scheduler) AddJob(j *Job) {
	s.Jobs <- j
}

func (s *Scheduler) GetFreeWorker() *Worker {
	select {
	case w := <-s.FreeWorkers:
		return w
	default:
		if s.size < s.maxSize {
			s.size += 1
			s.lastID += 1
			w := NewWorker(s.lastID)
			s.Workers[w.id] = w
			s.wg.Add(1)
			go w.Work(s)

			return w
		}
	}
	return nil
}

func (s *Scheduler) Routine() {
	for j := range s.Jobs {
		w := s.GetFreeWorker()
		if w == nil {
			w = <-s.FreeWorkers
		}
		w.job <- j
	}
	for s.size > 0 {
		w := <-s.FreeWorkers
		w.stop <- 1
	}
}

func (s *Scheduler) PrintWorkers() {
	fmt.Printf("  ID -> Status\n")
	for _, worker := range s.Workers {
		switch worker.status {
		case FREE:
			fmt.Printf("*%3d -> FREE\n", worker.id)
		case WORKING:
			fmt.Printf("*%3d -> WORKING\n", worker.id)
		}
	}
}

func (s *Scheduler) KillFree() {
	fmt.Println("killing free workers")
	for len(s.FreeWorkers) > 0 {
		w := <-s.FreeWorkers
		w.stop <- 1
		delete(s.Workers, w.id)
	}
}

/////--------Run
func Read(scheduler *Scheduler) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()
		words := strings.Fields(text)
		if len(words) == 0 {
			continue
		}

		command := words[0]
		switch command {
		case "exit":
			return
		case "ps":
			scheduler.PrintWorkers()
		case "kill":
			scheduler.KillFree()
		default:
			job, err := NewJob(command)
			if err != nil {
				fmt.Println("Invalid input:", err)
				continue
			}
			scheduler.AddJob(job)
		}
	}
	close(scheduler.Jobs)
}

func Run(poolSize int) {
	scheduler := NewScheduler(poolSize)

	go scheduler.Routine()
	Read(scheduler)

	scheduler.wg.Wait()
}

func main() {
	Run(5)
}
