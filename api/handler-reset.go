package api

import (
	"net/http"
	"sync/atomic"
)

// ResetHandler returns an HTTP handler that resets the counter to 0
func ResetHandler(counter *atomic.Int32) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		counter.Store(0) // reset safely
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("Hits reset"))
	}
}
