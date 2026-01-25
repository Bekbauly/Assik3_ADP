package handlers

import (
	"encoding/json"
	"net/http"
)

// Login plug
func Login(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"token": "mock-jwt-token-12345",
		"user":  "student@uni.edu",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
