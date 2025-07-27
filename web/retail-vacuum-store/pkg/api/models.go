package api

import (
	"time"
)

// Product represents a vacuum cleaner product
type Product struct {
	ID               int                    `json:"id" db:"id"`
	SKU              string                 `json:"sku" db:"sku"`
	Name             string                 `json:"name" db:"name"`
	Slug             string                 `json:"slug" db:"slug"`
	Description      string                 `json:"description" db:"description"`
	ShortDescription string                 `json:"short_description" db:"short_description"`
	CategoryID       *int                   `json:"category_id" db:"category_id"`
	BrandID          *int                   `json:"brand_id" db:"brand_id"`
	ProductType      string                 `json:"product_type" db:"product_type"`
	BasePrice        float64                `json:"base_price" db:"base_price"`
	SalePrice        *float64               `json:"sale_price" db:"sale_price"`
	Weight           *float64               `json:"weight" db:"weight"`
	WarrantyMonths   int                    `json:"warranty_months" db:"warranty_months"`
	IsActive         bool                   `json:"is_active" db:"is_active"`
	IsFeatured       bool                   `json:"is_featured" db:"is_featured"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" db:"updated_at"`
	
	// Joined fields
	CategoryName     *string                `json:"category_name,omitempty" db:"category_name"`
	BrandName        *string                `json:"brand_name,omitempty" db:"brand_name"`
	Images           []ProductImage         `json:"images,omitempty"`
	Attributes       map[string]interface{} `json:"attributes,omitempty"`
	InventoryQty     int                    `json:"inventory_qty,omitempty" db:"inventory_qty"`
	AvgRating        *float64               `json:"avg_rating,omitempty" db:"avg_rating"`
	ReviewCount      int                    `json:"review_count,omitempty" db:"review_count"`
}

// ProductImage represents product images
type ProductImage struct {
	ID           int    `json:"id" db:"id"`
	ProductID    int    `json:"product_id" db:"product_id"`
	URL          string `json:"url" db:"url"`
	AltText      string `json:"alt_text" db:"alt_text"`
	DisplayOrder int    `json:"display_order" db:"display_order"`
	IsPrimary    bool   `json:"is_primary" db:"is_primary"`
}

// Category represents product categories
type Category struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Slug         string `json:"slug" db:"slug"`
	Description  string `json:"description" db:"description"`
	ParentID     *int   `json:"parent_id" db:"parent_id"`
	DisplayOrder int    `json:"display_order" db:"display_order"`
	IsActive     bool   `json:"is_active" db:"is_active"`
}

// Brand represents vacuum cleaner brands
type Brand struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Slug        string `json:"slug" db:"slug"`
	Description string `json:"description" db:"description"`
	LogoURL     string `json:"logo_url" db:"logo_url"`
	IsActive    bool   `json:"is_active" db:"is_active"`
}

// Cart represents shopping cart
type Cart struct {
	ID         int        `json:"id" db:"id"`
	CustomerID *int       `json:"customer_id" db:"customer_id"`
	SessionID  *string    `json:"session_id" db:"session_id"`
	Items      []CartItem `json:"items,omitempty"`
	Subtotal   float64    `json:"subtotal"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

// CartItem represents items in shopping cart
type CartItem struct {
	ID          int     `json:"id" db:"id"`
	CartID      int     `json:"cart_id" db:"cart_id"`
	ProductID   int     `json:"product_id" db:"product_id"`
	Quantity    int     `json:"quantity" db:"quantity"`
	UnitPrice   float64 `json:"unit_price" db:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
	ProductName string  `json:"product_name,omitempty" db:"product_name"`
	ProductSKU  string  `json:"product_sku,omitempty" db:"product_sku"`
	ImageURL    string  `json:"image_url,omitempty" db:"image_url"`
}

