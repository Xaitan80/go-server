package api

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// HitsHandler returns a handler function using the provided config
func HitsHandler(cfg *atomic.Int32) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		hits := cfg.Load()
		fmt.Fprintf(w, "Hits: %d", hits)
	}
}
