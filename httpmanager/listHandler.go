package httpmanager

import (
	"encoding/json"
	"net/http"
)

// ListHandler A new worker registers
func (m* Manager) ListHandler(w http.ResponseWriter, r *http.Request) {
	ret, err := m.platformManager.ListFunction()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonRet, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(jsonRet)
	return
}

