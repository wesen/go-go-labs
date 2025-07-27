// Admin Dashboard JavaScript
class AdminDashboard {
    constructor() {
        this.currentView = 'dashboard';
        this.isAuthenticated = false;
        this.adminUser = null;
        this.orders = [];
        this.products = [];
        this.customers = [];
        this.inventory = [];
        
        this.init();
    }

    init() {
        this.checkAuth();
        this.setupEventListeners();
        this.loadDashboard();
    }

    checkAuth() {
        // Check if admin is logged in (check cookie or localStorage)
        const token = this.getCookie('admin_session');
        if (token) {
            this.isAuthenticated = true;
            this.loadUserProfile();
        } else {
            this.showLogin();
        }
    }

    setupEventListeners() {
        // Navigation
        document.addEventListener('click', (e) => {
            if (e.target.matches('.vm-admin-nav-link')) {
                e.preventDefault();
                this.handleNavigation(e.target);
            }
        });

        // Login form
        document.addEventListener('submit', (e) => {
            if (e.target.matches('#admin-login-form')) {
                e.preventDefault();
                this.handleLogin(e.target);
            }
        });

        // Logout
        document.addEventListener('click', (e) => {
            if (e.target.matches('#logout-btn') || e.target.closest('#logout-btn')) {
                e.preventDefault();
                this.handleLogout();
            }
        });

        // Modal handlers
        document.addEventListener('click', (e) => {
            if (e.target.matches('[data-bs-toggle="modal"]')) {
                this.handleModalOpen(e.target);
            }
        });

        // Form submissions
        document.addEventListener('submit', (e) => {
            if (e.target.matches('.admin-form')) {
                e.preventDefault();
                this.handleFormSubmit(e.target);
            }
        });

        // Status updates
        document.addEventListener('change', (e) => {
            if (e.target.matches('.order-status-select')) {
                this.handleOrderStatusChange(e.target);
            }
        });

        // Search and filters
        document.addEventListener('input', (e) => {
            if (e.target.matches('.search-input')) {
                this.debounce(() => this.handleSearch(e.target), 300)();
            }
        });

        // Pagination
        document.addEventListener('click', (e) => {
            if (e.target.matches('.page-link')) {
                e.preventDefault();
                this.handlePagination(e.target);
            }
        });
    }

    async handleLogin(form) {
        const formData = new FormData(form);
        const credentials = {
            username: formData.get('username'),
            password: formData.get('password')
        };

        try {
            const response = await this.apiCall('/api/admin/login', 'POST', credentials);
            
            if (response.success) {
                this.isAuthenticated = true;
                this.adminUser = response.data.user;
                this.hideLogin();
                this.loadDashboard();
                this.showAlert('success', 'Login successful');
            } else {
                this.showAlert('danger', response.error || 'Login failed');
            }
        } catch (error) {
            this.showAlert('danger', 'Login failed: ' + error.message);
        }
    }

    async handleLogout() {
        try {
            await this.apiCall('/api/admin/logout', 'POST');
            this.isAuthenticated = false;
            this.adminUser = null;
            this.showLogin();
            this.showAlert('info', 'Logged out successfully');
        } catch (error) {
            console.error('Logout error:', error);
        }
    }

    handleNavigation(navLink) {
        // Remove active class from all nav links
        document.querySelectorAll('.vm-admin-nav-link').forEach(link => {
            link.classList.remove('active');
        });
        
        // Add active class to clicked nav link
        navLink.classList.add('active');
        
        const navText = navLink.textContent.trim().toLowerCase();
        
        switch (navText) {
            case 'dashboard':
                this.loadDashboard();
                break;
            case 'products':
                this.loadProducts();
                break;
            case 'orders':
                this.loadOrders();
                break;
            case 'customers':
                this.loadCustomers();
                break;
            case 'inventory':
                this.loadInventory();
                break;
            case 'analytics':
                this.loadAnalytics();
                break;
            default:
                this.loadDashboard();
        }
    }

