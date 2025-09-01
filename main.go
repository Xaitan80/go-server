package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	_ "github.com/lib/pq"
	"github.com/xaitan80/go-server/api"
	"github.com/xaitan80/go-server/app"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

// Middleware that increments the fileserver hit counter
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	const port = "8080"

	// Database setup
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	// Create the serve mux
	mux := http.NewServeMux()

	// Instance to track hits
	apiCfg := &apiConfig{}

	// Admin endpoints
	mux.Handle("/admin/metrics", api.HitsHandler(&apiCfg.fileserverHits))
	mux.Handle("/admin/reset", api.ResetHandler(&apiCfg.fileserverHits))

	// Health endpoint
	mux.HandleFunc("/api/healthz", api.ReadinessHandler)

	// Fileserver with middleware
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(app.FileServerHandler()))

	// API endpoint: chirp validation
	mux.HandleFunc("/api/validate_chirp", api.ValidateChirpHandler)

	// Start the server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
