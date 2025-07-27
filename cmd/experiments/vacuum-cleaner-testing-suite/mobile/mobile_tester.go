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

type MobileTester struct {
	baseURL    string
	httpClient *http.Client
	results    []MobileTestResult
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

type DeviceProfile struct {
	Name      string
	UserAgent string
	Viewport  string
	Type      string
}

func NewMobileTester(baseURL string) *MobileTester {
	return &MobileTester{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		results: make([]MobileTestResult, 0),
	}
}

func (m *MobileTester) RunTests(ctx context.Context) error {
	log.Info().Msg("Starting mobile responsiveness tests for VacuumMart")
	
	// Define device profiles for testing
	devices := []DeviceProfile{
		{
			Name:      "iPhone 14 Pro",
			UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
			Viewport:  "393x852",
			Type:      "mobile",
		},
		{
			Name:      "iPhone SE",
			UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/604.1",
			Viewport:  "375x667",
			Type:      "mobile",
		},
		{
			Name:      "Samsung Galaxy S23",
			UserAgent: "Mozilla/5.0 (Linux; Android 13; SM-S911B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
			Viewport:  "360x800",
			Type:      "mobile",
		},
		{
			Name:      "iPad Pro",
			UserAgent: "Mozilla/5.0 (iPad; CPU OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
			Viewport:  "1024x1366",
			Type:      "tablet",
		},
		{
			Name:      "iPad Mini",
			UserAgent: "Mozilla/5.0 (iPad; CPU OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/604.1",
			Viewport:  "768x1024",
			Type:      "tablet",
		},
		{
			Name:      "Android Tablet",
			UserAgent: "Mozilla/5.0 (Linux; Android 12; SM-T870) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
			Viewport:  "800x1280",
			Type:      "tablet",
		},
	}
	
	// Test pages that should be mobile-friendly
	testPages := []string{
		"/",
		"/products",
		"/product/1",
		"/cart",
		"/checkout",
		"/search",
		"/category/robot",
		"/login",
		"/register",
	}
	
	// Run tests for each device and page combination
	for _, device := range devices {
		for _, page := range testPages {
			result := m.testPageOnDevice(ctx, page, device)
			m.results = append(m.results, result)
		}
	}
	
	// Run specific mobile functionality tests
	if err := m.testMobileSpecificFeatures(ctx, devices); err != nil {
		log.Error().Err(err).Msg("Mobile-specific feature tests failed")
	}
	
	// Test touch interactions
	if err := m.testTouchInteractions(ctx, devices); err != nil {
		log.Error().Err(err).Msg("Touch interaction tests failed")
	}
	
	// Test mobile performance
	if err := m.testMobilePerformance(ctx, devices); err != nil {
		log.Error().Err(err).Msg("Mobile performance tests failed")
	}
	
	m.generateMobileReport()
	return nil
}

func (m *MobileTester) testPageOnDevice(ctx context.Context, page string, device DeviceProfile) MobileTestResult {
	start := time.Now()
	
	result := MobileTestResult{
		TestName:   fmt.Sprintf("%s on %s", page, device.Name),
		DeviceType: device.Type,
		Viewport:   device.Viewport,
		UserAgent:  device.UserAgent,
		URL:        page,
		Issues:     make([]string, 0),
		Headers:    make(map[string]string),
		Timestamp:  start,
	}
	
	// Make request with mobile user agent
	url := m.baseURL + page
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("Failed to create request: %v", err))
		return result
	}
	
	// Set mobile user agent and viewport headers
	req.Header.Set("User-Agent", device.UserAgent)
	req.Header.Set("Viewport", device.Viewport)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	
	resp, err := m.httpClient.Do(req)
	loadTime := time.Since(start)
	result.LoadTime = loadTime
	
	if err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("Request failed: %v", err))
		return result
	}
	defer resp.Body.Close()
	
	result.ResponseCode = resp.StatusCode
	
	// Store response headers
	for name, values := range resp.Header {
		if len(values) > 0 {
			result.Headers[name] = values[0]
		}
	}
	
	// Read response body
	body := make([]byte, 10*1024*1024) // 10MB max
	n, _ := resp.Body.Read(body)
	result.ContentSize = n
	responseBody := string(body[:n])
	
	// Analyze mobile-friendliness
	m.analyzeMobileFriendliness(&result, responseBody, device)
	
	// Check performance criteria
	m.checkMobilePerformance(&result, loadTime, n)
	
	// Check for mobile-specific issues
	m.checkMobileIssues(&result, responseBody, resp.Header)
	
	result.Passed = len(result.Issues) == 0
	
	log.Info().
		Str("page", page).
		Str("device", device.Name).
		Dur("load_time", loadTime).
		Int("response_code", resp.StatusCode).
		Int("content_size", n).
		Bool("passed", result.Passed).
		Strs("issues", result.Issues).
		Msg("Mobile test completed")
	
	return result
}

