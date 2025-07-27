// Vacuum Store JavaScript Application

class VacuumStore {
    constructor() {
        this.cart = JSON.parse(localStorage.getItem('vacuumCart')) || [];
        this.init();
    }

    init() {
        this.updateCartDisplay();
        this.bindEvents();
        console.log('Vacuum Store application initialized');
    }

    bindEvents() {
        // Add to cart buttons
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('add-to-cart')) {
                e.preventDefault();
                const productId = e.target.dataset.productId;
                const productName = e.target.dataset.productName;
                const productPrice = parseFloat(e.target.dataset.productPrice);
                
                this.addToCart(productId, productName, productPrice);
            }
        });

        // Remove from cart buttons
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('remove-from-cart')) {
                e.preventDefault();
                const productId = e.target.dataset.productId;
                this.removeFromCart(productId);
            }
        });

        // Clear cart button
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('clear-cart')) {
                e.preventDefault();
                this.clearCart();
            }
        });
    }

    addToCart(productId, productName, productPrice) {
        const existingItem = this.cart.find(item => item.id === productId);
        
        if (existingItem) {
            existingItem.quantity += 1;
        } else {
            this.cart.push({
                id: productId,
                name: productName,
                price: productPrice,
                quantity: 1
            });
        }

        this.saveCart();
        this.updateCartDisplay();
        this.showNotification(`${productName} added to cart!`, 'success');
    }

    removeFromCart(productId) {
        const index = this.cart.findIndex(item => item.id === productId);
        if (index > -1) {
            const item = this.cart[index];
            if (item.quantity > 1) {
                item.quantity -= 1;
            } else {
                this.cart.splice(index, 1);
            }
            
            this.saveCart();
            this.updateCartDisplay();
            this.showNotification('Item removed from cart', 'info');
        }
    }

    clearCart() {
        this.cart = [];
        this.saveCart();
        this.updateCartDisplay();
        this.showNotification('Cart cleared', 'info');
    }

    saveCart() {
        localStorage.setItem('vacuumCart', JSON.stringify(this.cart));
    }

    updateCartDisplay() {
        const cartCount = this.cart.reduce((total, item) => total + item.quantity, 0);
        const cartCountElements = document.querySelectorAll('.cart-count');
        
        cartCountElements.forEach(element => {
            element.textContent = cartCount;
            element.style.display = cartCount > 0 ? 'flex' : 'none';
        });

        // Update cart total
        const cartTotal = this.cart.reduce((total, item) => total + (item.price * item.quantity), 0);
        const cartTotalElements = document.querySelectorAll('.cart-total');
        
        cartTotalElements.forEach(element => {
            element.textContent = `$${cartTotal.toFixed(2)}`;
        });
    }

    showNotification(message, type = 'info') {
        // Create notification element
        const notification = document.createElement('div');
        notification.className = `alert alert-${type} alert-dismissible fade show position-fixed`;
        notification.style.cssText = 'top: 20px; right: 20px; z-index: 9999; min-width: 300px;';
        notification.innerHTML = `
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        `;

        document.body.appendChild(notification);

        // Auto-remove after 3 seconds
        setTimeout(() => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 3000);
    }

    // API methods for future backend integration
    async loadProducts() {
        try {
            const response = await fetch('/api/products');
            if (!response.ok) throw new Error('Failed to load products');
            return await response.json();
        } catch (error) {
            console.error('Error loading products:', error);
            this.showNotification('Failed to load products', 'danger');
            return [];
        }
    }

    async submitOrder(orderData) {
        try {
            const response = await fetch('/api/orders', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(orderData)
            });
            
            if (!response.ok) throw new Error('Failed to submit order');
            return await response.json();
        } catch (error) {
            console.error('Error submitting order:', error);
            this.showNotification('Failed to submit order', 'danger');
            throw error;
        }
    }
}

// Initialize the application when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    window.vacuumStore = new VacuumStore();
});

// Utility functions
function formatCurrency(amount) {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD'
    }).format(amount);
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
