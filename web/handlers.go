package web

import (
	// "g-proxy/config"
	"html/template"
	"net/http"
)

// Make handlers for all html pages
func WebPagesHandler(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/": // For index.html
		renderTemplate(w, "index.html")
	case "/index.html":
		renderTemplate(w, "index.html")
	case "/login": // For login.html + LoginHandler
		renderTemplate(w, "login.html")
	case "/auth": // Authetication auth.html + AuthHandler
		renderTemplate(w, "auth.html")
	default:
		http.NotFound(w, req)
	}
}

// Функція для рендерингу шаблонів
func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("web/templates/" + tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Отримуємо статистику з proxy.Server
// config := r.Context().Value("config").(config.Config)
