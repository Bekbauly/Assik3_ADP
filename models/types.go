package models

import "time"

// Product - merchandise in the catalog
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	InStock     bool      `json:"in_stock"`
	Description string    `json:"description"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
}

// User - user in the system
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Hidden from JSON responses
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// CartItem - item in users shopping cart
type CartItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Name      string  `json:"name"`
}

// Cart - users shopping cart
type Cart struct {
	UserID int        `json:"user_id"`
	Items  []CartItem `json:"items"`
	Total  float64    `json:"total"`
}

// Order - purchase oreder
type Order struct {
	ID        int         `json:"id"`
	UserID    int         `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// OrderItem - single item in  and order
type OrderItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Name      string  `json:"name"`
}

// LoginRequest - login credentials
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest - registration data
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// AddToCartRequest - request to add items te cart
type AddToCartRequest struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// RemoveFromCartRequest - request to remove item from cart
type RemoveFromCartRequest struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
}

// CreateOrderRequest - request to create new order
type CreateOrderRequest struct {
	UserID int `json:"user_id"`
}

// UpdateOrderStatusRequest - request to update order status
type UpdateOrderStatusRequest struct {
	OrderID int    `json:"user_id"`
	Status  string `json:"status"`
}

// FilterProductsRequest - request to update order status
type FilterProductsRequest struct {
	Category string `json:"category"`
}
