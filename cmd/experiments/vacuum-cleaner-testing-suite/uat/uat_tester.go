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

type UATTester struct {
	baseURL    string
	httpClient *http.Client
	scenarios  []UATScenario
}

type UATScenario struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	BusinessValue   string          `json:"business_value"`
	UserStory       string          `json:"user_story"`
	AcceptanceCriteria []string     `json:"acceptance_criteria"`
	TestSteps       []UATStep       `json:"test_steps"`
	Success         bool            `json:"success"`
	Duration        time.Duration   `json:"duration"`
	ErrorMsg        string          `json:"error_msg,omitempty"`
	Timestamp       time.Time       `json:"timestamp"`
}

type UATStep struct {
	StepNumber  int           `json:"step_number"`
	Action      string        `json:"action"`
	Expected    string        `json:"expected"`
	Actual      string        `json:"actual"`
	Passed      bool          `json:"passed"`
	Duration    time.Duration `json:"duration"`
	Screenshots []string      `json:"screenshots,omitempty"`
}

func NewUATTester(baseURL string) *UATTester {
	return &UATTester{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		scenarios: make([]UATScenario, 0),
	}
}

func (u *UATTester) RunTests(ctx context.Context) error {
	log.Info().Msg("Starting User Acceptance Tests based on VacuumMart business requirements")
	
	// Define UAT scenarios based on business requirements
	scenarios := []func(context.Context) UATScenario{
		u.testCustomerProductBrowsing,
		u.testVacuumSearchAndFiltering,
		u.testShoppingCartFunctionality,
		u.testCheckoutProcess,
		u.testUserAccountManagement,
		u.testProductComparison,
		u.testInventoryManagement,
		u.testOrderTracking,
		u.testCustomerReviews,
		u.testAdminDashboard,
		u.testSalesReporting,
		u.testMobileShoppingExperience,
	}
	
	for _, scenarioFunc := range scenarios {
		start := time.Now()
		scenario := scenarioFunc(ctx)
		scenario.Duration = time.Since(start)
		scenario.Timestamp = start
		
		u.scenarios = append(u.scenarios, scenario)
		
		log.Info().
			Str("scenario", scenario.ID).
			Str("name", scenario.Name).
			Bool("success", scenario.Success).
			Dur("duration", scenario.Duration).
			Msg("UAT scenario completed")
	}
	
	u.generateUATReport()
	return nil
}

func (u *UATTester) testCustomerProductBrowsing(ctx context.Context) UATScenario {
	scenario := UATScenario{
		ID:          "UAT-001",
		Name:        "Customer Product Browsing",
		Description: "Customer can browse vacuum cleaners by category and brand",
		BusinessValue: "Customers can easily find vacuum cleaners that meet their needs",
		UserStory:   "As a customer, I want to browse vacuum cleaners by category so that I can find the right type for my home",
		AcceptanceCriteria: []string{
			"Customer can view all vacuum categories (upright, canister, robot, handheld)",
			"Customer can filter by brand (Dyson, Hoover, Roomba, Shark)",
			"Product listings show key information (name, price, rating, image)",
			"Customer can sort products by price, rating, and name",
			"Product pages load within 3 seconds",
		},
		TestSteps: make([]UATStep, 0),
	}
	
	// Step 1: Access product catalog
	step1 := u.executeUATStep(ctx, 1, "Navigate to product catalog", "GET", "/api/products", "Product catalog loads successfully")
	scenario.TestSteps = append(scenario.TestSteps, step1)
	
	// Step 2: Browse by category
	categories := []string{"upright", "canister", "robot", "handheld"}
	for i, category := range categories {
		stepNum := i + 2
		action := fmt.Sprintf("Browse %s vacuum category", category)
		endpoint := fmt.Sprintf("/api/products?category=%s", category)
		expected := fmt.Sprintf("%s vacuums displayed", category)
		step := u.executeUATStep(ctx, stepNum, action, "GET", endpoint, expected)
		scenario.TestSteps = append(scenario.TestSteps, step)
	}
	
	// Step 6: Filter by brand
	brands := []string{"Dyson", "Hoover", "Roomba"}
	for i, brand := range brands {
		stepNum := i + 6
		action := fmt.Sprintf("Filter by %s brand", brand)
		endpoint := fmt.Sprintf("/api/products?brand=%s", brand)
		expected := fmt.Sprintf("%s products displayed", brand)
		step := u.executeUATStep(ctx, stepNum, action, "GET", endpoint, expected)
		scenario.TestSteps = append(scenario.TestSteps, step)
	}
	
	// Step 9: Sort products
	step9 := u.executeUATStep(ctx, 9, "Sort by price ascending", "GET", "/api/products?sort=price_asc", "Products sorted by price low to high")
	scenario.TestSteps = append(scenario.TestSteps, step9)
	
	scenario.Success = u.allUATStepsSuccessful(scenario.TestSteps)
	return scenario
}

