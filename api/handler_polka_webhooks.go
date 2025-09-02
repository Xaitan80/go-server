package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/xaitan80/go-server/internal/auth"
	"github.com/xaitan80/go-server/internal/config"
	"github.com/xaitan80/go-server/internal/database"
)

// PolkaWebhookRequest represents the shape of incoming webhook requests from Polka
type PolkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

// PolkaWebhooksHandler handles POST /api/polka/webhooks
func PolkaWebhooksHandler(queries *database.Queries, cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
			return
		}

		// Check Polka API key
		key, err := auth.GetAPIKey(r.Header)
		if err != nil || key != cfg.PolkaKey {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid or missing Polka API key"})
			return
		}

		// Decode webhook request
		var req PolkaWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
			return
		}

		// Ignore all events except "user.upgraded"
		if req.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Parse user ID
		userID, err := uuid.Parse(req.Data.UserID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid user ID"})
			return
		}

		// Upgrade user to Chirpy Red
		if err := queries.UpgradeUserToChirpyRed(r.Context(), userID); err != nil {
			// Assume ErrNoRows means user not found
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "User not found"})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to upgrade user"})
			return
		}

		// Success: respond with 204 No Content
		w.WriteHeader(http.StatusNoContent)
	}
}
