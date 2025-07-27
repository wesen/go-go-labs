package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// Test result structures
type TestResult struct {
	TestName     string        `json:"test_name"`
	Endpoint     string        `json:"endpoint"`
	Method       string        `json:"method"`
	StatusCode   int           `json:"status_code"`
	Duration     time.Duration `json:"duration"`
	Success      bool          `json:"success"`
	ErrorMsg     string        `json:"error_msg,omitempty"`
	ResponseSize int           `json:"response_size"`
}

type TestScenario struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Steps       []ScenarioStep  `json:"steps"`
	Success     bool            `json:"success"`
	Duration    time.Duration   `json:"duration"`
	ErrorMsg    string          `json:"error_msg,omitempty"`
}

type ScenarioStep struct {
	StepName string        `json:"step_name"`
	Action   string        `json:"action"`
	URL      string        `json:"url"`
	Method   string        `json:"method"`
	Expected string        `json:"expected"`
	Actual   string        `json:"actual"`
	Success  bool          `json:"success"`
	Duration time.Duration `json:"duration"`
}

type LoadTestResult struct {
	TestName        string                   `json:"test_name"`
	TotalRequests   int                      `json:"total_requests"`
	ConcurrentUsers int                      `json:"concurrent_users"`
	Duration        time.Duration            `json:"duration"`
	SuccessfulReqs  int                      `json:"successful_requests"`
	FailedReqs      int                      `json:"failed_requests"`
	AvgResponseTime time.Duration            `json:"avg_response_time"`
	MinResponseTime time.Duration            `json:"min_response_time"`
	MaxResponseTime time.Duration            `json:"max_response_time"`
	RequestsPerSec  float64                  `json:"requests_per_second"`
	Percentiles     map[string]time.Duration `json:"percentiles"`
	ErrorTypes      map[string]int           `json:"error_types"`
}

type SecurityTestResult struct {
	TestName        string    `json:"test_name"`
	Category        string    `json:"category"`
	Severity        string    `json:"severity"`
	Endpoint        string    `json:"endpoint"`
	Attack          string    `json:"attack"`
	Vulnerable      bool      `json:"vulnerable"`
	Description     string    `json:"description"`
	Recommendation  string    `json:"recommendation"`
	ResponseCode    int       `json:"response_code"`
	ResponseSnippet string    `json:"response_snippet"`
	Timestamp       time.Time `json:"timestamp"`
}

type MobileTestResult struct {
	TestName     string            `json:"test_name"`
	DeviceType   string            `json:"device_type"`
	Viewport     string            `json:"viewport"`
	UserAgent    string            `json:"user_agent"`
	URL          string            `json:"url"`
	LoadTime     time.Duration     `json:"load_time"`
	ResponseCode int               `json:"response_code"`
	ContentSize  int               `json:"content_size"`
	Issues       []string          `json:"issues"`
	Passed       bool              `json:"passed"`
	Timestamp    time.Time         `json:"timestamp"`
	Headers      map[string]string `json:"headers"`
}

type UATScenario struct {
	ID                 string          `json:"id"`
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	BusinessValue      string          `json:"business_value"`
	UserStory          string          `json:"user_story"`
	AcceptanceCriteria []string        `json:"acceptance_criteria"`
	TestSteps          []UATStep       `json:"test_steps"`
	Success            bool            `json:"success"`
	Duration           time.Duration   `json:"duration"`
	ErrorMsg           string          `json:"error_msg,omitempty"`
	Timestamp          time.Time       `json:"timestamp"`
}

type UATStep struct {
	StepNumber int           `json:"step_number"`
	Action     string        `json:"action"`
	Expected   string        `json:"expected"`
	Actual     string        `json:"actual"`
	Passed     bool          `json:"passed"`
	Duration   time.Duration `json:"duration"`
}

// API Tester
type APITester struct {
	baseURL    string
	httpClient *http.Client
	results    []TestResult
}

