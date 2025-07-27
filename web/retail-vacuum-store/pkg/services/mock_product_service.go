package services

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/pkg/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// MockProductService provides an in-memory implementation for demo purposes
type MockProductService struct {
	logger     zerolog.Logger
	products   map[int64]*models.Product
	categories map[int64]*models.Category
	brands     map[int64]*models.Brand
	nextID     int64
	mu         sync.RWMutex
}

// NewMockProductService creates a new mock product service with sample data
func NewMockProductService(logger zerolog.Logger) *MockProductService {
	service := &MockProductService{
		logger:     logger.With().Str("service", "mock-product").Logger(),
		products:   make(map[int64]*models.Product),
		categories: make(map[int64]*models.Category),
		brands:     make(map[int64]*models.Brand),
		nextID:     1,
		mu:         sync.RWMutex{},
	}
	
	service.seedData()
	return service
}

func (s *MockProductService) seedData() {
	now := time.Now()
	
	// Seed categories
	categories := []*models.Category{
		{ID: 1, Name: "Upright Vacuums", Slug: "upright-vacuums", Description: strPtr("Traditional upright vacuum cleaners"), IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Name: "Canister Vacuums", Slug: "canister-vacuums", Description: strPtr("Canister-style vacuum cleaners"), IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: 3, Name: "Robot Vacuums", Slug: "robot-vacuums", Description: strPtr("Automated robot vacuum cleaners"), IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: 4, Name: "Handheld Vacuums", Slug: "handheld-vacuums", Description: strPtr("Portable handheld vacuum cleaners"), IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: 5, Name: "Accessories", Slug: "accessories", Description: strPtr("Vacuum accessories and attachments"), IsActive: true, CreatedAt: now, UpdatedAt: now},
	}
	
	for _, cat := range categories {
		s.categories[cat.ID] = cat
	}
	
	// Seed brands
	brands := []*models.Brand{
		{ID: 1, Name: "Dyson", Slug: "dyson", Description: strPtr("Premium vacuum cleaner manufacturer"), IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Name: "Shark", Slug: "shark", Description: strPtr("Popular vacuum cleaner brand"), IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: 3, Name: "Bissell", Slug: "bissell", Description: strPtr("Trusted carpet cleaning solutions"), IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: 4, Name: "Hoover", Slug: "hoover", Description: strPtr("Classic vacuum cleaner brand"), IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: 5, Name: "Miele", Slug: "miele", Description: strPtr("German engineering excellence"), IsActive: true, CreatedAt: now, UpdatedAt: now},
	}
	
	for _, brand := range brands {
		s.brands[brand.ID] = brand
	}
	
	// Seed products
	products := []*models.Product{
		{
			ID: 1, SKU: "DYS-V15-001", Name: "Dyson V15 Detect", Slug: "dyson-v15-detect",
			Description: strPtr("Advanced cordless vacuum with laser dust detection"),
			ShortDescription: strPtr("Cordless vacuum with laser detection"),
			CategoryID: int64Ptr(1), BrandID: int64Ptr(1),
			ProductType: "vacuum", BasePrice: 749.99, SalePrice: float64Ptr(699.99),
			Weight: float64Ptr(6.8), WarrantyMonths: 24,
			IsActive: true, IsFeatured: true,
			CreatedAt: now, UpdatedAt: now,
			Inventory: &models.Inventory{
				ID: 1, ProductID: 1, QuantityOnHand: 25, QuantityReserved: 5,
				QuantityAvailable: 20, ReorderLevel: 5, ReorderQuantity: 20,
				CreatedAt: now, UpdatedAt: now,
			},
			Images: []models.ProductImage{
				{ID: 1, ProductID: 1, URL: "/static/images/dyson-v15.jpg", AltText: strPtr("Dyson V15 Detect"), IsPrimary: true, DisplayOrder: 0, CreatedAt: now},
			},
		},
		{
			ID: 2, SKU: "SHK-NAV-001", Name: "Shark Navigator", Slug: "shark-navigator",
			Description: strPtr("Powerful upright vacuum with lift-away technology"),
			ShortDescription: strPtr("Upright vacuum with lift-away"),
			CategoryID: int64Ptr(1), BrandID: int64Ptr(2),
			ProductType: "vacuum", BasePrice: 179.99,
			Weight: float64Ptr(15.5), WarrantyMonths: 12,
			IsActive: true, IsFeatured: true,
			CreatedAt: now, UpdatedAt: now,
			Inventory: &models.Inventory{
				ID: 2, ProductID: 2, QuantityOnHand: 15, QuantityReserved: 2,
				QuantityAvailable: 13, ReorderLevel: 5, ReorderQuantity: 15,
				CreatedAt: now, UpdatedAt: now,
			},
			Images: []models.ProductImage{
				{ID: 2, ProductID: 2, URL: "/static/images/shark-navigator.jpg", AltText: strPtr("Shark Navigator"), IsPrimary: true, DisplayOrder: 0, CreatedAt: now},
			},
		},
		{
			ID: 3, SKU: "BIS-CRS-001", Name: "Bissell CrossWave Pet Pro", Slug: "bissell-crosswave-pet-pro",
			Description: strPtr("Multi-surface cleaner perfect for pet owners"),
			ShortDescription: strPtr("Multi-surface pet cleaner"),
			CategoryID: int64Ptr(1), BrandID: int64Ptr(3),
			ProductType: "vacuum", BasePrice: 249.99,
			Weight: float64Ptr(11.2), WarrantyMonths: 12,
			IsActive: true, IsFeatured: false,
			CreatedAt: now, UpdatedAt: now,
			Inventory: &models.Inventory{
				ID: 3, ProductID: 3, QuantityOnHand: 8, QuantityReserved: 1,
				QuantityAvailable: 7, ReorderLevel: 3, ReorderQuantity: 10,
				CreatedAt: now, UpdatedAt: now,
			},
			Images: []models.ProductImage{
				{ID: 3, ProductID: 3, URL: "/static/images/bissell-crosswave.jpg", AltText: strPtr("Bissell CrossWave Pet Pro"), IsPrimary: true, DisplayOrder: 0, CreatedAt: now},
			},
		},
		{
			ID: 4, SKU: "ROB-VAC-001", Name: "iRobot Roomba i7+", Slug: "irobot-roomba-i7-plus",
			Description: strPtr("Smart robot vacuum with automatic dirt disposal"),
			ShortDescription: strPtr("Smart robot vacuum"),
			CategoryID: int64Ptr(3), BrandID: int64Ptr(1),
			ProductType: "vacuum", BasePrice: 799.99, SalePrice: float64Ptr(749.99),
			Weight: float64Ptr(7.4), WarrantyMonths: 12,
			IsActive: true, IsFeatured: true,
			CreatedAt: now, UpdatedAt: now,
			Inventory: &models.Inventory{
				ID: 4, ProductID: 4, QuantityOnHand: 12, QuantityReserved: 3,
				QuantityAvailable: 9, ReorderLevel: 5, ReorderQuantity: 10,
				CreatedAt: now, UpdatedAt: now,
			},
			Images: []models.ProductImage{
				{ID: 4, ProductID: 4, URL: "/static/images/roomba-i7.jpg", AltText: strPtr("iRobot Roomba i7+"), IsPrimary: true, DisplayOrder: 0, CreatedAt: now},
			},
		},
		{
			ID: 5, SKU: "MIE-CAT-001", Name: "Miele Complete C3", Slug: "miele-complete-c3",
			Description: strPtr("Premium canister vacuum with HEPA filtration"),
			ShortDescription: strPtr("Premium canister with HEPA"),
			CategoryID: int64Ptr(2), BrandID: int64Ptr(5),
			ProductType: "vacuum", BasePrice: 399.99,
			Weight: float64Ptr(18.7), WarrantyMonths: 36,
			IsActive: true, IsFeatured: false,
			CreatedAt: now, UpdatedAt: now,
			Inventory: &models.Inventory{
				ID: 5, ProductID: 5, QuantityOnHand: 6, QuantityReserved: 0,
				QuantityAvailable: 6, ReorderLevel: 3, ReorderQuantity: 8,
				CreatedAt: now, UpdatedAt: now,
			},
			Images: []models.ProductImage{
				{ID: 5, ProductID: 5, URL: "/static/images/miele-c3.jpg", AltText: strPtr("Miele Complete C3"), IsPrimary: true, DisplayOrder: 0, CreatedAt: now},
			},
		},
	}
	
	for _, product := range products {
		// Set category and brand references
		if product.CategoryID != nil {
			if cat, exists := s.categories[*product.CategoryID]; exists {
				product.Category = cat
			}
		}
		if product.BrandID != nil {
			if brand, exists := s.brands[*product.BrandID]; exists {
				product.Brand = brand
			}
		}
		s.products[product.ID] = product
	}
	
	s.nextID = 6
}

