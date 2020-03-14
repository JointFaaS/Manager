package httpmanager

import (
	"io/ioutil"
	"net/http"

	"github.com/JointFaaS/Manager/worker"
)

// InvokeHandler invokes a function
func (m* Manager) InvokeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not support method", http.StatusBadRequest)
		return
	}
	r.ParseForm()
	funcName := r.FormValue("funcName")
	args, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Fail to read Payload", http.StatusBadRequest)
		return
	}
	resCh := make(chan *worker.Worker)
	m.scheduler.GetWorker(funcName, resCh)
	worker := <- resCh

	// prom metrics
	fnRequests.WithLabelValues(funcName).Inc()

	if worker == nil {
		res, err := m.platformManager.InvokeFunction(funcName, args)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(res)
	} else {
		res, err := worker.CallFunction(funcName, args)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(res)
	}

	return
}

