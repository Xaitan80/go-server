package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	// Serve files from the current directory
	fs := http.FileServer(http.Dir("."))

	// Handle root URL by serving files
	mux := http.NewServeMux()
	mux.Handle("/", fs)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
