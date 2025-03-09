package main

import (
	"crypto/tls"
	"g-proxy/config"
	"g-proxy/ctrl"
	"g-proxy/proxy"
	"g-proxy/web"
	"log"
	"net/http"
)

func main() {
	// Load config from JSON
	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Loading config file error: %v", err)
	}

	// Proxy server:
	proxySrv := http.Server{
		Addr:         config.ProxyPort,
		Handler:      http.HandlerFunc(proxy.ProxyHandler),
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
	}
	log.Printf("Started https proxy server on port %s", config.ProxyPort)
	go proxySrv.ListenAndServeTLS(config.CertFile, config.KeyFile)

	// Admin console:
	controlSrv := http.Server{
		Addr:         config.AdminPort,
		Handler:      http.HandlerFunc(ctrl.CtrlHandler),
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
	}
	log.Printf("Admin console started on port %s\n", config.AdminPort)
	go controlSrv.ListenAndServeTLS(config.CertFile, config.KeyFile)

	// WEB Server:
	webSrv := http.Server{
		Addr:         config.WebPort,
		Handler:      http.HandlerFunc(web.WebPagesHandler),
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
	}
	log.Printf("Web started on %s\n\n", config.WebPort)
	go webSrv.ListenAndServeTLS(config.CertFile, config.KeyFile)

	select {} // Block Main flow
}
