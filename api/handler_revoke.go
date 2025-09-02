package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/xaitan80/go-server/internal/auth"
	"github.com/xaitan80/go-server/internal/database"
)

// RevokeHandler handles POST /api/revoke
func RevokeHandler(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
			return
		}

		// Get the refresh token from Authorization header
		tokenStr, err := auth.GetBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Missing or invalid token"})
			return
		}

		now := time.Now()

		// Revoke the token in the database
		err = queries.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
			Token:     tokenStr,
			RevokedAt: sql.NullTime{Time: now, Valid: true}, // wrap time.Time in sql.NullTime
			UpdatedAt: now,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to revoke token"})
			return
		}

		// 204 No Content response
		w.WriteHeader(http.StatusNoContent)
	}
}
