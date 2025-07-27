package integration

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

// PaymentIntegrationTestSuite tests payment processing workflows with mock payment providers
type PaymentIntegrationTestSuite struct {
	suite.Suite
	server         *httptest.Server
	client         *http.Client
	sessionID      string
	mockStripe     *MockStripeServer
	mockPayPal     *MockPayPalServer
}

type MockStripeServer struct {
	server *httptest.Server
	client *http.Client
}

type MockPayPalServer struct {
	server *httptest.Server
	client *http.Client
}

type PaymentRequest struct {
	Method       string                 `json:"method"`
	Amount       float64                `json:"amount"`
	Currency     string                 `json:"currency"`
	CardToken    string                 `json:"card_token,omitempty"`
	PayPalToken  string                 `json:"paypal_token,omitempty"`
	BillingInfo  fixtures.BillingInfo   `json:"billing_info"`
	OrderID      string                 `json:"order_id"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type PaymentResponse struct {
	TransactionID string                 `json:"transaction_id"`
	Status        string                 `json:"status"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	Fees          float64                `json:"fees,omitempty"`
	ErrorCode     string                 `json:"error_code,omitempty"`
	ErrorMessage  string                 `json:"error_message,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
}

func (suite *PaymentIntegrationTestSuite) SetupSuite() {
	// Setup mock payment providers
	suite.mockStripe = suite.createMockStripeServer()
	suite.mockPayPal = suite.createMockPayPalServer()
	
	// Setup main server with payment provider endpoints
	suite.server = httptest.NewServer(suite.createTestRouterWithPayments())
	suite.client = &http.Client{
		Timeout: 30 * time.Second,
	}
}

func (suite *PaymentIntegrationTestSuite) TearDownSuite() {
	suite.server.Close()
	suite.mockStripe.server.Close()
	suite.mockPayPal.server.Close()
}

func (suite *PaymentIntegrationTestSuite) SetupTest() {
	suite.sessionID = suite.createTestSession()
	suite.seedTestData()
}

func (suite *PaymentIntegrationTestSuite) TearDownTest() {
	suite.cleanupTestData()
}

// Credit Card Payment Tests

func (suite *PaymentIntegrationTestSuite) TestCreditCardPayment_Success() {
	ctx := context.Background()
	
	// Setup order
	orderID := suite.createTestOrder(ctx, 149.99)
	
	// Process payment
	paymentReq := PaymentRequest{
		Method:   "credit_card",
		Amount:   149.99,
		Currency: "USD",
		CardToken: "test_card_visa_4242424242424242",
		BillingInfo: fixtures.BillingInfo{
			FirstName: "John",
			LastName:  "Doe",
			Address:   "123 Main St",
			City:      "Anytown",
			State:     "CA",
			ZipCode:   "12345",
			Country:   "US",
		},
		OrderID: orderID,
	}

	response := suite.processPayment(ctx, paymentReq)
	
	assert.Equal(suite.T(), "succeeded", response.Status)
	assert.NotEmpty(suite.T(), response.TransactionID)
	assert.Equal(suite.T(), 149.99, response.Amount)
	assert.Equal(suite.T(), "USD", response.Currency)
	assert.Greater(suite.T(), response.Fees, 0.0) // Should have processing fees
	
	// Verify order status updated
	orderStatus := suite.getOrderStatus(ctx, orderID)
	assert.Equal(suite.T(), "paid", orderStatus)
}

func (suite *PaymentIntegrationTestSuite) TestCreditCardPayment_DeclinedCard() {
	ctx := context.Background()
	
	orderID := suite.createTestOrder(ctx, 75.50)
	
	paymentReq := PaymentRequest{
		Method:    "credit_card",
		Amount:    75.50,
		Currency:  "USD",
		CardToken: "test_card_declined", // Mock declined card
		BillingInfo: fixtures.BillingInfo{
			FirstName: "Jane",
			LastName:  "Smith",
		},
		OrderID: orderID,
	}

	response := suite.processPayment(ctx, paymentReq)
	
	assert.Equal(suite.T(), "failed", response.Status)
	assert.Equal(suite.T(), "card_declined", response.ErrorCode)
	assert.NotEmpty(suite.T(), response.ErrorMessage)
	
	// Verify order status remains pending
	orderStatus := suite.getOrderStatus(ctx, orderID)
	assert.Equal(suite.T(), "pending_payment", orderStatus)
}

func (suite *PaymentIntegrationTestSuite) TestCreditCardPayment_InsufficientFunds() {
	ctx := context.Background()
	
	orderID := suite.createTestOrder(ctx, 999.99)
	
	paymentReq := PaymentRequest{
		Method:    "credit_card",
		Amount:    999.99,
		Currency:  "USD",
		CardToken: "test_card_insufficient_funds",
		BillingInfo: fixtures.BillingInfo{
			FirstName: "Bob",
			LastName:  "Wilson",
		},
		OrderID: orderID,
	}

	response := suite.processPayment(ctx, paymentReq)
	
	assert.Equal(suite.T(), "failed", response.Status)
	assert.Equal(suite.T(), "insufficient_funds", response.ErrorCode)
}

func (suite *PaymentIntegrationTestSuite) TestCreditCardPayment_ExpiredCard() {
	ctx := context.Background()
	
	orderID := suite.createTestOrder(ctx, 299.99)
	
	paymentReq := PaymentRequest{
		Method:    "credit_card",
		Amount:    299.99,
		Currency:  "USD",
		CardToken: "test_card_expired",
		BillingInfo: fixtures.BillingInfo{},
		OrderID:   orderID,
	}

	response := suite.processPayment(ctx, paymentReq)
	
	assert.Equal(suite.T(), "failed", response.Status)
	assert.Equal(suite.T(), "expired_card", response.ErrorCode)
}

// PayPal Payment Tests

func (suite *PaymentIntegrationTestSuite) TestPayPalPayment_Success() {
	ctx := context.Background()
	
	orderID := suite.createTestOrder(ctx, 199.99)
	
	// First get PayPal approval URL
	approvalURL := suite.initiatePayPalPayment(ctx, orderID, 199.99)
	assert.NotEmpty(suite.T(), approvalURL)
	
	// Simulate user approval (mock)
	paypalToken := suite.simulatePayPalApproval(ctx, approvalURL)
	
	// Complete payment
	paymentReq := PaymentRequest{
		Method:      "paypal",
		Amount:      199.99,
		Currency:    "USD",
		PayPalToken: paypalToken,
		OrderID:     orderID,
	}

	response := suite.processPayment(ctx, paymentReq)
	
	assert.Equal(suite.T(), "succeeded", response.Status)
	assert.NotEmpty(suite.T(), response.TransactionID)
	assert.Equal(suite.T(), 199.99, response.Amount)
}

func (suite *PaymentIntegrationTestSuite) TestPayPalPayment_UserCancelled() {
	ctx := context.Background()
	
	orderID := suite.createTestOrder(ctx, 89.99)
	
	approvalURL := suite.initiatePayPalPayment(ctx, orderID, 89.99)
	
	// Simulate user cancellation
	paymentReq := PaymentRequest{
		Method:      "paypal",
		Amount:      89.99,
		Currency:    "USD",
		PayPalToken: "cancelled_token",
		OrderID:     orderID,
	}

	response := suite.processPayment(ctx, paymentReq)
	
	assert.Equal(suite.T(), "cancelled", response.Status)
	assert.Equal(suite.T(), "user_cancelled", response.ErrorCode)
}

// Payment Validation Tests

func (suite *PaymentIntegrationTestSuite) TestPaymentValidation_AmountMismatch() {
	ctx := context.Background()
	
	orderID := suite.createTestOrder(ctx, 150.00)
	
	// Try to pay different amount than order total
	paymentReq := PaymentRequest{
		Method:    "credit_card",
		Amount:    100.00, // Different from order amount
		Currency:  "USD",
		CardToken: "test_card_visa_4242424242424242",
		OrderID:   orderID,
	}

	response := suite.processPayment(ctx, paymentReq)
	
	assert.Equal(suite.T(), "failed", response.Status)
	assert.Equal(suite.T(), "amount_mismatch", response.ErrorCode)
}

func (suite *PaymentIntegrationTestSuite) TestPaymentValidation_InvalidCurrency() {
	ctx := context.Background()
	
	orderID := suite.createTestOrder(ctx, 75.00)
	
	paymentReq := PaymentRequest{
		Method:    "credit_card",
		Amount:    75.00,
		Currency:  "INVALID", // Invalid currency
		CardToken: "test_card_visa_4242424242424242",
		OrderID:   orderID,
	}

	response := suite.processPayment(ctx, paymentReq)
	
	assert.Equal(suite.T(), "failed", response.Status)
	assert.Equal(suite.T(), "invalid_currency", response.ErrorCode)
}

func (suite *PaymentIntegrationTestSuite) TestPaymentValidation_DuplicatePayment() {
	ctx := context.Background()
	
	orderID := suite.createTestOrder(ctx, 250.00)
	
	paymentReq := PaymentRequest{
		Method:    "credit_card",
		Amount:    250.00,
		Currency:  "USD",
		CardToken: "test_card_visa_4242424242424242",
		OrderID:   orderID,
	}

	// First payment should succeed
	response1 := suite.processPayment(ctx, paymentReq)
	assert.Equal(suite.T(), "succeeded", response1.Status)
	
	// Second payment should be rejected
	response2 := suite.processPayment(ctx, paymentReq)
	assert.Equal(suite.T(), "failed", response2.Status)
	assert.Equal(suite.T(), "duplicate_payment", response2.ErrorCode)
}

// Refund Tests

func (suite *PaymentIntegrationTestSuite) TestRefund_FullRefund() {
	ctx := context.Background()
	
	// First make a successful payment
	orderID := suite.createTestOrder(ctx, 399.99)
	
	paymentReq := PaymentRequest{
		Method:    "credit_card",
		Amount:    399.99,
		Currency:  "USD",
		CardToken: "test_card_visa_4242424242424242",
		OrderID:   orderID,
	}

	paymentResponse := suite.processPayment(ctx, paymentReq)
	require.Equal(suite.T(), "succeeded", paymentResponse.Status)
	
	// Now process refund
	refundReq := struct {
		TransactionID string  `json:"transaction_id"`
		Amount        float64 `json:"amount"`
		Reason        string  `json:"reason"`
	}{
		TransactionID: paymentResponse.TransactionID,
		Amount:        399.99,
		Reason:        "customer_request",
	}

	refundResponse := suite.processRefund(ctx, refundReq)
	
	assert.Equal(suite.T(), "succeeded", refundResponse.Status)
	assert.Equal(suite.T(), 399.99, refundResponse.Amount)
	assert.NotEmpty(suite.T(), refundResponse.TransactionID)
}

func (suite *PaymentIntegrationTestSuite) TestRefund_PartialRefund() {
	ctx := context.Background()
	
	// Make payment
	orderID := suite.createTestOrder(ctx, 500.00)
	
	paymentResponse := suite.processPayment(ctx, PaymentRequest{
		Method:    "credit_card",
		Amount:    500.00,
		Currency:  "USD",
		CardToken: "test_card_visa_4242424242424242",
		OrderID:   orderID,
	})
	require.Equal(suite.T(), "succeeded", paymentResponse.Status)
	
	// Partial refund
	refundReq := struct {
		TransactionID string  `json:"transaction_id"`
		Amount        float64 `json:"amount"`
		Reason        string  `json:"reason"`
	}{
		TransactionID: paymentResponse.TransactionID,
		Amount:        150.00, // Partial amount
		Reason:        "damaged_item",
	}

	refundResponse := suite.processRefund(ctx, refundReq)
	
	assert.Equal(suite.T(), "succeeded", refundResponse.Status)
	assert.Equal(suite.T(), 150.00, refundResponse.Amount)
}

// Webhook Tests

func (suite *PaymentIntegrationTestSuite) TestWebhook_PaymentSucceeded() {
	ctx := context.Background()
	
	webhookPayload := map[string]interface{}{
		"type": "payment.succeeded",
		"data": map[string]interface{}{
			"transaction_id": "test_tx_123",
			"order_id":       "test_order_456",
			"amount":         299.99,
			"currency":       "USD",
		},
	}

	response := suite.sendWebhook(ctx, "stripe", webhookPayload)
	assert.Equal(suite.T(), http.StatusOK, response.StatusCode)
	
	// Verify order status was updated
	orderStatus := suite.getOrderStatus(ctx, "test_order_456")
	assert.Equal(suite.T(), "paid", orderStatus)
}

func (suite *PaymentIntegrationTestSuite) TestWebhook_PaymentFailed() {
	ctx := context.Background()
	
	webhookPayload := map[string]interface{}{
		"type": "payment.failed",
		"data": map[string]interface{}{
			"transaction_id": "test_tx_124",
			"order_id":       "test_order_457",
			"error_code":     "insufficient_funds",
			"error_message":  "Your card has insufficient funds",
		},
	}

	response := suite.sendWebhook(ctx, "stripe", webhookPayload)
	assert.Equal(suite.T(), http.StatusOK, response.StatusCode)
	
	// Verify order status
	orderStatus := suite.getOrderStatus(ctx, "test_order_457")
	assert.Equal(suite.T(), "payment_failed", orderStatus)
}

// Currency and Internationalization Tests

func (suite *PaymentIntegrationTestSuite) TestPayment_EURCurrency() {
	ctx := context.Background()
	
	orderID := suite.createTestOrderWithCurrency(ctx, 99.99, "EUR")
	
	paymentReq := PaymentRequest{
		Method:    "credit_card",
		Amount:    99.99,
		Currency:  "EUR",
		CardToken: "test_card_visa_4242424242424242",
		OrderID:   orderID,
	}

	response := suite.processPayment(ctx, paymentReq)
	
	assert.Equal(suite.T(), "succeeded", response.Status)
	assert.Equal(suite.T(), "EUR", response.Currency)
}

func (suite *PaymentIntegrationTestSuite) TestPayment_GBPCurrency() {
	ctx := context.Background()
	
	orderID := suite.createTestOrderWithCurrency(ctx, 75.50, "GBP")
	
	paymentReq := PaymentRequest{
		Method:    "credit_card",
		Amount:    75.50,
		Currency:  "GBP",
		CardToken: "test_card_visa_4242424242424242",
		OrderID:   orderID,
	}

	response := suite.processPayment(ctx, paymentReq)
	
	assert.Equal(suite.T(), "succeeded", response.Status)
	assert.Equal(suite.T(), "GBP", response.Currency)
}

// Security Tests

func (suite *PaymentIntegrationTestSuite) TestPaymentSecurity_TokenValidation() {
	ctx := context.Background()
	
	orderID := suite.createTestOrder(ctx, 100.00)
	
	// Try with invalid card token
	paymentReq := PaymentRequest{
		Method:    "credit_card",
		Amount:    100.00,
		Currency:  "USD",
		CardToken: "invalid_token_format",
		OrderID:   orderID,
	}

	response := suite.processPayment(ctx, paymentReq)
	
	assert.Equal(suite.T(), "failed", response.Status)
	assert.Equal(suite.T(), "invalid_token", response.ErrorCode)
}

func (suite *PaymentIntegrationTestSuite) TestPaymentSecurity_SignatureValidation() {
	ctx := context.Background()
	
	// Test webhook with invalid signature
	webhookPayload := map[string]interface{}{
		"type": "payment.succeeded",
		"data": map[string]interface{}{
			"transaction_id": "test_tx_125",
		},
	}

	response := suite.sendWebhookWithInvalidSignature(ctx, webhookPayload)
	assert.Equal(suite.T(), http.StatusUnauthorized, response.StatusCode)
}

// Helper Methods

func (suite *PaymentIntegrationTestSuite) createMockStripeServer() *MockStripeServer {
	handler := http.NewServeMux()
	
	// Mock charge endpoint
	handler.HandleFunc("/v1/charges", func(w http.ResponseWriter, r *http.Request) {
		// Simulate different card behaviors based on token
		cardToken := r.FormValue("source")
		
		response := map[string]interface{}{}
		
		switch cardToken {
		case "test_card_visa_4242424242424242":
			response = map[string]interface{}{
				"id":       "ch_test_123",
				"status":   "succeeded",
				"amount":   r.FormValue("amount"),
				"currency": r.FormValue("currency"),
			}
		case "test_card_declined":
			w.WriteHeader(http.StatusPaymentRequired)
			response = map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "card_declined",
					"message": "Your card was declined.",
				},
			}
		case "test_card_insufficient_funds":
			w.WriteHeader(http.StatusPaymentRequired)
			response = map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "insufficient_funds",
					"message": "Your card has insufficient funds.",
				},
			}
		case "test_card_expired":
			w.WriteHeader(http.StatusPaymentRequired)
			response = map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "expired_card",
					"message": "Your card has expired.",
				},
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			response = map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "invalid_token",
					"message": "Invalid card token.",
				},
			}
		}
		
		json.NewEncoder(w).Encode(response)
	})

	server := httptest.NewServer(handler)
	return &MockStripeServer{
		server: server,
		client: &http.Client{},
	}
}

func (suite *PaymentIntegrationTestSuite) createMockPayPalServer() *MockPayPalServer {
	handler := http.NewServeMux()
	
	// Mock PayPal payment creation
	handler.HandleFunc("/v1/payments/payment", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"id": "PAY-test123",
			"links": []map[string]interface{}{
				{
					"href":   "https://www.sandbox.paypal.com/cgi-bin/webscr?cmd=_express-checkout&token=test_token",
					"rel":    "approval_url",
					"method": "REDIRECT",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	})

	server := httptest.NewServer(handler)
	return &MockPayPalServer{
		server: server,
		client: &http.Client{},
	}
}

func (suite *PaymentIntegrationTestSuite) createTestRouterWithPayments() http.Handler {
	// Create router with payment endpoints
	mux := http.NewServeMux()
	
	// Payment processing endpoint
	mux.HandleFunc("/api/payments", suite.handlePayment)
	mux.HandleFunc("/api/refunds", suite.handleRefund)
	mux.HandleFunc("/api/webhooks/stripe", suite.handleStripeWebhook)
	mux.HandleFunc("/api/webhooks/paypal", suite.handlePayPalWebhook)
	
	return mux
}

func (suite *PaymentIntegrationTestSuite) handlePayment(w http.ResponseWriter, r *http.Request) {
	// Mock payment processing logic
	var req PaymentRequest
	json.NewDecoder(r.Body).Decode(&req)
	
	// Simulate payment processing
	response := PaymentResponse{
		TransactionID: fmt.Sprintf("tx_%d", time.Now().Unix()),
		Status:        "succeeded",
		Amount:        req.Amount,
		Currency:      req.Currency,
		Fees:          req.Amount * 0.029, // 2.9% fee
		CreatedAt:     time.Now(),
	}
	
	json.NewEncoder(w).Encode(response)
}

func (suite *PaymentIntegrationTestSuite) handleRefund(w http.ResponseWriter, r *http.Request) {
	// Mock refund processing
	var req struct {
		TransactionID string  `json:"transaction_id"`
		Amount        float64 `json:"amount"`
		Reason        string  `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	
	response := struct {
		TransactionID string    `json:"transaction_id"`
		Status        string    `json:"status"`
		Amount        float64   `json:"amount"`
		CreatedAt     time.Time `json:"created_at"`
	}{
		TransactionID: fmt.Sprintf("re_%d", time.Now().Unix()),
		Status:        "succeeded",
		Amount:        req.Amount,
		CreatedAt:     time.Now(),
	}
	
	json.NewEncoder(w).Encode(response)
}

