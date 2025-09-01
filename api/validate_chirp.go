package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

// bad words to filter
var badWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

// Request struct for incoming JSON
type chirpRequest struct {
	Body string `json:"body"`
}

// Response struct
type chirpResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

// ValidateChirpHandler handles POST /api/validate_chirp
func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		return
	}

	// Decode JSON body
	var req chirpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		return
	}

	// Validate length
	if len(req.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Chirp is too long"})
		return
	}

	// Split into words and filter bad words
	words := strings.Fields(req.Body)
	for i, w := range words {
		lower := strings.ToLower(w)
		for _, bad := range badWords {
			if lower == bad {
				words[i] = "****"
			}
		}
	}

	cleaned := strings.Join(words, " ")

	// Return cleaned text as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirpResponse{CleanedBody: cleaned})
}
