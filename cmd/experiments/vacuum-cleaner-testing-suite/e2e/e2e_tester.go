package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type E2ETester struct {
	baseURL    string
	httpClient *http.Client
	scenarios  []TestScenario
}

type TestScenario struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Steps       []ScenarioStep `json:"steps"`
	Success     bool          `json:"success"`
	Duration    time.Duration `json:"duration"`
	ErrorMsg    string        `json:"error_msg,omitempty"`
}

type ScenarioStep struct {
	StepName    string `json:"step_name"`
	Action      string `json:"action"`
	URL         string `json:"url"`
	Method      string `json:"method"`
	Expected    string `json:"expected"`
	Actual      string `json:"actual"`
	Success     bool   `json:"success"`
	Duration    time.Duration `json:"duration"`
}

type Customer struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewE2ETester(baseURL string) *E2ETester {
	return &E2ETester{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		scenarios: make([]TestScenario, 0),
	}
}

func (e *E2ETester) RunTests(ctx context.Context) error {
	log.Info().Msg("Starting E2E tests for VacuumMart shopping flows")
	
	// Test scenarios
	scenarios := []func(context.Context) error{
		e.testCompleteShoppingFlow,
		e.testGuestCheckoutFlow,
		e.testProductSearchToPurchase,
		e.testCartAbandonmentRecovery,
		e.testMultipleProductPurchase,
		e.testInventoryConstraints,
		e.testUserAccountManagement,
		e.testOrderTrackingFlow,
		e.testProductReviewFlow,
		e.testCompareProductsFlow,
	}
	
	for i, scenario := range scenarios {
		scenarioName := fmt.Sprintf("Scenario %d", i+1)
		start := time.Now()
		
		log.Info().Str("scenario", scenarioName).Msg("Starting E2E scenario")
		
		if err := scenario(ctx); err != nil {
			log.Error().Err(err).Str("scenario", scenarioName).Msg("E2E scenario failed")
			// Continue with other scenarios instead of failing completely
		}
		
		duration := time.Since(start)
		log.Info().Str("scenario", scenarioName).Dur("duration", duration).Msg("E2E scenario completed")
	}
	
	e.generateReport()
	return nil
}

func (e *E2ETester) testCompleteShoppingFlow(ctx context.Context) error {
	scenario := TestScenario{
		Name:        "Complete Shopping Flow",
		Description: "User browses products, adds items to cart, and completes checkout",
		Steps:       make([]ScenarioStep, 0),
	}
	start := time.Now()
	
	// Step 1: Visit homepage
	step1 := e.executeStep(ctx, "Visit Homepage", "GET", "/", "Load homepage successfully", "200 OK")
	scenario.Steps = append(scenario.Steps, step1)
	
	// Step 2: Browse product catalog
	step2 := e.executeStep(ctx, "Browse Products", "GET", "/api/products", "Load product catalog", "200 OK")
	scenario.Steps = append(scenario.Steps, step2)
	
	// Step 3: Search for specific product
	step3 := e.executeStep(ctx, "Search for Dyson", "GET", "/api/products/search?q=dyson", "Find Dyson products", "200 OK")
	scenario.Steps = append(scenario.Steps, step3)
	
	// Step 4: View product details
	step4 := e.executeStep(ctx, "View Product Details", "GET", "/api/products/1", "Load product details", "200 OK")
	scenario.Steps = append(scenario.Steps, step4)
	
	// Step 5: Add product to cart
	step5 := e.executeStep(ctx, "Add to Cart", "POST", "/api/cart/items", "Add product to cart", "200 OK")
	scenario.Steps = append(scenario.Steps, step5)
	
	// Step 6: View cart
	step6 := e.executeStep(ctx, "View Cart", "GET", "/api/cart", "View cart contents", "200 OK")
	scenario.Steps = append(scenario.Steps, step6)
	
	// Step 7: Update quantity
	step7 := e.executeStep(ctx, "Update Quantity", "PUT", "/api/cart/items/1", "Update item quantity", "200 OK")
	scenario.Steps = append(scenario.Steps, step7)
	
	// Step 8: Proceed to checkout
	step8 := e.executeStep(ctx, "Proceed to Checkout", "GET", "/checkout", "Load checkout page", "200 OK")
	scenario.Steps = append(scenario.Steps, step8)
	
	// Step 9: Create order
	step9 := e.executeStep(ctx, "Create Order", "POST", "/api/orders", "Create order successfully", "201 Created")
	scenario.Steps = append(scenario.Steps, step9)
	
	// Step 10: Confirm order
	step10 := e.executeStep(ctx, "Confirm Order", "GET", "/api/orders/1", "View order confirmation", "200 OK")
	scenario.Steps = append(scenario.Steps, step10)
	
	scenario.Duration = time.Since(start)
	scenario.Success = e.allStepsSuccessful(scenario.Steps)
	e.scenarios = append(e.scenarios, scenario)
	
	return nil
}

