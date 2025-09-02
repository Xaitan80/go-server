package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/xaitan80/go-server/internal/auth"
	"github.com/xaitan80/go-server/internal/database"
)

// Request struct for updating user
type updateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// UpdateUserHandler handles PUT /api/users
func UpdateUserHandler(queries *database.Queries, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow PUT
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
			return
		}

		// Extract user ID from access token
		userID, err := auth.GetUserIDFromHeader(r.Header, jwtSecret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid or missing token"})
			return
		}

		// Decode request body
		var req updateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
			return
		}

		if req.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Email is required"})
			return
		}

		// Prepare hashed password if provided
		var hashedPassword sql.NullString
		if req.Password != "" {
			hash, err := auth.HashPassword(req.Password)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to hash password"})
				return
			}
			hashedPassword = sql.NullString{
				String: hash,
				Valid:  true,
			}
		} else {
			// Keep existing password by setting Valid = false
			hashedPassword = sql.NullString{
				Valid: false,
			}
		}

		// Update user in database
		updatedUser, err := queries.UpdateUser(r.Context(), database.UpdateUserParams{
			ID:             userID,
			Email:          req.Email,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to update user"})
			return
		}

		// Respond with updated user (omit password)
		resp := struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		}{
			ID:        updatedUser.ID.String(),
			Email:     updatedUser.Email,
			CreatedAt: updatedUser.CreatedAt.Format(time.RFC3339),
			UpdatedAt: updatedUser.UpdatedAt.Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