func (u *UATTester) testVacuumSearchAndFiltering(ctx context.Context) UATScenario {
	scenario := UATScenario{
		ID:          "UAT-002",
		Name:        "Vacuum Search and Filtering",
		Description: "Customer can search for specific vacuum features and filter results",
		BusinessValue: "Customers can quickly find vacuums with specific features they need",
		UserStory:   "As a customer, I want to search for vacuums with specific features so that I can find one that meets my cleaning needs",
		AcceptanceCriteria: []string{
			"Customer can search by keywords (pet hair, cordless, lightweight)",
			"Search results are relevant and accurate",
			"Customer can combine multiple filters",
			"Price range filtering works correctly",
			"Search returns results within 2 seconds",
		},
		TestSteps: make([]UATStep, 0),
	}
	
	// Test keyword searches
	searchTerms := []string{"pet hair", "cordless", "lightweight", "robot"}
	for i, term := range searchTerms {
		stepNum := i + 1
		action := fmt.Sprintf("Search for '%s' vacuums", term)
		endpoint := fmt.Sprintf("/api/products/search?q=%s", strings.ReplaceAll(term, " ", "+"))
		expected := fmt.Sprintf("Relevant results for '%s'", term)
		step := u.executeUATStep(ctx, stepNum, action, "GET", endpoint, expected)
		scenario.TestSteps = append(scenario.TestSteps, step)
	}
	
	// Test price filtering
	step5 := u.executeUATStep(ctx, 5, "Filter by price range $100-$500", "GET", "/api/products?min_price=100&max_price=500", "Products in price range displayed")
	scenario.TestSteps = append(scenario.TestSteps, step5)
	
	// Test combined filtering
	step6 := u.executeUATStep(ctx, 6, "Combine category and brand filter", "GET", "/api/products?category=robot&brand=Roomba", "Robot vacuums from Roomba displayed")
	scenario.TestSteps = append(scenario.TestSteps, step6)
	
	scenario.Success = u.allUATStepsSuccessful(scenario.TestSteps)
	return scenario
}

func (u *UATTester) testShoppingCartFunctionality(ctx context.Context) UATScenario {
	scenario := UATScenario{
		ID:          "UAT-003",
		Name:        "Shopping Cart Functionality",
		Description: "Customer can add, modify, and remove items from shopping cart",
		BusinessValue: "Customers can manage their purchases before checkout",
		UserStory:   "As a customer, I want to add items to my cart and modify quantities so that I can purchase multiple items",
		AcceptanceCriteria: []string{
			"Customer can add products to cart",
			"Customer can view cart contents",
			"Customer can update item quantities",
			"Customer can remove items from cart",
			"Cart total is calculated correctly",
			"Cart persists during browsing session",
		},
		TestSteps: make([]UATStep, 0),
	}
	
	// Test adding to cart
	step1 := u.executeUATStep(ctx, 1, "Add Dyson V15 to cart", "POST", "/api/cart/items", "Item added to cart successfully")
	scenario.TestSteps = append(scenario.TestSteps, step1)
	
	// Test viewing cart
	step2 := u.executeUATStep(ctx, 2, "View cart contents", "GET", "/api/cart", "Cart shows added item")
	scenario.TestSteps = append(scenario.TestSteps, step2)
	
	// Test updating quantity
	step3 := u.executeUATStep(ctx, 3, "Update quantity to 2", "PUT", "/api/cart/items/1", "Quantity updated to 2")
	scenario.TestSteps = append(scenario.TestSteps, step3)
	
	// Test adding another item
	step4 := u.executeUATStep(ctx, 4, "Add Roomba to cart", "POST", "/api/cart/items", "Second item added to cart")
	scenario.TestSteps = append(scenario.TestSteps, step4)
	
	// Test removing item
	step5 := u.executeUATStep(ctx, 5, "Remove first item", "DELETE", "/api/cart/items/1", "Item removed from cart")
	scenario.TestSteps = append(scenario.TestSteps, step5)
	
	scenario.Success = u.allUATStepsSuccessful(scenario.TestSteps)
	return scenario
}

