package main

import (
	"fmt"
	"net/http"
)

// hitsHandler is now a method on *apiConfig
func (cfg *apiConfig) hitsHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		hits := cfg.fileserverHits.Load() // safe read
		fmt.Fprintf(w, "Hits: %d", hits)  // write as plain text

	}
}
