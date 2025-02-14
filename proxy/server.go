package proxy

import (
	"g-proxy/config"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server struct {
	config config.Config
}

func NewServer(config config.Config) *Server {
	// Reserving memory to Server struct & init fields
	s := &Server{
		config: config,
	}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Отримуємо URL цільового сервера з запиту клієнта
	targetURL, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Створюємо зворотний проксі для кожного запиту
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Перенаправляємо запит
	proxy.ServeHTTP(w, r)

}
