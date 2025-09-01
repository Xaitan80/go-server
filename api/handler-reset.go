package api

import (
	"encoding/json"
	"net/http"

	"github.com/xaitan80/go-server/internal/database"
)

// ResetHandler deletes all users if platform is "dev"
func ResetHandler(queries *database.Queries, platform string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
			return
		}

		if platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Forbidden"})
			return
		}

		if err := queries.DeleteAllUsers(r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to delete users"})
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"deleted":true}`))
	}
}
