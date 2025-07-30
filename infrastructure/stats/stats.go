package stats

import "net/http"

// BuildWrapStats returns a function that wraps a generic handler, so that all request urls are published to be processed later
func BuildWrapStats(ch chan<- string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ch <- r.URL.String()
			next.ServeHTTP(w, r)
		})
	}
}
