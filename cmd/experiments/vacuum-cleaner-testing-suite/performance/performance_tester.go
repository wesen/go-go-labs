package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type PerformanceTester struct {
	baseURL    string
	httpClient *http.Client
	results    []LoadTestResult
}

type LoadTestResult struct {
	TestName        string                 `json:"test_name"`
	TotalRequests   int                    `json:"total_requests"`
	ConcurrentUsers int                    `json:"concurrent_users"`
	Duration        time.Duration          `json:"duration"`
	SuccessfulReqs  int                    `json:"successful_requests"`
	FailedReqs      int                    `json:"failed_requests"`
	AvgResponseTime time.Duration          `json:"avg_response_time"`
	MinResponseTime time.Duration          `json:"min_response_time"`
	MaxResponseTime time.Duration          `json:"max_response_time"`
	RequestsPerSec  float64                `json:"requests_per_second"`
	Percentiles     map[string]time.Duration `json:"percentiles"`
	ErrorTypes      map[string]int         `json:"error_types"`
}

type RequestResult struct {
	Duration    time.Duration
	StatusCode  int
	Success     bool
	Error       string
}

func NewPerformanceTester(baseURL string) *PerformanceTester {
	return &PerformanceTester{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		results: make([]LoadTestResult, 0),
	}
}

func (p *PerformanceTester) RunTests(ctx context.Context) error {
	log.Info().Msg("Starting performance tests for VacuumMart e-commerce platform")
	
	// Test scenarios with different load patterns
	testScenarios := []struct {
		name            string
		endpoint        string
		concurrentUsers int
		duration        time.Duration
	}{
		{"Homepage Load Test", "/", 50, 30 * time.Second},
		{"Product Catalog Load Test", "/api/products", 100, 60 * time.Second},
		{"Product Search Load Test", "/api/products/search?q=dyson", 75, 45 * time.Second},
		{"Single Product Load Test", "/api/products/1", 200, 30 * time.Second},
		{"Cart Operations Load Test", "/api/cart", 50, 60 * time.Second},
		{"High Traffic Spike Test", "/api/products", 500, 15 * time.Second},
		{"Sustained Load Test", "/api/products", 25, 300 * time.Second}, // 5 minutes
	}
	
	for _, scenario := range testScenarios {
		log.Info().
			Str("test", scenario.name).
			Str("endpoint", scenario.endpoint).
			Int("users", scenario.concurrentUsers).
			Dur("duration", scenario.duration).
			Msg("Starting load test scenario")
		
		result, err := p.runLoadTest(ctx, scenario.name, scenario.endpoint, scenario.concurrentUsers, scenario.duration)
		if err != nil {
			log.Error().Err(err).Str("test", scenario.name).Msg("Load test failed")
			continue
		}
		
		p.results = append(p.results, result)
		p.logTestResult(result)
	}
	
	// Run specialized performance tests
	if err := p.runStressTest(ctx); err != nil {
		log.Error().Err(err).Msg("Stress test failed")
	}
	
	if err := p.runSpikeTest(ctx); err != nil {
		log.Error().Err(err).Msg("Spike test failed")
	}
	
	if err := p.runEnduranceTest(ctx); err != nil {
		log.Error().Err(err).Msg("Endurance test failed")
	}
	
	p.generateReport()
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
	requestResults := make([]RequestResult, 0)
	
	// Create context with timeout
	testCtx, cancel := context.WithTimeout(ctx, duration+10*time.Second)
	defer cancel()
	
	// Use errgroup for concurrent requests
	g, gCtx := errgroup.WithContext(testCtx)
	g.SetLimit(concurrentUsers) // Limit concurrent goroutines
	
	// Channel to signal when to stop generating requests
	stopCh := make(chan struct{})
	
	// Start timer to stop test after duration
	go func() {
		time.Sleep(duration)
		close(stopCh)
	}()
	
	// Generate requests continuously until duration expires
	for i := 0; i < concurrentUsers; i++ {
		g.Go(func() error {
			for {
				select {
				case <-stopCh:
					return nil
				case <-gCtx.Done():
					return gCtx.Err()
				default:
					reqResult := p.makeRequest(gCtx, endpoint)
					
					mu.Lock()
					requestResults = append(requestResults, reqResult)
					result.TotalRequests++
					if reqResult.Success {
						result.SuccessfulReqs++
					} else {
						result.FailedReqs++
						result.ErrorTypes[reqResult.Error]++
					}
					mu.Unlock()
					
					// Small delay between requests to simulate realistic user behavior
					time.Sleep(100 * time.Millisecond)
				}
			}
		})
	}
	
	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil && err != context.DeadlineExceeded {
		return result, fmt.Errorf("load test failed: %w", err)
	}
	
	// Calculate statistics
	p.calculateStatistics(&result, requestResults)
	
	return result, nil
}

