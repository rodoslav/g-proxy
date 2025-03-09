package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var hopByHopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Connection",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	// log.Printf("Reqwested Sheme: %s", r.URL.Scheme)
	if r.URL.Scheme == "" && r.URL.Port() == "443" {
		r.URL.Scheme = "https"
	}
	log.Printf("Processing: %s", r.URL.String())

	if r.Method == http.MethodConnect {
		handleTunneling(w, r)
	} else {
		handleHTTP(w, r)
	}

	// Створюємо ReverseProxy для пересилання запитів
	//	proxySrv := httputil.NewSingleHostReverseProxy(r.URL)

	// Виконуємо запит через ReverseProxy
	//	proxySrv.ServeHTTP(w, r)

}

func handleTunneling(w http.ResponseWriter, req *http.Request) {

	// Clear hop-by-hop headers
	clearHopHeaders(req.Header)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	// Connection to Client (tunnel)
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		log.Println("Hijack error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client_conn.Close()

	// Connection so Target-server
	log.Println("Connecting to ", req.URL.Host)
	dest_conn, err := net.DialTimeout("tcp", req.URL.Host, 10*time.Second)
	if err != nil {
		log.Println("Connect error:", err)
		writeRawResponse(client_conn, http.StatusServiceUnavailable, req)
		return
	}
	defer dest_conn.Close()

	writeRawResponse(client_conn, http.StatusOK, req)

	log.Println("Transferring:", req.RemoteAddr, "->", req.URL.Host)
	go func() {
		io.Copy(dest_conn, client_conn)
		dest_conn.Close()
	}()
	io.Copy(client_conn, dest_conn)
	log.Println("Done:", req.RemoteAddr, "->", req.URL.Host)

}

func writeRawResponse(conn net.Conn, statusCode int, req *http.Request) {
	if _, err := fmt.Fprintf(
		conn, "HTTP/%d.%d %03d %s\n\n", req.ProtoMajor,
		req.ProtoMinor, statusCode, http.StatusText(statusCode)); err != nil {
		log.Println("Writing response filed:", err)
	}
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	clearHopHeaders(req.Header)
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func clearHopHeaders(head http.Header) {
	for _, hd := range hopByHopHeaders {
		head.Del(hd)
	}
}
