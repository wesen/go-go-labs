package config

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// ConnectDatabase establishes a connection to the PostgreSQL database
func ConnectDatabase(cfg *DatabaseConfig) (*sqlx.DB, error) {
	if cfg == nil {
		return nil, errors.New("database configuration is nil")
	}

	// Set defaults
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == 0 {
		cfg.Port = 5432
	}
	if cfg.User == "" {
		cfg.User = "postgres"
	}
	if cfg.Name == "" {
		cfg.Name = "vacuum_store"
	}
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 25
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 10
	}
	if cfg.MaxLifetime == "" {
		cfg.MaxLifetime = "1h"
	}

	// Build connection string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	log.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("user", cfg.User).
		Str("database", cfg.Name).
		Str("ssl_mode", cfg.SSLMode).
		Msg("Connecting to database")

	// Connect to database
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	if lifetime, err := time.ParseDuration(cfg.MaxLifetime); err == nil {
		db.SetConnMaxLifetime(lifetime)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, errors.Wrap(err, "failed to ping database")
	}

	log.Info().Msg("Successfully connected to database")
	return db, nil
}

// MigrateDatabase runs database migrations
func MigrateDatabase(db *sqlx.DB) error {
	log.Info().Msg("Running database migrations")

	// Check if schema exists by trying to query a table
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'products')").Scan(&exists)
	if err != nil {
		return errors.Wrap(err, "failed to check if tables exist")
	}

	if !exists {
		log.Info().Msg("Tables do not exist, creating schema")
		
		// Read and execute schema file
		schemaSQL := `-- Minimal schema for the vacuum store
		
		CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			slug VARCHAR(100) UNIQUE NOT NULL,
			description TEXT,
			parent_id INTEGER REFERENCES categories(id),
			display_order INTEGER DEFAULT 0,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS brands (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL,
			slug VARCHAR(100) UNIQUE NOT NULL,
			description TEXT,
			logo_url VARCHAR(500),
			website_url VARCHAR(500),
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			sku VARCHAR(100) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) UNIQUE NOT NULL,
			description TEXT,
			short_description TEXT,
			category_id INTEGER REFERENCES categories(id),
			brand_id INTEGER REFERENCES brands(id),
			product_type VARCHAR(20) NOT NULL CHECK (product_type IN ('vacuum', 'accessory', 'part', 'bundle')),
			base_price DECIMAL(10,2) NOT NULL,
			sale_price DECIMAL(10,2),
			weight DECIMAL(8,2),
			dimensions_length DECIMAL(8,2),
			dimensions_width DECIMAL(8,2),
			dimensions_height DECIMAL(8,2),
			warranty_months INTEGER DEFAULT 12,
			is_active BOOLEAN DEFAULT true,
			is_featured BOOLEAN DEFAULT false,
			meta_title VARCHAR(255),
			meta_description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS customers (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255),
			first_name VARCHAR(100),
			last_name VARCHAR(100),
			phone VARCHAR(20),
			date_of_birth DATE,
			is_active BOOLEAN DEFAULT true,
			email_verified BOOLEAN DEFAULT false,
			marketing_opt_in BOOLEAN DEFAULT false,
			last_login TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS carts (
			id SERIAL PRIMARY KEY,
			customer_id INTEGER REFERENCES customers(id) ON DELETE CASCADE,
			session_id VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT cart_owner_check CHECK (customer_id IS NOT NULL OR session_id IS NOT NULL)
		);

		CREATE TABLE IF NOT EXISTS cart_items (
			id SERIAL PRIMARY KEY,
			cart_id INTEGER NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
			product_id INTEGER NOT NULL REFERENCES products(id),
			quantity INTEGER NOT NULL DEFAULT 1,
			unit_price DECIMAL(10,2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(cart_id, product_id)
		);

		CREATE TABLE IF NOT EXISTS inventory (
			id SERIAL PRIMARY KEY,
			product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			warehouse_location VARCHAR(100),
			quantity_on_hand INTEGER NOT NULL DEFAULT 0,
			quantity_reserved INTEGER NOT NULL DEFAULT 0,
			quantity_available INTEGER GENERATED ALWAYS AS (quantity_on_hand - quantity_reserved) STORED,
			reorder_level INTEGER DEFAULT 10,
			reorder_quantity INTEGER DEFAULT 50,
			cost_per_unit DECIMAL(10,2),
			last_received TIMESTAMP,
			last_sold TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(product_id, warehouse_location)
		);

		CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
		CREATE INDEX IF NOT EXISTS idx_products_brand ON products(brand_id);
		CREATE INDEX IF NOT EXISTS idx_products_active ON products(is_active);
		CREATE INDEX IF NOT EXISTS idx_products_featured ON products(is_featured);
		CREATE INDEX IF NOT EXISTS idx_cart_items_cart ON cart_items(cart_id);
		CREATE INDEX IF NOT EXISTS idx_cart_items_product ON cart_items(product_id);
		CREATE INDEX IF NOT EXISTS idx_carts_customer ON carts(customer_id);
		CREATE INDEX IF NOT EXISTS idx_carts_session ON carts(session_id);
		`

		if _, err := db.Exec(schemaSQL); err != nil {
			return errors.Wrap(err, "failed to create schema")
		}

		log.Info().Msg("Schema created successfully")
	} else {
		log.Info().Msg("Tables already exist, skipping schema creation")
	}

	return nil
}