// Customer represents registered customers
type Customer struct {
	ID              int       `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	FirstName       string    `json:"first_name" db:"first_name"`
	LastName        string    `json:"last_name" db:"last_name"`
	Phone           string    `json:"phone" db:"phone"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	EmailVerified   bool      `json:"email_verified" db:"email_verified"`
	MarketingOptIn  bool      `json:"marketing_opt_in" db:"marketing_opt_in"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// Order represents customer orders
type Order struct {
	ID              string      `json:"id" db:"id"`
	OrderNumber     string      `json:"order_number" db:"order_number"`
	CustomerID      *int        `json:"customer_id" db:"customer_id"`
	GuestEmail      *string     `json:"guest_email" db:"guest_email"`
	Status          string      `json:"status" db:"status"`
	Subtotal        float64     `json:"subtotal" db:"subtotal"`
	TaxAmount       float64     `json:"tax_amount" db:"tax_amount"`
	ShippingAmount  float64     `json:"shipping_amount" db:"shipping_amount"`
	DiscountAmount  float64     `json:"discount_amount" db:"discount_amount"`
	TotalAmount     float64     `json:"total_amount" db:"total_amount"`
	PaymentStatus   string      `json:"payment_status" db:"payment_status"`
	ShippingMethod  string      `json:"shipping_method" db:"shipping_method"`
	TrackingNumber  *string     `json:"tracking_number" db:"tracking_number"`
	Items           []OrderItem `json:"items,omitempty"`
	CreatedAt       time.Time   `json:"created_at" db:"created_at"`
	ShippedAt       *time.Time  `json:"shipped_at" db:"shipped_at"`
	DeliveredAt     *time.Time  `json:"delivered_at" db:"delivered_at"`
}

// OrderItem represents items in an order
type OrderItem struct {
	ID         int     `json:"id" db:"id"`
	OrderID    string  `json:"order_id" db:"order_id"`
	ProductID  int     `json:"product_id" db:"product_id"`
	SKU        string  `json:"sku" db:"sku"`
	Name       string  `json:"name" db:"name"`
	Quantity   int     `json:"quantity" db:"quantity"`
	UnitPrice  float64 `json:"unit_price" db:"unit_price"`
	TotalPrice float64 `json:"total_price" db:"total_price"`
}

// ProductFilters represents search and filter parameters
type ProductFilters struct {
	Search      string   `json:"search"`
	CategoryID  *int     `json:"category_id"`
	BrandID     *int     `json:"brand_id"`
	MinPrice    *float64 `json:"min_price"`
	MaxPrice    *float64 `json:"max_price"`
	ProductType string   `json:"product_type"`
	IsFeatured  *bool    `json:"is_featured"`
	SortBy      string   `json:"sort_by"` // price_asc, price_desc, name_asc, name_desc, newest, rating
	Page        int      `json:"page"`
	Limit       int      `json:"limit"`
}

// ProductSearchResponse represents paginated product results
type ProductSearchResponse struct {
	Products    []Product `json:"products"`
	Total       int       `json:"total"`
	Page        int       `json:"page"`
	Limit       int       `json:"limit"`
	TotalPages  int       `json:"total_pages"`
	Categories  []Category `json:"categories,omitempty"`
	Brands      []Brand   `json:"brands,omitempty"`
}

// APIResponse represents standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ==============================================
// Admin Models
// ==============================================

// Admin User Models
type AdminUser struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Admin Product Models
type AdminProduct struct {
	ID               int       `json:"id"`
	SKU              string    `json:"sku"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	Description      string    `json:"description"`
	ShortDescription string    `json:"short_description"`
	CategoryID       *int      `json:"category_id"`
	CategoryName     string    `json:"category_name,omitempty"`
	BrandID          *int      `json:"brand_id"`
	BrandName        string    `json:"brand_name,omitempty"`
	ProductType      string    `json:"product_type"`
	BasePrice        float64   `json:"base_price"`
	SalePrice        *float64  `json:"sale_price"`
	Weight           *float64  `json:"weight"`
	IsActive         bool      `json:"is_active"`
	IsFeatured       bool      `json:"is_featured"`
	Stock            int       `json:"stock,omitempty"`
	LowStockAlert    bool      `json:"low_stock_alert,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type AdminProductRequest struct {
	SKU              string   `json:"sku"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	ShortDescription string   `json:"short_description"`
	CategoryID       *int     `json:"category_id"`
	BrandID          *int     `json:"brand_id"`
	ProductType      string   `json:"product_type"`
	BasePrice        float64  `json:"base_price"`
	SalePrice        *float64 `json:"sale_price"`
	Weight           *float64 `json:"weight"`
	IsActive         bool     `json:"is_active"`
	IsFeatured       bool     `json:"is_featured"`
}

