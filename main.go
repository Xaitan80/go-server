package main

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/xaitan80/go-server/api"
	"github.com/xaitan80/go-server/app"
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

	mux.Handle("/admin/metrics", api.HitsHandler(&apiCfg.fileserverHits))
	mux.Handle("/admin/reset", api.ResetHandler(&apiCfg.fileserverHits))

	// health endpoint
	mux.HandleFunc("/api/healthz", api.ReadinessHandler)

	// fileserver with middleware

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(app.FileServerHandler()))

	// api handlers:
	mux.HandleFunc("/api/validate_chirp", api.ValidateChirpHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
