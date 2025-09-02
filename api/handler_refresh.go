package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/xaitan80/go-server/internal/auth"
	"github.com/xaitan80/go-server/internal/database"
)

// RefreshHandler creates a new access token for a valid refresh token
func RefreshHandler(queries *database.Queries, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract refresh token from Authorization header
		tokenStr, err := auth.GetBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Missing or invalid token"})
			return
		}

		// Look up the token in DB
		rt, err := queries.GetUserFromRefreshToken(r.Context(), tokenStr)
		if err != nil || rt.RevokedAt.Valid || rt.ExpiresAt.Before(time.Now()) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid or expired refresh token"})
			return
		}

		// Generate new JWT access token (expires in 1 hour)
		accessToken, err := auth.MakeJWT(rt.UserID, jwtSecret, time.Hour)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to generate token"})
			return
		}

		// Respond with new access token
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Token string `json:"token"`
		}{Token: accessToken})
	}
}
