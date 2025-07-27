package services

import (
	"context"
	"sync"
	"time"

	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/pkg/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// MockCartService provides an in-memory implementation for demo purposes
type MockCartService struct {
	logger  zerolog.Logger
	carts   map[int64]*models.Cart
	nextID  int64
	mu      sync.RWMutex
}

// NewMockCartService creates a new mock cart service
func NewMockCartService(logger zerolog.Logger) *MockCartService {
	return &MockCartService{
		logger: logger.With().Str("service", "mock-cart").Logger(),
		carts:  make(map[int64]*models.Cart),
		nextID: 1,
		mu:     sync.RWMutex{},
	}
}

// GetOrCreateCart gets an existing cart or creates a new one
func (s *MockCartService) GetOrCreateCart(ctx context.Context, customerID *int64, sessionID *string) (*models.Cart, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if customerID == nil && sessionID == nil {
		return nil, errors.New("either customer_id or session_id must be provided")
	}
	
	// Try to find existing cart
	for _, cart := range s.carts {
		if customerID != nil && cart.CustomerID != nil && *cart.CustomerID == *customerID {
			return cart, nil
		}
		if sessionID != nil && cart.SessionID != nil && *cart.SessionID == *sessionID {
			return cart, nil
		}
	}
	
	// Create new cart if none exists
	now := time.Now()
	cart := &models.Cart{
		ID:         s.nextID,
		CustomerID: customerID,
		SessionID:  sessionID,
		CreatedAt:  now,
		UpdatedAt:  now,
		Items:      []models.CartItem{},
	}
	
	s.carts[cart.ID] = cart
	s.nextID++
	
	return cart, nil
}

// GetCart retrieves a cart by ID
func (s *MockCartService) GetCart(ctx context.Context, cartID int64) (*models.Cart, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	cart, exists := s.carts[cartID]
	if !exists {
		return nil, errors.New("cart not found")
	}
	
	return cart, nil
}

// AddToCart adds an item to the cart or updates quantity if item already exists
func (s *MockCartService) AddToCart(ctx context.Context, cartID int64, req models.AddToCartRequest) (*models.Cart, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	cart, exists := s.carts[cartID]
	if !exists {
		return nil, errors.New("cart not found")
	}
	
	// For demo purposes, we'll use a simple price of $100 for any product
	unitPrice := 100.0
	
	// Check if item already exists in cart
	for i := range cart.Items {
		if cart.Items[i].ProductID == req.ProductID {
			// Update existing item
			cart.Items[i].Quantity += req.Quantity
			cart.Items[i].UnitPrice = unitPrice
			cart.Items[i].UpdatedAt = time.Now()
			cart.UpdatedAt = time.Now()
			return cart, nil
		}
	}
	
	// Add new item
	now := time.Now()
	newItem := models.CartItem{
		ID:        int64(len(cart.Items) + 1),
		CartID:    cartID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UnitPrice: unitPrice,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	cart.Items = append(cart.Items, newItem)
	cart.UpdatedAt = now
	
	return cart, nil
}

// UpdateCartItem updates the quantity of a cart item
func (s *MockCartService) UpdateCartItem(ctx context.Context, cartID, productID int64, req models.UpdateCartItemRequest) (*models.Cart, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	cart, exists := s.carts[cartID]
	if !exists {
		return nil, errors.New("cart not found")
	}
	
	if req.Quantity == 0 {
		// Remove item if quantity is 0
		for i, item := range cart.Items {
			if item.ProductID == productID {
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
				cart.UpdatedAt = time.Now()
				return cart, nil
			}
		}
		return nil, errors.New("cart item not found")
	}
	
	// Update item quantity
	for i := range cart.Items {
		if cart.Items[i].ProductID == productID {
			cart.Items[i].Quantity = req.Quantity
			cart.Items[i].UpdatedAt = time.Now()
			cart.UpdatedAt = time.Now()
			return cart, nil
		}
	}
	
	return nil, errors.New("cart item not found")
}

// RemoveFromCart removes an item from the cart
func (s *MockCartService) RemoveFromCart(ctx context.Context, cartID, productID int64) (*models.Cart, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	cart, exists := s.carts[cartID]
	if !exists {
		return nil, errors.New("cart not found")
	}
	
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			cart.UpdatedAt = time.Now()
			return cart, nil
		}
	}
	
	return nil, errors.New("cart item not found")
}

