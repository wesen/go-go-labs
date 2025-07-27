// VacuumMart Main JavaScript
document.addEventListener('DOMContentLoaded', function() {
    // Initialize all interactive features
    initializeCart();
    initializeProductCards();
    initializeSearch();
    initializeFilters();
    initializeWishlist();
    initializeQuickView();
});

// Cart functionality
function initializeCart() {
    const cartButtons = document.querySelectorAll('.vm-btn-cart');
    const cartCount = document.querySelector('.badge');
    let cartItems = JSON.parse(localStorage.getItem('cartItems') || '[]');
    
    // Update cart count on page load
    updateCartCount();
    
    cartButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            e.preventDefault();
            
            // Get product data from the card
            const productCard = this.closest('.vm-product-card') || this.closest('.card');
            const productData = {
                id: Date.now(), // Simple ID generation
                name: productCard.querySelector('.card-title').textContent,
                price: productCard.querySelector('.vm-product-price').textContent,
                image: productCard.querySelector('img').src,
                quantity: 1
            };
            
            // Add animation effect
            this.classList.add('added');
            this.innerHTML = '<i class="fas fa-check"></i> Added!';
            
            // Add to cart
            cartItems.push(productData);
            localStorage.setItem('cartItems', JSON.stringify(cartItems));
            updateCartCount();
            
            // Reset button after animation
            setTimeout(() => {
                this.classList.remove('added');
                this.innerHTML = 'Add to Cart';
            }, 2000);
            
            // Show toast notification
            showToast('Product added to cart!', 'success');
        });
    });
    
    function updateCartCount() {
        if (cartCount) {
            cartCount.textContent = cartItems.length;
        }
    }
}

// Product card interactions
function initializeProductCards() {
    const productCards = document.querySelectorAll('.vm-product-card');
    
    productCards.forEach(card => {
        // Quick view on image click
        const image = card.querySelector('.card-img-top');
        if (image) {
            image.addEventListener('click', function() {
                const productName = card.querySelector('.card-title').textContent;
                showQuickView(productName);
            });
        }
        
        // Wishlist toggle
        const wishlistBtn = card.querySelector('.fa-heart').closest('button');
        if (wishlistBtn) {
            wishlistBtn.addEventListener('click', function(e) {
                e.preventDefault();
                toggleWishlist(this);
            });
        }
    });
}

// Search functionality
function initializeSearch() {
    const searchInput = document.querySelector('.vm-search-input');
    const searchBtn = document.querySelector('.vm-search-btn');
    
    if (searchInput) {
        // Search on enter key
        searchInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                performSearch(this.value);
            }
        });
        
        // Auto-complete suggestions (mock)
        searchInput.addEventListener('input', function() {
            const query = this.value;
            if (query.length > 2) {
                showSearchSuggestions(query);
            }
        });
    }
    
    if (searchBtn) {
        searchBtn.addEventListener('click', function() {
            const query = searchInput.value;
            performSearch(query);
        });
    }
}

// Filter functionality
function initializeFilters() {
    const filterCheckboxes = document.querySelectorAll('.vm-filter-option input[type="checkbox"]');
    const sortSelect = document.querySelector('select.form-select');
    
    filterCheckboxes.forEach(checkbox => {
        checkbox.addEventListener('change', function() {
            applyFilters();
        });
    });
    
    if (sortSelect) {
        sortSelect.addEventListener('change', function() {
            applySorting(this.value);
        });
    }
}

// Wishlist functionality
function initializeWishlist() {
    let wishlist = JSON.parse(localStorage.getItem('wishlist') || '[]');
    
    // Update wishlist UI on page load
    updateWishlistUI();
    
    function updateWishlistUI() {
        const heartIcons = document.querySelectorAll('.fa-heart');
        heartIcons.forEach(icon => {
            const productCard = icon.closest('.vm-product-card') || icon.closest('.card');
            if (productCard) {
                const productName = productCard.querySelector('.card-title').textContent;
                if (wishlist.includes(productName)) {
                    icon.classList.remove('far');
                    icon.classList.add('fas', 'text-danger');
                }
            }
        });
    }
    
    window.toggleWishlist = function(button) {
        const icon = button.querySelector('.fa-heart');
        const productCard = button.closest('.vm-product-card') || button.closest('.card');
        const productName = productCard.querySelector('.card-title').textContent;
        
        if (icon.classList.contains('far')) {
            // Add to wishlist
            icon.classList.remove('far');
            icon.classList.add('fas', 'text-danger');
            wishlist.push(productName);
            showToast('Added to wishlist!', 'success');
        } else {
            // Remove from wishlist
            icon.classList.remove('fas', 'text-danger');
            icon.classList.add('far');
            wishlist = wishlist.filter(item => item !== productName);
            showToast('Removed from wishlist', 'info');
        }
        
        localStorage.setItem('wishlist', JSON.stringify(wishlist));
    };
}

