package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ProductHandler struct {
	productService ProductService
	logger         zerolog.Logger
}

type ProductService interface {
	GetProduct(ctx context.Context, id int64) (*models.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*models.Product, error)
	SearchProducts(ctx context.Context, filter models.ProductFilter) (*models.ProductSearchResult, error)
	CreateProduct(ctx context.Context, req models.ProductCreateRequest) (*models.Product, error)
	UpdateProduct(ctx context.Context, id int64, req models.ProductUpdateRequest) (*models.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
	GetCategories(ctx context.Context) ([]*models.CategoryTreeNode, error)
	GetBrands(ctx context.Context) ([]models.Brand, error)
}

func NewProductHandler(productService ProductService, logger zerolog.Logger) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		logger:         logger.With().Str("handler", "product").Logger(),
	}
}

// GetProduct retrieves a single product by ID
func (h *ProductHandler) GetProduct(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}
	
	product, err := h.productService.GetProduct(c.Request().Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Int64("product_id", id).Msg("Failed to get product")
		if err.Error() == "product not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	return c.JSON(http.StatusOK, product)
}

// GetProductBySKU retrieves a single product by SKU
func (h *ProductHandler) GetProductBySKU(c echo.Context) error {
	sku := c.Param("sku")
	if sku == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "SKU is required"})
	}
	
	product, err := h.productService.GetProductBySKU(c.Request().Context(), sku)
	if err != nil {
		h.logger.Error().Err(err).Str("sku", sku).Msg("Failed to get product by SKU")
		if err.Error() == "product not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	return c.JSON(http.StatusOK, product)
}

// SearchProducts searches for products with filters and pagination
func (h *ProductHandler) SearchProducts(c echo.Context) error {
	filter := models.ProductFilter{
		Limit:  20, // Default limit
		Offset: 0,  // Default offset
	}
	
	// Parse query parameters
	if categoryID := c.QueryParam("category_id"); categoryID != "" {
		id, err := strconv.ParseInt(categoryID, 10, 64)
		if err == nil {
			filter.CategoryID = &id
		}
	}
	
	if brandID := c.QueryParam("brand_id"); brandID != "" {
		id, err := strconv.ParseInt(brandID, 10, 64)
		if err == nil {
			filter.BrandID = &id
		}
	}
	
	if productType := c.QueryParam("product_type"); productType != "" {
		filter.ProductType = &productType
	}
	
	if minPrice := c.QueryParam("min_price"); minPrice != "" {
		price, err := strconv.ParseFloat(minPrice, 64)
		if err == nil {
			filter.MinPrice = &price
		}
	}
	
	if maxPrice := c.QueryParam("max_price"); maxPrice != "" {
		price, err := strconv.ParseFloat(maxPrice, 64)
		if err == nil {
			filter.MaxPrice = &price
		}
	}
	
	if featured := c.QueryParam("featured"); featured != "" {
		isFeatured, err := strconv.ParseBool(featured)
		if err == nil {
			filter.IsFeatured = &isFeatured
		}
	}
	
	if inStock := c.QueryParam("in_stock"); inStock != "" {
		isInStock, err := strconv.ParseBool(inStock)
		if err == nil {
			filter.InStock = &isInStock
		}
	}
	
	filter.SearchTerm = c.QueryParam("q")
	filter.SortBy = c.QueryParam("sort_by")
	filter.SortOrder = c.QueryParam("sort_order")
	
	if limit := c.QueryParam("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			filter.Limit = l
		}
	}
	
	if offset := c.QueryParam("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}
	
	// Parse page parameter (alternative to offset)
	if page := c.QueryParam("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filter.Offset = (p - 1) * filter.Limit
		}
	}
	
	// Parse attribute filters (e.g., ?attr_power=1200&attr_type=upright)
	filter.Attributes = make(map[string]string)
	for key, values := range c.QueryParams() {
		if len(key) > 5 && key[:5] == "attr_" && len(values) > 0 {
			attrSlug := key[5:]
			filter.Attributes[attrSlug] = values[0]
		}
	}
	
	result, err := h.productService.SearchProducts(c.Request().Context(), filter)
	if err != nil {
		h.logger.Error().Err(err).Interface("filter", filter).Msg("Failed to search products")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	return c.JSON(http.StatusOK, result)
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req models.ProductCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	// Validate request (you might want to use a validation library here)
	if req.SKU == "" || req.Name == "" || req.ProductType == "" || req.BasePrice < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
	}
	
	product, err := h.productService.CreateProduct(c.Request().Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Interface("request", req).Msg("Failed to create product")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	return c.JSON(http.StatusCreated, product)
}

