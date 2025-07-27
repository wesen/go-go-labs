package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type ProductService struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

func NewProductService(db *sqlx.DB, logger zerolog.Logger) *ProductService {
	return &ProductService{
		db:     db,
		logger: logger.With().Str("service", "product").Logger(),
	}
}

// GetProduct retrieves a product by ID with all related data
func (s *ProductService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	product := &models.Product{}
	
	query := `
		SELECT p.*, c.name as category_name, c.slug as category_slug,
		       b.name as brand_name, b.slug as brand_slug
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN brands b ON p.brand_id = b.id
		WHERE p.id = $1 AND p.is_active = true`
	
	err := s.db.GetContext(ctx, product, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, errors.Wrap(err, "failed to get product")
	}
	
	// Load related data
	if err := s.loadProductRelations(ctx, product); err != nil {
		return nil, errors.Wrap(err, "failed to load product relations")
	}
	
	return product, nil
}

// GetProductBySKU retrieves a product by SKU
func (s *ProductService) GetProductBySKU(ctx context.Context, sku string) (*models.Product, error) {
	product := &models.Product{}
	
	query := `
		SELECT p.*, c.name as category_name, c.slug as category_slug,
		       b.name as brand_name, b.slug as brand_slug
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN brands b ON p.brand_id = b.id
		WHERE p.sku = $1 AND p.is_active = true`
	
	err := s.db.GetContext(ctx, product, query, sku)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, errors.Wrap(err, "failed to get product by SKU")
	}
	
	if err := s.loadProductRelations(ctx, product); err != nil {
		return nil, errors.Wrap(err, "failed to load product relations")
	}
	
	return product, nil
}