func (e *E2ETester) testGuestCheckoutFlow(ctx context.Context) error {
	scenario := TestScenario{
		Name:        "Guest Checkout Flow",
		Description: "Guest user completes purchase without creating account",
		Steps:       make([]ScenarioStep, 0),
	}
	start := time.Now()
	
	// Guest checkout specific steps
	step1 := e.executeStep(ctx, "Browse as Guest", "GET", "/api/products", "Load products as guest", "200 OK")
	scenario.Steps = append(scenario.Steps, step1)
	
	step2 := e.executeStep(ctx, "Add to Cart (Guest)", "POST", "/api/cart/items", "Add item as guest", "200 OK")
	scenario.Steps = append(scenario.Steps, step2)
	
	step3 := e.executeStep(ctx, "Guest Checkout", "POST", "/api/checkout/guest", "Complete guest checkout", "200 OK")
	scenario.Steps = append(scenario.Steps, step3)
	
	scenario.Duration = time.Since(start)
	scenario.Success = e.allStepsSuccessful(scenario.Steps)
	e.scenarios = append(e.scenarios, scenario)
	
	return nil
}

func (e *E2ETester) testProductSearchToPurchase(ctx context.Context) error {
	scenario := TestScenario{
		Name:        "Product Search to Purchase",
		Description: "User searches for specific vacuum type and completes purchase",
		Steps:       make([]ScenarioStep, 0),
	}
	start := time.Now()
	
	searchTerms := []string{"robot vacuum", "pet hair", "cordless"}
	
	for _, term := range searchTerms {
		stepName := fmt.Sprintf("Search for '%s'", term)
		url := fmt.Sprintf("/api/products/search?q=%s", strings.ReplaceAll(term, " ", "+"))
		step := e.executeStep(ctx, stepName, "GET", url, "Find relevant products", "200 OK")
		scenario.Steps = append(scenario.Steps, step)
	}
	
	// Filter by category
	step := e.executeStep(ctx, "Filter by Robot Category", "GET", "/api/products?category=robot", "Filter robot vacuums", "200 OK")
	scenario.Steps = append(scenario.Steps, step)
	
	// Sort by price
	step = e.executeStep(ctx, "Sort by Price", "GET", "/api/products?category=robot&sort=price_asc", "Sort by price", "200 OK")
	scenario.Steps = append(scenario.Steps, step)
	
	scenario.Duration = time.Since(start)
	scenario.Success = e.allStepsSuccessful(scenario.Steps)
	e.scenarios = append(e.scenarios, scenario)
	
	return nil
}

func (e *E2ETester) testCartAbandonmentRecovery(ctx context.Context) error {
	scenario := TestScenario{
		Name:        "Cart Abandonment Recovery",
		Description: "Test cart persistence and recovery mechanisms",
		Steps:       make([]ScenarioStep, 0),
	}
	start := time.Now()
	
	// Add items to cart
	step1 := e.executeStep(ctx, "Add Items to Cart", "POST", "/api/cart/items", "Add multiple items", "200 OK")
	scenario.Steps = append(scenario.Steps, step1)
	
	// Simulate leaving site (save cart state)
	step2 := e.executeStep(ctx, "Save Cart State", "POST", "/api/cart/save", "Save cart for later", "200 OK")
	scenario.Steps = append(scenario.Steps, step2)
	
	// Return and restore cart
	step3 := e.executeStep(ctx, "Restore Cart", "GET", "/api/cart/restore", "Restore saved cart", "200 OK")
	scenario.Steps = append(scenario.Steps, step3)
	
	scenario.Duration = time.Since(start)
	scenario.Success = e.allStepsSuccessful(scenario.Steps)
	e.scenarios = append(e.scenarios, scenario)
	
	return nil
}

func (e *E2ETester) testMultipleProductPurchase(ctx context.Context) error {
	scenario := TestScenario{
		Name:        "Multiple Product Purchase",
		Description: "Purchase multiple different vacuum cleaners in one order",
		Steps:       make([]ScenarioStep, 0),
	}
	start := time.Now()
	
	productIDs := []int{1, 2, 3, 5} // Different vacuum products
	
	for _, id := range productIDs {
		stepName := fmt.Sprintf("Add Product %d to Cart", id)
		url := "/api/cart/items"
		step := e.executeStep(ctx, stepName, "POST", url, "Add product to cart", "200 OK")
		scenario.Steps = append(scenario.Steps, step)
	}
	
	// Calculate total
	step := e.executeStep(ctx, "Calculate Total", "GET", "/api/cart/total", "Calculate cart total", "200 OK")
	scenario.Steps = append(scenario.Steps, step)
	
	// Apply discount code
	step = e.executeStep(ctx, "Apply Discount", "POST", "/api/cart/discount", "Apply discount code", "200 OK")
	scenario.Steps = append(scenario.Steps, step)
	
	scenario.Duration = time.Since(start)
	scenario.Success = e.allStepsSuccessful(scenario.Steps)
	e.scenarios = append(e.scenarios, scenario)
	
	return nil
}

