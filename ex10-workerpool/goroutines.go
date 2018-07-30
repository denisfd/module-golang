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
	WORKING = iota
	FREE
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
	Workers     map[int]*Worker //id -> Worker
	FreeWorkers chan *Worker    //it helps to pass testcases, cause it saves the order
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

func (w *Worker) Work(j *Job) {
	w.status = WORKING
	w.job <- j
}

func (w *Worker) Routine(scheduler *Scheduler) {
	for {
		select {
		case <-w.stop:
			fmt.Printf("worker:%d stopping\n", w.id)
			scheduler.size -= 1
			scheduler.wg.Done()
			return
		case j := <-w.job:
			fmt.Printf("worker:%d sleep:%v\n", w.id, j.time)
			time.Sleep(time.Duration(j.time) * time.Millisecond)
			scheduler.wg.Done()
			w.status = FREE
			scheduler.FreeWorkers <- w
		}
	}
}

/////--------Scheduler
func NewScheduler(poolSize int) *Scheduler {
	s := &Scheduler{}

	s.maxSize = poolSize
	s.Jobs = make(chan *Job, (10 + poolSize))
	s.Workers = make(map[int]*Worker)
	s.FreeWorkers = make(chan *Worker, poolSize)

	return s
}

func (s *Scheduler) AddJob(j *Job) {
	s.wg.Add(1)
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
			go w.Routine(s)
			return w
		}
	}
	return nil
}

func (s *Scheduler) Routine() {
	for {
		if len(s.Jobs) > 0 {
			if w := s.GetFreeWorker(); w != nil {
				w.Work(<-s.Jobs)
			}
		} else {
			for len(s.FreeWorkers) > 0 {
				w := <-s.FreeWorkers
				w.stop <- 1
				delete(s.Workers, w.id)
			}
		}
	}
}

/////--------Run
func Read(scheduler *Scheduler) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		text := scanner.Text()
		words := strings.Fields(text)

		if len(words) == 0 {
			return
		}

		command := words[0]

		job, err := NewJob(command)
		if err != nil {
			fmt.Println("Invalid input:", err)
			continue
		}

		scheduler.AddJob(job)
	}
}

func Run(poolSize int) {
	scheduler := NewScheduler(poolSize)

	go scheduler.Routine()
	Read(scheduler)

	fmt.Printf("AWAITING\n")
	scheduler.wg.Wait()
	fmt.Printf("AWAITED\n")
}

func main() {
	Run(1)
}
