// Customer account page functionality
class AccountManager {
    constructor() {
        this.customer = null;
        this.orders = [];
        
        this.profileForm = document.querySelector('#profile-form');
        this.ordersContainer = document.querySelector('#orders-container');
        
        this.loadCustomerData();
        this.initializeEventListeners();
    }

    async loadCustomerData() {
        try {
            // Load profile and orders in parallel
            const [profileResponse, ordersResponse] = await Promise.all([
                vacuumAPI.getProfile(),
                vacuumAPI.getOrders()
            ]);
            
            this.customer = profileResponse.data;
            this.orders = ordersResponse.data;
            
            this.renderProfile();
            this.renderOrders();
            
        } catch (error) {
            console.error('Failed to load customer data:', error);
            this.showError('Failed to load account information');
        }
    }

    initializeEventListeners() {
        // Profile form submission
        if (this.profileForm) {
            this.profileForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.updateProfile();
            });
        }

        // Tab switching
        const tabLinks = document.querySelectorAll('[data-bs-toggle="tab"]');
        tabLinks.forEach(link => {
            link.addEventListener('shown.bs.tab', (e) => {
                const targetTab = e.target.getAttribute('href');
                if (targetTab === '#orders' && this.orders.length === 0) {
                    this.loadOrders();
                }
            });
        });
    }

    renderProfile() {
        if (!this.profileForm || !this.customer) return;

        // Populate form fields
        const fields = {
            'profile_first_name': this.customer.first_name || '',
            'profile_last_name': this.customer.last_name || '',
            'profile_email': this.customer.email || '',
            'profile_phone': this.customer.phone || ''
        };

        Object.entries(fields).forEach(([fieldId, value]) => {
            const field = document.querySelector(`#${fieldId}`);
            if (field) {
                field.value = value;
            }
        });
    }

    async updateProfile() {
        const formData = new FormData(this.profileForm);
        const profileData = {
            first_name: formData.get('first_name'),
            last_name: formData.get('last_name'),
            email: formData.get('email'),
            phone: formData.get('phone')
        };

        // Show loading state
        const submitBtn = this.profileForm.querySelector('button[type="submit"]');
        const originalText = submitBtn.innerHTML;
        submitBtn.disabled = true;
        submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin me-2"></i>Saving...';

        try {
            const response = await vacuumAPI.updateProfile(profileData);
            this.customer = response.data;
            
            this.showSuccess('Profile updated successfully');
            
        } catch (error) {
            console.error('Failed to update profile:', error);
            this.showError('Failed to update profile');
        } finally {
            // Restore button
            submitBtn.disabled = false;
            submitBtn.innerHTML = originalText;
        }
    }

    renderOrders() {
        if (!this.ordersContainer) return;

        if (!this.orders || this.orders.length === 0) {
            this.ordersContainer.innerHTML = `
                <div class="text-center py-5">
                    <i class="fas fa-shopping-bag fa-3x text-muted mb-3"></i>
                    <h5 class="text-muted">No orders yet</h5>
                    <p class="text-muted">When you place orders, they will appear here.</p>
                    <a href="/products" class="btn btn-primary">Start Shopping</a>
                </div>
            `;
            return;
        }

        const ordersHTML = this.orders.map(order => this.renderOrderCard(order)).join('');
        this.ordersContainer.innerHTML = ordersHTML;
    }

    renderOrderCard(order) {
        const statusBadge = this.getStatusBadge(order.status);
        const paymentBadge = this.getPaymentBadge(order.payment_status);

        return `
            <div class="card mb-3">
                <div class="card-header">
                    <div class="row align-items-center">
                        <div class="col-md-6">
                            <h6 class="mb-0">Order ${order.order_number}</h6>
                            <small class="text-muted">Placed on ${new Date(order.created_at).toLocaleDateString()}</small>
                        </div>
                        <div class="col-md-6 text-md-end">
                            ${statusBadge}
                            ${paymentBadge}
                        </div>
                    </div>
                </div>
                <div class="card-body">
                    <div class="row">
                        <div class="col-md-8">
                            ${order.items ? this.renderOrderItems(order.items) : ''}
                        </div>
                        <div class="col-md-4">
                            <div class="order-summary">
                                <div class="d-flex justify-content-between">
                                    <span>Total:</span>
                                    <span class="fw-bold">${formatPrice(order.total_amount)}</span>
                                </div>
                                ${order.shipping_method ? `
                                    <div class="d-flex justify-content-between">
                                        <span>Shipping:</span>
                                        <span>${order.shipping_method}</span>
                                    </div>
                                ` : ''}
                                ${order.tracking_number ? `
                                    <div class="d-flex justify-content-between">
                                        <span>Tracking:</span>
                                        <span class="font-monospace">${order.tracking_number}</span>
                                    </div>
                                ` : ''}
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card-footer bg-transparent">
                    <div class="d-flex gap-2">
                        <button class="btn btn-outline-primary btn-sm" onclick="window.print()">
                            <i class="fas fa-print me-1"></i>Print Receipt
                        </button>
                        ${order.status === 'delivered' ? `
                            <button class="btn btn-outline-secondary btn-sm">
                                <i class="fas fa-star me-1"></i>Leave Review
                            </button>
                        ` : ''}
                        ${order.status === 'shipped' && order.tracking_number ? `
                            <a href="#" class="btn btn-outline-info btn-sm">
                                <i class="fas fa-truck me-1"></i>Track Package
                            </a>
                        ` : ''}
                    </div>
                </div>
            </div>
        `;
    }

    renderOrderItems(items) {
        if (!items || items.length === 0) {
            return '<p class="text-muted">No items found</p>';
        }

        return items.slice(0, 3).map(item => `
            <div class="d-flex align-items-center mb-2">
                <div class="me-3">
                    <div class="fw-bold">${item.name}</div>
                    <small class="text-muted">Qty: ${item.quantity} × ${formatPrice(item.unit_price)}</small>
                </div>
            </div>
        `).join('') + 
        (items.length > 3 ? `<small class="text-muted">+${items.length - 3} more items</small>` : '');
    }

    getStatusBadge(status) {
        const statusConfig = {
            'pending': { class: 'warning', text: 'Pending' },
            'confirmed': { class: 'info', text: 'Confirmed' },
            'processing': { class: 'primary', text: 'Processing' },
            'shipped': { class: 'success', text: 'Shipped' },
            'delivered': { class: 'success', text: 'Delivered' },
            'cancelled': { class: 'danger', text: 'Cancelled' },
            'refunded': { class: 'secondary', text: 'Refunded' }
        };

        const config = statusConfig[status] || { class: 'secondary', text: status };
        return `<span class="badge bg-${config.class} me-2">${config.text}</span>`;
    }

    getPaymentBadge(paymentStatus) {
        const paymentConfig = {
            'pending': { class: 'warning', text: 'Payment Pending' },
            'paid': { class: 'success', text: 'Paid' },
            'failed': { class: 'danger', text: 'Payment Failed' },
            'refunded': { class: 'info', text: 'Refunded' },
            'partially_refunded': { class: 'info', text: 'Partially Refunded' }
        };

        const config = paymentConfig[paymentStatus] || { class: 'secondary', text: paymentStatus };
        return `<span class="badge bg-${config.class}">${config.text}</span>`;
    }

    async loadOrders() {
        try {
            const response = await vacuumAPI.getOrders();
            this.orders = response.data;
            this.renderOrders();
        } catch (error) {
            console.error('Failed to load orders:', error);
            this.showError('Failed to load order history');
        }
    }

    showSuccess(message) {
        const toast = document.createElement('div');
        toast.className = 'toast align-items-center text-white bg-success border-0';
        toast.setAttribute('role', 'alert');
        toast.innerHTML = `
            <div class="d-flex">
                <div class="toast-body">${message}</div>
                <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
            </div>
        `;
        
        document.body.appendChild(toast);
        const bsToast = new bootstrap.Toast(toast);
        bsToast.show();
        
        toast.addEventListener('hidden.bs.toast', () => {
            toast.remove();
        });
    }

    showError(message) {
        const toast = document.createElement('div');
        toast.className = 'toast align-items-center text-white bg-danger border-0';
        toast.setAttribute('role', 'alert');
        toast.innerHTML = `
            <div class="d-flex">
                <div class="toast-body">${message}</div>
                <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
            </div>
        `;
        
        document.body.appendChild(toast);
        const bsToast = new bootstrap.Toast(toast);
        bsToast.show();
        
        toast.addEventListener('hidden.bs.toast', () => {
            toast.remove();
        });
    }
}

// Initialize account manager
document.addEventListener('DOMContentLoaded', () => {
    if (window.location.pathname === '/account') {
        window.accountManager = new AccountManager();
    }
});
