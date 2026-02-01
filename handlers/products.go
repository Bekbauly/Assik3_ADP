package handlers

import (
	"Assik3_ADP/models"
	"Assik3_ADP/storage"
	"encoding/json"
	"net/http"
)

// GetProducts returns the complete product catalog
func GetProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	products := storage.GetAllProducts()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"products": products,
		"count":    len(products),
	})
}

// CreateProduct adds a new product to the catalog (Admin only)
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name        string  `json:"name"`
		Category    string  `json:"category"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Name == "" || req.Category == "" || req.Price <= 0 {
		http.Error(w, "Name, category, and price are required", http.StatusBadRequest)
		return
	}

	// Validate category
	validCategories := map[string]bool{
		"Hoodies":         true,
		"Varsity Jackets": true,
		"T-shirts":        true,
		"Accessories":     true,
		"Outerwear":       true,
	}

	if !validCategories[req.Category] {
		http.Error(w, "Invalid category", http.StatusBadRequest)
		return
	}

	product, err := storage.CreateProduct(req.Name, req.Category, req.Description, req.Price, req.Stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"product": product,
		"message": "Product created successfully",
	})
}

// FilterProducts filters products by category
func FilterProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.FilterProductsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Category == "" {
		http.Error(w, "Category is required", http.StatusBadRequest)
		return
	}

	products := storage.FilterProductsByCategory(req.Category)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"products": products,
		"count":    len(products),
		"category": req.Category,
	})
}
