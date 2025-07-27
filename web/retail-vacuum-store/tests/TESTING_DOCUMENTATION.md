# E-commerce Testing Suite Documentation

## Overview

This comprehensive testing suite validates all aspects of the vacuum cleaner retail e-commerce platform, from individual API endpoints to complete user workflows, security vulnerabilities, and performance characteristics.

## Test Structure

```
tests/
├── api/                     # API endpoint testing
│   ├── products_test.go     # Product catalog API tests
│   └── cart_test.go         # Shopping cart API tests
├── e2e/                     # End-to-end workflow testing
│   └── shopping_flow_test.go # Complete shopping journeys
├── integration/             # Integration testing
│   └── payment_test.go      # Payment processing integration
├── security/                # Security vulnerability testing
│   └── security_test.go     # XSS, SQL injection, CSRF, etc.
├── load/                    # Performance and load testing
│   └── performance_test.go  # Load testing and benchmarks
├── fixtures/                # Test data and helpers
│   ├── products.go          # Product test fixtures
│   └── extended_fixtures.go # Extended test data sets
├── utils/                   # Test utilities and helpers
│   └── test_helpers.go      # Database, HTTP, and utility functions
├── run_tests.sh            # Comprehensive test runner script
├── TEST_PLAN.md            # Detailed testing strategy
└── TESTING_DOCUMENTATION.md # This file
```

## Running Tests

### Quick Start

```bash
# Run all tests with default settings
./tests/run_tests.sh

# Run specific test suite
./tests/run_tests.sh api

# Run with verbose output and coverage
./tests/run_tests.sh --verbose --coverage

# Run performance tests (requires special setup)
./tests/run_tests.sh --performance
```

### Available Test Suites

1. **API Tests** (`api`)
   - Product catalog endpoints
   - Shopping cart operations
   - Search functionality
   - Error handling

2. **End-to-End Tests** (`e2e`)
   - Complete shopping workflows
   - User account management
   - Order processing
   - Error recovery scenarios

3. **Integration Tests** (`integration`)
   - Payment processing
   - External service integration
   - Database transactions
   - Webhook processing

4. **Security Tests** (`security`)
   - SQL injection prevention
   - XSS protection
   - CSRF validation
   - Input sanitization
   - Authentication security

5. **Performance Tests** (`load`)
   - Load testing scenarios
   - Stress testing
   - Endurance testing
   - Memory leak detection

### Test Configuration Options

```bash
# Environment variables
export VERBOSE=true          # Enable verbose output
export COVERAGE=true         # Enable coverage reporting
export RACE_DETECTION=true   # Enable race condition detection
export PARALLEL=true         # Enable parallel test execution
export TIMEOUT=300s          # Set test timeout
export TEST_PATTERN=".*"     # Filter tests by pattern

# Command line options
--verbose                    # Enable verbose output
--coverage / --no-coverage   # Control coverage reporting
--race / --no-race          # Control race detection
--performance               # Include performance tests
--no-parallel              # Disable parallel execution
--timeout DURATION          # Set custom timeout
--pattern PATTERN           # Run only matching tests
```

## Test Categories

### Unit Tests
**Location**: Throughout the codebase alongside source files  
**Purpose**: Test individual functions and methods in isolation  
**Coverage Goal**: >90% line coverage  
**Execution Time**: <30 seconds  

### API Tests
**Location**: `tests/api/`  
**Purpose**: Validate HTTP API endpoints  
**Coverage Goal**: 100% endpoint coverage  
**Key Scenarios**:
- Product CRUD operations
- Shopping cart management
- Search and filtering
- Pagination and sorting
- Error responses

### Integration Tests
**Location**: `tests/integration/`  
**Purpose**: Test component interactions  
**Coverage Goal**: All external integrations  
**Key Scenarios**:
- Payment gateway integration
- Database transactions
- Email notifications
- Third-party APIs

### End-to-End Tests
**Location**: `tests/e2e/`  
**Purpose**: Validate complete user workflows  
**Coverage Goal**: All critical user paths  
**Key Scenarios**:
- Guest checkout process
- User registration and login
- Order management
- Cart abandonment and recovery

### Security Tests
**Location**: `tests/security/`  
**Purpose**: Identify security vulnerabilities  
**Coverage Goal**: OWASP Top 10 coverage  
**Key Scenarios**:
- SQL injection attempts
- XSS payload testing
- CSRF protection validation
- Authentication bypass attempts
- Input validation testing

### Performance Tests
**Location**: `tests/load/`  
**Purpose**: Validate system performance  
**Coverage Goal**: All performance requirements  
**Key Scenarios**:
- Concurrent user simulation
- API response time validation
- Database performance testing
- Memory usage monitoring

## Test Data Management

### Fixtures
The `fixtures/` directory contains pre-defined test data:

- **Products**: Various vacuum cleaner models with different attributes
- **Categories**: Product categorization hierarchy
- **Brands**: Manufacturer information
- **Customers**: Test user accounts and profiles
- **Orders**: Sample transaction data
- **Coupons**: Discount codes and promotions

### Database Setup
Tests use isolated SQLite in-memory databases:

```go
// Example database setup
func DatabaseTestSetup(t *testing.T) *TestDatabase {
    testDB := &TestDatabase{
        DSN:      ":memory:",
        TestName: t.Name(),
    }
    // ... setup code
}
```

### Mock Services
External services are mocked for reliable testing:

- Payment gateways (Stripe, PayPal)
- Email services
- Inventory management systems
- Analytics services

