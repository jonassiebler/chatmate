# Test Directory

This directory contains external test applications and test data following Go community standards for the ChatMate Go CLI.

## Structure

- **`fixtures/`** - Test data and sample files used by tests across the project
- **`integration/`** - Integration tests that test complete workflows and CLI behavior  
- **`benchmarks/`** - Performance benchmarks for the CLI application

## Root Level Test Files

- **`chatmate_validation_test.go`** - Validation tests for chatmate files and project structure

## Running Tests

```bash
# Run all tests (recommended)
go test ./...

# Run specific test suites
go test ./cmd/...              # Unit tests for CLI commands
go test ./internal/manager/... # Manager package tests
go test ./pkg/...              # Package tests (security, utils)
go test ./test/integration/... # Integration tests
go test ./test/...             # Validation and external tests

# Run with coverage
go test ./... -cover

# Run benchmarks
go test -bench=. ./test/benchmarks/...

# Run specific test patterns
go test -run TestHire ./...    # All tests matching "TestHire"
```

## Test Organization

The project follows Go's standard testing conventions with a clean separation:

### Unit Tests (Package Level)
- **`cmd/*_test.go`** - CLI command unit tests (structure, flags, validation)
- **`internal/manager/*_test.go`** - Manager business logic tests
- **`pkg/security/*_test.go`** - Security validation tests (96.4% coverage)
- **`pkg/utils/*_test.go`** - Utility function tests

### Integration Tests
- **`test/integration/`** - End-to-end CLI behavior and workflows
- **`test/chatmate_validation_test.go`** - Project structure validation

### Performance Tests
- **`test/benchmarks/`** - Performance benchmarks and optimization tests

### Test Helpers
- **`internal/testing/helpers/`** - Shared test utilities and environment setup

## Go-Only Approach

This project uses a **Go-only** approach with no shell script dependencies:
- All functionality implemented in Go
- Cross-platform compatibility through Go's standard library
- Professional CLI using Cobra framework
- Comprehensive test coverage using Go's testing package
