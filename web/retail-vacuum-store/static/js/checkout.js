// Checkout page functionality
class CheckoutManager {
    constructor() {
        this.cart = null;
        this.orderSummaryContainer = document.querySelector('#order-summary');
        this.checkoutForm = document.querySelector('#checkout-form');
        
        this.loadCart();
        this.initializeEventListeners();
    }

    async loadCart() {
        try {
            const response = await vacuumAPI.getCart();
            this.cart = response.data;
            this.renderOrderSummary();
            
            // Redirect to cart if empty
            if (!this.cart || !this.cart.items || this.cart.items.length === 0) {
                window.location.href = '/cart';
            }
        } catch (error) {
            console.error('Failed to load cart:', error);
            this.showError('Failed to load cart information');
        }
    }

    initializeEventListeners() {
        if (this.checkoutForm) {
            this.checkoutForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.processCheckout();
            });
        }

        // Form validation
        const formInputs = this.checkoutForm?.querySelectorAll('input[required]');
        formInputs?.forEach(input => {
            input.addEventListener('blur', () => {
                this.validateField(input);
            });
        });
    }

    renderOrderSummary() {
        if (!this.orderSummaryContainer || !this.cart) return;

        const subtotal = this.cart.items.reduce((sum, item) => sum + item.total_price, 0);
        const shipping = subtotal >= 99 ? 0 : 15.99;
        const tax = subtotal * 0.08;
        const total = subtotal + shipping + tax;

        const itemsHTML = this.cart.items.map(item => `
            <div class="d-flex justify-content-between align-items-center mb-2">
                <div class="d-flex align-items-center">
                    <img src="${item.image_url || '/static/images/placeholder-vacuum.jpg'}" 
                         alt="${item.product_name}" 
                         class="me-2"
                         style="width: 40px; height: 40px; object-fit: cover;">
                    <div>
                        <div class="fw-bold small">${item.product_name}</div>
                        <div class="text-muted small">Qty: ${item.quantity}</div>
                    </div>
                </div>
                <span class="fw-bold">${formatPrice(item.total_price)}</span>
            </div>
        `).join('');

        this.orderSummaryContainer.innerHTML = `
            <div class="order-items mb-3">
                ${itemsHTML}
            </div>
            
            <hr>
            
            <div class="order-totals">
                <div class="d-flex justify-content-between mb-2">
                    <span>Subtotal:</span>
                    <span>${formatPrice(subtotal)}</span>
                </div>
                <div class="d-flex justify-content-between mb-2">
                    <span>Shipping:</span>
                    <span>${shipping === 0 ? 'Free' : formatPrice(shipping)}</span>
                </div>
                <div class="d-flex justify-content-between mb-2">
                    <span>Tax:</span>
                    <span>${formatPrice(tax)}</span>
                </div>
                <hr>
                <div class="d-flex justify-content-between fw-bold h5">
                    <span>Total:</span>
                    <span>${formatPrice(total)}</span>
                </div>
            </div>
            
            <div class="security-notice mt-3">
                <small class="text-muted">
                    <i class="fas fa-lock me-1"></i>
                    Your order information is encrypted and secure
                </small>
            </div>
        `;
    }

    async processCheckout() {
        if (!this.validateForm()) {
            return;
        }

        const formData = new FormData(this.checkoutForm);
        const orderData = {
            email: formData.get('email'),
            shipping_address: {
                first_name: formData.get('first_name'),
                last_name: formData.get('last_name'),
                address_line1: formData.get('address_line1'),
                city: formData.get('city'),
                state: formData.get('state'),
                postal_code: formData.get('postal_code')
            },
            payment_method: formData.get('payment_method')
        };

        // Show loading state
        const submitBtn = this.checkoutForm.querySelector('button[type="submit"]');
        const originalText = submitBtn.innerHTML;
        submitBtn.disabled = true;
        submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin me-2"></i>Processing...';

        try {
            const response = await vacuumAPI.checkout(orderData);
            
            // Show success and redirect
            this.showOrderSuccess(response.data);
            
        } catch (error) {
            console.error('Checkout failed:', error);
            this.showError('Checkout failed. Please try again.');
            
            // Restore button
            submitBtn.disabled = false;
            submitBtn.innerHTML = originalText;
        }
    }

    validateForm() {
        const requiredFields = this.checkoutForm.querySelectorAll('input[required]');
        let isValid = true;

        requiredFields.forEach(field => {
            if (!this.validateField(field)) {
                isValid = false;
            }
        });

        return isValid;
    }

    validateField(field) {
        const value = field.value.trim();
        let isValid = true;
        let errorMessage = '';

        // Remove existing error
        this.clearFieldError(field);

        if (!value) {
            isValid = false;
            errorMessage = 'This field is required';
        } else if (field.type === 'email') {
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(value)) {
                isValid = false;
                errorMessage = 'Please enter a valid email address';
            }
        } else if (field.name === 'postal_code') {
            const zipRegex = /^\d{5}(-\d{4})?$/;
            if (!zipRegex.test(value)) {
                isValid = false;
                errorMessage = 'Please enter a valid ZIP code';
            }
        }

        if (!isValid) {
            this.showFieldError(field, errorMessage);
        }

        return isValid;
    }

    showFieldError(field, message) {
        field.classList.add('is-invalid');
        
        const errorDiv = document.createElement('div');
        errorDiv.className = 'invalid-feedback';
        errorDiv.textContent = message;
        
        field.parentNode.appendChild(errorDiv);
    }

    clearFieldError(field) {
        field.classList.remove('is-invalid');
        
        const errorDiv = field.parentNode.querySelector('.invalid-feedback');
        if (errorDiv) {
            errorDiv.remove();
        }
    }

    showOrderSuccess(order) {
        // Replace page content with success message
        document.body.innerHTML = `
            <div class="container mt-5">
                <div class="row justify-content-center">
                    <div class="col-md-8 text-center">
                        <div class="card">
                            <div class="card-body py-5">
                                <i class="fas fa-check-circle text-success" style="font-size: 4rem;"></i>
                                <h2 class="mt-3 mb-4">Order Placed Successfully!</h2>
                                <p class="lead">Thank you for your order. Your order number is:</p>
                                <h4 class="text-primary mb-4">${order.order_number}</h4>
                                <p>You will receive an email confirmation shortly.</p>
                                <div class="d-flex gap-3 justify-content-center mt-4">
                                    <a href="/" class="btn btn-primary">Continue Shopping</a>
                                    <a href="/account" class="btn btn-outline-primary">View Orders</a>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
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

// Initialize checkout manager
document.addEventListener('DOMContentLoaded', () => {
    if (window.location.pathname === '/checkout') {
        window.checkoutManager = new CheckoutManager();
    }
});