func (u *UATTester) testCheckoutProcess(ctx context.Context) UATScenario {
	scenario := UATScenario{
		ID:          "UAT-004",
		Name:        "Checkout Process",
		Description: "Customer can complete purchase through checkout process",
		BusinessValue: "Customers can successfully purchase vacuum cleaners",
		UserStory:   "As a customer, I want to checkout and pay for my items so that I can receive my vacuum cleaner",
		AcceptanceCriteria: []string{
			"Customer can proceed to checkout from cart",
			"Customer can enter shipping information",
			"Customer can select payment method",
			"Order total is displayed correctly",
			"Customer receives order confirmation",
			"Order is recorded in system",
		},
		TestSteps: make([]UATStep, 0),
	}
	
	// Test checkout initiation
	step1 := u.executeUATStep(ctx, 1, "Proceed to checkout", "GET", "/checkout", "Checkout page loads")
	scenario.TestSteps = append(scenario.TestSteps, step1)
	
	// Test order creation
	step2 := u.executeUATStep(ctx, 2, "Create order", "POST", "/api/orders", "Order created successfully")
	scenario.TestSteps = append(scenario.TestSteps, step2)
	
	// Test order confirmation
	step3 := u.executeUATStep(ctx, 3, "View order confirmation", "GET", "/api/orders/1", "Order confirmation displayed")
	scenario.TestSteps = append(scenario.TestSteps, step3)
	
	scenario.Success = u.allUATStepsSuccessful(scenario.TestSteps)
	return scenario
}

func (u *UATTester) testUserAccountManagement(ctx context.Context) UATScenario {
	scenario := UATScenario{
		ID:          "UAT-005",
		Name:        "User Account Management",
		Description: "Customer can create and manage their account",
		BusinessValue: "Customers can track orders and save preferences",
		UserStory:   "As a customer, I want to create an account so that I can track my orders and save my preferences",
		AcceptanceCriteria: []string{
			"Customer can register for new account",
			"Customer can login to existing account",
			"Customer can view order history",
			"Customer can update profile information",
			"Password requirements are enforced",
		},
		TestSteps: make([]UATStep, 0),
	}
	
	// Test user registration
	step1 := u.executeUATStep(ctx, 1, "Register new account", "POST", "/api/users/register", "Account created successfully")
	scenario.TestSteps = append(scenario.TestSteps, step1)
	
	// Test user login
	step2 := u.executeUATStep(ctx, 2, "Login to account", "POST", "/api/users/login", "Login successful")
	scenario.TestSteps = append(scenario.TestSteps, step2)
	
	// Test profile viewing
	step3 := u.executeUATStep(ctx, 3, "View user profile", "GET", "/api/users/profile", "Profile information displayed")
	scenario.TestSteps = append(scenario.TestSteps, step3)
	
	scenario.Success = u.allUATStepsSuccessful(scenario.TestSteps)
	return scenario
}

func (u *UATTester) testProductComparison(ctx context.Context) UATScenario {
	scenario := UATScenario{
		ID:          "UAT-006",
		Name:        "Product Comparison",
		Description: "Customer can compare vacuum cleaner features side-by-side",
		BusinessValue: "Customers can make informed purchase decisions",
		UserStory:   "As a customer, I want to compare vacuum cleaners side-by-side so that I can choose the best one for my needs",
		AcceptanceCriteria: []string{
			"Customer can add products to comparison",
			"Comparison shows key features side-by-side",
			"Customer can compare up to 3 products",
			"Customer can remove products from comparison",
			"Comparison is visually clear and helpful",
		},
		TestSteps: make([]UATStep, 0),
	}
	
	// Test adding to comparison
	step1 := u.executeUATStep(ctx, 1, "Add products to comparison", "POST", "/api/compare/products", "Products added to comparison")
	scenario.TestSteps = append(scenario.TestSteps, step1)
	
	// Test viewing comparison
	step2 := u.executeUATStep(ctx, 2, "View product comparison", "GET", "/api/compare", "Comparison table displayed")
	scenario.TestSteps = append(scenario.TestSteps, step2)
	
	scenario.Success = u.allUATStepsSuccessful(scenario.TestSteps)
	return scenario
}

