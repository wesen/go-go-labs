package utils

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	
	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/tests/fixtures"
)

// TestDatabase provides utilities for test database management
type TestDatabase struct {
	DB       *sql.DB
	DSN      string
	TestName string
}

// TestHTTPClient provides utilities for HTTP testing
type TestHTTPClient struct {
	Server    *httptest.Server
	Client    *http.Client
	SessionID string
	Headers   map[string]string
}

// TestDataManager handles test data lifecycle
type TestDataManager struct {
	db       *TestDatabase
	fixtures map[string]interface{}
	cleanup  []func() error
}

// DatabaseTestSetup creates and initializes a test database
func DatabaseTestSetup(t *testing.T) *TestDatabase {
	// Create in-memory SQLite database for tests
	testDB := &TestDatabase{
		DSN:      ":memory:",
		TestName: t.Name(),
	}

	db, err := sql.Open("sqlite3", testDB.DSN)
	require.NoError(t, err)
	
	testDB.DB = db

	// Load and execute schema
	err = testDB.LoadSchema()
	require.NoError(t, err)

	return testDB
}

// LoadSchema loads the database schema from schema.sql
func (td *TestDatabase) LoadSchema() error {
	schemaPath := filepath.Join("..", "..", "db", "schema.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		// Fallback to relative path
		schemaPath = "../../db/schema.sql"
		schema, err = os.ReadFile(schemaPath)
		if err != nil {
			return fmt.Errorf("failed to load schema: %w", err)
		}
	}

	// Execute schema (split by semicolons for multiple statements)
	statements := strings.Split(string(schema), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err := td.DB.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute schema statement: %w", err)
		}
	}

	return nil
}

