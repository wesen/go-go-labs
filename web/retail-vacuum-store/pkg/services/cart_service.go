package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type CartService struct {
	db            *sqlx.DB
	logger        zerolog.Logger
	productService *ProductService
}

func NewCartService(db *sqlx.DB, logger zerolog.Logger, productService *ProductService) *CartService {
	return &CartService{
		db:            db,
		logger:        logger.With().Str("service", "cart").Logger(),
		productService: productService,
	}
}

// GetOrCreateCart gets an existing cart or creates a new one
func (s *CartService) GetOrCreateCart(ctx context.Context, customerID *int64, sessionID *string) (*models.Cart, error) {
	if customerID == nil && sessionID == nil {
		return nil, errors.New("either customer_id or session_id must be provided")
	}
	
	// Try to find existing cart
	var cart *models.Cart
	var err error
	
	if customerID != nil {
		cart, err = s.getCartByCustomerID(ctx, *customerID)
	} else {
		cart, err = s.getCartBySessionID(ctx, *sessionID)
	}
	
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "failed to get existing cart")
	}
	
	// Create new cart if none exists
	if cart == nil {
		cart, err = s.createCart(ctx, customerID, sessionID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create cart")
		}
	}
	
	// Load cart items
	if err := s.loadCartItems(ctx, cart); err != nil {
		return nil, errors.Wrap(err, "failed to load cart items")
	}
	
	return cart, nil
}

// GetCart retrieves a cart by ID
func (s *CartService) GetCart(ctx context.Context, cartID int64) (*models.Cart, error) {
	cart := &models.Cart{}
	err := s.db.GetContext(ctx, cart, `
		SELECT id, customer_id, session_id, created_at, updated_at
		FROM carts
		WHERE id = $1`, cartID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("cart not found")
		}
		return nil, errors.Wrap(err, "failed to get cart")
	}
	
	if err := s.loadCartItems(ctx, cart); err != nil {
		return nil, errors.Wrap(err, "failed to load cart items")
	}
	
	return cart, nil
}

// AddToCart adds an item to the cart or updates quantity if item already exists
func (s *CartService) AddToCart(ctx context.Context, cartID int64, req models.AddToCartRequest) (*models.Cart, error) {
	// Verify product exists and is active
	product, err := s.productService.GetProduct(ctx, req.ProductID)
	if err != nil {
		return nil, errors.Wrap(err, "product not found")
	}
	
	if !product.IsActive {
		return nil, errors.New("product is not available")
	}
	
	// Check inventory
	if product.Inventory == nil || product.Inventory.QuantityAvailable < req.Quantity {
		return nil, errors.New("insufficient inventory")
	}
	
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()
	
	// Check if item already exists in cart
	var existingQuantity int
	err = tx.GetContext(ctx, &existingQuantity, `
		SELECT quantity FROM cart_items WHERE cart_id = $1 AND product_id = $2`,
		cartID, req.ProductID)
	
	currentPrice := product.GetCurrentPrice()
	
	if err == sql.ErrNoRows {
		// Add new item
		_, err = tx.ExecContext(ctx, `
			INSERT INTO cart_items (cart_id, product_id, quantity, unit_price)
			VALUES ($1, $2, $3, $4)`,
			cartID, req.ProductID, req.Quantity, currentPrice)
		if err != nil {
			return nil, errors.Wrap(err, "failed to add item to cart")
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to check existing cart item")
	} else {
		// Update existing item
		newQuantity := existingQuantity + req.Quantity
		
		// Check total quantity against inventory
		if product.Inventory.QuantityAvailable < newQuantity {
			return nil, errors.New("insufficient inventory for requested quantity")
		}
		
		_, err = tx.ExecContext(ctx, `
			UPDATE cart_items 
			SET quantity = $1, unit_price = $2, updated_at = CURRENT_TIMESTAMP
			WHERE cart_id = $3 AND product_id = $4`,
			newQuantity, currentPrice, cartID, req.ProductID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to update cart item")
		}
	}
	
	// Update cart timestamp
	_, err = tx.ExecContext(ctx, `
		UPDATE carts SET updated_at = CURRENT_TIMESTAMP WHERE id = $1`, cartID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update cart timestamp")
	}
	
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "failed to commit transaction")
	}
	
	return s.GetCart(ctx, cartID)
}

