#!/bin/bash

# ChatMate Testing Framework Runner
# Comprehensive test execution script with Go testing and coverage analysis

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results tracking
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# Logging
LOG_FILE="$PROJECT_ROOT/test-results.log"
COVERAGE_FILE="$PROJECT_ROOT/coverage.out"
COVERAGE_HTML="$PROJECT_ROOT/coverage.html"

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1" | tee -a "$LOG_FILE"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "$LOG_FILE"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$LOG_FILE"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"
}

# Initialize logging
init_logging() {
    echo "ChatMate Test Run - $(date)" > "$LOG_FILE"
    echo "======================================" >> "$LOG_FILE"
}

# Setup test environment
setup_test_environment() {
    log_info "Setting up test environment..."
    
    cd "$PROJECT_ROOT" || exit 1
    
    # Check for required tools
    if ! command -v go >/dev/null 2>&1; then
        log_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    # Verify Go version
    local go_version
    go_version=$(go version | awk '{print $3}' | sed 's/go//')
    log_info "Using Go version: $go_version"
    
    # Install/update dependencies
    log_info "Installing Go dependencies..."
    go mod download
    go mod tidy
    
    log_success "Test environment ready"
}

# Run unit tests
run_unit_tests() {
    log_info "Running Go unit tests..."
    
    # Run tests with coverage
    local test_output
    test_output=$(go test -v -coverprofile="$COVERAGE_FILE" ./... 2>&1) || true
    
    echo "$test_output" | tee -a "$LOG_FILE"
    
    # Parse results from test output
    local unit_total unit_passed unit_failed unit_skipped
    
    # Count top-level test functions only (=== RUN TestFunctionName, not sub-tests)
    unit_total=$(echo "$test_output" | grep -E "^=== RUN   Test[A-Z]" | grep -v "^=== RUN.*/" | wc -l | tr -d ' ')
    unit_passed=$(echo "$test_output" | grep -E "^--- PASS: Test[A-Z]" | grep -v "^--- PASS:.*/" | wc -l | tr -d ' ')
    unit_failed=$(echo "$test_output" | grep -E "^--- FAIL: Test[A-Z]" | grep -v "^--- FAIL:.*/" | wc -l | tr -d ' ')
    unit_skipped=$(echo "$test_output" | grep -E "^--- SKIP: Test[A-Z]" | grep -v "^--- SKIP:.*/" | wc -l | tr -d ' ')
    
    # Ensure we have valid numbers (fallback to 0 if empty)
    unit_total=${unit_total:-0}
    unit_passed=${unit_passed:-0}
    unit_failed=${unit_failed:-0}
    unit_skipped=${unit_skipped:-0}
    
    # Convert empty strings to 0
    [[ -z "$unit_total" ]] && unit_total=0
    [[ -z "$unit_passed" ]] && unit_passed=0
    [[ -z "$unit_failed" ]] && unit_failed=0
    [[ -z "$unit_skipped" ]] && unit_skipped=0
    
    TOTAL_TESTS=$((TOTAL_TESTS + unit_total))
    PASSED_TESTS=$((PASSED_TESTS + unit_passed))
    FAILED_TESTS=$((FAILED_TESTS + unit_failed))
    SKIPPED_TESTS=$((SKIPPED_TESTS + unit_skipped))
    
    if [[ $unit_failed -eq 0 ]]; then
        log_success "Unit tests completed: $unit_passed/$unit_total passed"
        if [[ $unit_skipped -gt 0 ]]; then
            log_warning "$unit_skipped tests were skipped"
        fi
    else
        log_error "Unit tests failed: $unit_failed/$unit_total failed"
    fi
    
    return $unit_failed
}

