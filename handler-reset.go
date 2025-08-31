package main

import (
	"net/http"
)

// resetHandler resets fileserverHits to 0
func (cfg *apiConfig) resetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Store(0) // safely reset counter
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("Hits reset"))
	}
}
