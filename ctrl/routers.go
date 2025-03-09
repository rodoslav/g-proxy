package ctrl

import "net/http"

// Registering Handlers
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", CtrlHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}
