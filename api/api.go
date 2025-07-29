package api

import (
	"encoding/json"
	"fizzbuzz-code-challenge/fizzbuzz"
	"fmt"
	"net/http"
	"strconv"
)

const (
	queryParamInt1        = "int1"
	queryParamInt2        = "int2"
	queryParamLimit       = "limit"
	queryParamStr1        = "str1"
	queryParamStr2        = "str2"
	errInvalidParamPrefix = "invalid query param: "
	errMissingParamPrefix = "missing query param: "
	errMarshallResponse   = "unable to write response"
)

// BuildFizzBuzzHandler returns a handler that validates the request, calls the fizzbuzz function and
// returns the content in json
func BuildFizzBuzzHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		int1, err := strconv.Atoi(r.URL.Query().Get(queryParamInt1))
		if err != nil {
			writeError(w, errInvalidParamPrefix+queryParamInt1, http.StatusBadRequest)

			return
		}

		int2, err := strconv.Atoi(r.URL.Query().Get(queryParamInt2))
		if err != nil {
			writeError(w, errInvalidParamPrefix+queryParamInt2, http.StatusBadRequest)

			return
		}

		limit, err := strconv.Atoi(r.URL.Query().Get(queryParamLimit))
		if err != nil {
			writeError(w, errInvalidParamPrefix+queryParamLimit, http.StatusBadRequest)

			return
		}

		str1 := r.URL.Query().Get(queryParamStr1)
		if str1 == "" && !r.URL.Query().Has(queryParamStr1) {
			writeError(w, errMissingParamPrefix+queryParamStr1, http.StatusBadRequest)

			return
		}

		str2 := r.URL.Query().Get(queryParamStr2)
		if str2 == "" && !r.URL.Query().Has(queryParamStr2) {
			writeError(w, errMissingParamPrefix+queryParamStr2, http.StatusBadRequest)

			return
		}

		result, err := fizzbuzz.FizzBuzz(int1, int2, limit, str1, str2)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)

			return
		}

		b, err := json.Marshal(result)
		if err != nil {
			writeError(w, errMarshallResponse, http.StatusInternalServerError)

			return
		}

		_, _ = w.Write(b)
	})
}

func writeError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = fmt.Fprintf(w, `{"error":"%s"}`, error)
}
