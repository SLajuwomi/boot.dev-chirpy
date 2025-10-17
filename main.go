package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hello world!")
	const filepathRoot = "."
	const port = "8080"
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))
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
