// Product detail page functionality
class ProductDetailManager {
    constructor() {
        this.productId = this.extractProductIdFromURL();
        this.product = null;
        this.currentImageIndex = 0;
        
        this.loadProductDetail();
    }

    extractProductIdFromURL() {
        const path = window.location.pathname;
        const segments = path.split('/');
        const productSlug = segments[segments.length - 1];
        
        // For now, try to extract ID from slug or use mock data
        // In real implementation, we'd either use slug-based lookup or extract ID
        return 1; // Mock ID for now
    }

    async loadProductDetail() {
        const loadingSpinner = document.querySelector('#loading-spinner');
        const productContainer = document.querySelector('#product-detail-container');
        const errorContainer = document.querySelector('#error-container');
        
        try {
            if (loadingSpinner) loadingSpinner.style.display = 'block';
            
            const response = await vacuumAPI.getProduct(this.productId);
            this.product = response.data;
            
            this.renderProductDetail();
            this.updateBreadcrumb();
            
            if (loadingSpinner) loadingSpinner.style.display = 'none';
            if (productContainer) productContainer.style.display = 'block';
            
        } catch (error) {
            console.error('Failed to load product:', error);
            
            if (loadingSpinner) loadingSpinner.style.display = 'none';
            if (errorContainer) errorContainer.style.display = 'block';
        }
    }