func (a *APITester) RunTests(ctx context.Context) error {
	log.Info().Msg("Starting API tests for VacuumMart e-commerce platform")
	
	// Test product catalog endpoints
	a.testEndpoint(ctx, "GET", "/api/products", nil, 200, "Get all products")
	a.testEndpoint(ctx, "GET", "/api/products/1", nil, 200, "Get product by ID")
	a.testEndpoint(ctx, "GET", "/api/products?category=upright", nil, 200, "Get upright vacuums")
	a.testEndpoint(ctx, "GET", "/api/products?brand=Dyson", nil, 200, "Get Dyson products")
	a.testEndpoint(ctx, "GET", "/api/products/search?q=pet+hair", nil, 200, "Search for pet hair vacuums")
	
	// Test shopping cart endpoints
	a.testEndpoint(ctx, "POST", "/api/cart/items", map[string]interface{}{"product_id": 1, "quantity": 1}, 200, "Add item to cart")
	a.testEndpoint(ctx, "GET", "/api/cart", nil, 200, "Get cart contents")
	a.testEndpoint(ctx, "PUT", "/api/cart/items/1", map[string]interface{}{"quantity": 2}, 200, "Update cart item")
	a.testEndpoint(ctx, "DELETE", "/api/cart/items/1", nil, 200, "Remove cart item")
	
	// Test order endpoints
	a.testEndpoint(ctx, "POST", "/api/orders", map[string]interface{}{"customer_id": 1, "items": []map[string]interface{}{{"product_id": 1, "quantity": 1}}}, 201, "Create order")
	a.testEndpoint(ctx, "GET", "/api/orders/1", nil, 200, "Get order details")
	
	a.generateAPIReport()
	return nil
}

func (a *APITester) testEndpoint(ctx context.Context, method, endpoint string, body interface{}, expectedStatus int, testName string) {
	start := time.Now()
	
	var reqBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}
	
	url := a.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		a.results = append(a.results, TestResult{TestName: testName, ErrorMsg: err.Error(), Duration: time.Since(start)})
		return
	}
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	resp, err := a.httpClient.Do(req)
	duration := time.Since(start)
	
	result := TestResult{
		TestName: testName,
		Endpoint: endpoint,
		Method:   method,
		Duration: duration,
		Success:  false,
	}
	
	if err != nil {
		result.ErrorMsg = err.Error()
	} else {
		defer resp.Body.Close()
		result.StatusCode = resp.StatusCode
		respBody, _ := io.ReadAll(resp.Body)
		result.ResponseSize = len(respBody)
		result.Success = resp.StatusCode == expectedStatus
	}
	
	a.results = append(a.results, result)
	log.Info().Str("test", testName).Bool("success", result.Success).Dur("duration", duration).Msg("API test completed")
}

func (a *APITester) generateAPIReport() {
	passed := 0
	for _, result := range a.results {
		if result.Success {
			passed++
		}
	}
	log.Info().Int("total", len(a.results)).Int("passed", passed).Int("failed", len(a.results)-passed).Msg("API test summary")
}

// E2E Tester
type E2ETester struct {
	baseURL    string
	httpClient *http.Client
	scenarios  []TestScenario
}

func (e *E2ETester) RunTests(ctx context.Context) error {
	log.Info().Msg("Starting E2E tests for VacuumMart shopping flows")
	
	scenario := TestScenario{
		Name:        "Complete Shopping Flow",
		Description: "User browses products, adds items to cart, and completes checkout",
		Steps:       make([]ScenarioStep, 0),
	}
	
	// Execute shopping flow steps
	steps := []struct{ name, method, endpoint, expected string }{
		{"Visit Homepage", "GET", "/", "200 OK"},
		{"Browse Products", "GET", "/api/products", "200 OK"},
		{"Search for Dyson", "GET", "/api/products/search?q=dyson", "200 OK"},
		{"View Product Details", "GET", "/api/products/1", "200 OK"},
		{"Add to Cart", "POST", "/api/cart/items", "200 OK"},
		{"View Cart", "GET", "/api/cart", "200 OK"},
		{"Proceed to Checkout", "GET", "/checkout", "200 OK"},
		{"Create Order", "POST", "/api/orders", "201 Created"},
	}
	
	for _, step := range steps {
		scenarioStep := e.executeStep(ctx, step.name, step.method, step.endpoint, step.expected)
		scenario.Steps = append(scenario.Steps, scenarioStep)
	}
	
	scenario.Success = e.allStepsSuccessful(scenario.Steps)
	e.scenarios = append(e.scenarios, scenario)
	
	log.Info().Str("scenario", scenario.Name).Bool("success", scenario.Success).Msg("E2E scenario completed")
	return nil
}

