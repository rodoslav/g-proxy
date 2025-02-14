package web

import "net/http"

// Make handlers for all html page

// For index.html
func HomeHandler(w http.ResponseWriter, r *http.Request) {}

// For login.html
func LoginHandler(w http.ResponseWriter, r *http.Request) {}

// For config.html
func ConfigHandler(w http.ResponseWriter, r *http.Request) {}
