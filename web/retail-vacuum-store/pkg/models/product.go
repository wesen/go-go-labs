package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Product represents a vacuum cleaner product
type Product struct {
	ID                  int64              `json:"id" db:"id"`
	SKU                 string             `json:"sku" db:"sku"`
	Name                string             `json:"name" db:"name"`
	Slug                string             `json:"slug" db:"slug"`
	Description         *string            `json:"description" db:"description"`
	ShortDescription    *string            `json:"short_description" db:"short_description"`
	CategoryID          *int64             `json:"category_id" db:"category_id"`
	BrandID             *int64             `json:"brand_id" db:"brand_id"`
	ProductType         string             `json:"product_type" db:"product_type"`
	BasePrice           float64            `json:"base_price" db:"base_price"`
	SalePrice           *float64           `json:"sale_price" db:"sale_price"`
	Weight              *float64           `json:"weight" db:"weight"`
	DimensionsLength    *float64           `json:"dimensions_length" db:"dimensions_length"`
	DimensionsWidth     *float64           `json:"dimensions_width" db:"dimensions_width"`
	DimensionsHeight    *float64           `json:"dimensions_height" db:"dimensions_height"`
	WarrantyMonths      int                `json:"warranty_months" db:"warranty_months"`
	IsActive            bool               `json:"is_active" db:"is_active"`
	IsFeatured          bool               `json:"is_featured" db:"is_featured"`
	MetaTitle           *string            `json:"meta_title" db:"meta_title"`
	MetaDescription     *string            `json:"meta_description" db:"meta_description"`
	CreatedAt           time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at" db:"updated_at"`
	
	// Relationships (loaded separately)
	Category            *Category          `json:"category,omitempty"`
	Brand               *Brand             `json:"brand,omitempty"`
	Images              []ProductImage     `json:"images,omitempty"`
	Attributes          []ProductAttribute `json:"attributes,omitempty"`
	Inventory           *Inventory         `json:"inventory,omitempty"`
}

// Category represents a product category with hierarchical support
type Category struct {
	ID           int64      `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Slug         string     `json:"slug" db:"slug"`
	Description  *string    `json:"description" db:"description"`
	ParentID     *int64     `json:"parent_id" db:"parent_id"`
	DisplayOrder int        `json:"display_order" db:"display_order"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	
	// Relationships
	Parent       *Category  `json:"parent,omitempty"`
	Children     []Category `json:"children,omitempty"`
}

// Brand represents a vacuum cleaner brand
type Brand struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description *string   `json:"description" db:"description"`
	LogoURL     *string   `json:"logo_url" db:"logo_url"`
	WebsiteURL  *string   `json:"website_url" db:"website_url"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ProductImage represents product images
type ProductImage struct {
	ID           int64     `json:"id" db:"id"`
	ProductID    int64     `json:"product_id" db:"product_id"`
	URL          string    `json:"url" db:"url"`
	AltText      *string   `json:"alt_text" db:"alt_text"`
	DisplayOrder int       `json:"display_order" db:"display_order"`
	IsPrimary    bool      `json:"is_primary" db:"is_primary"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ProductAttribute represents product attributes and their values
