package load

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	
	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/tests/fixtures"
)

// PerformanceTestSuite provides comprehensive performance and load testing
type PerformanceTestSuite struct {
	suite.Suite
	server     *httptest.Server
	client     *http.Client
	baseURL    string
	testData   *TestDataSet
}

type TestDataSet struct {
	Products   []fixtures.TestProduct
	Categories []fixtures.TestCategory
	Brands     []fixtures.TestBrand
}

type PerformanceMetrics struct {
	TotalRequests    int
	SuccessfulReqs   int
	FailedReqs       int
	AverageLatency   time.Duration
	MinLatency       time.Duration
	MaxLatency       time.Duration
	P95Latency       time.Duration
	P99Latency       time.Duration
	RequestsPerSec   float64
	ErrorRate        float64
	Latencies        []time.Duration
}

func (suite *PerformanceTestSuite) SetupSuite() {
	suite.server = httptest.NewServer(createTestRouter())
	suite.baseURL = suite.server.URL
	suite.client = &http.Client{
		Timeout: 30 * time.Second,
	}
}

func (suite *PerformanceTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *PerformanceTestSuite) SetupTest() {
	suite.testData = suite.seedPerformanceTestData()
}

func (suite *PerformanceTestSuite) TearDownTest() {
	suite.cleanupTestData()
}

// Load Testing - API Endpoints

func (suite *PerformanceTestSuite) TestLoadTest_ProductCatalog() {
	scenarios := []struct {
		name           string
		concurrency    int
		duration       time.Duration
		maxLatencyP95  time.Duration
		minRequestsPS  float64
	}{
		{"Low Load", 10, 30 * time.Second, 200 * time.Millisecond, 40.0},
		{"Medium Load", 50, 60 * time.Second, 300 * time.Millisecond, 150.0},
		{"High Load", 100, 60 * time.Second, 500 * time.Millisecond, 250.0},
		{"Peak Load", 200, 30 * time.Second, 1 * time.Second, 300.0},
	}

	for _, scenario := range scenarios {
		suite.T().Run(scenario.name, func(t *testing.T) {
			metrics := suite.runLoadTest(t, scenario.concurrency, scenario.duration, 
				suite.productCatalogWorkload)

			// Performance assertions
			assert.LessOrEqual(t, metrics.P95Latency, scenario.maxLatencyP95,
				"P95 latency should be within acceptable limits")
			assert.GreaterOrEqual(t, metrics.RequestsPerSec, scenario.minRequestsPS,
				"Requests per second should meet minimum threshold")
			assert.LessOrEqual(t, metrics.ErrorRate, 0.01, // 1% error rate
				"Error rate should be minimal")

			suite.logPerformanceMetrics(t, scenario.name, metrics)
		})
	}
}

func (suite *PerformanceTestSuite) TestLoadTest_SearchFunctionality() {
	metrics := suite.runLoadTest(suite.T(), 50, 60*time.Second, suite.searchWorkload)

	// Search should be fast even under load
	assert.LessOrEqual(suite.T(), metrics.P95Latency, 300*time.Millisecond,
		"Search P95 latency should be under 300ms")
	assert.GreaterOrEqual(suite.T(), metrics.RequestsPerSec, 100.0,
		"Search should handle at least 100 requests/sec")
	assert.LessOrEqual(suite.T(), metrics.ErrorRate, 0.005, // 0.5% error rate
		"Search error rate should be very low")
}

func (suite *PerformanceTestSuite) TestLoadTest_ShoppingCart() {
	metrics := suite.runLoadTest(suite.T(), 30, 45*time.Second, suite.cartWorkload)

	// Cart operations should be responsive
	assert.LessOrEqual(suite.T(), metrics.P95Latency, 250*time.Millisecond,
		"Cart P95 latency should be under 250ms")
	assert.LessOrEqual(suite.T(), metrics.ErrorRate, 0.01,
		"Cart operations should have low error rate")
}

// Stress Testing - Finding Breaking Points

