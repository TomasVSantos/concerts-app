package handlers

import (
	"encoding/json"
	"net/http"
)

type errorPayload struct {
    Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
    writeJSON(w, status, errorPayload{Error: err.Error()})
}

func readJSON(r *http.Request, v interface{}) error {
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    return dec.Decode(v)
}

func redirect(w http.ResponseWriter, r *http.Request, path string) {
    http.Redirect(w, r, path, http.StatusFound)
}