func (u *UATTester) testInventoryManagement(ctx context.Context) UATScenario {
	scenario := UATScenario{
		ID:          "UAT-007",
		Name:        "Inventory Management",
		Description: "Admin can manage product inventory levels",
		BusinessValue: "Store can maintain accurate stock levels",
		UserStory:   "As a store manager, I want to manage inventory so that customers see accurate stock information",
		AcceptanceCriteria: []string{
			"Admin can view current inventory levels",
			"Admin can update stock quantities",
			"Out-of-stock products are marked clearly",
			"Low stock warnings are displayed",
			"Inventory changes are logged",
		},
		TestSteps: make([]UATStep, 0),
	}
	
	// Test inventory viewing
	step1 := u.executeUATStep(ctx, 1, "View inventory status", "GET", "/api/admin/inventory", "Inventory levels displayed")
	scenario.TestSteps = append(scenario.TestSteps, step1)
	
	// Test inventory update
	step2 := u.executeUATStep(ctx, 2, "Update product inventory", "PUT", "/api/admin/products/1/inventory", "Inventory updated successfully")
	scenario.TestSteps = append(scenario.TestSteps, step2)
	
	scenario.Success = u.allUATStepsSuccessful(scenario.TestSteps)
	return scenario
}

func (u *UATTester) testOrderTracking(ctx context.Context) UATScenario {
	scenario := UATScenario{
		ID:          "UAT-008",
		Name:        "Order Tracking",
		Description: "Customer can track order status and history",
		BusinessValue: "Customers stay informed about their purchase status",
		UserStory:   "As a customer, I want to track my order so that I know when to expect delivery",
		AcceptanceCriteria: []string{
			"Customer can view order status",
			"Status updates are timely and accurate",
			"Customer can view order history",
			"Tracking information is clear",
			"Customer receives status notifications",
		},
		TestSteps: make([]UATStep, 0),
	}
	
	// Test order tracking
	step1 := u.executeUATStep(ctx, 1, "View order tracking", "GET", "/api/orders/1/tracking", "Tracking information displayed")
	scenario.TestSteps = append(scenario.TestSteps, step1)
	
	// Test order history
	step2 := u.executeUATStep(ctx, 2, "View order history", "GET", "/api/orders?customer_id=1", "Order history displayed")
	scenario.TestSteps = append(scenario.TestSteps, step2)
	
	scenario.Success = u.allUATStepsSuccessful(scenario.TestSteps)
	return scenario
}

func (u *UATTester) testCustomerReviews(ctx context.Context) UATScenario {
	scenario := UATScenario{
		ID:          "UAT-009",
		Name:        "Customer Reviews",
		Description: "Customer can read and write product reviews",
		BusinessValue: "Reviews help customers make purchase decisions",
		UserStory:   "As a customer, I want to read reviews so that I can learn from other customers' experiences",
		AcceptanceCriteria: []string{
			"Customer can view product reviews",
			"Customer can write reviews for purchased products",
			"Reviews show rating and detailed feedback",
			"Reviews are moderated appropriately",
			"Review summary is helpful",
		},
		TestSteps: make([]UATStep, 0),
	}
	
	// Test viewing reviews
	step1 := u.executeUATStep(ctx, 1, "View product reviews", "GET", "/api/products/1/reviews", "Reviews displayed")
	scenario.TestSteps = append(scenario.TestSteps, step1)
	
	// Test adding review
	step2 := u.executeUATStep(ctx, 2, "Add product review", "POST", "/api/products/1/reviews", "Review submitted successfully")
	scenario.TestSteps = append(scenario.TestSteps, step2)
	
	scenario.Success = u.allUATStepsSuccessful(scenario.TestSteps)
	return scenario
}

func (u *UATTester) testAdminDashboard(ctx context.Context) UATScenario {
	scenario := UATScenario{
		ID:          "UAT-010",
		Name:        "Admin Dashboard",
		Description: "Admin can access comprehensive management dashboard",
		BusinessValue: "Store management can monitor business operations",
		UserStory:   "As a store manager, I want a dashboard so that I can monitor sales and operations",
		AcceptanceCriteria: []string{
			"Admin can access dashboard",
			"Dashboard shows key metrics",
			"Admin can manage orders",
			"Admin can manage products",
			"Dashboard is responsive and fast",
		},
		TestSteps: make([]UATStep, 0),
	}
	
	// Test dashboard access
	step1 := u.executeUATStep(ctx, 1, "Access admin dashboard", "GET", "/admin", "Dashboard loads successfully")
	scenario.TestSteps = append(scenario.TestSteps, step1)
	
	// Test order management
	step2 := u.executeUATStep(ctx, 2, "View all orders", "GET", "/api/admin/orders", "Orders list displayed")
	scenario.TestSteps = append(scenario.TestSteps, step2)
	
	scenario.Success = u.allUATStepsSuccessful(scenario.TestSteps)
	return scenario
}

