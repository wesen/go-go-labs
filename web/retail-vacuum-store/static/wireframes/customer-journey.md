# Customer Journey Wireframes

## Overview
This document outlines the complete customer journey from browsing to checkout completion, focusing on conversion optimization and user experience.

## 1. Browse Phase

### Homepage Entry Points
```
┌─────────────────────────────────────────┐
│ NAVIGATION BAR                          │
│ [Logo] [Categories] [Search] [Cart]     │
├─────────────────────────────────────────┤
│                                         │
│ HERO SECTION                           │
│ ┌─────────────┐ ┌─────────────────┐    │
│ │             │ │ Find Your       │    │
│ │ Hero Image  │ │ Perfect Vacuum  │    │
│ │             │ │ [Shop Now]      │    │
│ │             │ │ [Compare]       │    │
│ └─────────────┘ └─────────────────┘    │
│                                         │
├─────────────────────────────────────────┤
│ CATEGORY GRID                          │
│ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐       │
│ │Cat 1│ │Cat 2│ │Cat 3│ │Cat 4│       │
│ └─────┘ └─────┘ └─────┘ └─────┘       │
├─────────────────────────────────────────┤
│ FEATURED PRODUCTS                      │
│ ┌─────┐ ┌─────┐ ┌─────┐               │
│ │Prod1│ │Prod2│ │Prod3│               │
│ │★★★★★│ │★★★★☆│ │★★★★★│               │
│ │$699 │ │$179 │ │$599 │               │
│ │[Add]│ │[Add]│ │[Add]│               │
│ └─────┘ └─────┘ └─────┘               │
└─────────────────────────────────────────┘
```

### Product Listing Page
```
┌─────────────────────────────────────────┐
│ BREADCRUMB: Home > Categories > Upright │
├─────────────────────────────────────────┤
│ ┌───────────┐ ┌─────────────────────────┐│
│ │ FILTERS   │ │ PRODUCT GRID            ││
│ │           │ │ ┌─────┐ ┌─────┐ ┌─────┐ ││
│ │ Price     │ │ │Prod │ │Prod │ │Prod │ ││
│ │ □ <$100   │ │ │Image│ │Image│ │Image│ ││
│ │ ☑ $100-300│ │ │Title│ │Title│ │Title│ ││
│ │ □ $300+   │ │ │★★★★★│ │★★★★☆│ │★★★★★│ ││
│ │           │ │ │$299 │ │$179 │ │$149 │ ││
│ │ Brand     │ │ │[Add]│ │[Add]│ │[Add]│ ││
│ │ ☑ Dyson   │ │ └─────┘ └─────┘ └─────┘ ││
│ │ □ Shark   │ │                         ││
│ │ □ Bissell │ │ [Load More Products]    ││
│ │           │ │                         ││
│ │ Features  │ │                         ││
│ │ □ HEPA    │ │                         ││
│ │ □ Pet Hair│ │                         ││
│ └───────────┘ └─────────────────────────┘│
└─────────────────────────────────────────┘
```

## 2. Product Details Phase

### Product Detail Page
```
┌─────────────────────────────────────────┐
│ BREADCRUMB: Home > Upright > Dyson V15  │
├─────────────────────────────────────────┤
│ ┌─────────────────┐ ┌─────────────────┐ │
│ │ PRODUCT GALLERY │ │ PRODUCT INFO    │ │
│ │ ┌─────────────┐ │ │ Dyson V15 Detect│ │
│ │ │ Main Image  │ │ │ ★★★★★ (248)    │ │
│ │ │             │ │ │                 │ │
│ │ └─────────────┘ │ │ $699.99 $799.99 │ │
│ │ [◯][◯][◯][◯]   │ │                 │ │
│ │                 │ │ ┌─────────────┐ │ │
│ │ [360° View]     │ │ │Qty: [1] [▼] │ │ │
│ │                 │ │ └─────────────┘ │ │
│ │                 │ │                 │ │
│ │                 │ │ [ADD TO CART]   │ │
│ │                 │ │ [♡ WISHLIST]    │ │
│ │                 │ │                 │ │
│ │                 │ │ ✓ Free Shipping │ │
│ │                 │ │ ✓ 30-Day Return │ │
│ │                 │ │ ✓ 2-Year Warranty│ │
│ └─────────────────┘ └─────────────────┘ │
├─────────────────────────────────────────┤
│ TABS: [Description] [Specs] [Reviews]   │
│ ┌─────────────────────────────────────┐ │
│ │ Product description and features... │ │
│ └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│ RELATED PRODUCTS                        │
│ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐       │
│ │Rel 1│ │Rel 2│ │Rel 3│ │Rel 4│       │
│ └─────┘ └─────┘ └─────┘ └─────┘       │
└─────────────────────────────────────────┘
```

## 3. Cart Phase

