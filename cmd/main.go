package main

import (
	"crypto/tls"
	"g-proxy/config"
	"g-proxy/ctrl"
	"g-proxy/proxy"
	"log"
	"net/http"
)

func main() {
	// Завантажуємо конфігурацію з файлу
	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Loading config file error: %v", err)
	}

	// Proxy server:
	proxySrv := http.Server{
		Addr:         config.ProxyPortTLS,
		Handler:      http.HandlerFunc(proxy.ProxyHandler),
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
	}
	log.Printf("Started https proxy server on port %s", config.ProxyPortTLS)
	go proxySrv.ListenAndServeTLS(config.CertFile, config.KeyFile)

	// Admin console:
	adminMux := http.NewServeMux()
	ctrl.RegisterHandlers(adminMux)
	log.Printf("Started web admin console on port %s\n", config.AdminPortTLS)
	go admConsoleSrv(config)

	// WEB Server:
	select {} // Блокуємо основний потік
}

func admConsoleSrv(conf config.Config) {
	err := http.ListenAndServeTLS(
		conf.AdminPortTLS,
		conf.CertFile,
		conf.KeyFile,
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctrl.PagesHandler(w, r)
			}))
	if err != nil {
		log.Printf("TLS Server error: %v", err)
	} else {
		log.Printf("TLS Server started.")
	}
}