// UpdateCartItem updates the quantity of a cart item
func (s *CartService) UpdateCartItem(ctx context.Context, cartID, productID int64, req models.UpdateCartItemRequest) (*models.Cart, error) {
	if req.Quantity == 0 {
		return s.RemoveFromCart(ctx, cartID, productID)
	}
	
	// Verify product exists and check inventory
	product, err := s.productService.GetProduct(ctx, productID)
	if err != nil {
		return nil, errors.Wrap(err, "product not found")
	}
	
	if product.Inventory == nil || product.Inventory.QuantityAvailable < req.Quantity {
		return nil, errors.New("insufficient inventory")
	}
	
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()
	
	// Update cart item
	result, err := tx.ExecContext(ctx, `
		UPDATE cart_items 
		SET quantity = $1, unit_price = $2, updated_at = CURRENT_TIMESTAMP
		WHERE cart_id = $3 AND product_id = $4`,
		req.Quantity, product.GetCurrentPrice(), cartID, productID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update cart item")
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get rows affected")
	}
	
	if rowsAffected == 0 {
		return nil, errors.New("cart item not found")
	}
	
	// Update cart timestamp
	_, err = tx.ExecContext(ctx, `
		UPDATE carts SET updated_at = CURRENT_TIMESTAMP WHERE id = $1`, cartID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update cart timestamp")
	}
	
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "failed to commit transaction")
	}
	
	return s.GetCart(ctx, cartID)
}

// RemoveFromCart removes an item from the cart
func (s *CartService) RemoveFromCart(ctx context.Context, cartID, productID int64) (*models.Cart, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()
	
	// Remove cart item
	result, err := tx.ExecContext(ctx, `
		DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2`,
		cartID, productID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to remove cart item")
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get rows affected")
	}
	
	if rowsAffected == 0 {
		return nil, errors.New("cart item not found")
	}
	
	// Update cart timestamp
	_, err = tx.ExecContext(ctx, `
		UPDATE carts SET updated_at = CURRENT_TIMESTAMP WHERE id = $1`, cartID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update cart timestamp")
	}
	
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "failed to commit transaction")
	}
	
	return s.GetCart(ctx, cartID)
}

// ClearCart removes all items from the cart
func (s *CartService) ClearCart(ctx context.Context, cartID int64) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()
	
	// Remove all cart items
	_, err = tx.ExecContext(ctx, `DELETE FROM cart_items WHERE cart_id = $1`, cartID)
	if err != nil {
		return errors.Wrap(err, "failed to clear cart items")
	}
	
	// Update cart timestamp
	_, err = tx.ExecContext(ctx, `
		UPDATE carts SET updated_at = CURRENT_TIMESTAMP WHERE id = $1`, cartID)
	if err != nil {
		return errors.Wrap(err, "failed to update cart timestamp")
	}
	
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}
	
	return nil
}

// MergeGuestCart merges a guest cart with a customer's existing cart
func (s *CartService) MergeGuestCart(ctx context.Context, customerID int64, sessionID string) (*models.Cart, error) {
	// Get guest cart
	guestCart, err := s.getCartBySessionID(ctx, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No guest cart to merge, just get or create customer cart
			return s.GetOrCreateCart(ctx, &customerID, nil)
		}
		return nil, errors.Wrap(err, "failed to get guest cart")
	}
	
	// Load guest cart items
	if err := s.loadCartItems(ctx, guestCart); err != nil {
		return nil, errors.Wrap(err, "failed to load guest cart items")
	}
	
	// Get or create customer cart
	customerCart, err := s.GetOrCreateCart(ctx, &customerID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get customer cart")
	}
	
	// Merge items from guest cart to customer cart
	for _, item := range guestCart.Items {
		req := models.AddToCartRequest{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
		
		customerCart, err = s.AddToCart(ctx, customerCart.ID, req)
		if err != nil {
			s.logger.Error().Err(err).Int64("product_id", item.ProductID).Msg("failed to merge cart item")
			// Continue with other items even if one fails
		}
	}
	
	// Delete guest cart
	_, err = s.db.ExecContext(ctx, `DELETE FROM carts WHERE id = $1`, guestCart.ID)
	if err != nil {
		s.logger.Error().Err(err).Int64("cart_id", guestCart.ID).Msg("failed to delete guest cart after merge")
	}
	
	return customerCart, nil
}

// GetCartSummary calculates and returns cart summary
func (s *CartService) GetCartSummary(ctx context.Context, cart *models.Cart) models.CartSummary {
	return cart.CalculateSummary()
}

