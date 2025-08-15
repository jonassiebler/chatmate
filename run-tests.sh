#!/bin/bash

# ChatMate Testing Framework Runner
# Comprehensive test execution script with reporting and coverage analysis

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TEST_DIR="$SCRIPT_DIR/tests"
BATS_LIBS_DIR="$TEST_DIR"

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
    
    # Ensure test directories exist
    mkdir -p "$TEST_DIR"/{unit,integration,fixtures,helpers}
    
    # Check for required tools
    local missing_tools=()
    
    if ! command -v bats >/dev/null 2>&1; then
        if [[ ! -x "$BATS_LIBS_DIR/bats-core/bin/bats" ]]; then
            missing_tools+=("bats")
        fi
    fi
    
    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        log_warning "Missing tools: ${missing_tools[*]}"
        log_info "Installing missing tools..."
        install_test_dependencies
    fi
    
    log_success "Test environment ready"
}

# Install test dependencies
install_test_dependencies() {
    log_info "Installing Bats testing framework..."
    
    # Install Bats Core
    if [[ ! -d "$BATS_LIBS_DIR/bats-core" ]]; then
        git clone https://github.com/bats-core/bats-core.git "$BATS_LIBS_DIR/bats-core"
    fi
    
    # Install Bats Support
    if [[ ! -d "$BATS_LIBS_DIR/bats-support" ]]; then
        git clone https://github.com/bats-core/bats-support.git "$BATS_LIBS_DIR/bats-support"
    fi
    
    # Install Bats Assert
    if [[ ! -d "$BATS_LIBS_DIR/bats-assert" ]]; then
        git clone https://github.com/bats-core/bats-assert.git "$BATS_LIBS_DIR/bats-assert"
    fi
    
    log_success "Bats framework installed"
}

# Get Bats executable path
get_bats_executable() {
    if command -v bats >/dev/null 2>&1; then
        echo "bats"
    elif [[ -x "$BATS_LIBS_DIR/bats-core/bin/bats" ]]; then
        echo "$BATS_LIBS_DIR/bats-core/bin/bats"
    else
        log_error "Bats executable not found"
        exit 1
    fi
}

