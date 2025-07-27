-- ==============================================
-- Database Views for Vacuum Cleaner Retail Store
-- ==============================================

-- Product catalog view with computed fields
CREATE VIEW product_catalog AS
SELECT 
    p.id,
    p.sku,
    p.name,
    p.slug,
    p.short_description,
    p.description,
    p.product_type,
    p.base_price,
    p.sale_price,
    COALESCE(p.sale_price, p.base_price) AS current_price,
    CASE 
        WHEN p.sale_price IS NOT NULL AND p.sale_price < p.base_price 
        THEN ROUND(((p.base_price - p.sale_price) / p.base_price * 100), 0)
        ELSE 0 
    END AS discount_percentage,
    p.weight,
    p.warranty_months,
    p.is_active,
    p.is_featured,
    c.name AS category_name,
    c.slug AS category_slug,
    b.name AS brand_name,
    b.slug AS brand_slug,
    pi.url AS primary_image_url,
    i.quantity_available,
    CASE WHEN i.quantity_available > 0 THEN true ELSE false END AS in_stock,
    CASE 
        WHEN i.quantity_available = 0 THEN 'out_of_stock'
        WHEN i.quantity_available <= i.reorder_level THEN 'low_stock'
        ELSE 'in_stock'
    END AS stock_status,
    AVG(pr.rating) AS average_rating,
    COUNT(pr.id) AS review_count,
    p.created_at,
    p.updated_at
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN brands b ON p.brand_id = b.id
LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.is_primary = true
LEFT JOIN inventory i ON p.id = i.product_id
LEFT JOIN product_reviews pr ON p.id = pr.product_id AND pr.is_approved = true
WHERE p.is_active = true
GROUP BY p.id, c.name, c.slug, b.name, b.slug, pi.url, i.quantity_available, i.reorder_level;

-- Order summary view with customer information
CREATE VIEW order_summary AS
SELECT 
    o.id,
    o.order_number,
    o.status,
    o.payment_status,
    o.total_amount,
    o.currency,
    o.shipping_method,
    o.tracking_number,
    CONCAT(o.billing_first_name, ' ', o.billing_last_name) AS billing_name,
    CONCAT(o.shipping_first_name, ' ', o.shipping_last_name) AS shipping_name,
    CONCAT(o.shipping_address_line1, ', ', o.shipping_city, ', ', o.shipping_state, ' ', o.shipping_postal_code) AS shipping_address,
    o.created_at AS order_date,
    o.shipped_at,
    o.delivered_at,
    c.email AS customer_email,
    c.first_name AS customer_first_name,
    c.last_name AS customer_last_name,
    COUNT(oi.id) AS item_count,
    o.guest_email
FROM orders o
LEFT JOIN customers c ON o.customer_id = c.id
LEFT JOIN order_items oi ON o.id = oi.order_id
GROUP BY o.id, c.email, c.first_name, c.last_name;

-- Inventory status view with reorder alerts
CREATE VIEW inventory_status AS
SELECT 
    i.id,
    p.sku,
    p.name AS product_name,
    i.warehouse_location,
    i.quantity_on_hand,
    i.quantity_reserved,
    i.quantity_available,
    i.reorder_level,
    i.reorder_quantity,
    i.cost_per_unit,
    CASE 
        WHEN i.quantity_available = 0 THEN 'OUT_OF_STOCK'
        WHEN i.quantity_available <= i.reorder_level THEN 'NEEDS_REORDER'
        WHEN i.quantity_available <= (i.reorder_level * 1.5) THEN 'LOW_STOCK'
        ELSE 'ADEQUATE'
    END AS stock_alert,
    i.last_received,
    i.last_sold,
    i.updated_at
FROM inventory i
JOIN products p ON i.product_id = p.id
WHERE p.is_active = true;

-- Customer order history view
CREATE VIEW customer_order_history AS
SELECT 
    c.id AS customer_id,
    c.email,
    c.first_name,
    c.last_name,
    COUNT(o.id) AS total_orders,
    SUM(o.total_amount) AS total_spent,
    AVG(o.total_amount) AS average_order_value,
    MAX(o.created_at) AS last_order_date,
    MIN(o.created_at) AS first_order_date,
    COUNT(CASE WHEN o.status = 'delivered' THEN 1 END) AS completed_orders,
    COUNT(CASE WHEN o.status = 'cancelled' THEN 1 END) AS cancelled_orders
FROM customers c
LEFT JOIN orders o ON c.id = o.customer_id
GROUP BY c.id, c.email, c.first_name, c.last_name;

-- Product performance view
CREATE VIEW product_performance AS
SELECT 
    p.id,
    p.sku,
    p.name,
    p.category_id,
    c.name AS category_name,
    p.brand_id,
    b.name AS brand_name,
    COUNT(oi.id) AS units_sold,
    SUM(oi.total_price) AS total_revenue,
    AVG(oi.unit_price) AS average_selling_price,
    COUNT(DISTINCT oi.order_id) AS orders_count,
    AVG(pr.rating) AS average_rating,
    COUNT(pr.id) AS review_count,
    MAX(oi.created_at) AS last_sold_date,
    i.quantity_available AS current_stock
FROM products p
LEFT JOIN order_items oi ON p.id = oi.product_id
LEFT JOIN orders o ON oi.order_id = o.id AND o.status NOT IN ('cancelled', 'refunded')
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN brands b ON p.brand_id = b.id
LEFT JOIN product_reviews pr ON p.id = pr.product_id AND pr.is_approved = true
LEFT JOIN inventory i ON p.id = i.product_id
WHERE p.is_active = true
GROUP BY p.id, p.sku, p.name, p.category_id, c.name, p.brand_id, b.name, i.quantity_available;

