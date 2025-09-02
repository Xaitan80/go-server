package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/xaitan80/go-server/internal/auth"
	"github.com/xaitan80/go-server/internal/database"
)

// Token TTL for access tokens
const accessTokenTTL = time.Hour * 24

// Request struct for login
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Response struct for login
type loginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

// LoginHandler handles POST /api/login
func LoginHandler(queries *database.Queries, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
			return
		}

		// Decode request
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
			return
		}

		// Get user by email
		user, err := queries.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid credentials"})
			return
		}

		// Check password hash
		if err := auth.CheckPasswordHash(req.Password, user.HashedPassword.String); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid credentials"})
			return
		}

		// Generate JWT access token
		accessToken, err := auth.MakeJWT(user.ID, jwtSecret, accessTokenTTL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to generate token"})
			return
		}

		// Return token + email
		resp := loginResponse{
			Email: user.Email,
			Token: accessToken,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
