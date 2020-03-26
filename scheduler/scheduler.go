package scheduler

import (
	"log"
	"math/rand"
	"time"

	"github.com/JointFaaS/Manager/worker"
)

type scheduleTask struct {
	f func()
}

type PlatformStorageManager interface {
	GetCodeURI(funcName string) (string, error)
	GetImage(funcName string) (string, error)
}
// Scheduler has serveral roles:
// dispatch tasks to workers
// maintain metrics
// alloc and release workers
type Scheduler struct {
	workers map[string]*worker.Worker

	funcToWorker map[string][]*worker.Worker

	tasks chan *scheduleTask

	funcInvokeMetrics map[string]int32

	platformStorageManager PlatformStorageManager
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
func New(p PlatformStorageManager) (*Scheduler, error) {
	s := &Scheduler{
		workers: make(map[string]*worker.Worker),
		funcToWorker: make(map[string][]*worker.Worker),
		tasks: make(chan *scheduleTask, 64),
		funcInvokeMetrics: map[string]int32{},
		platformStorageManager: p,
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
	go func() {
		for {
			s.tasks <- &scheduleTask{
				f: func() {
					for funcName, times := range s.funcInvokeMetrics {
						if times == 0 {
							continue
						}
						s.funcInvokeMetrics[funcName] = 0
						workers, isPresent := s.funcToWorker[funcName]
						if isPresent == false {
							s.funcToWorker[funcName] = make([]*worker.Worker, 0)
						} else if len(workers) > 0 {
							continue
						}
						for _, worker := range s.workers {
							if worker.HasFunction(funcName) {
								continue
							}
							// TODO: exception handle
							image, err := s.platformStorageManager.GetImage(funcName)
							if err != nil {
								return
							}
							codeURI, err := s.platformStorageManager.GetCodeURI(funcName)
							if err != nil {
								return
							}
							err = worker.InitFunction(funcName, image, codeURI)
							if err != nil {
								return
							}
							s.funcToWorker[funcName] = append(s.funcToWorker[funcName], worker)
							break
							// only run once, and golang range is random
						}
					}
				},
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

func (s *Scheduler) GetWorker(funcName string, resCh chan *worker.Worker) {
	s.tasks <- &scheduleTask{
		f: func() {
			times, isPresent := s.funcInvokeMetrics[funcName]
			if isPresent == true {
				s.funcInvokeMetrics[funcName] = times + 1
			}else {
				s.funcInvokeMetrics[funcName] = 1
			}

			funcWorkers, isPresent := s.funcToWorker[funcName]
			if isPresent == false || len(funcWorkers) == 0 {
				resCh <- nil
			} else {
				resCh <- funcWorkers[rand.Intn(len(funcWorkers))]
			}
		},
	}
}

func (s *Scheduler) GetWorkerMust(funcName string, resCh chan *worker.Worker) {
	s.tasks <- &scheduleTask{
		f: func() {
			times, isPresent := s.funcInvokeMetrics[funcName]
			if isPresent == true {
				s.funcInvokeMetrics[funcName] = times + 1
			}else {
				s.funcInvokeMetrics[funcName] = 1
			}

			funcWorkers, isPresent := s.funcToWorker[funcName]
			if isPresent == false || len(funcWorkers) == 0 {
				for _, worker := range s.workers {
					// TODO: exception handle
					image, err := s.platformStorageManager.GetImage(funcName)
					if err != nil {
						return
					}
					codeURI, err := s.platformStorageManager.GetCodeURI(funcName)
					if err != nil {
						return
					}
					err = worker.InitFunction(funcName, image, codeURI)
					if err != nil {
						return
					}
					resCh <- worker
					s.funcToWorker[funcName] = append(s.funcToWorker[funcName], worker)
					break
					// only run once, and golang range is random
				}
				resCh <- nil
			} else {
				resCh <- funcWorkers[rand.Intn(len(funcWorkers))]
			}
		},
	}
}