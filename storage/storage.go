package storage

import (
	"Assik3_ADP/models"
	"errors"
	"log"
	"sync"
	"time"
)

// Storage represents the in-memory database
type Storage struct {
	Users    map[int]*models.User
	Products map[int]*models.Product
	Carts    map[int]*models.Cart
	Orders   map[int]*models.Order
	mu       sync.RWMutex // Mutex for thread-safe operations
}

var (
	store      *Storage
	userIDSeq  int
	prodIDSeq  int
	orderIDSeq int
)

// OrderProcessingChannel is used for background order processing
var OrderProcessingChannel = make(chan *models.Order, 100)

// InitStorage initializes the in-memory storage with sample data
func InitStorage() {
	store = &Storage{
		Users:    make(map[int]*models.User),
		Products: make(map[int]*models.Product),
		Carts:    make(map[int]*models.Cart),
		Orders:   make(map[int]*models.Order),
	}

	// Seed initial users
	seedUsers()

	// Seed initial products
	seedProducts()

	log.Println("Storage initialized successfully")
}

// seedUsers creates initial user accounts
func seedUsers() {
	users := []*models.User{
		{ID: 1, Email: "student@uni.edu", Password: "password123", Role: "student", CreatedAt: time.Now()},
		{ID: 2, Email: "admin@uni.edu", Password: "admin123", Role: "admin", CreatedAt: time.Now()},
	}

	for _, user := range users {
		store.Users[user.ID] = user
		if user.ID > userIDSeq {
			userIDSeq = user.ID
		}
	}

	log.Printf("Seeded %d users", len(users))
}

// seedProducts creates initial product catalog
func seedProducts() {
	products := []*models.Product{
		{
			ID:          1,
			Name:        "Varsity Jacket '26",
			Category:    "Outerwear",
			Price:       25000,
			InStock:     true,
			Stock:       15,
			Description: "Premium varsity jacket with university logo",
			CreatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "Uni Hoodie Grey",
			Category:    "Hoodies",
			Price:       12000,
			InStock:     true,
			Stock:       30,
			Description: "Comfortable grey hoodie with embroidered logo",
			CreatedAt:   time.Now(),
		},
		{
			ID:          3,
			Name:        "Logo Cap",
			Category:    "Accessories",
			Price:       5000,
			InStock:     true,
			Stock:       50,
			Description: "Baseball cap with university branding",
			CreatedAt:   time.Now(),
		},
		{
			ID:          4,
			Name:        "Classic T-Shirt White",
			Category:    "T-shirts",
			Price:       7000,
			InStock:     true,
			Stock:       40,
			Description: "White cotton t-shirt with small logo",
			CreatedAt:   time.Now(),
		},
		{
			ID:          5,
			Name:        "Bomber Jacket Navy",
			Category:    "Outerwear",
			Price:       28000,
			InStock:     false,
			Stock:       0,
			Description: "Limited edition navy bomber jacket",
			CreatedAt:   time.Now(),
		},
	}

	for _, product := range products {
		store.Products[product.ID] = product
		if product.ID > prodIDSeq {
			prodIDSeq = product.ID
		}
	}

	log.Printf("Seeded %d products", len(products))
}

// User Operations

// CreateUser adds a new user to the system
func CreateUser(email, password, role string) (*models.User, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	// Check if email already exists
	for _, user := range store.Users {
		if user.Email == email {
			return nil, errors.New("email already registered")
		}
	}

	userIDSeq++
	user := &models.User{
		ID:        userIDSeq,
		Email:     email,
		Password:  password,
		Role:      role,
		CreatedAt: time.Now(),
	}

	store.Users[user.ID] = user
	return user, nil
}

// AuthenticateUser validates login credentials
func AuthenticateUser(email, password string) (*models.User, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	for _, user := range store.Users {
		if user.Email == email && user.Password == password {
			return user, nil
		}
	}

	return nil, errors.New("invalid credentials")
}

// Product Operations

// GetAllProducts returns all products in the catalog
func GetAllProducts() []*models.Product {
	store.mu.RLock()
	defer store.mu.RUnlock()

	products := make([]*models.Product, 0, len(store.Products))
	for _, product := range store.Products {
		products = append(products, product)
	}

	return products
}

// GetProductByID retrieves a product by its ID
func GetProductByID(id int) (*models.Product, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	product, exists := store.Products[id]
	if !exists {
		return nil, errors.New("product not found")
	}

	return product, nil
}

// CreateProduct adds a new product to the catalog
func CreateProduct(name, category, description string, price float64, stock int) (*models.Product, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	prodIDSeq++
	product := &models.Product{
		ID:          prodIDSeq,
		Name:        name,
		Category:    category,
		Price:       price,
		InStock:     stock > 0,
		Stock:       stock,
		Description: description,
		CreatedAt:   time.Now(),
	}

	store.Products[product.ID] = product
	return product, nil
}