// ClearCart removes all items from the cart
func (s *MockCartService) ClearCart(ctx context.Context, cartID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	cart, exists := s.carts[cartID]
	if !exists {
		return errors.New("cart not found")
	}
	
	cart.Items = []models.CartItem{}
	cart.UpdatedAt = time.Now()
	
	return nil
}

// MergeGuestCart merges a guest cart with a customer's existing cart
func (s *MockCartService) MergeGuestCart(ctx context.Context, customerID int64, sessionID string) (*models.Cart, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Find guest cart
	var guestCart *models.Cart
	for _, cart := range s.carts {
		if cart.SessionID != nil && *cart.SessionID == sessionID && cart.CustomerID == nil {
			guestCart = cart
			break
		}
	}
	
	if guestCart == nil {
		// No guest cart to merge, just get or create customer cart
		s.mu.Unlock()
		return s.GetOrCreateCart(context.Background(), &customerID, nil)
	}
	
	// Find or create customer cart
	var customerCart *models.Cart
	for _, cart := range s.carts {
		if cart.CustomerID != nil && *cart.CustomerID == customerID {
			customerCart = cart
			break
		}
	}
	
	if customerCart == nil {
		// Create new customer cart
		now := time.Now()
		customerCart = &models.Cart{
			ID:         s.nextID,
			CustomerID: &customerID,
			CreatedAt:  now,
			UpdatedAt:  now,
			Items:      []models.CartItem{},
		}
		s.carts[customerCart.ID] = customerCart
		s.nextID++
	}
	
	// Merge items from guest cart to customer cart
	for _, guestItem := range guestCart.Items {
		found := false
		for i := range customerCart.Items {
			if customerCart.Items[i].ProductID == guestItem.ProductID {
				// Update existing item
				customerCart.Items[i].Quantity += guestItem.Quantity
				customerCart.Items[i].UpdatedAt = time.Now()
				found = true
				break
			}
		}
		
		if !found {
			// Add new item to customer cart
			newItem := guestItem
			newItem.ID = int64(len(customerCart.Items) + 1)
			newItem.CartID = customerCart.ID
			newItem.UpdatedAt = time.Now()
			customerCart.Items = append(customerCart.Items, newItem)
		}
	}
	
	customerCart.UpdatedAt = time.Now()
	
	// Delete guest cart
	delete(s.carts, guestCart.ID)
	
	return customerCart, nil
}

// GetCartSummary calculates and returns cart summary
func (s *MockCartService) GetCartSummary(ctx context.Context, cart *models.Cart) models.CartSummary {
	return cart.CalculateSummary()
}

// ReserveCartItems reserves inventory for cart items (mock implementation)
func (s *MockCartService) ReserveCartItems(ctx context.Context, cartID int64) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	_, exists := s.carts[cartID]
	if !exists {
		return errors.New("cart not found")
	}
	
	// In mock implementation, we always succeed
	s.logger.Info().Int64("cart_id", cartID).Msg("Mock inventory reservation successful")
	return nil
}

// ReleaseCartReservations releases inventory reservations for cart items (mock implementation)
func (s *MockCartService) ReleaseCartReservations(ctx context.Context, cartID int64) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	_, exists := s.carts[cartID]
	if !exists {
		return errors.New("cart not found")
	}
	
	// In mock implementation, we always succeed
	s.logger.Info().Int64("cart_id", cartID).Msg("Mock inventory reservation release successful")
	return nil
}

// CleanupAbandonedCarts removes old guest carts (mock implementation)
func (s *MockCartService) CleanupAbandonedCarts(ctx context.Context, olderThan time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	cutoffTime := time.Now().Add(-olderThan)
	
	for id, cart := range s.carts {
		if cart.UpdatedAt.Before(cutoffTime) && cart.CustomerID == nil {
			delete(s.carts, id)
		}
	}
	
	return nil
}
