package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/xaitan80/go-server/internal/auth"
	"github.com/xaitan80/go-server/internal/database"
)

// DeleteChirpHandler handles DELETE /api/chirps/{id}
func DeleteChirpHandler(queries *database.Queries, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract chirp ID from URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid chirp ID"})
			return
		}
		chirpIDStr := parts[3]
		chirpID, err := uuid.Parse(chirpIDStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid chirp ID"})
			return
		}

		// Get user ID from JWT in Authorization header
		userID, err := auth.GetUserIDFromHeader(r.Header, jwtSecret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid or missing token"})
			return
		}

		// Fetch chirp to check ownership
		chirp, err := queries.GetChirpByID(r.Context(), chirpID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Chirp not found"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to fetch chirp"})
			}
			return
		}

		// Check ownership
		if chirp.UserID != userID {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "You are not the author of this chirp"})
			return
		}

		// Delete the chirp
		if err := queries.DeleteChirp(r.Context(), chirpID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to delete chirp"})
			return
		}

		w.WriteHeader(http.StatusNoContent) // 204
	}
}
