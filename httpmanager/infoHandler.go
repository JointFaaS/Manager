package httpmanager

import (
	"encoding/json"
	"net/http"
)

type infoInputBody struct {
}

type infoOutputBody struct {
	Used        int32 `json:"used"`
	Total       int32 `json:"total"`
	UnitRequest int32 `json:"unitRequest"`
	UnitPrice   int32 `json:"unitPrice"`
}

// InfoHandler returns the basic info of the cloud
func (m *Manager) InfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not support method", http.StatusBadRequest)
		return
	}
	// var req priceInputBody
	// err := json.NewDecoder(r.Body).Decode(&req)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	res := infoOutputBody{}
	jsonRet, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	w.Write(jsonRet)

	return
}
