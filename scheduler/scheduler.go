package scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	/* key/value : funcName/workers[], means that which workers have the func's metadata (images and so on) */
	funcToWorker map[string][]*worker.Worker

	tasks chan *scheduleTask

	/* key/value : funcName/invoke times, to record how many times that a function is invoked */
	funcInvokeMetrics map[string]int32

	platformStorageManager PlatformStorageManager

	/* key/value : workerId/bias, to record the percentile of the possibility that a worker will be scheduled */
	//TODO: bias is worker granular but not worker/function granular
	bias map[string]int
}

//The total bias value of all workers
const totalBudget = 100

//PromoteBias promote the bias of a designate worker. It also reduce the other worker's bias in proportion
func (s *Scheduler) PromoteBias(workerID string, promoteValue int) ([]byte, error) {
	_, isPresent := s.bias[workerID]
	if isPresent == false {
		fmt.Println("error: worker not found")
		return nil, errors.New("worker not found")
	}

	if promoteValue < 0 {
		fmt.Println("error: promote bias should be positive")
		return nil, errors.New("negative promote bias")
	} else if promoteValue == 0 {
		return json.Marshal(s.bias)
	}

	var reduceProportion float32 = float32(totalBudget) / float32(totalBudget+promoteValue)
	for wID, workerBias := range s.bias {
		if wID == workerID {
			workerBias += promoteValue
		}
		s.bias[wID] = int(float32(workerBias) * reduceProportion)
		fmt.Printf("[liu] worker %s bias is %d, promoteValue is %d, workerID is %s\n", wID, int(float32(workerBias)*reduceProportion), promoteValue, workerID)
	}
	b, _ := json.Marshal(s.bias)
	return b, nil
}

// RegisterWorker alloc a new worker and add it to s.worker
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

				//initialize the bias vector
				var newBias int = totalBudget / int(len(s.workers))
				for wID := range s.bias {
					s.bias[wID] = newBias
				}
				s.bias[workerID] = newBias
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

		bias: map[string]int{},

		platformStorageManager: p,
	}
	return s, nil
}

func (s *Scheduler) Work() {
	// process task
	go func() {
		for {
			t := <-s.tasks
			t.f()
		}
	}()
	// produce schedule task regularly. The task register the metadata of every function to each worker.
	go func() {
		const totalBudget = 100

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
						fmt.Printf("[liu] workers: %v\n", s.workers)
						for _, worker := range s.workers {
							fmt.Printf("[liu] assigned worker: %v\n", worker)
							if worker.HasFunction(funcName) {
								continue
							}
							//Init the function if the function has not registered to the worker.
							// TODO: exception handle
							image, err := s.platformStorageManager.GetImage(funcName)
							if err != nil {
								return
							}
							codeURI, err := s.platformStorageManager.GetCodeURI(funcName)
							if err != nil {
								return
							}
							ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
							defer cancel()
							err = worker.InitFunction(ctx, funcName, image, codeURI)
							if err != nil {
								return
							}
							s.funcToWorker[funcName] = append(s.funcToWorker[funcName], worker)
							fmt.Printf("[liu] append func %s to worker %v, funcWorkers: %v\n", funcName, worker, s.funcToWorker[funcName])
							//break 				//commented by liu
							// only run once, and golang range is random. liu: why only run once
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
			//TODO: funcInvokeMetrics now is useless because it will not change scheduler's behavior.
			times, isPresent := s.funcInvokeMetrics[funcName]
			if isPresent == true {
				s.funcInvokeMetrics[funcName] = times + 1
			} else {
				s.funcInvokeMetrics[funcName] = 1
			}

			funcWorkers, isPresent := s.funcToWorker[funcName]
			if isPresent == false || len(funcWorkers) == 0 {
				resCh <- nil
			} else {
				droppoint := rand.Intn(totalBudget)
				currentTotalBias := 0
				var lastWorker *worker.Worker
				for _, targetWorker := range funcWorkers {
					currentTotalBias += s.bias[targetWorker.GetId()]
					if currentTotalBias >= droppoint {
						resCh <- targetWorker
						return
					}
					lastWorker = targetWorker
				}
				resCh <- lastWorker
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
			} else {
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
					ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
					defer cancel()
					err = worker.InitFunction(ctx, funcName, image, codeURI)
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

func (s *Scheduler) getWorkerIndex() string {
	droppoint := rand.Intn(totalBudget)
	currentTotalBias := 0

	//The actual total bias may be less than totalBudget. So we need to record the last worker
	var lastWorker string
	for workerID, bias := range s.bias {
		currentTotalBias += bias
		if currentTotalBias > droppoint {
			return workerID
		}
		lastWorker = workerID
	}
	return lastWorker
}
