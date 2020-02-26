package scheduler

import (
	"log"
	"math/rand"

	"github.com/JointFaaS/Manager/worker"
)

type scheduleTask struct {
	f func()
}

// Scheduler has serveral roles:
// dispatch tasks to workers
// maintain metrics
// alloc and release workers
type Scheduler struct {
	workers map[string]*worker.Worker

	funcToWorker map[string][]*worker.Worker

	tasks chan *scheduleTask
}

// RegisterWorker 
func (s *Scheduler) RegisterWorker(workerID string, workerAddr string) {
	s.tasks <- &scheduleTask{
		f: func() {
			_, isPresent := s.workers[workerID]
			if isPresent == false {
				newWorker, err := worker.New(workerAddr, workerID)
				log.Printf("New worker: %s %s", workerAddr, workerID)
				if err != nil {
					return
				}
				s.workers[workerID] = newWorker
			}
		},
	}
}

// New returns a scheduler
func New() (*Scheduler, error) {
	s := &Scheduler{
		workers: make(map[string]*worker.Worker),
		funcToWorker: make(map[string][]*worker.Worker),
		tasks: make(chan *scheduleTask, 64),
	}
	return s, nil
}

func (s *Scheduler) Work() {
	go func() {
		for {
			t := <- s.tasks
			t.f()
		}
	}()
}

func (s *Scheduler) GetWorker(funcName string, resCh chan *worker.Worker) {
	s.tasks <- &scheduleTask{
		f: func() {
			funcWorkers, isPresent := s.funcToWorker[funcName]
			if isPresent == false || len(funcWorkers) == 0 {
				resCh <- nil
			} else {
				resCh <- funcWorkers[rand.Intn(len(funcWorkers))]
			}
		},
	}
}