# Run unit tests
run_unit_tests() {
    log_info "Running unit tests..."
    
    local bats_cmd
    bats_cmd=$(get_bats_executable)
    
    local test_files=("$TEST_DIR/unit"/*.bats)
    
    if [[ ! -f "${test_files[0]}" ]]; then
        log_warning "No unit test files found"
        return 0
    fi
    
    local unit_results
    unit_results=$("$bats_cmd" --tap "${test_files[@]}" 2>&1) || true
    
    echo "$unit_results" | tee -a "$LOG_FILE"
    
    # Parse results
    local unit_total unit_passed unit_failed
    unit_total=$(echo "$unit_results" | grep -c "^ok\|^not ok" || echo "0")
    unit_passed=$(echo "$unit_results" | grep -c "^ok" || echo "0")
    unit_failed=$(echo "$unit_results" | grep -c "^not ok" || echo "0")
    
    # Ensure we have valid numbers
    unit_total=${unit_total:-0}
    unit_passed=${unit_passed:-0}
    unit_failed=${unit_failed:-0}
    
    TOTAL_TESTS=$((TOTAL_TESTS + unit_total))
    PASSED_TESTS=$((PASSED_TESTS + unit_passed))
    FAILED_TESTS=$((FAILED_TESTS + unit_failed))
    
    if [[ $unit_failed -eq 0 ]]; then
        log_success "Unit tests completed: $unit_passed/$unit_total passed"
    else
        log_error "Unit tests failed: $unit_failed/$unit_total failed"
    fi
}

# Run integration tests
run_integration_tests() {
    log_info "Running integration tests..."
    
    local bats_cmd
    bats_cmd=$(get_bats_executable)
    
    local test_files=("$TEST_DIR/integration"/*.bats)
    
    if [[ ! -f "${test_files[0]}" ]]; then
        log_warning "No integration test files found"
        return 0
    fi
    
    local integration_results
    integration_results=$("$bats_cmd" --tap "${test_files[@]}" 2>&1) || true
    
    echo "$integration_results" | tee -a "$LOG_FILE"
    
    # Parse results
    local int_total int_passed int_failed
    int_total=$(echo "$integration_results" | grep -c "^ok\|^not ok" || echo "0")
    int_passed=$(echo "$integration_results" | grep -c "^ok" || echo "0")
    int_failed=$(echo "$integration_results" | grep -c "^not ok" || echo "0")
    
    # Ensure we have valid numbers
    int_total=${int_total:-0}
    int_passed=${int_passed:-0}
    int_failed=${int_failed:-0}
    
    TOTAL_TESTS=$((TOTAL_TESTS + int_total))
    PASSED_TESTS=$((PASSED_TESTS + int_passed))
    FAILED_TESTS=$((FAILED_TESTS + int_failed))
    
    if [[ $int_failed -eq 0 ]]; then
        log_success "Integration tests completed: $int_passed/$int_total passed"
    else
        log_error "Integration tests failed: $int_failed/$int_total failed"
    fi
}

# Run markdown linting
run_markdown_tests() {
    log_info "Running markdown quality tests..."
    
    if command -v markdownlint >/dev/null 2>&1; then
        local lint_config="$PROJECT_ROOT/.markdownlint.json"
        
        # Create markdownlint config if it doesn't exist
        if [[ ! -f "$lint_config" ]]; then
            cat > "$lint_config" << 'EOF'
{
  "MD013": false,
  "MD041": false,
  "MD033": false,
  "MD034": false,
  "MD032": false
}
EOF
        fi
        
        if markdownlint "**/*.md" --config "$lint_config" >> "$LOG_FILE" 2>&1; then
            log_success "Markdown linting passed"
        else
            log_warning "Markdown linting found issues (see log for details)"
        fi
    else
        log_warning "markdownlint not available, skipping markdown tests"
    fi
}

# Run shellcheck if available
run_shell_tests() {
    log_info "Running shell script quality tests..."
    
    if command -v shellcheck >/dev/null 2>&1; then
        local shell_files=("hire.sh")
        
        for file in "${shell_files[@]}"; do
            if [[ -f "$file" ]]; then
                if shellcheck "$file" >> "$LOG_FILE" 2>&1; then
                    log_success "Shell check passed for $file"
                else
                    log_warning "Shell check found issues in $file"
                fi
            fi
        done
    else
        log_warning "shellcheck not available, skipping shell script tests"
    fi
}

# Generate test report
generate_report() {
    log_info "Generating test report..."
    
    echo "" | tee -a "$LOG_FILE"
    echo "======================================" | tee -a "$LOG_FILE"
    echo "ChatMate Test Results Summary" | tee -a "$LOG_FILE"
    echo "======================================" | tee -a "$LOG_FILE"
    echo "Total Tests: $TOTAL_TESTS" | tee -a "$LOG_FILE"
    echo "Passed: $PASSED_TESTS" | tee -a "$LOG_FILE"
    echo "Failed: $FAILED_TESTS" | tee -a "$LOG_FILE"
    echo "Skipped: $SKIPPED_TESTS" | tee -a "$LOG_FILE"
    
    local success_rate
    if [[ $TOTAL_TESTS -gt 0 ]]; then
        success_rate=$((PASSED_TESTS * 100 / TOTAL_TESTS))
        echo "Success Rate: ${success_rate}%" | tee -a "$LOG_FILE"
    fi
    
    echo "Log File: $LOG_FILE" | tee -a "$LOG_FILE"
    echo "Timestamp: $(date)" | tee -a "$LOG_FILE"
    echo "======================================" | tee -a "$LOG_FILE"
    
    if [[ $FAILED_TESTS -eq 0 ]]; then
        log_success "All tests passed! 🎉"
        return 0
    else
        log_error "Some tests failed. Check the log for details."
        return 1
    fi
}

# Cleanup test environment
cleanup() {
    log_info "Cleaning up test environment..."
    
    # Remove temporary files if any
    find /tmp -name "chatmate-test-*" -type d -exec rm -rf {} + 2>/dev/null || true
    
    log_success "Cleanup completed"
}

# Main execution
main() {
    local run_unit=true
    local run_integration=true
    local run_quality=true
    local verbose=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --unit-only)
                run_integration=false
                run_quality=false
                shift
                ;;
            --integration-only)
                run_unit=false
                run_quality=false
                shift
                ;;
            --quality-only)
                run_unit=false
                run_integration=false
                shift
                ;;
            --verbose|-v)
                verbose=true
                shift
                ;;
            --help|-h)
                echo "Usage: $0 [OPTIONS]"
                echo "Options:"
                echo "  --unit-only        Run only unit tests"
                echo "  --integration-only Run only integration tests"
                echo "  --quality-only     Run only quality tests"
                echo "  --verbose, -v      Verbose output"
                echo "  --help, -h         Show this help"
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    # Initialize
    init_logging
    log_info "Starting ChatMate test suite..."
    
    # Setup
    setup_test_environment
    
    # Run tests based on options
    if [[ "$run_unit" == true ]]; then
        run_unit_tests
    fi
    
    if [[ "$run_integration" == true ]]; then
        run_integration_tests
    fi
    
    if [[ "$run_quality" == true ]]; then
        run_markdown_tests
        run_shell_tests
    fi
    
    # Generate report and cleanup
    if generate_report; then
        cleanup
        exit 0
    else
        cleanup
        exit 1
    fi
}

# Handle script interruption
trap cleanup EXIT

# Run main function with all arguments
main "$@"
