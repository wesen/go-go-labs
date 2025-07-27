package handlers

import (
	"strings"
	"time"

	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/pkg/api"
)

func getMockProducts(filters api.ProductFilters) []api.Product {
	allProducts := []api.Product{
		{
			ID:               1,
			SKU:              "SHARK-NV356E",
			Name:             "Shark Navigator Lift-Away Professional",
			Slug:             "shark-navigator-lift-away-professional",
			Description:      "Powerful upright vacuum with detachable canister for portable cleaning. Anti-Allergen Complete Seal Technology and HEPA filter.",
			ShortDescription: "Professional upright vacuum with lift-away feature",
			CategoryID:       intPtr(1),
			BrandID:          intPtr(1),
			ProductType:      "vacuum",
			BasePrice:        179.99,
			SalePrice:        floatPtr(149.99),
			Weight:           floatPtr(13.7),
			WarrantyMonths:   60,
			IsActive:         true,
			IsFeatured:       true,
			CreatedAt:        time.Now().AddDate(0, -2, 0),
			UpdatedAt:        time.Now(),
			CategoryName:     stringPtr("Upright Vacuums"),
			BrandName:        stringPtr("Shark"),
			InventoryQty:     45,
			AvgRating:        floatPtr(4.5),
			ReviewCount:      1234,
			Images: []api.ProductImage{
				{ID: 1, ProductID: 1, URL: "/static/images/shark-nv356e-1.jpg", AltText: "Shark Navigator front view", DisplayOrder: 1, IsPrimary: true},
				{ID: 2, ProductID: 1, URL: "/static/images/shark-nv356e-2.jpg", AltText: "Shark Navigator side view", DisplayOrder: 2, IsPrimary: false},
			},
		},
		{
			ID:               2,
			SKU:              "DYSON-V15",
			Name:             "Dyson V15 Detect Absolute",
			Slug:             "dyson-v15-detect-absolute",
			Description:      "Cordless vacuum with laser dust detection and LCD screen. Most powerful suction of any cordless vacuum.",
			ShortDescription: "Cordless stick vacuum with laser detection",
			CategoryID:       intPtr(2),
			BrandID:          intPtr(2),
			ProductType:      "vacuum",
			BasePrice:        749.99,
			Weight:           floatPtr(6.8),
			WarrantyMonths:   24,
			IsActive:         true,
			IsFeatured:       true,
			CreatedAt:        time.Now().AddDate(0, -1, 0),
			UpdatedAt:        time.Now(),
			CategoryName:     stringPtr("Cordless Vacuums"),
			BrandName:        stringPtr("Dyson"),
			InventoryQty:     23,
			AvgRating:        floatPtr(4.7),
			ReviewCount:      892,
			Images: []api.ProductImage{
				{ID: 3, ProductID: 2, URL: "/static/images/dyson-v15-1.jpg", AltText: "Dyson V15 main unit", DisplayOrder: 1, IsPrimary: true},
				{ID: 4, ProductID: 2, URL: "/static/images/dyson-v15-2.jpg", AltText: "Dyson V15 with attachments", DisplayOrder: 2, IsPrimary: false},
			},
		},
		{
			ID:               3,
			SKU:              "BISSELL-2252",
			Name:             "Bissell CrossWave Pet Pro All-in-One",
			Slug:             "bissell-crosswave-pet-pro",
			Description:      "Wet dry vacuum that vacuums and washes floors at the same time. Specially designed for homes with pets.",
			ShortDescription: "Multi-surface wet dry vacuum for pet owners",
			CategoryID:       intPtr(3),
			BrandID:          intPtr(3),
			ProductType:      "vacuum",
			BasePrice:        249.99,
			SalePrice:        floatPtr(199.99),
			Weight:           floatPtr(11.5),
			WarrantyMonths:   24,
			IsActive:         true,
			IsFeatured:       false,
			CreatedAt:        time.Now().AddDate(0, -3, 0),
			UpdatedAt:        time.Now(),
			CategoryName:     stringPtr("Wet/Dry Vacuums"),
			BrandName:        stringPtr("Bissell"),
			InventoryQty:     67,
			AvgRating:        floatPtr(4.3),
			ReviewCount:      543,
			Images: []api.ProductImage{
				{ID: 5, ProductID: 3, URL: "/static/images/bissell-2252-1.jpg", AltText: "Bissell CrossWave Pet Pro", DisplayOrder: 1, IsPrimary: true},
			},
		},
		{
			ID:               4,
			SKU:              "HOOVER-UH71230",
			Name:             "Hoover WindTunnel MAX Bagged Upright",
			Slug:             "hoover-windtunnel-max-bagged",
			Description:      "Traditional bagged upright vacuum with WindTunnel technology for powerful suction and allergen filtration.",
			ShortDescription: "Bagged upright vacuum with WindTunnel technology",
			CategoryID:       intPtr(1),
			BrandID:          intPtr(4),
			ProductType:      "vacuum",
			BasePrice:        129.99,
			Weight:           floatPtr(16.2),
			WarrantyMonths:   12,
			IsActive:         true,
			IsFeatured:       false,
			CreatedAt:        time.Now().AddDate(0, -4, 0),
			UpdatedAt:        time.Now(),
			CategoryName:     stringPtr("Upright Vacuums"),
			BrandName:        stringPtr("Hoover"),
			InventoryQty:     89,
			AvgRating:        floatPtr(4.1),
			ReviewCount:      267,
			Images: []api.ProductImage{
				{ID: 6, ProductID: 4, URL: "/static/images/hoover-uh71230-1.jpg", AltText: "Hoover WindTunnel MAX", DisplayOrder: 1, IsPrimary: true},
			},
		},
		{
			ID:               5,
			SKU:              "TINECO-A11",
			Name:             "Tineco A11 Master Cordless Stick Vacuum",
			Slug:             "tineco-a11-master-cordless",
			Description:      "Lightweight cordless stick vacuum with powerful suction and multiple accessories for versatile cleaning.",
			ShortDescription: "Lightweight cordless stick vacuum",
			CategoryID:       intPtr(2),
			BrandID:          intPtr(5),
			ProductType:      "vacuum",
			BasePrice:        89.99,
			SalePrice:        floatPtr(69.99),
			Weight:           floatPtr(3.3),
			WarrantyMonths:   12,
			IsActive:         true,
			IsFeatured:       false,
			CreatedAt:        time.Now().AddDate(0, -2, -15),
			UpdatedAt:        time.Now(),
			CategoryName:     stringPtr("Cordless Vacuums"),
			BrandName:        stringPtr("Tineco"),
			InventoryQty:     156,
			AvgRating:        floatPtr(4.2),
			ReviewCount:      445,
			Images: []api.ProductImage{
				{ID: 7, ProductID: 5, URL: "/static/images/tineco-a11-1.jpg", AltText: "Tineco A11 Master", DisplayOrder: 1, IsPrimary: true},
			},
		},
		{
			ID:               6,
			SKU:              "SHARK-IQ-R101AE",
			Name:             "Shark IQ Robot Self-Empty XL",
			Slug:             "shark-iq-robot-self-empty-xl",
			Description:      "Self-emptying robot vacuum with IQ Navigation for whole-home mapping and self-cleaning brushroll.",
			ShortDescription: "Self-emptying robot vacuum with mapping",
			CategoryID:       intPtr(4),
			BrandID:          intPtr(1),
			ProductType:      "vacuum",
			BasePrice:        499.99,
			SalePrice:        floatPtr(399.99),
			Weight:           floatPtr(5.51),
			WarrantyMonths:   12,
			IsActive:         true,
			IsFeatured:       true,
			CreatedAt:        time.Now().AddDate(0, -1, -10),
			UpdatedAt:        time.Now(),
			CategoryName:     stringPtr("Robot Vacuums"),
			BrandName:        stringPtr("Shark"),
			InventoryQty:     34,
			AvgRating:        floatPtr(4.4),
			ReviewCount:      789,
			Images: []api.ProductImage{
				{ID: 8, ProductID: 6, URL: "/static/images/shark-iq-r101ae-1.jpg", AltText: "Shark IQ Robot", DisplayOrder: 1, IsPrimary: true},
			},
		},
	}

	// Apply filters
	var filteredProducts []api.Product
	for _, product := range allProducts {
		// Apply search filter
		if filters.Search != "" {
			searchTerm := strings.ToLower(filters.Search)
			if !strings.Contains(strings.ToLower(product.Name), searchTerm) &&
				!strings.Contains(strings.ToLower(product.Description), searchTerm) {
				continue
			}
		}

		// Apply category filter
		if filters.CategoryID != nil && product.CategoryID != nil && *product.CategoryID != *filters.CategoryID {
			continue
		}

		// Apply brand filter
		if filters.BrandID != nil && product.BrandID != nil && *product.BrandID != *filters.BrandID {
			continue
		}

		// Apply price filters
		currentPrice := product.BasePrice
		if product.SalePrice != nil {
			currentPrice = *product.SalePrice
		}
		
		if filters.MinPrice != nil && currentPrice < *filters.MinPrice {
			continue
		}
		
		if filters.MaxPrice != nil && currentPrice > *filters.MaxPrice {
			continue
		}

		// Apply featured filter
		if filters.IsFeatured != nil && product.IsFeatured != *filters.IsFeatured {
			continue
		}

		filteredProducts = append(filteredProducts, product)
	}

	// Apply pagination
	start := (filters.Page - 1) * filters.Limit
	end := start + filters.Limit
	
	if start >= len(filteredProducts) {
		return []api.Product{}
	}
	
	if end > len(filteredProducts) {
		end = len(filteredProducts)
	}

	return filteredProducts[start:end]
}

