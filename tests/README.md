# ChatMate Testing Framework

Comprehensive testing infrastructure for the ChatMate project, ensuring reliability, quality, and cross-platform compatibility.

## Overview

This testing framework provides:

- ✅ **Unit Testing** - Individual component validation
- ✅ **Integration Testing** - End-to-end workflow verification  
- ✅ **Quality Assurance** - Code quality and formatting checks
- ✅ **Cross-Platform Testing** - macOS, Linux, and Windows compatibility
- ✅ **Security Scanning** - Sensitive data detection
- ✅ **Automated CI/CD** - GitHub Actions integration

## Quick Start

### Prerequisites

- **Bash** (for running tests)
- **Git** (for framework installation)
- **Node.js** (optional, for additional quality tools)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/jonassiebler/chatmate.git
   cd chatmate
   ```

2. **Install testing dependencies:**
   ```bash
   npm install  # Installs markdownlint and other tools
   ```

3. **Run all tests:**
   ```bash
   ./run-tests.sh
   ```

## Test Framework Architecture

```
tests/
├── unit/                    # Unit tests for individual components
│   ├── hire_script.bats     # Installation script tests
│   └── chatmate_validation.bats  # Chatmate file validation tests
├── integration/             # End-to-end integration tests
│   └── integration_tests.bats    # Full workflow tests
├── fixtures/                # Test data and mock files
│   ├── sample.chatmode.md   # Valid test chatmate
│   ├── invalid_*.md         # Invalid test cases
│   └── empty.chatmode.md    # Edge case test file
├── helpers/                 # Shared test utilities
│   └── test_helper.bash     # Common test functions
├── bats-core/              # Bats testing framework (auto-installed)
├── bats-support/           # Bats support library (auto-installed)
├── bats-assert/            # Bats assertion library (auto-installed)
└── test.config             # Testing configuration
```

## Testing Categories

### 1. Unit Tests

**Purpose**: Test individual components in isolation

**Location**: `tests/unit/`

**What's Tested**:
- `hire.sh` script functionality
- Chatmate file format validation
- OS detection and path resolution
- Error handling scenarios

**Example**:
```bash
./run-tests.sh --unit-only
```

### 2. Integration Tests

**Purpose**: Test complete workflows and system interactions

**Location**: `tests/integration/`

**What's Tested**:
- Full installation process
- Cross-platform compatibility
- File copying and permissions
- Repository structure validation
- Security scanning

**Example**:
```bash
./run-tests.sh --integration-only
```

### 3. Quality Tests

**Purpose**: Ensure code quality and formatting standards

**What's Tested**:
- Markdown formatting (markdownlint)
- Shell script quality (shellcheck)
- File structure compliance
- Content validation

**Example**:
```bash
./run-tests.sh --quality-only
```

## Running Tests

### Basic Usage

```bash
# Run all tests
./run-tests.sh

# Run specific test categories
./run-tests.sh --unit-only
./run-tests.sh --integration-only
./run-tests.sh --quality-only

# Verbose output
./run-tests.sh --verbose
```

### Manual Bats Usage

```bash
# Run specific test file
tests/bats-core/bin/bats tests/unit/hire_script.bats