-- Top selling products view
CREATE VIEW top_selling_products AS
SELECT 
    pp.*,
    ROW_NUMBER() OVER (ORDER BY pp.units_sold DESC) AS sales_rank,
    ROW_NUMBER() OVER (ORDER BY pp.total_revenue DESC) AS revenue_rank
FROM product_performance pp
WHERE pp.units_sold > 0
ORDER BY pp.units_sold DESC;

-- Low stock alerts view
CREATE VIEW low_stock_alerts AS
SELECT 
    p.id,
    p.sku,
    p.name,
    i.warehouse_location,
    i.quantity_available,
    i.reorder_level,
    i.reorder_quantity,
    CASE 
        WHEN i.quantity_available = 0 THEN 'CRITICAL'
        WHEN i.quantity_available <= (i.reorder_level * 0.5) THEN 'URGENT'
        WHEN i.quantity_available <= i.reorder_level THEN 'MODERATE'
        ELSE 'LOW'
    END AS priority,
    pp.units_sold AS sales_velocity,
    i.last_received,
    i.updated_at
FROM products p
JOIN inventory i ON p.id = i.product_id
LEFT JOIN product_performance pp ON p.id = pp.id
WHERE p.is_active = true 
    AND i.quantity_available <= i.reorder_level
ORDER BY 
    CASE 
        WHEN i.quantity_available = 0 THEN 1
        WHEN i.quantity_available <= (i.reorder_level * 0.5) THEN 2
        WHEN i.quantity_available <= i.reorder_level THEN 3
        ELSE 4
    END,
    pp.units_sold DESC NULLS LAST;

-- Category performance view
CREATE VIEW category_performance AS
SELECT 
    c.id,
    c.name,
    c.slug,
    COUNT(p.id) AS product_count,
    COUNT(CASE WHEN p.is_active = true THEN 1 END) AS active_product_count,
    SUM(pp.units_sold) AS total_units_sold,
    SUM(pp.total_revenue) AS total_revenue,
    AVG(pp.average_selling_price) AS avg_product_price,
    AVG(pp.average_rating) AS avg_category_rating,
    SUM(i.quantity_available) AS total_inventory
FROM categories c
LEFT JOIN products p ON c.id = p.category_id
LEFT JOIN product_performance pp ON p.id = pp.id
LEFT JOIN inventory i ON p.id = i.product_id
GROUP BY c.id, c.name, c.slug
ORDER BY total_revenue DESC NULLS LAST;

-- Brand performance view
CREATE VIEW brand_performance AS
SELECT 
    b.id,
    b.name,
    b.slug,
    COUNT(p.id) AS product_count,
    COUNT(CASE WHEN p.is_active = true THEN 1 END) AS active_product_count,
    SUM(pp.units_sold) AS total_units_sold,
    SUM(pp.total_revenue) AS total_revenue,
    AVG(pp.average_selling_price) AS avg_product_price,
    AVG(pp.average_rating) AS avg_brand_rating,
    SUM(i.quantity_available) AS total_inventory
FROM brands b
LEFT JOIN products p ON b.id = p.brand_id
LEFT JOIN product_performance pp ON p.id = pp.id
LEFT JOIN inventory i ON p.id = i.product_id
GROUP BY b.id, b.name, b.slug
ORDER BY total_revenue DESC NULLS LAST;

-- Recent orders view for admin dashboard
CREATE VIEW recent_orders AS
SELECT 
    o.id,
    o.order_number,
    o.status,
    o.total_amount,
    COALESCE(c.first_name || ' ' || c.last_name, 'Guest') AS customer_name,
    COALESCE(c.email, o.guest_email) AS customer_email,
    o.created_at,
    COUNT(oi.id) AS item_count
FROM orders o
LEFT JOIN customers c ON o.customer_id = c.id
LEFT JOIN order_items oi ON o.id = oi.order_id
GROUP BY o.id, o.order_number, o.status, o.total_amount, c.first_name, c.last_name, c.email, o.guest_email, o.created_at
ORDER BY o.created_at DESC
LIMIT 50;

-- Admin dashboard stats view
CREATE VIEW admin_dashboard_stats AS
SELECT 
    (SELECT COUNT(*) FROM orders WHERE created_at >= CURRENT_DATE) AS orders_today,
    (SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE created_at >= CURRENT_DATE AND status NOT IN ('cancelled', 'refunded')) AS revenue_today,
    (SELECT COUNT(*) FROM orders WHERE created_at >= CURRENT_DATE - INTERVAL '7 days') AS orders_this_week,
    (SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE created_at >= CURRENT_DATE - INTERVAL '7 days' AND status NOT IN ('cancelled', 'refunded')) AS revenue_this_week,
    (SELECT COUNT(*) FROM orders WHERE created_at >= CURRENT_DATE - INTERVAL '30 days') AS orders_this_month,
    (SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE created_at >= CURRENT_DATE - INTERVAL '30 days' AND status NOT IN ('cancelled', 'refunded')) AS revenue_this_month,
    (SELECT COUNT(*) FROM products WHERE is_active = true) AS active_products,
    (SELECT COUNT(*) FROM customers WHERE created_at >= CURRENT_DATE - INTERVAL '7 days') AS new_customers_this_week,
    (SELECT COUNT(*) FROM inventory WHERE quantity_available <= reorder_level) AS low_stock_items,
    (SELECT COUNT(*) FROM orders WHERE status = 'pending') AS pending_orders;

-- Create indexes on views that might be used for filtering
CREATE INDEX IF NOT EXISTS idx_product_catalog_category ON products(category_id) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_product_catalog_brand ON products(brand_id) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_product_catalog_featured ON products(is_featured) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_product_catalog_price ON products(COALESCE(sale_price, base_price)) WHERE is_active = true;