func getMockProductDetail(productID int) *api.Product {
	products := getMockProducts(api.ProductFilters{Limit: 100})
	for _, product := range products {
		if product.ID == productID {
			// Add more detailed information for single product view
			product.Attributes = map[string]interface{}{
				"power":         "1200W",
				"cord_length":   "25 feet",
				"dust_capacity": "2.2 liters",
				"filtration":    "HEPA",
				"noise_level":   "75 dB",
				"pet_friendly":  true,
			}
			return &product
		}
	}
	return nil
}

func getMockCategories() []api.Category {
	return []api.Category{
		{ID: 1, Name: "Upright Vacuums", Slug: "upright-vacuums", Description: "Traditional upright vacuum cleaners", DisplayOrder: 1, IsActive: true},
		{ID: 2, Name: "Cordless Vacuums", Slug: "cordless-vacuums", Description: "Battery-powered stick and handheld vacuums", DisplayOrder: 2, IsActive: true},
		{ID: 3, Name: "Wet/Dry Vacuums", Slug: "wet-dry-vacuums", Description: "Multi-surface cleaning vacuums", DisplayOrder: 3, IsActive: true},
		{ID: 4, Name: "Robot Vacuums", Slug: "robot-vacuums", Description: "Autonomous robotic vacuum cleaners", DisplayOrder: 4, IsActive: true},
		{ID: 5, Name: "Canister Vacuums", Slug: "canister-vacuums", Description: "Canister-style vacuum cleaners", DisplayOrder: 5, IsActive: true},
		{ID: 6, Name: "Accessories", Slug: "accessories", Description: "Vacuum accessories and parts", DisplayOrder: 6, IsActive: true},
	}
}

