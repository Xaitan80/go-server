package api

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// HitsHandler returns an HTML page showing the current hits
func HitsHandler(cfg *atomic.Int32) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		hits := cfg.Load() // safe read of counter

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		html := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`, hits)

		fmt.Fprint(w, html)
	}
}