// GetProduct retrieves a product by ID
func (s *MockProductService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	product, exists := s.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	
	if !product.IsActive {
		return nil, errors.New("product not found")
	}
	
	return product, nil
}

// GetProductBySKU retrieves a product by SKU
func (s *MockProductService) GetProductBySKU(ctx context.Context, sku string) (*models.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for _, product := range s.products {
		if product.SKU == sku && product.IsActive {
			return product, nil
		}
	}
	
	return nil, errors.New("product not found")
}

// SearchProducts searches for products with filters
func (s *MockProductService) SearchProducts(ctx context.Context, filter models.ProductFilter) (*models.ProductSearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var results []models.Product
	
	for _, product := range s.products {
		if !product.IsActive {
			continue
		}
		
		// Apply filters
		if filter.CategoryID != nil && (product.CategoryID == nil || *product.CategoryID != *filter.CategoryID) {
			continue
		}
		
		if filter.BrandID != nil && (product.BrandID == nil || *product.BrandID != *filter.BrandID) {
			continue
		}
		
		if filter.ProductType != nil && product.ProductType != *filter.ProductType {
			continue
		}
		
		if filter.IsFeatured != nil && product.IsFeatured != *filter.IsFeatured {
			continue
		}
		
		if filter.MinPrice != nil && product.GetCurrentPrice() < *filter.MinPrice {
			continue
		}
		
		if filter.MaxPrice != nil && product.GetCurrentPrice() > *filter.MaxPrice {
			continue
		}
		
		if filter.InStock != nil && *filter.InStock && !product.IsInStock() {
			continue
		}
		
		if filter.SearchTerm != "" {
			searchMatch := false
			term := filter.SearchTerm
			if containsIgnoreCase(product.Name, term) ||
			   containsIgnoreCase(product.SKU, term) ||
			   (product.Description != nil && containsIgnoreCase(*product.Description, term)) {
				searchMatch = true
			}
			if !searchMatch {
				continue
			}
		}
		
		results = append(results, *product)
	}
	
	// Apply pagination
	total := int64(len(results))
	start := filter.Offset
	end := start + filter.Limit
	
	if start > len(results) {
		results = []models.Product{}
	} else {
		if end > len(results) {
			end = len(results)
		}
		results = results[start:end]
	}
	
	pageSize := filter.Limit
	if pageSize == 0 {
		pageSize = 20
	}
	
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}
	
	return &models.ProductSearchResult{
		Products: results,
		Total:    total,
		Page:     (filter.Offset / pageSize) + 1,
		PageSize: pageSize,
		Pages:    pages,
	}, nil
}