func (m *MobileTester) analyzeMobileFriendliness(result *MobileTestResult, body string, device DeviceProfile) {
	bodyLower := strings.ToLower(body)
	
	// Check for viewport meta tag
	if !strings.Contains(bodyLower, `<meta name="viewport"`) {
		result.Issues = append(result.Issues, "Missing viewport meta tag")
	}
	
	// Check for responsive design indicators
	if !strings.Contains(bodyLower, "@media") && !strings.Contains(bodyLower, "responsive") {
		result.Issues = append(result.Issues, "No responsive design indicators found")
	}
	
	// Check for mobile-friendly CSS frameworks
	mobileFrameworks := []string{"bootstrap", "foundation", "bulma", "tailwind"}
	foundFramework := false
	for _, framework := range mobileFrameworks {
		if strings.Contains(bodyLower, framework) {
			foundFramework = true
			break
		}
	}
	if !foundFramework {
		result.Issues = append(result.Issues, "No mobile-friendly CSS framework detected")
	}
	
	// Check for touch-friendly elements
	if !strings.Contains(bodyLower, "touch") && !strings.Contains(bodyLower, "tap") {
		result.Issues = append(result.Issues, "No touch-specific interaction elements found")
	}
	
	// Check for mobile-optimized images
	if strings.Contains(bodyLower, "<img") && !strings.Contains(bodyLower, "srcset") {
		result.Issues = append(result.Issues, "Images not optimized for mobile (no srcset)")
	}
	
	// Check for hamburger menu or mobile navigation
	mobileNavIndicators := []string{"hamburger", "menu-toggle", "nav-toggle", "mobile-menu"}
	foundMobileNav := false
	for _, indicator := range mobileNavIndicators {
		if strings.Contains(bodyLower, indicator) {
			foundMobileNav = true
			break
		}
	}
	if !foundMobileNav && result.DeviceType == "mobile" {
		result.Issues = append(result.Issues, "No mobile navigation pattern detected")
	}
	
	// Check for excessive horizontal scrolling indicators
	if strings.Contains(bodyLower, "overflow-x: scroll") || strings.Contains(bodyLower, "white-space: nowrap") {
		result.Issues = append(result.Issues, "Potential horizontal scrolling detected")
	}
}

func (m *MobileTester) checkMobilePerformance(result *MobileTestResult, loadTime time.Duration, contentSize int) {
	// Performance thresholds for mobile
	maxLoadTime := 3 * time.Second
	maxContentSize := 2 * 1024 * 1024 // 2MB
	
	if loadTime > maxLoadTime {
		result.Issues = append(result.Issues, fmt.Sprintf("Slow load time: %v (target: <%v)", loadTime, maxLoadTime))
	}
	
	if contentSize > maxContentSize {
		result.Issues = append(result.Issues, fmt.Sprintf("Large content size: %d bytes (target: <%d)", contentSize, maxContentSize))
	}
	
	// Check for mobile-specific performance issues
	if result.DeviceType == "mobile" {
		mobileMaxLoadTime := 2 * time.Second
		if loadTime > mobileMaxLoadTime {
			result.Issues = append(result.Issues, fmt.Sprintf("Too slow for mobile: %v (mobile target: <%v)", loadTime, mobileMaxLoadTime))
		}
	}
}

