package models

import (
	"time"
)

// Cart represents a shopping cart
type Cart struct {
	ID         int64      `json:"id" db:"id"`
	CustomerID *int64     `json:"customer_id" db:"customer_id"`
	SessionID  *string    `json:"session_id" db:"session_id"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	
	// Relationships
	Items      []CartItem `json:"items,omitempty"`
	Customer   *Customer  `json:"customer,omitempty"`
}

// CartItem represents an item in a shopping cart
type CartItem struct {
	ID        int64     `json:"id" db:"id"`
	CartID    int64     `json:"cart_id" db:"cart_id"`
	ProductID int64     `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	UnitPrice float64   `json:"unit_price" db:"unit_price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	
	// Relationships
	Product   *Product  `json:"product,omitempty"`
}

// Customer represents a customer account
type Customer struct {
	ID              int64      `json:"id" db:"id"`
	Email           string     `json:"email" db:"email"`
	PasswordHash    *string    `json:"-" db:"password_hash"`
	FirstName       *string    `json:"first_name" db:"first_name"`
	LastName        *string    `json:"last_name" db:"last_name"`
	Phone           *string    `json:"phone" db:"phone"`
	DateOfBirth     *time.Time `json:"date_of_birth" db:"date_of_birth"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	EmailVerified   bool       `json:"email_verified" db:"email_verified"`
	MarketingOptIn  bool       `json:"marketing_opt_in" db:"marketing_opt_in"`
	LastLogin       *time.Time `json:"last_login" db:"last_login"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	
	// Relationships
	Addresses       []CustomerAddress `json:"addresses,omitempty"`
}

// CustomerAddress represents customer shipping/billing addresses
type CustomerAddress struct {
	ID           int64     `json:"id" db:"id"`
	CustomerID   int64     `json:"customer_id" db:"customer_id"`
	Type         string    `json:"type" db:"type"`
	FirstName    *string   `json:"first_name" db:"first_name"`
	LastName     *string   `json:"last_name" db:"last_name"`
	Company      *string   `json:"company" db:"company"`
	AddressLine1 string    `json:"address_line1" db:"address_line1"`
	AddressLine2 *string   `json:"address_line2" db:"address_line2"`
	City         string    `json:"city" db:"city"`
	State        *string   `json:"state" db:"state"`
	PostalCode   string    `json:"postal_code" db:"postal_code"`
	Country      string    `json:"country" db:"country"`
	IsDefault    bool      `json:"is_default" db:"is_default"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CartSummary represents cart totals and summary information
type CartSummary struct {
	ItemCount     int     `json:"item_count"`
	Subtotal      float64 `json:"subtotal"`
	TaxAmount     float64 `json:"tax_amount"`
	ShippingCost  float64 `json:"shipping_cost"`
	DiscountAmount float64 `json:"discount_amount"`
	Total         float64 `json:"total"`
	Currency      string  `json:"currency"`
}

// AddToCartRequest represents a request to add an item to cart
type AddToCartRequest struct {
	ProductID int64 `json:"product_id" validate:"required"`
	Quantity  int   `json:"quantity" validate:"required,min=1"`
}

// UpdateCartItemRequest represents a request to update cart item quantity
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,min=0"`
}

// CartResponse represents the full cart response with items and summary
type CartResponse struct {
	Cart    Cart        `json:"cart"`
	Summary CartSummary `json:"summary"`
}

// CustomerCreateRequest represents data needed to create a customer account
type CustomerCreateRequest struct {
	Email          string     `json:"email" validate:"required,email"`
	Password       string     `json:"password" validate:"required,min=8"`
	FirstName      *string    `json:"first_name"`
	LastName       *string    `json:"last_name"`
	Phone          *string    `json:"phone"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	MarketingOptIn bool       `json:"marketing_opt_in"`
}

// CustomerUpdateRequest represents data needed to update a customer account
type CustomerUpdateRequest struct {
	FirstName      *string    `json:"first_name"`
	LastName       *string    `json:"last_name"`
	Phone          *string    `json:"phone"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	MarketingOptIn *bool      `json:"marketing_opt_in"`
}

// CustomerAddressCreateRequest represents data needed to create a customer address
type CustomerAddressCreateRequest struct {
	Type         string  `json:"type" validate:"required,oneof=billing shipping"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	Company      *string `json:"company"`
	AddressLine1 string  `json:"address_line1" validate:"required"`
	AddressLine2 *string `json:"address_line2"`
	City         string  `json:"city" validate:"required"`
	State        *string `json:"state"`
	PostalCode   string  `json:"postal_code" validate:"required"`
	Country      string  `json:"country" validate:"required"`
	IsDefault    bool    `json:"is_default"`
}

// GetTotalPrice returns the total price for this cart item
func (ci *CartItem) GetTotalPrice() float64 {
	return ci.UnitPrice * float64(ci.Quantity)
}

// GetFullName returns the customer's full name
func (c *Customer) GetFullName() string {
	var fullName string
	if c.FirstName != nil {
		fullName = *c.FirstName
	}
	if c.LastName != nil {
		if fullName != "" {
			fullName += " "
		}
		fullName += *c.LastName
	}
	return fullName
}

// GetFormattedAddress returns a formatted address string
func (ca *CustomerAddress) GetFormattedAddress() string {
	address := ca.AddressLine1
	if ca.AddressLine2 != nil && *ca.AddressLine2 != "" {
		address += ", " + *ca.AddressLine2
	}
	address += ", " + ca.City
	if ca.State != nil && *ca.State != "" {
		address += ", " + *ca.State
	}
	address += " " + ca.PostalCode
	if ca.Country != "US" {
		address += ", " + ca.Country
	}
	return address
}

// CalculateSummary calculates the cart summary from cart items
func (c *Cart) CalculateSummary() CartSummary {
	summary := CartSummary{
		Currency: "USD",
	}
	
	for _, item := range c.Items {
		summary.ItemCount += item.Quantity
		summary.Subtotal += item.GetTotalPrice()
	}
	
	// Calculate tax (simplified - would normally use address-based calculation)
	summary.TaxAmount = summary.Subtotal * 0.08 // 8% tax rate
	
	// Calculate shipping (simplified - would normally use weight/size/distance)
	if summary.Subtotal > 0 {
		if summary.Subtotal >= 50 {
			summary.ShippingCost = 0 // Free shipping over $50
		} else {
			summary.ShippingCost = 9.99 // Standard shipping
		}
	}
	
	summary.Total = summary.Subtotal + summary.TaxAmount + summary.ShippingCost - summary.DiscountAmount
	
	return summary
}

// IsEmpty returns true if the cart has no items
func (c *Cart) IsEmpty() bool {
	return len(c.Items) == 0
}

// GetItemCount returns the total number of items in the cart
func (c *Cart) GetItemCount() int {
	count := 0
	for _, item := range c.Items {
		count += item.Quantity
	}
	return count
}

// HasProduct returns true if the cart contains the specified product
func (c *Cart) HasProduct(productID int64) bool {
	for _, item := range c.Items {
		if item.ProductID == productID {
			return true
		}
	}
	return false
}

// GetItemByProduct returns the cart item for the specified product, or nil if not found
func (c *Cart) GetItemByProduct(productID int64) *CartItem {
	for i := range c.Items {
		if c.Items[i].ProductID == productID {
			return &c.Items[i]
		}
	}
	return nil
}
