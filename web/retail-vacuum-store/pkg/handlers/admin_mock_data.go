package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/pkg/api"
)

// Mock data functions for admin API

func getMockAdminProducts(filters api.AdminProductFilters) []api.AdminProduct {
	products := []api.AdminProduct{
		{
			ID:               1,
			SKU:              "DYS-V15-001",
			Name:             "Dyson V15 Detect Absolute",
			Slug:             "dyson-v15-detect-absolute",
			Description:      "Advanced cordless vacuum with laser detection technology",
			ShortDescription: "Powerful cordless vacuum with intelligent suction",
			CategoryID:       &[]int{1}[0],
			CategoryName:     "Cordless Vacuums",
			BrandID:          &[]int{1}[0],
			BrandName:        "Dyson",
			ProductType:      "vacuum",
			BasePrice:        749.99,
			SalePrice:        &[]float64{699.99}[0],
			Weight:           &[]float64{6.8}[0],
			IsActive:         true,
			IsFeatured:       true,
			Stock:            25,
			LowStockAlert:    false,
			CreatedAt:        time.Now().AddDate(0, -2, 0),
			UpdatedAt:        time.Now().AddDate(0, 0, -1),
		},
		{
			ID:               2,
			SKU:              "SHK-NAV-001",
			Name:             "Shark Navigator Lift-Away Professional",
			Slug:             "shark-navigator-lift-away-professional",
			Description:      "Lightweight upright vacuum with lift-away technology",
			ShortDescription: "Versatile upright vacuum with detachable canister",
			CategoryID:       &[]int{2}[0],
			CategoryName:     "Upright Vacuums",
			BrandID:          &[]int{2}[0],
			BrandName:        "Shark",
			ProductType:      "vacuum",
			BasePrice:        179.99,
			SalePrice:        nil,
			Weight:           &[]float64{12.5}[0],
			IsActive:         true,
			IsFeatured:       false,
			Stock:            45,
			LowStockAlert:    false,
			CreatedAt:        time.Now().AddDate(0, -1, -15),
			UpdatedAt:        time.Now().AddDate(0, 0, -3),
		},
		{
			ID:               3,
			SKU:              "IRB-I7-001",
			Name:             "iRobot Roomba i7+",
			Slug:             "irobot-roomba-i7-plus",
			Description:      "Smart robot vacuum with automatic dirt disposal",
			ShortDescription: "Self-emptying robot vacuum with smart mapping",
			CategoryID:       &[]int{3}[0],
			CategoryName:     "Robot Vacuums",
			BrandID:          &[]int{3}[0],
			BrandName:        "iRobot",
			ProductType:      "vacuum",
			BasePrice:        799.99,
			SalePrice:        &[]float64{599.99}[0],
			Weight:           &[]float64{7.4}[0],
			IsActive:         true,
			IsFeatured:       true,
			Stock:            8,
			LowStockAlert:    true,
			CreatedAt:        time.Now().AddDate(0, -3, 0),
			UpdatedAt:        time.Now().AddDate(0, 0, -2),
		},
		{
			ID:               4,
			SKU:              "BSL-PET-001",
			Name:             "Bissell Pet Hair Eraser Turbo",
			Slug:             "bissell-pet-hair-eraser-turbo",
			Description:      "Specialized vacuum for pet hair removal",
			ShortDescription: "Powerful pet hair vacuum with specialized tools",
			CategoryID:       &[]int{2}[0],
			CategoryName:     "Upright Vacuums",
			BrandID:          &[]int{4}[0],
			BrandName:        "Bissell",
			ProductType:      "vacuum",
			BasePrice:        149.99,
			SalePrice:        nil,
			Weight:           &[]float64{17.2}[0],
			IsActive:         true,
			IsFeatured:       false,
			Stock:            3,
			LowStockAlert:    true,
			CreatedAt:        time.Now().AddDate(0, -1, -8),
			UpdatedAt:        time.Now().AddDate(0, 0, -1),
		},
	}

	// Apply filters
	var filtered []api.AdminProduct
	for _, product := range products {
		if filters.Search != "" {
			// Simple search in name and SKU
			if !containsIgnoreCase(product.Name, filters.Search) && 
			   !containsIgnoreCase(product.SKU, filters.Search) {
				continue
			}
		}

		if filters.CategoryID != nil && product.CategoryID != nil && *product.CategoryID != *filters.CategoryID {
			continue
		}

		if filters.BrandID != nil && product.BrandID != nil && *product.BrandID != *filters.BrandID {
			continue
		}

		if filters.Status != nil {
			if *filters.Status == "active" && !product.IsActive {
				continue
			}
			if *filters.Status == "inactive" && product.IsActive {
				continue
			}
		}

		filtered = append(filtered, product)
	}

	return filtered
}

