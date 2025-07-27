package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/pkg/api"
)

// SetupAdminAPIRoutes configures admin-specific API routes
func SetupAdminAPIRoutes() http.Handler {
	mux := http.NewServeMux()

	// Admin authentication and session management
	mux.HandleFunc("/api/admin/login", handleAdminLogin)
	mux.HandleFunc("/api/admin/logout", handleAdminLogout)
	mux.HandleFunc("/api/admin/profile", handleAdminProfile)

	// Product management
	mux.HandleFunc("/api/admin/products", handleAdminProducts)
	mux.HandleFunc("/api/admin/products/", handleAdminProductDetail)
	mux.HandleFunc("/api/admin/products/bulk", handleAdminProductsBulk)

	// Order management
	mux.HandleFunc("/api/admin/orders", handleAdminOrders)
	mux.HandleFunc("/api/admin/orders/", handleAdminOrderDetail)
	mux.HandleFunc("/api/admin/orders/status", handleAdminOrderStatus)

	// Customer management
	mux.HandleFunc("/api/admin/customers", handleAdminCustomers)
	mux.HandleFunc("/api/admin/customers/", handleAdminCustomerDetail)

	// Inventory management
	mux.HandleFunc("/api/admin/inventory", handleAdminInventory)
	mux.HandleFunc("/api/admin/inventory/", handleAdminInventoryDetail)
	mux.HandleFunc("/api/admin/inventory/movements", handleAdminInventoryMovements)
	mux.HandleFunc("/api/admin/inventory/low-stock", handleAdminLowStock)

	// Analytics and reporting
	mux.HandleFunc("/api/admin/analytics/dashboard", handleAdminDashboard)
	mux.HandleFunc("/api/admin/analytics/sales", handleAdminSalesAnalytics)
	mux.HandleFunc("/api/admin/analytics/products", handleAdminProductAnalytics)

	// Admin user management
	mux.HandleFunc("/api/admin/users", handleAdminUsers)
	mux.HandleFunc("/api/admin/users/", handleAdminUserDetail)

	// Apply admin middleware (authentication, CORS, logging)
	return adminMiddleware(corsMiddleware(loggingMiddleware(mux)))
}

// Admin Authentication
func handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Mock authentication - in real implementation, verify against database
	if req.Username == "admin" && req.Password == "admin123" {
		// Generate session token
		token := generateSessionToken()
		
		// Set secure cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "admin_session",
			Value:    token,
			HttpOnly: true,
			Secure:   false, // Set to true in production with HTTPS
			SameSite: http.SameSiteStrictMode,
			MaxAge:   86400, // 24 hours
		})

		writeJSONResponse(w, http.StatusOK, api.APIResponse{
			Success: true,
			Data: map[string]interface{}{
				"user": api.AdminUser{
					ID:        1,
					Username:  "admin",
					Email:     "admin@vacuummart.com",
					FirstName: "Admin",
					LastName:  "User",
					Role:      "admin",
					IsActive:  true,
				},
				"token": token,
			},
			Message: "Login successful",
		})
	} else {
		writeJSONError(w, http.StatusUnauthorized, "Invalid credentials")
	}
}

func handleAdminLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}

// Product Management
func handleAdminProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetAdminProducts(w, r)
	case http.MethodPost:
		handleCreateProduct(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleGetAdminProducts(w http.ResponseWriter, r *http.Request) {
	filters := parseAdminProductFilters(r)
	
	// Mock data - in real implementation, query database
	products := getMockAdminProducts(filters)
	
	response := api.AdminProductListResponse{
		Products:   products,
		Total:      len(products),
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: (len(products) + filters.Limit - 1) / filters.Limit,
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    response,
	})
}

func handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req api.AdminProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" || req.SKU == "" || req.BasePrice == 0 {
		writeJSONError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Mock product creation
	product := api.AdminProduct{
		ID:               1,
		SKU:              req.SKU,
		Name:             req.Name,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		CategoryID:       req.CategoryID,
		BrandID:          req.BrandID,
		ProductType:      req.ProductType,
		BasePrice:        req.BasePrice,
		SalePrice:        req.SalePrice,
		Weight:           req.Weight,
		IsActive:         req.IsActive,
		IsFeatured:       req.IsFeatured,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	writeJSONResponse(w, http.StatusCreated, api.APIResponse{
		Success: true,
		Data:    product,
		Message: "Product created successfully",
	})
}

