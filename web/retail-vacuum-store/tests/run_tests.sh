#!/bin/bash

# E-commerce Testing Suite Runner
# Comprehensive test execution script for the vacuum cleaner retail store

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
VERBOSE=${VERBOSE:-false}
COVERAGE=${COVERAGE:-true}
RACE_DETECTION=${RACE_DETECTION:-true}
PARALLEL=${PARALLEL:-true}
TIMEOUT=${TIMEOUT:-300s}
TEST_PATTERN=${TEST_PATTERN:-".*"}

# Directories
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEST_DIR="$PROJECT_ROOT/tests"
COVERAGE_DIR="$PROJECT_ROOT/coverage"
REPORTS_DIR="$PROJECT_ROOT/test_reports"

echo -e "${BLUE}=== E-commerce Testing Suite ===${NC}"
echo "Project Root: $PROJECT_ROOT"
echo "Test Directory: $TEST_DIR"
echo ""

# Create output directories
mkdir -p "$COVERAGE_DIR" "$REPORTS_DIR"

# Function to print section headers
print_section() {
    echo -e "\n${BLUE}=== $1 ===${NC}"
}

# Function to print success message
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# Function to print error message
print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Function to print warning message
print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Function to run a test suite
run_test_suite() {
    local suite_name="$1"
    local test_path="$2"
    local extra_flags="$3"
    
    echo -e "\n${YELLOW}Running $suite_name tests...${NC}"
    
    local cmd="go test"
    local flags=""
    
    # Add common flags
    if [ "$VERBOSE" = "true" ]; then
        flags="$flags -v"
    fi
    
    if [ "$RACE_DETECTION" = "true" ]; then
        flags="$flags -race"
    fi
    
    if [ "$PARALLEL" = "true" ]; then
        flags="$flags -parallel 4"
    fi
    
    flags="$flags -timeout $TIMEOUT"
    flags="$flags -run $TEST_PATTERN"
    
    # Add coverage if enabled
    if [ "$COVERAGE" = "true" ]; then
        local coverage_file="$COVERAGE_DIR/${suite_name}_coverage.out"
        flags="$flags -coverprofile=$coverage_file"
    fi
    
    # Add any extra flags
    flags="$flags $extra_flags"
    
    # Run the test
    if $cmd $flags "$test_path"; then
        print_success "$suite_name tests passed"
        return 0
    else
        print_error "$suite_name tests failed"
        return 1
    fi
}

# Function to check prerequisites
check_prerequisites() {
    print_section "Checking Prerequisites"
    
    # Check Go installation
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    print_success "Go $(go version | cut -d' ' -f3) found"
    
    # Check if in correct directory
    if [ ! -f "$PROJECT_ROOT/go.mod" ]; then
        print_error "Not in a Go module directory"
        exit 1
    fi
    print_success "Go module found"
    
    # Check test dependencies
    echo "Checking test dependencies..."
    if ! go list -m github.com/stretchr/testify &> /dev/null; then
        print_warning "Installing testify..."
        go get github.com/stretchr/testify
    fi
    print_success "Test dependencies available"
}

# Function to lint code
run_linting() {
    print_section "Code Linting"
    
    # Format check
    echo "Checking code formatting..."
    if ! gofmt -l . | grep -v vendor/ | grep -q .; then
        print_success "Code is properly formatted"
    else
        print_warning "Code formatting issues found"
        if [ "$VERBOSE" = "true" ]; then
            gofmt -l . | grep -v vendor/
        fi
    fi
    
    # Vet check
    echo "Running go vet..."
    if go vet ./...; then
        print_success "Go vet passed"
    else
        print_error "Go vet found issues"
    fi
    
    # Golangci-lint if available
    if command -v golangci-lint &> /dev/null; then
        echo "Running golangci-lint..."
        if golangci-lint run; then
            print_success "Linting passed"
        else
            print_warning "Linting found issues"
        fi
    else
        print_warning "golangci-lint not available"
    fi
}