func (m *MobileTester) checkMobileIssues(result *MobileTestResult, body string, headers http.Header) {
	bodyLower := strings.ToLower(body)
	
	// Check for Flash content (not supported on mobile)
	if strings.Contains(bodyLower, "flash") || strings.Contains(bodyLower, ".swf") {
		result.Issues = append(result.Issues, "Flash content detected (not mobile compatible)")
	}
	
	// Check for fixed-width layouts
	if strings.Contains(bodyLower, "width: 1200px") || strings.Contains(bodyLower, "min-width: 1024px") {
		result.Issues = append(result.Issues, "Fixed-width layout detected (not mobile-friendly)")
	}
	
	// Check for missing mobile-specific headers
	if headers.Get("Cache-Control") == "" {
		result.Issues = append(result.Issues, "Missing Cache-Control header (important for mobile)")
	}
	
	// Check for compression
	if headers.Get("Content-Encoding") == "" && result.ContentSize > 1024 {
		result.Issues = append(result.Issues, "Content not compressed (important for mobile)")
	}
	
	// Check for mobile app store links
	if result.DeviceType == "mobile" {
		hasAppStoreLink := strings.Contains(bodyLower, "app store") || strings.Contains(bodyLower, "google play")
		if !hasAppStoreLink {
			// This is informational, not necessarily an issue
			log.Info().Msg("No mobile app store links found (consider adding if mobile app exists)")
		}
	}
	
	// Check for popup issues on mobile
	if strings.Contains(bodyLower, "popup") || strings.Contains(bodyLower, "modal") {
		result.Issues = append(result.Issues, "Popups/modals detected (ensure mobile-friendly)")
	}
}

func (m *MobileTester) testMobileSpecificFeatures(ctx context.Context, devices []DeviceProfile) error {
	log.Info().Msg("Testing mobile-specific features")
	
	mobileDevices := make([]DeviceProfile, 0)
	for _, device := range devices {
		if device.Type == "mobile" {
			mobileDevices = append(mobileDevices, device)
		}
	}
	
	// Test mobile-specific endpoints
	mobileEndpoints := []string{
		"/api/mobile/version",
		"/api/mobile/config",
		"/mobile",
		"/m",
	}
	
	for _, endpoint := range mobileEndpoints {
		for _, device := range mobileDevices {
			result := m.testMobileEndpoint(ctx, endpoint, device)
			m.results = append(m.results, result)
		}
	}
	
	return nil
}

func (m *MobileTester) testMobileEndpoint(ctx context.Context, endpoint string, device DeviceProfile) MobileTestResult {
	start := time.Now()
	
	result := MobileTestResult{
		TestName:   fmt.Sprintf("Mobile feature: %s on %s", endpoint, device.Name),
		DeviceType: device.Type,
		Viewport:   device.Viewport,
		UserAgent:  device.UserAgent,
		URL:        endpoint,
		Issues:     make([]string, 0),
		Timestamp:  start,
	}
	
	url := m.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("Failed to create request: %v", err))
		return result
	}
	
	req.Header.Set("User-Agent", device.UserAgent)
	
	resp, err := m.httpClient.Do(req)
	result.LoadTime = time.Since(start)
	
	if err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("Request failed: %v", err))
		return result
	}
	defer resp.Body.Close()
	
	result.ResponseCode = resp.StatusCode
	
	// Mobile-specific endpoints should either work (200) or not exist (404)
	// 500 errors indicate potential issues
	if resp.StatusCode == 500 {
		result.Issues = append(result.Issues, "Server error on mobile endpoint")
	}
	
	result.Passed = len(result.Issues) == 0
	
	return result
}

func (m *MobileTester) testTouchInteractions(ctx context.Context, devices []DeviceProfile) error {
	log.Info().Msg("Testing touch interaction compatibility")
	
	// Test pages with interactive elements
	interactivePages := []string{
		"/cart",
		"/products",
		"/search",
		"/checkout",
	}
	
	for _, page := range interactivePages {
		for _, device := range devices {
			if device.Type == "mobile" || device.Type == "tablet" {
				result := m.testTouchFriendliness(ctx, page, device)
				m.results = append(m.results, result)
			}
		}
	}
	
	return nil
}