type AdminProductFilters struct {
	Search     string  `json:"search"`
	CategoryID *int    `json:"category_id"`
	BrandID    *int    `json:"brand_id"`
	Status     *string `json:"status"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	SortBy     string  `json:"sort_by"`
	SortDir    string  `json:"sort_dir"`
}

type AdminProductListResponse struct {
	Products   []AdminProduct `json:"products"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// Admin Order Models
type AdminOrder struct {
	ID                int               `json:"id"`
	OrderNumber       string            `json:"order_number"`
	CustomerID        *int              `json:"customer_id"`
	CustomerName      string            `json:"customer_name,omitempty"`
	CustomerEmail     string            `json:"customer_email"`
	Status            string            `json:"status"`
	PaymentStatus     string            `json:"payment_status"`
	Subtotal          float64           `json:"subtotal"`
	TaxAmount         float64           `json:"tax_amount"`
	ShippingAmount    float64           `json:"shipping_amount"`
	TotalAmount       float64           `json:"total_amount"`
	ShippingMethod    string            `json:"shipping_method"`
	TrackingNumber    *string           `json:"tracking_number"`
	Items             []AdminOrderItem  `json:"items,omitempty"`
	StatusHistory     []OrderStatusEntry `json:"status_history,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	ShippedAt         *time.Time        `json:"shipped_at"`
	DeliveredAt       *time.Time        `json:"delivered_at"`
}

type AdminOrderItem struct {
	ID          int     `json:"id"`
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	SKU         string  `json:"sku"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

type OrderStatusEntry struct {
	Status    string     `json:"status"`
	Notes     string     `json:"notes"`
	CreatedBy *string    `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
}

type AdminOrderFilters struct {
	Status      *string `json:"status"`
	CustomerID  *int    `json:"customer_id"`
	OrderNumber *string `json:"order_number"`
	Page        int     `json:"page"`
	Limit       int     `json:"limit"`
	SortBy      string  `json:"sort_by"`
	SortDir     string  `json:"sort_dir"`
}

type AdminOrderListResponse struct {
	Orders     []AdminOrder `json:"orders"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
	TotalPages int          `json:"total_pages"`
}

// Admin Customer Models
type AdminCustomer struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Phone          string    `json:"phone"`
	IsActive       bool      `json:"is_active"`
	EmailVerified  bool      `json:"email_verified"`
	MarketingOptIn bool      `json:"marketing_opt_in"`
	OrderCount     int       `json:"order_count"`
	TotalSpent     float64   `json:"total_spent"`
	LastLogin      *time.Time `json:"last_login"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AdminCustomerFilters struct {
	Search  string `json:"search"`
	Status  string `json:"status"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	SortBy  string `json:"sort_by"`
	SortDir string `json:"sort_dir"`
}

type AdminCustomerListResponse struct {
	Customers  []AdminCustomer `json:"customers"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// Inventory Models
type InventoryItem struct {
	ID                int       `json:"id"`
	ProductID         int       `json:"product_id"`
	ProductName       string    `json:"product_name"`
	SKU               string    `json:"sku"`
	WarehouseLocation string    `json:"warehouse_location"`
	QuantityOnHand    int       `json:"quantity_on_hand"`
	QuantityReserved  int       `json:"quantity_reserved"`
	QuantityAvailable int       `json:"quantity_available"`
	ReorderLevel      int       `json:"reorder_level"`
	ReorderQuantity   int       `json:"reorder_quantity"`
	IsLowStock        bool      `json:"is_low_stock"`
	LastReceived      *time.Time `json:"last_received"`
	LastSold          *time.Time `json:"last_sold"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type InventoryFilters struct {
	Search       string `json:"search"`
	LowStockOnly bool   `json:"low_stock_only"`
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
	SortBy       string `json:"sort_by"`
	SortDir      string `json:"sort_dir"`
}

type InventoryListResponse struct {
	Items      []InventoryItem `json:"items"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// Dashboard Models
type AdminDashboard struct {
	Stats          DashboardStats    `json:"stats"`
	RecentOrders   []AdminOrder      `json:"recent_orders"`
	TopProducts    []TopProduct      `json:"top_products"`
	SalesChart     SalesChartData    `json:"sales_chart"`
	LowStockAlerts []LowStockAlert   `json:"low_stock_alerts"`
}

type DashboardStats struct {
	TodaysSales     float64 `json:"todays_sales"`
	TodaysOrders    int     `json:"todays_orders"`
	ProductsInStock int     `json:"products_in_stock"`
	ActiveCustomers int     `json:"active_customers"`
	LowStockItems   int     `json:"low_stock_items"`
}

type TopProduct struct {
	ProductID   int     `json:"product_id"`
	Name        string  `json:"name"`
	UnitsSold   int     `json:"units_sold"`
	Revenue     float64 `json:"revenue"`
	Percentage  float64 `json:"percentage"`
}

type SalesChartData struct {
	Labels []string  `json:"labels"`
	Data   []float64 `json:"data"`
}

type LowStockAlert struct {
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	SKU         string `json:"sku"`
	CurrentStock int   `json:"current_stock"`
	ReorderLevel int   `json:"reorder_level"`
	Status       string `json:"status"` // "critical", "low", "warning"
}