# Run all unit tests
tests/bats-core/bin/bats tests/unit/*.bats

# TAP output format
tests/bats-core/bin/bats --tap tests/unit/*.bats
```

### Using NPM Scripts

```bash
# Run all tests
npm test

# Run specific categories
npm run test:shell
npm run test:markdown
npm run test:integration

# Setup testing environment
npm run setup
```

## Test Configuration

Configure testing behavior by editing `tests/test.config`:

```bash
# Enable/disable test categories
UNIT_TESTS=true
INTEGRATION_TESTS=true
QUALITY_TESTS=true

# Set coverage thresholds
COVERAGE_THRESHOLD=80

# Platform testing
TEST_MACOS=true
TEST_LINUX=true
TEST_WINDOWS=true
```

## Writing New Tests

### Unit Test Example

Create a new file in `tests/unit/` with `.bats` extension:

```bash
#!/usr/bin/env bats

load '../helpers/test_helper'

@test "my new test case" {
    # Arrange
    local test_file="$TEST_TEMP_DIR/test.md"
    echo "# Test Content" > "$test_file"
    
    # Act
    run grep "Test Content" "$test_file"
    
    # Assert
    assert_success
    assert_output "# Test Content"
}
```

### Integration Test Example

```bash
#!/usr/bin/env bats

load '../helpers/test_helper'

@test "complete workflow test" {
    # Setup test environment
    simulate_os "linux"
    local prompts_dir="$MOCK_HOME/.config/Code/User/prompts"
    
    # Test installation
    run bash -c "
        mkdir -p '$prompts_dir'
        cp -v mates/*.md '$prompts_dir'/
    "
    assert_success
    
    # Verify results
    assert verify_installation "$prompts_dir" "$(count_files "mates" "*.chatmode.md")"
}
```

## Helper Functions

The testing framework provides several helper functions in `tests/helpers/test_helper.bash`:

```bash
# File and directory helpers
count_files "$dir" "$pattern"          # Count files matching pattern
file_contains "$file" "$content"       # Check if file contains content
is_valid_chatmate "$file"             # Validate chatmate format

# Environment helpers
simulate_os "$os_type"                 # Mock different operating systems
create_test_chatmate "$name" "$content" # Create test chatmate files
verify_installation "$dir" "$count"    # Verify installation results

# Test environment
setup()    # Run before each test
teardown() # Run after each test
```

## Continuous Integration

### GitHub Actions

Tests automatically run on:
- **Push** to `main` or `dev` branches
- **Pull requests** to `main` or `dev` branches

### Workflow Jobs

1. **test-framework**: Core Bats testing on Ubuntu and macOS
2. **security-scan**: Security and sensitive data checks
3. **markdown-quality**: Markdown formatting validation
4. **shell-quality**: Shell script quality checks
5. **cross-platform-test**: Multi-platform compatibility
6. **test-coverage**: Coverage analysis and reporting

### Viewing Results

Check test results in:
- GitHub Actions tab in your repository
- Test artifacts uploaded after each run
- `test-results.log` file generated locally

## Troubleshooting

### Common Issues

**Bats not found:**
```bash
# Install Bats manually
./run-tests.sh  # This will auto-install Bats
```

**Permission denied:**
```bash
chmod +x run-tests.sh
chmod +x hire.sh
```

**Tests failing on macOS:**
```bash
# Install additional tools via Homebrew
brew install shellcheck
```

**Node.js dependencies:**
```bash
npm install  # Install markdownlint and other tools
```

### Debug Mode

Run tests with verbose output:
```bash
./run-tests.sh --verbose
```

Check detailed logs:
```bash
cat test-results.log
```

## Test Coverage

Current testing covers:

- ✅ **Installation Script** - All OS types and scenarios
- ✅ **Chatmate Validation** - Format, content, and quality
- ✅ **File Operations** - Copying, permissions, and paths
- ✅ **Error Handling** - Graceful failure scenarios
- ✅ **Cross-Platform** - macOS, Linux, Windows compatibility
- ✅ **Security** - Sensitive data detection
- ✅ **Quality** - Code formatting and standards

### Coverage Goals

- **Unit Test Coverage**: >95%
- **Integration Coverage**: >90%
- **Platform Coverage**: 100% (macOS, Linux, Windows)
- **Quality Gates**: 100% compliance

## Contributing

### Adding New Tests

1. **Identify the component** to test
2. **Choose the appropriate category** (unit/integration/quality)
3. **Write the test** following existing patterns
4. **Use helper functions** for consistency
5. **Test your test** locally
6. **Submit a pull request**

### Test Standards

- **One concept per test** - Keep tests focused
- **Descriptive names** - Test names should explain what's being tested
- **Arrange-Act-Assert** - Follow the AAA pattern
- **Use helpers** - Leverage shared utilities
- **Clean up** - Tests should not affect each other

## License

This testing framework is part of the ChatMate project and follows the same MIT License.