func (m *MobileTester) testTouchFriendliness(ctx context.Context, page string, device DeviceProfile) MobileTestResult {
	start := time.Now()
	
	result := MobileTestResult{
		TestName:   fmt.Sprintf("Touch test: %s on %s", page, device.Name),
		DeviceType: device.Type,
		Viewport:   device.Viewport,
		UserAgent:  device.UserAgent,
		URL:        page,
		Issues:     make([]string, 0),
		Timestamp:  start,
	}
	
	url := m.baseURL + page
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("Failed to create request: %v", err))
		return result
	}
	
	req.Header.Set("User-Agent", device.UserAgent)
	
	resp, err := m.httpClient.Do(req)
	result.LoadTime = time.Since(start)
	
	if err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("Request failed: %v", err))
		return result
	}
	defer resp.Body.Close()
	
	result.ResponseCode = resp.StatusCode
	
	body := make([]byte, 1024*1024) // 1MB
	n, _ := resp.Body.Read(body)
	result.ContentSize = n
	responseBody := string(body[:n])
	
	m.analyzeTouchFriendliness(&result, responseBody)
	
	result.Passed = len(result.Issues) == 0
	
	return result
}

func (m *MobileTester) analyzeTouchFriendliness(result *MobileTestResult, body string) {
	bodyLower := strings.ToLower(body)
	
	// Check for touch event handlers
	touchEvents := []string{"touchstart", "touchend", "touchmove", "onclick"}
	foundTouchEvents := false
	for _, event := range touchEvents {
		if strings.Contains(bodyLower, event) {
			foundTouchEvents = true
			break
		}
	}
	if !foundTouchEvents {
		result.Issues = append(result.Issues, "No touch event handlers found")
	}
	
	// Check for appropriately sized touch targets
	if strings.Contains(bodyLower, "font-size: 8px") || strings.Contains(bodyLower, "font-size: 10px") {
		result.Issues = append(result.Issues, "Text may be too small for touch devices")
	}
	
	// Check for hover-dependent interactions (problematic on touch devices)
	hoverIndicators := []string{":hover", "onmouseover", "onmouseout"}
	for _, indicator := range hoverIndicators {
		if strings.Contains(bodyLower, indicator) {
			result.Issues = append(result.Issues, "Hover-dependent interactions detected (problematic on touch devices)")
			break
		}
	}
	
	// Check for touch-friendly form elements
	if strings.Contains(bodyLower, "<form") {
		if !strings.Contains(bodyLower, `type="tel"`) && !strings.Contains(bodyLower, `type="email"`) {
			result.Issues = append(result.Issues, "Forms may not be optimized for mobile input")
		}
	}
}

func (m *MobileTester) testMobilePerformance(ctx context.Context, devices []DeviceProfile) error {
	log.Info().Msg("Testing mobile performance characteristics")
	
	performancePages := []string{
		"/",
		"/products",
		"/search?q=vacuum",
	}
	
	for _, page := range performancePages {
		for _, device := range devices {
			if device.Type == "mobile" {
				result := m.testPagePerformance(ctx, page, device)
				m.results = append(m.results, result)
			}
		}
	}
	
	return nil
}

