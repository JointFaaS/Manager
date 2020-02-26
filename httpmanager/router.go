package httpmanager

import "net/http"

// SetRouter initialise the http handler
func (m * Manager) SetRouter() {
	http.HandleFunc("/createfunction", m.UploadHandler)
	http.HandleFunc("/invoke", m.InvokeHandler)
	http.HandleFunc("/register", m.RegisterHandler)
	http.HandleFunc("/list", m.ListHandler)
	http.HandleFunc("/get", m.GetHandler)
}