func (e *E2ETester) executeStep(ctx context.Context, stepName, method, endpoint, expected string) ScenarioStep {
	start := time.Now()
	
	url := e.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return ScenarioStep{StepName: stepName, Success: false, Duration: time.Since(start)}
	}
	
	resp, err := e.httpClient.Do(req)
	duration := time.Since(start)
	
	step := ScenarioStep{
		StepName: stepName,
		Action:   fmt.Sprintf("%s %s", method, endpoint),
		URL:      endpoint,
		Method:   method,
		Expected: expected,
		Duration: duration,
		Success:  false,
	}
	
	if err != nil {
		step.Actual = fmt.Sprintf("Error: %v", err)
	} else {
		defer resp.Body.Close()
		step.Actual = fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		step.Success = resp.StatusCode >= 200 && resp.StatusCode < 300
	}
	
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

// Performance Tester
type PerformanceTester struct {
	baseURL    string
	httpClient *http.Client
	results    []LoadTestResult
}

func (p *PerformanceTester) RunTests(ctx context.Context) error {
	log.Info().Msg("Starting performance tests for VacuumMart")
	
	result, _ := p.runLoadTest(ctx, "Basic Load Test", "/api/products", 10, 30*time.Second)
	p.results = append(p.results, result)
	
	log.Info().Str("test", result.TestName).Float64("rps", result.RequestsPerSec).Dur("avg_time", result.AvgResponseTime).Msg("Performance test completed")
	return nil
}

func (p *PerformanceTester) runLoadTest(ctx context.Context, testName, endpoint string, concurrentUsers int, duration time.Duration) (LoadTestResult, error) {
	result := LoadTestResult{
		TestName:        testName,
		ConcurrentUsers: concurrentUsers,
		Duration:        duration,
		Percentiles:     make(map[string]time.Duration),
		ErrorTypes:      make(map[string]int),
	}
	
	start := time.Now()
	endTime := start.Add(duration)
	
	var mu sync.Mutex
	var requestResults []time.Duration
	
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(concurrentUsers)
	
	for time.Now().Before(endTime) {
		g.Go(func() error {
			reqStart := time.Now()
			url := p.baseURL + endpoint
			req, err := http.NewRequestWithContext(gCtx, "GET", url, nil)
			if err != nil {
				mu.Lock()
				result.FailedReqs++
				result.TotalRequests++
				mu.Unlock()
				return nil
			}
			
			resp, err := p.httpClient.Do(req)
			reqDuration := time.Since(reqStart)
			
			mu.Lock()
			result.TotalRequests++
			requestResults = append(requestResults, reqDuration)
			if err != nil || resp.StatusCode >= 400 {
				result.FailedReqs++
			} else {
				result.SuccessfulReqs++
			}
			mu.Unlock()
			
			if resp != nil {
				resp.Body.Close()
			}
			
			return nil
		})
		
		time.Sleep(100 * time.Millisecond) // Throttle request generation
	}
	
	g.Wait()
	
	// Calculate statistics
	if len(requestResults) > 0 {
		var total time.Duration
		result.MinResponseTime = requestResults[0]
		result.MaxResponseTime = requestResults[0]
		
		for _, duration := range requestResults {
			total += duration
			if duration < result.MinResponseTime {
				result.MinResponseTime = duration
			}
			if duration > result.MaxResponseTime {
				result.MaxResponseTime = duration
			}
		}
		
		result.AvgResponseTime = total / time.Duration(len(requestResults))
		result.RequestsPerSec = float64(result.TotalRequests) / result.Duration.Seconds()
	}
	
	return result, nil
}

// Security Tester
type SecurityTester struct {
	baseURL    string
	httpClient *http.Client
	results    []SecurityTestResult
}