func (m *MobileTester) testPagePerformance(ctx context.Context, page string, device DeviceProfile) MobileTestResult {
	start := time.Now()
	
	result := MobileTestResult{
		TestName:   fmt.Sprintf("Performance: %s on %s", page, device.Name),
		DeviceType: device.Type,
		Viewport:   device.Viewport,
		UserAgent:  device.UserAgent,
		URL:        page,
		Issues:     make([]string, 0),
		Timestamp:  start,
	}
	
	// Simulate slower mobile connection
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	
	url := m.baseURL + page
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("Failed to create request: %v", err))
		return result
	}
	
	req.Header.Set("User-Agent", device.UserAgent)
	
	resp, err := m.httpClient.Do(req)
	result.LoadTime = time.Since(start)
	
	if err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("Request failed: %v", err))
		return result
	}
	defer resp.Body.Close()
	
	result.ResponseCode = resp.StatusCode
	
	body := make([]byte, 5*1024*1024) // 5MB max
	n, _ := resp.Body.Read(body)
	result.ContentSize = n
	
	// Mobile performance thresholds
	maxMobileLoadTime := 2 * time.Second
	maxMobileContentSize := 1 * 1024 * 1024 // 1MB
	
	if result.LoadTime > maxMobileLoadTime {
		result.Issues = append(result.Issues, fmt.Sprintf("Slow mobile load time: %v", result.LoadTime))
	}
	
	if result.ContentSize > maxMobileContentSize {
		result.Issues = append(result.Issues, fmt.Sprintf("Large mobile content size: %d bytes", result.ContentSize))
	}
	
	result.Passed = len(result.Issues) == 0
	
	return result
}

func (m *MobileTester) generateMobileReport() {
	log.Info().Msg("Generating mobile test report")
	
	totalTests := len(m.results)
	passedTests := 0
	failedTests := 0
	
	deviceTypeResults := make(map[string]int)
	commonIssues := make(map[string]int)
	
	var totalLoadTime time.Duration
	var totalContentSize int
	
	for _, result := range m.results {
		if result.Passed {
			passedTests++
		} else {
			failedTests++
		}
		
		deviceTypeResults[result.DeviceType]++
		totalLoadTime += result.LoadTime
		totalContentSize += result.ContentSize
		
		for _, issue := range result.Issues {
			commonIssues[issue]++
		}
	}
	
	avgLoadTime := totalLoadTime / time.Duration(totalTests)
	avgContentSize := totalContentSize / totalTests
	
	log.Info().
		Int("total_tests", totalTests).
		Int("passed", passedTests).
		Int("failed", failedTests).
		Float64("success_rate", float64(passedTests)/float64(totalTests)*100).
		Dur("avg_load_time", avgLoadTime).
		Int("avg_content_size", avgContentSize).
		Interface("device_breakdown", deviceTypeResults).
		Msg("Mobile test summary")
	
	// Log most common issues
	log.Info().Interface("common_issues", commonIssues).Msg("Most common mobile issues")
	
	// Save detailed results
	reportData := map[string]interface{}{
		"summary": map[string]interface{}{
			"total_tests":       totalTests,
			"passed_tests":      passedTests,
			"failed_tests":      failedTests,
			"success_rate":      float64(passedTests) / float64(totalTests) * 100,
			"avg_load_time_ms":  avgLoadTime.Milliseconds(),
			"avg_content_size":  avgContentSize,
		},
		"device_breakdown":    deviceTypeResults,
		"common_issues":       commonIssues,
		"detailed_results":    m.results,
		"test_timestamp":      time.Now().Format(time.RFC3339),
		"recommendations":     m.getMobileRecommendations(),
	}
	
	if jsonData, err := json.MarshalIndent(reportData, "", "  "); err == nil {
		log.Info().Str("report", string(jsonData)).Msg("Detailed mobile test report")
	}
}

func (m *MobileTester) getMobileRecommendations() []string {
	return []string{
		"Implement responsive design with proper viewport meta tag",
		"Use mobile-friendly CSS frameworks (Bootstrap, Tailwind, etc.)",
		"Optimize images for mobile with srcset and appropriate sizing",
		"Implement touch-friendly navigation (hamburger menu, etc.)",
		"Ensure touch targets are at least 44px for accessibility",
		"Use appropriate input types for mobile keyboards (tel, email, etc.)",
		"Implement mobile-specific performance optimizations",
		"Test on real devices, not just browser dev tools",
		"Consider implementing a mobile app or PWA",
		"Use mobile-first design approach",
		"Implement proper caching for mobile networks",
		"Minimize the use of hover-dependent interactions",
		"Ensure forms are optimized for mobile input",
		"Test across different mobile browsers and OS versions",
	}
}
