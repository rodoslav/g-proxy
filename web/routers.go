package web

import "net/http"

// Registering Handlers
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", WebPagesHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}