func (e *E2ETester) testInventoryConstraints(ctx context.Context) error {
	scenario := TestScenario{
		Name:        "Inventory Constraints",
		Description: "Test behavior when products are out of stock or low inventory",
		Steps:       make([]ScenarioStep, 0),
	}
	start := time.Now()
	
	// Try to add out-of-stock item
	step1 := e.executeStep(ctx, "Add Out-of-Stock Item", "POST", "/api/cart/items", "Should fail gracefully", "400 Bad Request")
	scenario.Steps = append(scenario.Steps, step1)
	
	// Check inventory status
	step2 := e.executeStep(ctx, "Check Inventory", "GET", "/api/products/1/inventory", "Check inventory status", "200 OK")
	scenario.Steps = append(scenario.Steps, step2)
	
	// Test low stock warning
	step3 := e.executeStep(ctx, "Low Stock Warning", "GET", "/api/products?low_stock=true", "Show low stock items", "200 OK")
	scenario.Steps = append(scenario.Steps, step3)
	
	scenario.Duration = time.Since(start)
	scenario.Success = e.allStepsSuccessful(scenario.Steps)
	e.scenarios = append(e.scenarios, scenario)
	
	return nil
}

func (e *E2ETester) testUserAccountManagement(ctx context.Context) error {
	scenario := TestScenario{
		Name:        "User Account Management",
		Description: "Test user registration, login, and profile management",
		Steps:       make([]ScenarioStep, 0),
	}
	start := time.Now()
	
	// Register new user
	step1 := e.executeStep(ctx, "User Registration", "POST", "/api/users/register", "Create new account", "201 Created")
	scenario.Steps = append(scenario.Steps, step1)
	
	// User login
	step2 := e.executeStep(ctx, "User Login", "POST", "/api/users/login", "Login successfully", "200 OK")
	scenario.Steps = append(scenario.Steps, step2)
	
	// View profile
	step3 := e.executeStep(ctx, "View Profile", "GET", "/api/users/profile", "Load user profile", "200 OK")
	scenario.Steps = append(scenario.Steps, step3)
	
	// Update profile
	step4 := e.executeStep(ctx, "Update Profile", "PUT", "/api/users/profile", "Update user info", "200 OK")
	scenario.Steps = append(scenario.Steps, step4)
	
	scenario.Duration = time.Since(start)
	scenario.Success = e.allStepsSuccessful(scenario.Steps)
	e.scenarios = append(e.scenarios, scenario)
	
	return nil
}

func (e *E2ETester) testOrderTrackingFlow(ctx context.Context) error {
	scenario := TestScenario{
		Name:        "Order Tracking Flow",
		Description: "Test order creation and tracking functionality",
		Steps:       make([]ScenarioStep, 0),
	}
	start := time.Now()
	
	// Create order
	step1 := e.executeStep(ctx, "Create Order", "POST", "/api/orders", "Create new order", "201 Created")
	scenario.Steps = append(scenario.Steps, step1)
	
	// Track order status
	step2 := e.executeStep(ctx, "Track Order", "GET", "/api/orders/1/tracking", "Get tracking info", "200 OK")
	scenario.Steps = append(scenario.Steps, step2)
	
	// Update order status (admin action)
	step3 := e.executeStep(ctx, "Update Status", "PUT", "/api/admin/orders/1/status", "Update order status", "200 OK")
	scenario.Steps = append(scenario.Steps, step3)
	
	// Check status history
	step4 := e.executeStep(ctx, "Status History", "GET", "/api/orders/1/history", "View status history", "200 OK")
	scenario.Steps = append(scenario.Steps, step4)
	
	scenario.Duration = time.Since(start)
	scenario.Success = e.allStepsSuccessful(scenario.Steps)
	e.scenarios = append(e.scenarios, scenario)
	
	return nil
}

func (e *E2ETester) testProductReviewFlow(ctx context.Context) error {
	scenario := TestScenario{
		Name:        "Product Review Flow",
		Description: "Test customer product review and rating functionality",
		Steps:       make([]ScenarioStep, 0),
	}
	start := time.Now()
	
	// View product reviews
	step1 := e.executeStep(ctx, "View Reviews", "GET", "/api/products/1/reviews", "Load product reviews", "200 OK")
	scenario.Steps = append(scenario.Steps, step1)
	
	// Add product review
	step2 := e.executeStep(ctx, "Add Review", "POST", "/api/products/1/reviews", "Submit review", "201 Created")
	scenario.Steps = append(scenario.Steps, step2)
	
	// Update review
	step3 := e.executeStep(ctx, "Update Review", "PUT", "/api/reviews/1", "Update review", "200 OK")
	scenario.Steps = append(scenario.Steps, step3)
	
	scenario.Duration = time.Since(start)
	scenario.Success = e.allStepsSuccessful(scenario.Steps)
	e.scenarios = append(e.scenarios, scenario)
	
	return nil
}