func getMockAdminProductDetail(productID int) *api.AdminProduct {
	products := getMockAdminProducts(api.AdminProductFilters{})
	for _, product := range products {
		if product.ID == productID {
			return &product
		}
	}
	return nil
}

func getMockAdminOrders(filters api.AdminOrderFilters) []api.AdminOrder {
	orders := []api.AdminOrder{
		{
			ID:             1,
			OrderNumber:    "ORD-2025-001",
			CustomerID:     &[]int{1}[0],
			CustomerName:   "John Smith",
			CustomerEmail:  "john.smith@email.com",
			Status:         "delivered",
			PaymentStatus:  "paid",
			Subtotal:       699.99,
			TaxAmount:      56.00,
			ShippingAmount: 15.99,
			TotalAmount:    771.98,
			ShippingMethod: "standard",
			TrackingNumber: &[]string{"1Z999AA1234567890"}[0],
			CreatedAt:      time.Now().AddDate(0, 0, -5),
			UpdatedAt:      time.Now().AddDate(0, 0, -2),
			ShippedAt:      &[]time.Time{time.Now().AddDate(0, 0, -3)}[0],
			DeliveredAt:    &[]time.Time{time.Now().AddDate(0, 0, -1)}[0],
		},
		{
			ID:             2,
			OrderNumber:    "ORD-2025-002",
			CustomerID:     &[]int{2}[0],
			CustomerName:   "Sarah Johnson",
			CustomerEmail:  "sarah.johnson@email.com",
			Status:         "processing",
			PaymentStatus:  "paid",
			Subtotal:       179.99,
			TaxAmount:      14.40,
			ShippingAmount: 9.99,
			TotalAmount:    204.38,
			ShippingMethod: "standard",
			TrackingNumber: nil,
			CreatedAt:      time.Now().AddDate(0, 0, -2),
			UpdatedAt:      time.Now().AddDate(0, 0, -1),
			ShippedAt:      nil,
			DeliveredAt:    nil,
		},
		{
			ID:             3,
			OrderNumber:    "ORD-2025-003",
			CustomerID:     &[]int{3}[0],
			CustomerName:   "Mike Davis",
			CustomerEmail:  "mike.davis@email.com",
			Status:         "shipped",
			PaymentStatus:  "paid",
			Subtotal:       599.99,
			TaxAmount:      48.00,
			ShippingAmount: 19.99,
			TotalAmount:    667.98,
			ShippingMethod: "express",
			TrackingNumber: &[]string{"1Z999AA1234567891"}[0],
			CreatedAt:      time.Now().AddDate(0, 0, -3),
			UpdatedAt:      time.Now().AddDate(0, 0, -1),
			ShippedAt:      &[]time.Time{time.Now().AddDate(0, 0, -1)}[0],
			DeliveredAt:    nil,
		},
		{
			ID:             4,
			OrderNumber:    "ORD-2025-004",
			CustomerID:     nil,
			CustomerName:   "",
			CustomerEmail:  "guest@email.com",
			Status:         "pending",
			PaymentStatus:  "pending",
			Subtotal:       149.99,
			TaxAmount:      12.00,
			ShippingAmount: 9.99,
			TotalAmount:    171.98,
			ShippingMethod: "standard",
			TrackingNumber: nil,
			CreatedAt:      time.Now().AddDate(0, 0, -1),
			UpdatedAt:      time.Now().AddDate(0, 0, -1),
			ShippedAt:      nil,
			DeliveredAt:    nil,
		},
	}

	// Apply filters
	var filtered []api.AdminOrder
	for _, order := range orders {
		if filters.Status != nil && order.Status != *filters.Status {
			continue
		}

		if filters.CustomerID != nil && (order.CustomerID == nil || *order.CustomerID != *filters.CustomerID) {
			continue
		}

		if filters.OrderNumber != nil && order.OrderNumber != *filters.OrderNumber {
			continue
		}

		filtered = append(filtered, order)
	}

	return filtered
}