    renderProductDetail() {
        const container = document.querySelector('#product-detail-container');
        if (!container || !this.product) return;

        const currentPrice = this.product.sale_price || this.product.base_price;
        const isOnSale = this.product.sale_price && this.product.sale_price < this.product.base_price;
        const savings = isOnSale ? this.product.base_price - this.product.sale_price : 0;

        container.innerHTML = `
            <div class="row">
                <div class="col-md-6">
                    ${this.renderImageGallery()}
                </div>
                
                <div class="col-md-6">
                    <div class="product-info">
                        <div class="mb-2">
                            <small class="text-muted">${this.product.brand_name || ''} • ${this.product.category_name || ''}</small>
                        </div>
                        
                        <h1 class="h3 mb-3">${this.product.name}</h1>
                        
                        <div class="mb-3">
                            ${this.product.avg_rating ? `
                                <div class="d-flex align-items-center mb-2">
                                    <div class="me-2">${formatRating(this.product.avg_rating)}</div>
                                    <span class="text-muted">${this.product.avg_rating.toFixed(1)} (${this.product.review_count} reviews)</span>
                                </div>
                            ` : ''}
                        </div>
                        
                        <div class="price-section mb-4">
                            <div class="d-flex align-items-center gap-3">
                                ${isOnSale ? `
                                    <span class="h4 text-danger mb-0">${formatPrice(currentPrice)}</span>
                                    <span class="text-decoration-line-through text-muted">${formatPrice(this.product.base_price)}</span>
                                    <span class="badge bg-danger">Save ${formatPrice(savings)}</span>
                                ` : `
                                    <span class="h4 mb-0">${formatPrice(currentPrice)}</span>
                                `}
                            </div>
                        </div>
                        
                        <div class="mb-4">
                            <p class="lead">${this.product.short_description}</p>
                        </div>
                        
                        <div class="stock-info mb-4">
                            ${this.product.inventory_qty > 0 ? `
                                <div class="alert alert-success">
                                    <i class="fas fa-check-circle me-2"></i>
                                    <strong>In Stock</strong> - ${this.product.inventory_qty} units available
                                </div>
                            ` : `
                                <div class="alert alert-danger">
                                    <i class="fas fa-exclamation-triangle me-2"></i>
                                    <strong>Out of Stock</strong>
                                </div>
                            `}
                        </div>
                        
                        <form class="add-to-cart-form mb-4">
                            <div class="row align-items-end">
                                <div class="col-auto">
                                    <label for="quantity" class="form-label">Quantity</label>
                                    <select class="form-select" id="quantity" name="quantity" style="width: 80px;">
                                        ${Array.from({length: Math.min(10, this.product.inventory_qty)}, (_, i) => 
                                            `<option value="${i + 1}">${i + 1}</option>`
                                        ).join('')}
                                    </select>
                                </div>
                                <div class="col">
                                    <button type="button" class="btn btn-primary btn-lg add-to-cart-btn" 
                                            data-product-id="${this.product.id}"
                                            ${this.product.inventory_qty <= 0 ? 'disabled' : ''}>
                                        <i class="fas fa-cart-plus me-2"></i>
                                        ${this.product.inventory_qty > 0 ? 'Add to Cart' : 'Out of Stock'}
                                    </button>
                                </div>
                            </div>
                        </form>
                        
                        <div class="product-features mb-4">
                            <h6>Key Features:</h6>
                            <ul class="list-unstyled">
                                <li><i class="fas fa-check text-success me-2"></i>${this.product.warranty_months} month warranty</li>
                                <li><i class="fas fa-check text-success me-2"></i>Free shipping on orders over $99</li>
                                <li><i class="fas fa-check text-success me-2"></i>30-day return policy</li>
                                <li><i class="fas fa-check text-success me-2"></i>Expert customer support</li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
            
            <div class="row mt-5">
                <div class="col-12">
                    <ul class="nav nav-tabs" id="productTabs" role="tablist">
                        <li class="nav-item" role="presentation">
                            <button class="nav-link active" id="description-tab" data-bs-toggle="tab" data-bs-target="#description" type="button">
                                Description
                            </button>
                        </li>
                        <li class="nav-item" role="presentation">
                            <button class="nav-link" id="specifications-tab" data-bs-toggle="tab" data-bs-target="#specifications" type="button">
                                Specifications
                            </button>
                        </li>
                        <li class="nav-item" role="presentation">
                            <button class="nav-link" id="reviews-tab" data-bs-toggle="tab" data-bs-target="#reviews" type="button">
                                Reviews (${this.product.review_count || 0})
                            </button>
                        </li>
                    </ul>
                    
                    <div class="tab-content mt-3" id="productTabContent">
                        <div class="tab-pane fade show active" id="description">
                            <div class="p-3">
                                <p>${this.product.description}</p>
                            </div>
                        </div>
                        
                        <div class="tab-pane fade" id="specifications">
                            <div class="p-3">
                                ${this.renderSpecifications()}
                            </div>
                        </div>
                        
                        <div class="tab-pane fade" id="reviews">
                            <div class="p-3">
                                <p class="text-muted">Reviews section coming soon.</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
    }

    renderImageGallery() {
        const images = this.product.images || [];
        if (images.length === 0) {
            return `
                <div class="product-image-gallery">
                    <img src="/static/images/placeholder-vacuum.jpg" 
                         alt="${this.product.name}" 
                         class="img-fluid main-product-image">
                </div>
            `;
        }

        const mainImage = images[this.currentImageIndex] || images[0];
        
        return `
            <div class="product-image-gallery">
                <div class="main-image-container mb-3">
                    <img src="${mainImage.url}" 
                         alt="${mainImage.alt_text || this.product.name}" 
                         class="img-fluid main-product-image"
                         id="main-product-image">
                </div>
                
                ${images.length > 1 ? `
                    <div class="thumbnail-container">
                        <div class="row g-2">
                            ${images.map((image, index) => `
                                <div class="col-3">
                                    <img src="${image.url}" 
                                         alt="${image.alt_text || this.product.name}"
                                         class="img-fluid thumbnail-image ${index === this.currentImageIndex ? 'active' : ''}"
                                         data-index="${index}"
                                         style="cursor: pointer; border: ${index === this.currentImageIndex ? '2px solid #007bff' : '1px solid #ddd'};">
                                </div>
                            `).join('')}
                        </div>
                    </div>
                ` : ''}
            </div>
        `;
    }

    renderSpecifications() {
        if (!this.product.attributes) {
            return '<p class="text-muted">No specifications available.</p>';
        }

        const specs = Object.entries(this.product.attributes).map(([key, value]) => {
            const label = key.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase());
            return `
                <div class="row mb-2">
                    <div class="col-4 fw-bold">${label}:</div>
                    <div class="col-8">${value}</div>
                </div>
            `;
        }).join('');

        return `<div class="specifications-table">${specs}</div>`;
    }

    updateBreadcrumb() {
        const breadcrumb = document.querySelector('#product-breadcrumb');
        if (breadcrumb && this.product) {
            breadcrumb.textContent = this.product.name;
        }
    }

    initializeImageGallery() {
        // Add click handlers for thumbnail images
        document.addEventListener('click', (e) => {
            if (e.target.matches('.thumbnail-image')) {
                const index = parseInt(e.target.dataset.index);
                this.switchImage(index);
            }
        });
    }

    switchImage(index) {
        if (!this.product.images || index >= this.product.images.length) return;
        
        this.currentImageIndex = index;
        
        const mainImage = document.querySelector('#main-product-image');
        const thumbnails = document.querySelectorAll('.thumbnail-image');
        
        if (mainImage) {
            const newImage = this.product.images[index];
            mainImage.src = newImage.url;
            mainImage.alt = newImage.alt_text || this.product.name;
        }
        
        thumbnails.forEach((thumb, i) => {
            thumb.style.border = i === index ? '2px solid #007bff' : '1px solid #ddd';
            thumb.classList.toggle('active', i === index);
        });
    }
}

// Initialize product detail manager
document.addEventListener('DOMContentLoaded', () => {
    if (window.location.pathname.startsWith('/products/') && window.location.pathname !== '/products') {
        window.productDetailManager = new ProductDetailManager();
    }
});
