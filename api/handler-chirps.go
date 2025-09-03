package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xaitan80/go-server/internal/auth"
	"github.com/xaitan80/go-server/internal/database"
)

// List of bad words to filter
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
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
			return
		}

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

		var req chirpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
			return
		}

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

		chirp, err := queries.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   cleaned,
			UserID: userID,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create chirp"})
			return
		}

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

// GetAllChirpsHandler handles GET /api/chirps
// GetAllChirpsHandler handles GET /api/chirps
func GetAllChirpsHandler(DB *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
			return
		}

		sortOrder := r.URL.Query().Get("sort")
		if sortOrder == "" {
			sortOrder = "asc"
		}

		authorIDStr := r.URL.Query().Get("author_id")
		var chirps []database.Chirp
		var err error

		if authorIDStr != "" {
			authorID, parseErr := uuid.Parse(authorIDStr)
			if parseErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid author_id"})
				return
			}
			chirps, err = DB.GetChirpsByAuthorID(r.Context(), authorID)
		} else {
			chirps, err = DB.ListChirps(r.Context())
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to fetch chirps"})
			return
		}

		// Sort chirps by created_at
		sort.Slice(chirps, func(i, j int) bool {
			if sortOrder == "desc" {
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
			}
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})

		resp := make([]ChirpResponse, len(chirps))
		for i, c := range chirps {
			resp[i] = ChirpResponse{
				ID:        c.ID.String(),
				Body:      c.Body,
				UserID:    c.UserID.String(),
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
