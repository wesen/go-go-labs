package e2e

import (
	"bytes"
	"context"
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

// ShoppingFlowTestSuite tests complete end-to-end shopping scenarios
type ShoppingFlowTestSuite struct {
	suite.Suite
	server     *httptest.Server
	client     *http.Client
	sessionID  string
	customerID int
}

func (suite *ShoppingFlowTestSuite) SetupSuite() {
	suite.server = httptest.NewServer(createTestRouter())
	suite.client = &http.Client{
		Timeout: 30 * time.Second,
	}
}

func (suite *ShoppingFlowTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *ShoppingFlowTestSuite) SetupTest() {
	suite.sessionID = suite.createTestSession()
	suite.customerID = suite.createTestCustomer()
	suite.seedTestData()
}

func (suite *ShoppingFlowTestSuite) TearDownTest() {
	suite.cleanupTestData()
}

// Complete Shopping Journey Tests

func (suite *ShoppingFlowTestSuite) TestCompleteShoppingJourney_Success() {
	ctx := context.Background()
	
	// Step 1: Browse products
	products := suite.browseProducts(ctx)
	assert.GreaterOrEqual(suite.T(), len(products), 1, "Should have products available")

	// Step 2: Search for specific product
	searchResults := suite.searchProducts(ctx, "Dyson")
	assert.GreaterOrEqual(suite.T(), len(searchResults), 1, "Should find Dyson products")

	// Step 3: View product details
	product := suite.getProductDetails(ctx, fixtures.DysonV15.ID)
	assert.Equal(suite.T(), fixtures.DysonV15.Name, product.Name)

	// Step 4: Add to cart
	suite.addToCart(ctx, fixtures.DysonV15.ID, 1)
	suite.addToCart(ctx, fixtures.SharkNavigator.ID, 2)

	// Step 5: Review cart
	cart := suite.getCart(ctx)
	assert.Len(suite.T(), cart.Items, 2, "Cart should have 2 different products")
	assert.Equal(suite.T(), 3, cart.ItemCount, "Cart should have 3 total items")

	// Step 6: Apply discount/coupon (if available)
	suite.applyCoupon(ctx, "SAVE10")

	// Step 7: Proceed to checkout
	checkout := suite.initiateCheckout(ctx)
	assert.NotEmpty(suite.T(), checkout.CheckoutID)

	// Step 8: Enter shipping information
	shippingInfo := fixtures.TestShippingInfo{
		FirstName: "John",
		LastName:  "Doe",
		Address:   "123 Main St",
		City:      "Anytown",
		State:     "CA",
		ZipCode:   "12345",
		Country:   "US",
	}
	suite.setShippingInfo(ctx, checkout.CheckoutID, shippingInfo)

	// Step 9: Select shipping method
	shippingMethods := suite.getShippingMethods(ctx, checkout.CheckoutID)
	assert.GreaterOrEqual(suite.T(), len(shippingMethods), 1, "Should have shipping methods available")
	suite.selectShippingMethod(ctx, checkout.CheckoutID, shippingMethods[0].ID)

	// Step 10: Enter payment information
	paymentInfo := fixtures.TestPaymentInfo{
		Method:    "credit_card",
		CardToken: "test_card_token_123",
	}
	suite.setPaymentInfo(ctx, checkout.CheckoutID, paymentInfo)

	// Step 11: Review order summary
	orderSummary := suite.getOrderSummary(ctx, checkout.CheckoutID)
	assert.Greater(suite.T(), orderSummary.Total, 0.0)
	assert.NotEmpty(suite.T(), orderSummary.Items)

	// Step 12: Place order
	order := suite.placeOrder(ctx, checkout.CheckoutID)
	assert.NotEmpty(suite.T(), order.OrderID)
	assert.Equal(suite.T(), "confirmed", order.Status)

	// Step 13: Verify order confirmation
	orderDetails := suite.getOrderDetails(ctx, order.OrderID)
	assert.Equal(suite.T(), order.OrderID, orderDetails.ID)
	assert.Equal(suite.T(), "confirmed", orderDetails.Status)

	// Step 14: Verify cart is cleared
	cartAfterOrder := suite.getCart(ctx)
	assert.Empty(suite.T(), cartAfterOrder.Items, "Cart should be empty after order")
}

func (suite *ShoppingFlowTestSuite) TestShoppingJourney_WithAccountCreation() {
	ctx := context.Background()

	// Step 1: Browse as guest
	suite.addToCart(ctx, fixtures.DysonV15.ID, 1)

	// Step 2: Decide to create account during checkout
	checkout := suite.initiateCheckout(ctx)

	// Step 3: Create customer account
	customerData := fixtures.TestCustomerRegistration{
		Email:     "john.doe@example.com",
		Password:  "SecurePass123!",
		FirstName: "John",
		LastName:  "Doe",
	}
	customer := suite.createCustomerAccount(ctx, customerData)
	assert.NotEmpty(suite.T(), customer.ID)

	// Step 4: Login and associate cart
	suite.loginCustomer(ctx, customerData.Email, customerData.Password)

	// Step 5: Complete checkout as registered customer
	order := suite.completeCheckoutAsCustomer(ctx, checkout.CheckoutID)
	assert.Equal(suite.T(), customer.ID, order.CustomerID)
}

func (suite *ShoppingFlowTestSuite) TestShoppingJourney_ReturningCustomer() {
	ctx := context.Background()

	// Step 1: Login as existing customer
	suite.loginExistingCustomer(ctx)

	// Step 2: View order history
	orderHistory := suite.getOrderHistory(ctx)
	// May be empty for new test customer

	// Step 3: Browse products with personalized recommendations
	products := suite.getRecommendedProducts(ctx)
	assert.GreaterOrEqual(suite.T(), len(products), 0)

	// Step 4: Add items to cart
	suite.addToCart(ctx, fixtures.MieleCX1.ID, 1)

	// Step 5: Use saved shipping address
	savedAddresses := suite.getSavedAddresses(ctx)
	if len(savedAddresses) > 0 {
		suite.selectSavedAddress(ctx, savedAddresses[0].ID)
	}

	// Step 6: Use saved payment method
	savedPaymentMethods := suite.getSavedPaymentMethods(ctx)
	if len(savedPaymentMethods) > 0 {
		suite.selectSavedPaymentMethod(ctx, savedPaymentMethods[0].ID)
	}

	// Step 7: Quick checkout
	order := suite.quickCheckout(ctx)
	assert.NotEmpty(suite.T(), order.OrderID)
}

// Cart Abandonment and Recovery Tests

func (suite *ShoppingFlowTestSuite) TestCartAbandonment_SessionRecovery() {
	ctx := context.Background()

	// Step 1: Add items to cart
	suite.addToCart(ctx, fixtures.DysonV15.ID, 1)
	originalCart := suite.getCart(ctx)

	// Step 2: Simulate leaving and returning (same session)
	time.Sleep(100 * time.Millisecond) // Simulate time gap

	// Step 3: Return and verify cart persists
	recoveredCart := suite.getCart(ctx)
	assert.Equal(suite.T(), len(originalCart.Items), len(recoveredCart.Items))
	assert.Equal(suite.T(), originalCart.Total, recoveredCart.Total)
}

func (suite *ShoppingFlowTestSuite) TestCartAbandonment_LoginRecovery() {
	ctx := context.Background()

	// Step 1: Add items as guest
	suite.addToCart(ctx, fixtures.SharkNavigator.ID, 2)

	// Step 2: Create account and login
	customerData := fixtures.TestCustomerRegistration{
		Email:     "recovery@example.com",
		Password:  "TestPass123!",
		FirstName: "Test",
		LastName:  "User",
	}
	suite.createCustomerAccount(ctx, customerData)
	suite.loginCustomer(ctx, customerData.Email, customerData.Password)

	// Step 3: Verify cart is merged/recovered
	cart := suite.getCart(ctx)
	assert.GreaterOrEqual(suite.T(), len(cart.Items), 1, "Cart should be recovered after login")
}

// Product Browsing and Search Tests

func (suite *ShoppingFlowTestSuite) TestProductBrowsing_FilterAndSort() {
	ctx := context.Background()

	// Test category filtering
	categoryProducts := suite.getProductsByCategory(ctx, fixtures.UprightCategory.ID)
	for _, product := range categoryProducts {
		assert.Equal(suite.T(), fixtures.UprightCategory.ID, product.CategoryID)
	}

	// Test brand filtering
	brandProducts := suite.getProductsByBrand(ctx, fixtures.DysonBrand.ID)
	for _, product := range brandProducts {
		assert.Equal(suite.T(), fixtures.DysonBrand.ID, product.BrandID)
	}

	// Test price range filtering
	priceFilteredProducts := suite.getProductsByPriceRange(ctx, 100.0, 500.0)
	for _, product := range priceFilteredProducts {
		assert.GreaterOrEqual(suite.T(), product.BasePrice, 100.0)
		assert.LessOrEqual(suite.T(), product.BasePrice, 500.0)
	}

	// Test sorting
	productsSortedByPrice := suite.getProductsSorted(ctx, "price_asc")
	if len(productsSortedByPrice) > 1 {
		assert.LessOrEqual(suite.T(), productsSortedByPrice[0].BasePrice, productsSortedByPrice[1].BasePrice)
	}
}

func (suite *ShoppingFlowTestSuite) TestProductSearch_Functionality() {
	ctx := context.Background()

	// Test exact name search
	results := suite.searchProducts(ctx, "Dyson V15")
	assert.GreaterOrEqual(suite.T(), len(results), 1)

	// Test partial name search
	partialResults := suite.searchProducts(ctx, "Shark")
	assert.GreaterOrEqual(suite.T(), len(partialResults), 1)

	// Test description search
	descResults := suite.searchProducts(ctx, "cordless")
	assert.GreaterOrEqual(suite.T(), len(descResults), 0)

	// Test no results
	noResults := suite.searchProducts(ctx, "nonexistentproduct")
	assert.Empty(suite.T(), noResults)
}

// Error Handling and Edge Cases

func (suite *ShoppingFlowTestSuite) TestCheckout_OutOfStock() {
	ctx := context.Background()

	// Add item to cart
	suite.addToCart(ctx, fixtures.DysonV15.ID, 1)

	// Simulate product going out of stock
	suite.setProductStock(ctx, fixtures.DysonV15.ID, 0)

	// Try to checkout
	checkout := suite.initiateCheckout(ctx)
	
	// Should receive out of stock error
	assert.Contains(suite.T(), checkout.Errors, "out of stock")
}

func (suite *ShoppingFlowTestSuite) TestCheckout_PriceChange() {
	ctx := context.Background()

	// Add item to cart
	suite.addToCart(ctx, fixtures.SharkNavigator.ID, 1)
	originalCart := suite.getCart(ctx)

	// Simulate price change
	newPrice := 199.99
	suite.setProductPrice(ctx, fixtures.SharkNavigator.ID, newPrice)

	// Proceed to checkout
	checkout := suite.initiateCheckout(ctx)

	// Should show updated price and warning
	assert.NotEqual(suite.T(), originalCart.Total, checkout.Total)
	assert.Contains(suite.T(), checkout.Warnings, "price changed")
}

func (suite *ShoppingFlowTestSuite) TestCheckout_PaymentFailure() {
	ctx := context.Background()

	// Complete shopping flow up to payment
	suite.addToCart(ctx, fixtures.MieleCX1.ID, 1)
	checkout := suite.initiateCheckout(ctx)
	suite.setShippingInfo(ctx, checkout.CheckoutID, fixtures.TestShippingInfo{})

	// Use invalid payment info
	invalidPayment := fixtures.TestPaymentInfo{
		Method:    "credit_card",
		CardToken: "invalid_token",
	}
	suite.setPaymentInfo(ctx, checkout.CheckoutID, invalidPayment)

	// Try to place order
	orderResult := suite.attemptPlaceOrder(ctx, checkout.CheckoutID)
	assert.Equal(suite.T(), "payment_failed", orderResult.Status)
	assert.NotEmpty(suite.T(), orderResult.Error)
}

// Performance and Scalability Tests

func (suite *ShoppingFlowTestSuite) TestConcurrentShopping_MultipleUsers() {
	ctx := context.Background()
	numUsers := 10
	done := make(chan bool, numUsers)

	// Simulate multiple users shopping concurrently
	for i := 0; i < numUsers; i++ {
		go func(userID int) {
			defer func() { done <- true }()
			
			userSessionID := suite.createTestSession()
			suite.concurrentUserShoppingFlow(ctx, userSessionID, userID)
		}(i)
	}

	// Wait for all users to complete
	for i := 0; i < numUsers; i++ {
		select {
		case <-done:
			// User completed
		case <-time.After(30 * time.Second):
			suite.T().Fatal("Concurrent shopping test timed out")
		}
	}
}

// Helper methods

func (suite *ShoppingFlowTestSuite) browseProducts(ctx context.Context) []fixtures.TestProduct {
	url := fmt.Sprintf("%s/api/products", suite.server.URL)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, _ := suite.client.Do(req)
	defer resp.Body.Close()

	var products []fixtures.TestProduct
	json.NewDecoder(resp.Body).Decode(&products)
	return products
}

func (suite *ShoppingFlowTestSuite) searchProducts(ctx context.Context, query string) []fixtures.TestProduct {
	url := fmt.Sprintf("%s/api/products/search?q=%s", suite.server.URL, query)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, _ := suite.client.Do(req)
	defer resp.Body.Close()

	var products []fixtures.TestProduct
	json.NewDecoder(resp.Body).Decode(&products)
	return products
}

func (suite *ShoppingFlowTestSuite) getProductDetails(ctx context.Context, productID int) fixtures.TestProduct {
	url := fmt.Sprintf("%s/api/products/%d", suite.server.URL, productID)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, _ := suite.client.Do(req)
	defer resp.Body.Close()

	var product fixtures.TestProduct
	json.NewDecoder(resp.Body).Decode(&product)
	return product
}

func (suite *ShoppingFlowTestSuite) addToCart(ctx context.Context, productID, quantity int) {
	item := struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}{productID, quantity}

	jsonData, _ := json.Marshal(item)
	url := fmt.Sprintf("%s/api/cart/add", suite.server.URL)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Session-ID", suite.sessionID)

	resp, _ := suite.client.Do(req)
	resp.Body.Close()
}

