package api

import (
	"encoding/json"
	"net/http"

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

		// Expect JSON body with token
		var req struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
			return
		}

		if req.Token == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Token is required"})
			return
		}

		// Revoke the refresh token
		if err := queries.RevokeRefreshToken(r.Context(), req.Token); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to revoke token"})
			return
		}

		w.WriteHeader(http.StatusNoContent) // 204 No Content on success
	}
}
