package handlers

import (
	"encoding/json"
	"fizzbuzz-code-challenge/services"
	"net/http"
)

// BuildStatsHandler returns a handler that responds with the most frequent request made to the server, if available
func BuildStatsHandler(ch <-chan string) http.Handler {
	type response struct {
		MostFrequent string `json:"most_frequent"`
	}

	s := services.NewStats(ch)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		req, err := s.MostFrequent()
		if err != nil {
			writeError(w, err.Error(), http.StatusNotFound)

			return
		}

		b, err := json.Marshal(response{MostFrequent: req})
		if err != nil {
			writeError(w, errMarshallResponse, http.StatusInternalServerError)

			return
		}

		_, _ = w.Write(b)
	})
}