func (p *PerformanceTester) makeRequest(ctx context.Context, endpoint string) RequestResult {
	start := time.Now()
	
	url := p.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return RequestResult{
			Duration: time.Since(start),
			Success:  false,
			Error:    "request_creation_error",
		}
	}
	
	resp, err := p.httpClient.Do(req)
	duration := time.Since(start)
	
	if err != nil {
		errorType := "network_error"
		if strings.Contains(err.Error(), "timeout") {
			errorType = "timeout_error"
		} else if strings.Contains(err.Error(), "connection refused") {
			errorType = "connection_error"
		}
		
		return RequestResult{
			Duration: duration,
			Success:  false,
			Error:    errorType,
		}
	}
	defer resp.Body.Close()
	
	success := resp.StatusCode >= 200 && resp.StatusCode < 400
	errorType := ""
	if !success {
		if resp.StatusCode >= 500 {
			errorType = "server_error"
		} else if resp.StatusCode >= 400 {
			errorType = "client_error"
		}
	}
	
	return RequestResult{
		Duration:   duration,
		StatusCode: resp.StatusCode,
		Success:    success,
		Error:      errorType,
	}
}

func (p *PerformanceTester) calculateStatistics(result *LoadTestResult, requestResults []RequestResult) {
	if len(requestResults) == 0 {
		return
	}
	
	// Calculate response time statistics
	var totalDuration time.Duration
	responseTimes := make([]time.Duration, len(requestResults))
	
	result.MinResponseTime = time.Duration(1<<63 - 1) // Max duration
	result.MaxResponseTime = 0
	
	for i, req := range requestResults {
		responseTimes[i] = req.Duration
		totalDuration += req.Duration
		
		if req.Duration < result.MinResponseTime {
			result.MinResponseTime = req.Duration
		}
		if req.Duration > result.MaxResponseTime {
			result.MaxResponseTime = req.Duration
		}
	}
	
	result.AvgResponseTime = totalDuration / time.Duration(len(requestResults))
	result.RequestsPerSec = float64(result.TotalRequests) / result.Duration.Seconds()
	
	// Calculate percentiles
	p.calculatePercentiles(result, responseTimes)
}

func (p *PerformanceTester) calculatePercentiles(result *LoadTestResult, responseTimes []time.Duration) {
	if len(responseTimes) == 0 {
		return
	}
	
	// Simple bubble sort for percentile calculation
	for i := 0; i < len(responseTimes); i++ {
		for j := 0; j < len(responseTimes)-1-i; j++ {
			if responseTimes[j] > responseTimes[j+1] {
				responseTimes[j], responseTimes[j+1] = responseTimes[j+1], responseTimes[j]
			}
		}
	}
	
	percentiles := map[string]float64{
		"p50":  0.50,
		"p75":  0.75,
		"p90":  0.90,
		"p95":  0.95,
		"p99":  0.99,
		"p999": 0.999,
	}
	
	for name, pct := range percentiles {
		index := int(float64(len(responseTimes)) * pct)
		if index >= len(responseTimes) {
			index = len(responseTimes) - 1
		}
		result.Percentiles[name] = responseTimes[index]
	}
}

func (p *PerformanceTester) runStressTest(ctx context.Context) error {
	log.Info().Msg("Running stress test - gradually increasing load")
	
	// Gradually increase load to find breaking point
	for users := 50; users <= 1000; users += 50 {
		log.Info().Int("concurrent_users", users).Msg("Testing with increasing load")
		
		result, err := p.runLoadTest(ctx, fmt.Sprintf("Stress Test - %d users", users), "/api/products", users, 30*time.Second)
		if err != nil {
			return err
		}
		
		p.results = append(p.results, result)
		
		// Check if error rate is too high (>5%) or average response time is too slow (>2s)
		errorRate := float64(result.FailedReqs) / float64(result.TotalRequests) * 100
		if errorRate > 5.0 || result.AvgResponseTime > 2*time.Second {
			log.Warn().
				Int("users", users).
				Float64("error_rate", errorRate).
				Dur("avg_response_time", result.AvgResponseTime).
				Msg("System showing stress - breaking point reached")
			break
		}
	}
	
	return nil
}