// SearchProducts searches for products with filters and pagination
func (s *ProductService) SearchProducts(ctx context.Context, filter models.ProductFilter) (*models.ProductSearchResult, error) {
	whereClause, args := s.buildWhereClause(filter)
	orderClause := s.buildOrderClause(filter.SortBy, filter.SortOrder)
	
	// Count total results
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT p.id)
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN brands b ON p.brand_id = b.id
		LEFT JOIN inventory i ON p.id = i.product_id
		LEFT JOIN product_attribute_values pav ON p.id = pav.product_id
		LEFT JOIN product_attributes pa ON pav.attribute_id = pa.id
		%s`, whereClause)
	
	var total int64
	err := s.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to count products")
	}
	
	// Get products
	query := fmt.Sprintf(`
		SELECT DISTINCT p.id, p.sku, p.name, p.slug, p.description, p.short_description,
		       p.category_id, p.brand_id, p.product_type, p.base_price, p.sale_price,
		       p.weight, p.dimensions_length, p.dimensions_width, p.dimensions_height,
		       p.warranty_months, p.is_active, p.is_featured, p.meta_title, p.meta_description,
		       p.created_at, p.updated_at,
		       c.name as category_name, c.slug as category_slug,
		       b.name as brand_name, b.slug as brand_slug
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN brands b ON p.brand_id = b.id
		LEFT JOIN inventory i ON p.id = i.product_id
		LEFT JOIN product_attribute_values pav ON p.id = pav.product_id
		LEFT JOIN product_attributes pa ON pav.attribute_id = pa.id
		%s
		%s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderClause, len(args)+1, len(args)+2)
	
	args = append(args, filter.Limit, filter.Offset)
	
	products := []models.Product{}
	err = s.db.SelectContext(ctx, &products, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to search products")
	}
	
	// Load related data for each product
	for i := range products {
		if err := s.loadProductRelations(ctx, &products[i]); err != nil {
			s.logger.Error().Err(err).Int64("product_id", products[i].ID).Msg("failed to load product relations")
		}
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
		Products: products,
		Total:    total,
		Page:     (filter.Offset / pageSize) + 1,
		PageSize: pageSize,
		Pages:    pages,
	}, nil
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, req models.ProductCreateRequest) (*models.Product, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()
	
	// Insert product
	query := `
		INSERT INTO products (sku, name, slug, description, short_description, category_id, brand_id,
		                     product_type, base_price, sale_price, weight, dimensions_length,
		                     dimensions_width, dimensions_height, warranty_months, is_active,
		                     is_featured, meta_title, meta_description)
		VALUES (:sku, :name, :slug, :description, :short_description, :category_id, :brand_id,
		        :product_type, :base_price, :sale_price, :weight, :dimensions_length,
		        :dimensions_width, :dimensions_height, :warranty_months, :is_active,
		        :is_featured, :meta_title, :meta_description)
		RETURNING id`
	
	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare product insert statement")
	}
	defer stmt.Close()
	
	var productID int64
	err = stmt.GetContext(ctx, &productID, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert product")
	}
	
	// Insert product images
	if len(req.Images) > 0 {
		for _, img := range req.Images {
			_, err = tx.ExecContext(ctx, `
				INSERT INTO product_images (product_id, url, alt_text, display_order, is_primary)
				VALUES ($1, $2, $3, $4, $5)`,
				productID, img.URL, img.AltText, img.DisplayOrder, img.IsPrimary)
			if err != nil {
				return nil, errors.Wrap(err, "failed to insert product image")
			}
		}
	}
	
	// Insert product attributes
	if len(req.Attributes) > 0 {
		for _, attr := range req.Attributes {
			_, err = tx.ExecContext(ctx, `
				INSERT INTO product_attribute_values (product_id, attribute_id, value)
				VALUES ($1, $2, $3)`,
				productID, attr.AttributeID, attr.Value)
			if err != nil {
				return nil, errors.Wrap(err, "failed to insert product attribute")
			}
		}
	}
	
	// Insert initial inventory if provided
	if req.InitialInventory != nil {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO inventory (product_id, warehouse_location, quantity_on_hand, reorder_level,
			                      reorder_quantity, cost_per_unit)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			productID, req.InitialInventory.WarehouseLocation, req.InitialInventory.QuantityOnHand,
			req.InitialInventory.ReorderLevel, req.InitialInventory.ReorderQuantity,
			req.InitialInventory.CostPerUnit)
		if err != nil {
			return nil, errors.Wrap(err, "failed to insert initial inventory")
		}
	}
	
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "failed to commit transaction")
	}
	
	return s.GetProduct(ctx, productID)
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, id int64, req models.ProductUpdateRequest) (*models.Product, error) {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1
	
	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Slug != nil {
		setParts = append(setParts, fmt.Sprintf("slug = $%d", argIndex))
		args = append(args, *req.Slug)
		argIndex++
	}
	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.ShortDescription != nil {
		setParts = append(setParts, fmt.Sprintf("short_description = $%d", argIndex))
		args = append(args, *req.ShortDescription)
		argIndex++
	}
	if req.CategoryID != nil {
		setParts = append(setParts, fmt.Sprintf("category_id = $%d", argIndex))
		args = append(args, *req.CategoryID)
		argIndex++
	}
	if req.BrandID != nil {
		setParts = append(setParts, fmt.Sprintf("brand_id = $%d", argIndex))
		args = append(args, *req.BrandID)
		argIndex++
	}
	if req.ProductType != nil {
		setParts = append(setParts, fmt.Sprintf("product_type = $%d", argIndex))
		args = append(args, *req.ProductType)
		argIndex++
	}
	if req.BasePrice != nil {
		setParts = append(setParts, fmt.Sprintf("base_price = $%d", argIndex))
		args = append(args, *req.BasePrice)
		argIndex++
	}
	if req.SalePrice != nil {
		setParts = append(setParts, fmt.Sprintf("sale_price = $%d", argIndex))
		args = append(args, *req.SalePrice)
		argIndex++
	}
	if req.Weight != nil {
		setParts = append(setParts, fmt.Sprintf("weight = $%d", argIndex))
		args = append(args, *req.Weight)
		argIndex++
	}
	if req.DimensionsLength != nil {
		setParts = append(setParts, fmt.Sprintf("dimensions_length = $%d", argIndex))
		args = append(args, *req.DimensionsLength)
		argIndex++
	}
	if req.DimensionsWidth != nil {
		setParts = append(setParts, fmt.Sprintf("dimensions_width = $%d", argIndex))
		args = append(args, *req.DimensionsWidth)
		argIndex++
	}
	if req.DimensionsHeight != nil {
		setParts = append(setParts, fmt.Sprintf("dimensions_height = $%d", argIndex))
		args = append(args, *req.DimensionsHeight)
		argIndex++
	}
	if req.WarrantyMonths != nil {
		setParts = append(setParts, fmt.Sprintf("warranty_months = $%d", argIndex))
		args = append(args, *req.WarrantyMonths)
		argIndex++
	}
	if req.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}
	if req.IsFeatured != nil {
		setParts = append(setParts, fmt.Sprintf("is_featured = $%d", argIndex))
		args = append(args, *req.IsFeatured)
		argIndex++
	}
	if req.MetaTitle != nil {
		setParts = append(setParts, fmt.Sprintf("meta_title = $%d", argIndex))
		args = append(args, *req.MetaTitle)
		argIndex++
	}
	if req.MetaDescription != nil {
		setParts = append(setParts, fmt.Sprintf("meta_description = $%d", argIndex))
		args = append(args, *req.MetaDescription)
		argIndex++
	}
	
	if len(setParts) == 0 {
		return s.GetProduct(ctx, id)
	}
	
	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)
	
	query := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIndex)
	
	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update product")
	}
	
	return s.GetProduct(ctx, id)
}

