# Vacuum Cleaner Retail Store API Documentation

## Overview

This API provides backend services for a vacuum cleaner retail store, including product catalog management and shopping cart functionality.

## Base URL

```
http://localhost:3001/api/v1
```

## Authentication

Currently, the API uses session-based authentication with the following headers:
- `X-Session-ID`: Session identifier for guest users
- `X-Customer-ID`: Customer identifier for logged-in users (for demo purposes)

## Product Catalog API

### Get All Products

```http
GET /products
```

Query Parameters:
- `q`: Search term
- `category_id`: Filter by category ID
- `brand_id`: Filter by brand ID
- `product_type`: Filter by product type (vacuum, accessory, part, bundle)
- `min_price`: Minimum price filter
- `max_price`: Maximum price filter
- `featured`: Filter featured products (true/false)
- `in_stock`: Filter in-stock products (true/false)
- `sort_by`: Sort field (name, price, created_at, updated_at, featured)
- `sort_order`: Sort direction (asc/desc)
- `limit`: Number of results per page (default: 20)
- `offset`: Number of results to skip
- `page`: Page number (alternative to offset)

### Get Product by ID

```http
GET /products/{id}
```

### Get Product by SKU

```http
GET /products/sku/{sku}
```

### Create Product

```http
POST /products
```

Request Body:
```json
{
  "sku": "DYS-V15-001",
  "name": "Dyson V15 Detect",
  "slug": "dyson-v15-detect",
  "description": "Advanced cordless vacuum with laser dust detection",
  "short_description": "Cordless vacuum with laser detection",
  "category_id": 1,
  "brand_id": 1,
  "product_type": "vacuum",
  "base_price": 749.99,
  "sale_price": 699.99,
  "weight": 6.8,
  "warranty_months": 24,
  "is_active": true,
  "is_featured": true,
  "images": [
    {
      "url": "/static/images/dyson-v15.jpg",
      "alt_text": "Dyson V15 Detect",
      "is_primary": true,
      "display_order": 0
    }
  ],
  "initial_inventory": {
    "quantity_on_hand": 25,
    "reorder_level": 5,
    "reorder_quantity": 20
  }
}
```

### Update Product

```http
PUT /products/{id}
```

Request Body (partial update):
```json
{
  "name": "Updated Product Name",
  "base_price": 799.99,
  "is_featured": false
}
```

### Delete Product

```http
DELETE /products/{id}
```

### Get Categories

```http
GET /products/categories
```

### Get Brands

```http
GET /products/brands
```

### Get Featured Products

```http
GET /products/featured
```

### Get Products by Category

```http
GET /products/category/{categoryId}
```

### Get Products by Brand

```http
GET /products/brand/{brandId}
```

## Shopping Cart API

### Get Cart

```http
GET /cart
```

Headers:
- `X-Session-ID`: Session identifier

### Add Item to Cart

```http
POST /cart/items
```

Headers:
- `X-Session-ID`: Session identifier

Request Body:
```json
{
  "product_id": 1,
  "quantity": 2
}
```

### Update Cart Item

```http
PUT /cart/items/{productId}
```

Headers:
- `X-Session-ID`: Session identifier

Request Body:
```json
{
  "quantity": 3
}
```

### Remove Item from Cart

```http
DELETE /cart/items/{productId}
```

Headers:
- `X-Session-ID`: Session identifier

### Clear Cart

```http
DELETE /cart
```

Headers:
- `X-Session-ID`: Session identifier

### Get Cart Summary

```http
GET /cart/summary
```

Headers:
- `X-Session-ID`: Session identifier

### Reserve Cart Items

```http
POST /cart/reserve
```

Headers:
- `X-Session-ID`: Session identifier

### Release Cart Reservations

```http
POST /cart/release
```

Headers:
- `X-Session-ID`: Session identifier

### Merge Guest Cart

```http
POST /cart/merge
```

Headers:
- `X-Session-ID`: Session identifier
- `X-Customer-ID`: Customer identifier

## Data Models

### Product

```json
{
  "id": 1,
  "sku": "DYS-V15-001",
  "name": "Dyson V15 Detect",
  "slug": "dyson-v15-detect",
  "description": "Advanced cordless vacuum with laser dust detection",
  "short_description": "Cordless vacuum with laser detection",
  "category_id": 1,
  "brand_id": 1,
  "product_type": "vacuum",
  "base_price": 749.99,
  "sale_price": 699.99,
  "weight": 6.8,
  "dimensions_length": null,
  "dimensions_width": null,
  "dimensions_height": null,
  "warranty_months": 24,
  "is_active": true,
  "is_featured": true,
  "meta_title": null,
  "meta_description": null,
  "created_at": "2025-07-25T13:03:13.47738918-04:00",
  "updated_at": "2025-07-25T13:03:13.47738918-04:00",
  "category": {
    "id": 1,
    "name": "Upright Vacuums",
    "slug": "upright-vacuums",
    "description": "Traditional upright vacuum cleaners"
  },
  "brand": {
    "id": 1,
    "name": "Dyson",
    "slug": "dyson",
    "description": "Premium vacuum cleaner manufacturer"
  },
  "images": [
    {
      "id": 1,
      "product_id": 1,
      "url": "/static/images/dyson-v15.jpg",
      "alt_text": "Dyson V15 Detect",
      "display_order": 0,
      "is_primary": true
    }
  ],
  "inventory": {
    "id": 1,
    "product_id": 1,
    "quantity_on_hand": 25,
    "quantity_reserved": 5,
    "quantity_available": 20,
    "reorder_level": 5,
    "reorder_quantity": 20
  }
}
```

### Cart

```json
{
  "cart": {
    "id": 1,
    "customer_id": null,
    "session_id": "demo-session-123",
    "created_at": "2025-07-25T13:03:35.628842932-04:00",
    "updated_at": "2025-07-25T13:03:35.628842932-04:00",
    "items": [
      {
        "id": 1,
        "cart_id": 1,
        "product_id": 1,
        "quantity": 2,
        "unit_price": 100,
        "created_at": "2025-07-25T13:03:40.039027512-04:00",
        "updated_at": "2025-07-25T13:03:40.039027512-04:00",
        "product": {
          // Full product object
        }
      }
    ]
  },
  "summary": {
    "item_count": 2,
    "subtotal": 200,
    "tax_amount": 16,
    "shipping_cost": 0,
    "discount_amount": 0,
    "total": 216,
    "currency": "USD"
  }
}
```

## Error Responses

All endpoints return appropriate HTTP status codes and error messages:

```json
{
  "error": "Product not found"
}
```

Common status codes:
- `200`: Success
- `201`: Created
- `204`: No Content
- `400`: Bad Request
- `404`: Not Found
- `500`: Internal Server Error

## Demo Server

To run the demo server:

```bash
cd web/retail-vacuum-store
go run ./cmd/simple-demo
```

The server will start on port 3001 with mock data for testing all API endpoints.

### Sample API Calls

```bash
# Get all products
curl http://localhost:3001/api/v1/products

# Search for robot vacuums
curl "http://localhost:3001/api/v1/products?q=robot"

# Get featured products
curl http://localhost:3001/api/v1/products/featured

# Get cart (will create new cart)
curl -H "X-Session-ID: demo-session" http://localhost:3001/api/v1/cart

# Add product to cart
curl -X POST -H "Content-Type: application/json" -H "X-Session-ID: demo-session" \
  -d '{"product_id": 1, "quantity": 2}' \
  http://localhost:3001/api/v1/cart/items
```