// CreateProduct creates a new product
func (s *MockProductService) CreateProduct(ctx context.Context, req models.ProductCreateRequest) (*models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Check for duplicate SKU
	for _, product := range s.products {
		if product.SKU == req.SKU {
			return nil, errors.New("product with this SKU already exists")
		}
	}
	
	now := time.Now()
	product := &models.Product{
		ID:                 s.nextID,
		SKU:                req.SKU,
		Name:               req.Name,
		Slug:               req.Slug,
		Description:        req.Description,
		ShortDescription:   req.ShortDescription,
		CategoryID:         req.CategoryID,
		BrandID:            req.BrandID,
		ProductType:        req.ProductType,
		BasePrice:          req.BasePrice,
		SalePrice:          req.SalePrice,
		Weight:             req.Weight,
		DimensionsLength:   req.DimensionsLength,
		DimensionsWidth:    req.DimensionsWidth,
		DimensionsHeight:   req.DimensionsHeight,
		WarrantyMonths:     req.WarrantyMonths,
		IsActive:           req.IsActive,
		IsFeatured:         req.IsFeatured,
		MetaTitle:          req.MetaTitle,
		MetaDescription:    req.MetaDescription,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	
	// Set category and brand references
	if product.CategoryID != nil {
		if cat, exists := s.categories[*product.CategoryID]; exists {
			product.Category = cat
		}
	}
	if product.BrandID != nil {
		if brand, exists := s.brands[*product.BrandID]; exists {
			product.Brand = brand
		}
	}
	
	// Add initial inventory if provided
	if req.InitialInventory != nil {
		product.Inventory = &models.Inventory{
			ID:                s.nextID,
			ProductID:         product.ID,
			WarehouseLocation: req.InitialInventory.WarehouseLocation,
			QuantityOnHand:    req.InitialInventory.QuantityOnHand,
			QuantityReserved:  0,
			QuantityAvailable: req.InitialInventory.QuantityOnHand,
			ReorderLevel:      req.InitialInventory.ReorderLevel,
			ReorderQuantity:   req.InitialInventory.ReorderQuantity,
			CostPerUnit:       req.InitialInventory.CostPerUnit,
			CreatedAt:         now,
			UpdatedAt:         now,
		}
	}
	
	s.products[product.ID] = product
	s.nextID++
	
	return product, nil
}

// UpdateProduct updates an existing product
func (s *MockProductService) UpdateProduct(ctx context.Context, id int64, req models.ProductUpdateRequest) (*models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	product, exists := s.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	
	// Update fields if provided
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Slug != nil {
		product.Slug = *req.Slug
	}
	if req.Description != nil {
		product.Description = req.Description
	}
	if req.ShortDescription != nil {
		product.ShortDescription = req.ShortDescription
	}
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
		if cat, exists := s.categories[*product.CategoryID]; exists {
			product.Category = cat
		}
	}
	if req.BrandID != nil {
		product.BrandID = req.BrandID
		if brand, exists := s.brands[*product.BrandID]; exists {
			product.Brand = brand
		}
	}
	if req.ProductType != nil {
		product.ProductType = *req.ProductType
	}
	if req.BasePrice != nil {
		product.BasePrice = *req.BasePrice
	}
	if req.SalePrice != nil {
		product.SalePrice = req.SalePrice
	}
	if req.Weight != nil {
		product.Weight = req.Weight
	}
	if req.DimensionsLength != nil {
		product.DimensionsLength = req.DimensionsLength
	}
	if req.DimensionsWidth != nil {
		product.DimensionsWidth = req.DimensionsWidth
	}
	if req.DimensionsHeight != nil {
		product.DimensionsHeight = req.DimensionsHeight
	}
	if req.WarrantyMonths != nil {
		product.WarrantyMonths = *req.WarrantyMonths
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
	}
	if req.MetaTitle != nil {
		product.MetaTitle = req.MetaTitle
	}
	if req.MetaDescription != nil {
		product.MetaDescription = req.MetaDescription
	}
	
	product.UpdatedAt = time.Now()
	
	return product, nil
}

// DeleteProduct soft deletes a product
func (s *MockProductService) DeleteProduct(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	product, exists := s.products[id]
	if !exists {
		return errors.New("product not found")
	}
	
	product.IsActive = false
	product.UpdatedAt = time.Now()
	
	return nil
}

// GetCategories retrieves all categories
func (s *MockProductService) GetCategories(ctx context.Context) ([]*models.CategoryTreeNode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var nodes []*models.CategoryTreeNode
	for _, category := range s.categories {
		if category.IsActive {
			node := &models.CategoryTreeNode{
				Category: *category,
				Children: []*models.CategoryTreeNode{},
			}
			nodes = append(nodes, node)
		}
	}
	
	return nodes, nil
}

// GetBrands retrieves all brands
func (s *MockProductService) GetBrands(ctx context.Context) ([]models.Brand, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var brands []models.Brand
	for _, brand := range s.brands {
		if brand.IsActive {
			brands = append(brands, *brand)
		}
	}
	
	return brands, nil
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func containsIgnoreCase(s, substr string) bool {
	// Simple case-insensitive search
	sLower := strings.ToLower(s)
	substrLower := strings.ToLower(substr)
	return strings.Contains(sLower, substrLower)
}