func (u *UATTester) testSalesReporting(ctx context.Context) UATScenario {
	scenario := UATScenario{
		ID:          "UAT-011",
		Name:        "Sales Reporting",
		Description: "Admin can generate sales reports and analytics",
		BusinessValue: "Business can track performance and make data-driven decisions",
		UserStory:   "As a business owner, I want sales reports so that I can understand business performance",
		AcceptanceCriteria: []string{
			"Admin can view sales analytics",
			"Reports show revenue and unit sales",
			"Data can be filtered by date range",
			"Top-selling products are highlighted",
			"Reports are accurate and timely",
		},
		TestSteps: make([]UATStep, 0),
	}
	
	// Test sales analytics
	step1 := u.executeUATStep(ctx, 1, "View sales analytics", "GET", "/api/admin/analytics/sales", "Sales data displayed")
	scenario.TestSteps = append(scenario.TestSteps, step1)
	
	scenario.Success = u.allUATStepsSuccessful(scenario.TestSteps)
	return scenario
}

func (u *UATTester) testMobileShoppingExperience(ctx context.Context) UATScenario {
	scenario := UATScenario{
		ID:          "UAT-012",
		Name:        "Mobile Shopping Experience",
		Description: "Customer can shop effectively on mobile devices",
		BusinessValue: "Mobile customers can complete purchases easily",
		UserStory:   "As a mobile customer, I want to browse and purchase vacuums on my phone so that I can shop anywhere",
		AcceptanceCriteria: []string{
			"Site is responsive on mobile devices",
			"Mobile navigation is user-friendly",
			"Touch interactions work properly",
			"Mobile checkout process is smooth",
			"Load times are acceptable on mobile",
		},
		TestSteps: make([]UATStep, 0),
	}
	
	// Test mobile homepage
	step1 := u.executeUATStepWithMobile(ctx, 1, "Access homepage on mobile", "GET", "/", "Mobile homepage loads properly")
	scenario.TestSteps = append(scenario.TestSteps, step1)
	
	// Test mobile product browsing
	step2 := u.executeUATStepWithMobile(ctx, 2, "Browse products on mobile", "GET", "/api/products", "Products display well on mobile")
	scenario.TestSteps = append(scenario.TestSteps, step2)
	
	scenario.Success = u.allUATStepsSuccessful(scenario.TestSteps)
	return scenario
}

func (u *UATTester) executeUATStep(ctx context.Context, stepNumber int, action, method, endpoint, expected string) UATStep {
	start := time.Now()
	
	step := UATStep{
		StepNumber: stepNumber,
		Action:     action,
		Expected:   expected,
		Passed:     false,
	}
	
	// Make HTTP request
	url := u.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		step.Actual = fmt.Sprintf("Error creating request: %v", err)
		step.Duration = time.Since(start)
		return step
	}
	
	resp, err := u.httpClient.Do(req)
	if err != nil {
		step.Actual = fmt.Sprintf("Request failed: %v", err)
		step.Duration = time.Since(start)
		return step
	}
	defer resp.Body.Close()
	
	step.Actual = fmt.Sprintf("HTTP %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	step.Duration = time.Since(start)
	
	// Simple success criteria: 2xx status codes
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		step.Passed = true
	}
	
	log.Info().
		Int("step", stepNumber).
		Str("action", action).
		Str("expected", expected).
		Str("actual", step.Actual).
		Bool("passed", step.Passed).
		Dur("duration", step.Duration).
		Msg("UAT step completed")
	
	return step
}