func getMockBrands() []api.Brand {
	return []api.Brand{
		{ID: 1, Name: "Shark", Slug: "shark", Description: "Innovative vacuum cleaning solutions", LogoURL: "/static/images/brands/shark-logo.png", IsActive: true},
		{ID: 2, Name: "Dyson", Slug: "dyson", Description: "Premium cordless and cyclonic vacuums", LogoURL: "/static/images/brands/dyson-logo.png", IsActive: true},
		{ID: 3, Name: "Bissell", Slug: "bissell", Description: "Home cleaning and pet care specialists", LogoURL: "/static/images/brands/bissell-logo.png", IsActive: true},
		{ID: 4, Name: "Hoover", Slug: "hoover", Description: "Trusted vacuum cleaner brand since 1908", LogoURL: "/static/images/brands/hoover-logo.png", IsActive: true},
		{ID: 5, Name: "Tineco", Slug: "tineco", Description: "Smart cleaning appliances", LogoURL: "/static/images/brands/tineco-logo.png", IsActive: true},
		{ID: 6, Name: "Black+Decker", Slug: "black-decker", Description: "Reliable home and garden tools", LogoURL: "/static/images/brands/blackdecker-logo.png", IsActive: true},
	}
}

func getMockCart() api.Cart {
	items := []api.CartItem{
		{
			ID:          1,
			CartID:      1,
			ProductID:   1,
			Quantity:    1,
			UnitPrice:   149.99,
			TotalPrice:  149.99,
			ProductName: "Shark Navigator Lift-Away Professional",
			ProductSKU:  "SHARK-NV356E",
			ImageURL:    "/static/images/shark-nv356e-1.jpg",
		},
	}

	var subtotal float64
	for _, item := range items {
		subtotal += item.TotalPrice
	}

	return api.Cart{
		ID:        1,
		SessionID: stringPtr("session-123"),
		Items:     items,
		Subtotal:  subtotal,
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now(),
	}
}

// Helper functions
func intPtr(i int) *int       { return &i }
func floatPtr(f float64) *float64 { return &f }
func stringPtr(s string) *string { return &s }
