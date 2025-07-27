-- ==============================================
-- Sample Data for Vacuum Cleaner Retail Store
-- ==============================================

-- Admin Users
INSERT INTO admin_users (username, email, password_hash, role, first_name, last_name) VALUES
('admin', 'admin@vacuumstore.com', '$2a$10$example_hash', 'admin', 'System', 'Administrator'),
('manager1', 'manager@vacuumstore.com', '$2a$10$example_hash', 'manager', 'Store', 'Manager'),
('staff1', 'staff@vacuumstore.com', '$2a$10$example_hash', 'staff', 'Sales', 'Associate');

-- Categories
INSERT INTO categories (name, slug, description, parent_id, display_order) VALUES
('Vacuum Cleaners', 'vacuum-cleaners', 'All types of vacuum cleaners', NULL, 1),
('Upright Vacuums', 'upright-vacuums', 'Traditional upright vacuum cleaners', 1, 1),
('Canister Vacuums', 'canister-vacuums', 'Portable canister-style vacuums', 1, 2),
('Robot Vacuums', 'robot-vacuums', 'Automated robotic vacuum cleaners', 1, 3),
('Handheld Vacuums', 'handheld-vacuums', 'Compact handheld cleaning devices', 1, 4),
('Accessories', 'accessories', 'Vacuum accessories and attachments', NULL, 2),
('Filters', 'filters', 'Replacement filters for all vacuum types', 6, 1),
('Bags', 'bags', 'Vacuum bags and disposable containers', 6, 2),
('Attachments', 'attachments', 'Cleaning attachments and tools', 6, 3),
('Parts', 'parts', 'Replacement parts and components', NULL, 3);

-- Brands
INSERT INTO brands (name, slug, description, website_url) VALUES
('Dyson', 'dyson', 'Premium cyclonic vacuum technology', 'https://www.dyson.com'),
('Shark', 'shark', 'Innovative cleaning solutions for every home', 'https://www.sharkclean.com'),
('Bissell', 'bissell', 'Deep cleaning solutions since 1876', 'https://www.bissell.com'),
('Hoover', 'hoover', 'Trusted cleaning solutions for over 100 years', 'https://www.hoover.com'),
('Miele', 'miele', 'German engineering for superior cleaning', 'https://www.miele.com'),
('Roomba', 'roomba', 'Revolutionary robotic cleaning technology', 'https://www.irobot.com'),
('Tineco', 'tineco', 'Smart cleaning appliances', 'https://www.tineco.com'),
('Oreck', 'oreck', 'Lightweight commercial-grade vacuums', 'https://www.oreck.com');

-- Product Attributes
INSERT INTO product_attributes (name, slug, attribute_type, unit, is_filterable) VALUES
('Power Rating', 'power-rating', 'number', 'watts', true),
('Weight', 'weight', 'number', 'lbs', true),
('Cord Length', 'cord-length', 'number', 'feet', true),
('Filtration Type', 'filtration-type', 'select', NULL, true),
('Bag Type', 'bag-type', 'select', NULL, true),
('Surface Type', 'surface-type', 'multiselect', NULL, true),
('Noise Level', 'noise-level', 'number', 'dB', true),
('Dustbin Capacity', 'dustbin-capacity', 'number', 'liters', true),
('Battery Life', 'battery-life', 'number', 'minutes', true),
('Warranty Period', 'warranty-period', 'number', 'years', false);

-- Attribute Values
INSERT INTO attribute_values (attribute_id, value, display_order) VALUES
-- Filtration Type
(4, 'HEPA', 1), (4, 'Cyclonic', 2), (4, 'Bag', 3), (4, 'Washable', 4),
-- Bag Type  
(5, 'Bagless', 1), (5, 'Type A', 2), (5, 'Type B', 3), (5, 'Universal', 4),
-- Surface Type
(6, 'Carpet', 1), (6, 'Hardwood', 2), (6, 'Tile', 3), (6, 'Pet Hair', 4), (6, 'Upholstery', 5);

-- Products
INSERT INTO products (sku, name, slug, description, short_description, category_id, brand_id, product_type, base_price, sale_price, weight, warranty_months, is_featured) VALUES
-- Dyson Products
('DYS-V15-001', 'Dyson V15 Detect Absolute', 'dyson-v15-detect-absolute', 
 'The most powerful, intelligent cordless vacuum. Laser reveals microscopic dust. LCD screen shows real-time proof of a deep clean.',
 'Powerful cordless with laser dust detection', 2, 1, 'vacuum', 749.99, 699.99, 6.8, 24, true),