### Shopping Cart Page
```
┌─────────────────────────────────────────┐
│ SHOPPING CART (3 items)                 │
├─────────────────────────────────────────┤
│ ┌───────────────────────┐ ┌───────────┐ │
│ │ CART ITEMS            │ │ SUMMARY   │ │
│ │ ┌───┐ Product Info    │ │           │ │
│ │ │IMG│ Dyson V15       │ │ Subtotal: │ │
│ │ └───┘ ★★★★★ $699.99   │ │ $929.97   │ │
│ │ [−][1][+] [Remove]    │ │           │ │
│ │                       │ │ Savings:  │ │
│ │ ┌───┐ Product Info    │ │ -$100.00  │ │
│ │ │IMG│ Shark Navigator │ │           │ │
│ │ └───┘ ★★★★☆ $179.99   │ │ Shipping: │ │
│ │ [−][1][+] [Remove]    │ │ FREE      │ │
│ │                       │ │           │ │
│ │ ┌───┐ Product Info    │ │ Tax:      │ │
│ │ │IMG│ HEPA Bags (2)   │ │ $66.40    │ │
│ │ └───┘ ★★★★★ $24.99 ea │ │           │ │
│ │ [−][2][+] [Remove]    │ │ Total:    │ │
│ │                       │ │ $896.37   │ │
│ │ ┌─────────────────┐   │ │           │ │
│ │ │ Promo Code: [__]│   │ │ [CHECKOUT]│ │
│ │ │ [Apply]         │   │ │           │ │
│ │ └─────────────────┘   │ │ Trust     │ │
│ │                       │ │ Signals   │ │
│ │ [Continue Shopping]   │ │ ✓ Secure  │ │
│ └───────────────────────┘ └───────────┘ │
└─────────────────────────────────────────┘
```

## 4. Checkout Phase

### Single Page Checkout
```
┌─────────────────────────────────────────┐
│ CHECKOUT - Secure & Fast                │
├─────────────────────────────────────────┤
│ Progress: [●]─[●]─[○]─[○]               │
│          Guest Shipping Payment Review   │
├─────────────────────────────────────────┤
│ ┌───────────────────────┐ ┌───────────┐ │
│ │ CHECKOUT FORM         │ │ ORDER     │ │
│ │                       │ │ SUMMARY   │ │
│ │ 1. CONTACT INFO       │ │           │ │
│ │ Email: [___________]  │ │ Dyson V15 │ │
│ │ □ Guest checkout      │ │ Qty: 1    │ │
│ │                       │ │ $699.99   │ │
│ │ 2. SHIPPING ADDRESS   │ │           │ │
│ │ First: [_____] Last:  │ │ Shark Nav │ │
│ │ [__________]          │ │ Qty: 1    │ │
│ │ Address: [__________] │ │ $179.99   │ │
│ │ City: [_____] State:  │ │           │ │
│ │ [__] ZIP: [_____]     │ │ HEPA Bags │ │
│ │                       │ │ Qty: 2    │ │
│ │ 3. SHIPPING METHOD    │ │ $49.98    │ │
│ │ ○ Standard (Free)     │ │           │ │
│ │ ○ Express ($9.99)     │ │ ───────── │ │
│ │ ○ Overnight ($19.99)  │ │ Subtotal: │ │
│ │                       │ │ $929.97   │ │
│ │ 4. PAYMENT METHOD     │ │ Tax:      │ │
│ │ ○ Credit Card         │ │ $66.40    │ │
│ │ ○ PayPal              │ │ ───────── │ │
│ │ ○ Apple Pay           │ │ Total:    │ │
│ │                       │ │ $996.37   │ │
│ │ Card: [____________]  │ │           │ │
│ │ Exp: [__/__] CVV:[_] │ │ [PLACE    │ │
│ │                       │ │  ORDER]   │ │
│ │ [◄ Back to Cart]      │ │           │ │
│ └───────────────────────┘ └───────────┘ │
└─────────────────────────────────────────┘
```

## 5. Conversion Optimization Elements

### Trust Signals Placement
- **Header**: SSL badges, customer service phone
- **Product Pages**: Reviews, ratings, warranties
- **Cart**: Security badges, return policy
- **Checkout**: Payment security, SSL indicators

### Urgency Elements
- **Product Cards**: "Only X left in stock"
- **Cart**: "Items in your cart are reserved for 15 minutes"
- **Checkout**: "X customers viewing this item"

### Social Proof
- **Homepage**: Customer testimonials
- **Product Pages**: Recent purchases, review highlights
- **Cart**: "Customers also bought" suggestions

## 6. Mobile Responsive Adaptations

### Mobile Navigation
```
┌─────────────────┐
│ [☰] VacuumMart │
│                 │
│ [Search______]  │
│ [🛒] [👤]       │
└─────────────────┘
```

### Mobile Product Grid
```
┌─────────────────┐
│ ┌─────┐ ┌─────┐ │
│ │Prod │ │Prod │ │
│ │     │ │     │ │
│ └─────┘ └─────┘ │
│ ┌─────┐ ┌─────┐ │
│ │Prod │ │Prod │ │
│ │     │ │     │ │
│ └─────┘ └─────┘ │
└─────────────────┘
```

### Mobile Checkout
```
┌─────────────────┐
│ Step 1 of 4     │
│ Contact Info    │
│                 │
│ Email:          │
│ [____________]  │
│                 │
│ [Continue]      │
│                 │
│ Order Summary   │
│ 3 items $896.37 │
│ [Show Details]  │
└─────────────────┘
```

## 7. Accessibility Considerations

### Screen Reader Support
- Proper heading hierarchy (h1-h6)
- Alt text for all product images
- ARIA labels for interactive elements
- Form field associations

### Keyboard Navigation
- Tab order through all interactive elements
- Skip navigation links
- Focus indicators for all focusable elements
- Keyboard shortcuts for common actions

### Color and Contrast
- WCAG AA compliant color ratios
- Color-blind friendly palette
- Information not conveyed by color alone
- High contrast mode support

This wireframe system ensures a seamless customer journey with multiple optimization touchpoints for maximum conversion rates.
