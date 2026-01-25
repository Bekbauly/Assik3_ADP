package main

import (
	"fmt"
	"net/http"
	"uni-merch/handlers"
)

func main() {
	// Маршруты
	http.HandleFunc("/products", handlers.GetProducts)
	http.HandleFunc("/login", handlers.Login)

	fmt.Println("Server running on http://localhost:8080")

	// Запуск
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error:", err)
	}
}
