package api

import (
	"encoding/json"
	"net/http"

	"strings"

	"github.com/google/uuid"
	"github.com/xaitan80/go-server/internal/database"
)

// Response struct for JSON output
type chirpResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserID    string `json:"user_id"`
}

// GetChirpHandler handles GET /api/chirps/{chirpID}
func GetChirpHandler(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
			return
		}

		// Extract chirpID from the path
		// Expected path: /api/chirps/{chirpID}
		parts := splitPath(r.URL.Path)
		if len(parts) != 3 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid path"})
			return
		}

		chirpID, err := uuid.Parse(parts[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid UUID"})
			return
		}

		// Fetch chirp from database
		chirp, err := queries.GetChirpByID(r.Context(), chirpID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Chirp not found"})
			return
		}

		// Build response with string UUIDs
		resp := chirpResponse{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: chirp.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// Helper function to split path
func splitPath(path string) []string {
	var parts []string
	for _, p := range strings.Split(path, "/") {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}
