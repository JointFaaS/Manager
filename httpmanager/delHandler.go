package httpmanager

import (
	"net/http"
)

// DelHandler deletes a function
func (m *Manager) DelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not support method", http.StatusBadRequest)
		return
	}
	r.ParseForm()
	funcName := r.FormValue("funcName")

	err := m.platformManager.DeleteFunction(funcName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = m.scheduler.DeleteFunction(funcName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Delete Successful"))

	return
}
