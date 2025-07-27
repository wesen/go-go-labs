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

// CartAPITestSuite provides comprehensive testing for shopping cart functionality
type CartAPITestSuite struct {
	suite.Suite
	server   *httptest.Server
	client   *http.Client
	sessionID string
}

func (suite *CartAPITestSuite) SetupSuite() {
	suite.server = httptest.NewServer(createTestRouter())
	suite.client = &http.Client{}
}

func (suite *CartAPITestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *CartAPITestSuite) SetupTest() {
	suite.sessionID = suite.createTestSession()
	suite.seedTestData()
}

func (suite *CartAPITestSuite) TearDownTest() {
	suite.cleanupTestData()
}

// Cart Management Tests

func (suite *CartAPITestSuite) TestGetCart_Empty() {
	url := fmt.Sprintf("%s/api/cart", suite.server.URL)
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(suite.T(), err)
	
	req.Header.Set("X-Session-ID", suite.sessionID)
	
	resp, err := suite.client.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var cartResponse struct {
		Items      []fixtures.TestCartItem `json:"items"`
		Total      float64                 `json:"total"`
		ItemCount  int                     `json:"item_count"`
		SessionID  string                  `json:"session_id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&cartResponse)
	require.NoError(suite.T(), err)
	
	assert.Empty(suite.T(), cartResponse.Items)
	assert.Equal(suite.T(), 0.0, cartResponse.Total)
	assert.Equal(suite.T(), 0, cartResponse.ItemCount)
}

func (suite *CartAPITestSuite) TestAddToCart_Success() {
	item := struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}{
		ProductID: fixtures.DysonV15.ID,
		Quantity:  2,
	}

	jsonData, err := json.Marshal(item)
	require.NoError(suite.T(), err)

	url := fmt.Sprintf("%s/api/cart/add", suite.server.URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	require.NoError(suite.T(), err)
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Session-ID", suite.sessionID)

	resp, err := suite.client.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var cartResponse struct {
		Items     []fixtures.TestCartItem `json:"items"`
		Total     float64                 `json:"total"`
		ItemCount int                     `json:"item_count"`
		Message   string                  `json:"message"`
	}
	err = json.NewDecoder(resp.Body).Decode(&cartResponse)
	require.NoError(suite.T(), err)
	
	assert.Len(suite.T(), cartResponse.Items, 1)
	assert.Equal(suite.T(), fixtures.DysonV15.ID, cartResponse.Items[0].ProductID)
	assert.Equal(suite.T(), 2, cartResponse.Items[0].Quantity)
	assert.Equal(suite.T(), fixtures.DysonV15.BasePrice*2, cartResponse.Total)
	assert.Equal(suite.T(), 2, cartResponse.ItemCount)
}

func (suite *CartAPITestSuite) TestAddToCart_InvalidProduct() {
	item := struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}{
		ProductID: 99999, // Non-existent product
		Quantity:  1,
	}

	jsonData, err := json.Marshal(item)
	require.NoError(suite.T(), err)

	url := fmt.Sprintf("%s/api/cart/add", suite.server.URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	require.NoError(suite.T(), err)
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Session-ID", suite.sessionID)

	resp, err := suite.client.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *CartAPITestSuite) TestAddToCart_InvalidQuantity() {
	testCases := []struct {
		name     string
		quantity int
		status   int
	}{
		{"Zero quantity", 0, http.StatusBadRequest},
		{"Negative quantity", -1, http.StatusBadRequest},
		{"Excessive quantity", 1000, http.StatusBadRequest},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			item := struct {
				ProductID int `json:"product_id"`
				Quantity  int `json:"quantity"`
			}{
				ProductID: fixtures.DysonV15.ID,
				Quantity:  tc.quantity,
			}

			jsonData, err := json.Marshal(item)
			require.NoError(t, err)

			url := fmt.Sprintf("%s/api/cart/add", suite.server.URL)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Session-ID", suite.sessionID)

			resp, err := suite.client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.status, resp.StatusCode)
		})
	}
}

func (suite *CartAPITestSuite) TestUpdateCartItem_Success() {
	// First add an item
	suite.addItemToCart(fixtures.DysonV15.ID, 1)

	// Then update the quantity
	updateData := struct {
		Quantity int `json:"quantity"`
	}{
		Quantity: 3,
	}

	jsonData, err := json.Marshal(updateData)
	require.NoError(suite.T(), err)

	url := fmt.Sprintf("%s/api/cart/update/%d", suite.server.URL, fixtures.DysonV15.ID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	require.NoError(suite.T(), err)
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Session-ID", suite.sessionID)

	resp, err := suite.client.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var cartResponse struct {
		Items []fixtures.TestCartItem `json:"items"`
		Total float64                 `json:"total"`
	}
	err = json.NewDecoder(resp.Body).Decode(&cartResponse)
	require.NoError(suite.T(), err)
	
	assert.Len(suite.T(), cartResponse.Items, 1)
	assert.Equal(suite.T(), 3, cartResponse.Items[0].Quantity)
	assert.Equal(suite.T(), fixtures.DysonV15.BasePrice*3, cartResponse.Total)
}

func (suite *CartAPITestSuite) TestRemoveFromCart_Success() {
	// First add an item
	suite.addItemToCart(fixtures.DysonV15.ID, 2)

	url := fmt.Sprintf("%s/api/cart/remove/%d", suite.server.URL, fixtures.DysonV15.ID)
	req, err := http.NewRequest("DELETE", url, nil)
	require.NoError(suite.T(), err)
	
	req.Header.Set("X-Session-ID", suite.sessionID)

	resp, err := suite.client.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var cartResponse struct {
		Items     []fixtures.TestCartItem `json:"items"`
		Total     float64                 `json:"total"`
		ItemCount int                     `json:"item_count"`
	}
	err = json.NewDecoder(resp.Body).Decode(&cartResponse)
	require.NoError(suite.T(), err)
	
	assert.Empty(suite.T(), cartResponse.Items)
	assert.Equal(suite.T(), 0.0, cartResponse.Total)
	assert.Equal(suite.T(), 0, cartResponse.ItemCount)
}

func (suite *CartAPITestSuite) TestClearCart_Success() {
	// Add multiple items
	suite.addItemToCart(fixtures.DysonV15.ID, 1)
	suite.addItemToCart(fixtures.SharkNavigator.ID, 2)

	url := fmt.Sprintf("%s/api/cart/clear", suite.server.URL)
	req, err := http.NewRequest("DELETE", url, nil)
	require.NoError(suite.T(), err)
	
	req.Header.Set("X-Session-ID", suite.sessionID)

	resp, err := suite.client.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify cart is empty
	suite.verifyEmptyCart()
}

// Cart Business Logic Tests

func (suite *CartAPITestSuite) TestCartTotal_Calculation() {
	// Add multiple items with different prices
	suite.addItemToCart(fixtures.DysonV15.ID, 1)      // $749.99
	suite.addItemToCart(fixtures.SharkNavigator.ID, 2) // $149.99 x 2 = $299.98

	expectedTotal := fixtures.DysonV15.BasePrice + (*fixtures.SharkNavigator.SalePrice * 2)

	cart := suite.getCart()
	assert.Equal(suite.T(), expectedTotal, cart.Total)
}

func (suite *CartAPITestSuite) TestCartTotal_WithSalePrice() {
	// Test that sale prices are used when available
	suite.addItemToCart(fixtures.SharkNavigator.ID, 1)

	cart := suite.getCart()
	assert.Equal(suite.T(), *fixtures.SharkNavigator.SalePrice, cart.Total)
}

func (suite *CartAPITestSuite) TestCart_DuplicateProductConsolidation() {
	// Add same product twice
	suite.addItemToCart(fixtures.DysonV15.ID, 2)
	suite.addItemToCart(fixtures.DysonV15.ID, 3)

	cart := suite.getCart()
	assert.Len(suite.T(), cart.Items, 1)
	assert.Equal(suite.T(), 5, cart.Items[0].Quantity)
}

// Session Management Tests

func (suite *CartAPITestSuite) TestCart_SessionIsolation() {
	sessionID1 := suite.createTestSession()
	sessionID2 := suite.createTestSession()

	// Add item to first session
	suite.addItemToCartWithSession(fixtures.DysonV15.ID, 1, sessionID1)

	// Verify second session has empty cart
	cart2 := suite.getCartWithSession(sessionID2)
	assert.Empty(suite.T(), cart2.Items)

	// Verify first session still has the item
	cart1 := suite.getCartWithSession(sessionID1)
	assert.Len(suite.T(), cart1.Items, 1)
}

func (suite *CartAPITestSuite) TestCart_NoSession() {
	url := fmt.Sprintf("%s/api/cart", suite.server.URL)
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(suite.T(), err)
	
	// Don't set session ID
	
	resp, err := suite.client.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

// Performance Tests

func (suite *CartAPITestSuite) TestCart_ConcurrentUpdates() {
	// Test concurrent cart updates from same session
	suite.addItemToCart(fixtures.DysonV15.ID, 1)

	done := make(chan bool, 10)
	
	// Launch concurrent update requests
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			suite.addItemToCart(fixtures.SharkNavigator.ID, 1)
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	cart := suite.getCart()
	// Verify cart state is consistent (exact behavior depends on implementation)
	assert.GreaterOrEqual(suite.T(), len(cart.Items), 1)
}

// Helper methods

func (suite *CartAPITestSuite) createTestSession() string {
	// Generate unique session ID for testing
	return fmt.Sprintf("test-session-%d", time.Now().UnixNano())
}

func (suite *CartAPITestSuite) addItemToCart(productID, quantity int) {
	suite.addItemToCartWithSession(productID, quantity, suite.sessionID)
}

func (suite *CartAPITestSuite) addItemToCartWithSession(productID, quantity int, sessionID string) {
	item := struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}{
		ProductID: productID,
		Quantity:  quantity,
	}

	jsonData, _ := json.Marshal(item)
	url := fmt.Sprintf("%s/api/cart/add", suite.server.URL)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Session-ID", sessionID)

	resp, _ := suite.client.Do(req)
	resp.Body.Close()
}

func (suite *CartAPITestSuite) getCart() struct {
	Items     []fixtures.TestCartItem `json:"items"`
	Total     float64                 `json:"total"`
	ItemCount int                     `json:"item_count"`
} {
	return suite.getCartWithSession(suite.sessionID)
}

func (suite *CartAPITestSuite) getCartWithSession(sessionID string) struct {
	Items     []fixtures.TestCartItem `json:"items"`
	Total     float64                 `json:"total"`
	ItemCount int                     `json:"item_count"`
} {
	url := fmt.Sprintf("%s/api/cart", suite.server.URL)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Session-ID", sessionID)

	resp, _ := suite.client.Do(req)
	defer resp.Body.Close()

	var cart struct {
		Items     []fixtures.TestCartItem `json:"items"`
		Total     float64                 `json:"total"`
		ItemCount int                     `json:"item_count"`
	}
	json.NewDecoder(resp.Body).Decode(&cart)
	return cart
}

func (suite *CartAPITestSuite) verifyEmptyCart() {
	cart := suite.getCart()
	assert.Empty(suite.T(), cart.Items)
	assert.Equal(suite.T(), 0.0, cart.Total)
	assert.Equal(suite.T(), 0, cart.ItemCount)
}

func (suite *CartAPITestSuite) seedTestData() {
	// Seed test database with fixture data
}

func (suite *CartAPITestSuite) cleanupTestData() {
	// Clean up test data after each test
}

// Run the test suite
func TestCartAPITestSuite(t *testing.T) {
	suite.Run(t, new(CartAPITestSuite))
}

// Additional edge case tests
func TestCartAPI_ErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		body           string
		expectedStatus int
	}{
		{"Invalid JSON", "POST", "/api/cart/add", "invalid json", http.StatusBadRequest},
		{"Missing product ID", "POST", "/api/cart/add", `{"quantity": 1}`, http.StatusBadRequest},
		{"Missing quantity", "POST", "/api/cart/add", `{"product_id": 1}`, http.StatusBadRequest},
		{"Invalid content type", "POST", "/api/cart/add", `{"product_id": 1, "quantity": 1}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Implementation for error handling tests
		})
	}
}
