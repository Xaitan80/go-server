package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()

	// health endpoint
	mux.HandleFunc("/healthz", readinessHandler)

	// file server endpoint
	mux.Handle("/app/", fileServerHandler())

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