func handleAdminProductDetail(w http.ResponseWriter, r *http.Request) {
	// Extract product ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/products/")
	productID, err := strconv.Atoi(path)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetAdminProductDetail(w, r, productID)
	case http.MethodPut:
		handleUpdateProduct(w, r, productID)
	case http.MethodDelete:
		handleDeleteProduct(w, r, productID)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleGetAdminProductDetail(w http.ResponseWriter, r *http.Request, productID int) {
	// Mock product detail
	product := getMockAdminProductDetail(productID)
	if product == nil {
		writeJSONError(w, http.StatusNotFound, "Product not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    product,
	})
}

// Order Management
func handleAdminOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	filters := parseAdminOrderFilters(r)
	orders := getMockAdminOrders(filters)
	
	response := api.AdminOrderListResponse{
		Orders:     orders,
		Total:      len(orders),
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: (len(orders) + filters.Limit - 1) / filters.Limit,
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    response,
	})
}

func handleAdminOrderDetail(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/orders/")
	orderID, err := strconv.Atoi(path)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetAdminOrderDetail(w, r, orderID)
	case http.MethodPut:
		handleUpdateOrder(w, r, orderID)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleAdminOrderStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		OrderID       int    `json:"order_id"`
		Status        string `json:"status"`
		Notes         string `json:"notes"`
		TrackingNumber string `json:"tracking_number,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate status
	validStatuses := []string{"pending", "confirmed", "processing", "shipped", "delivered", "cancelled", "refunded"}
	isValid := false
	for _, status := range validStatuses {
		if req.Status == status {
			isValid = true
			break
		}
	}

	if !isValid {
		writeJSONError(w, http.StatusBadRequest, "Invalid order status")
		return
	}

	// Mock order status update
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Message: fmt.Sprintf("Order %d status updated to %s", req.OrderID, req.Status),
	})
}

// Dashboard Analytics
func handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	dashboard := api.AdminDashboard{
		Stats: api.DashboardStats{
			TodaysSales:     24573.50,
			TodaysOrders:    142,
			ProductsInStock: 1247,
			ActiveCustomers: 8934,
			LowStockItems:   23,
		},
		RecentOrders: getMockRecentOrders(),
		TopProducts:  getMockTopProducts(),
		SalesChart:   getMockSalesChart(),
		LowStockAlerts: getMockLowStockAlerts(),
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    dashboard,
	})
}

// Inventory Management
func handleAdminInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	filters := parseInventoryFilters(r)
	inventory := getMockInventory(filters)
	
	response := api.InventoryListResponse{
		Items:      inventory,
		Total:      len(inventory),
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: (len(inventory) + filters.Limit - 1) / filters.Limit,
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    response,
	})
}

func handleAdminLowStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	lowStockItems := getMockLowStockItems()
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    lowStockItems,
	})
}

// Customer Management
func handleAdminCustomers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	filters := parseCustomerFilters(r)
	customers := getMockAdminCustomers(filters)
	
	response := api.AdminCustomerListResponse{
		Customers:  customers,
		Total:      len(customers),
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: (len(customers) + filters.Limit - 1) / filters.Limit,
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    response,
	})
}

// Helper functions
func adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for login endpoint
		if r.URL.Path == "/api/admin/login" {
			next.ServeHTTP(w, r)
			return
		}

		// Check for admin session cookie
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Value == "" {
			writeJSONError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// In real implementation, validate session token
		if !isValidSessionToken(cookie.Value) {
			writeJSONError(w, http.StatusUnauthorized, "Invalid session")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func generateSessionToken() string {
	// Mock token generation - use proper JWT or secure random in production
	return fmt.Sprintf("admin_session_%d", time.Now().Unix())
}

func isValidSessionToken(token string) bool {
	// Mock validation - implement proper session validation in production
	return strings.HasPrefix(token, "admin_session_")
}

func parseAdminProductFilters(r *http.Request) api.AdminProductFilters {
	query := r.URL.Query()
	
	filters := api.AdminProductFilters{
		Search:  query.Get("search"),
		Page:    1,
		Limit:   20,
		SortBy:  "updated_at",
		SortDir: "desc",
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

	if status := query.Get("status"); status != "" {
		filters.Status = &status
	}

	if sortBy := query.Get("sort_by"); sortBy != "" {
		filters.SortBy = sortBy
	}

	if sortDir := query.Get("sort_dir"); sortDir != "" {
		filters.SortDir = sortDir
	}

	return filters
}

func parseAdminOrderFilters(r *http.Request) api.AdminOrderFilters {
	query := r.URL.Query()
	
	filters := api.AdminOrderFilters{
		Page:    1,
		Limit:   20,
		SortBy:  "created_at",
		SortDir: "desc",
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

	if status := query.Get("status"); status != "" {
		filters.Status = &status
	}

	if customerID := query.Get("customer_id"); customerID != "" {
		if c, err := strconv.Atoi(customerID); err == nil {
			filters.CustomerID = &c
		}
	}

	if orderNumber := query.Get("order_number"); orderNumber != "" {
		filters.OrderNumber = &orderNumber
	}

	return filters
}

// Additional handler functions

func handleAdminProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Mock admin profile
		profile := api.AdminUser{
			ID:        1,
			Username:  "admin",
			Email:     "admin@vacuummart.com",
			FirstName: "Admin",
			LastName:  "User",
			Role:      "admin",
			IsActive:  true,
		}
		writeJSONResponse(w, http.StatusOK, api.APIResponse{
			Success: true,
			Data:    profile,
		})
	case http.MethodPut:
		// Handle profile update
		writeJSONResponse(w, http.StatusOK, api.APIResponse{
			Success: true,
			Message: "Profile updated successfully",
		})
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleAdminProductsBulk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Message: "Bulk operation completed successfully",
	})
}

func handleUpdateProduct(w http.ResponseWriter, r *http.Request, productID int) {
	var req api.AdminProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Message: fmt.Sprintf("Product %d updated successfully", productID),
	})
}

func handleDeleteProduct(w http.ResponseWriter, r *http.Request, productID int) {
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Message: fmt.Sprintf("Product %d deleted successfully", productID),
	})
}

func handleGetAdminOrderDetail(w http.ResponseWriter, r *http.Request, orderID int) {
	order := getMockAdminOrderDetail(orderID)
	if order == nil {
		writeJSONError(w, http.StatusNotFound, "Order not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    order,
	})
}

func handleUpdateOrder(w http.ResponseWriter, r *http.Request, orderID int) {
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Message: fmt.Sprintf("Order %d updated successfully", orderID),
	})
}

func handleAdminCustomerDetail(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/customers/")
	customerID, err := strconv.Atoi(path)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	// Mock customer detail
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    map[string]interface{}{"id": customerID},
	})
}

func handleAdminInventoryDetail(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/inventory/")
	inventoryID, err := strconv.Atoi(path)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    map[string]interface{}{"id": inventoryID},
	})
}

func handleAdminInventoryMovements(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    []interface{}{},
	})
}

func handleAdminSalesAnalytics(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    getMockSalesChart(),
	})
}

func handleAdminProductAnalytics(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    getMockTopProducts(),
	})
}

func handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Mock admin users
	users := []api.AdminUser{
		{
			ID:        1,
			Username:  "admin",
			Email:     "admin@vacuummart.com",
			FirstName: "Admin",
			LastName:  "User",
			Role:      "admin",
			IsActive:  true,
		},
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    users,
	})
}

func handleAdminUserDetail(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/users/")
	userID, err := strconv.Atoi(path)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	writeJSONResponse(w, http.StatusOK, api.APIResponse{
		Success: true,
		Data:    map[string]interface{}{"id": userID},
	})
}