## Performance Benchmarks

### API Response Times
- **Product Listing**: <200ms (95th percentile)
- **Product Search**: <100ms (95th percentile)
- **Cart Operations**: <50ms (95th percentile)
- **Checkout Process**: <500ms (95th percentile)

### Load Testing Targets
- **Low Load**: 10 concurrent users, >40 req/sec
- **Medium Load**: 50 concurrent users, >150 req/sec
- **High Load**: 100 concurrent users, >250 req/sec
- **Peak Load**: 200 concurrent users, >300 req/sec

### Error Rate Thresholds
- **Normal Load**: <1% error rate
- **High Load**: <5% error rate
- **Stress Test**: <10% error rate

## Security Testing

### OWASP Top 10 Coverage
1. **Injection**: SQL injection testing across all inputs
2. **Broken Authentication**: Session management validation
3. **Sensitive Data Exposure**: Data protection verification
4. **XML External Entities**: Input parsing security
5. **Broken Access Control**: Authorization testing
6. **Security Misconfiguration**: Configuration validation
7. **Cross-Site Scripting**: XSS prevention testing
8. **Insecure Deserialization**: Input validation
9. **Known Vulnerabilities**: Dependency scanning
10. **Insufficient Logging**: Audit trail validation

### Input Validation Testing
All user inputs are tested with:
- Malicious payloads
- SQL injection attempts
- XSS vectors
- Buffer overflow attempts
- Format string attacks
- Path traversal attempts

## Continuous Integration

### Pre-commit Hooks
```bash
# Install pre-commit hooks
make install-hooks

# Hooks include:
# - Go formatting check
# - Unit test execution
# - Basic linting
# - Import organization
```

### CI/CD Pipeline
1. **On Pull Request**:
   - Full test suite execution
   - Coverage report generation
   - Security scan
   - Performance regression check

2. **On Merge**:
   - Production deployment tests
   - End-to-end validation
   - Performance benchmarking
   - Security vulnerability scan

### Test Reporting
Tests generate comprehensive reports:

- **Coverage Reports**: HTML coverage visualization
- **Performance Reports**: Response time metrics
- **Security Reports**: Vulnerability assessments
- **Test Execution Reports**: Pass/fail summaries

## Troubleshooting

### Common Issues

#### Test Database Issues
```bash
# Clear test database
rm -f tests/test_*.db

# Reset database schema
go test ./tests/utils -run TestDatabaseSetup
```

#### Coverage Issues
```bash
# Generate fresh coverage report
rm -rf coverage/
./tests/run_tests.sh --coverage

# View coverage report
open test_reports/coverage.html
```

#### Performance Test Issues
```bash
# Run performance tests in isolation
./tests/run_tests.sh performance --no-parallel

# Check system resources
top
free -m
```

#### Mock Service Issues
```bash
# Reset mock services
pkill -f "mock-server"
./tests/start_mocks.sh
```

### Debug Mode
```bash
# Run tests with maximum verbosity
VERBOSE=true DEBUG=true ./tests/run_tests.sh --verbose

# Run specific test with debugger
go test -v ./tests/api -run TestProductAPI_Specific
```

### Log Analysis
```bash
# View test logs
tail -f test_reports/test_execution.log

# Search for specific errors
grep -r "ERROR" test_reports/

# Analyze performance metrics
cat test_reports/performance_metrics.json | jq '.latency_p95'
```

## Best Practices

### Writing Tests
1. **Isolation**: Each test should be independent
2. **Clarity**: Test names should describe the scenario
3. **Coverage**: Test both success and failure paths
4. **Data**: Use fixtures for consistent test data
5. **Cleanup**: Always clean up test artifacts

### Test Organization
1. **Grouping**: Related tests in the same file
2. **Naming**: Consistent naming conventions
3. **Documentation**: Clear test documentation
4. **Maintenance**: Regular test review and updates

### Performance Considerations
1. **Parallel Execution**: Enable parallel tests where safe
2. **Resource Management**: Clean up resources promptly
3. **Test Data**: Use minimal required data sets
4. **Timeouts**: Set appropriate test timeouts

## Extending the Test Suite

### Adding New Tests
1. Choose appropriate test category
2. Use existing fixtures and utilities
3. Follow naming conventions
4. Include both positive and negative scenarios
5. Update test documentation

### Adding Test Data
1. Add fixtures to `fixtures/` directory
2. Use factory patterns for dynamic data
3. Include cleanup procedures
4. Document data relationships

### Custom Test Utilities
1. Add utilities to `utils/` directory
2. Include comprehensive documentation
3. Write unit tests for utilities
4. Follow project coding standards

## Metrics and Reporting

### Key Metrics Tracked
- **Test Coverage**: Line and branch coverage
- **Test Execution Time**: Individual and suite times
- **Error Rates**: Pass/fail ratios by category
- **Performance Metrics**: Response times and throughput
- **Security Findings**: Vulnerability counts and types

### Report Generation
```bash
# Generate all reports
./tests/run_tests.sh --generate-reports

# View reports
open test_reports/index.html
```

### Metrics Dashboard
Test metrics are available via:
- HTML reports in `test_reports/`
- JSON metrics in `test_reports/metrics.json`
- CI/CD integration endpoints
- Real-time monitoring dashboards

## Conclusion

This comprehensive testing suite provides confidence in the e-commerce platform's reliability, security, and performance. Regular execution of all test categories ensures high quality software delivery and rapid identification of regressions.

For questions or issues with the testing suite, refer to the troubleshooting section or contact the development team.
