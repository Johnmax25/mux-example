package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

var guitars = map[string]string{
	"Fender": "$1000",
	"Gibson": "$2000",
	"Ibanez": "$3000",
	"Squier": "$1899",
	
}

func main() {
	isReady := &atomic.Value{}
	isReady.Store(false)
	go func() {
		log.Printf("Readyz probe is negative by default...")
		time.Sleep(10 * time.Second)
		isReady.Store(true)
		log.Printf("Readyz probe is positive.")
	}()

	r := mux.NewRouter()
	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/healthz", healthz)
	r.HandleFunc("/readyz", readyz(isReady))

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><head><meta name=\"WAS\" content=\"d10fd3a22bec0058e0b36df7d607385a\"></head><body><h1>John's Guitar Store</h1>"))
	w.Write([]byte("<h2><i>- Currently on Sale:</i></h2><ul>"))
	for k, v := range guitars {
		w.Write([]byte(fmt.Sprintf("<li>%s: %s</li>", k, v)))
	}
	w.Write([]byte("</ul></body></html>"))
}
// healthz is a liveness probe.
func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	
}
// readyz is a readiness probe.
func readyz(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}