# Function to run unit tests
run_unit_tests() {
    print_section "Unit Tests"
    
    # Find all packages with unit tests
    local unit_packages
    unit_packages=$(find "$PROJECT_ROOT" -name "*_test.go" -not -path "*/tests/*" -exec dirname {} \; | sort -u)
    
    if [ -n "$unit_packages" ]; then
        for package in $unit_packages; do
            local rel_path
            rel_path=$(realpath --relative-to="$PROJECT_ROOT" "$package")
            run_test_suite "unit-$(basename "$package")" "./$rel_path" ""
        done
    else
        print_warning "No unit tests found"
    fi
}

# Function to run API tests
run_api_tests() {
    print_section "API Tests"
    
    if [ -d "$TEST_DIR/api" ]; then
        run_test_suite "api" "$TEST_DIR/api" ""
    else
        print_warning "API tests directory not found"
    fi
}

# Function to run E2E tests
run_e2e_tests() {
    print_section "End-to-End Tests"
    
    if [ -d "$TEST_DIR/e2e" ]; then
        # E2E tests might need more time
        local old_timeout="$TIMEOUT"
        TIMEOUT="600s"
        run_test_suite "e2e" "$TEST_DIR/e2e" ""
        TIMEOUT="$old_timeout"
    else
        print_warning "E2E tests directory not found"
    fi
}

# Function to run integration tests
run_integration_tests() {
    print_section "Integration Tests"
    
    if [ -d "$TEST_DIR/integration" ]; then
        run_test_suite "integration" "$TEST_DIR/integration" ""
    else
        print_warning "Integration tests directory not found"
    fi
}

# Function to run security tests
run_security_tests() {
    print_section "Security Tests"
    
    if [ -d "$TEST_DIR/security" ]; then
        run_test_suite "security" "$TEST_DIR/security" ""
    else
        print_warning "Security tests directory not found"
    fi
}

# Function to run performance tests
run_performance_tests() {
    print_section "Performance Tests"
    
    if [ -d "$TEST_DIR/load" ]; then
        # Performance tests need special handling
        print_warning "Performance tests require special setup - running basic tests only"
        run_test_suite "performance" "$TEST_DIR/load" "-short"
    else
        print_warning "Performance tests directory not found"
    fi
}

# Function to run benchmarks
run_benchmarks() {
    print_section "Benchmarks"
    
    echo "Running benchmarks..."
    local bench_output="$REPORTS_DIR/benchmarks.txt"
    
    if go test -bench=. -benchmem ./... > "$bench_output" 2>&1; then
        print_success "Benchmarks completed"
        if [ "$VERBOSE" = "true" ]; then
            cat "$bench_output"
        fi
    else
        print_warning "Some benchmarks failed"
    fi
}

