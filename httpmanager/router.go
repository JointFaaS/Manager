package httpmanager

import "net/http"

func (m *Manager) setRouter() {
	m.server.HandleFunc("/create", m.UploadHandler)
	m.server.HandleFunc("/delete", m.DelHandler)
	m.server.HandleFunc("/invoke", m.InvokeHandler)
	m.server.HandleFunc("/register", m.RegisterHandler)
	m.server.HandleFunc("/list", m.ListHandler)
	m.server.HandleFunc("/get", m.GetHandler)
	m.server.HandleFunc("/info", m.InfoHandler)
}

// ListenAndServe starts the Manager main process
func (m *Manager) ListenAndServe() error {
	return http.ListenAndServe(":"+m.port, m.server)
}