// ReserveCartItems reserves inventory for cart items (called before checkout)
func (s *CartService) ReserveCartItems(ctx context.Context, cartID int64) error {
	cart, err := s.GetCart(ctx, cartID)
	if err != nil {
		return errors.Wrap(err, "failed to get cart")
	}
	
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()
	
	for _, item := range cart.Items {
		// Check if we can reserve the requested quantity
		var availableQuantity int
		err = tx.GetContext(ctx, &availableQuantity, `
			SELECT quantity_available FROM inventory WHERE product_id = $1`,
			item.ProductID)
		if err != nil {
			return errors.Wrapf(err, "failed to check inventory for product %d", item.ProductID)
		}
		
		if availableQuantity < item.Quantity {
			return errors.Errorf("insufficient inventory for product %d", item.ProductID)
		}
		
		// Reserve the quantity
		_, err = tx.ExecContext(ctx, `
			UPDATE inventory 
			SET quantity_reserved = quantity_reserved + $1, updated_at = CURRENT_TIMESTAMP
			WHERE product_id = $2`,
			item.Quantity, item.ProductID)
		if err != nil {
			return errors.Wrapf(err, "failed to reserve inventory for product %d", item.ProductID)
		}
	}
	
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit inventory reservation")
	}
	
	return nil
}

// ReleaseCartReservations releases inventory reservations for cart items
func (s *CartService) ReleaseCartReservations(ctx context.Context, cartID int64) error {
	cart, err := s.GetCart(ctx, cartID)
	if err != nil {
		return errors.Wrap(err, "failed to get cart")
	}
	
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()
	
	for _, item := range cart.Items {
		_, err = tx.ExecContext(ctx, `
			UPDATE inventory 
			SET quantity_reserved = GREATEST(0, quantity_reserved - $1), updated_at = CURRENT_TIMESTAMP
			WHERE product_id = $2`,
			item.Quantity, item.ProductID)
		if err != nil {
			return errors.Wrapf(err, "failed to release inventory reservation for product %d", item.ProductID)
		}
	}
	
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit inventory release")
	}
	
	return nil
}

// CleanupAbandonedCarts removes old guest carts and empty carts
func (s *CartService) CleanupAbandonedCarts(ctx context.Context, olderThan time.Duration) error {
	cutoffTime := time.Now().Add(-olderThan)
	
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM carts 
		WHERE updated_at < $1 
		AND (
			(session_id IS NOT NULL AND customer_id IS NULL) 
			OR NOT EXISTS (SELECT 1 FROM cart_items WHERE cart_id = carts.id)
		)`, cutoffTime)
	
	if err != nil {
		return errors.Wrap(err, "failed to cleanup abandoned carts")
	}
	
	return nil
}

// getCartByCustomerID retrieves a cart by customer ID
func (s *CartService) getCartByCustomerID(ctx context.Context, customerID int64) (*models.Cart, error) {
	cart := &models.Cart{}
	err := s.db.GetContext(ctx, cart, `
		SELECT id, customer_id, session_id, created_at, updated_at
		FROM carts
		WHERE customer_id = $1
		ORDER BY updated_at DESC
		LIMIT 1`, customerID)
	return cart, err
}

// getCartBySessionID retrieves a cart by session ID
func (s *CartService) getCartBySessionID(ctx context.Context, sessionID string) (*models.Cart, error) {
	cart := &models.Cart{}
	err := s.db.GetContext(ctx, cart, `
		SELECT id, customer_id, session_id, created_at, updated_at
		FROM carts
		WHERE session_id = $1 AND customer_id IS NULL
		ORDER BY updated_at DESC
		LIMIT 1`, sessionID)
	return cart, err
}

// createCart creates a new cart
func (s *CartService) createCart(ctx context.Context, customerID *int64, sessionID *string) (*models.Cart, error) {
	cart := &models.Cart{}
	err := s.db.GetContext(ctx, cart, `
		INSERT INTO carts (customer_id, session_id)
		VALUES ($1, $2)
		RETURNING id, customer_id, session_id, created_at, updated_at`,
		customerID, sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cart")
	}
	
	cart.Items = []models.CartItem{}
	return cart, nil
}

// loadCartItems loads cart items with product details
func (s *CartService) loadCartItems(ctx context.Context, cart *models.Cart) error {
	items := []models.CartItem{}
	err := s.db.SelectContext(ctx, &items, `
		SELECT ci.id, ci.cart_id, ci.product_id, ci.quantity, ci.unit_price, 
		       ci.created_at, ci.updated_at
		FROM cart_items ci
		WHERE ci.cart_id = $1
		ORDER BY ci.created_at`, cart.ID)
	if err != nil {
		return errors.Wrap(err, "failed to load cart items")
	}
	
	// Load product details for each item
	for i := range items {
		product, err := s.productService.GetProduct(ctx, items[i].ProductID)
		if err != nil {
			s.logger.Error().Err(err).Int64("product_id", items[i].ProductID).Msg("failed to load product for cart item")
			continue
		}
		items[i].Product = product
	}
	
	cart.Items = items
	return nil
}
