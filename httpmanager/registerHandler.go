package httpmanager

import (
	"encoding/json"
	"net/http"
	"strings"
)

type workerRegistrationBody struct {
	WorkerPort string `json:"workerPort"`
	WorkerID   string `json:"workerID"`
}

type workerRegistrationResponseBody struct {
	Region          string `json:"region"`
	JointfaasEnv    string `json:"jointfaasEnv"`
	AccessKeyID     string `json:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret"`
	CenterStorage   string `json:"centerStorage"`
}

// RegisterHandler a new worker registers
func (m *Manager) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not support method", http.StatusBadRequest)
		return
	}
	var req workerRegistrationBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	workerAddr := strings.Split(r.RemoteAddr, ":")[0] + ":" + req.WorkerPort
	m.scheduler.RegisterWorker(req.WorkerID, workerAddr)

	resb, _ := json.Marshal(m.registrationResponse)
	w.Write(resb)

	return
}
