# VacuumMart E-commerce Testing Suite

A comprehensive testing framework for the VacuumMart vacuum cleaner retail platform. This suite provides automated testing across multiple dimensions including API functionality, end-to-end user flows, performance, security, mobile responsiveness, and user acceptance testing.

## 🎯 Overview

This testing suite is designed to validate the complete VacuumMart e-commerce platform which includes:
- Product catalog management for vacuum cleaners (Dyson, Hoover, Roomba, Shark brands)
- Shopping cart and checkout functionality
- User account management
- Admin dashboard and inventory management
- Mobile-responsive web interface
- Payment processing integration

## 🏗️ Architecture

The testing suite is built in Go using a modular architecture:

```
vacuum-cleaner-testing-suite/
├── main.go                    # CLI entry point and command definitions
├── testers.go                 # Factory functions and interfaces
├── api/                       # API testing module
│   └── api_tester.go         # Product catalog, cart, and order API tests
├── e2e/                       # End-to-end testing module
│   └── e2e_tester.go         # Complete shopping flow scenarios
├── performance/               # Performance testing module
│   └── performance_tester.go # Load testing and stress testing
├── security/                  # Security testing module
│   └── security_tester.go    # XSS, SQL injection, and security tests
├── mobile/                    # Mobile testing module
│   └── mobile_tester.go      # Mobile responsiveness and touch testing
└── uat/                       # User Acceptance Testing module
    └── uat_tester.go         # Business requirement validation
```

## 🚀 Features

### 1. API Testing (`vacuum-testing-suite api`)
- **Product Catalog APIs**: Test product listing, filtering, search, and details
- **Shopping Cart APIs**: Add, update, remove items, calculate totals
- **Order Management**: Create orders, track status, order history
- **Search & Filtering**: Category, brand, price, and keyword filtering
- **Admin APIs**: Inventory management, sales analytics

### 2. End-to-End Testing (`vacuum-testing-suite e2e`)
- **Complete Shopping Flow**: Browse → Search → Add to Cart → Checkout
- **Guest Checkout**: Purchase without account creation
- **Product Search to Purchase**: Feature-based product discovery
- **Cart Management**: Add multiple items, update quantities, persist cart
- **User Account Flow**: Registration, login, profile management
- **Order Tracking**: Status updates and history
- **Product Reviews**: Read and write customer feedback
- **Product Comparison**: Side-by-side feature comparison

### 3. Performance Testing (`vacuum-testing-suite performance`)
- **Load Testing**: Concurrent user simulation (50-500 users)
- **Stress Testing**: Find system breaking points
- **Spike Testing**: Handle sudden traffic surges
- **Endurance Testing**: Extended duration stability
- **Response Time Analysis**: P50, P95, P99 percentiles
- **Throughput Measurement**: Requests per second under load

### 4. Security Testing (`vacuum-testing-suite security`)
- **XSS Protection**: Cross-site scripting vulnerability testing
- **SQL Injection**: Database injection attack prevention
- **Input Validation**: Boundary testing and malformed data
- **Authentication Bypass**: Unauthorized access attempts
- **CSRF Protection**: Cross-site request forgery prevention
- **Directory Traversal**: File system access protection
- **HTTP Security**: Security headers and configurations
- **Business Logic Flaws**: Negative quantities, price manipulation

### 5. Mobile Responsiveness Testing (`vacuum-testing-suite mobile`)
- **Multi-Device Testing**: iPhone, Android, iPad compatibility
- **Viewport Responsiveness**: Various screen sizes and orientations
- **Touch Interaction**: Touch-friendly interface validation
- **Mobile Performance**: Load time optimization for mobile networks
- **Mobile-Specific Features**: App store links, mobile navigation
- **Cross-Browser Testing**: Safari, Chrome, Firefox mobile

### 6. User Acceptance Testing (`vacuum-testing-suite uat`)
Based on actual business requirements for VacuumMart:
- **Customer Product Browsing**: Category and brand navigation
- **Vacuum Search & Filtering**: Feature-based product discovery
- **Shopping Cart Functionality**: Complete cart management
- **Checkout Process**: End-to-end purchase completion
- **User Account Management**: Registration and profile management
- **Product Comparison**: Side-by-side vacuum comparison
- **Inventory Management**: Stock level management
- **Order Tracking**: Purchase status monitoring
- **Customer Reviews**: Product feedback system
- **Admin Dashboard**: Management interface validation
- **Sales Reporting**: Business analytics
- **Mobile Shopping Experience**: Mobile commerce validation

## 🛠️ Installation & Setup