// FilterProductsByCategory returns products filtered by category
func FilterProductsByCategory(category string) []*models.Product {
	store.mu.RLock()
	defer store.mu.RUnlock()

	products := make([]*models.Product, 0)
	for _, product := range store.Products {
		if product.Category == category {
			products = append(products, product)
		}
	}

	return products
}

// Cart Operations

// AddToCart adds an item to user's cart
func AddToCart(userID, productID, quantity int) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	product, exists := store.Products[productID]
	if !exists {
		return errors.New("product not found")
	}

	if !product.InStock || product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	cart, exists := store.Carts[userID]
	if !exists {
		cart = &models.Cart{
			UserID: userID,
			Items:  []models.CartItem{},
			Total:  0,
		}
		store.Carts[userID] = cart
	}

	// Check if item already in cart
	found := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items[i].Quantity += quantity
			found = true
			break
		}
	}

	if !found {
		cart.Items = append(cart.Items, models.CartItem{
			ProductID: productID,
			Quantity:  quantity,
			Price:     product.Price,
			Name:      product.Name,
		})
	}

	// Recalculate total
	cart.Total = 0
	for _, item := range cart.Items {
		cart.Total += item.Price * float64(item.Quantity)
	}

	return nil
}

// GetCart retrieves user's cart
func GetCart(userID int) (*models.Cart, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	cart, exists := store.Carts[userID]
	if !exists {
		return &models.Cart{
			UserID: userID,
			Items:  []models.CartItem{},
			Total:  0,
		}, nil
	}

	return cart, nil
}

// RemoveFromCart removes an item from user's cart
func RemoveFromCart(userID, productID int) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	cart, exists := store.Carts[userID]
	if !exists {
		return errors.New("cart not found")
	}

	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)

			// Recalculate total
			cart.Total = 0
			for _, item := range cart.Items {
				cart.Total += item.Price * float64(item.Quantity)
			}

			return nil
		}
	}

	return errors.New("item not found in cart")
}

// Order Operations

// CreateOrder creates a new order from user's cart
func CreateOrder(userID int) (*models.Order, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	cart, exists := store.Carts[userID]
	if !exists || len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	orderIDSeq++
	order := &models.Order{
		ID:        orderIDSeq,
		UserID:    userID,
		Items:     make([]models.OrderItem, 0),
		Total:     cart.Total,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Convert cart items to order items
	for _, cartItem := range cart.Items {
		order.Items = append(order.Items, models.OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     cartItem.Price,
			Name:      cartItem.Name,
		})

		// Update product stock
		if product, exists := store.Products[cartItem.ProductID]; exists {
			product.Stock -= cartItem.Quantity
			if product.Stock <= 0 {
				product.InStock = false
				product.Stock = 0
			}
		}
	}

	store.Orders[order.ID] = order

	// Clear cart
	store.Carts[userID] = &models.Cart{
		UserID: userID,
		Items:  []models.CartItem{},
		Total:  0,
	}

	// Send order to processing channel
	go func() {
		OrderProcessingChannel <- order
	}()

	return order, nil
}

// GetUserOrders retrieves all orders for a user
func GetUserOrders(userID int) []*models.Order {
	store.mu.RLock()
	defer store.mu.RUnlock()

	orders := make([]*models.Order, 0)
	for _, order := range store.Orders {
		if order.UserID == userID {
			orders = append(orders, order)
		}
	}

	return orders
}

// UpdateOrderStatus updates the status of an order
func UpdateOrderStatus(orderID int, status string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	order, exists := store.Orders[orderID]
	if !exists {
		return errors.New("order not found")
	}

	order.Status = status
	order.UpdatedAt = time.Now()

	return nil
}

// StartOrderProcessor runs a background goroutine to process orders
func StartOrderProcessor() {
	log.Println("Order processor started")

	for order := range OrderProcessingChannel {
		// Simulate order processing
		log.Printf("Processing order #%d for user #%d", order.ID, order.UserID)

		time.Sleep(2 * time.Second) // Simulate processing time

		// Update order status
		store.mu.Lock()
		if storedOrder, exists := store.Orders[order.ID]; exists {
			storedOrder.Status = "processing"
			storedOrder.UpdatedAt = time.Now()
			log.Printf("Order #%d status updated to 'processing'", order.ID)
		}
		store.mu.Unlock()

		// Simulate further processing
		time.Sleep(3 * time.Second)

		store.mu.Lock()
		if storedOrder, exists := store.Orders[order.ID]; exists {
			storedOrder.Status = "completed"
			storedOrder.UpdatedAt = time.Now()
			log.Printf("Order #%d completed", order.ID)
		}
		store.mu.Unlock()
	}
}
