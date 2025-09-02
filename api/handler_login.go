package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/xaitan80/go-server/internal/auth"
	"github.com/xaitan80/go-server/internal/database"
)

// Request struct for login
type loginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int64  `json:"expires_in_seconds,omitempty"`
}

// Response struct for returning the user (no password)
type loginResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
}

// LoginHandler handles POST /api/login
func LoginHandler(queries *database.Queries, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
			return
		}

		// Decode request body
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
			return
		}

		if req.Email == "" || req.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Email and password are required"})
			return
		}

		// Fetch the user by email
		user, err := queries.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Incorrect email or password"})
			return
		}

		// Compare password with stored hash
		if err := auth.CheckPasswordHash(req.Password, user.HashedPassword); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Incorrect email or password"})
			return
		}

		// --- JWT creation ---
		expires := time.Hour // default 1 hour
		if req.ExpiresInSeconds > 0 {
			if req.ExpiresInSeconds > 3600 {
				req.ExpiresInSeconds = 3600
			}
			expires = time.Duration(req.ExpiresInSeconds) * time.Second
		}

		token, err := auth.MakeJWT(user.ID, jwtSecret, expires)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to generate token"})
			return
		}

		// --- Refresh token creation ---
		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to generate refresh token"})
			return
		}

		// Save refresh token in DB
		_, err = queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			UserID:    user.ID,
			Token:     refreshToken,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to save refresh token"})
			return
		}

		// Build response including JWT and refresh token
		resp := loginResponse{
			ID:           user.ID.String(),
			Email:        user.Email,
			CreatedAt:    user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
			Token:        token,
			RefreshToken: refreshToken,
			IsChirpyRed:  user.IsChirpyRed,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
