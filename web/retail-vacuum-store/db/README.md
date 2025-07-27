# Vacuum Cleaner Retail Store - Database Design

## Overview

This database schema is designed for a comprehensive vacuum cleaner retail e-commerce platform, supporting product catalog management, inventory tracking, customer accounts, order processing, and administrative functions.

## Core Design Principles

### 1. Scalability
- Proper indexing on frequently queried columns
- Normalized design to reduce data redundancy
- Partitioning-ready structure for large tables (orders, inventory_movements)
- Efficient foreign key relationships

### 2. Data Integrity
- Comprehensive constraints and check conditions
- Foreign key relationships with appropriate CASCADE rules
- Generated columns for computed values (available inventory)
- Audit trails through timestamps and status histories

### 3. E-commerce Best Practices
- Flexible product attribute system for varying product specifications
- Snapshot data in orders (SKU, name, price) for historical accuracy
- Support for guest checkouts and registered customers
- Comprehensive pricing structure (base price, sale price, discounts)
- Multi-address support for customers
- Inventory reservation system

## Table Structure

### Product Management

#### `products`
Core product information including SKU, pricing, dimensions, and metadata.
- **Key Features**: Flexible product types (vacuum, accessory, part, bundle)
- **SEO Ready**: Slugs, meta titles, descriptions
- **Pricing**: Base price with optional sale pricing

#### `categories`
Hierarchical product categorization with self-referencing parent relationships.
- **Hierarchical**: Support for nested categories
- **SEO Friendly**: Slugs for URL generation
- **Flexible**: Display ordering and active/inactive states

#### `brands`
Manufacturer/brand information with logo and website links.

#### `product_attributes` & `product_attribute_values`
Flexible attribute system supporting various data types:
- Text, number, boolean, select, multiselect
- Units for numeric values (watts, pounds, etc.)
- Filterable attributes for search/browse functionality

#### `product_relationships`
Flexible product relationships supporting:
- Compatible parts and accessories
- Product bundles with quantities
- Upsell and cross-sell recommendations

#### `product_images`
Multiple images per product with ordering and primary image designation.

### Inventory Management

#### `inventory`
Real-time inventory tracking with:
- **Available Quantity**: Generated column (on_hand - reserved)
- **Multi-warehouse**: Support for different warehouse locations
- **Reorder Points**: Automatic reorder level tracking
- **Cost Tracking**: Per-unit cost for margin analysis

#### `inventory_movements`
Comprehensive audit trail of all inventory changes:
- Purchase receipts, sales, adjustments, transfers, returns
- Reference tracking to orders, purchase orders, etc.
- User attribution for accountability

### Customer Management

#### `customers`
Customer account information with:
- **Authentication**: Password hash storage
- **Preferences**: Marketing opt-in, email verification
- **Demographics**: Optional birth date for personalization

#### `customer_addresses`
Multiple addresses per customer:
- Separate billing and shipping addresses
- Default address designation
- Complete address validation fields

### Shopping & Orders

#### `carts` & `cart_items`
Shopping cart functionality supporting:
- **Guest Carts**: Session-based for non-registered users
- **Persistent Carts**: Customer-associated for registered users
- **Price Snapshots**: Current pricing captured in cart

#### `orders` & `order_items`
Comprehensive order management:
- **Address Snapshots**: Billing and shipping addresses stored with order
- **Pricing Breakdown**: Subtotal, tax, shipping, discounts, total
- **Status Tracking**: Full order lifecycle management
- **Guest Support**: Orders without customer accounts
- **Historical Data**: Product information snapshot at order time

#### `order_status_history`
Complete audit trail of order status changes with admin attribution.

### Promotions & Marketing

#### `coupons` & `coupon_usage`
Flexible discount system:
- **Discount Types**: Percentage, fixed amount, free shipping
- **Usage Limits**: Total and per-customer restrictions
- **Date Ranges**: Start and expiration dates
- **Usage Tracking**: Complete audit trail

### Reviews & Social Proof

#### `product_reviews`
Customer review system with:
- **Verified Purchases**: Reviews linked to actual orders
- **Moderation**: Approval system for review quality
- **Helpfulness**: Community voting on review quality
- **Rating System**: 1-5 star ratings

### Administration

#### `admin_users`
Administrative user management with:
- **Role-based Access**: Admin, manager, staff roles
- **Security**: Password hashing and session tracking
- **User Management**: Active/inactive states

## Key Features

### 1. Flexible Product Attributes
The attribute system allows for unlimited product specifications:
```sql
-- Example: Power rating for vacuum cleaners
INSERT INTO product_attributes (name, slug, attribute_type, unit) 
VALUES ('Power Rating', 'power-rating', 'number', 'watts');

-- Example: Filtration type options
INSERT INTO product_attributes (name, slug, attribute_type) 
VALUES ('Filtration Type', 'filtration-type', 'select');
INSERT INTO attribute_values (attribute_id, value) 
VALUES (2, 'HEPA'), (2, 'Cyclonic'), (2, 'Bag');
```

### 2. Inventory Reservation System
Prevents overselling by tracking reserved quantities:
```sql
-- Available quantity is automatically calculated
SELECT quantity_available FROM inventory WHERE product_id = 123;
```

### 3. Order Number Generation
Automatic order number generation with date prefix:
```sql
-- Generates: ORD-20241225-000001
```

### 4. Audit Trails
Comprehensive tracking of changes:
- Inventory movements with reference tracking
- Order status history with admin attribution
- Automatic timestamp updates

## Indexes and Performance

Key indexes for optimal performance:
- Product searches by category, brand, SKU
- Order lookups by customer, status, date
- Inventory queries by availability and reorder levels
- Customer email lookups for authentication

## Future Enhancements

### Phase 2 Considerations
1. **Multi-currency Support**: Currency conversion and regional pricing
2. **Advanced Inventory**: Serial number tracking, batch management
3. **Loyalty Program**: Points system and tier management
4. **Advanced Analytics**: Sales reporting, customer behavior tracking
5. **Multi-tenant**: Support for multiple store instances

### Scaling Considerations
1. **Partitioning**: Large tables (orders, inventory_movements) by date
2. **Read Replicas**: Separate read/write database instances
3. **Caching Layer**: Redis for frequently accessed data
4. **Search Engine**: Elasticsearch for advanced product search

## Security Considerations

1. **Password Security**: bcrypt hashing for all passwords
2. **Data Privacy**: Customer data encryption at rest
3. **Access Control**: Role-based permissions for admin users
4. **Audit Logging**: Complete activity tracking
5. **SQL Injection**: Prepared statements required for all queries

## Data Migration Strategy

1. **Version Control**: All schema changes in numbered migration files
2. **Rollback Support**: Down migrations for all changes
3. **Data Seeding**: Sample data for development and testing
4. **Production Safety**: Non-destructive migrations only

This schema provides a solid foundation for a scalable, secure vacuum cleaner retail platform while maintaining flexibility for future enhancements.
