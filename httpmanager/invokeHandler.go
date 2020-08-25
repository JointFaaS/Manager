package httpmanager

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/JointFaaS/Manager/worker"
)

type invokeBody struct {
	FuncName string `json:"funcName"`
	Args     string `json:"args"`

	// 'true' : may use native serverless, optimize cold-boot
	// 'false' : prevent manager to use native serverless. escape from the limits of specified platform.
	EnableNative string `json:"enableNative"`
}

// InvokeHandler invokes a function
func (m *Manager) InvokeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not support method", http.StatusBadRequest)
		return
	}
	var req invokeBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	funcName := req.FuncName
	args, err := base64.StdEncoding.DecodeString(req.Args)
	if err != nil {
		http.Error(w, "Fail to read Payload", http.StatusBadRequest)
		return
	}
	resCh := make(chan *worker.Worker)

	var worker *worker.Worker = nil
	if req.EnableNative == "true" {
		m.scheduler.GetWorker(funcName, resCh)
		worker = <-resCh
	} else if req.EnableNative == "false" {
		m.scheduler.GetWorkerMust(funcName, resCh)
		worker = <-resCh
		if worker == nil {
			http.Error(w, "No available worker", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Invalid enableNative, it must be 'true' or 'false'", http.StatusBadRequest)
		return
	}

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
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*300)
		defer cancel()
		res, err := worker.CallFunction(ctx, funcName, args)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(res)
	}

	return
}