('DYS-V8-001', 'Dyson V8 Animal', 'dyson-v8-animal',
 'Designed for homes with pets. Powerful fade-free suction. Transforms to a handheld for versatile cleaning.',
 'Cordless vacuum designed for pet owners', 2, 1, 'vacuum', 449.99, NULL, 5.6, 24, false),

-- Shark Products  
('SHK-NV356', 'Shark Navigator Lift-Away Professional', 'shark-navigator-professional',
 'Never loses suction with Anti-Allergen Complete Seal Technology and HEPA filter.',
 'Professional upright with lift-away canister', 2, 2, 'vacuum', 179.99, 149.99, 13.7, 60, true),

('SHK-RV1001AE', 'Shark IQ Robot Self-Empty XL', 'shark-iq-robot-self-empty',
 'Self-emptying robot vacuum with IQ Navigation and app control. No-sweep zones.',
 'Self-emptying robot with smart navigation', 4, 2, 'vacuum', 649.99, 549.99, 5.8, 12, true),

-- Bissell Products
('BIS-2252', 'Bissell CrossWave Pet Pro All-in-One', 'bissell-crosswave-pet-pro',
 'Vacuum and wash your floors at the same time. Safe for pets and reduces pet odors.',
 'Wet/dry vacuum for multi-surface cleaning', 2, 3, 'vacuum', 249.99, 199.99, 11.5, 24, false),

-- Roomba Products
('IRB-J7PLUS', 'iRobot Roomba j7+ Self-Emptying Robot', 'roomba-j7-plus-self-empty',
 'Avoids pet waste and gets stuck less. Empties itself for 60 days. Smart mapping.',
 'Self-emptying robot with pet waste avoidance', 4, 6, 'vacuum', 849.99, 799.99, 7.5, 12, true),

-- Accessories and Parts
('ACC-HEPA-UNI', 'Universal HEPA Filter Set', 'universal-hepa-filter-set',
 'High-efficiency particulate air filter compatible with most vacuum brands.',
 'HEPA filtration for cleaner air', 7, NULL, 'accessory', 29.99, NULL, 0.5, 6, false),

('ACC-PET-TOOL', 'Pet Hair Removal Tool', 'pet-hair-removal-tool',
 'Specialized attachment for removing embedded pet hair from upholstery and stairs.',
 'Removes stubborn pet hair effectively', 9, NULL, 'accessory', 19.99, 14.99, 0.3, 12, false),

('PART-DYS-WAND', 'Dyson Replacement Wand', 'dyson-replacement-wand',
 'Original replacement wand compatible with Dyson V8, V10, V11, and V15 models.',
 'OEM replacement wand for Dyson cordless', 10, 1, 'part', 89.99, NULL, 1.2, 12, false);

-- Product Attribute Values
INSERT INTO product_attribute_values (product_id, attribute_id, value) VALUES
-- Dyson V15 Detect
(1, 1, '230'), (1, 2, '6.8'), (1, 4, 'Cyclonic'), (1, 5, 'Bagless'), (1, 6, 'Carpet,Hardwood,Tile,Pet Hair'), (1, 9, '60'),
-- Dyson V8 Animal  
(2, 1, '115'), (2, 2, '5.6'), (2, 4, 'Cyclonic'), (2, 5, 'Bagless'), (2, 6, 'Carpet,Hardwood,Pet Hair'), (2, 9, '40'),
-- Shark Navigator
(3, 1, '1200'), (3, 2, '13.7'), (3, 3, '25'), (3, 4, 'HEPA'), (3, 5, 'Bagless'), (3, 6, 'Carpet,Hardwood,Tile'),
-- Shark Robot
(4, 2, '5.8'), (4, 4, 'HEPA'), (4, 5, 'Bagless'), (4, 6, 'Carpet,Hardwood,Tile,Pet Hair'), (4, 9, '90'),
-- Bissell CrossWave
(5, 1, '560'), (5, 2, '11.5'), (5, 3, '25'), (5, 4, 'Washable'), (5, 5, 'Bagless'), (5, 6, 'Hardwood,Tile'),
-- Roomba j7+
(6, 2, '7.5'), (6, 4, 'HEPA'), (6, 5, 'Bagless'), (6, 6, 'Carpet,Hardwood,Tile,Pet Hair'), (6, 9, '75');

