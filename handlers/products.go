package handlers

import (
	"Assik3_ADP/models"
	"encoding/json"
	"net/http"
)

// GetProducts - returns list of product
func GetProducts(w http.ResponseWriter, r *http.Request) {
	// Database simulation
	catalog := []models.Product{
		{ID: 1, Name: "Varsity Jacket '26", Category: "Outerwear", Price: 25000, InStock: true},
		{ID: 2, Name: "Uni Hoodie Grey", Category: "Top Base", Price: 12000, InStock: true},
		{ID: 3, Name: "Logo Cap", Category: "Accessories", Price: 5000, InStock: true},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(catalog)
}

// CreateProduct plug
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": "Product created"}`))
}
