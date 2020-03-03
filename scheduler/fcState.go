package scheduler

import "github.com/JointFaaS/Manager/worker"

// this file is about state refreshing
// When a function is removed or updated, the scheduler should
// stop the function instances which has the old version.
// Futhermore, it's better if we can start new version and then stop the old.

// DeleteFunction is invoked when the user deletes a function
func (s *Scheduler) DeleteFunction(funcName string) error {
	resCh := make(chan error)
	s.tasks <- &scheduleTask{
		f: func() {
			funcWorkers, isPresent := s.funcToWorker[funcName]
			if isPresent == false || len(funcWorkers) == 0 {
				resCh <- nil
			} else {
				s.funcToWorker[funcName] = make([]*worker.Worker, 0)
				resCh <- nil
				// TODO
				// Do we need to inform the workers? or will the worker release these resources after a while?
			}
		},
	}
	err := <-resCh
	return err
}
