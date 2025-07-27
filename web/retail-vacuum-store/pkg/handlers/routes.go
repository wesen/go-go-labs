package handlers

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// SetupRoutes configures all HTTP routes
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Static files - serve from local directory for development
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	mux.Handle("/static/", staticHandler)

	// API routes
	apiHandler := SetupAPIRoutes()
	mux.Handle("/api/", apiHandler)
	
	// Admin API routes
	adminAPIHandler := SetupAdminAPIRoutes()
	mux.Handle("/api/admin/", adminAPIHandler)

	// Health check
	mux.HandleFunc("/health", healthHandler)

	// Main HTML routes
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/products", productsHandler)
	mux.HandleFunc("/products/", productDetailHandler)
	mux.HandleFunc("/cart", cartHandler)
	mux.HandleFunc("/checkout", checkoutHandler)
	mux.HandleFunc("/account", accountHandler)
	mux.HandleFunc("/admin", adminHandler)

	// Wrap with logging middleware
	return loggingMiddleware(mux)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"retail-vacuum-store"}`))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vacuum Store - Premium Cleaning Solutions</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="/static/css/style.css" rel="stylesheet">
</head>
<body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
        <div class="container">
            <a class="navbar-brand" href="/">🌪️ Vacuum Store</a>
            <div class="navbar-nav ms-auto">
                <a class="nav-link" href="/">Home</a>
                <a class="nav-link" href="/products">Products</a>
                <a class="nav-link" href="/cart">Cart</a>
                <a class="nav-link" href="/admin">Admin</a>
            </div>
        </div>
    </nav>
    
    <div class="container mt-4">
        <div class="jumbotron bg-light p-5 rounded">
            <h1 class="display-4">Welcome to Vacuum Store</h1>
            <p class="lead">Your one-stop shop for premium vacuum cleaners and cleaning solutions.</p>
            <a class="btn btn-primary btn-lg" href="/products" role="button">Shop Now</a>
        </div>
    </div>
    
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <script src="/static/js/app.js"></script>
</body>
</html>
	`))
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Products - VacuumMart</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <link href="/static/css/style.css" rel="stylesheet">
</head>
<body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
        <div class="container">
            <a class="navbar-brand" href="/">🌪️ VacuumMart</a>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarNav">
                <ul class="navbar-nav me-auto">
                    <li class="nav-item"><a class="nav-link" href="/">Home</a></li>
                    <li class="nav-item"><a class="nav-link active" href="/products">Products</a></li>
                </ul>
                <ul class="navbar-nav">
                    <li class="nav-item">
                        <a class="nav-link position-relative" href="/cart">
                            <i class="fas fa-shopping-cart"></i> Cart
                            <span id="cart-count" class="badge bg-danger position-absolute top-0 start-100 translate-middle" style="display: none;">0</span>
                        </a>
                    </li>
                    <li class="nav-item"><a class="nav-link" href="/account">Account</a></li>
                </ul>
            </div>
        </div>
    </nav>
    
    <div class="container mt-4">
        <div class="row">
            <div class="col-lg-3 mb-4">
                <div class="card">
                    <div class="card-header">
                        <h5 class="card-title mb-0">Filters</h5>
                    </div>
                    <div class="card-body">
                        <form id="filters-form">
                            <div class="mb-3">
                                <label for="category-filter" class="form-label">Category</label>
                                <select class="form-select" id="category-filter" name="category_id">
                                    <option value="">All Categories</option>
                                </select>
                            </div>
                            
                            <div class="mb-3">
                                <label for="brand-filter" class="form-label">Brand</label>
                                <select class="form-select" id="brand-filter" name="brand_id">
                                    <option value="">All Brands</option>
                                </select>
                            </div>
                            
                            <div class="mb-3">
                                <label class="form-label">Price Range</label>
                                <div class="row">
                                    <div class="col-6">
                                        <input type="number" class="form-control" name="min_price" placeholder="Min" min="0">
                                    </div>
                                    <div class="col-6">
                                        <input type="number" class="form-control" name="max_price" placeholder="Max" min="0">
                                    </div>
                                </div>
                            </div>
                            
                            <button type="button" class="btn btn-outline-secondary w-100" onclick="this.closest('form').reset(); productManager.applyFilters();">
                                Clear Filters
                            </button>
                        </form>
                    </div>
                </div>
            </div>
            
            <div class="col-lg-9">
                <div class="d-flex justify-content-between align-items-center mb-4">
                    <h1>Our Vacuum Collection</h1>
                    
                    <div class="d-flex gap-3">
                        <div class="search-box">
                            <input type="text" id="search-input" class="form-control" placeholder="Search products..." style="width: 250px;">
                        </div>
                        
                        <select id="sort-select" class="form-select" style="width: auto;">
                            <option value="newest">Newest First</option>
                            <option value="price_asc">Price: Low to High</option>
                            <option value="price_desc">Price: High to Low</option>
                            <option value="name_asc">Name: A to Z</option>
                            <option value="rating">Highest Rated</option>
                        </select>
                    </div>
                </div>
                
                <div id="loading-spinner" class="text-center py-5" style="display: none;">
                    <div class="spinner-border" role="status">
                        <span class="visually-hidden">Loading...</span>
                    </div>
                </div>
                
                <div class="row" id="products-container">
                    <!-- Products will be loaded here -->
                </div>
                
                <div id="pagination-container" class="mt-4">
                    <!-- Pagination will be loaded here -->
                </div>
            </div>
        </div>
    </div>
    
    <footer class="bg-dark text-light mt-5 py-4">
        <div class="container">
            <div class="text-center">
                <p>&copy; 2025 VacuumMart. All rights reserved.</p>
            </div>
        </div>
    </footer>
    
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <script src="/static/js/api.js"></script>
    <script src="/static/js/products.js"></script>