func (e *E2ETester) testCompareProductsFlow(ctx context.Context) error {
	scenario := TestScenario{
		Name:        "Compare Products Flow",
		Description: "Test product comparison functionality",
		Steps:       make([]ScenarioStep, 0),
	}
	start := time.Now()
	
	// Add products to comparison
	step1 := e.executeStep(ctx, "Add to Compare", "POST", "/api/compare/products", "Add products to compare", "200 OK")
	scenario.Steps = append(scenario.Steps, step1)
	
	// View comparison
	step2 := e.executeStep(ctx, "View Comparison", "GET", "/api/compare", "View product comparison", "200 OK")
	scenario.Steps = append(scenario.Steps, step2)
	
	// Remove from comparison
	step3 := e.executeStep(ctx, "Remove from Compare", "DELETE", "/api/compare/products/1", "Remove product", "200 OK")
	scenario.Steps = append(scenario.Steps, step3)
	
	scenario.Duration = time.Since(start)
	scenario.Success = e.allStepsSuccessful(scenario.Steps)
	e.scenarios = append(e.scenarios, scenario)
	
	return nil
}

func (e *E2ETester) executeStep(ctx context.Context, stepName, method, endpoint, expected, actualExpected string) ScenarioStep {
	start := time.Now()
	
	step := ScenarioStep{
		StepName: stepName,
		Action:   fmt.Sprintf("%s %s", method, endpoint),
		URL:      endpoint,
		Method:   method,
		Expected: expected,
		Success:  false,
	}
	
	// Make HTTP request
	url := e.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		step.Actual = fmt.Sprintf("Error creating request: %v", err)
		step.Duration = time.Since(start)
		return step
	}
	
	resp, err := e.httpClient.Do(req)
	if err != nil {
		step.Actual = fmt.Sprintf("Request failed: %v", err)
		step.Duration = time.Since(start)
		return step
	}
	defer resp.Body.Close()
	
	step.Actual = fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	step.Duration = time.Since(start)
	
	// Check if response matches expected (simplified check)
	if strings.Contains(actualExpected, fmt.Sprintf("%d", resp.StatusCode)) {
		step.Success = true
	}
	
	log.Info().
		Str("step", stepName).
		Str("method", method).
		Str("endpoint", endpoint).
		Str("expected", expected).
		Str("actual", step.Actual).
		Bool("success", step.Success).
		Dur("duration", step.Duration).
		Msg("E2E step completed")
	
	return step
}

func (e *E2ETester) allStepsSuccessful(steps []ScenarioStep) bool {
	for _, step := range steps {
		if !step.Success {
			return false
		}
	}
	return true
}

func (e *E2ETester) generateReport() {
	log.Info().Msg("Generating E2E test report")
	
	totalScenarios := len(e.scenarios)
	passedScenarios := 0
	failedScenarios := 0
	totalDuration := time.Duration(0)
	totalSteps := 0
	
	for _, scenario := range e.scenarios {
		if scenario.Success {
			passedScenarios++
		} else {
			failedScenarios++
		}
		totalDuration += scenario.Duration
		totalSteps += len(scenario.Steps)
	}
	
	log.Info().
		Int("total_scenarios", totalScenarios).
		Int("passed", passedScenarios).
		Int("failed", failedScenarios).
		Int("total_steps", totalSteps).
		Dur("total_duration", totalDuration).
		Dur("avg_scenario_duration", totalDuration/time.Duration(totalScenarios)).
		Msg("E2E test summary")
	
	// Save detailed results
	reportData := map[string]interface{}{
		"summary": map[string]interface{}{
			"total_scenarios":      totalScenarios,
			"passed_scenarios":     passedScenarios,
			"failed_scenarios":     failedScenarios,
			"success_rate":         float64(passedScenarios) / float64(totalScenarios) * 100,
			"total_steps":          totalSteps,
			"total_duration_ms":    totalDuration.Milliseconds(),
			"avg_duration_ms":      totalDuration.Milliseconds() / int64(totalScenarios),
		},
		"detailed_scenarios": e.scenarios,
		"test_timestamp":     time.Now().Format(time.RFC3339),
	}
	
	if jsonData, err := json.MarshalIndent(reportData, "", "  "); err == nil {
		log.Info().Str("report", string(jsonData)).Msg("Detailed E2E test report")
	}
}