# Run integration tests
run_integration_tests() {
    log_info "Running Go integration tests..."
    
    # Run integration tests with longer timeout
    local integration_output
    integration_output=$(go test -v -timeout=30m -tags=integration ./... 2>&1) || true
    
    echo "$integration_output" | tee -a "$LOG_FILE"
    
    # Parse results from integration test output
    local integration_total integration_passed integration_failed integration_skipped
    
    # Count integration tests (look for test functions with Integration in name)
    integration_total=$(echo "$integration_output" | grep -E "^=== RUN   .*Integration|^=== RUN   .*Test.*Suite" | wc -l | tr -d ' ')
    integration_passed=$(echo "$integration_output" | grep -E "^--- PASS: .*Integration|^--- PASS: .*Test.*Suite" | wc -l | tr -d ' ')
    integration_failed=$(echo "$integration_output" | grep -E "^--- FAIL: .*Integration|^--- FAIL: .*Test.*Suite" | wc -l | tr -d ' ')
    integration_skipped=$(echo "$integration_output" | grep -E "^--- SKIP: .*Integration|^--- SKIP: .*Test.*Suite" | wc -l | tr -d ' ')
    
    # Ensure we have valid numbers (fallback to 0 if empty)
    integration_total=${integration_total:-0}
    integration_passed=${integration_passed:-0}
    integration_failed=${integration_failed:-0}
    integration_skipped=${integration_skipped:-0}
    
    # Convert empty strings to 0
    [[ -z "$integration_total" ]] && integration_total=0
    [[ -z "$integration_passed" ]] && integration_passed=0
    [[ -z "$integration_failed" ]] && integration_failed=0
    [[ -z "$integration_skipped" ]] && integration_skipped=0
    
    TOTAL_TESTS=$((TOTAL_TESTS + integration_total))
    PASSED_TESTS=$((PASSED_TESTS + integration_passed))
    FAILED_TESTS=$((FAILED_TESTS + integration_failed))
    SKIPPED_TESTS=$((SKIPPED_TESTS + integration_skipped))
    
    if [[ $integration_failed -eq 0 ]]; then
        log_success "Integration tests completed: $integration_passed/$integration_total passed"
        if [[ $integration_skipped -gt 0 ]]; then
            log_warning "$integration_skipped integration tests were skipped"
        fi
    else
        log_error "Integration tests failed: $integration_failed/$integration_total failed"
    fi
    
    return $integration_failed
}

# Run benchmark tests
run_benchmark_tests() {
    log_info "Running Go benchmark tests..."
    
    local benchmark_output
    benchmark_output=$(go test -bench=. -benchmem ./... 2>&1) || true
    
    echo "$benchmark_output" | tee -a "$LOG_FILE"
    
    local benchmark_count
    benchmark_count=$(echo "$benchmark_output" | grep -c "^Benchmark" || echo "0")
    
    if [[ $benchmark_count -gt 0 ]]; then
        log_success "Benchmark tests completed: $benchmark_count benchmarks run"
    else
        log_info "No benchmark tests found"
    fi
}

# Generate coverage report
generate_coverage_report() {
    if [[ -f "$COVERAGE_FILE" ]]; then
        log_info "Generating coverage report..."
        
        # Generate HTML coverage report
        go tool cover -html="$COVERAGE_FILE" -o "$COVERAGE_HTML"
        
        # Get coverage percentage
        local coverage_percent
        coverage_percent=$(go tool cover -func="$COVERAGE_FILE" | tail -1 | awk '{print $3}')
        
        log_info "Code coverage: $coverage_percent"
        log_success "Coverage report saved to: $COVERAGE_HTML"
        
        # Parse coverage percentage for threshold check
        local coverage_num
        coverage_num=$(echo "$coverage_percent" | sed 's/%//')
        
        if (( $(echo "$coverage_num >= 80" | bc -l) )); then
            log_success "Coverage meets threshold (>= 80%)"
        elif (( $(echo "$coverage_num >= 60" | bc -l) )); then
            log_warning "Coverage is moderate ($coverage_percent), consider improving"
        else
            log_warning "Coverage is low ($coverage_percent), improvement recommended"
        fi
    else
        log_warning "No coverage data available"
    fi
}

# Run security checks
run_security_checks() {
    log_info "Running security checks..."
    
    # Check for common security issues
    local security_issues=0
    
    # Check for hardcoded credentials/secrets
    if grep -r -i "password\|secret\|key\|token" --include="*.go" . >/dev/null 2>&1; then
        log_warning "Potential hardcoded credentials found - review manually"
        ((security_issues++))
    fi
    
    # Check for unsafe operations
    if grep -r "unsafe\." --include="*.go" . >/dev/null 2>&1; then
        log_warning "Unsafe operations found - review manually"
        ((security_issues++))
    fi
    
    # Check for SQL injection patterns (basic check)
    if grep -r "fmt\.Sprintf.*%s.*sql\|fmt\.Sprintf.*%v.*sql" --include="*.go" . >/dev/null 2>&1; then
        log_warning "Potential SQL injection patterns found - review manually"
        ((security_issues++))
    fi
    
    if [[ $security_issues -eq 0 ]]; then
        log_success "Security checks passed"
    else
        log_warning "Found $security_issues potential security issues"
    fi
}

# Run code quality checks
run_quality_checks() {
    log_info "Running code quality checks..."
    
    # Go vet
    log_info "Running go vet..."
    if go vet ./...; then
        log_success "go vet passed"
    else
        log_error "go vet found issues"
        return 1
    fi
    
    # Go fmt check
    log_info "Checking go fmt..."
    local fmt_issues
    fmt_issues=$(gofmt -l . | grep -v vendor | grep -v .git || true)
    
    if [[ -z "$fmt_issues" ]]; then
        log_success "go fmt check passed"
    else
        log_warning "go fmt issues found in: $fmt_issues"
    fi
    
    # Check for gosec if available
    if command -v gosec >/dev/null 2>&1; then
        log_info "Running gosec security analyzer..."
        if gosec ./...; then
            log_success "gosec security check passed"
        else
            log_warning "gosec found potential security issues"
        fi
    else
        log_info "gosec not available - skipping advanced security analysis"
    fi
}