// SeedDatabase populates the database with initial data
func SeedDatabase(db *sqlx.DB) error {
	log.Info().Msg("Seeding database with initial data")

	// Check if we already have data
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return errors.Wrap(err, "failed to check existing data")
	}

	if count > 0 {
		log.Info().Int("products", count).Msg("Database already has data, skipping seed")
		return nil
	}

	// Insert sample categories
	categories := []map[string]interface{}{
		{"name": "Upright Vacuums", "slug": "upright-vacuums", "description": "Traditional upright vacuum cleaners"},
		{"name": "Canister Vacuums", "slug": "canister-vacuums", "description": "Canister-style vacuum cleaners"},
		{"name": "Robot Vacuums", "slug": "robot-vacuums", "description": "Automated robot vacuum cleaners"},
		{"name": "Handheld Vacuums", "slug": "handheld-vacuums", "description": "Portable handheld vacuum cleaners"},
		{"name": "Accessories", "slug": "accessories", "description": "Vacuum accessories and attachments"},
		{"name": "Parts", "slug": "parts", "description": "Replacement parts for vacuum cleaners"},
	}

	for _, cat := range categories {
		_, err := db.NamedExec(`
			INSERT INTO categories (name, slug, description)
			VALUES (:name, :slug, :description)
			ON CONFLICT (slug) DO NOTHING`,
			cat)
		if err != nil {
			log.Warn().Err(err).Str("category", cat["name"].(string)).Msg("Failed to insert category")
		}
	}

	// Insert sample brands
	brands := []map[string]interface{}{
		{"name": "Dyson", "slug": "dyson", "description": "Premium vacuum cleaner manufacturer"},
		{"name": "Shark", "slug": "shark", "description": "Popular vacuum cleaner brand"},
		{"name": "Bissell", "slug": "bissell", "description": "Trusted carpet cleaning solutions"},
		{"name": "Hoover", "slug": "hoover", "description": "Classic vacuum cleaner brand"},
		{"name": "Miele", "slug": "miele", "description": "German engineering excellence"},
		{"name": "Eureka", "slug": "eureka", "description": "Affordable vacuum solutions"},
	}

	for _, brand := range brands {
		_, err := db.NamedExec(`
			INSERT INTO brands (name, slug, description)
			VALUES (:name, :slug, :description)
			ON CONFLICT (slug) DO NOTHING`,
			brand)
		if err != nil {
			log.Warn().Err(err).Str("brand", brand["name"].(string)).Msg("Failed to insert brand")
		}
	}

	// Insert sample products
	products := []map[string]interface{}{
		{
			"sku": "DYS-V15-001", "name": "Dyson V15 Detect", "slug": "dyson-v15-detect",
			"description": "Advanced cordless vacuum with laser dust detection", "base_price": 749.99,
			"product_type": "vacuum", "warranty_months": 24, "is_featured": true,
		},
		{
			"sku": "SHK-NAV-001", "name": "Shark Navigator", "slug": "shark-navigator",
			"description": "Powerful upright vacuum with lift-away technology", "base_price": 179.99,
			"product_type": "vacuum", "warranty_months": 12, "is_featured": true,
		},
		{
			"sku": "BIS-CRS-001", "name": "Bissell CrossWave Pet Pro", "slug": "bissell-crosswave-pet-pro",
			"description": "Multi-surface cleaner for pet owners", "base_price": 249.99,
			"product_type": "vacuum", "warranty_months": 12,
		},
		{
			"sku": "HOV-LIN-001", "name": "Hoover Linx", "slug": "hoover-linx",
			"description": "Lightweight cordless stick vacuum", "base_price": 99.99,
			"product_type": "vacuum", "warranty_months": 12,
		},
		{
			"sku": "MIE-CAT-001", "name": "Miele Complete C3", "slug": "miele-complete-c3",
			"description": "Premium canister vacuum with HEPA filtration", "base_price": 399.99,
			"product_type": "vacuum", "warranty_months": 36,
		},
	}

	for _, product := range products {
		_, err := db.NamedExec(`
			INSERT INTO products (sku, name, slug, description, base_price, product_type, warranty_months, is_featured)
			VALUES (:sku, :name, :slug, :description, :base_price, :product_type, :warranty_months, :is_featured)
			ON CONFLICT (sku) DO NOTHING`,
			product)
		if err != nil {
			log.Warn().Err(err).Str("product", product["name"].(string)).Msg("Failed to insert product")
		} else {
			// Add inventory for the product
			var productID int
			err = db.QueryRow("SELECT id FROM products WHERE sku = $1", product["sku"]).Scan(&productID)
			if err == nil {
				_, err = db.Exec(`
					INSERT INTO inventory (product_id, quantity_on_hand, quantity_reserved)
					VALUES ($1, $2, 0)
					ON CONFLICT (product_id, warehouse_location) DO NOTHING`,
					productID, 50) // Start with 50 units in stock
				if err != nil {
					log.Warn().Err(err).Int("product_id", productID).Msg("Failed to create inventory")
				}
			}
		}
	}

	log.Info().Msg("Database seeded successfully")
	return nil
}
