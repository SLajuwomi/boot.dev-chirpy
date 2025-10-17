package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hello world!")
	newServeMux := http.NewServeMux()
	newServer := http.Server{
		Handler: newServeMux,
		Addr:    ":8080",
	}
	err := newServer.ListenAndServe()
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
		os.Exit(1)
	}
}