func (suite *PaymentIntegrationTestSuite) handleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	// Mock webhook processing
	w.WriteHeader(http.StatusOK)
}

func (suite *PaymentIntegrationTestSuite) handlePayPalWebhook(w http.ResponseWriter, r *http.Request) {
	// Mock webhook processing
	w.WriteHeader(http.StatusOK)
}

func (suite *PaymentIntegrationTestSuite) processPayment(ctx context.Context, req PaymentRequest) PaymentResponse {
	jsonData, _ := json.Marshal(req)
	url := fmt.Sprintf("%s/api/payments", suite.server.URL)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	
	resp, _ := suite.client.Do(httpReq)
	defer resp.Body.Close()
	
	var response PaymentResponse
	json.NewDecoder(resp.Body).Decode(&response)
	return response
}

func (suite *PaymentIntegrationTestSuite) processRefund(ctx context.Context, req interface{}) struct {
	TransactionID string    `json:"transaction_id"`
	Status        string    `json:"status"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
} {
	jsonData, _ := json.Marshal(req)
	url := fmt.Sprintf("%s/api/refunds", suite.server.URL)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	
	resp, _ := suite.client.Do(httpReq)
	defer resp.Body.Close()
	
	var response struct {
		TransactionID string    `json:"transaction_id"`
		Status        string    `json:"status"`
		Amount        float64   `json:"amount"`
		CreatedAt     time.Time `json:"created_at"`
	}
	json.NewDecoder(resp.Body).Decode(&response)
	return response
}

func (suite *PaymentIntegrationTestSuite) sendWebhook(ctx context.Context, provider string, payload map[string]interface{}) *http.Response {
	jsonData, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/api/webhooks/%s", suite.server.URL, provider)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := suite.client.Do(req)
	return resp
}

func (suite *PaymentIntegrationTestSuite) sendWebhookWithInvalidSignature(ctx context.Context, payload map[string]interface{}) *http.Response {
	jsonData, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/api/webhooks/stripe", suite.server.URL)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Stripe-Signature", "invalid_signature")
	
	resp, _ := suite.client.Do(req)
	return resp
}

func (suite *PaymentIntegrationTestSuite) createTestOrder(ctx context.Context, amount float64) string {
	return suite.createTestOrderWithCurrency(ctx, amount, "USD")
}

func (suite *PaymentIntegrationTestSuite) createTestOrderWithCurrency(ctx context.Context, amount float64, currency string) string {
	return fmt.Sprintf("order_%d", time.Now().Unix())
}

func (suite *PaymentIntegrationTestSuite) getOrderStatus(ctx context.Context, orderID string) string {
	// Mock order status lookup
	return "pending_payment"
}

func (suite *PaymentIntegrationTestSuite) initiatePayPalPayment(ctx context.Context, orderID string, amount float64) string {
	// Mock PayPal payment initiation
	return "https://www.sandbox.paypal.com/cgi-bin/webscr?cmd=_express-checkout&token=test_token"
}

func (suite *PaymentIntegrationTestSuite) simulatePayPalApproval(ctx context.Context, approvalURL string) string {
	// Mock user approval
	return "approved_paypal_token"
}

func (suite *PaymentIntegrationTestSuite) createTestSession() string {
	return fmt.Sprintf("payment-test-session-%d", time.Now().UnixNano())
}

func (suite *PaymentIntegrationTestSuite) seedTestData() {
	// Seed test data
}

func (suite *PaymentIntegrationTestSuite) cleanupTestData() {
	// Cleanup test data
}

// Run the test suite
func TestPaymentIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentIntegrationTestSuite))
}
