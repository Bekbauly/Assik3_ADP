package handlers

import (
	"Assik3_ADP/models"
	"Assik3_ADP/storage"
	"encoding/json"
	"fmt"
	"net/http"
)

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate user
	user, err := storage.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate mock JWT token
	token := fmt.Sprintf("jwt-token-%d-%s", user.ID, user.Role)

	response := map[string]interface{}{
		"token":   token,
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"message": "Login successful",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Default role is "student"
	if req.Role == "" {
		req.Role = "student"
	}

	// Validate role
	if req.Role != "student" && req.Role != "admin" {
		http.Error(w, "Invalid role. Must be 'student' or 'admin'", http.StatusBadRequest)
		return
	}

	// Create user
	user, err := storage.CreateUser(req.Email, req.Password, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	response := map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"message": "Registration successful",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
