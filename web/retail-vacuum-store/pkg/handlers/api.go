package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/pkg/api"
	"github.com/rs/zerolog/log"
)

// SetupAPIRoutes configures all API routes
func SetupAPIRoutes() http.Handler {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/products", handleProducts)
	mux.HandleFunc("/api/products/", handleProductDetail)
	mux.HandleFunc("/api/categories", handleCategories)
	mux.HandleFunc("/api/brands", handleBrands)
	mux.HandleFunc("/api/cart", handleCart)
	mux.HandleFunc("/api/cart/items", handleCartItems)
	mux.HandleFunc("/api/checkout", handleCheckout)
	mux.HandleFunc("/api/orders", handleOrders)
	mux.HandleFunc("/api/customers/profile", handleCustomerProfile)
	mux.HandleFunc("/api/health", handleAPIHealth)

	// Apply CORS and logging middleware
	return corsMiddleware(loggingMiddleware(mux))
}

func handleAPIHealth(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data: map[string]string{
			"status":  "ok",
			"service": "retail-vacuum-store-api",
		},
	})
}

func handleProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse query parameters
	filters := parseProductFilters(r)
	
	// Mock data for now - in real implementation, this would query the database
	products := getMockProducts(filters)
	
	response := api.ProductSearchResponse{
		Products:   products,
		Total:      len(products),
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: (len(products) + filters.Limit - 1) / filters.Limit,
		Categories: getMockCategories(),
		Brands:     getMockBrands(),
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    response,
	})
}

func handleProductDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract product ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/products/")
	productID, err := strconv.Atoi(path)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	// Mock product detail - in real implementation, query database
	product := getMockProductDetail(productID)
	if product == nil {
		writeJSONError(w, http.StatusNotFound, "Product not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    product,
	})
}

func handleCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	categories := getMockCategories()
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    categories,
	})
}

func handleBrands(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	brands := getMockBrands()
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    brands,
	})
}

func handleCart(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetCart(w, r)
	case http.MethodPost:
		handleUpdateCart(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleGetCart(w http.ResponseWriter, r *http.Request) {
	// Get cart from session/cookie - mock for now
	cart := getMockCart()
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    cart,
	})
}

func handleUpdateCart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Mock cart update
	cart := getMockCart()
	cart.Items = append(cart.Items, api.CartItem{
		ProductID:   req.ProductID,
		Quantity:    req.Quantity,
		UnitPrice:   299.99,
		TotalPrice:  299.99 * float64(req.Quantity),
		ProductName: "Premium Vacuum Model",
		ProductSKU:  "VAC-001",
		ImageURL:    "/static/images/vacuum1.jpg",
	})

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    cart,
		Message: "Item added to cart",
	})
}

func handleCartItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		handleUpdateCartItem(w, r)
	case http.MethodDelete:
		handleRemoveCartItem(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleUpdateCartItem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ItemID   int `json:"item_id"`
		Quantity int `json:"quantity"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Mock cart item update
	cart := getMockCart()
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    cart,
		Message: "Cart item updated",
	})
}

func handleRemoveCartItem(w http.ResponseWriter, r *http.Request) {
	itemID := r.URL.Query().Get("item_id")
	if itemID == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing item_id")
		return
	}

	// Mock cart item removal
	cart := getMockCart()
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    cart,
		Message: "Item removed from cart",
	})
}

func handleCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Email           string  `json:"email"`
		ShippingAddress struct {
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			AddressLine1 string `json:"address_line1"`
			City         string `json:"city"`
			State        string `json:"state"`
			PostalCode   string `json:"postal_code"`
		} `json:"shipping_address"`
		PaymentMethod string `json:"payment_method"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Mock order creation
	order := api.Order{
		OrderNumber:    "ORD-20250725-000001",
		GuestEmail:     &req.Email,
		Status:         "confirmed",
		Subtotal:       299.99,
		TaxAmount:      24.00,
		ShippingAmount: 15.99,
		TotalAmount:    339.98,
		PaymentStatus:  "paid",
		ShippingMethod: "standard",
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    order,
		Message: "Order placed successfully",
	})
}

func handleOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Mock order history
	orders := []api.Order{
		{
			OrderNumber:    "ORD-20250720-000001",
			Status:         "delivered",
			TotalAmount:    299.99,
			PaymentStatus:  "paid",
			ShippingMethod: "standard",
		},
		{
			OrderNumber:    "ORD-20250715-000002",
			Status:         "shipped",
			TotalAmount:    449.99,
			PaymentStatus:  "paid",
			ShippingMethod: "express",
		},
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    orders,
	})
}

func handleCustomerProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetProfile(w, r)
	case http.MethodPut:
		handleUpdateProfile(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleGetProfile(w http.ResponseWriter, r *http.Request) {
	// Mock customer profile
	customer := api.Customer{
		ID:             1,
		Email:          "customer@example.com",
		FirstName:      "John",
		LastName:       "Doe",
		Phone:          "555-0123",
		EmailVerified:  true,
		MarketingOptIn: true,
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    customer,
	})
}

func handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req api.Customer
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Mock profile update
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    req,
		Message: "Profile updated successfully",
	})
}

// Helper functions

func parseProductFilters(r *http.Request) api.ProductFilters {
	query := r.URL.Query()
	
	filters := api.ProductFilters{
		Search:  query.Get("search"),
		Page:    1,
		Limit:   12,
		SortBy:  "newest",
	}

	if page := query.Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filters.Page = p
		}
	}

	if limit := query.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			filters.Limit = l
		}
	}

	if categoryID := query.Get("category_id"); categoryID != "" {
		if c, err := strconv.Atoi(categoryID); err == nil {
			filters.CategoryID = &c
		}
	}

	if brandID := query.Get("brand_id"); brandID != "" {
		if b, err := strconv.Atoi(brandID); err == nil {
			filters.BrandID = &b
		}
	}

	if minPrice := query.Get("min_price"); minPrice != "" {
		if mp, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filters.MinPrice = &mp
		}
	}

	if maxPrice := query.Get("max_price"); maxPrice != "" {
		if mp, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filters.MaxPrice = &mp
		}
	}

	if sortBy := query.Get("sort_by"); sortBy != "" {
		filters.SortBy = sortBy
	}

	return filters
}

func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error().Err(err).Msg("Failed to encode JSON response")
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSONResponse(w, status, api.APIResponse{
		Success: false,
		Error:   message,
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
