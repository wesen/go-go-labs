-- ==============================================
-- Complete Database Setup for Vacuum Cleaner Retail Store
-- ==============================================
-- Run this script to set up the complete database schema,
-- sample data, and views for the vacuum cleaner retail store.
--
-- Usage:
-- psql -U username -d database_name -f setup.sql
-- ==============================================

\echo 'Starting database setup for Vacuum Cleaner Retail Store...'

-- Create schema and tables
\echo 'Creating database schema...'
\i schema.sql

-- Create views for common queries
\echo 'Creating database views...'
\i views.sql

-- Insert sample data for development/testing
\echo 'Inserting sample data...'
\i sample_data.sql

\echo 'Database setup completed successfully!'
\echo ''
\echo '=== Quick Stats ==='
SELECT 
    (SELECT COUNT(*) FROM products) AS total_products,
    (SELECT COUNT(*) FROM categories) AS total_categories,
    (SELECT COUNT(*) FROM brands) AS total_brands,
    (SELECT COUNT(*) FROM customers) AS total_customers,
    (SELECT COUNT(*) FROM orders) AS total_orders;

\echo ''
\echo '=== Inventory Summary ==='
SELECT 
    product_name,
    quantity_available,
    stock_alert
FROM inventory_status 
ORDER BY 
    CASE stock_alert 
        WHEN 'OUT_OF_STOCK' THEN 1
        WHEN 'NEEDS_REORDER' THEN 2
        WHEN 'LOW_STOCK' THEN 3
        ELSE 4
    END,
    product_name;

\echo ''
\echo 'Setup complete! The database is ready for use.'