### Prerequisites
- Go 1.21 or later
- VacuumMart server running (default: http://localhost:8080)

### Build
```bash
# Navigate to the testing suite directory
cd cmd/experiments/vacuum-cleaner-testing-suite

# Build the testing application
go build -o vacuum-testing-suite .
```

### Configuration
The testing suite connects to the VacuumMart server at `http://localhost:8080` by default. You can specify a different URL:

```bash
./vacuum-testing-suite --base-url http://your-server:port
```

## 📋 Usage

### Run All Tests
```bash
./vacuum-testing-suite all --log-level info
```

### Run Individual Test Suites
```bash
# API Tests
./vacuum-testing-suite api

# End-to-End Tests
./vacuum-testing-suite e2e

# Performance Tests
./vacuum-testing-suite performance

# Security Tests
./vacuum-testing-suite security

# Mobile Tests
./vacuum-testing-suite mobile

# User Acceptance Tests
./vacuum-testing-suite uat
```

### Configuration Options
```bash
# Set log level (trace, debug, info, warn, error)
./vacuum-testing-suite all --log-level debug

# Specify different base URL
./vacuum-testing-suite api --base-url https://staging.vacuummart.com

# Run with verbose output
./vacuum-testing-suite performance --log-level trace
```

## 📊 Test Reports

Each test suite generates comprehensive reports including:

### API Testing Reports
- Total requests and success rates
- Response time statistics (min, max, average, percentiles)
- Error categorization
- Endpoint-specific performance metrics

### Performance Reports
- Load test results with concurrent user simulation
- Response time percentiles (P50, P75, P90, P95, P99)
- Throughput measurements (requests per second)
- Error rates and failure categorization
- Stress test breaking points

### Security Reports
- Vulnerability categorization (Critical, High, Medium, Low)
- Risk scoring based on severity and impact
- Detailed attack vectors and recommendations
- Compliance with security best practices

### Mobile Reports
- Device-specific compatibility results
- Performance metrics across different devices
- Touch interaction validation
- Responsive design compliance

### UAT Reports
- Business requirement validation
- User story completion rates
- Critical path testing results
- Business impact analysis

## 🎯 Test Scenarios

### Critical User Journeys
1. **New Customer Purchase**:
   - Browse vacuum categories
   - Search for specific features (pet hair, cordless)
   - Compare products
   - Add to cart and checkout
   - Create account during checkout

2. **Returning Customer**:
   - Login to account
   - View order history
   - Reorder previous items
   - Update profile information

3. **Admin Operations**:
   - Manage inventory levels
   - Process orders
   - View sales analytics
   - Update product information

### Performance Scenarios
- **Normal Load**: 25-50 concurrent users
- **Peak Traffic**: 100-200 concurrent users
- **Black Friday Spike**: 500+ concurrent users
- **Sustained Load**: 5-minute endurance testing

### Security Test Cases
- **Input Validation**: XSS, SQL injection, buffer overflow
- **Authentication**: Bypass attempts, token validation
- **Authorization**: Privilege escalation, access control
- **Business Logic**: Price manipulation, inventory exploits

## 🔧 Customization

### Adding New Test Cases
1. Implement the `TestRunner` interface
2. Add factory function in `testers.go`
3. Register new command in `main.go`

### Modifying Test Parameters
- Update user simulation counts in performance tests
- Add new device profiles for mobile testing
- Extend security payload lists
- Add business requirements to UAT scenarios

## 📈 Continuous Integration

The testing suite is designed for CI/CD integration:

```bash
# Run in CI environment
./vacuum-testing-suite all --base-url $STAGING_URL --log-level warn

# Exit codes
# 0: All tests passed
# 1: Some tests failed
# 2: Configuration error
```

## 🤝 Contributing

1. Follow Go coding standards
2. Add comprehensive logging for new test cases
3. Update documentation for new features
4. Ensure test isolation and repeatability

## 📝 Test Coverage

### API Endpoints Tested
- ✅ GET /api/products (with filtering and sorting)
- ✅ GET /api/products/{id}
- ✅ GET /api/products/search
- ✅ POST /api/cart/items
- ✅ PUT /api/cart/items/{id}
- ✅ DELETE /api/cart/items/{id}
- ✅ POST /api/orders
- ✅ GET /api/orders/{id}
- ✅ GET /api/admin/inventory
- ✅ PUT /api/admin/products/{id}/inventory

### Business Scenarios Tested
- ✅ Product browsing and filtering
- ✅ Shopping cart management
- ✅ Checkout process
- ✅ User account management
- ✅ Order tracking
- ✅ Admin dashboard functionality
- ✅ Mobile responsiveness
- ✅ Security vulnerabilities

### Performance Benchmarks
- **Target Response Time**: < 2 seconds for product pages
- **Concurrent Users**: Support for 100+ simultaneous users
- **Mobile Performance**: < 3 seconds load time on mobile
- **Availability**: 99.9% uptime during normal operations

## 🚨 Known Limitations

1. **Mock Payment Processing**: Payment integration testing uses mock endpoints
2. **Real Device Testing**: Mobile tests simulate device characteristics
3. **Database State**: Tests may require database reset between runs
4. **External Dependencies**: Some tests require specific test data

## 📞 Support

For issues with the testing suite:
1. Check test server connectivity
2. Verify VacuumMart application is running
3. Review log output for detailed error information
4. Ensure proper test data is available

## 🔄 Version History

- **v1.0.0**: Initial comprehensive testing suite
  - API testing implementation
  - E2E scenario coverage
  - Performance and security testing
  - Mobile responsiveness validation
  - UAT business requirement testing