func getMockAdminOrderDetail(orderID int) *api.AdminOrder {
	orders := getMockAdminOrders(api.AdminOrderFilters{})
	for _, order := range orders {
		if order.ID == orderID {
			// Add order items and status history
			order.Items = []api.AdminOrderItem{
				{
					ID:          1,
					ProductID:   1,
					ProductName: "Dyson V15 Detect Absolute",
					SKU:         "DYS-V15-001",
					Quantity:    1,
					UnitPrice:   699.99,
					TotalPrice:  699.99,
				},
			}
			order.StatusHistory = []api.OrderStatusEntry{
				{
					Status:    "pending",
					Notes:     "Order placed",
					CreatedBy: nil,
					CreatedAt: order.CreatedAt,
				},
				{
					Status:    order.Status,
					Notes:     "Order status updated",
					CreatedBy: &[]string{"admin"}[0],
					CreatedAt: order.UpdatedAt,
				},
			}
			return &order
		}
	}
	return nil
}

func getMockAdminCustomers(filters api.AdminCustomerFilters) []api.AdminCustomer {
	customers := []api.AdminCustomer{
		{
			ID:             1,
			Email:          "john.smith@email.com",
			FirstName:      "John",
			LastName:       "Smith",
			Phone:          "555-0123",
			IsActive:       true,
			EmailVerified:  true,
			MarketingOptIn: true,
			OrderCount:     3,
			TotalSpent:     1299.97,
			LastLogin:      &[]time.Time{time.Now().AddDate(0, 0, -1)}[0],
			CreatedAt:      time.Now().AddDate(0, -6, 0),
			UpdatedAt:      time.Now().AddDate(0, 0, -1),
		},
		{
			ID:             2,
			Email:          "sarah.johnson@email.com",
			FirstName:      "Sarah",
			LastName:       "Johnson",
			Phone:          "555-0124",
			IsActive:       true,
			EmailVerified:  true,
			MarketingOptIn: false,
			OrderCount:     1,
			TotalSpent:     204.38,
			LastLogin:      &[]time.Time{time.Now().AddDate(0, 0, -3)}[0],
			CreatedAt:      time.Now().AddDate(0, -3, 0),
			UpdatedAt:      time.Now().AddDate(0, 0, -2),
		},
		{
			ID:             3,
			Email:          "mike.davis@email.com",
			FirstName:      "Mike",
			LastName:       "Davis",
			Phone:          "555-0125",
			IsActive:       true,
			EmailVerified:  false,
			MarketingOptIn: true,
			OrderCount:     2,
			TotalSpent:     899.98,
			LastLogin:      &[]time.Time{time.Now().AddDate(0, 0, -7)}[0],
			CreatedAt:      time.Now().AddDate(0, -2, 0),
			UpdatedAt:      time.Now().AddDate(0, 0, -4),
		},
	}

	// Apply filters
	var filtered []api.AdminCustomer
	for _, customer := range customers {
		if filters.Search != "" {
			if !containsIgnoreCase(customer.Email, filters.Search) &&
			   !containsIgnoreCase(customer.FirstName, filters.Search) &&
			   !containsIgnoreCase(customer.LastName, filters.Search) {
				continue
			}
		}

		if filters.Status == "active" && !customer.IsActive {
			continue
		}
		if filters.Status == "inactive" && customer.IsActive {
			continue
		}

		filtered = append(filtered, customer)
	}

	return filtered
}

