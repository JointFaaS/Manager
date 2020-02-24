package httpmanager

import "net/http"

// SetRouter initialise the http handler
func (m * Manager) SetRouter() {
	http.HandleFunc("/createfunction", m.UploadHandler)
	http.HandleFunc("/invoke", m.InvokeHandler)
}