-- Product Images (URLs would be actual image paths in production)
INSERT INTO product_images (product_id, url, alt_text, display_order, is_primary) VALUES
(1, '/static/images/products/dyson-v15-detect-1.jpg', 'Dyson V15 Detect front view', 1, true),
(1, '/static/images/products/dyson-v15-detect-2.jpg', 'Dyson V15 Detect accessories', 2, false),
(2, '/static/images/products/dyson-v8-animal-1.jpg', 'Dyson V8 Animal main unit', 1, true),
(3, '/static/images/products/shark-navigator-1.jpg', 'Shark Navigator Professional', 1, true),
(4, '/static/images/products/shark-robot-1.jpg', 'Shark IQ Robot vacuum', 1, true),
(5, '/static/images/products/bissell-crosswave-1.jpg', 'Bissell CrossWave Pet Pro', 1, true),
(6, '/static/images/products/roomba-j7-1.jpg', 'iRobot Roomba j7+ on dock', 1, true);

-- Product Relationships
INSERT INTO product_relationships (product_id, related_product_id, relationship_type, display_order) VALUES
-- Accessories for Dyson V15
(1, 7, 'accessory', 1), (1, 8, 'accessory', 2),
-- Parts for Dyson V8  
(2, 9, 'part', 1),
-- Cross-sell recommendations
(1, 2, 'cross_sell', 1), (2, 1, 'cross_sell', 1),
(3, 4, 'cross_sell', 1), (4, 3, 'cross_sell', 1),
-- Upsell recommendations
(2, 1, 'upsell', 1), (3, 4, 'upsell', 1);

-- Inventory
INSERT INTO inventory (product_id, warehouse_location, quantity_on_hand, quantity_reserved, reorder_level, reorder_quantity, cost_per_unit) VALUES
(1, 'Main Warehouse', 25, 3, 10, 20, 450.00),
(2, 'Main Warehouse', 40, 5, 15, 25, 270.00),
(3, 'Main Warehouse', 60, 8, 20, 30, 95.00),
(4, 'Main Warehouse', 15, 2, 8, 15, 390.00),
(5, 'Main Warehouse', 30, 4, 12, 20, 150.00),
(6, 'Main Warehouse', 12, 1, 6, 12, 510.00),
(7, 'Main Warehouse', 200, 15, 50, 100, 15.00),
(8, 'Main Warehouse', 150, 10, 30, 75, 8.00),
(9, 'Main Warehouse', 45, 2, 15, 25, 45.00);

-- Sample Customers
INSERT INTO customers (email, password_hash, first_name, last_name, phone, marketing_opt_in, email_verified) VALUES
('john.doe@email.com', '$2a$10$example_hash', 'John', 'Doe', '555-0123', true, true),
('jane.smith@email.com', '$2a$10$example_hash', 'Jane', 'Smith', '555-0124', false, true),
('bob.johnson@email.com', '$2a$10$example_hash', 'Bob', 'Johnson', '555-0125', true, false);

-- Customer Addresses  
INSERT INTO customer_addresses (customer_id, type, first_name, last_name, address_line1, city, state, postal_code, country, is_default) VALUES
(1, 'billing', 'John', 'Doe', '123 Main St', 'Anytown', 'CA', '90210', 'US', true),
(1, 'shipping', 'John', 'Doe', '123 Main St', 'Anytown', 'CA', '90210', 'US', true),
(2, 'billing', 'Jane', 'Smith', '456 Oak Ave', 'Springfield', 'IL', '62701', 'US', true),
(2, 'shipping', 'Jane', 'Smith', '789 Elm St', 'Springfield', 'IL', '62701', 'US', false),
(3, 'billing', 'Bob', 'Johnson', '321 Pine Rd', 'Portland', 'OR', '97205', 'US', true);