</body>
</html>
	`))
}

func cartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Shopping Cart - VacuumMart</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <link href="/static/css/style.css" rel="stylesheet">
</head>
<body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
        <div class="container">
            <a class="navbar-brand" href="/">🌪️ VacuumMart</a>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarNav">
                <ul class="navbar-nav me-auto">
                    <li class="nav-item"><a class="nav-link" href="/">Home</a></li>
                    <li class="nav-item"><a class="nav-link" href="/products">Products</a></li>
                </ul>
                <ul class="navbar-nav">
                    <li class="nav-item">
                        <a class="nav-link position-relative active" href="/cart">
                            <i class="fas fa-shopping-cart"></i> Cart
                            <span id="cart-count" class="badge bg-danger position-absolute top-0 start-100 translate-middle" style="display: none;">0</span>
                        </a>
                    </li>
                    <li class="nav-item"><a class="nav-link" href="/account">Account</a></li>
                </ul>
            </div>
        </div>
    </nav>
    
    <div class="container mt-4">
        <nav aria-label="breadcrumb">
            <ol class="breadcrumb">
                <li class="breadcrumb-item"><a href="/">Home</a></li>
                <li class="breadcrumb-item active" aria-current="page">Shopping Cart</li>
            </ol>
        </nav>
        
        <div class="row">
            <div class="col-lg-8">
                <div class="card">
                    <div class="card-header d-flex justify-content-between align-items-center">
                        <h4 class="mb-0">Shopping Cart</h4>
                        <span id="cart-item-count" class="badge bg-secondary">0 items</span>
                    </div>
                    <div class="card-body">
                        <div id="cart-items-container">
                            <!-- Cart items will be loaded here -->
                        </div>
                        
                        <div id="empty-cart" style="display: none;" class="text-center py-5">
                            <i class="fas fa-shopping-cart fa-3x text-muted mb-3"></i>
                            <h5 class="text-muted">Your cart is empty</h5>
                            <p class="text-muted mb-4">Add some products to get started!</p>
                            <a href="/products" class="btn btn-primary">Continue Shopping</a>
                        </div>
                    </div>
                </div>
            </div>
            
            <div class="col-lg-4">
                <div class="card">
                    <div class="card-header">
                        <h5 class="mb-0">Order Summary</h5>
                    </div>
                    <div class="card-body">
                        <div class="d-flex justify-content-between mb-2">
                            <span>Subtotal:</span>
                            <span id="cart-subtotal">$0.00</span>
                        </div>
                        <div class="d-flex justify-content-between mb-2">
                            <span>Shipping:</span>
                            <span id="cart-shipping">Free</span>
                        </div>
                        <div class="d-flex justify-content-between mb-2">
                            <span>Tax:</span>
                            <span id="cart-tax">$0.00</span>
                        </div>
                        <hr>
                        <div class="d-flex justify-content-between fw-bold mb-3">
                            <span>Total:</span>
                            <span id="cart-total">$0.00</span>
                        </div>
                        
                        <div class="d-grid gap-2">
                            <a href="/checkout" id="checkout-btn" class="btn btn-success btn-lg" style="display: none;">
                                <i class="fas fa-lock me-2"></i>Secure Checkout
                            </a>
                            <a href="/products" class="btn btn-outline-primary">Continue Shopping</a>
                        </div>
                        
                        <div class="mt-3 text-center">
                            <small class="text-muted">
                                <i class="fas fa-shield-alt me-1"></i>
                                Secure SSL encrypted checkout
                            </small>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    
    <footer class="bg-dark text-light mt-5 py-4">
        <div class="container">
            <div class="text-center">
                <p>&copy; 2025 VacuumMart. All rights reserved.</p>
            </div>
        </div>
    </footer>
    
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <script src="/static/js/api.js"></script>
    <script src="/static/js/cart.js"></script>
</body>
</html>
	`))
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the admin dashboard HTML file
	http.ServeFile(w, r, "./static/admin-dashboard.html")
}

func productDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Product Details - VacuumMart</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <link href="/static/css/style.css" rel="stylesheet">
</head>
<body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
        <div class="container">
            <a class="navbar-brand" href="/">🌪️ VacuumMart</a>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarNav">
                <ul class="navbar-nav me-auto">
                    <li class="nav-item"><a class="nav-link" href="/">Home</a></li>
                    <li class="nav-item"><a class="nav-link" href="/products">Products</a></li>
                </ul>
                <ul class="navbar-nav">
                    <li class="nav-item">
                        <a class="nav-link position-relative" href="/cart">
                            <i class="fas fa-shopping-cart"></i> Cart
                            <span id="cart-count" class="badge bg-danger position-absolute top-0 start-100 translate-middle" style="display: none;">0</span>
                        </a>
                    </li>
                    <li class="nav-item"><a class="nav-link" href="/account">Account</a></li>
                </ul>
            </div>
        </div>
    </nav>
    
    <div class="container mt-4">
        <nav aria-label="breadcrumb">
            <ol class="breadcrumb">
                <li class="breadcrumb-item"><a href="/">Home</a></li>
                <li class="breadcrumb-item"><a href="/products">Products</a></li>
                <li class="breadcrumb-item active" aria-current="page" id="product-breadcrumb">Loading...</li>
            </ol>
        </nav>
        
        <div id="loading-spinner" class="text-center py-5">
            <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
        
        <div id="product-detail-container" style="display: none;">
            <!-- Product details will be loaded here -->
        </div>
        
        <div id="error-container" style="display: none;" class="text-center py-5">
            <i class="fas fa-exclamation-triangle fa-3x text-warning mb-3"></i>
            <h4 class="text-warning">Product Not Found</h4>
            <p class="text-muted">The product you're looking for doesn't exist.</p>
            <a href="/products" class="btn btn-primary">Browse Products</a>
        </div>
    </div>
    
    <footer class="bg-dark text-light mt-5 py-4">
        <div class="container">
            <div class="text-center">
                <p>&copy; 2025 VacuumMart. All rights reserved.</p>
            </div>
        </div>
    </footer>
    
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <script src="/static/js/api.js"></script>
    <script src="/static/js/product-detail.js"></script>
