package web

import "net/http"

// Registering Handlers
func RegisterHandlers() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/config", ConfigHandler)
}