func (p *PerformanceTester) runSpikeTest(ctx context.Context) error {
	log.Info().Msg("Running spike test - sudden traffic surge")
	
	// Normal load followed by sudden spike
	normalResult, err := p.runLoadTest(ctx, "Spike Test - Normal Load", "/api/products", 25, 30*time.Second)
	if err != nil {
		return err
	}
	p.results = append(p.results, normalResult)
	
	// Sudden spike
	spikeResult, err := p.runLoadTest(ctx, "Spike Test - Traffic Spike", "/api/products", 500, 15*time.Second)
	if err != nil {
		return err
	}
	p.results = append(p.results, spikeResult)
	
	// Recovery
	recoveryResult, err := p.runLoadTest(ctx, "Spike Test - Recovery", "/api/products", 25, 30*time.Second)
	if err != nil {
		return err
	}
	p.results = append(p.results, recoveryResult)
	
	return nil
}

func (p *PerformanceTester) runEnduranceTest(ctx context.Context) error {
	log.Info().Msg("Running endurance test - extended duration")
	
	// Test system stability over extended period
	result, err := p.runLoadTest(ctx, "Endurance Test", "/api/products", 50, 600*time.Second) // 10 minutes
	if err != nil {
		return err
	}
	
	p.results = append(p.results, result)
	return nil
}

func (p *PerformanceTester) logTestResult(result LoadTestResult) {
	log.Info().
		Str("test", result.TestName).
		Int("total_requests", result.TotalRequests).
		Int("successful", result.SuccessfulReqs).
		Int("failed", result.FailedReqs).
		Float64("success_rate", float64(result.SuccessfulReqs)/float64(result.TotalRequests)*100).
		Dur("avg_response_time", result.AvgResponseTime).
		Dur("p95_response_time", result.Percentiles["p95"]).
		Float64("requests_per_sec", result.RequestsPerSec).
		Msg("Performance test completed")
}

func (p *PerformanceTester) generateReport() {
	log.Info().Msg("Generating performance test report")
	
	totalTests := len(p.results)
	totalRequests := 0
	totalSuccessful := 0
	totalFailed := 0
	
	for _, result := range p.results {
		totalRequests += result.TotalRequests
		totalSuccessful += result.SuccessfulReqs
		totalFailed += result.FailedReqs
	}
	
	overallSuccessRate := float64(totalSuccessful) / float64(totalRequests) * 100
	
	log.Info().
		Int("total_tests", totalTests).
		Int("total_requests", totalRequests).
		Int("total_successful", totalSuccessful).
		Int("total_failed", totalFailed).
		Float64("overall_success_rate", overallSuccessRate).
		Msg("Performance test summary")
	
	// Identify performance issues
	p.analyzePerformanceIssues()
	
	// Save detailed results
	reportData := map[string]interface{}{
		"summary": map[string]interface{}{
			"total_tests":           totalTests,
			"total_requests":        totalRequests,
			"total_successful":      totalSuccessful,
			"total_failed":          totalFailed,
			"overall_success_rate":  overallSuccessRate,
		},
		"detailed_results": p.results,
		"test_timestamp":   time.Now().Format(time.RFC3339),
		"performance_analysis": p.getPerformanceAnalysis(),
	}
	
	if jsonData, err := json.MarshalIndent(reportData, "", "  "); err == nil {
		log.Info().Str("report", string(jsonData)).Msg("Detailed performance test report")
	}
}

func (p *PerformanceTester) analyzePerformanceIssues() {
	log.Info().Msg("Analyzing performance issues")
	
	for _, result := range p.results {
		issues := make([]string, 0)
		
		// Check response time issues
		if result.AvgResponseTime > 1*time.Second {
			issues = append(issues, "High average response time")
		}
		
		if result.Percentiles["p95"] > 2*time.Second {
			issues = append(issues, "High P95 response time")
		}
		
		// Check error rate issues
		errorRate := float64(result.FailedReqs) / float64(result.TotalRequests) * 100
		if errorRate > 1.0 {
			issues = append(issues, fmt.Sprintf("High error rate: %.2f%%", errorRate))
		}
		
		// Check throughput issues
		if result.RequestsPerSec < 10.0 {
			issues = append(issues, "Low throughput")
		}
		
		if len(issues) > 0 {
			log.Warn().
				Str("test", result.TestName).
				Strs("issues", issues).
				Msg("Performance issues detected")
		}
	}
}

func (p *PerformanceTester) getPerformanceAnalysis() map[string]interface{} {
	analysis := map[string]interface{}{
		"recommendations": []string{},
		"concerns": []string{},
		"benchmarks": map[string]interface{}{},
	}
	
	// Add performance recommendations based on results
	recommendations := []string{
		"Consider implementing response caching for product catalog",
		"Add database connection pooling for better concurrency",
		"Implement rate limiting to prevent abuse",
		"Consider CDN for static assets",
		"Add monitoring and alerting for response times",
		"Optimize database queries with proper indexing",
		"Consider horizontal scaling for high traffic periods",
	}
	
	analysis["recommendations"] = recommendations
	
	return analysis
}