// SeedData populates the test database with fixture data
func (td *TestDatabase) SeedData() error {
	// Seed categories
	categories := fixtures.CreateTestCategories()
	for _, category := range categories {
		_, err := td.DB.Exec(`
			INSERT INTO categories (id, name, slug, description, parent_id, is_active)
			VALUES (?, ?, ?, ?, ?, ?)`,
			category.ID, category.Name, category.Slug, category.Description,
			category.ParentID, category.IsActive)
		if err != nil {
			return fmt.Errorf("failed to seed category: %w", err)
		}
	}

	// Seed brands
	brands := fixtures.CreateTestBrands()
	for _, brand := range brands {
		_, err := td.DB.Exec(`
			INSERT INTO brands (id, name, slug, description, is_active)
			VALUES (?, ?, ?, ?, ?)`,
			brand.ID, brand.Name, brand.Slug, brand.Description, brand.IsActive)
		if err != nil {
			return fmt.Errorf("failed to seed brand: %w", err)
		}
	}

	// Seed products
	products := fixtures.CreateTestProducts()
	for _, product := range products {
		_, err := td.DB.Exec(`
			INSERT INTO products (id, sku, name, description, base_price, sale_price, 
			                    category_id, brand_id, is_active, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			product.ID, product.SKU, product.Name, product.Description,
			product.BasePrice, product.SalePrice, product.CategoryID,
			product.BrandID, product.IsActive, product.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to seed product: %w", err)
		}
	}

	return nil
}

// CleanData removes all test data from the database
func (td *TestDatabase) CleanData() error {
	tables := []string{
		"order_items", "orders", "cart_items", "products", 
		"categories", "brands", "customers", "admin_users",
	}

	for _, table := range tables {
		_, err := td.DB.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("failed to clean table %s: %w", table, err)
		}
	}

	return nil
}

// Close closes the test database connection
func (td *TestDatabase) Close() error {
	if td.DB != nil {
		return td.DB.Close()
	}
	return nil
}

// HTTPTestSetup creates a test HTTP client with server
func HTTPTestSetup(handler http.Handler) *TestHTTPClient {
	server := httptest.NewServer(handler)
	
	return &TestHTTPClient{
		Server: server,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		SessionID: GenerateTestSessionID(),
		Headers:   make(map[string]string),
	}
}

// SetHeader sets a default header for all requests
func (tc *TestHTTPClient) SetHeader(key, value string) {
	tc.Headers[key] = value
}

// GET performs a GET request with default headers
func (tc *TestHTTPClient) GET(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", tc.Server.URL+path, nil)
	if err != nil {
		return nil, err
	}

	tc.addDefaultHeaders(req)
	return tc.Client.Do(req)
}

// POST performs a POST request with JSON body
func (tc *TestHTTPClient) POST(path string, body interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", tc.Server.URL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	tc.addDefaultHeaders(req)
	return tc.Client.Do(req)
}

// PUT performs a PUT request with JSON body
func (tc *TestHTTPClient) PUT(path string, body interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", tc.Server.URL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	tc.addDefaultHeaders(req)
	return tc.Client.Do(req)
}

// DELETE performs a DELETE request
func (tc *TestHTTPClient) DELETE(path string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", tc.Server.URL+path, nil)
	if err != nil {
		return nil, err
	}

	tc.addDefaultHeaders(req)
	return tc.Client.Do(req)
}

// addDefaultHeaders adds session ID and other default headers
func (tc *TestHTTPClient) addDefaultHeaders(req *http.Request) {
	req.Header.Set("X-Session-ID", tc.SessionID)
	
	for key, value := range tc.Headers {
		req.Header.Set(key, value)
	}
}

// Close closes the test server
func (tc *TestHTTPClient) Close() {
	tc.Server.Close()
}

// ParseJSONResponse parses HTTP response body into a struct
func ParseJSONResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

// ReadResponseBody reads the entire response body as string
func ReadResponseBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// TestDataManagerSetup creates a test data manager
func TestDataManagerSetup(t *testing.T, db *TestDatabase) *TestDataManager {
	return &TestDataManager{
		db:       db,
		fixtures: make(map[string]interface{}),
		cleanup:  make([]func() error, 0),
	}
}

// CreateProduct creates a test product and tracks it for cleanup
func (tdm *TestDataManager) CreateProduct(product fixtures.TestProduct) error {
	_, err := tdm.db.DB.Exec(`
		INSERT INTO products (id, sku, name, description, base_price, sale_price, 
		                    category_id, brand_id, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		product.ID, product.SKU, product.Name, product.Description,
		product.BasePrice, product.SalePrice, product.CategoryID,
		product.BrandID, product.IsActive, product.CreatedAt)
	
	if err != nil {
		return err
	}

	// Track for cleanup
	tdm.cleanup = append(tdm.cleanup, func() error {
		_, err := tdm.db.DB.Exec("DELETE FROM products WHERE id = ?", product.ID)
		return err
	})

	return nil
}

// CreateCategory creates a test category and tracks it for cleanup
func (tdm *TestDataManager) CreateCategory(category fixtures.TestCategory) error {
	_, err := tdm.db.DB.Exec(`
		INSERT INTO categories (id, name, slug, description, parent_id, is_active)
		VALUES (?, ?, ?, ?, ?, ?)`,
		category.ID, category.Name, category.Slug, category.Description,
		category.ParentID, category.IsActive)
	
	if err != nil {
		return err
	}

	tdm.cleanup = append(tdm.cleanup, func() error {
		_, err := tdm.db.DB.Exec("DELETE FROM categories WHERE id = ?", category.ID)
		return err
	})

	return nil
}

// CreateOrder creates a test order and tracks it for cleanup
func (tdm *TestDataManager) CreateOrder(order fixtures.TestOrder) error {
	// First create the order
	_, err := tdm.db.DB.Exec(`
		INSERT INTO orders (id, customer_id, total_amount, status, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		order.ID, order.CustomerID, order.TotalAmount, order.Status, order.CreatedAt)
	
	if err != nil {
		return err
	}

	// Then create order items
	for _, item := range order.Items {
		_, err := tdm.db.DB.Exec(`
			INSERT INTO order_items (order_id, product_id, quantity, unit_price)
			VALUES (?, ?, ?, ?)`,
			order.ID, item.ProductID, item.Quantity, 0) // Price would be calculated
		if err != nil {
			return err
		}
	}

	tdm.cleanup = append(tdm.cleanup, func() error {
		_, err1 := tdm.db.DB.Exec("DELETE FROM order_items WHERE order_id = ?", order.ID)
		_, err2 := tdm.db.DB.Exec("DELETE FROM orders WHERE id = ?", order.ID)
		if err1 != nil {
			return err1
		}
		return err2
	})

	return nil
}

// Cleanup runs all cleanup functions
func (tdm *TestDataManager) Cleanup() error {
	var lastErr error
	// Run cleanup in reverse order
	for i := len(tdm.cleanup) - 1; i >= 0; i-- {
		if err := tdm.cleanup[i](); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Utility Functions

// GenerateTestSessionID generates a unique session ID for testing
func GenerateTestSessionID() string {
	return fmt.Sprintf("test-session-%d-%d", time.Now().Unix(), rand.Intn(10000))
}

// GenerateTestSKU generates a unique SKU for testing
func GenerateTestSKU() string {
	return fmt.Sprintf("TEST-SKU-%d", time.Now().UnixNano())
}

// RandomString generates a random string of specified length
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// RandomEmail generates a random email address for testing
func RandomEmail() string {
	return fmt.Sprintf("test_%s@example.com", RandomString(8))
}

// RandomPrice generates a random price between min and max
func RandomPrice(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// WaitForCondition waits for a condition to be true with timeout
func WaitForCondition(condition func() bool, timeout time.Duration, interval time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(interval)
	}
	return false
}

// AssertJSONResponse asserts that response contains expected JSON fields
func AssertJSONResponse(t *testing.T, resp *http.Response, expectedFields map[string]interface{}) {
	defer resp.Body.Close()
	
	var responseData map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseData)
	require.NoError(t, err)

	for field, expectedValue := range expectedFields {
		actualValue, exists := responseData[field]
		require.True(t, exists, "Field %s should exist in response", field)
		require.Equal(t, expectedValue, actualValue, "Field %s should match expected value", field)
	}
}

// MockServer creates a mock HTTP server for external services
func MockServer(responses map[string]func(http.ResponseWriter, *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	
	for path, handler := range responses {
		mux.HandleFunc(path, handler)
	}

	return httptest.NewServer(mux)
}

// LoadTestFixture loads test data from a JSON file
func LoadTestFixture(filename string, target interface{}) error {
	fixturePath := filepath.Join("fixtures", filename)
	data, err := os.ReadFile(fixturePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

// SaveTestResults saves test results to a file for analysis
func SaveTestResults(filename string, data interface{}) error {
	resultsDir := "test_results"
	err := os.MkdirAll(resultsDir, 0755)
	if err != nil {
		return err
	}

	filePath := filepath.Join(resultsDir, filename)
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}

// TimeMeasure measures execution time of a function
func TimeMeasure(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}

// MemorySnapshot captures memory usage for leak detection
type MemorySnapshot struct {
	Timestamp time.Time
	// In a real implementation, this would contain actual memory metrics
	AllocMB   float64
	SysMB     float64
	NumGC     uint32
}

// TakeMemorySnapshot captures current memory usage
func TakeMemorySnapshot() MemorySnapshot {
	// In a real implementation, this would use runtime.MemStats
	return MemorySnapshot{
		Timestamp: time.Now(),
		AllocMB:   rand.Float64() * 100, // Mock data
		SysMB:     rand.Float64() * 200,
		NumGC:     rand.Uint32(),
	}
}

// CompareMemorySnapshots compares two memory snapshots
func CompareMemorySnapshots(before, after MemorySnapshot) bool {
	// Simple leak detection: if memory increased by more than 50MB
	return (after.AllocMB - before.AllocMB) < 50.0
}