// DeleteProduct soft deletes a product (sets is_active to false)
func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, "UPDATE products SET is_active = false WHERE id = $1", id)
	if err != nil {
		return errors.Wrap(err, "failed to delete product")
	}
	return nil
}

// GetCategories retrieves all categories with hierarchy
func (s *ProductService) GetCategories(ctx context.Context) ([]*models.CategoryTreeNode, error) {
	categories := []models.Category{}
	err := s.db.SelectContext(ctx, &categories, `
		SELECT id, name, slug, description, parent_id, display_order, is_active, created_at, updated_at
		FROM categories
		WHERE is_active = true
		ORDER BY display_order, name`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get categories")
	}
	
	return s.buildCategoryTree(categories), nil
}

// GetBrands retrieves all active brands
func (s *ProductService) GetBrands(ctx context.Context) ([]models.Brand, error) {
	brands := []models.Brand{}
	err := s.db.SelectContext(ctx, &brands, `
		SELECT id, name, slug, description, logo_url, website_url, is_active, created_at, updated_at
		FROM brands
		WHERE is_active = true
		ORDER BY name`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get brands")
	}
	return brands, nil
}

// loadProductRelations loads images, attributes, and inventory for a product
func (s *ProductService) loadProductRelations(ctx context.Context, product *models.Product) error {
	// Load images
	err := s.db.SelectContext(ctx, &product.Images, `
		SELECT id, product_id, url, alt_text, display_order, is_primary, created_at
		FROM product_images
		WHERE product_id = $1
		ORDER BY display_order, is_primary DESC`, product.ID)
	if err != nil {
		return errors.Wrap(err, "failed to load product images")
	}
	
	// Load attributes
	err = s.db.SelectContext(ctx, &product.Attributes, `
		SELECT pa.id, pa.name, pa.slug, pa.attribute_type, pa.unit, pa.is_filterable,
		       pa.display_order, pav.value
		FROM product_attribute_values pav
		JOIN product_attributes pa ON pav.attribute_id = pa.id
		WHERE pav.product_id = $1
		ORDER BY pa.display_order`, product.ID)
	if err != nil {
		return errors.Wrap(err, "failed to load product attributes")
	}
	
	// Load inventory
	inventory := &models.Inventory{}
	err = s.db.GetContext(ctx, inventory, `
		SELECT id, product_id, warehouse_location, quantity_on_hand, quantity_reserved,
		       quantity_available, reorder_level, reorder_quantity, cost_per_unit,
		       last_received, last_sold, created_at, updated_at
		FROM inventory
		WHERE product_id = $1`, product.ID)
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrap(err, "failed to load product inventory")
	} else if err == nil {
		product.Inventory = inventory
	}
	
	return nil
}