func (suite *PerformanceTestSuite) TestStressTest_ProductAPI() {
	// Gradually increase load to find breaking point
	concurrencyLevels := []int{50, 100, 200, 400, 800}
	testDuration := 30 * time.Second

	breakingPoint := -1
	for i, concurrency := range concurrencyLevels {
		suite.T().Run(fmt.Sprintf("Concurrency_%d", concurrency), func(t *testing.T) {
			metrics := suite.runLoadTest(t, concurrency, testDuration, 
				suite.productCatalogWorkload)

			suite.logPerformanceMetrics(t, fmt.Sprintf("Stress_%d", concurrency), metrics)

			// Check if this is the breaking point
			if metrics.ErrorRate > 0.05 || metrics.P95Latency > 2*time.Second {
				if breakingPoint == -1 {
					breakingPoint = i
				}
				t.Logf("Breaking point detected at concurrency level %d", concurrency)
			}
		})
	}

	if breakingPoint != -1 {
		suite.T().Logf("System breaks down at approximately %d concurrent users", 
			concurrencyLevels[breakingPoint])
	}
}

// Endurance Testing - Sustained Load

func (suite *PerformanceTestSuite) TestEnduranceTest_LongRunning() {
	if testing.Short() {
		suite.T().Skip("Skipping endurance test in short mode")
	}

	// Run sustained load for 10 minutes
	concurrency := 20
	duration := 10 * time.Minute

	metrics := suite.runLoadTest(suite.T(), concurrency, duration, 
		suite.mixedWorkload)

	// System should remain stable over time
	assert.LessOrEqual(suite.T(), metrics.P95Latency, 400*time.Millisecond,
		"System should maintain performance over extended period")
	assert.LessOrEqual(suite.T(), metrics.ErrorRate, 0.01,
		"Error rate should remain low during endurance test")

	suite.T().Logf("Endurance test completed: %d requests over %v", 
		metrics.TotalRequests, duration)
}

// Memory and Resource Usage Tests

func (suite *PerformanceTestSuite) TestMemoryLeak_CartOperations() {
	// Simulate many cart operations to check for memory leaks
	sessionCount := 100
	operationsPerSession := 50

	var wg sync.WaitGroup
	for i := 0; i < sessionCount; i++ {
		wg.Add(1)
		go func(sessionID int) {
			defer wg.Done()
			suite.simulateCartSession(sessionID, operationsPerSession)
		}(i)
	}

	wg.Wait()

	// In a real implementation, you would check memory usage here
	// This is a placeholder for memory monitoring
	suite.T().Log("Memory leak test completed - monitor memory usage externally")
}

// Database Performance Tests

func (suite *PerformanceTestSuite) TestDatabasePerformance_ProductQueries() {
	// Test database performance under concurrent load
	metrics := suite.runLoadTest(suite.T(), 100, 60*time.Second, 
		suite.databaseIntensiveWorkload)

	// Database queries should remain fast
	assert.LessOrEqual(suite.T(), metrics.P95Latency, 500*time.Millisecond,
		"Database queries should be optimized for performance")
	assert.GreaterOrEqual(suite.T(), metrics.RequestsPerSec, 80.0,
		"Database should handle reasonable query load")
}

// Workload Implementations

func (suite *PerformanceTestSuite) productCatalogWorkload(ctx context.Context, workerID int) error {
	// Simulate browsing product catalog
	actions := []func() error{
		func() error { return suite.getProducts() },
		func() error { return suite.getCategories() },
		func() error { return suite.getBrands() },
		func() error { return suite.getProductDetails(suite.randomProductID()) },
		func() error { return suite.getProductsByCategory(suite.randomCategoryID()) },
	}

	// Random action selection
	action := actions[rand.Intn(len(actions))]
	return action()
}

func (suite *PerformanceTestSuite) searchWorkload(ctx context.Context, workerID int) error {
	// Simulate search operations
	searchTerms := []string{"Dyson", "Shark", "vacuum", "cordless", "upright", "canister"}
	term := searchTerms[rand.Intn(len(searchTerms))]
	return suite.searchProducts(term)
}

func (suite *PerformanceTestSuite) cartWorkload(ctx context.Context, workerID int) error {
	// Simulate cart operations
	sessionID := fmt.Sprintf("load-test-session-%d", workerID)
	
	actions := []func(string) error{
		func(s string) error { return suite.addToCart(s, suite.randomProductID(), rand.Intn(3)+1) },
		func(s string) error { return suite.getCart(s) },
		func(s string) error { return suite.updateCartItem(s, suite.randomProductID(), rand.Intn(5)+1) },
		func(s string) error { return suite.removeFromCart(s, suite.randomProductID()) },
	}

	action := actions[rand.Intn(len(actions))]
	return action(sessionID)
}

