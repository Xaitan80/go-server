package main

import (
	"database/sql"
	"encoding/json"
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

// helper for method-based routing
func methodHandler(handlers map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h, ok := handlers[r.Method]; ok {
			h(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Method not allowed"})
		}
	}
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

	queries := database.New(db)
	apiCfg := &apiConfig{Platform: os.Getenv("PLATFORM")}

	mux := http.NewServeMux()

	// --- Admin Endpoints ---
	mux.Handle("/admin/metrics", api.HitsHandler(&apiCfg.fileserverHits))
	mux.HandleFunc("/admin/reset", api.ResetHandler(queries, apiCfg.Platform))

	// --- Health Endpoint ---
	mux.HandleFunc("/api/healthz", api.ReadinessHandler)

	// --- Fileserver ---
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(app.FileServerHandler()))

	// --- API Endpoints ---
	// /api/chirps handles both GET (all) and POST (create)
	mux.HandleFunc("/api/chirps", methodHandler(map[string]http.HandlerFunc{
		http.MethodPost: api.ChirpsHandler(queries),
		http.MethodGet:  api.GetAllChirpsHandler(queries),
	}))

	// /api/chirps/{id} for single chirp
	mux.HandleFunc("/api/chirps/", api.GetChirpHandler(queries))

	// /api/users
	mux.HandleFunc("/api/users", api.CreateUserHandler(queries))
	mux.HandleFunc("/api/login", api.LoginHandler(queries))

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
