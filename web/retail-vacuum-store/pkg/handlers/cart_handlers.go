package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type CartHandler struct {
	cartService CartService
	logger      zerolog.Logger
}

type CartService interface {
	GetOrCreateCart(ctx context.Context, customerID *int64, sessionID *string) (*models.Cart, error)
	GetCart(ctx context.Context, cartID int64) (*models.Cart, error)
	AddToCart(ctx context.Context, cartID int64, req models.AddToCartRequest) (*models.Cart, error)
	UpdateCartItem(ctx context.Context, cartID, productID int64, req models.UpdateCartItemRequest) (*models.Cart, error)
	RemoveFromCart(ctx context.Context, cartID, productID int64) (*models.Cart, error)
	ClearCart(ctx context.Context, cartID int64) error
	MergeGuestCart(ctx context.Context, customerID int64, sessionID string) (*models.Cart, error)
	GetCartSummary(ctx context.Context, cart *models.Cart) models.CartSummary
	ReserveCartItems(ctx context.Context, cartID int64) error
	ReleaseCartReservations(ctx context.Context, cartID int64) error
}

func NewCartHandler(cartService CartService, logger zerolog.Logger) *CartHandler {
	return &CartHandler{
		cartService: cartService,
		logger:      logger.With().Str("handler", "cart").Logger(),
	}
}

// GetCart retrieves the current cart with items and summary
func (h *CartHandler) GetCart(c echo.Context) error {
	cart, err := h.getOrCreateCartFromContext(c)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get cart")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	summary := h.cartService.GetCartSummary(c.Request().Context(), cart)
	
	response := models.CartResponse{
		Cart:    *cart,
		Summary: summary,
	}
	
	return c.JSON(http.StatusOK, response)
}

// AddToCart adds an item to the cart
func (h *CartHandler) AddToCart(c echo.Context) error {
	var req models.AddToCartRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	if req.ProductID <= 0 || req.Quantity <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID or quantity"})
	}
	
	cart, err := h.getOrCreateCartFromContext(c)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get cart")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	updatedCart, err := h.cartService.AddToCart(c.Request().Context(), cart.ID, req)
	if err != nil {
		h.logger.Error().Err(err).Interface("request", req).Msg("Failed to add item to cart")
		
		switch err.Error() {
		case "product not found":
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
		case "product is not available":
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Product is not available"})
		case "insufficient inventory", "insufficient inventory for requested quantity":
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Insufficient inventory"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
	}
	
	summary := h.cartService.GetCartSummary(c.Request().Context(), updatedCart)
	
	response := models.CartResponse{
		Cart:    *updatedCart,
		Summary: summary,
	}
	
	return c.JSON(http.StatusOK, response)
}