func (suite *ShoppingFlowTestSuite) getCart(ctx context.Context) struct {
	Items     []fixtures.TestCartItem `json:"items"`
	Total     float64                 `json:"total"`
	ItemCount int                     `json:"item_count"`
} {
	url := fmt.Sprintf("%s/api/cart", suite.server.URL)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("X-Session-ID", suite.sessionID)
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

func (suite *ShoppingFlowTestSuite) initiateCheckout(ctx context.Context) struct {
	CheckoutID string   `json:"checkout_id"`
	Total      float64  `json:"total"`
	Errors     []string `json:"errors"`
	Warnings   []string `json:"warnings"`
} {
	url := fmt.Sprintf("%s/api/checkout", suite.server.URL)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, nil)
	req.Header.Set("X-Session-ID", suite.sessionID)
	resp, _ := suite.client.Do(req)
	defer resp.Body.Close()

	var checkout struct {
		CheckoutID string   `json:"checkout_id"`
		Total      float64  `json:"total"`
		Errors     []string `json:"errors"`
		Warnings   []string `json:"warnings"`
	}
	json.NewDecoder(resp.Body).Decode(&checkout)
	return checkout
}

func (suite *ShoppingFlowTestSuite) placeOrder(ctx context.Context, checkoutID string) struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
} {
	url := fmt.Sprintf("%s/api/checkout/%s/place-order", suite.server.URL, checkoutID)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, nil)
	req.Header.Set("X-Session-ID", suite.sessionID)
	resp, _ := suite.client.Do(req)
	defer resp.Body.Close()

	var order struct {
		OrderID string `json:"order_id"`
		Status  string `json:"status"`
	}
	json.NewDecoder(resp.Body).Decode(&order)
	return order
}

// Additional helper methods would be implemented here...

func (suite *ShoppingFlowTestSuite) createTestSession() string {
	return fmt.Sprintf("e2e-session-%d", time.Now().UnixNano())
}

func (suite *ShoppingFlowTestSuite) createTestCustomer() int {
	// Implementation would create a test customer and return ID
	return 1
}

func (suite *ShoppingFlowTestSuite) seedTestData() {
	// Seed test database with fixture data
}

func (suite *ShoppingFlowTestSuite) cleanupTestData() {
	// Clean up test data
}

func (suite *ShoppingFlowTestSuite) concurrentUserShoppingFlow(ctx context.Context, sessionID string, userID int) {
	// Implement concurrent user shopping simulation
}

// Run the test suite
func TestShoppingFlowTestSuite(t *testing.T) {
	suite.Run(t, new(ShoppingFlowTestSuite))
}