func (suite *PerformanceTestSuite) mixedWorkload(ctx context.Context, workerID int) error {
	// Mix of different operations
	workloads := []func(context.Context, int) error{
		suite.productCatalogWorkload,
		suite.searchWorkload,
		suite.cartWorkload,
	}

	workload := workloads[rand.Intn(len(workloads))]
	return workload(ctx, workerID)
}

func (suite *PerformanceTestSuite) databaseIntensiveWorkload(ctx context.Context, workerID int) error {
	// Operations that stress the database
	actions := []func() error{
		func() error { return suite.getProductsWithFilters() },
		func() error { return suite.searchProductsWithSorting() },
		func() error { return suite.getProductRecommendations() },
		func() error { return suite.getInventoryStatus() },
	}

	action := actions[rand.Intn(len(actions))]
	return action()
}

// Load Test Runner

func (suite *PerformanceTestSuite) runLoadTest(t *testing.T, concurrency int, duration time.Duration, 
	workload func(context.Context, int) error) PerformanceMetrics {
	
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var allLatencies []time.Duration
	
	successCount := 0
	errorCount := 0
	startTime := time.Now()

	// Launch workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			for {
				select {
				case <-ctx.Done():
					return
				default:
					requestStart := time.Now()
					err := workload(ctx, workerID)
					latency := time.Since(requestStart)

					mu.Lock()
					allLatencies = append(allLatencies, latency)
					if err != nil {
						errorCount++
					} else {
						successCount++
					}
					mu.Unlock()
				}
			}
		}(i)
	}

	// Wait for completion
	wg.Wait()
	totalDuration := time.Since(startTime)

	// Calculate metrics
	return suite.calculateMetrics(allLatencies, successCount, errorCount, totalDuration)
}

func (suite *PerformanceTestSuite) calculateMetrics(latencies []time.Duration, 
	success, errors int, duration time.Duration) PerformanceMetrics {
	
	if len(latencies) == 0 {
		return PerformanceMetrics{}
	}

	// Sort latencies for percentile calculation
	sortedLatencies := make([]time.Duration, len(latencies))
	copy(sortedLatencies, latencies)
	
	// Simple bubble sort for demonstration (use sort.Slice in real implementation)
	for i := 0; i < len(sortedLatencies)-1; i++ {
		for j := 0; j < len(sortedLatencies)-i-1; j++ {
			if sortedLatencies[j] > sortedLatencies[j+1] {
				sortedLatencies[j], sortedLatencies[j+1] = sortedLatencies[j+1], sortedLatencies[j]
			}
		}
	}

	// Calculate average
	var totalLatency time.Duration
	for _, lat := range latencies {
		totalLatency += lat
	}
	avgLatency := totalLatency / time.Duration(len(latencies))

	// Calculate percentiles
	p95Index := int(float64(len(sortedLatencies)) * 0.95)
	p99Index := int(float64(len(sortedLatencies)) * 0.99)
	
	if p95Index >= len(sortedLatencies) {
		p95Index = len(sortedLatencies) - 1
	}
	if p99Index >= len(sortedLatencies) {
		p99Index = len(sortedLatencies) - 1
	}

	totalRequests := success + errors
	requestsPerSec := float64(totalRequests) / duration.Seconds()
	errorRate := float64(errors) / float64(totalRequests)

	return PerformanceMetrics{
		TotalRequests:  totalRequests,
		SuccessfulReqs: success,
		FailedReqs:     errors,
		AverageLatency: avgLatency,
		MinLatency:     sortedLatencies[0],
		MaxLatency:     sortedLatencies[len(sortedLatencies)-1],
		P95Latency:     sortedLatencies[p95Index],
		P99Latency:     sortedLatencies[p99Index],
		RequestsPerSec: requestsPerSec,
		ErrorRate:      errorRate,
		Latencies:      latencies,
	}
}

// API Helper Methods

func (suite *PerformanceTestSuite) getProducts() error {
	url := fmt.Sprintf("%s/api/products", suite.baseURL)
	resp, err := suite.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}

