// API client for VacuumMart frontend
class VacuumMartAPI {
    constructor(baseURL = '/api') {
        this.baseURL = baseURL;
    }

    // Generic API request method
    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const config = {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            ...options
        };

        try {
            const response = await fetch(url, config);
            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || `HTTP ${response.status}`);
            }
            
            return data;
        } catch (error) {
            console.error('API request failed:', error);
            throw error;
        }
    }

    // Product API methods
    async getProducts(filters = {}) {
        const params = new URLSearchParams();
        Object.entries(filters).forEach(([key, value]) => {
            if (value !== null && value !== undefined && value !== '') {
                params.append(key, value);
            }
        });
        
        const endpoint = params.toString() ? `/products?${params}` : '/products';
        return this.request(endpoint);
    }

    async getProduct(id) {
        return this.request(`/products/${id}`);
    }

    async getCategories() {
        return this.request('/categories');
    }

    async getBrands() {
        return this.request('/brands');
    }

    // Cart API methods
    async getCart() {
        return this.request('/cart');
    }

    async addToCart(productId, quantity = 1) {
        return this.request('/cart', {
            method: 'POST',
            body: JSON.stringify({
                product_id: productId,
                quantity: quantity
            })
        });
    }

    async updateCartItem(itemId, quantity) {
        return this.request('/cart/items', {
            method: 'PUT',
            body: JSON.stringify({
                item_id: itemId,
                quantity: quantity
            })
        });
    }

    async removeCartItem(itemId) {
        return this.request(`/cart/items?item_id=${itemId}`, {
            method: 'DELETE'
        });
    }

    // Checkout API methods
    async checkout(orderData) {
        return this.request('/checkout', {
            method: 'POST',
            body: JSON.stringify(orderData)
        });
    }

    // Customer API methods
    async getProfile() {
        return this.request('/customers/profile');
    }

    async updateProfile(profileData) {
        return this.request('/customers/profile', {
            method: 'PUT',
            body: JSON.stringify(profileData)
        });
    }

    async getOrders() {
        return this.request('/orders');
    }
}

// Global API instance
window.vacuumAPI = new VacuumMartAPI();

// Cart management utilities
class CartManager {
    constructor() {
        this.cart = null;
        this.cartCountElement = document.querySelector('#cart-count');
        this.miniCartElement = document.querySelector('#mini-cart');
        this.loadCart();
    }

    async loadCart() {
        try {
            const response = await vacuumAPI.getCart();
            this.cart = response.data;
            this.updateCartUI();
        } catch (error) {
            console.error('Failed to load cart:', error);
        }
    }

    async addItem(productId, quantity = 1) {
        try {
            const response = await vacuumAPI.addToCart(productId, quantity);
            this.cart = response.data;
            this.updateCartUI();
            this.showAddToCartMessage(response.message);
            return true;
        } catch (error) {
            console.error('Failed to add item to cart:', error);
            this.showError('Failed to add item to cart');
            return false;
        }
    }

    async updateItem(itemId, quantity) {
        try {
            const response = await vacuumAPI.updateCartItem(itemId, quantity);
            this.cart = response.data;
            this.updateCartUI();
            return true;
        } catch (error) {
            console.error('Failed to update cart item:', error);
            this.showError('Failed to update cart item');
            return false;
        }
    }

    async removeItem(itemId) {
        try {
            const response = await vacuumAPI.removeCartItem(itemId);
            this.cart = response.data;
            this.updateCartUI();
            return true;
        } catch (error) {
            console.error('Failed to remove cart item:', error);
            this.showError('Failed to remove cart item');
            return false;
        }
    }

    updateCartUI() {
        // Update cart count badge
        if (this.cartCountElement) {
            const itemCount = this.cart ? this.cart.items.length : 0;
            this.cartCountElement.textContent = itemCount;
            this.cartCountElement.style.display = itemCount > 0 ? 'inline' : 'none';
        }

        // Update mini cart
        if (this.miniCartElement) {
            this.renderMiniCart();
        }
    }

    renderMiniCart() {
        if (!this.cart || !this.cart.items.length) {
            this.miniCartElement.innerHTML = '<p class="text-muted">Your cart is empty</p>';
            return;
        }

        const itemsHTML = this.cart.items.map(item => `
            <div class="mini-cart-item d-flex align-items-center mb-3">
                <img src="${item.image_url}" alt="${item.product_name}" class="me-3" style="width: 50px; height: 50px; object-fit: cover;">
                <div class="flex-grow-1">
                    <div class="fw-bold">${item.product_name}</div>
                    <div class="text-muted small">Qty: ${item.quantity} × $${item.unit_price}</div>
                </div>
                <div class="fw-bold">$${item.total_price.toFixed(2)}</div>
            </div>
        `).join('');

        this.miniCartElement.innerHTML = `
            ${itemsHTML}
            <hr>
            <div class="d-flex justify-content-between fw-bold">
                <span>Subtotal:</span>
                <span>$${this.cart.subtotal.toFixed(2)}</span>
            </div>
            <div class="d-grid gap-2 mt-3">
                <a href="/cart" class="btn btn-outline-primary">View Cart</a>
                <a href="/checkout" class="btn btn-primary">Checkout</a>
            </div>
        `;
    }

    showAddToCartMessage(message) {
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

// Initialize cart manager
document.addEventListener('DOMContentLoaded', () => {
    window.cartManager = new CartManager();
});

// Utility functions
function formatPrice(price) {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD'
    }).format(price);
}

function formatRating(rating) {
    const stars = Math.round(rating * 2) / 2; // Round to nearest 0.5
    let html = '';
    for (let i = 1; i <= 5; i++) {
        if (i <= stars) {
            html += '<i class="fas fa-star text-warning"></i>';
        } else if (i - 0.5 <= stars) {
            html += '<i class="fas fa-star-half-alt text-warning"></i>';
        } else {
            html += '<i class="far fa-star text-warning"></i>';
        }
    }
    return html;
}

function debounce(func, wait) {
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
