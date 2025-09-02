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
	"github.com/xaitan80/go-server/internal/config"
	"github.com/xaitan80/go-server/internal/database"
)

// Middleware that increments the fileserver hit counter
func middlewareMetricsInc(hits *atomic.Int32, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
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

	// Global API config
	apiCfg := &config.APIConfig{
		Platform:  os.Getenv("PLATFORM"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		PolkaKey:  os.Getenv("POLKA_KEY"),
	}

	// Fileserver hit counter
	var fileserverHits atomic.Int32

	mux := http.NewServeMux()

	// --- Admin Endpoints ---
	mux.Handle("/admin/metrics", api.HitsHandler(&fileserverHits))
	mux.HandleFunc("/admin/reset", api.ResetHandler(queries, apiCfg.Platform))

	// --- Health Endpoint ---
	mux.HandleFunc("/api/healthz", api.ReadinessHandler)

	// --- Fileserver ---
	mux.Handle("/app/", middlewareMetricsInc(&fileserverHits, app.FileServerHandler()))

	// --- API Endpoints ---
	// /api/chirps handles GET (all) and POST (create)
	mux.HandleFunc("/api/chirps", methodHandler(map[string]http.HandlerFunc{
		http.MethodPost: api.ChirpsHandler(queries, apiCfg.JWTSecret),
		http.MethodGet:  api.GetAllChirpsHandler(queries),
	}))

	// /api/chirps/{id} for GET single chirp and DELETE chirp
	mux.HandleFunc("/api/chirps/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.GetChirpHandler(queries)(w, r)
		case http.MethodDelete:
			api.DeleteChirpHandler(queries, apiCfg.JWTSecret)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Method not allowed"})
		}
	})

	// /api/users handles POST (create) and PUT (update)
	mux.HandleFunc("/api/users", methodHandler(map[string]http.HandlerFunc{
		http.MethodPost: api.CreateUserHandler(queries),
		http.MethodPut:  api.UpdateUserHandler(queries, apiCfg.JWTSecret),
	}))

	// /api/login
	mux.HandleFunc("/api/login", api.LoginHandler(queries, apiCfg.JWTSecret))

	// refresh and revoke
	mux.HandleFunc("/api/refresh", api.RefreshHandler(queries, apiCfg.JWTSecret))
	mux.HandleFunc("/api/revoke", api.RevokeHandler(queries))

	// /api/polka/webhooks
	mux.HandleFunc("/api/polka/webhooks", api.PolkaWebhooksHandler(queries, apiCfg))

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
