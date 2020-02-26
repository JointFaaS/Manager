package httpmanager

import (
	"encoding/json"
	"net/http"
	"strings"
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
	workerAddr := strings.Split(r.RemoteAddr, ":")[0] + req.WorkerPort
	m.scheduler.RegisterWorker(req.WorkerID, workerAddr)

	return
}

