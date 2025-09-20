package util

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON General Response helper
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return nil
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error: could not encode JSON response: %v", err)
		return err
	}

	return nil
}

// JSON Error Response helper
func WriteError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, map[string]string{
		"error": message,
	})
}