func (u *UATTester) executeUATStepWithMobile(ctx context.Context, stepNumber int, action, method, endpoint, expected string) UATStep {
	start := time.Now()
	
	step := UATStep{
		StepNumber: stepNumber,
		Action:     action,
		Expected:   expected,
		Passed:     false,
	}
	
	// Make HTTP request with mobile user agent
	url := u.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		step.Actual = fmt.Sprintf("Error creating request: %v", err)
		step.Duration = time.Since(start)
		return step
	}
	
	// Set mobile user agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1")
	
	resp, err := u.httpClient.Do(req)
	if err != nil {
		step.Actual = fmt.Sprintf("Request failed: %v", err)
		step.Duration = time.Since(start)
		return step
	}
	defer resp.Body.Close()
	
	step.Actual = fmt.Sprintf("HTTP %d %s (Mobile)", resp.StatusCode, http.StatusText(resp.StatusCode))
	step.Duration = time.Since(start)
	
	// Success criteria for mobile
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		step.Passed = true
	}
	
	return step
}

func (u *UATTester) allUATStepsSuccessful(steps []UATStep) bool {
	for _, step := range steps {
		if !step.Passed {
			return false
		}
	}
	return true
}

func (u *UATTester) generateUATReport() {
	log.Info().Msg("Generating User Acceptance Test report")
	
	totalScenarios := len(u.scenarios)
	passedScenarios := 0
	failedScenarios := 0
	totalSteps := 0
	passedSteps := 0
	totalDuration := time.Duration(0)
	
	for _, scenario := range u.scenarios {
		if scenario.Success {
			passedScenarios++
		} else {
			failedScenarios++
		}
		totalDuration += scenario.Duration
		totalSteps += len(scenario.TestSteps)
		
		for _, step := range scenario.TestSteps {
			if step.Passed {
				passedSteps++
			}
		}
	}
	
	log.Info().
		Int("total_scenarios", totalScenarios).
		Int("passed_scenarios", passedScenarios).
		Int("failed_scenarios", failedScenarios).
		Int("total_steps", totalSteps).
		Int("passed_steps", passedSteps).
		Float64("scenario_success_rate", float64(passedScenarios)/float64(totalScenarios)*100).
		Float64("step_success_rate", float64(passedSteps)/float64(totalSteps)*100).
		Dur("total_duration", totalDuration).
		Msg("UAT summary")
	
	// Identify failed scenarios for follow-up
	failedScenarios := make([]string, 0)
	for _, scenario := range u.scenarios {
		if !scenario.Success {
			failedScenarios = append(failedScenarios, scenario.ID+": "+scenario.Name)
		}
	}
	
	if len(failedScenarios) > 0 {
		log.Warn().Strs("failed_scenarios", failedScenarios).Msg("UAT scenarios requiring attention")
	}
	
	// Save detailed results
	reportData := map[string]interface{}{
		"summary": map[string]interface{}{
			"total_scenarios":      totalScenarios,
			"passed_scenarios":     passedScenarios,
			"failed_scenarios":     failedScenarios,
			"scenario_success_rate": float64(passedScenarios) / float64(totalScenarios) * 100,
			"total_steps":          totalSteps,
			"passed_steps":         passedSteps,
			"step_success_rate":    float64(passedSteps) / float64(totalSteps) * 100,
			"total_duration_ms":    totalDuration.Milliseconds(),
		},
		"detailed_scenarios": u.scenarios,
		"test_timestamp":     time.Now().Format(time.RFC3339),
		"business_impact":    u.getBusinessImpactAnalysis(),
	}
	
	if jsonData, err := json.MarshalIndent(reportData, "", "  "); err == nil {
		log.Info().Str("report", string(jsonData)).Msg("Detailed UAT report")
	}
}

func (u *UATTester) getBusinessImpactAnalysis() map[string]interface{} {
	return map[string]interface{}{
		"critical_features": []string{
			"Product browsing and search",
			"Shopping cart functionality",
			"Checkout process",
			"User account management",
		},
		"business_risks": []string{
			"Failed checkout process leads to lost sales",
			"Poor mobile experience reduces customer base",
			"Inventory issues cause customer dissatisfaction",
			"Security vulnerabilities damage brand trust",
		},
		"success_metrics": []string{
			"Customer conversion rate",
			"Average order value",
			"Customer satisfaction scores",
			"Mobile traffic percentage",
		},
		"recommendations": []string{
			"Prioritize critical path fixes (browsing → cart → checkout)",
			"Implement comprehensive monitoring for key user journeys",
			"Regular UAT testing with real users",
			"A/B testing for conversion optimization",
			"Performance monitoring especially for mobile users",
		},
	}
}
