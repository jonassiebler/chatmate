# Test Directory

This directory contains external test applications and test data following Go community standards.

## Structure

- **`fixtures/`** - Test data and sample files used by tests across the project
- **`integration/`** - Integration tests that test complete workflows and system behavior  
- **`benchmarks/`** - Performance benchmarks for the application

## Root Level Test Files

- **`chatmate_validation_test.go`** - Validation tests for chatmate files and project structure
- **`hire_script_test.go`** - Tests for the hire.sh installation script

## Running Tests

```bash
# Run all tests
go test ./...

# Run only tests in test/ directory
go test ./test/...

# Run integration tests
go test ./test/integration/...

# Run benchmarks
go test -bench=. ./test/benchmarks/...

# Using the test runner script
./run-tests.sh              # Run all tests
./run-tests.sh --integration # Run integration tests only
./run-tests.sh --benchmark   # Run benchmark tests only
```

## Test Organization

Tests follow Go's standard testing conventions:
- Test files end with `_test.go`
- Test functions start with `Test`
- Benchmark functions start with `Benchmark`
- Example functions start with `Example`

For unit tests that are tightly coupled with specific packages, those remain alongside the source code in their respective packages (e.g., `cmd/`, `pkg/`, `internal/`).