-- Sample Orders
INSERT INTO orders (customer_id, status, billing_first_name, billing_last_name, billing_address_line1, billing_city, billing_state, billing_postal_code, billing_country, shipping_first_name, shipping_last_name, shipping_address_line1, shipping_city, shipping_state, shipping_postal_code, shipping_country, subtotal, tax_amount, shipping_amount, total_amount, payment_status, shipping_method) VALUES
(1, 'delivered', 'John', 'Doe', '123 Main St', 'Anytown', 'CA', '90210', 'US', 'John', 'Doe', '123 Main St', 'Anytown', 'CA', '90210', 'US', 699.99, 56.00, 0.00, 755.99, 'paid', 'Free Shipping'),
(2, 'shipped', 'Jane', 'Smith', '456 Oak Ave', 'Springfield', 'IL', '62701', 'US', 'Jane', 'Smith', '789 Elm St', 'Springfield', 'IL', '62701', 'US', 149.99, 12.00, 9.99, 171.98, 'paid', 'Standard Shipping'),
(3, 'processing', 'Bob', 'Johnson', '321 Pine Rd', 'Portland', 'OR', '97205', 'US', 'Bob', 'Johnson', '321 Pine Rd', 'Portland', 'OR', '97205', 'US', 799.99, 64.00, 0.00, 863.99, 'paid', 'Free Shipping');

-- Order Items
INSERT INTO order_items (order_id, product_id, sku, name, quantity, unit_price, total_price) VALUES
(1, 1, 'DYS-V15-001', 'Dyson V15 Detect Absolute', 1, 699.99, 699.99),
(2, 3, 'SHK-NV356', 'Shark Navigator Lift-Away Professional', 1, 149.99, 149.99),
(3, 6, 'IRB-J7PLUS', 'iRobot Roomba j7+ Self-Emptying Robot', 1, 799.99, 799.99);

-- Order Status History
INSERT INTO order_status_history (order_id, status, notes, created_by) VALUES
(1, 'pending', 'Order received', 1),
(1, 'confirmed', 'Payment processed', 1),
(1, 'processing', 'Preparing for shipment', 2),
(1, 'shipped', 'Shipped via UPS', 2),
(1, 'delivered', 'Delivered successfully', 1),
(2, 'pending', 'Order received', 1),
(2, 'confirmed', 'Payment processed', 1),
(2, 'processing', 'Preparing for shipment', 2),
(2, 'shipped', 'Shipped via FedEx', 2),
(3, 'pending', 'Order received', 1),
(3, 'confirmed', 'Payment processed', 1),
(3, 'processing', 'Preparing for shipment', 2);

-- Coupons
INSERT INTO coupons (code, name, description, discount_type, discount_value, minimum_order_amount, usage_limit, is_active, expires_at) VALUES
('WELCOME10', 'Welcome 10% Off', 'Get 10% off your first order', 'percentage', 10.00, 100.00, 1000, true, '2025-12-31 23:59:59'),
('FREESHIP50', 'Free Shipping Over $50', 'Free shipping on orders over $50', 'free_shipping', 0.00, 50.00, NULL, true, '2025-12-31 23:59:59'),
('SPRING25', 'Spring Sale $25 Off', 'Save $25 on orders over $200', 'fixed_amount', 25.00, 200.00, 500, true, '2025-06-30 23:59:59');

-- Product Reviews
INSERT INTO product_reviews (product_id, customer_id, order_id, rating, title, review_text, is_verified_purchase, is_approved) VALUES
(1, 1, 1, 5, 'Amazing suction power!', 'This vacuum is incredible. The laser really shows how much dirt was hiding. Worth every penny.', true, true),
(3, 2, 2, 4, 'Great value for money', 'Solid vacuum cleaner. Good suction and the lift-away feature is handy for stairs.', true, true),
(6, 3, 3, 5, 'Life-changing robot vacuum', 'This robot vacuum is amazing. The pet waste avoidance actually works! No more accidents.', true, true);

-- Inventory Movements (sample transactions)
INSERT INTO inventory_movements (product_id, movement_type, quantity, reference_id, reference_type, notes, created_by) VALUES
-- Initial stock receipt
(1, 'purchase', 50, NULL, 'purchase_order', 'Initial inventory - Purchase Order #PO2024001', 1),
(2, 'purchase', 75, NULL, 'purchase_order', 'Initial inventory - Purchase Order #PO2024002', 1),
(3, 'purchase', 100, NULL, 'purchase_order', 'Initial inventory - Purchase Order #PO2024003', 1),

-- Sales from orders
(1, 'sale', -1, 1, 'order', 'Sale to customer - Order #1', 2),
(3, 'sale', -1, 2, 'order', 'Sale to customer - Order #2', 2),
(6, 'sale', -1, 3, 'order', 'Sale to customer - Order #3', 2),

-- Inventory adjustments
(2, 'adjustment', -2, NULL, 'adjustment', 'Damaged units removed from inventory', 1),
(7, 'purchase', 250, NULL, 'purchase_order', 'Filter restock - Purchase Order #PO2024004', 1);
