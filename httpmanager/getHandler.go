package httpmanager

import (
	"encoding/json"
	"net/http"
)

// GetHandler returns a function metadata
func (m* Manager) GetHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	w.Write(jsonRet)
	return
}

