package web

import (
	"encoding/json"
	"net/http"
)

func EncodeJSON(w http.ResponseWriter, data interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func DecodeJSON(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}
