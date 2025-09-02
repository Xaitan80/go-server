package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	_ "github.com/google/uuid"
	"github.com/xaitan80/go-server/internal/auth"
	"github.com/xaitan80/go-server/internal/database"
)

// bad words to filter
var badWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

// Request struct for incoming JSON
type chirpRequest struct {
	Body string `json:"body"`
}

// Response struct for JSON
type ChirpResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}

// ChirpsHandler handles POST /api/chirps
func ChirpsHandler(queries *database.Queries, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only accept POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
			return
		}

		// Extract JWT from Authorization header
		tokenString, err := auth.GetBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Missing or invalid token"})
			return
		}

		userID, err := auth.ValidateJWT(tokenString, jwtSecret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid token"})
			return
		}

		// Decode JSON body
		var req chirpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
			return
		}

		// Validate length
		if len(req.Body) > 140 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Chirp is too long"})
			return
		}

		// Filter bad words
		words := strings.Fields(req.Body)
		for i, w := range words {
			lower := strings.ToLower(w)
			for _, bad := range badWords {
				if lower == bad {
					words[i] = "****"
				}
			}
		}
		cleaned := strings.Join(words, " ")

		// Insert into database
		chirp, err := queries.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   cleaned,
			UserID: userID,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create chirp"})
			return
		}

		// Build response
		resp := ChirpResponse{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