</body>
</html>
	`))
}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Checkout - VacuumMart</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <link href="/static/css/style.css" rel="stylesheet">
</head>
<body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
        <div class="container">
            <a class="navbar-brand" href="/">🌪️ VacuumMart</a>
            <div class="navbar-nav ms-auto">
                <a class="nav-link" href="/">Home</a>
                <a class="nav-link" href="/products">Products</a>
                <a class="nav-link" href="/cart">Cart</a>
            </div>
        </div>
    </nav>
    
    <div class="container mt-4">
        <div class="row">
            <div class="col-lg-8">
                <div class="card">
                    <div class="card-header">
                        <h4 class="mb-0">Checkout</h4>
                    </div>
                    <div class="card-body">
                        <form id="checkout-form">
                            <div class="row mb-4">
                                <div class="col-md-6">
                                    <h5>Contact Information</h5>
                                    <div class="mb-3">
                                        <label for="email" class="form-label">Email Address</label>
                                        <input type="email" class="form-control" id="email" name="email" required>
                                    </div>
                                </div>
                            </div>
                            
                            <div class="row mb-4">
                                <div class="col-12">
                                    <h5>Shipping Address</h5>
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="first_name" class="form-label">First Name</label>
                                    <input type="text" class="form-control" id="first_name" name="first_name" required>
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="last_name" class="form-label">Last Name</label>
                                    <input type="text" class="form-control" id="last_name" name="last_name" required>
                                </div>
                                <div class="col-12 mb-3">
                                    <label for="address_line1" class="form-label">Address</label>
                                    <input type="text" class="form-control" id="address_line1" name="address_line1" required>
                                </div>
                                <div class="col-md-6 mb-3">
                                    <label for="city" class="form-label">City</label>
                                    <input type="text" class="form-control" id="city" name="city" required>
                                </div>
                                <div class="col-md-3 mb-3">
                                    <label for="state" class="form-label">State</label>
                                    <input type="text" class="form-control" id="state" name="state" required>
                                </div>
                                <div class="col-md-3 mb-3">
                                    <label for="postal_code" class="form-label">ZIP Code</label>
                                    <input type="text" class="form-control" id="postal_code" name="postal_code" required>
                                </div>
                            </div>
                            
                            <div class="row mb-4">
                                <div class="col-12">
                                    <h5>Payment Method</h5>
                                </div>
                                <div class="col-12">
                                    <div class="form-check">
                                        <input class="form-check-input" type="radio" name="payment_method" id="credit_card" value="credit_card" checked>
                                        <label class="form-check-label" for="credit_card">
                                            <i class="fas fa-credit-card me-2"></i>Credit Card
                                        </label>
                                    </div>
                                    <div class="form-check">
                                        <input class="form-check-input" type="radio" name="payment_method" id="paypal" value="paypal">
                                        <label class="form-check-label" for="paypal">
                                            <i class="fab fa-paypal me-2"></i>PayPal
                                        </label>
                                    </div>
                                </div>
                            </div>
                            
                            <button type="submit" class="btn btn-success btn-lg w-100">
                                <i class="fas fa-lock me-2"></i>Place Order
                            </button>
                        </form>
                    </div>
                </div>
            </div>
            
            <div class="col-lg-4">
                <div class="card">
                    <div class="card-header">
                        <h5 class="mb-0">Order Summary</h5>
                    </div>
                    <div class="card-body" id="order-summary">
                        <!-- Order summary will be loaded here -->
                    </div>
                </div>
            </div>
        </div>
    </div>
    
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <script src="/static/js/api.js"></script>
    <script src="/static/js/checkout.js"></script>
</body>
</html>
	`))
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My Account - VacuumMart</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <link href="/static/css/style.css" rel="stylesheet">
</head>
<body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
        <div class="container">
            <a class="navbar-brand" href="/">🌪️ VacuumMart</a>
            <div class="navbar-nav ms-auto">
                <a class="nav-link" href="/">Home</a>
                <a class="nav-link" href="/products">Products</a>
                <a class="nav-link" href="/cart">Cart</a>
                <a class="nav-link active" href="/account">Account</a>
            </div>
        </div>
    </nav>
    
    <div class="container mt-4">
        <div class="row">
            <div class="col-md-3">
                <div class="list-group">
                    <a href="#profile" class="list-group-item list-group-item-action active" data-bs-toggle="tab">
                        <i class="fas fa-user me-2"></i>Profile
                    </a>
                    <a href="#orders" class="list-group-item list-group-item-action" data-bs-toggle="tab">
                        <i class="fas fa-shopping-bag me-2"></i>Order History
                    </a>
                    <a href="#addresses" class="list-group-item list-group-item-action" data-bs-toggle="tab">
                        <i class="fas fa-map-marker-alt me-2"></i>Addresses
                    </a>
                </div>
            </div>
            
            <div class="col-md-9">
                <div class="tab-content">
                    <div class="tab-pane fade show active" id="profile">
                        <div class="card">
                            <div class="card-header">
                                <h4 class="mb-0">Profile Information</h4>
                            </div>
                            <div class="card-body">
                                <form id="profile-form">
                                    <div class="row">
                                        <div class="col-md-6 mb-3">
                                            <label for="profile_first_name" class="form-label">First Name</label>
                                            <input type="text" class="form-control" id="profile_first_name" name="first_name">
                                        </div>
                                        <div class="col-md-6 mb-3">
                                            <label for="profile_last_name" class="form-label">Last Name</label>
                                            <input type="text" class="form-control" id="profile_last_name" name="last_name">
                                        </div>
                                        <div class="col-md-6 mb-3">
                                            <label for="profile_email" class="form-label">Email</label>
                                            <input type="email" class="form-control" id="profile_email" name="email">
                                        </div>
                                        <div class="col-md-6 mb-3">
                                            <label for="profile_phone" class="form-label">Phone</label>
                                            <input type="tel" class="form-control" id="profile_phone" name="phone">
                                        </div>
                                    </div>
                                    <button type="submit" class="btn btn-primary">Save Changes</button>
                                </form>
                            </div>
                        </div>
                    </div>
                    
                    <div class="tab-pane fade" id="orders">
                        <div class="card">
                            <div class="card-header">
                                <h4 class="mb-0">Order History</h4>
                            </div>
                            <div class="card-body">
                                <div id="orders-container">
                                    <!-- Orders will be loaded here -->
                                </div>
                            </div>
                        </div>
                    </div>
                    
                    <div class="tab-pane fade" id="addresses">
                        <div class="card">
                            <div class="card-header">
                                <h4 class="mb-0">Saved Addresses</h4>
                            </div>
                            <div class="card-body">
                                <p class="text-muted">Address management coming soon.</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <script src="/static/js/api.js"></script>
    <script src="/static/js/account.js"></script>
</body>
</html>
	`))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("remote_addr", r.RemoteAddr).
			Msg("HTTP request")
		
		next.ServeHTTP(w, r)
	})
}
