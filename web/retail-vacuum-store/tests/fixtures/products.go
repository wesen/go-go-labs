package fixtures

import (
	"time"
)

// Product test fixtures for consistent testing across test suites

type TestProduct struct {
	ID          int     `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	BasePrice   float64 `json:"base_price"`
	SalePrice   *float64 `json:"sale_price,omitempty"`
	CategoryID  int     `json:"category_id"`
	BrandID     int     `json:"brand_id"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type TestCategory struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ParentID    *int   `json:"parent_id,omitempty"`
	IsActive    bool   `json:"is_active"`
}

type TestBrand struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type TestCartItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type TestOrder struct {
	ID          int           `json:"id"`
	CustomerID  int           `json:"customer_id"`
	TotalAmount float64       `json:"total_amount"`
	Status      string        `json:"status"`
	Items       []TestCartItem `json:"items"`
	CreatedAt   time.Time     `json:"created_at"`
}

// Test data constants
var (
	// Categories
	VacuumCategory = TestCategory{
		ID:          1,
		Name:        "Vacuum Cleaners",
		Slug:        "vacuum-cleaners",
		Description: "All types of vacuum cleaners",
		IsActive:    true,
	}

	UprightCategory = TestCategory{
		ID:          2,
		Name:        "Upright Vacuums",
		Slug:        "upright-vacuums", 
		Description: "Traditional upright vacuum cleaners",
		ParentID:    &[]int{1}[0],
		IsActive:    true,
	}

	CanisterCategory = TestCategory{
		ID:          3,
		Name:        "Canister Vacuums",
		Slug:        "canister-vacuums",
		Description: "Portable canister vacuum cleaners",
		ParentID:    &[]int{1}[0],
		IsActive:    true,
	}

	// Brands
	DysonBrand = TestBrand{
		ID:          1,
		Name:        "Dyson",
		Slug:        "dyson",
		Description: "Premium vacuum cleaner manufacturer",
		IsActive:    true,
	}

	SharkBrand = TestBrand{
		ID:          2,
		Name:        "Shark",
		Slug:        "shark",
		Description: "Innovative cleaning solutions",
		IsActive:    true,
	}

	// Products
	DysonV15 = TestProduct{
		ID:          1,
		SKU:         "DYS-V15-001",
		Name:        "Dyson V15 Detect",
		Description: "Cordless vacuum with laser dust detection",
		BasePrice:   749.99,
		CategoryID:  1,
		BrandID:     1,
		IsActive:    true,
		CreatedAt:   time.Now().AddDate(0, -1, 0),
	}

	SharkNavigator = TestProduct{
		ID:          2,
		SKU:         "SHK-NAV-001",
		Name:        "Shark Navigator Lift-Away",
		Description: "Lightweight upright vacuum with lift-away canister",
		BasePrice:   179.99,
		SalePrice:   &[]float64{149.99}[0],
		CategoryID:  2,
		BrandID:     2,
		IsActive:    true,
		CreatedAt:   time.Now().AddDate(0, -2, 0),
	}

	MieleCX1 = TestProduct{
		ID:          3,
		SKU:         "MIE-CX1-001",
		Name:        "Miele Complete C3",
		Description: "Canister vacuum with HEPA filtration",
		BasePrice:   399.99,
		CategoryID:  3,
		BrandID:     3,
		IsActive:    true,
		CreatedAt:   time.Now().AddDate(0, -3, 0),
	}

	// Test cart scenarios
	SingleItemCart = []TestCartItem{
		{ProductID: 1, Quantity: 1},
	}

	MultipleItemCart = []TestCartItem{
		{ProductID: 1, Quantity: 2},
		{ProductID: 2, Quantity: 1},
		{ProductID: 3, Quantity: 3},
	}

	LargeQuantityCart = []TestCartItem{
		{ProductID: 1, Quantity: 10},
	}
)

// Helper functions for test data generation
func CreateTestProducts() []TestProduct {
	return []TestProduct{DysonV15, SharkNavigator, MieleCX1}
}

func CreateTestCategories() []TestCategory {
	return []TestCategory{VacuumCategory, UprightCategory, CanisterCategory}
}

func CreateTestBrands() []TestBrand {
	return []TestBrand{DysonBrand, SharkBrand}
}

func CreateTestCart(scenario string) []TestCartItem {
	switch scenario {
	case "single":
		return SingleItemCart
	case "multiple":
		return MultipleItemCart
	case "large":
		return LargeQuantityCart
	default:
		return SingleItemCart
	}
}