func (s *SecurityTester) RunTests(ctx context.Context) error {
	log.Info().Msg("Starting security tests for VacuumMart")
	
	// Test XSS vulnerabilities
	xssPayloads := []string{"<script>alert('XSS')</script>", "<img src=x onerror=alert('XSS')>"}
	for _, payload := range xssPayloads {
		s.testXSSEndpoint(ctx, "/api/products/search", "q", payload, "Search XSS")
	}
	
	// Test SQL injection
	sqlPayloads := []string{"' OR '1'='1", "'; DROP TABLE products;--"}
	for _, payload := range sqlPayloads {
		s.testSQLInjectionEndpoint(ctx, "/api/products/search", "q", payload, "Search SQL Injection")
	}
	
	s.generateSecurityReport()
	return nil
}

func (s *SecurityTester) testXSSEndpoint(ctx context.Context, endpoint, param, payload, testName string) {
	s.testSecurityEndpoint(ctx, "GET", endpoint, param, payload, testName, "XSS", "HIGH")
}

func (s *SecurityTester) testSQLInjectionEndpoint(ctx context.Context, endpoint, param, payload, testName string) {
	s.testSecurityEndpoint(ctx, "GET", endpoint, param, payload, testName, "SQL Injection", "CRITICAL")
}

func (s *SecurityTester) testSecurityEndpoint(ctx context.Context, method, endpoint, param, payload, testName, category, severity string) {
	start := time.Now()
	
	url := fmt.Sprintf("%s%s?%s=%s", s.baseURL, endpoint, param, url.QueryEscape(payload))
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	
	buf := make([]byte, 1024)
	n, _ := resp.Body.Read(buf)
	responseBody := string(buf[:n])
	
	vulnerable := strings.Contains(responseBody, payload) || resp.StatusCode == 500
	
	result := SecurityTestResult{
		TestName:        testName,
		Category:        category,
		Severity:        severity,
		Endpoint:        endpoint,
		Attack:          payload,
		Vulnerable:      vulnerable,
		Description:     fmt.Sprintf("Potential %s vulnerability detected", category),
		Recommendation:  "Implement proper input sanitization",
		ResponseCode:    resp.StatusCode,
		ResponseSnippet: responseBody[:min(len(responseBody), 200)],
		Timestamp:       start,
	}
	
	s.results = append(s.results, result)
	log.Info().Str("test", testName).Bool("vulnerable", vulnerable).Msg("Security test completed")
}

func (s *SecurityTester) generateSecurityReport() {
	vulnerabilities := 0
	for _, result := range s.results {
		if result.Vulnerable {
			vulnerabilities++
		}
	}
	log.Info().Int("total", len(s.results)).Int("vulnerabilities", vulnerabilities).Msg("Security test summary")
}

// Mobile Tester
type MobileTester struct {
	baseURL    string
	httpClient *http.Client
	results    []MobileTestResult
}

