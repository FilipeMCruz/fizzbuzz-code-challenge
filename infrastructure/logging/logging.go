package logging

import (
	"log/slog"
	"net/http"
)

// WrapLogging wrap the handler so that all requests passed are logged.
func WrapLogging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Incoming request", "remote address", r.RemoteAddr, "http method", r.Method, "url", r.URL)

		handler.ServeHTTP(w, r)
	})
}