func getMockInventory(filters api.InventoryFilters) []api.InventoryItem {
	inventory := []api.InventoryItem{
		{
			ID:                1,
			ProductID:         1,
			ProductName:       "Dyson V15 Detect Absolute",
			SKU:               "DYS-V15-001",
			WarehouseLocation: "Main",
			QuantityOnHand:    25,
			QuantityReserved:  2,
			QuantityAvailable: 23,
			ReorderLevel:      10,
			ReorderQuantity:   25,
			IsLowStock:        false,
			LastReceived:      &[]time.Time{time.Now().AddDate(0, 0, -5)}[0],
			LastSold:          &[]time.Time{time.Now().AddDate(0, 0, -1)}[0],
			UpdatedAt:         time.Now().AddDate(0, 0, -1),
		},
		{
			ID:                2,
			ProductID:         2,
			ProductName:       "Shark Navigator Lift-Away Professional",
			SKU:               "SHK-NAV-001",
			WarehouseLocation: "Main",
			QuantityOnHand:    45,
			QuantityReserved:  3,
			QuantityAvailable: 42,
			ReorderLevel:      15,
			ReorderQuantity:   50,
			IsLowStock:        false,
			LastReceived:      &[]time.Time{time.Now().AddDate(0, 0, -10)}[0],
			LastSold:          &[]time.Time{time.Now().AddDate(0, 0, -2)}[0],
			UpdatedAt:         time.Now().AddDate(0, 0, -2),
		},
		{
			ID:                3,
			ProductID:         3,
			ProductName:       "iRobot Roomba i7+",
			SKU:               "IRB-I7-001",
			WarehouseLocation: "Main",
			QuantityOnHand:    8,
			QuantityReserved:  1,
			QuantityAvailable: 7,
			ReorderLevel:      10,
			ReorderQuantity:   20,
			IsLowStock:        true,
			LastReceived:      &[]time.Time{time.Now().AddDate(0, 0, -15)}[0],
			LastSold:          &[]time.Time{time.Now().AddDate(0, 0, -1)}[0],
			UpdatedAt:         time.Now().AddDate(0, 0, -1),
		},
		{
			ID:                4,
			ProductID:         4,
			ProductName:       "Bissell Pet Hair Eraser Turbo",
			SKU:               "BSL-PET-001",
			WarehouseLocation: "Main",
			QuantityOnHand:    3,
			QuantityReserved:  0,
			QuantityAvailable: 3,
			ReorderLevel:      5,
			ReorderQuantity:   15,
			IsLowStock:        true,
			LastReceived:      &[]time.Time{time.Now().AddDate(0, 0, -20)}[0],
			LastSold:          &[]time.Time{time.Now().AddDate(0, 0, -3)}[0],
			UpdatedAt:         time.Now().AddDate(0, 0, -3),
		},
	}

	// Apply filters
	var filtered []api.InventoryItem
	for _, item := range inventory {
		if filters.Search != "" {
			if !containsIgnoreCase(item.ProductName, filters.Search) &&
			   !containsIgnoreCase(item.SKU, filters.Search) {
				continue
			}
		}

		if filters.LowStockOnly && !item.IsLowStock {
			continue
		}

		filtered = append(filtered, item)
	}

	return filtered
}

func getMockLowStockItems() []api.LowStockAlert {
	return []api.LowStockAlert{
		{
			ProductID:    3,
			ProductName:  "iRobot Roomba i7+",
			SKU:          "IRB-I7-001",
			CurrentStock: 7,
			ReorderLevel: 10,
			Status:       "warning",
		},
		{
			ProductID:    4,
			ProductName:  "Bissell Pet Hair Eraser Turbo",
			SKU:          "BSL-PET-001",
			CurrentStock: 3,
			ReorderLevel: 5,
			Status:       "critical",
		},
		{
			ProductID:    5,
			ProductName:  "Dyson V15 Filters (Pack of 2)",
			SKU:          "DYS-FLT-001",
			CurrentStock: 2,
			ReorderLevel: 5,
			Status:       "critical",
		},
	}
}

func getMockRecentOrders() []api.AdminOrder {
	return getMockAdminOrders(api.AdminOrderFilters{Limit: 5})[:5]
}

func getMockTopProducts() []api.TopProduct {
	return []api.TopProduct{
		{
			ProductID:  1,
			Name:       "Dyson V15 Detect Absolute",
			UnitsSold:  87,
			Revenue:    60899.13,
			Percentage: 87.0,
		},
		{
			ProductID:  2,
			Name:       "Shark Navigator Lift-Away Professional",
			UnitsSold:  64,
			Revenue:    11519.36,
			Percentage: 64.0,
		},
		{
			ProductID:  3,
			Name:       "iRobot Roomba i7+",
			UnitsSold:  45,
			Revenue:    26999.55,
			Percentage: 45.0,
		},
		{
			ProductID:  4,
			Name:       "Bissell Pet Hair Eraser Turbo",
			UnitsSold:  32,
			Revenue:    4799.68,
			Percentage: 32.0,
		},
	}
}

func getMockSalesChart() api.SalesChartData {
	return api.SalesChartData{
		Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
		Data:   []float64{12500, 15800, 18200, 14600, 22100, 28900, 19500},
	}
}

func getMockLowStockAlerts() []api.LowStockAlert {
	return getMockLowStockItems()
}

// Helper function for case-insensitive string contains
func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}

// Helper functions for parsing filters
func parseInventoryFilters(r *http.Request) api.InventoryFilters {
	query := r.URL.Query()
	
	filters := api.InventoryFilters{
		Search:       query.Get("search"),
		LowStockOnly: query.Get("low_stock_only") == "true",
		Page:         1,
		Limit:        20,
		SortBy:       "updated_at",
		SortDir:      "desc",
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

	return filters
}

func parseCustomerFilters(r *http.Request) api.AdminCustomerFilters {
	query := r.URL.Query()
	
	filters := api.AdminCustomerFilters{
		Search:  query.Get("search"),
		Status:  query.Get("status"),
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

	return filters
}
