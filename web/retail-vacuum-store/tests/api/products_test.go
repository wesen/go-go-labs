package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	
	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/tests/fixtures"
)

// ProductAPITestSuite provides comprehensive testing for product API endpoints
type ProductAPITestSuite struct {
	suite.Suite
	server *httptest.Server
	client *http.Client
}

func (suite *ProductAPITestSuite) SetupSuite() {
	// Initialize test server and database
	suite.server = httptest.NewServer(createTestRouter())
	suite.client = &http.Client{}
}

func (suite *ProductAPITestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *ProductAPITestSuite) SetupTest() {
	// Clean and seed database for each test
	suite.seedTestData()
}

func (suite *ProductAPITestSuite) TearDownTest() {
	// Clean up test data
	suite.cleanupTestData()
}

// Test Product Catalog Endpoints

func (suite *ProductAPITestSuite) TestGetProducts_Success() {
	resp, err := suite.client.Get(suite.server.URL + "/api/products")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.Equal(suite.T(), "application/json", resp.Header.Get("Content-Type"))

	var products []fixtures.TestProduct
	err = json.NewDecoder(resp.Body).Decode(&products)
	require.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(products), 1)
}

func (suite *ProductAPITestSuite) TestGetProducts_WithPagination() {
	url := fmt.Sprintf("%s/api/products?page=1&limit=2", suite.server.URL)
	resp, err := suite.client.Get(url)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response struct {
		Products []fixtures.TestProduct `json:"products"`
		Page     int                    `json:"page"`
		Limit    int                    `json:"limit"`
		Total    int                    `json:"total"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), 1, response.Page)
	assert.Equal(suite.T(), 2, response.Limit)
	assert.LessOrEqual(suite.T(), len(response.Products), 2)
}

func (suite *ProductAPITestSuite) TestGetProducts_WithFilters() {
	// Test category filter
	url := fmt.Sprintf("%s/api/products?category_id=1", suite.server.URL)
	resp, err := suite.client.Get(url)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var products []fixtures.TestProduct
	err = json.NewDecoder(resp.Body).Decode(&products)
	require.NoError(suite.T(), err)
	
	for _, product := range products {
		assert.Equal(suite.T(), 1, product.CategoryID)
	}
}

func (suite *ProductAPITestSuite) TestGetProducts_WithPriceRange() {
	url := fmt.Sprintf("%s/api/products?min_price=100&max_price=500", suite.server.URL)
	resp, err := suite.client.Get(url)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var products []fixtures.TestProduct
	err = json.NewDecoder(resp.Body).Decode(&products)
	require.NoError(suite.T(), err)
	
	for _, product := range products {
		assert.GreaterOrEqual(suite.T(), product.BasePrice, 100.0)
		assert.LessOrEqual(suite.T(), product.BasePrice, 500.0)
	}
}

func (suite *ProductAPITestSuite) TestGetProduct_Success() {
	productID := fixtures.DysonV15.ID
	url := fmt.Sprintf("%s/api/products/%d", suite.server.URL, productID)
	
	resp, err := suite.client.Get(url)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var product fixtures.TestProduct
	err = json.NewDecoder(resp.Body).Decode(&product)
	require.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), productID, product.ID)
	assert.Equal(suite.T(), fixtures.DysonV15.Name, product.Name)
	assert.Equal(suite.T(), fixtures.DysonV15.SKU, product.SKU)
}

func (suite *ProductAPITestSuite) TestGetProduct_NotFound() {
	url := fmt.Sprintf("%s/api/products/99999", suite.server.URL)
	
	resp, err := suite.client.Get(url)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)

	var errorResponse struct {
		Error string `json:"error"`
	}
	err = json.NewDecoder(resp.Body).Decode(&errorResponse)
	require.NoError(suite.T(), err)
	assert.Contains(suite.T(), errorResponse.Error, "not found")
}

func (suite *ProductAPITestSuite) TestGetProduct_InvalidID() {
	url := fmt.Sprintf("%s/api/products/invalid", suite.server.URL)
	
	resp, err := suite.client.Get(url)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

// Test Search Functionality

func (suite *ProductAPITestSuite) TestSearchProducts_ByName() {
	url := fmt.Sprintf("%s/api/products/search?q=Dyson", suite.server.URL)
	
	resp, err := suite.client.Get(url)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var products []fixtures.TestProduct
	err = json.NewDecoder(resp.Body).Decode(&products)
	require.NoError(suite.T(), err)
	
	for _, product := range products {
		assert.Contains(suite.T(), product.Name, "Dyson")
	}
}

func (suite *ProductAPITestSuite) TestSearchProducts_EmptyQuery() {
	url := fmt.Sprintf("%s/api/products/search?q=", suite.server.URL)
	
	resp, err := suite.client.Get(url)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *ProductAPITestSuite) TestSearchProducts_NoResults() {
	url := fmt.Sprintf("%s/api/products/search?q=nonexistent", suite.server.URL)
	
	resp, err := suite.client.Get(url)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var products []fixtures.TestProduct
	err = json.NewDecoder(resp.Body).Decode(&products)
	require.NoError(suite.T(), err)
	assert.Empty(suite.T(), products)
}

// Test Categories API

func (suite *ProductAPITestSuite) TestGetCategories_Success() {
	resp, err := suite.client.Get(suite.server.URL + "/api/categories")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var categories []fixtures.TestCategory
	err = json.NewDecoder(resp.Body).Decode(&categories)
	require.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(categories), 1)
}

func (suite *ProductAPITestSuite) TestGetCategory_Success() {
	categoryID := fixtures.VacuumCategory.ID
	url := fmt.Sprintf("%s/api/categories/%d", suite.server.URL, categoryID)
	
	resp, err := suite.client.Get(url)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var category fixtures.TestCategory
	err = json.NewDecoder(resp.Body).Decode(&category)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), categoryID, category.ID)
}

// Test Brands API

func (suite *ProductAPITestSuite) TestGetBrands_Success() {
	resp, err := suite.client.Get(suite.server.URL + "/api/brands")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var brands []fixtures.TestBrand
	err = json.NewDecoder(resp.Body).Decode(&brands)
	require.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(brands), 1)
}

// Performance Tests

func (suite *ProductAPITestSuite) TestGetProducts_Performance() {
	// Test response time under normal load
	start := time.Now()
	
	resp, err := suite.client.Get(suite.server.URL + "/api/products")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()
	
	duration := time.Since(start)
	assert.Less(suite.T(), duration, 200*time.Millisecond, "Product API should respond within 200ms")
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

// Helper methods

func (suite *ProductAPITestSuite) seedTestData() {
	// Seed test database with fixture data
	// Implementation would depend on actual database setup
}

func (suite *ProductAPITestSuite) cleanupTestData() {
	// Clean up test data after each test
	// Implementation would depend on actual database setup
}

func createTestRouter() http.Handler {
	// Create test HTTP router with all endpoints
	// This would integrate with the actual handlers
	return http.NewServeMux()
}

// Run the test suite
func TestProductAPITestSuite(t *testing.T) {
	suite.Run(t, new(ProductAPITestSuite))
}

// Additional individual tests for edge cases
func TestProductAPI_ConcurrentRequests(t *testing.T) {
	// Test concurrent access to product API
	t.Parallel()
	
	// Implementation for concurrent request testing
}

func TestProductAPI_InvalidInputs(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{"Invalid page number", "/api/products?page=-1", http.StatusBadRequest},
		{"Invalid limit", "/api/products?limit=0", http.StatusBadRequest},
		{"Invalid price range", "/api/products?min_price=500&max_price=100", http.StatusBadRequest},
		{"SQL injection attempt", "/api/products?category_id=1; DROP TABLE products;", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Implementation for testing invalid inputs
		})
	}
}
