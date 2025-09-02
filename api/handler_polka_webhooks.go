package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/xaitan80/go-server/internal/database"
)

type PolkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

// PolkaWebhooksHandler handles POST /api/polka/webhooks
func PolkaWebhooksHandler(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PolkaWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		// Ignore everything except user.upgraded events
		if req.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Convert user_id string → uuid.UUID
		userID, err := uuid.Parse(req.Data.UserID)
		if err != nil {
			http.Error(w, `{"error":"invalid user_id"}`, http.StatusBadRequest)
			return
		}

		// Try to upgrade the user
		err = queries.UpgradeUserToChirpyRed(r.Context(), userID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
				return
			}
			http.Error(w, `{"error":"failed to upgrade user"}`, http.StatusInternalServerError)
			return
		}

		// Successful upgrade
		w.WriteHeader(http.StatusNoContent)
	}
}
