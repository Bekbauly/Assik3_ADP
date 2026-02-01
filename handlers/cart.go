package handlers

import (
	"Assik3_ADP/models"
	"Assik3_ADP/storage"
	"encoding/json"
	"fmt"
	"net/http"
)

// AddToCart adds a product to user's shopping cart
func AddToCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.UserID <= 0 || req.ProductID <= 0 || req.Quantity <= 0 {
		http.Error(w, "Invalid user ID, product ID, or quantity", http.StatusBadRequest)
		return
	}

	// Add to cart
	if err := storage.AddToCart(req.UserID, req.ProductID, req.Quantity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get updated cart
	cart, _ := storage.GetCart(req.UserID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Item added to cart",
		"cart":    cart,
	})
}

// ViewCart retrieves user's current shopping cart
func ViewCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from query parameter
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var userID int
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cart, err := storage.GetCart(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// RemoveFromCart removes a product from user's cart
func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RemoveFromCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.UserID <= 0 || req.ProductID <= 0 {
		http.Error(w, "Invalid user ID or product ID", http.StatusBadRequest)
		return
	}

	// Remove from cart
	if err := storage.RemoveFromCart(req.UserID, req.ProductID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get updated cart
	cart, _ := storage.GetCart(req.UserID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Item removed from cart",
		"cart":    cart,
	})
}
