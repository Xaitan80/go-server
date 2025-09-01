package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/xaitan80/go-server/api"
	"github.com/xaitan80/go-server/app"
	"github.com/xaitan80/go-server/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	Platform       string
}

// Middleware that increments the fileserver hit counter
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on OS environment variables")
	}

	const port = "8080"

	// Database setup
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	// Create SQLC queries instance
	queries := database.New(db)

	// Create the serve mux
	mux := http.NewServeMux()

	// Instance to track hits and platform
	apiCfg := &apiConfig{
		Platform: os.Getenv("PLATFORM"),
	}

	// Admin endpoints
	mux.Handle("/admin/metrics", api.HitsHandler(&apiCfg.fileserverHits))
	mux.HandleFunc("/admin/reset", api.ResetHandler(queries, apiCfg.Platform))

	// Health endpoint
	mux.HandleFunc("/api/healthz", api.ReadinessHandler)

	// Fileserver with middleware
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(app.FileServerHandler()))

	// API endpoints
	mux.HandleFunc("/api/validate_chirp", api.ValidateChirpHandler)
	mux.HandleFunc("/api/users", api.CreateUserHandler(queries))

	// Start the server
	srv := &http.Server{
		Addr: ":" + port,

		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