func (m *MobileTester) RunTests(ctx context.Context) error {
	log.Info().Msg("Starting mobile responsiveness tests for VacuumMart")
	
	devices := []struct{ name, userAgent, viewport string }{
		{"iPhone 14", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15", "393x852"},
		{"Samsung Galaxy", "Mozilla/5.0 (Linux; Android 13; SM-S911B) AppleWebKit/537.36", "360x800"},
		{"iPad Pro", "Mozilla/5.0 (iPad; CPU OS 16_0 like Mac OS X) AppleWebKit/605.1.15", "1024x1366"},
	}
	
	pages := []string{"/", "/products", "/cart", "/checkout"}
	
	for _, device := range devices {
		for _, page := range pages {
			result := m.testPageOnDevice(ctx, page, device.name, device.userAgent, device.viewport)
			m.results = append(m.results, result)
		}
	}
	
	m.generateMobileReport()
	return nil
}

func (m *MobileTester) testPageOnDevice(ctx context.Context, page, deviceName, userAgent, viewport string) MobileTestResult {
	start := time.Now()
	
	url := m.baseURL + page
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return MobileTestResult{TestName: fmt.Sprintf("%s on %s", page, deviceName), Issues: []string{err.Error()}}
	}
	
	req.Header.Set("User-Agent", userAgent)
	
	resp, err := m.httpClient.Do(req)
	loadTime := time.Since(start)
	
	result := MobileTestResult{
		TestName:    fmt.Sprintf("%s on %s", page, deviceName),
		DeviceType:  "mobile",
		Viewport:    viewport,
		UserAgent:   userAgent,
		URL:         page,
		LoadTime:    loadTime,
		Issues:      make([]string, 0),
		Timestamp:   start,
		Headers:     make(map[string]string),
	}
	
	if err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("Request failed: %v", err))
		return result
	}
	defer resp.Body.Close()
	
	result.ResponseCode = resp.StatusCode
	body, _ := io.ReadAll(resp.Body)
	result.ContentSize = len(body)
	
	// Check mobile-friendliness
	bodyStr := string(body)
	if !strings.Contains(bodyStr, "viewport") {
		result.Issues = append(result.Issues, "Missing viewport meta tag")
	}
	if loadTime > 3*time.Second {
		result.Issues = append(result.Issues, "Slow load time for mobile")
	}
	
	result.Passed = len(result.Issues) == 0
	log.Info().Str("test", result.TestName).Bool("passed", result.Passed).Dur("load_time", loadTime).Msg("Mobile test completed")
	
	return result
}

func (m *MobileTester) generateMobileReport() {
	passed := 0
	for _, result := range m.results {
		if result.Passed {
			passed++
		}
	}
	log.Info().Int("total", len(m.results)).Int("passed", passed).Int("failed", len(m.results)-passed).Msg("Mobile test summary")
}

// UAT Tester
type UATTester struct {
	baseURL    string
	httpClient *http.Client
	scenarios  []UATScenario
}

func (u *UATTester) RunTests(ctx context.Context) error {
	log.Info().Msg("Starting User Acceptance Tests for VacuumMart")
	
	scenario := UATScenario{
		ID:          "UAT-001",
		Name:        "Customer Product Browsing",
		Description: "Customer can browse vacuum cleaners by category and brand",
		UserStory:   "As a customer, I want to browse vacuum cleaners by category",
		AcceptanceCriteria: []string{
			"Customer can view all vacuum categories",
			"Customer can filter by brand",
			"Product listings show key information",
		},
		TestSteps: make([]UATStep, 0),
	}
	
	steps := []struct{ action, endpoint, expected string }{
		{"Navigate to product catalog", "/api/products", "Product catalog loads successfully"},
		{"Browse upright vacuums", "/api/products?category=upright", "Upright vacuums displayed"},
		{"Filter by Dyson brand", "/api/products?brand=Dyson", "Dyson products displayed"},
	}
	
	for i, step := range steps {
		uatStep := u.executeUATStep(ctx, i+1, step.action, "GET", step.endpoint, step.expected)
		scenario.TestSteps = append(scenario.TestSteps, uatStep)
	}
	
	scenario.Success = u.allUATStepsSuccessful(scenario.TestSteps)
	u.scenarios = append(u.scenarios, scenario)
	
	log.Info().Str("scenario", scenario.ID).Bool("success", scenario.Success).Msg("UAT scenario completed")
	return nil
}

func (u *UATTester) executeUATStep(ctx context.Context, stepNumber int, action, method, endpoint, expected string) UATStep {
	start := time.Now()
	
	url := u.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return UATStep{StepNumber: stepNumber, Action: action, Expected: expected, Actual: err.Error(), Passed: false, Duration: time.Since(start)}
	}
	
	resp, err := u.httpClient.Do(req)
	duration := time.Since(start)
	
	step := UATStep{
		StepNumber: stepNumber,
		Action:     action,
		Expected:   expected,
		Duration:   duration,
		Passed:     false,
	}
	
	if err != nil {
		step.Actual = fmt.Sprintf("Error: %v", err)
	} else {
		defer resp.Body.Close()
		step.Actual = fmt.Sprintf("HTTP %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		step.Passed = resp.StatusCode >= 200 && resp.StatusCode < 300
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
