package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func customHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) numOfRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetRequests(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %v", err)
		w.WriteHeader(500)
		return
	}

	// type failed struct {
	// 	Error string `json:"error"`
	// }

	// type valid struct {
	// 	Valid bool `json:"valid"`
	// }

	type returnVals struct {
		Error string `json:"error"`
		Valid bool   `json:"valid"`
	}

	respBody := returnVals{
		Error: "Something went wrong",
		Valid: false,
	}
	// log.Printf("The text: %v\n the length: %v\n", params.Body, len(params.Body))
	if len(params.Body) > 140 {
		// log.Printf("here")
		respBody.Error = "Chirp is too long"
	} else {
		respBody.Valid = true
	}
	w.Header().Set("Content-Type", "application/json")

	dat, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(400)
		w.Write(dat)
		return
	}
	if respBody.Valid {
		w.WriteHeader(200)

	} else {
		w.WriteHeader(400)
	}
	w.Write(dat)

}
func main() {
	fmt.Println("Hello world!")
	const filepathRoot = "."
	const port = "8080"
	var apiCfg apiConfig
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /admin/metrics", apiCfg.numOfRequests)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetRequests)
	mux.HandleFunc("GET /api/healthz", customHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
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