    async loadDashboard() {
        this.currentView = 'dashboard';
        this.setPageTitle('Dashboard');
        
        try {
            const response = await this.apiCall('/api/admin/analytics/dashboard');
            if (response.success) {
                this.renderDashboard(response.data);
            }
        } catch (error) {
            this.showAlert('danger', 'Failed to load dashboard: ' + error.message);
        }
    }

    async loadProducts() {
        this.currentView = 'products';
        this.setPageTitle('Products');
        
        try {
            const response = await this.apiCall('/api/admin/products?page=1&limit=20');
            if (response.success) {
                this.products = response.data.products;
                this.renderProducts(response.data);
            }
        } catch (error) {
            this.showAlert('danger', 'Failed to load products: ' + error.message);
        }
    }

    async loadOrders() {
        this.currentView = 'orders';
        this.setPageTitle('Orders');
        
        try {
            const response = await this.apiCall('/api/admin/orders?page=1&limit=20');
            if (response.success) {
                this.orders = response.data.orders;
                this.renderOrders(response.data);
            }
        } catch (error) {
            this.showAlert('danger', 'Failed to load orders: ' + error.message);
        }
    }

    async loadCustomers() {
        this.currentView = 'customers';
        this.setPageTitle('Customers');
        
        try {
            const response = await this.apiCall('/api/admin/customers?page=1&limit=20');
            if (response.success) {
                this.customers = response.data.customers;
                this.renderCustomers(response.data);
            }
        } catch (error) {
            this.showAlert('danger', 'Failed to load customers: ' + error.message);
        }
    }

    async loadInventory() {
        this.currentView = 'inventory';
        this.setPageTitle('Inventory');
        
        try {
            const response = await this.apiCall('/api/admin/inventory?page=1&limit=20');
            if (response.success) {
                this.inventory = response.data.items;
                this.renderInventory(response.data);
            }
        } catch (error) {
            this.showAlert('danger', 'Failed to load inventory: ' + error.message);
        }
    }

    renderDashboard(data) {
        const content = `
            <!-- Stats Cards -->
            <div class="row g-4 mb-4">
                <div class="col-md-3">
                    <div class="vm-stat-card">
                        <div class="d-flex justify-content-between align-items-center">
                            <div>
                                <div class="vm-stat-number">$${data.stats.todays_sales.toLocaleString()}</div>
                                <div class="vm-stat-label">Today's Sales</div>
                            </div>
                            <i class="fas fa-dollar-sign fa-2x opacity-75"></i>
                        </div>
                    </div>
                </div>
                <div class="col-md-3">
                    <div class="vm-stat-card" style="background: linear-gradient(135deg, var(--vm-accent), #059669);">
                        <div class="d-flex justify-content-between align-items-center">
                            <div>
                                <div class="vm-stat-number">${data.stats.todays_orders}</div>
                                <div class="vm-stat-label">Orders Today</div>
                            </div>
                            <i class="fas fa-shopping-cart fa-2x opacity-75"></i>
                        </div>
                    </div>
                </div>
                <div class="col-md-3">
                    <div class="vm-stat-card" style="background: linear-gradient(135deg, var(--vm-warning), #d97706);">
                        <div class="d-flex justify-content-between align-items-center">
                            <div>
                                <div class="vm-stat-number">${data.stats.products_in_stock}</div>
                                <div class="vm-stat-label">Products in Stock</div>
                            </div>
                            <i class="fas fa-boxes fa-2x opacity-75"></i>
                        </div>
                        <div class="mt-2">
                            <small class="text-white opacity-75">
                                <i class="fas fa-exclamation-triangle"></i> ${data.stats.low_stock_items} low stock items
                            </small>
                        </div>
                    </div>
                </div>
                <div class="col-md-3">
                    <div class="vm-stat-card" style="background: linear-gradient(135deg, #6366f1, #4f46e5);">
                        <div class="d-flex justify-content-between align-items-center">
                            <div>
                                <div class="vm-stat-number">${data.stats.active_customers}</div>
                                <div class="vm-stat-label">Active Customers</div>
                            </div>
                            <i class="fas fa-users fa-2x opacity-75"></i>
                        </div>
                    </div>
                </div>
            </div>
            
            <div class="row g-4">
                <!-- Recent Orders -->
                <div class="col-lg-8">
                    <div class="card vm-admin-card">
                        <div class="card-header d-flex justify-content-between align-items-center">
                            <h5 class="mb-0">Recent Orders</h5>
                            <button class="btn btn-sm btn-outline-primary" onclick="adminDashboard.loadOrders()">View All</button>
                        </div>
                        <div class="card-body p-0">
                            <div class="table-responsive">
                                <table class="table table-hover mb-0">
                                    <thead class="table-light">
                                        <tr>
                                            <th>Order ID</th>
                                            <th>Customer</th>
                                            <th>Total</th>
                                            <th>Status</th>
                                            <th>Actions</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        ${data.recent_orders.slice(0, 5).map(order => `
                                            <tr>
                                                <td>#${order.order_number}</td>
                                                <td>${order.customer_name || order.customer_email}</td>
                                                <td>$${order.total_amount.toFixed(2)}</td>
                                                <td><span class="badge bg-${this.getStatusColor(order.status)}">${order.status}</span></td>
                                                <td>
                                                    <button class="btn btn-sm btn-outline-primary" onclick="adminDashboard.viewOrderDetail(${order.id})">View</button>
                                                </td>
                                            </tr>
                                        `).join('')}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
                
                <!-- Low Stock Alerts -->
                <div class="col-lg-4">
                    <div class="card vm-admin-card">
                        <div class="card-header bg-warning text-dark">
                            <h6 class="mb-0">
                                <i class="fas fa-exclamation-triangle me-2"></i>Low Stock Alert
                            </h6>
                        </div>
                        <div class="card-body">
                            ${data.low_stock_alerts.map(alert => `
                                <div class="mb-3">
                                    <div class="d-flex justify-content-between align-items-center">
                                        <span>${alert.product_name}</span>
                                        <span class="badge bg-${alert.status === 'critical' ? 'danger' : 'warning'}">${alert.current_stock} left</span>
                                    </div>
                                </div>
                            `).join('')}
                            <button class="btn btn-sm btn-outline-warning w-100" onclick="adminDashboard.loadInventory()">
                                Manage Inventory
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `;
        
        this.setMainContent(content);
    }

