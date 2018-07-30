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
	FreeWorkers chan *Worker
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
		j := <-w.job
		fmt.Printf("worker:%d sleep:%.1f\n", w.id, j.time)
		time.Sleep(time.Millisecond * time.Duration(int(j.time*1000)))
		s.FreeWorkers <- w
	}
}

/////--------Scheduler
func NewScheduler(poolSize int) *Scheduler {
	s := &Scheduler{}

	s.maxSize = poolSize
	s.Jobs = make(chan *Job, (10 + poolSize))
	s.FreeWorkers = make(chan *Worker, poolSize)

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
		for ; w == nil; w = s.GetFreeWorker() {
		}
		w.job <- j
	}
	for i := 0; i < s.size; i++ {
		w := <-s.FreeWorkers
		fmt.Printf("worker:%d stopping\n", w.id)
		s.wg.Done()
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
			close(scheduler.Jobs)
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

	scheduler.wg.Wait()
}

func main() {
	Run(2)
}