# Function to generate coverage report
generate_coverage_report() {
    if [ "$COVERAGE" != "true" ]; then
        return 0
    fi
    
    print_section "Coverage Report"
    
    # Merge coverage files
    echo "Merging coverage files..."
    local merged_coverage="$COVERAGE_DIR/merged_coverage.out"
    
    # Create merged coverage file
    echo "mode: atomic" > "$merged_coverage"
    
    # Merge all coverage files
    for coverage_file in "$COVERAGE_DIR"/*_coverage.out; do
        if [ -f "$coverage_file" ]; then
            tail -n +2 "$coverage_file" >> "$merged_coverage"
        fi
    done
    
    # Generate HTML report
    echo "Generating HTML coverage report..."
    local html_report="$REPORTS_DIR/coverage.html"
    
    if go tool cover -html="$merged_coverage" -o "$html_report"; then
        print_success "Coverage report generated: $html_report"
    else
        print_error "Failed to generate coverage report"
    fi
    
    # Show coverage summary
    echo "Coverage summary:"
    go tool cover -func="$merged_coverage" | tail -1
}

# Function to run all tests
run_all_tests() {
    local failed_suites=()
    
    # Track start time
    local start_time
    start_time=$(date +%s)
    
    # Run test suites
    run_unit_tests || failed_suites+=("unit")
    run_api_tests || failed_suites+=("api")
    run_integration_tests || failed_suites+=("integration")
    run_security_tests || failed_suites+=("security")
    run_e2e_tests || failed_suites+=("e2e")
    
    # Run performance tests separately (optional)
    if [ "${RUN_PERFORMANCE:-false}" = "true" ]; then
        run_performance_tests || failed_suites+=("performance")
    fi
    
    # Generate reports
    generate_coverage_report
    
    # Calculate total time
    local end_time
    end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    # Summary
    print_section "Test Summary"
    echo "Total test duration: ${duration}s"
    
    if [ ${#failed_suites[@]} -eq 0 ]; then
        print_success "All test suites passed!"
        echo -e "\n${GREEN}🎉 Testing completed successfully!${NC}"
        return 0
    else
        print_error "Failed test suites: ${failed_suites[*]}"
        echo -e "\n${RED}❌ Testing completed with failures!${NC}"
        return 1
    fi
}

# Function to clean up test artifacts
cleanup() {
    print_section "Cleanup"
    
    echo "Cleaning up test artifacts..."
    
    # Remove temporary files
    find . -name "*.test" -delete
    find . -name "test_*.db" -delete
    
    print_success "Cleanup completed"
}

# Function to show help
show_help() {
    cat << EOF
E-commerce Testing Suite Runner

Usage: $0 [OPTIONS] [TEST_SUITE]

Options:
    -v, --verbose       Enable verbose output
    -c, --coverage      Enable coverage reporting (default: true)
    -r, --race          Enable race detection (default: true)
    -p, --performance   Include performance tests
    --no-parallel       Disable parallel test execution
    --timeout DURATION  Set test timeout (default: 300s)
    --pattern PATTERN   Run only tests matching pattern
    -h, --help          Show this help message

Test Suites:
    all                 Run all test suites (default)
    unit                Run unit tests only
    api                 Run API tests only
    integration         Run integration tests only
    security            Run security tests only
    e2e                 Run end-to-end tests only
    performance         Run performance tests only
    benchmarks          Run benchmarks only
    lint                Run linting only

Examples:
    $0                          # Run all tests
    $0 --verbose api            # Run API tests with verbose output
    $0 --performance --timeout 600s  # Include performance tests with longer timeout
    $0 --pattern "TestCart.*"   # Run only cart-related tests

Environment Variables:
    VERBOSE=true                # Enable verbose output
    COVERAGE=false              # Disable coverage
    RACE_DETECTION=false        # Disable race detection
    RUN_PERFORMANCE=true        # Include performance tests
    TEST_PATTERN=".*"           # Test pattern filter

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        --no-coverage)
            COVERAGE=false
            shift
            ;;
        -r|--race)
            RACE_DETECTION=true
            shift
            ;;
        --no-race)
            RACE_DETECTION=false
            shift
            ;;
        -p|--performance)
            RUN_PERFORMANCE=true
            shift
            ;;
        --no-parallel)
            PARALLEL=false
            shift
            ;;
        --timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        --pattern)
            TEST_PATTERN="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            TEST_SUITE="$1"
            shift
            ;;
    esac
done

# Set default test suite
TEST_SUITE=${TEST_SUITE:-"all"}

# Main execution
main() {
    check_prerequisites
    
    case "$TEST_SUITE" in
        all)
            run_linting
            run_all_tests
            ;;
        unit)
            run_unit_tests
            ;;
        api)
            run_api_tests
            ;;
        integration)
            run_integration_tests
            ;;
        security)
            run_security_tests
            ;;
        e2e)
            run_e2e_tests
            ;;
        performance)
            run_performance_tests
            ;;
        benchmarks)
            run_benchmarks
            ;;
        lint)
            run_linting
            ;;
        *)
            print_error "Unknown test suite: $TEST_SUITE"
            show_help
            exit 1
            ;;
    esac
    
    local exit_code=$?
    
    # Always run cleanup
    cleanup
    
    exit $exit_code
}

# Handle script interruption
trap cleanup EXIT

# Run main function
main "$@"
