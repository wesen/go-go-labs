package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type APITester struct {
	baseURL    string
	httpClient *http.Client
	results    []TestResult
}

type TestResult struct {
	TestName    string        `json:"test_name"`
	Endpoint    string        `json:"endpoint"`
	Method      string        `json:"method"`
	StatusCode  int           `json:"status_code"`
	Duration    time.Duration `json:"duration"`
	Success     bool          `json:"success"`
	ErrorMsg    string        `json:"error_msg,omitempty"`
	ResponseSize int          `json:"response_size"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Brand       string  `json:"brand"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	InStock     bool    `json:"in_stock"`
	ImageURL    string  `json:"image_url"`
}

type CartItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Order struct {
	ID          int         `json:"id"`
	CustomerID  int         `json:"customer_id"`
	Items       []CartItem  `json:"items"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
}

func NewAPITester(baseURL string) *APITester {
	return &APITester{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		results: make([]TestResult, 0),
	}
}

func (a *APITester) RunTests(ctx context.Context) error {
	log.Info().Msg("Starting API tests for VacuumMart e-commerce platform")
	
	// Test product catalog endpoints
	if err := a.testProductCatalogAPIs(ctx); err != nil {
		return fmt.Errorf("product catalog API tests failed: %w", err)
	}
	
	// Test shopping cart endpoints
	if err := a.testShoppingCartAPIs(ctx); err != nil {
		return fmt.Errorf("shopping cart API tests failed: %w", err)
	}
	
	// Test order management endpoints
	if err := a.testOrderManagementAPIs(ctx); err != nil {
		return fmt.Errorf("order management API tests failed: %w", err)
	}
	
	// Test search and filtering endpoints
	if err := a.testSearchFilterAPIs(ctx); err != nil {
		return fmt.Errorf("search/filter API tests failed: %w", err)
	}
	
	// Test admin endpoints
	if err := a.testAdminAPIs(ctx); err != nil {
		return fmt.Errorf("admin API tests failed: %w", err)
	}
	
	// Generate test report
	a.generateReport()
	
	return nil
}

func (a *APITester) testProductCatalogAPIs(ctx context.Context) error {
	log.Info().Msg("Testing product catalog APIs")
	
	// Test get all products
	if err := a.testEndpoint(ctx, "GET", "/api/products", nil, 200, "Get all products"); err != nil {
		return err
	}
	
	// Test get product by ID
	if err := a.testEndpoint(ctx, "GET", "/api/products/1", nil, 200, "Get product by ID"); err != nil {
		return err
	}
	
	// Test get products by category
	categories := []string{"upright", "canister", "robot", "handheld"}
	for _, category := range categories {
		endpoint := fmt.Sprintf("/api/products?category=%s", category)
		testName := fmt.Sprintf("Get products by category: %s", category)
		if err := a.testEndpoint(ctx, "GET", endpoint, nil, 200, testName); err != nil {
			return err
		}
	}
	
	// Test get products by brand
	brands := []string{"Dyson", "Hoover", "Roomba", "Shark"}
	for _, brand := range brands {
		endpoint := fmt.Sprintf("/api/products?brand=%s", brand)
		testName := fmt.Sprintf("Get products by brand: %s", brand)
		if err := a.testEndpoint(ctx, "GET", endpoint, nil, 200, testName); err != nil {
			return err
		}
	}
	
	// Test price filtering
	if err := a.testEndpoint(ctx, "GET", "/api/products?min_price=100&max_price=500", nil, 200, "Price filtering"); err != nil {
		return err
	}
	
	// Test invalid product ID
	if err := a.testEndpoint(ctx, "GET", "/api/products/999999", nil, 404, "Invalid product ID should return 404"); err != nil {
		return err
	}
	
	return nil
}

func (a *APITester) testShoppingCartAPIs(ctx context.Context) error {
	log.Info().Msg("Testing shopping cart APIs")
	
	// Test create cart
	if err := a.testEndpoint(ctx, "POST", "/api/cart", nil, 201, "Create new cart"); err != nil {
		return err
	}
	
	// Test add item to cart
	cartItem := CartItem{ProductID: 1, Quantity: 2}
	if err := a.testEndpoint(ctx, "POST", "/api/cart/items", cartItem, 200, "Add item to cart"); err != nil {
		return err
	}
	
	// Test get cart contents
	if err := a.testEndpoint(ctx, "GET", "/api/cart", nil, 200, "Get cart contents"); err != nil {
		return err
	}
	
	// Test update cart item quantity
	updateItem := CartItem{ProductID: 1, Quantity: 3}
	if err := a.testEndpoint(ctx, "PUT", "/api/cart/items/1", updateItem, 200, "Update cart item quantity"); err != nil {
		return err
	}
	
	// Test remove item from cart
	if err := a.testEndpoint(ctx, "DELETE", "/api/cart/items/1", nil, 200, "Remove item from cart"); err != nil {
		return err
	}
	
	// Test clear cart
	if err := a.testEndpoint(ctx, "DELETE", "/api/cart", nil, 200, "Clear cart"); err != nil {
		return err
	}
	
	return nil
}

func (a *APITester) testOrderManagementAPIs(ctx context.Context) error {
	log.Info().Msg("Testing order management APIs")
	
	// Test create order
	order := Order{
		CustomerID: 1,
		Items: []CartItem{
			{ProductID: 1, Quantity: 1},
			{ProductID: 2, Quantity: 2},
		},
		TotalAmount: 599.99,
		Status:      "pending",
	}
	if err := a.testEndpoint(ctx, "POST", "/api/orders", order, 201, "Create new order"); err != nil {
		return err
	}
	
	// Test get order by ID
	if err := a.testEndpoint(ctx, "GET", "/api/orders/1", nil, 200, "Get order by ID"); err != nil {
		return err
	}
	
	// Test get orders by customer
	if err := a.testEndpoint(ctx, "GET", "/api/orders?customer_id=1", nil, 200, "Get orders by customer"); err != nil {
		return err
	}
	
	// Test update order status
	statusUpdate := map[string]string{"status": "confirmed"}
	if err := a.testEndpoint(ctx, "PUT", "/api/orders/1/status", statusUpdate, 200, "Update order status"); err != nil {
		return err
	}
	
	return nil
}

func (a *APITester) testSearchFilterAPIs(ctx context.Context) error {
	log.Info().Msg("Testing search and filter APIs")
	
	// Test search products
	searchTerms := []string{"dyson", "robot", "cordless", "pet hair"}
	for _, term := range searchTerms {
		endpoint := fmt.Sprintf("/api/products/search?q=%s", term)
		testName := fmt.Sprintf("Search products: %s", term)
		if err := a.testEndpoint(ctx, "GET", endpoint, nil, 200, testName); err != nil {
			return err
		}
	}
	
	// Test complex filtering
	complexFilter := "/api/products?category=upright&brand=Dyson&min_price=200&max_price=800&in_stock=true"
	if err := a.testEndpoint(ctx, "GET", complexFilter, nil, 200, "Complex product filtering"); err != nil {
		return err
	}
	
	// Test sorting
	sortOptions := []string{"price_asc", "price_desc", "name_asc", "rating_desc"}
	for _, sort := range sortOptions {
		endpoint := fmt.Sprintf("/api/products?sort=%s", sort)
		testName := fmt.Sprintf("Sort products: %s", sort)
		if err := a.testEndpoint(ctx, "GET", endpoint, nil, 200, testName); err != nil {
			return err
		}
	}
	
	return nil
}

func (a *APITester) testAdminAPIs(ctx context.Context) error {
	log.Info().Msg("Testing admin APIs")
	
	// Test get all orders (admin)
	if err := a.testEndpoint(ctx, "GET", "/api/admin/orders", nil, 200, "Admin: Get all orders"); err != nil {
		return err
	}
	
	// Test inventory management
	if err := a.testEndpoint(ctx, "GET", "/api/admin/inventory", nil, 200, "Admin: Get inventory status"); err != nil {
		return err
	}
	
	// Test update product inventory
	inventoryUpdate := map[string]interface{}{
		"stock_quantity": 50,
		"in_stock":       true,
	}
	if err := a.testEndpoint(ctx, "PUT", "/api/admin/products/1/inventory", inventoryUpdate, 200, "Admin: Update product inventory"); err != nil {
		return err
	}
	
	// Test sales analytics
	if err := a.testEndpoint(ctx, "GET", "/api/admin/analytics/sales", nil, 200, "Admin: Sales analytics"); err != nil {
		return err
	}
	
	return nil
}

func (a *APITester) testEndpoint(ctx context.Context, method, endpoint string, body interface{}, expectedStatus int, testName string) error {
	start := time.Now()
	
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}
	
	url := a.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	resp, err := a.httpClient.Do(req)
	duration := time.Since(start)
	
	result := TestResult{
		TestName:   testName,
		Endpoint:   endpoint,
		Method:     method,
		Duration:   duration,
		Success:    false,
	}
	
	if err != nil {
		result.ErrorMsg = err.Error()
		a.results = append(a.results, result)
		log.Error().Err(err).Str("endpoint", endpoint).Str("method", method).Msg("Request failed")
		return nil // Don't fail the entire test suite for one endpoint
	}
	
	defer resp.Body.Close()
	
	result.StatusCode = resp.StatusCode
	
	// Read response body to get size
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		result.ErrorMsg = fmt.Sprintf("failed to read response body: %v", err)
	} else {
		result.ResponseSize = len(respBody)
	}
	
	if resp.StatusCode == expectedStatus {
		result.Success = true
		log.Info().
			Str("test", testName).
			Str("endpoint", endpoint).
			Str("method", method).
			Int("status", resp.StatusCode).
			Dur("duration", duration).
			Int("response_size", result.ResponseSize).
			Msg("API test passed")
	} else {
		result.ErrorMsg = fmt.Sprintf("expected status %d, got %d", expectedStatus, resp.StatusCode)
		log.Warn().
			Str("test", testName).
			Str("endpoint", endpoint).
			Str("method", method).
			Int("expected_status", expectedStatus).
			Int("actual_status", resp.StatusCode).
			Msg("API test failed")
	}
	
	a.results = append(a.results, result)
	return nil
}

func (a *APITester) generateReport() {
	log.Info().Msg("Generating API test report")
	
	totalTests := len(a.results)
	passedTests := 0
	failedTests := 0
	totalDuration := time.Duration(0)
	totalResponseSize := 0
	
	for _, result := range a.results {
		if result.Success {
			passedTests++
		} else {
			failedTests++
		}
		totalDuration += result.Duration
		totalResponseSize += result.ResponseSize
	}
	
	log.Info().
		Int("total_tests", totalTests).
		Int("passed", passedTests).
		Int("failed", failedTests).
		Dur("total_duration", totalDuration).
		Dur("avg_duration", totalDuration/time.Duration(totalTests)).
		Int("total_response_size", totalResponseSize).
		Msg("API test summary")
	
	// Save detailed results to JSON file
	reportData := map[string]interface{}{
		"summary": map[string]interface{}{
			"total_tests":         totalTests,
			"passed_tests":        passedTests,
			"failed_tests":        failedTests,
			"success_rate":        float64(passedTests) / float64(totalTests) * 100,
			"total_duration_ms":   totalDuration.Milliseconds(),
			"avg_duration_ms":     totalDuration.Milliseconds() / int64(totalTests),
			"total_response_size": totalResponseSize,
		},
		"detailed_results": a.results,
		"test_timestamp":   time.Now().Format(time.RFC3339),
	}
	
	if jsonData, err := json.MarshalIndent(reportData, "", "  "); err == nil {
		log.Info().Str("report", string(jsonData)).Msg("Detailed API test report")
	}
}
