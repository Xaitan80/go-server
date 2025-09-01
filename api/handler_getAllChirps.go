package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/xaitan80/go-server/internal/database"
)

// Response struct for multiple chirps
type chirpsResponse struct {
	ID        string `json:"id"`
	Body      string `json:"body"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetAllChirpsHandler handles GET /api/chirps
func GetAllChirpsHandler(DB *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
			return
		}

		chirps, err := DB.GetAllChirps(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to fetch chirps"})
			return
		}

		// Map SQLC chirps to API response
		resp := make([]chirpsResponse, len(chirps))
		for i, c := range chirps {
			resp[i] = chirpsResponse{
				ID:        c.ID.String(),
				Body:      c.Body,
				UserID:    c.UserID.String(),
				CreatedAt: c.CreatedAt.Format(time.RFC3339),
				UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