// Quick view modal
function initializeQuickView() {
    // Create quick view modal if it doesn't exist
    if (!document.getElementById('quickViewModal')) {
        const modalHTML = `
            <div class="modal fade vm-quick-view" id="quickViewModal" tabindex="-1">
                <div class="modal-dialog modal-lg">
                    <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title">Quick View</h5>
                            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                        </div>
                        <div class="modal-body">
                            <div class="row">
                                <div class="col-md-6">
                                    <img src="" class="img-fluid" id="quickViewImage" alt="Product Image">
                                </div>
                                <div class="col-md-6">
                                    <h3 id="quickViewTitle"></h3>
                                    <div class="vm-review-stars mb-3">
                                        ★★★★★ <small class="text-muted">(Reviews)</small>
                                    </div>
                                    <p id="quickViewDescription"></p>
                                    <div class="mb-3">
                                        <span class="vm-product-price" id="quickViewPrice"></span>
                                    </div>
                                    <button class="btn vm-btn-cart">Add to Cart</button>
                                    <button class="btn btn-outline-secondary ms-2">View Details</button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
        document.body.insertAdjacentHTML('beforeend', modalHTML);
    }
}

// Utility functions
function showQuickView(productName) {
    const modal = new bootstrap.Modal(document.getElementById('quickViewModal'));
    
    // Update modal content (mock data)
    document.getElementById('quickViewTitle').textContent = productName;
    document.getElementById('quickViewImage').src = 'https://via.placeholder.com/400x300';
    document.getElementById('quickViewPrice').textContent = '$299.99';
    document.getElementById('quickViewDescription').textContent = 'High-performance vacuum cleaner with advanced features and superior cleaning power.';
    
    modal.show();
}

function performSearch(query) {
    console.log('Searching for:', query);
    // In a real application, this would trigger a search API call
    showToast(`Searching for "${query}"...`, 'info');
}

function showSearchSuggestions(query) {
    // Mock search suggestions
    const suggestions = ['Dyson V15', 'Shark Navigator', 'Robot Vacuum', 'Cordless Vacuum'];
    const filtered = suggestions.filter(s => s.toLowerCase().includes(query.toLowerCase()));
    console.log('Search suggestions:', filtered);
}

function applyFilters() {
    const activeFilters = [];
    const checkboxes = document.querySelectorAll('.vm-filter-option input[type="checkbox"]:checked');
    
    checkboxes.forEach(checkbox => {
        activeFilters.push(checkbox.nextElementSibling.textContent);
    });
    
    console.log('Active filters:', activeFilters);
    showToast('Filters applied', 'info');
}

function applySorting(sortValue) {
    console.log('Sorting by:', sortValue);
    showToast(`Sorted by ${sortValue}`, 'info');
}

function showToast(message, type = 'info') {
    // Create toast container if it doesn't exist
    let toastContainer = document.getElementById('toastContainer');
    if (!toastContainer) {
        toastContainer = document.createElement('div');
        toastContainer.id = 'toastContainer';
        toastContainer.className = 'toast-container position-fixed top-0 end-0 p-3';
        toastContainer.style.zIndex = '9999';
        document.body.appendChild(toastContainer);
    }
    
    // Create toast
    const toastHTML = `
        <div class="toast align-items-center text-bg-${type} border-0" role="alert">
            <div class="d-flex">
                <div class="toast-body">${message}</div>
                <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
            </div>
        </div>
    `;
    
    toastContainer.insertAdjacentHTML('beforeend', toastHTML);
    
    // Show the toast
    const toastElement = toastContainer.lastElementChild;
    const toast = new bootstrap.Toast(toastElement, { delay: 3000 });
    toast.show();
    
    // Remove toast element after it's hidden
    toastElement.addEventListener('hidden.bs.toast', function() {
        this.remove();
    });
}

// Quantity controls for cart page
document.addEventListener('click', function(e) {
    if (e.target.matches('.btn-outline-secondary') && e.target.textContent === '-') {
        const input = e.target.nextElementSibling;
        if (input && input.type === 'number') {
            const currentValue = parseInt(input.value);
            if (currentValue > 1) {
                input.value = currentValue - 1;
                updateCartTotal();
            }
        }
    }
    
    if (e.target.matches('.btn-outline-secondary') && e.target.textContent === '+') {
        const input = e.target.previousElementSibling;
        if (input && input.type === 'number') {
            const currentValue = parseInt(input.value);
            input.value = currentValue + 1;
            updateCartTotal();
        }
    }
});

function updateCartTotal() {
    // Mock cart total update
    console.log('Updating cart total...');
}

// Smooth scrolling for anchor links
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            });
        }
    });
});

// Loading states for buttons
document.addEventListener('click', function(e) {
    if (e.target.classList.contains('vm-btn-cart') || 
        e.target.textContent.includes('Checkout') ||
        e.target.textContent.includes('Calculate')) {
        
        const originalText = e.target.innerHTML;
        e.target.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Loading...';
        e.target.disabled = true;
        
        setTimeout(() => {
            e.target.innerHTML = originalText;
            e.target.disabled = false;
        }, 1500);
    }
});