type ProductAttribute struct {
	ID            int64  `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	Slug          string `json:"slug" db:"slug"`
	AttributeType string `json:"attribute_type" db:"attribute_type"`
	Unit          *string `json:"unit" db:"unit"`
	IsFilterable  bool   `json:"is_filterable" db:"is_filterable"`
	DisplayOrder  int    `json:"display_order" db:"display_order"`
	Value         string `json:"value" db:"value"`
}

// Inventory represents product inventory information
type Inventory struct {
	ID                int64      `json:"id" db:"id"`
	ProductID         int64      `json:"product_id" db:"product_id"`
	WarehouseLocation *string    `json:"warehouse_location" db:"warehouse_location"`
	QuantityOnHand    int        `json:"quantity_on_hand" db:"quantity_on_hand"`
	QuantityReserved  int        `json:"quantity_reserved" db:"quantity_reserved"`
	QuantityAvailable int        `json:"quantity_available" db:"quantity_available"`
	ReorderLevel      int        `json:"reorder_level" db:"reorder_level"`
	ReorderQuantity   int        `json:"reorder_quantity" db:"reorder_quantity"`
	CostPerUnit       *float64   `json:"cost_per_unit" db:"cost_per_unit"`
	LastReceived      *time.Time `json:"last_received" db:"last_received"`
	LastSold          *time.Time `json:"last_sold" db:"last_sold"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// ProductFilter represents filters for product searches
type ProductFilter struct {
	CategoryID       *int64             `json:"category_id"`
	BrandID          *int64             `json:"brand_id"`
	ProductType      *string            `json:"product_type"`
	MinPrice         *float64           `json:"min_price"`
	MaxPrice         *float64           `json:"max_price"`
	IsActive         *bool              `json:"is_active"`
	IsFeatured       *bool              `json:"is_featured"`
	SearchTerm       string             `json:"search_term"`
	Attributes       map[string]string  `json:"attributes"`
	InStock          *bool              `json:"in_stock"`
	SortBy           string             `json:"sort_by"`
	SortOrder        string             `json:"sort_order"`
	Limit            int                `json:"limit"`
	Offset           int                `json:"offset"`
}

// ProductCreateRequest represents the data needed to create a new product
type ProductCreateRequest struct {
	SKU                 string                       `json:"sku" validate:"required"`
	Name                string                       `json:"name" validate:"required"`
	Slug                string                       `json:"slug" validate:"required"`
	Description         *string                      `json:"description"`
	ShortDescription    *string                      `json:"short_description"`
	CategoryID          *int64                       `json:"category_id"`
	BrandID             *int64                       `json:"brand_id"`
	ProductType         string                       `json:"product_type" validate:"required,oneof=vacuum accessory part bundle"`
	BasePrice           float64                      `json:"base_price" validate:"required,min=0"`
	SalePrice           *float64                     `json:"sale_price" validate:"omitempty,min=0"`
	Weight              *float64                     `json:"weight" validate:"omitempty,min=0"`
	DimensionsLength    *float64                     `json:"dimensions_length" validate:"omitempty,min=0"`
	DimensionsWidth     *float64                     `json:"dimensions_width" validate:"omitempty,min=0"`
	DimensionsHeight    *float64                     `json:"dimensions_height" validate:"omitempty,min=0"`
	WarrantyMonths      int                          `json:"warranty_months" validate:"min=0"`
	IsActive            bool                         `json:"is_active"`
	IsFeatured          bool                         `json:"is_featured"`
	MetaTitle           *string                      `json:"meta_title"`
	MetaDescription     *string                      `json:"meta_description"`
	Images              []ProductImageCreateRequest  `json:"images"`
	Attributes          []ProductAttributeRequest    `json:"attributes"`
	InitialInventory    *InventoryCreateRequest      `json:"initial_inventory"`
}

// ProductUpdateRequest represents the data needed to update a product
type ProductUpdateRequest struct {
	Name                *string                     `json:"name"`
	Slug                *string                     `json:"slug"`
	Description         *string                     `json:"description"`
	ShortDescription    *string                     `json:"short_description"`
	CategoryID          *int64                      `json:"category_id"`
	BrandID             *int64                      `json:"brand_id"`
	ProductType         *string                     `json:"product_type" validate:"omitempty,oneof=vacuum accessory part bundle"`
	BasePrice           *float64                    `json:"base_price" validate:"omitempty,min=0"`
	SalePrice           *float64                    `json:"sale_price" validate:"omitempty,min=0"`
	Weight              *float64                    `json:"weight" validate:"omitempty,min=0"`
	DimensionsLength    *float64                    `json:"dimensions_length" validate:"omitempty,min=0"`
	DimensionsWidth     *float64                    `json:"dimensions_width" validate:"omitempty,min=0"`
	DimensionsHeight    *float64                    `json:"dimensions_height" validate:"omitempty,min=0"`
	WarrantyMonths      *int                        `json:"warranty_months" validate:"omitempty,min=0"`
	IsActive            *bool                       `json:"is_active"`
	IsFeatured          *bool                       `json:"is_featured"`
	MetaTitle           *string                     `json:"meta_title"`
	MetaDescription     *string                     `json:"meta_description"`
}

// ProductImageCreateRequest represents the data needed to create a product image
type ProductImageCreateRequest struct {
	URL          string  `json:"url" validate:"required,url"`
	AltText      *string `json:"alt_text"`
	DisplayOrder int     `json:"display_order"`
	IsPrimary    bool    `json:"is_primary"`
}

// ProductAttributeRequest represents attribute values for a product
type ProductAttributeRequest struct {
	AttributeID int64  `json:"attribute_id" validate:"required"`
	Value       string `json:"value" validate:"required"`
}

// InventoryCreateRequest represents initial inventory data
type InventoryCreateRequest struct {
	WarehouseLocation *string  `json:"warehouse_location"`
	QuantityOnHand    int      `json:"quantity_on_hand" validate:"min=0"`
	ReorderLevel      int      `json:"reorder_level" validate:"min=0"`
	ReorderQuantity   int      `json:"reorder_quantity" validate:"min=0"`
	CostPerUnit       *float64 `json:"cost_per_unit" validate:"omitempty,min=0"`
}

// SearchResult represents search results with pagination
type ProductSearchResult struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	Pages    int       `json:"pages"`
}

// CategoryTreeNode represents a category with its children for hierarchical display
type CategoryTreeNode struct {
	Category
	Children []*CategoryTreeNode `json:"children"`
}

// Scan implements the sql.Scanner interface for custom JSON fields
func (p *ProductFilter) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, p)
	case string:
		return json.Unmarshal([]byte(v), p)
	default:
		return fmt.Errorf("cannot scan %T into ProductFilter", value)
	}
}

// Value implements the driver.Valuer interface for custom JSON fields
func (p ProductFilter) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// GetCurrentPrice returns the current effective price (sale price if available, otherwise base price)
func (p *Product) GetCurrentPrice() float64 {
	if p.SalePrice != nil && *p.SalePrice > 0 {
		return *p.SalePrice
	}
	return p.BasePrice
}

// IsOnSale returns true if the product has a sale price
func (p *Product) IsOnSale() bool {
	return p.SalePrice != nil && *p.SalePrice > 0 && *p.SalePrice < p.BasePrice
}

// GetPrimaryImage returns the primary image URL or the first image if no primary is set
func (p *Product) GetPrimaryImage() *ProductImage {
	if len(p.Images) == 0 {
		return nil
	}
	
	for _, img := range p.Images {
		if img.IsPrimary {
			return &img
		}
	}
	
	return &p.Images[0]
}

// IsInStock returns true if the product has available inventory
func (p *Product) IsInStock() bool {
	return p.Inventory != nil && p.Inventory.QuantityAvailable > 0
}

// GetDimensions returns formatted dimensions string
func (p *Product) GetDimensions() string {
	if p.DimensionsLength == nil || p.DimensionsWidth == nil || p.DimensionsHeight == nil {
		return ""
	}
	return fmt.Sprintf("%.1f × %.1f × %.1f", *p.DimensionsLength, *p.DimensionsWidth, *p.DimensionsHeight)
}