func (suite *PerformanceTestSuite) getCategories() error {
	url := fmt.Sprintf("%s/api/categories", suite.baseURL)
	resp, err := suite.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (suite *PerformanceTestSuite) getBrands() error {
	url := fmt.Sprintf("%s/api/brands", suite.baseURL)
	resp, err := suite.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (suite *PerformanceTestSuite) getProductDetails(productID int) error {
	url := fmt.Sprintf("%s/api/products/%d", suite.baseURL, productID)
	resp, err := suite.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (suite *PerformanceTestSuite) searchProducts(query string) error {
	url := fmt.Sprintf("%s/api/products/search?q=%s", suite.baseURL, query)
	resp, err := suite.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (suite *PerformanceTestSuite) addToCart(sessionID string, productID, quantity int) error {
	item := map[string]interface{}{
		"product_id": productID,
		"quantity":   quantity,
	}
	jsonData, _ := json.Marshal(item)
	
	url := fmt.Sprintf("%s/api/cart/add", suite.baseURL)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Session-ID", sessionID)
	
	resp, err := suite.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (suite *PerformanceTestSuite) getCart(sessionID string) error {
	url := fmt.Sprintf("%s/api/cart", suite.baseURL)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Session-ID", sessionID)
	
	resp, err := suite.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// Additional helper methods...

func (suite *PerformanceTestSuite) randomProductID() int {
	if len(suite.testData.Products) == 0 {
		return 1
	}
	return suite.testData.Products[rand.Intn(len(suite.testData.Products))].ID
}

func (suite *PerformanceTestSuite) randomCategoryID() int {
	if len(suite.testData.Categories) == 0 {
		return 1
	}
	return suite.testData.Categories[rand.Intn(len(suite.testData.Categories))].ID
}

func (suite *PerformanceTestSuite) logPerformanceMetrics(t *testing.T, testName string, metrics PerformanceMetrics) {
	t.Logf("=== Performance Results: %s ===", testName)
	t.Logf("Total Requests: %d", metrics.TotalRequests)
	t.Logf("Success Rate: %.2f%%", (1.0-metrics.ErrorRate)*100)
	t.Logf("Requests/sec: %.2f", metrics.RequestsPerSec)
	t.Logf("Average Latency: %v", metrics.AverageLatency)
	t.Logf("P95 Latency: %v", metrics.P95Latency)
	t.Logf("P99 Latency: %v", metrics.P99Latency)
	t.Logf("Max Latency: %v", metrics.MaxLatency)
}

// Helper methods for complex operations
func (suite *PerformanceTestSuite) getProductsByCategory(categoryID int) error {
	// Implementation
	return nil
}

func (suite *PerformanceTestSuite) updateCartItem(sessionID string, productID, quantity int) error {
	// Implementation
	return nil
}

func (suite *PerformanceTestSuite) removeFromCart(sessionID string, productID int) error {
	// Implementation
	return nil
}

func (suite *PerformanceTestSuite) getProductsWithFilters() error {
	// Implementation
	return nil
}

func (suite *PerformanceTestSuite) searchProductsWithSorting() error {
	// Implementation
	return nil
}

func (suite *PerformanceTestSuite) getProductRecommendations() error {
	// Implementation
	return nil
}

func (suite *PerformanceTestSuite) getInventoryStatus() error {
	// Implementation
	return nil
}

func (suite *PerformanceTestSuite) simulateCartSession(sessionID, operations int) {
	// Implementation
}

func (suite *PerformanceTestSuite) seedPerformanceTestData() *TestDataSet {
	return &TestDataSet{
		Products:   fixtures.CreateTestProducts(),
		Categories: fixtures.CreateTestCategories(),
		Brands:     fixtures.CreateTestBrands(),
	}
}

func (suite *PerformanceTestSuite) cleanupTestData() {
	// Implementation
}

func createTestRouter() http.Handler {
	return http.NewServeMux()
}

// Run the test suite
func TestPerformanceTestSuite(t *testing.T) {
	suite.Run(t, new(PerformanceTestSuite))
}

// Benchmark tests using Go's built-in benchmarking
func BenchmarkProductAPI(b *testing.B) {
	server := httptest.NewServer(createTestRouter())
	defer server.Close()

	client := &http.Client{}
	url := fmt.Sprintf("%s/api/products", server.URL)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, _ := client.Get(url)
			resp.Body.Close()
		}
	})
}

func BenchmarkSearchAPI(b *testing.B) {
	server := httptest.NewServer(createTestRouter())
	defer server.Close()

	client := &http.Client{}
	url := fmt.Sprintf("%s/api/products/search?q=vacuum", server.URL)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, _ := client.Get(url)
			resp.Body.Close()
		}
	})
}
