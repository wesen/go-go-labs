// Shopping Cart functionality
class ShoppingCart {
    constructor() {
        this.cart = this.loadCart();
        this.promoCode = null;
        this.shippingCost = 0;
        this.taxRate = 0.08; // 8% tax rate
        this.freeShippingThreshold = 99;
        
        this.init();
    }
    
    init() {
        this.bindEvents();
        this.renderCart();
        this.loadRecommendedProducts();
    }
    
    bindEvents() {
        // Clear cart
        document.getElementById('clearCartBtn').addEventListener('click', () => this.clearCart());
        
        // Promo code
        document.getElementById('applyPromoBtn').addEventListener('click', () => this.applyPromoCode());
        document.getElementById('promoCode').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.applyPromoCode();
        });
        
        // Shipping calculator
        document.getElementById('calculateShippingBtn').addEventListener('click', () => this.calculateShipping());
        document.getElementById('zipCode').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.calculateShipping();
        });
        
        // Checkout
        document.getElementById('checkoutBtn').addEventListener('click', () => this.proceedToCheckout());
        
        // Quantity changes and item removal will be bound dynamically
    }
    
    loadCart() {
        return JSON.parse(localStorage.getItem('cart') || '{"items": [], "lastUpdated": null}');
    }
    
    saveCart() {
        this.cart.lastUpdated = new Date().toISOString();
        localStorage.setItem('cart', JSON.stringify(this.cart));
        this.updateCartCounts();
    }
    
    renderCart() {
        this.updateCartCounts();
        
        if (this.cart.items.length === 0) {
            this.showEmptyCart();
        } else {
            this.showCartItems();
            this.calculateTotals();
        }
    }
    
    showEmptyCart() {
        document.getElementById('emptyCartState').style.display = 'block';
        document.getElementById('cartItemsList').style.display = 'none';
        document.getElementById('recommendedSection').style.display = 'none';
        document.getElementById('checkoutBtn').disabled = true;
    }
    
    showCartItems() {
        document.getElementById('emptyCartState').style.display = 'none';
        document.getElementById('cartItemsList').style.display = 'block';
        document.getElementById('recommendedSection').style.display = 'block';
        document.getElementById('checkoutBtn').disabled = false;
        
        const container = document.getElementById('cartItemsList');
        const itemsHtml = this.cart.items.map((item, index) => this.createCartItemHtml(item, index)).join('');
        container.innerHTML = itemsHtml;
        
        // Bind item-specific events
        this.bindCartItemEvents();
    }
    
    createCartItemHtml(item, index) {
        const totalPrice = item.price * item.quantity;
        const originalPrice = item.originalPrice || item.price;
        const isOnSale = item.originalPrice && item.originalPrice > item.price;
        const savings = isOnSale ? (item.originalPrice - item.price) * item.quantity : 0;
        
        return `
            <div class="cart-item border-bottom pb-3 mb-3" data-index="${index}">
                <div class="row align-items-center">
                    <div class="col-md-2">
                        <img src="${item.image}" alt="${item.name}" class="img-fluid rounded" style="max-height: 100px;">
                    </div>
                    <div class="col-md-4">
                        <h6 class="mb-1">
                            <a href="/products/${item.slug || '#'}" class="text-decoration-none">${item.name}</a>
                        </h6>
                        <p class="text-muted small mb-1">SKU: ${item.sku || 'N/A'}</p>
                        ${item.attributes ? `<div class="small text-muted">${item.attributes.join(', ')}</div>` : ''}
                        <div class="price-info">
                            <span class="fw-bold text-primary">$${item.price.toFixed(2)}</span>
                            ${isOnSale ? `<span class="text-decoration-line-through text-muted ms-2">$${originalPrice.toFixed(2)}</span>` : ''}
                        </div>
                        ${savings > 0 ? `<div class="text-success small">Save $${savings.toFixed(2)}</div>` : ''}
                    </div>
                    <div class="col-md-3">
                        <label class="form-label small">Quantity:</label>
                        <div class="input-group input-group-sm">
                            <button class="btn btn-outline-secondary decrease-qty" type="button" data-index="${index}">-</button>
                            <input type="number" class="form-control text-center qty-input" value="${item.quantity}" min="1" max="10" data-index="${index}">
                            <button class="btn btn-outline-secondary increase-qty" type="button" data-index="${index}">+</button>
                        </div>
                        <div class="small text-muted mt-1">
                            ${item.stock ? (item.stock > 10 ? 'In Stock' : `Only ${item.stock} left`) : 'In Stock'}
                        </div>
                    </div>
                    <div class="col-md-2 text-end">
                        <div class="fw-bold fs-5">$${totalPrice.toFixed(2)}</div>
                        <button class="btn btn-outline-danger btn-sm mt-2 remove-item" data-index="${index}">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                    <div class="col-md-1 text-end">
                        <button class="btn btn-outline-secondary btn-sm save-later" data-index="${index}" title="Save for Later">
                            <i class="far fa-bookmark"></i>
                        </button>
                    </div>
                </div>
            </div>
        `;
    }
    
    bindCartItemEvents() {
        // Quantity decrease
        document.querySelectorAll('.decrease-qty').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const index = parseInt(e.target.getAttribute('data-index'));
                this.updateQuantity(index, -1);
            });
        });
        
        // Quantity increase
        document.querySelectorAll('.increase-qty').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const index = parseInt(e.target.getAttribute('data-index'));
                this.updateQuantity(index, 1);
            });
        });
        
        // Direct quantity input
        document.querySelectorAll('.qty-input').forEach(input => {
            input.addEventListener('change', (e) => {
                const index = parseInt(e.target.getAttribute('data-index'));
                const newQty = parseInt(e.target.value);
                this.setQuantity(index, newQty);
            });
        });
        
        // Remove item
        document.querySelectorAll('.remove-item').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const index = parseInt(e.target.getAttribute('data-index'));
                this.removeItem(index);
            });
        });
        
        // Save for later
        document.querySelectorAll('.save-later').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const index = parseInt(e.target.getAttribute('data-index'));
                this.saveForLater(index);
            });
        });
    }
    
    updateQuantity(index, change) {
        const item = this.cart.items[index];
        const newQty = Math.max(1, Math.min(item.quantity + change, item.stock || 10));
        this.setQuantity(index, newQty);
    }
    
    setQuantity(index, quantity) {
        if (quantity < 1) {
            this.removeItem(index);
            return;
        }
        
        this.cart.items[index].quantity = quantity;
        this.saveCart();
        this.renderCart();
        this.showToast('Cart updated', 'success');
    }
    
    removeItem(index) {
        const item = this.cart.items[index];
        
        if (confirm(`Remove "${item.name}" from your cart?`)) {
            this.cart.items.splice(index, 1);
            this.saveCart();
            this.renderCart();
            this.showToast('Item removed from cart', 'info');
        }
    }
    
    saveForLater(index) {
        const item = this.cart.items[index];
        
        // Get saved items from localStorage
        let savedItems = JSON.parse(localStorage.getItem('savedItems') || '[]');
        
        // Check if item already saved
        if (!savedItems.find(saved => saved.id === item.id)) {
            savedItems.push(item);
            localStorage.setItem('savedItems', JSON.stringify(savedItems));
        }
        
        // Remove from cart
        this.cart.items.splice(index, 1);
        this.saveCart();
        this.renderCart();
        this.showToast('Item saved for later', 'success');
    }
    
    clearCart() {
        if (this.cart.items.length === 0) return;
        
        if (confirm('Are you sure you want to clear your cart?')) {
            this.cart.items = [];
            this.saveCart();
            this.renderCart();
            this.showToast('Cart cleared', 'info');
        }
    }
    
    applyPromoCode() {
        const code = document.getElementById('promoCode').value.trim().toUpperCase();
        const messageDiv = document.getElementById('promoMessage');
        
        if (!code) {
            this.showPromoMessage('Please enter a promo code', 'error');
            return;
        }
        
        // Mock promo codes
        const promoCodes = {
            'SAVE10': { type: 'percentage', value: 10, description: '10% off your order' },
            'NEWCUSTOMER': { type: 'percentage', value: 15, description: '15% off for new customers' },
            'FREESHIP': { type: 'shipping', value: 0, description: 'Free shipping' },
            'SAVE25': { type: 'fixed', value: 25, description: '$25 off your order' }
        };
        
        if (promoCodes[code]) {
            this.promoCode = { code, ...promoCodes[code] };
            this.showPromoMessage(`Applied: ${this.promoCode.description}`, 'success');
            this.calculateTotals();
            
            // Disable the input and button
            document.getElementById('promoCode').disabled = true;
            document.getElementById('applyPromoBtn').textContent = 'Applied';
            document.getElementById('applyPromoBtn').disabled = true;
        } else {
            this.showPromoMessage('Invalid promo code', 'error');
        }
    }
    
    showPromoMessage(message, type) {
        const messageDiv = document.getElementById('promoMessage');
        messageDiv.textContent = message;
        messageDiv.className = `mt-2 small ${type === 'success' ? 'text-success' : 'text-danger'}`;
        messageDiv.style.display = 'block';
    }
    
    calculateShipping() {
        const zipCode = document.getElementById('zipCode').value.trim();
        
        if (!zipCode) {
            this.showToast('Please enter a zip code', 'error');
            return;
        }
        
        // Mock shipping calculation
        const subtotal = this.getSubtotal();
        
        if (subtotal >= this.freeShippingThreshold) {
            this.shippingCost = 0;
        } else {
            // Simple shipping calculation based on zip code
            const zone = this.getShippingZone(zipCode);
            this.shippingCost = zone.cost;
        }
        
        this.calculateTotals();
        this.showToast(`Shipping calculated for ${zipCode}`, 'success');
    }
    
    getShippingZone(zipCode) {
        // Mock shipping zones
        const firstDigit = parseInt(zipCode.charAt(0));
        
        if (firstDigit >= 0 && firstDigit <= 3) {
            return { name: 'East Coast', cost: 5.99 };
        } else if (firstDigit >= 4 && firstDigit <= 6) {
            return { name: 'Central', cost: 7.99 };
        } else {
            return { name: 'West Coast', cost: 9.99 };
        }
    }
    
    getSubtotal() {
        return this.cart.items.reduce((total, item) => total + (item.price * item.quantity), 0);
    }
    
    getDiscount() {
        if (!this.promoCode) return 0;
        
        const subtotal = this.getSubtotal();
        
        switch (this.promoCode.type) {
            case 'percentage':
                return subtotal * (this.promoCode.value / 100);
            case 'fixed':
                return Math.min(this.promoCode.value, subtotal);
            default:
                return 0;
        }
    }
    
    getTax() {
        const subtotal = this.getSubtotal();
        const discount = this.getDiscount();
        const taxableAmount = subtotal - discount;
        return taxableAmount * this.taxRate;
    }
    
    getShippingCost() {
        if (this.promoCode && this.promoCode.type === 'shipping') {
            return 0;
        }
        return this.shippingCost;
    }
    
    getTotal() {
        const subtotal = this.getSubtotal();
        const discount = this.getDiscount();
        const tax = this.getTax();
        const shipping = this.getShippingCost();
        
        return subtotal - discount + tax + shipping;
    }
    
    calculateTotals() {
        const subtotal = this.getSubtotal();
        const discount = this.getDiscount();
        const tax = this.getTax();
        const shipping = this.getShippingCost();
        const total = this.getTotal();
        
        // Update display
        document.getElementById('subtotal').textContent = `$${subtotal.toFixed(2)}`;
        document.getElementById('tax').textContent = `$${tax.toFixed(2)}`;
        document.getElementById('total').textContent = `$${total.toFixed(2)}`;
        
        // Shipping
        if (shipping === 0) {
            document.getElementById('shippingCost').textContent = 'FREE';
            document.getElementById('shippingCost').className = 'text-success fw-bold';
        } else {
            document.getElementById('shippingCost').textContent = `$${shipping.toFixed(2)}`;
            document.getElementById('shippingCost').className = '';
        }
        
        // Discount
        if (discount > 0) {
            document.getElementById('discountAmount').textContent = `-$${discount.toFixed(2)}`;
            document.getElementById('discountRow').style.display = 'flex';
        } else {
            document.getElementById('discountRow').style.display = 'none';
        }
        
        // Calculate total savings
        const totalSavings = this.cart.items.reduce((total, item) => {
            if (item.originalPrice && item.originalPrice > item.price) {
                return total + ((item.originalPrice - item.price) * item.quantity);
            }
            return total;
        }, 0) + discount;
        
        if (totalSavings > 0) {
            document.getElementById('totalSavings').textContent = totalSavings.toFixed(2);
            document.getElementById('savingsBanner').style.display = 'block';
        } else {
            document.getElementById('savingsBanner').style.display = 'none';
        }
        
        // Free shipping progress
        if (subtotal < this.freeShippingThreshold && shipping > 0) {
            const remaining = this.freeShippingThreshold - subtotal;
            this.showFreeShippingProgress(remaining);
        }
    }
    
    showFreeShippingProgress(remaining) {
        // Create or update free shipping progress message
        let progressDiv = document.getElementById('freeShippingProgress');
        if (!progressDiv) {
            progressDiv = document.createElement('div');
            progressDiv.id = 'freeShippingProgress';
            progressDiv.className = 'alert alert-info';
            document.querySelector('.order-totals').prepend(progressDiv);
        }
        
        progressDiv.innerHTML = `
            <small>
                <i class="fas fa-shipping-fast me-1"></i>
                Add $${remaining.toFixed(2)} more for <strong>FREE shipping</strong>!
            </small>
        `;
    }
    
    updateCartCounts() {
        const itemCount = this.cart.items.length;
        const totalQty = this.cart.items.reduce((total, item) => total + item.quantity, 0);
        
        document.getElementById('cartItemCount').textContent = itemCount;
        document.getElementById('navCartCount').textContent = totalQty;
        document.getElementById('summaryItemCount').textContent = totalQty;
    }
    
    async loadRecommendedProducts() {
        // Mock recommended products based on cart contents
        const recommendations = [
            {
                id: 101,
                name: "HEPA Filter Replacement Pack",
                price: 29.99,
                originalPrice: 39.99,
                image: "https://via.placeholder.com/200x200?text=HEPA+Filter",
                rating: 4.8,
                reviewCount: 156
            },
            {
                id: 102,
                name: "Vacuum Storage Bag",
                price: 19.99,
                image: "https://via.placeholder.com/200x200?text=Storage+Bag",
                rating: 4.5,
                reviewCount: 89
            },
            {
                id: 103,
                name: "Cleaning Solution Concentrate",
                price: 24.99,
                originalPrice: 29.99,
                image: "https://via.placeholder.com/200x200?text=Cleaning+Solution",
                rating: 4.6,
                reviewCount: 234
            }
        ];
        
        if (this.cart.items.length > 0) {
            this.renderRecommendedProducts(recommendations);
        }
    }
    
    renderRecommendedProducts(products) {
        const container = document.getElementById('recommendedProducts');
        
        const productsHtml = products.map(product => {
            const isOnSale = product.originalPrice && product.originalPrice > product.price;
            
            return `
                <div class="col-md-4 mb-3">
                    <div class="card h-100">
                        <div class="position-relative">
                            <img src="${product.image}" class="card-img-top" alt="${product.name}" style="height: 150px; object-fit: cover;">
                            ${isOnSale ? '<span class="vm-product-badge bg-danger text-white">Sale</span>' : ''}
                        </div>
                        <div class="card-body d-flex flex-column">
                            <h6 class="card-title">${product.name}</h6>
                            <div class="vm-review-stars small mb-2">
                                ${this.renderStars(product.rating)} (${product.reviewCount})
                            </div>
                            <div class="mt-auto">
                                <div class="mb-2">
                                    <span class="fw-bold text-primary">$${product.price.toFixed(2)}</span>
                                    ${isOnSale ? `<span class="text-decoration-line-through text-muted ms-1 small">$${product.originalPrice.toFixed(2)}</span>` : ''}
                                </div>
                                <button class="btn vm-btn-cart btn-sm w-100 add-recommended" data-product='${JSON.stringify(product)}'>
                                    Add to Cart
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            `;
        }).join('');
        
        container.innerHTML = productsHtml;
        
        // Bind add to cart events for recommended products
        container.querySelectorAll('.add-recommended').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const product = JSON.parse(e.target.getAttribute('data-product'));
                this.addRecommendedToCart(product, e.target);
            });
        });
    }
    
    addRecommendedToCart(product, button) {
        // Check if product already in cart
        const existingItem = this.cart.items.find(item => item.id === product.id);
        
        if (existingItem) {
            existingItem.quantity += 1;
        } else {
            this.cart.items.push({
                id: product.id,
                name: product.name,
                price: product.price,
                originalPrice: product.originalPrice,
                image: product.image,
                quantity: 1
            });
        }
        
        this.saveCart();
        this.renderCart();
        
        // Update button state
        button.innerHTML = '<i class="fas fa-check me-1"></i>Added!';
        button.classList.remove('vm-btn-cart');
        button.classList.add('btn-success');
        
        setTimeout(() => {
            button.innerHTML = 'Add to Cart';
            button.classList.remove('btn-success');
            button.classList.add('vm-btn-cart');
        }, 2000);
        
        this.showToast('Item added to cart!', 'success');
    }
    
    renderStars(rating) {
        const fullStars = Math.floor(rating);
        const hasHalfStar = rating % 1 >= 0.5;
        let stars = '';
        
        for (let i = 0; i < fullStars; i++) {
            stars += '★';
        }
        
        if (hasHalfStar) {
            stars += '☆';
        }
        
        for (let i = fullStars + (hasHalfStar ? 1 : 0); i < 5; i++) {
            stars += '☆';
        }
        
        return `<span class="text-warning">${stars}</span>`;
    }
    
    proceedToCheckout() {
        if (this.cart.items.length === 0) {
            this.showToast('Your cart is empty', 'error');
            return;
        }
        
        // Store cart summary for checkout
        const checkoutData = {
            items: this.cart.items,
            subtotal: this.getSubtotal(),
            discount: this.getDiscount(),
            tax: this.getTax(),
            shipping: this.getShippingCost(),
            total: this.getTotal(),
            promoCode: this.promoCode
        };
        
        localStorage.setItem('checkoutData', JSON.stringify(checkoutData));
        
        // Redirect to checkout
        window.location.href = '/checkout';
    }
    
    showToast(message, type = 'info') {
        // Create toast element
        const toast = document.createElement('div');
        toast.className = `toast align-items-center text-white bg-${type === 'success' ? 'success' : type === 'error' ? 'danger' : 'primary'} border-0`;
        toast.setAttribute('role', 'alert');
        toast.innerHTML = `
            <div class="d-flex">
                <div class="toast-body">${message}</div>
                <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
            </div>
        `;
        
        // Add to page
        let toastContainer = document.getElementById('toastContainer');
        if (!toastContainer) {
            toastContainer = document.createElement('div');
            toastContainer.id = 'toastContainer';
            toastContainer.className = 'toast-container position-fixed top-0 end-0 p-3';
            toastContainer.style.zIndex = '1055';
            document.body.appendChild(toastContainer);
        }
        
        toastContainer.appendChild(toast);
        
        // Show toast
        const bsToast = new bootstrap.Toast(toast);
        bsToast.show();
        
        // Remove from DOM after hidden
        toast.addEventListener('hidden.bs.toast', () => {
            toast.remove();
        });
    }
}

// Initialize the shopping cart
document.addEventListener('DOMContentLoaded', () => {
    new ShoppingCart();
});