// UpdateProduct updates an existing product
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}
	
	var req models.ProductUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	product, err := h.productService.UpdateProduct(c.Request().Context(), id, req)
	if err != nil {
		h.logger.Error().Err(err).Int64("product_id", id).Interface("request", req).Msg("Failed to update product")
		if err.Error() == "product not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	return c.JSON(http.StatusOK, product)
}

// DeleteProduct soft deletes a product
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}
	
	err = h.productService.DeleteProduct(c.Request().Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Int64("product_id", id).Msg("Failed to delete product")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	return c.NoContent(http.StatusNoContent)
}

// GetCategories retrieves all categories with hierarchy
func (h *ProductHandler) GetCategories(c echo.Context) error {
	categories, err := h.productService.GetCategories(c.Request().Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get categories")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	return c.JSON(http.StatusOK, categories)
}

// GetBrands retrieves all brands
func (h *ProductHandler) GetBrands(c echo.Context) error {
	brands, err := h.productService.GetBrands(c.Request().Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get brands")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	return c.JSON(http.StatusOK, brands)
}

// GetFeaturedProducts retrieves featured products
func (h *ProductHandler) GetFeaturedProducts(c echo.Context) error {
	limit := 10 // Default limit for featured products
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}
	
	featured := true
	filter := models.ProductFilter{
		IsFeatured: &featured,
		Limit:      limit,
		Offset:     0,
		SortBy:     "created_at",
		SortOrder:  "desc",
	}
	
	result, err := h.productService.SearchProducts(c.Request().Context(), filter)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get featured products")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	return c.JSON(http.StatusOK, result.Products)
}

// GetProductsByCategory retrieves products in a specific category
func (h *ProductHandler) GetProductsByCategory(c echo.Context) error {
	categoryID, err := strconv.ParseInt(c.Param("categoryId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}
	
	limit := 20
	offset := 0
	
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	
	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}
	
	if page := c.QueryParam("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			offset = (p - 1) * limit
		}
	}
	
	filter := models.ProductFilter{
		CategoryID: &categoryID,
		Limit:      limit,
		Offset:     offset,
		SortBy:     c.QueryParam("sort_by"),
		SortOrder:  c.QueryParam("sort_order"),
	}
	
	result, err := h.productService.SearchProducts(c.Request().Context(), filter)
	if err != nil {
		h.logger.Error().Err(err).Int64("category_id", categoryID).Msg("Failed to get products by category")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	return c.JSON(http.StatusOK, result)
}

// GetProductsByBrand retrieves products for a specific brand
func (h *ProductHandler) GetProductsByBrand(c echo.Context) error {
	brandID, err := strconv.ParseInt(c.Param("brandId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid brand ID"})
	}
	
	limit := 20
	offset := 0
	
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	
	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}
	
	if page := c.QueryParam("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			offset = (p - 1) * limit
		}
	}
	
	filter := models.ProductFilter{
		BrandID:   &brandID,
		Limit:     limit,
		Offset:    offset,
		SortBy:    c.QueryParam("sort_by"),
		SortOrder: c.QueryParam("sort_order"),
	}
	
	result, err := h.productService.SearchProducts(c.Request().Context(), filter)
	if err != nil {
		h.logger.Error().Err(err).Int64("brand_id", brandID).Msg("Failed to get products by brand")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	return c.JSON(http.StatusOK, result)
}
