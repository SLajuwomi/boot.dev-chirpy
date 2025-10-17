package main

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(rw, r)	
	}) 
}

func customHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) numOfRequests(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetRequests(rw http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}

func main() {
	fmt.Println("Hello world!")
	const filepathRoot = "."
	const port = "8080"
	var apiCfg apiConfig
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("/metrics", apiCfg.numOfRequests)
	mux.HandleFunc("/reset", apiCfg.resetRequests)
	mux.HandleFunc("/healthz", customHandler)
	newServer := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	err := newServer.ListenAndServe()
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
		os.Exit(1)
	}
}