// UpdateCartItem updates the quantity of a cart item
func (h *CartHandler) UpdateCartItem(c echo.Context) error {
	productID, err := strconv.ParseInt(c.Param("productId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}
	
	var req models.UpdateCartItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	if req.Quantity < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid quantity"})
	}
	
	cart, err := h.getOrCreateCartFromContext(c)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get cart")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	updatedCart, err := h.cartService.UpdateCartItem(c.Request().Context(), cart.ID, productID, req)
	if err != nil {
		h.logger.Error().Err(err).Int64("product_id", productID).Interface("request", req).Msg("Failed to update cart item")
		
		switch err.Error() {
		case "product not found":
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
		case "cart item not found":
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Cart item not found"})
		case "insufficient inventory":
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Insufficient inventory"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
	}
	
	summary := h.cartService.GetCartSummary(c.Request().Context(), updatedCart)
	
	response := models.CartResponse{
		Cart:    *updatedCart,
		Summary: summary,
	}
	
	return c.JSON(http.StatusOK, response)
}

// RemoveFromCart removes an item from the cart
func (h *CartHandler) RemoveFromCart(c echo.Context) error {
	productID, err := strconv.ParseInt(c.Param("productId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}
	
	cart, err := h.getOrCreateCartFromContext(c)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get cart")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	updatedCart, err := h.cartService.RemoveFromCart(c.Request().Context(), cart.ID, productID)
	if err != nil {
		h.logger.Error().Err(err).Int64("product_id", productID).Msg("Failed to remove item from cart")
		
		if err.Error() == "cart item not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Cart item not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	summary := h.cartService.GetCartSummary(c.Request().Context(), updatedCart)
	
	response := models.CartResponse{
		Cart:    *updatedCart,
		Summary: summary,
	}
	
	return c.JSON(http.StatusOK, response)
}

// ClearCart removes all items from the cart
func (h *CartHandler) ClearCart(c echo.Context) error {
	cart, err := h.getOrCreateCartFromContext(c)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get cart")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	err = h.cartService.ClearCart(c.Request().Context(), cart.ID)
	if err != nil {
		h.logger.Error().Err(err).Int64("cart_id", cart.ID).Msg("Failed to clear cart")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	// Return empty cart
	emptyCart, err := h.cartService.GetCart(c.Request().Context(), cart.ID)
	if err != nil {
		h.logger.Error().Err(err).Int64("cart_id", cart.ID).Msg("Failed to get cleared cart")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	summary := h.cartService.GetCartSummary(c.Request().Context(), emptyCart)
	
	response := models.CartResponse{
		Cart:    *emptyCart,
		Summary: summary,
	}
	
	return c.JSON(http.StatusOK, response)
}

// MergeGuestCart merges a guest cart with customer cart (typically called after login)
func (h *CartHandler) MergeGuestCart(c echo.Context) error {
	// This would typically be called after authentication
	customerID := h.getCustomerIDFromContext(c)
	if customerID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
	}
	
	sessionID := h.getSessionIDFromContext(c)
	if sessionID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Session ID required"})
	}
	
	cart, err := h.cartService.MergeGuestCart(c.Request().Context(), *customerID, sessionID)
	if err != nil {
		h.logger.Error().Err(err).Int64("customer_id", *customerID).Str("session_id", sessionID).Msg("Failed to merge guest cart")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	summary := h.cartService.GetCartSummary(c.Request().Context(), cart)
	
	response := models.CartResponse{
		Cart:    *cart,
		Summary: summary,
	}
	
	return c.JSON(http.StatusOK, response)
}

// GetCartSummary returns just the cart summary (for quick checks)
func (h *CartHandler) GetCartSummary(c echo.Context) error {
	cart, err := h.getOrCreateCartFromContext(c)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get cart")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	summary := h.cartService.GetCartSummary(c.Request().Context(), cart)
	
	return c.JSON(http.StatusOK, summary)
}

// ReserveCartItems reserves inventory for cart items (called before checkout)
func (h *CartHandler) ReserveCartItems(c echo.Context) error {
	cart, err := h.getOrCreateCartFromContext(c)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get cart")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	err = h.cartService.ReserveCartItems(c.Request().Context(), cart.ID)
	if err != nil {
		h.logger.Error().Err(err).Int64("cart_id", cart.ID).Msg("Failed to reserve cart items")
		
		if err.Error() == "insufficient inventory for product" || 
		   err.Error()[:36] == "insufficient inventory for product " {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Insufficient inventory for one or more items"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "Inventory reserved successfully"})
}

// ReleaseCartReservations releases inventory reservations for cart items
func (h *CartHandler) ReleaseCartReservations(c echo.Context) error {
	cart, err := h.getOrCreateCartFromContext(c)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get cart")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	err = h.cartService.ReleaseCartReservations(c.Request().Context(), cart.ID)
	if err != nil {
		h.logger.Error().Err(err).Int64("cart_id", cart.ID).Msg("Failed to release cart reservations")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "Inventory reservations released successfully"})
}

// Helper functions

// getOrCreateCartFromContext retrieves or creates a cart based on the current session/customer context
func (h *CartHandler) getOrCreateCartFromContext(c echo.Context) (*models.Cart, error) {
	customerID := h.getCustomerIDFromContext(c)
	sessionID := h.getSessionIDFromContext(c)
	
	if customerID == nil && sessionID == "" {
		// Generate a new session ID if none exists
		sessionID = h.generateSessionID()
		h.setSessionIDInContext(c, sessionID)
	}
	
	var customerIDPtr *int64
	var sessionIDPtr *string
	
	if customerID != nil {
		customerIDPtr = customerID
	}
	if sessionID != "" {
		sessionIDPtr = &sessionID
	}
	
	return h.cartService.GetOrCreateCart(c.Request().Context(), customerIDPtr, sessionIDPtr)
}

// getCustomerIDFromContext extracts customer ID from authentication context
func (h *CartHandler) getCustomerIDFromContext(c echo.Context) *int64 {
	// This would typically extract from JWT token or session
	// For now, we'll check a header (in real implementation, use proper auth)
	if customerIDStr := c.Request().Header.Get("X-Customer-ID"); customerIDStr != "" {
		if customerID, err := strconv.ParseInt(customerIDStr, 10, 64); err == nil {
			return &customerID
		}
	}
	return nil
}

// getSessionIDFromContext extracts session ID from cookies or headers
func (h *CartHandler) getSessionIDFromContext(c echo.Context) string {
	// Try to get session ID from cookie first
	if cookie, err := c.Cookie("session_id"); err == nil {
		return cookie.Value
	}
	
	// Fallback to header
	return c.Request().Header.Get("X-Session-ID")
}

// setSessionIDInContext sets session ID in response cookie
func (h *CartHandler) setSessionIDInContext(c echo.Context, sessionID string) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 30, // 30 days
	}
	c.SetCookie(cookie)
}

// generateSessionID generates a new session ID
func (h *CartHandler) generateSessionID() string {
	// In a real implementation, use a proper session ID generator
	// For now, use a simple timestamp-based approach
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}
