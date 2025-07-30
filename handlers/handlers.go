package handlers

import (
	"fmt"
	"net/http"
)

const (
	errInvalidParamPrefix = "invalid query param: "
	errMissingParamPrefix = "missing query param: "
	errMarshallResponse   = "unable to write response"
)

func writeError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = fmt.Fprintf(w, `{"error":"%s"}`, error)
}