# Print test summary
print_summary() {
    echo ""
    echo "======================================" | tee -a "$LOG_FILE"
    echo "TEST SUMMARY" | tee -a "$LOG_FILE"
    echo "======================================" | tee -a "$LOG_FILE"
    
    log_info "Total tests run: $TOTAL_TESTS"
    log_success "Tests passed: $PASSED_TESTS"
    
    if [[ $FAILED_TESTS -gt 0 ]]; then
        log_error "Tests failed: $FAILED_TESTS"
    fi
    
    if [[ $SKIPPED_TESTS -gt 0 ]]; then
        log_warning "Tests skipped: $SKIPPED_TESTS"
    fi
    
    local success_rate
    if [[ $TOTAL_TESTS -gt 0 ]]; then
        # Use bc for accurate arithmetic to avoid bash integer overflow
        success_rate=$(echo "scale=1; $PASSED_TESTS * 100 / $TOTAL_TESTS" | bc 2>/dev/null || echo "0")
        # Round to nearest integer
        success_rate=$(printf "%.0f" "$success_rate" 2>/dev/null || echo "0")
        log_info "Success rate: ${success_rate}%"
    else
        log_info "Success rate: N/A (no tests run)"
    fi
    
    echo "======================================" | tee -a "$LOG_FILE"
    
    if [[ $FAILED_TESTS -eq 0 ]]; then
        log_success "ALL TESTS PASSED!"
        return 0
    else
        log_error "SOME TESTS FAILED!"
        return 1
    fi
}

# Show usage information
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -u, --unit           Run only unit tests"
    echo "  -i, --integration    Run only integration tests"
    echo "  -b, --benchmark      Run benchmark tests"
    echo "  -c, --coverage       Generate coverage report only"
    echo "  -q, --quality        Run quality checks only"
    echo "  -s, --security       Run security checks only"
    echo "  -a, --all            Run all tests (default)"
    echo "  -h, --help           Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                   # Run all tests"
    echo "  $0 --unit           # Run only unit tests"
    echo "  $0 --integration    # Run only integration tests"
    echo "  $0 --coverage       # Generate coverage report"
}

# Clean up function
cleanup() {
    log_info "Cleaning up temporary files..."
    # Add cleanup logic here if needed
}

# Main execution
main() {
    local run_unit=true
    local run_integration=true
    local run_benchmark=false
    local run_coverage=true
    local run_quality=false
    local run_security=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -u|--unit)
                run_unit=true
                run_integration=false
                run_benchmark=false
                run_coverage=false
                run_quality=false
                run_security=false
                shift
                ;;
            -i|--integration)
                run_unit=false
                run_integration=true
                run_benchmark=false
                run_coverage=false
                run_quality=false
                run_security=false
                shift
                ;;
            -b|--benchmark)
                run_unit=false
                run_integration=false
                run_benchmark=true
                run_coverage=false
                run_quality=false
                run_security=false
                shift
                ;;
            -c|--coverage)
                run_unit=false
                run_integration=false
                run_benchmark=false
                run_coverage=true
                run_quality=false
                run_security=false
                shift
                ;;
            -q|--quality)
                run_unit=false
                run_integration=false
                run_benchmark=false
                run_coverage=false
                run_quality=true
                run_security=false
                shift
                ;;
            -s|--security)
                run_unit=false
                run_integration=false
                run_benchmark=false
                run_coverage=false
                run_quality=false
                run_security=true
                shift
                ;;
            -a|--all)
                run_unit=true
                run_integration=true
                run_benchmark=true
                run_coverage=true
                run_quality=true
                run_security=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    # Set up signal handlers for cleanup
    trap cleanup EXIT
    
    # Initialize
    init_logging
    setup_test_environment
    
    local exit_code=0
    
    # Run tests based on options
    if [[ "$run_quality" == "true" ]]; then
        if ! run_quality_checks; then
            exit_code=1
        fi
    fi
    
    if [[ "$run_security" == "true" ]]; then
        run_security_checks
    fi
    
    if [[ "$run_unit" == "true" ]]; then
        if ! run_unit_tests; then
            exit_code=1
        fi
    fi
    
    if [[ "$run_integration" == "true" ]]; then
        if ! run_integration_tests; then
            exit_code=1
        fi
    fi
    
    if [[ "$run_benchmark" == "true" ]]; then
        run_benchmark_tests
    fi
    
    if [[ "$run_coverage" == "true" ]]; then
        generate_coverage_report
    fi
    
    # Print summary
    if ! print_summary; then
        exit_code=1
    fi
    
    exit $exit_code
}

# Run main function
main "$@"
