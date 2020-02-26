package httpmanager

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/JointFaaS/Manager/worker"
)

type registrationBody struct {
	WorkerPort string `json:"workerPort"`
	WorkerID string `json:"workerID"`
}

// RegisterHandler a new worker registers
func (m* Manager) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not support method", http.StatusBadRequest)
		return
	}
	var req registrationBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, isPresent := m.workers[req.WorkerID]
	if isPresent == false {
		newWorker, err := worker.New(strings.Split(r.RemoteAddr, ":")[0] + req.WorkerPort, req.WorkerID)
		log.Printf("New worker: %s %s %s", r.RemoteAddr, req.WorkerPort, req.WorkerID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		m.workers[req.WorkerID] = newWorker
	}
	return
}