    renderProducts(data) {
        const content = `
            <div class="d-flex justify-content-between align-items-center mb-4">
                <h2>Products</h2>
                <div class="d-flex gap-2">
                    <input type="text" class="form-control search-input" placeholder="Search products..." style="width: 250px;">
                    <button class="btn btn-primary" onclick="adminDashboard.showProductModal()">
                        <i class="fas fa-plus me-1"></i> Add Product
                    </button>
                </div>
            </div>
            
            <div class="card vm-admin-card">
                <div class="card-body p-0">
                    <div class="table-responsive">
                        <table class="table table-hover mb-0">
                            <thead class="table-light">
                                <tr>
                                    <th>SKU</th>
                                    <th>Name</th>
                                    <th>Category</th>
                                    <th>Brand</th>
                                    <th>Price</th>
                                    <th>Stock</th>
                                    <th>Status</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                ${data.products.map(product => `
                                    <tr>
                                        <td><code>${product.sku}</code></td>
                                        <td>
                                            <div class="fw-medium">${product.name}</div>
                                            ${product.low_stock_alert ? '<small class="text-danger"><i class="fas fa-exclamation-triangle"></i> Low Stock</small>' : ''}
                                        </td>
                                        <td>${product.category_name || '-'}</td>
                                        <td>${product.brand_name || '-'}</td>
                                        <td>
                                            <div>$${product.base_price.toFixed(2)}</div>
                                            ${product.sale_price ? `<small class="text-success">Sale: $${product.sale_price.toFixed(2)}</small>` : ''}
                                        </td>
                                        <td>
                                            <span class="badge bg-${product.stock > 10 ? 'success' : product.stock > 0 ? 'warning' : 'danger'}">
                                                ${product.stock || 0}
                                            </span>
                                        </td>
                                        <td>
                                            <span class="badge bg-${product.is_active ? 'success' : 'secondary'}">
                                                ${product.is_active ? 'Active' : 'Inactive'}
                                            </span>
                                        </td>
                                        <td>
                                            <div class="btn-group btn-group-sm">
                                                <button class="btn btn-outline-primary" onclick="adminDashboard.editProduct(${product.id})">
                                                    <i class="fas fa-edit"></i>
                                                </button>
                                                <button class="btn btn-outline-danger" onclick="adminDashboard.deleteProduct(${product.id})">
                                                    <i class="fas fa-trash"></i>
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
            
            ${this.renderPagination(data)}
        `;
        
        this.setMainContent(content);
    }

    renderOrders(data) {
        const content = `
            <div class="d-flex justify-content-between align-items-center mb-4">
                <h2>Orders</h2>
                <div class="d-flex gap-2">
                    <select class="form-select" onchange="adminDashboard.filterOrders(this.value)" style="width: auto;">
                        <option value="">All Orders</option>
                        <option value="pending">Pending</option>
                        <option value="confirmed">Confirmed</option>
                        <option value="processing">Processing</option>
                        <option value="shipped">Shipped</option>
                        <option value="delivered">Delivered</option>
                        <option value="cancelled">Cancelled</option>
                    </select>
                    <input type="text" class="form-control search-input" placeholder="Search orders..." style="width: 250px;">
                </div>
            </div>
            
            <div class="card vm-admin-card">
                <div class="card-body p-0">
                    <div class="table-responsive">
                        <table class="table table-hover mb-0">
                            <thead class="table-light">
                                <tr>
                                    <th>Order #</th>
                                    <th>Customer</th>
                                    <th>Date</th>
                                    <th>Total</th>
                                    <th>Payment</th>
                                    <th>Status</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                ${data.orders.map(order => `
                                    <tr>
                                        <td><code>#${order.order_number}</code></td>
                                        <td>
                                            <div class="fw-medium">${order.customer_name || 'Guest'}</div>
                                            <small class="text-muted">${order.customer_email}</small>
                                        </td>
                                        <td>${new Date(order.created_at).toLocaleDateString()}</td>
                                        <td>$${order.total_amount.toFixed(2)}</td>
                                        <td>
                                            <span class="badge bg-${this.getPaymentStatusColor(order.payment_status)}">
                                                ${order.payment_status}
                                            </span>
                                        </td>
                                        <td>
                                            <select class="form-select form-select-sm order-status-select" data-order-id="${order.id}">
                                                <option value="pending" ${order.status === 'pending' ? 'selected' : ''}>Pending</option>
                                                <option value="confirmed" ${order.status === 'confirmed' ? 'selected' : ''}>Confirmed</option>
                                                <option value="processing" ${order.status === 'processing' ? 'selected' : ''}>Processing</option>
                                                <option value="shipped" ${order.status === 'shipped' ? 'selected' : ''}>Shipped</option>
                                                <option value="delivered" ${order.status === 'delivered' ? 'selected' : ''}>Delivered</option>
                                                <option value="cancelled" ${order.status === 'cancelled' ? 'selected' : ''}>Cancelled</option>
                                            </select>
                                        </td>
                                        <td>
                                            <button class="btn btn-sm btn-outline-primary" onclick="adminDashboard.viewOrderDetail(${order.id})">
                                                <i class="fas fa-eye"></i> View
                                            </button>
                                        </td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
            
            ${this.renderPagination(data)}
        `;
        
        this.setMainContent(content);
    }

    renderCustomers(data) {
        const content = `
            <div class="d-flex justify-content-between align-items-center mb-4">
                <h2>Customers</h2>
                <div class="d-flex gap-2">
                    <select class="form-select" style="width: auto;">
                        <option value="">All Customers</option>
                        <option value="active">Active</option>
                        <option value="inactive">Inactive</option>
                    </select>
                    <input type="text" class="form-control search-input" placeholder="Search customers..." style="width: 250px;">
                </div>
            </div>
            
            <div class="card vm-admin-card">
                <div class="card-body p-0">
                    <div class="table-responsive">
                        <table class="table table-hover mb-0">
                            <thead class="table-light">
                                <tr>
                                    <th>Customer</th>
                                    <th>Email</th>
                                    <th>Phone</th>
                                    <th>Orders</th>
                                    <th>Total Spent</th>
                                    <th>Status</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                ${data.customers.map(customer => `
                                    <tr>
                                        <td>
                                            <div class="fw-medium">${customer.first_name} ${customer.last_name}</div>
                                            <small class="text-muted">
                                                ${customer.email_verified ? '<i class="fas fa-check-circle text-success"></i>' : '<i class="fas fa-exclamation-circle text-warning"></i>'}
                                                ${customer.email_verified ? 'Verified' : 'Unverified'}
                                            </small>
                                        </td>
                                        <td>${customer.email}</td>
                                        <td>${customer.phone || '-'}</td>
                                        <td>${customer.order_count}</td>
                                        <td>$${customer.total_spent.toFixed(2)}</td>
                                        <td>
                                            <span class="badge bg-${customer.is_active ? 'success' : 'secondary'}">
                                                ${customer.is_active ? 'Active' : 'Inactive'}
                                            </span>
                                        </td>
                                        <td>
                                            <div class="btn-group btn-group-sm">
                                                <button class="btn btn-outline-primary" onclick="adminDashboard.viewCustomer(${customer.id})">
                                                    <i class="fas fa-eye"></i>
                                                </button>
                                                <button class="btn btn-outline-secondary" onclick="adminDashboard.emailCustomer(${customer.id})">
                                                    <i class="fas fa-envelope"></i>
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
            
            ${this.renderPagination(data)}
        `;
        
        this.setMainContent(content);
    }

    renderInventory(data) {
        const content = `
            <div class="d-flex justify-content-between align-items-center mb-4">
                <h2>Inventory</h2>
                <div class="d-flex gap-2">
                    <label class="form-check">
                        <input type="checkbox" class="form-check-input" onchange="adminDashboard.filterLowStock(this.checked)">
                        <span class="form-check-label">Low Stock Only</span>
                    </label>
                    <input type="text" class="form-control search-input" placeholder="Search inventory..." style="width: 250px;">
                    <button class="btn btn-primary" onclick="adminDashboard.bulkInventoryUpdate()">
                        <i class="fas fa-upload me-1"></i> Bulk Update
                    </button>
                </div>
            </div>
            
            <div class="card vm-admin-card">
                <div class="card-body p-0">
                    <div class="table-responsive">
                        <table class="table table-hover mb-0">
                            <thead class="table-light">
                                <tr>
                                    <th>SKU</th>
                                    <th>Product</th>
                                    <th>Location</th>
                                    <th>On Hand</th>
                                    <th>Reserved</th>
                                    <th>Available</th>
                                    <th>Reorder Level</th>
                                    <th>Status</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                ${data.items.map(item => `
                                    <tr class="${item.is_low_stock ? 'table-warning' : ''}">
                                        <td><code>${item.sku}</code></td>
                                        <td>
                                            <div class="fw-medium">${item.product_name}</div>
                                            ${item.is_low_stock ? '<small class="text-danger"><i class="fas fa-exclamation-triangle"></i> Low Stock</small>' : ''}
                                        </td>
                                        <td>${item.warehouse_location}</td>
                                        <td>${item.quantity_on_hand}</td>
                                        <td>${item.quantity_reserved}</td>
                                        <td>
                                            <span class="badge bg-${item.quantity_available > item.reorder_level ? 'success' : item.quantity_available > 0 ? 'warning' : 'danger'}">
                                                ${item.quantity_available}
                                            </span>
                                        </td>
                                        <td>${item.reorder_level}</td>
                                        <td>
                                            ${item.is_low_stock ? 
                                                '<span class="badge bg-warning">Low Stock</span>' : 
                                                '<span class="badge bg-success">In Stock</span>'
                                            }
                                        </td>
                                        <td>
                                            <div class="btn-group btn-group-sm">
                                                <button class="btn btn-outline-primary" onclick="adminDashboard.adjustInventory(${item.id})">
                                                    <i class="fas fa-edit"></i>
                                                </button>
                                                <button class="btn btn-outline-success" onclick="adminDashboard.reorderInventory(${item.id})">
                                                    <i class="fas fa-shopping-cart"></i>
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
            
            ${this.renderPagination(data)}
        `;
        
        this.setMainContent(content);
    }

    async handleOrderStatusChange(select) {
        const orderId = select.dataset.orderId;
        const newStatus = select.value;
        const originalStatus = select.dataset.originalStatus || select.value;
        
        try {
            const response = await this.apiCall('/api/admin/orders/status', 'PUT', {
                order_id: parseInt(orderId),
                status: newStatus,
                notes: `Status changed from ${originalStatus} to ${newStatus}`
            });
            
            if (response.success) {
                this.showAlert('success', 'Order status updated successfully');
                select.dataset.originalStatus = newStatus;
            } else {
                select.value = originalStatus; // Revert on error
                this.showAlert('danger', response.error || 'Failed to update order status');
            }
        } catch (error) {
            select.value = originalStatus; // Revert on error
            this.showAlert('danger', 'Failed to update order status: ' + error.message);
        }
    }

    async viewOrderDetail(orderId) {
        try {
            const response = await this.apiCall(`/api/admin/orders/${orderId}`);
            if (response.success) {
                this.showOrderDetailModal(response.data);
            }
        } catch (error) {
            this.showAlert('danger', 'Failed to load order details: ' + error.message);
        }
    }

    showOrderDetailModal(order) {
        const modalContent = `
            <div class="modal fade" id="orderDetailModal" tabindex="-1">
                <div class="modal-dialog modal-lg">
                    <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title">Order #${order.order_number}</h5>
                            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                        </div>
                        <div class="modal-body">
                            <div class="row">
                                <div class="col-md-6">
                                    <h6>Customer Information</h6>
                                    <p><strong>Name:</strong> ${order.customer_name || 'Guest'}<br>
                                    <strong>Email:</strong> ${order.customer_email}<br>
                                    <strong>Order Date:</strong> ${new Date(order.created_at).toLocaleString()}</p>
                                </div>
                                <div class="col-md-6">
                                    <h6>Order Summary</h6>
                                    <p><strong>Status:</strong> <span class="badge bg-${this.getStatusColor(order.status)}">${order.status}</span><br>
                                    <strong>Payment:</strong> <span class="badge bg-${this.getPaymentStatusColor(order.payment_status)}">${order.payment_status}</span><br>
                                    <strong>Total:</strong> $${order.total_amount.toFixed(2)}</p>
                                </div>
                            </div>
                            
                            <h6>Order Items</h6>
                            <div class="table-responsive">
                                <table class="table table-sm">
                                    <thead>
                                        <tr>
                                            <th>Product</th>
                                            <th>SKU</th>
                                            <th>Qty</th>
                                            <th>Price</th>
                                            <th>Total</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        ${(order.items || []).map(item => `
                                            <tr>
                                                <td>${item.product_name}</td>
                                                <td><code>${item.sku}</code></td>
                                                <td>${item.quantity}</td>
                                                <td>$${item.unit_price.toFixed(2)}</td>
                                                <td>$${item.total_price.toFixed(2)}</td>
                                            </tr>
                                        `).join('')}
                                    </tbody>
                                </table>
                            </div>
                            
                            <h6>Status History</h6>
                            <div class="timeline">
                                ${(order.status_history || []).map(entry => `
                                    <div class="timeline-item">
                                        <div class="timeline-marker"></div>
                                        <div class="timeline-content">
                                            <h6>${entry.status}</h6>
                                            <p>${entry.notes}</p>
                                            <small class="text-muted">${new Date(entry.created_at).toLocaleString()}</small>
                                        </div>
                                    </div>
                                `).join('')}
                            </div>
                        </div>
                        <div class="modal-footer">
                            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                            <button type="button" class="btn btn-primary" onclick="adminDashboard.printOrder(${order.id})">
                                <i class="fas fa-print"></i> Print
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `;
        
        // Remove existing modal if any
        const existingModal = document.getElementById('orderDetailModal');
        if (existingModal) {
            existingModal.remove();
        }
        
        // Add modal to body
        document.body.insertAdjacentHTML('beforeend', modalContent);
        
        // Show modal
        const modal = new bootstrap.Modal(document.getElementById('orderDetailModal'));
        modal.show();
    }

    renderPagination(data) {
        if (data.total_pages <= 1) return '';
        
        let pages = [];
        for (let i = 1; i <= data.total_pages; i++) {
            pages.push(i);
        }
        
        return `
            <nav aria-label="Page navigation" class="mt-4">
                <ul class="pagination justify-content-center">
                    <li class="page-item ${data.page === 1 ? 'disabled' : ''}">
                        <a class="page-link" href="#" data-page="${data.page - 1}">Previous</a>
                    </li>
                    ${pages.map(page => `
                        <li class="page-item ${page === data.page ? 'active' : ''}">
                            <a class="page-link" href="#" data-page="${page}">${page}</a>
                        </li>
                    `).join('')}
                    <li class="page-item ${data.page === data.total_pages ? 'disabled' : ''}">
                        <a class="page-link" href="#" data-page="${data.page + 1}">Next</a>
                    </li>
                </ul>
            </nav>
        `;
    }

    // Utility methods
    setPageTitle(title) {
        const titleElement = document.querySelector('.container-fluid h2');
        if (titleElement) {
            titleElement.textContent = title;
        }
    }

    setMainContent(content) {
        const mainContent = document.querySelector('.col-md-9.col-lg-10 .p-4');
        if (mainContent) {
            mainContent.innerHTML = content;
        }
    }

    getStatusColor(status) {
        const colors = {
            'pending': 'warning',
            'confirmed': 'info',
            'processing': 'primary',
            'shipped': 'info',
            'delivered': 'success',
            'cancelled': 'danger',
            'refunded': 'secondary'
        };
        return colors[status] || 'secondary';
    }

    getPaymentStatusColor(status) {
        const colors = {
            'pending': 'warning',
            'paid': 'success',
            'failed': 'danger',
            'refunded': 'secondary',
            'partially_refunded': 'warning'
        };
        return colors[status] || 'secondary';
    }

    showAlert(type, message) {
        const alertContainer = document.querySelector('.alert-container') || document.body;
        const alertElement = document.createElement('div');
        alertElement.className = `alert alert-${type} alert-dismissible fade show position-fixed`;
        alertElement.style.cssText = 'top: 20px; right: 20px; z-index: 9999; min-width: 300px;';
        alertElement.innerHTML = `
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        `;
        
        alertContainer.appendChild(alertElement);
        
        // Auto-dismiss after 5 seconds
        setTimeout(() => {
            if (alertElement.parentNode) {
                alertElement.remove();
            }
        }, 5000);
    }

    showLogin() {
        const loginContent = `
            <div class="row justify-content-center">
                <div class="col-md-4">
                    <div class="card">
                        <div class="card-header text-center">
                            <h4><i class="fas fa-wind me-2"></i>VacuumMart Admin</h4>
                        </div>
                        <div class="card-body">
                            <form id="admin-login-form">
                                <div class="mb-3">
                                    <label class="form-label">Username</label>
                                    <input type="text" class="form-control" name="username" required>
                                </div>
                                <div class="mb-3">
                                    <label class="form-label">Password</label>
                                    <input type="password" class="form-control" name="password" required>
                                </div>
                                <button type="submit" class="btn btn-primary w-100">Login</button>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        `;
        
        document.querySelector('.container-fluid').innerHTML = loginContent;
    }

    hideLogin() {
        // Reload the page to show the full admin interface
        window.location.reload();
    }

    async apiCall(endpoint, method = 'GET', data = null) {
        const options = {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include' // Include cookies
        };
        
        if (data && method !== 'GET') {
            options.body = JSON.stringify(data);
        }
        
        const response = await fetch(endpoint, options);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        return await response.json();
    }

    getCookie(name) {
        const value = `; ${document.cookie}`;
        const parts = value.split(`; ${name}=`);
        if (parts.length === 2) return parts.pop().split(';').shift();
        return null;
    }

    debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }

    async loadUserProfile() {
        try {
            const response = await this.apiCall('/api/admin/profile');
            if (response.success) {
                this.adminUser = response.data;
            }
        } catch (error) {
            console.error('Failed to load user profile:', error);
        }
    }
}

// Initialize admin dashboard when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.adminDashboard = new AdminDashboard();
});
