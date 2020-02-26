package httpmanager

import (
	"encoding/json"
	"net/http"
)

// GetHandler returns a function metadata
func (m* Manager) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not support method", http.StatusBadRequest)
		return
	}
	r.ParseForm()
	funcName := r.FormValue("funcName")

	res, err := m.platformManager.GetFunction(funcName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonRet, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(jsonRet)
	return
}

