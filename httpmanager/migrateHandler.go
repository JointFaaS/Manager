package httpmanager

import (
	"encoding/json"
	"net/http"
)

type migrateBody struct {
	WorkerID    string `json:"workerID"`
	PromoteBias int    `json:"promoteBias"`
}

// MigrateHandler change a function's bias. TODO: return the bias of all workers to the frontend
func (m *Manager) MigrateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not support method", http.StatusBadRequest)
		return
	}

	var req migrateBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workerID := req.WorkerID
	promoteBias := req.PromoteBias
	biasMap, err := m.scheduler.PromoteBias(workerID, promoteBias)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(biasMap)
	return
}
