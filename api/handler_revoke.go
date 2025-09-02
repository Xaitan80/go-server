package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/xaitan80/go-server/internal/auth"
	"github.com/xaitan80/go-server/internal/database"
)

// RevokeHandler revokes a refresh token
func RevokeHandler(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract refresh token from header

		tokenStr, err := auth.GetBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Missing or invalid token"})
			return
		}

		// Revoke the refresh token
		params := database.RevokeRefreshTokenParams{
			Token:     tokenStr,
			RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
		}

		if err := queries.RevokeRefreshToken(r.Context(), params); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to revoke token"})
			return
		}

		// Successful revoke returns 204 No Content
		w.WriteHeader(http.StatusNoContent)
	}
}
