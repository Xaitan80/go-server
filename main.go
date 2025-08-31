package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// increment the counter
		cfg.fileserverHits.Add(1)

		// call the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	const port = "8080"

	mux := http.NewServeMux()

	apiCfg := &apiConfig{} // instance

	// hit endpoint

	mux.Handle("/metrics", apiCfg.hitsHandler())
	mux.Handle("/reset", apiCfg.resetHandler())

	// health endpoint
	mux.HandleFunc("/healthz", readinessHandler)

	// fileserver with middleware
	fs := http.FileServer(http.Dir("."))
	handler := http.StripPrefix("/app", fs)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
