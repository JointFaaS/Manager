package httpmanager

import (
	"encoding/json"
	"net/http"
)

type priceInputBody struct {
}

type priceOutputBody struct {
	Used        int32 `json:"used"`
	Total       int32 `json:"total"`
	UnitRequest int32 `json:"unitRequest"`
	UnitPrice   int32 `json:"unitPrice"`
}

// PriceHandler invokes a function
func (m *Manager) PriceHandler(w http.ResponseWriter, r *http.Request) {
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
	res := priceOutputBody{}
	jsonRet, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(jsonRet)

	return
}
