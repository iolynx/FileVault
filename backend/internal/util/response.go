package util

import (
	"encoding/json"
	"net/http"
)

// JSON General Response helper
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// JSON Error Response helper
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{
		"error": message,
	})
}
