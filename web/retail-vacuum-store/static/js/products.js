// Products page functionality
class ProductsPage {
    constructor() {
        this.currentFilters = {
            category: [],
            price: [],
            brand: [],
            feature: [],
            search: '',
            sortBy: 'name',
            page: 1,
            pageSize: 12
        };
        this.currentView = 'grid';
        this.products = [];
        this.totalPages = 0;
        
        this.init();
    }
    
    init() {
        this.bindEvents();
        this.loadProducts();
        this.loadCartCount();
    }
    
    bindEvents() {
        // Search functionality
        const searchInput = document.getElementById('searchInput');
        const searchBtn = document.getElementById('searchBtn');
        
        searchBtn.addEventListener('click', () => this.handleSearch());
        searchInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.handleSearch();
        });
        
        // Filter checkboxes
        document.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
            checkbox.addEventListener('change', () => this.handleFilterChange());
        });
        
        // Clear filters
        document.getElementById('clearFilters').addEventListener('click', () => this.clearAllFilters());
        
        // Sort dropdown
        document.getElementById('sortBy').addEventListener('change', (e) => {
            this.currentFilters.sortBy = e.target.value;
            this.currentFilters.page = 1;
            this.loadProducts();
        });
        
        // View toggle
        document.getElementById('gridView').addEventListener('click', () => this.switchView('grid'));
        document.getElementById('listView').addEventListener('click', () => this.switchView('list'));
    }
    
    handleSearch() {
        const searchInput = document.getElementById('searchInput');
        this.currentFilters.search = searchInput.value.trim();
        this.currentFilters.page = 1;
        this.loadProducts();
    }
    
    handleFilterChange() {
        // Collect all selected filters
        this.currentFilters.category = this.getSelectedValues('category');
        this.currentFilters.price = this.getSelectedValues('price');
        this.currentFilters.brand = this.getSelectedValues('brand');
        this.currentFilters.feature = this.getSelectedValues('feature');
        this.currentFilters.page = 1;
        
        this.loadProducts();
    }
    
    getSelectedValues(filterName) {
        return Array.from(document.querySelectorAll(`input[name="${filterName}"]:checked`))
            .map(input => input.value);
    }
    
    clearAllFilters() {
        // Uncheck all filter checkboxes
        document.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
            checkbox.checked = false;
        });
        
        // Reset search
        document.getElementById('searchInput').value = '';
        
        // Reset filters
        this.currentFilters = {
            category: [],
            price: [],
            brand: [],
            feature: [],
            search: '',
            sortBy: 'name',
            page: 1,
            pageSize: 12
        };
        
        this.loadProducts();
    }
    
    switchView(view) {
        this.currentView = view;
        
        // Update button states
        document.getElementById('gridView').classList.toggle('active', view === 'grid');
        document.getElementById('listView').classList.toggle('active', view === 'list');
        
        // Toggle visibility
        document.getElementById('productsGrid').style.display = view === 'grid' ? 'block' : 'none';
        document.getElementById('productsList').classList.toggle('d-none', view !== 'list');
        
        this.renderProducts();
    }
    
    async loadProducts() {
        document.getElementById('loadingState').style.display = 'block';
        document.getElementById('productsGrid').style.display = 'none';
        document.getElementById('productsList').classList.add('d-none');
        
        try {
            // Build query parameters
            const params = new URLSearchParams();
            
            if (this.currentFilters.search) {
                params.append('search', this.currentFilters.search);
            }
            
            this.currentFilters.category.forEach(cat => params.append('category', cat));
            this.currentFilters.brand.forEach(brand => params.append('brand', brand));
            this.currentFilters.feature.forEach(feature => params.append('feature', feature));
            this.currentFilters.price.forEach(price => params.append('price', price));
            
            params.append('sort', this.currentFilters.sortBy);
            params.append('page', this.currentFilters.page);
            params.append('pageSize', this.currentFilters.pageSize);
            
            // Mock API call - replace with actual API endpoint when backend is ready
            const response = await this.mockApiCall(params);
            
            this.products = response.products;
            this.totalPages = response.pages;
            
            this.renderProducts();
            this.renderPagination();
            this.updateProductCount(response.total);
            
        } catch (error) {
            console.error('Error loading products:', error);
            this.showError('Failed to load products. Please try again.');
        } finally {
            document.getElementById('loadingState').style.display = 'none';
        }
    }
    
    // Mock API response - replace with actual API integration
    async mockApiCall(params) {
        // Simulate API delay
        await new Promise(resolve => setTimeout(resolve, 500));
        
        // Mock products data
        const mockProducts = [
            {
                id: 1,
                name: "Dyson V15 Detect",
                slug: "dyson-v15-detect",
                basePrice: 749.99,
                salePrice: 649.99,
                category: "stick",
                brand: "dyson",
                shortDescription: "Powerful cordless vacuum with laser detection",
                images: [{ url: "https://via.placeholder.com/300x300?text=Dyson+V15", isPrimary: true }],
                attributes: ["cordless", "hepa", "smart"],
                inventory: { quantityAvailable: 15 },
                rating: 4.8,
                reviewCount: 248
            },
            {
                id: 2,
                name: "Shark Navigator Lift-Away",
                slug: "shark-navigator-lift-away",
                basePrice: 179.99,
                salePrice: null,
                category: "upright",
                brand: "shark",
                shortDescription: "Versatile upright vacuum with detachable canister",
                images: [{ url: "https://via.placeholder.com/300x300?text=Shark+Navigator", isPrimary: true }],
                attributes: ["pet"],
                inventory: { quantityAvailable: 8 },
                rating: 4.5,
                reviewCount: 156
            },
            {
                id: 3,
                name: "Roomba j7+",
                slug: "roomba-j7-plus",
                basePrice: 849.99,
                salePrice: 799.99,
                category: "robot",
                brand: "irobot",
                shortDescription: "Smart robot vacuum that avoids obstacles",
                images: [{ url: "https://via.placeholder.com/300x300?text=Roomba+j7+", isPrimary: true }],
                attributes: ["smart", "pet"],
                inventory: { quantityAvailable: 5 },
                rating: 4.6,
                reviewCount: 189
            },
            {
                id: 4,
                name: "Bissell Pet Hair Eraser",
                slug: "bissell-pet-hair-eraser",
                basePrice: 89.99,
                salePrice: null,
                category: "handheld",
                brand: "bissell",
                shortDescription: "Specialized handheld vacuum for pet hair",
                images: [{ url: "https://via.placeholder.com/300x300?text=Bissell+Pet", isPrimary: true }],
                attributes: ["pet", "cordless"],
                inventory: { quantityAvailable: 12 },
                rating: 4.3,
                reviewCount: 94
            },
            {
                id: 5,
                name: "Hoover Linx Cordless",
                slug: "hoover-linx-cordless",
                basePrice: 99.99,
                salePrice: 79.99,
                category: "stick",
                brand: "hoover",
                shortDescription: "Lightweight cordless stick vacuum",
                images: [{ url: "https://via.placeholder.com/300x300?text=Hoover+Linx", isPrimary: true }],
                attributes: ["cordless"],
                inventory: { quantityAvailable: 20 },
                rating: 4.1,
                reviewCount: 67
            },
            {
                id: 6,
                name: "Miele Complete C3",
                slug: "miele-complete-c3",
                basePrice: 399.99,
                salePrice: null,
                category: "canister",
                brand: "miele",
                shortDescription: "Premium canister vacuum with HEPA filtration",
                images: [{ url: "https://via.placeholder.com/300x300?text=Miele+C3", isPrimary: true }],
                attributes: ["hepa"],
                inventory: { quantityAvailable: 7 },
                rating: 4.7,
                reviewCount: 203
            }
        ];
        
        // Apply filters
        let filteredProducts = mockProducts.slice();
        
        // Search filter
        if (params.get('search')) {
            const searchTerm = params.get('search').toLowerCase();
            filteredProducts = filteredProducts.filter(product => 
                product.name.toLowerCase().includes(searchTerm) ||
                product.shortDescription.toLowerCase().includes(searchTerm)
            );
        }
        
        // Category filter
        const categories = params.getAll('category');
        if (categories.length > 0) {
            filteredProducts = filteredProducts.filter(product => 
                categories.includes(product.category)
            );
        }
        
        // Brand filter
        const brands = params.getAll('brand');
        if (brands.length > 0) {
            filteredProducts = filteredProducts.filter(product => 
                brands.includes(product.brand)
            );
        }
        
        // Feature filter
        const features = params.getAll('feature');
        if (features.length > 0) {
            filteredProducts = filteredProducts.filter(product => 
                features.some(feature => product.attributes.includes(feature))
            );
        }
        
        // Price filter
        const priceRanges = params.getAll('price');
        if (priceRanges.length > 0) {
            filteredProducts = filteredProducts.filter(product => {
                const price = product.salePrice || product.basePrice;
                return priceRanges.some(range => {
                    switch (range) {
                        case '0-100': return price < 100;
                        case '100-300': return price >= 100 && price < 300;
                        case '300-500': return price >= 300 && price < 500;
                        case '500+': return price >= 500;
                        default: return true;
                    }
                });
            });
        }
        
        // Sort
        const sortBy = params.get('sort') || 'name';
        filteredProducts.sort((a, b) => {
            switch (sortBy) {
                case 'price-low':
                    return (a.salePrice || a.basePrice) - (b.salePrice || b.basePrice);
                case 'price-high':
                    return (b.salePrice || b.basePrice) - (a.salePrice || a.basePrice);
                case 'rating':
                    return b.rating - a.rating;
                case 'newest':
                    return b.id - a.id;
                default:
                    return a.name.localeCompare(b.name);
            }
        });
        
        // Pagination
        const page = parseInt(params.get('page')) || 1;
        const pageSize = parseInt(params.get('pageSize')) || 12;
        const total = filteredProducts.length;
        const pages = Math.ceil(total / pageSize);
        const startIndex = (page - 1) * pageSize;
        const endIndex = startIndex + pageSize;
        
        return {
            products: filteredProducts.slice(startIndex, endIndex),
            total,
            page,
            pageSize,
            pages
        };
    }
    
    renderProducts() {
        if (this.currentView === 'grid') {
            this.renderGridView();
        } else {
            this.renderListView();
        }
    }
    
    renderGridView() {
        const container = document.getElementById('productsGrid');
        
        if (this.products.length === 0) {
            container.innerHTML = `
                <div class="col-12 text-center py-5">
                    <i class="fas fa-search fa-3x text-muted mb-3"></i>
                    <h3>No products found</h3>
                    <p class="text-muted">Try adjusting your filters or search terms.</p>
                </div>
            `;
            container.style.display = 'block';
            return;
        }
        
        container.innerHTML = this.products.map(product => this.createProductCard(product)).join('');
        container.style.display = 'block';
    }
    
    renderListView() {
        const container = document.getElementById('productsList');
        
        if (this.products.length === 0) {
            container.innerHTML = `
                <div class="text-center py-5">
                    <i class="fas fa-search fa-3x text-muted mb-3"></i>
                    <h3>No products found</h3>
                    <p class="text-muted">Try adjusting your filters or search terms.</p>
                </div>
            `;
            container.classList.remove('d-none');
            return;
        }
        
        container.innerHTML = this.products.map(product => this.createProductListItem(product)).join('');
        container.classList.remove('d-none');
    }
    
    createProductCard(product) {
        const currentPrice = product.salePrice || product.basePrice;
        const isOnSale = product.salePrice && product.salePrice < product.basePrice;
        const primaryImage = product.images[0];
        const inStock = product.inventory.quantityAvailable > 0;
        
        return `
            <div class="col-lg-4 col-md-6 mb-4">
                <div class="vm-product-card card h-100">
                    <div class="position-relative">
                        <img src="${primaryImage.url}" class="card-img-top" alt="${product.name}" style="height: 200px; object-fit: cover;">
                        ${product.attributes.includes('smart') ? '<span class="vm-product-badge bg-primary text-white">Smart Home</span>' : ''}
                        ${product.attributes.includes('pet') ? '<span class="vm-product-badge bg-success text-white">Pet Specialist</span>' : ''}
                        ${isOnSale ? '<span class="vm-product-badge bg-danger text-white">Sale</span>' : ''}
                        <button class="btn btn-sm btn-light position-absolute wishlist-btn" style="top: 10px; right: 10px;" data-product-id="${product.id}">
                            <i class="far fa-heart"></i>
                        </button>
                    </div>
                    <div class="card-body d-flex flex-column">
                        <h5 class="card-title">
                            <a href="/products/${product.slug}" class="text-decoration-none">${product.name}</a>
                        </h5>
                        <div class="vm-review-stars mb-2">
                            ${this.renderStars(product.rating)} 
                            <small class="text-muted">(${product.reviewCount} reviews)</small>
                        </div>
                        <p class="card-text flex-grow-1">${product.shortDescription}</p>
                        <div class="mt-auto">
                            <div class="d-flex justify-content-between align-items-center mb-2">
                                <div>
                                    <span class="vm-product-price fs-5 fw-bold">$${currentPrice.toFixed(2)}</span>
                                    ${isOnSale ? `<span class="vm-product-price-original text-decoration-line-through text-muted ms-2">$${product.basePrice.toFixed(2)}</span>` : ''}
                                </div>
                            </div>
                            ${inStock ? 
                                `<button class="btn vm-btn-cart w-100 add-to-cart-btn" data-product-id="${product.id}">Add to Cart</button>` :
                                `<button class="btn btn-secondary w-100" disabled>Out of Stock</button>`
                            }
                        </div>
                    </div>
                </div>
            </div>
        `;
    }
    
    createProductListItem(product) {
        const currentPrice = product.salePrice || product.basePrice;
        const isOnSale = product.salePrice && product.salePrice < product.basePrice;
        const primaryImage = product.images[0];
        const inStock = product.inventory.quantityAvailable > 0;
        
        return `
            <div class="card mb-3">
                <div class="row g-0">
                    <div class="col-md-3">
                        <div class="position-relative">
                            <img src="${primaryImage.url}" class="img-fluid rounded-start" alt="${product.name}" style="height: 200px; width: 100%; object-fit: cover;">
                            ${isOnSale ? '<span class="vm-product-badge bg-danger text-white">Sale</span>' : ''}
                        </div>
                    </div>
                    <div class="col-md-9">
                        <div class="card-body h-100 d-flex flex-column">
                            <div class="row">
                                <div class="col-md-8">
                                    <h5 class="card-title">
                                        <a href="/products/${product.slug}" class="text-decoration-none">${product.name}</a>
                                    </h5>
                                    <div class="vm-review-stars mb-2">
                                        ${this.renderStars(product.rating)} 
                                        <small class="text-muted">(${product.reviewCount} reviews)</small>
                                    </div>
                                    <p class="card-text">${product.shortDescription}</p>
                                    <div class="mb-2">
                                        ${product.attributes.map(attr => `<span class="badge bg-light text-dark me-1">${attr}</span>`).join('')}
                                    </div>
                                </div>
                                <div class="col-md-4 d-flex flex-column justify-content-between">
                                    <div class="text-end">
                                        <div class="mb-2">
                                            <span class="vm-product-price fs-4 fw-bold">$${currentPrice.toFixed(2)}</span>
                                            ${isOnSale ? `<div class="vm-product-price-original text-decoration-line-through text-muted">$${product.basePrice.toFixed(2)}</div>` : ''}
                                        </div>
                                        <p class="text-muted small mb-3">${inStock ? `${product.inventory.quantityAvailable} in stock` : 'Out of stock'}</p>
                                    </div>
                                    <div>
                                        ${inStock ? 
                                            `<button class="btn vm-btn-cart w-100 mb-2 add-to-cart-btn" data-product-id="${product.id}">Add to Cart</button>` :
                                            `<button class="btn btn-secondary w-100 mb-2" disabled>Out of Stock</button>`
                                        }
                                        <button class="btn btn-outline-secondary w-100 wishlist-btn" data-product-id="${product.id}">
                                            <i class="far fa-heart me-1"></i>Wishlist
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
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
    
    renderPagination() {
        const container = document.getElementById('pagination');
        
        if (this.totalPages <= 1) {
            container.innerHTML = '';
            return;
        }
        
        let pagination = '';
        
        // Previous button
        pagination += `
            <li class="page-item ${this.currentFilters.page === 1 ? 'disabled' : ''}">
                <a class="page-link" href="#" data-page="${this.currentFilters.page - 1}">Previous</a>
            </li>
        `;
        
        // Page numbers
        const startPage = Math.max(1, this.currentFilters.page - 2);
        const endPage = Math.min(this.totalPages, this.currentFilters.page + 2);
        
        for (let i = startPage; i <= endPage; i++) {
            pagination += `
                <li class="page-item ${i === this.currentFilters.page ? 'active' : ''}">
                    <a class="page-link" href="#" data-page="${i}">${i}</a>
                </li>
            `;
        }
        
        // Next button
        pagination += `
            <li class="page-item ${this.currentFilters.page === this.totalPages ? 'disabled' : ''}">
                <a class="page-link" href="#" data-page="${this.currentFilters.page + 1}">Next</a>
            </li>
        `;
        
        container.innerHTML = pagination;
        
        // Bind pagination click events
        container.querySelectorAll('a.page-link').forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                const page = parseInt(e.target.getAttribute('data-page'));
                if (page && page !== this.currentFilters.page) {
                    this.currentFilters.page = page;
                    this.loadProducts();
                }
            });
        });
    }
    
    updateProductCount(total) {
        const countElement = document.getElementById('productCount');
        countElement.textContent = `${total} products found`;
    }
    
    showError(message) {
        const container = document.getElementById('productsGrid');
        container.innerHTML = `
            <div class="col-12 text-center py-5">
                <i class="fas fa-exclamation-triangle fa-3x text-warning mb-3"></i>
                <h3>Error</h3>
                <p class="text-muted">${message}</p>
                <button class="btn btn-primary" onclick="location.reload()">Retry</button>
            </div>
        `;
        container.style.display = 'block';
    }
    
    async loadCartCount() {
        // Load cart count from localStorage or API
        const cart = JSON.parse(localStorage.getItem('cart') || '{"items": []}');
        const count = cart.items.reduce((total, item) => total + item.quantity, 0);
        document.getElementById('cartCount').textContent = count;
        
        if (count > 0) {
            this.updateMiniCart(cart);
        }
    }
    
    updateMiniCart(cart) {
        const miniCart = document.getElementById('miniCart');
        
        if (cart.items.length === 0) {
            miniCart.innerHTML = '<div class="p-3 text-center text-muted">Your cart is empty</div>';
            return;
        }
        
        const cartItems = cart.items.map(item => `
            <div class="vm-mini-cart-item border-bottom pb-2 mb-2">
                <div class="d-flex">
                    <img src="${item.image}" alt="${item.name}" style="width: 50px; height: 50px; object-fit: cover;" class="me-2">
                    <div class="flex-grow-1">
                        <div class="fw-bold">${item.name}</div>
                        <div class="text-muted small">Qty: ${item.quantity}</div>
                        <div class="text-primary fw-bold">$${(item.price * item.quantity).toFixed(2)}</div>
                    </div>
                </div>
            </div>
        `).join('');
        
        const total = cart.items.reduce((sum, item) => sum + (item.price * item.quantity), 0);
        
        miniCart.innerHTML = `
            ${cartItems}
            <div class="vm-mini-cart-total p-3">
                <div class="d-flex justify-content-between mb-2">
                    <span class="fw-bold">Total: $${total.toFixed(2)}</span>
                </div>
                <a href="/cart" class="btn btn-primary w-100">View Cart</a>
            </div>
        `;
    }
}

// Initialize the products page
document.addEventListener('DOMContentLoaded', () => {
    const productsPage = new ProductsPage();
    
    // Handle add to cart clicks
    document.addEventListener('click', (e) => {
        if (e.target.classList.contains('add-to-cart-btn') || e.target.closest('.add-to-cart-btn')) {
            const button = e.target.classList.contains('add-to-cart-btn') ? e.target : e.target.closest('.add-to-cart-btn');
            const productId = parseInt(button.getAttribute('data-product-id'));
            
            // Find the product
            const product = productsPage.products.find(p => p.id === productId);
            if (product) {
                addToCart(product);
                
                // Update button state
                button.innerHTML = '<i class="fas fa-check me-1"></i>Added!';
                button.classList.remove('vm-btn-cart');
                button.classList.add('btn-success');
                
                setTimeout(() => {
                    button.innerHTML = 'Add to Cart';
                    button.classList.remove('btn-success');
                    button.classList.add('vm-btn-cart');
                }, 2000);
            }
        }
        
        // Handle wishlist clicks
        if (e.target.classList.contains('wishlist-btn') || e.target.closest('.wishlist-btn')) {
            const button = e.target.classList.contains('wishlist-btn') ? e.target : e.target.closest('.wishlist-btn');
            const icon = button.querySelector('i');
            
            // Toggle wishlist state
            if (icon.classList.contains('far')) {
                icon.classList.remove('far');
                icon.classList.add('fas');
                button.classList.add('text-danger');
            } else {
                icon.classList.remove('fas');
                icon.classList.add('far');
                button.classList.remove('text-danger');
            }
        }
    });
});

// Cart functionality
function addToCart(product) {
    const cart = JSON.parse(localStorage.getItem('cart') || '{"items": []}');
    
    const existingItem = cart.items.find(item => item.id === product.id);
    
    if (existingItem) {
        existingItem.quantity += 1;
    } else {
        cart.items.push({
            id: product.id,
            name: product.name,
            price: product.salePrice || product.basePrice,
            image: product.images[0].url,
            quantity: 1
        });
    }
    
    localStorage.setItem('cart', JSON.stringify(cart));
    
    // Update cart count
    const count = cart.items.reduce((total, item) => total + item.quantity, 0);
    document.getElementById('cartCount').textContent = count;
    
    // Update mini cart
    if (window.productsPage) {
        window.productsPage.updateMiniCart(cart);
    }
    
    // Show success message
    showToast('Product added to cart!', 'success');
}

function showToast(message, type = 'info') {
    // Create toast element
    const toast = document.createElement('div');
    toast.className = `toast align-items-center text-white bg-${type === 'success' ? 'success' : 'primary'} border-0`;
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