// buildWhereClause builds the WHERE clause for product search
func (s *ProductService) buildWhereClause(filter models.ProductFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1
	
	// Always filter active products
	conditions = append(conditions, "p.is_active = true")
	
	if filter.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("p.category_id = $%d", argIndex))
		args = append(args, *filter.CategoryID)
		argIndex++
	}
	
	if filter.BrandID != nil {
		conditions = append(conditions, fmt.Sprintf("p.brand_id = $%d", argIndex))
		args = append(args, *filter.BrandID)
		argIndex++
	}
	
	if filter.ProductType != nil {
		conditions = append(conditions, fmt.Sprintf("p.product_type = $%d", argIndex))
		args = append(args, *filter.ProductType)
		argIndex++
	}
	
	if filter.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("COALESCE(p.sale_price, p.base_price) >= $%d", argIndex))
		args = append(args, *filter.MinPrice)
		argIndex++
	}
	
	if filter.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("COALESCE(p.sale_price, p.base_price) <= $%d", argIndex))
		args = append(args, *filter.MaxPrice)
		argIndex++
	}
	
	if filter.IsFeatured != nil {
		conditions = append(conditions, fmt.Sprintf("p.is_featured = $%d", argIndex))
		args = append(args, *filter.IsFeatured)
		argIndex++
	}
	
	if filter.InStock != nil && *filter.InStock {
		conditions = append(conditions, "i.quantity_available > 0")
	}
	
	if filter.SearchTerm != "" {
		conditions = append(conditions, fmt.Sprintf("(p.name ILIKE $%d OR p.description ILIKE $%d OR p.sku ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, "%"+filter.SearchTerm+"%")
		argIndex++
	}
	
	// Handle attribute filters
	for attrSlug, value := range filter.Attributes {
		conditions = append(conditions, fmt.Sprintf("EXISTS (SELECT 1 FROM product_attribute_values pav2 JOIN product_attributes pa2 ON pav2.attribute_id = pa2.id WHERE pav2.product_id = p.id AND pa2.slug = $%d AND pav2.value = $%d)", argIndex, argIndex+1))
		args = append(args, attrSlug, value)
		argIndex += 2
	}
	
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	
	return whereClause, args
}

// buildOrderClause builds the ORDER BY clause for product search
func (s *ProductService) buildOrderClause(sortBy, sortOrder string) string {
	validSortFields := map[string]string{
		"name":       "p.name",
		"price":      "COALESCE(p.sale_price, p.base_price)",
		"created_at": "p.created_at",
		"updated_at": "p.updated_at",
		"featured":   "p.is_featured",
	}
	
	field, ok := validSortFields[sortBy]
	if !ok {
		field = "p.created_at"
	}
	
	order := "ASC"
	if sortOrder == "desc" {
		order = "DESC"
	}
	
	return fmt.Sprintf("ORDER BY %s %s", field, order)
}

// buildCategoryTree builds a hierarchical tree from flat category list
func (s *ProductService) buildCategoryTree(categories []models.Category) []*models.CategoryTreeNode {
	nodeMap := make(map[int64]*models.CategoryTreeNode)
	var roots []*models.CategoryTreeNode
	
	// Create nodes
	for _, cat := range categories {
		node := &models.CategoryTreeNode{
			Category: cat,
			Children: []*models.CategoryTreeNode{},
		}
		nodeMap[cat.ID] = node
	}
	
	// Build tree
	for _, node := range nodeMap {
		if node.ParentID == nil {
			roots = append(roots, node)
		} else {
			if parent, exists := nodeMap[*node.ParentID]; exists {
				parent.Children = append(parent.Children, node)
			}
		}
	}
	
	return roots
}
