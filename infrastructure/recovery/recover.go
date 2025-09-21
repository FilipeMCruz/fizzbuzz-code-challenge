package recovery

import (
	"log/slog"
	"net/http"
)

// WrapRecovery wraps the handler so that all panics can be recovered from, a 500 code is returned.
func WrapRecovery(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				slog.Warn("Recovered from panic", "error", r)

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
