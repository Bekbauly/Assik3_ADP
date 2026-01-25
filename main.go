package main

import (
	"Assik3_ADP/handlers"
	"fmt"
	"net/http"
)

func main() {
	// routes
	http.HandleFunc("/products", handlers.GetProducts)
	http.HandleFunc("/login", handlers.Login)

	fmt.Println("Server running on http://localhost:8080")

	// run
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error:", err)
	}
}
