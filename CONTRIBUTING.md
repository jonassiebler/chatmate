
# Contributing to Chatmate

Welcome to the Chatmate project! 🎉 We're excited that you want to contribute. This guide will help you get started with development, testing, and contributing to the codebase.

## Table of Contents

- [Quick Start](#quick-start)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Code Style](#code-style)
- [Security](#security)
- [Documentation](#documentation)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)

## Quick Start

1. **Fork and clone the repository**

   ```bash
   git clone https://github.com/yourusername/chatmate.git
   cd chatmate
   ```


2. **Install dependencies**

   ```bash
   go mod download
   ```


3. **Build the project**

   ```bash
   go build -o chatmate
   ```


4. **Run tests**

   ```bash
   ./run-tests.sh
   ```


5. **Try your changes**

   ```bash
   ./chatmate --help
   ```

## Development Setup

### Prerequisites

- **Go 1.21 or later** - [Download Go](https://golang.org/dl/)
- **Git** - Version control
- **Make** (optional) - Build automation
- **Bash** - For test scripts (Windows users: use Git Bash or WSL)

### Environment Setup

1. **Set up your development environment**
   ```bash
   # Clone your fork
   git clone https://github.com/yourusername/chatmate.git
   cd chatmate

   # Add upstream remote
   git remote add upstream https://github.com/original/chatmate.git

   # Install dependencies
   go mod download
   ```

2. **Build and test**
   ```bash
   # Build the binary
   go build -o chatmate

   # Run all tests
   ./run-tests.sh

   # Run specific test suite
   ./tests/bats-core/bin/bats tests/unit/
   ```

3. **Set up pre-commit hooks** (recommended)
   ```bash
   # Install pre-commit (if you have it)
   pre-commit install

   # Or manually run checks before committing
   go fmt ./...
   go vet ./...
   go test ./...
   ```

## Project Structure

```
chatmate/
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go            # Root command and global flags
│   ├── hire.go            # Main hire command
│   ├── list.go            # List available chatmates
│   ├── status.go          # Check installation status
│   ├── config.go          # Configuration management
│   ├── tutorial.go        # Interactive tutorials
│   └── completion.go      # Shell completion
├── internal/              # Private application code
│   └── manager/           # Core chatmate management logic
│       └── chatmate.go    # Main manager implementation
├── pkg/                   # Public packages (reusable)
│   ├── security/          # Security validation and scanning
│   │   └── validation.go  # Security checks and validation
│   └── utils/             # Utility functions
│       └── paths.go       # Cross-platform path handling
├── docs/                  # Documentation
│   ├── man/              # Unix man pages
│   ├── USER_GUIDE.md     # Comprehensive user guide
│   ├── INSTALLATION.md   # Installation instructions
│   └── QUICK_START.md    # Getting started quickly
├── scripts/              # Build and utility scripts
│   ├── generate-man-pages.go    # Man page generation
│   ├── install-man-pages.sh     # Man page installation
│   └── install-completions.sh   # Shell completion setup
├── tests/                # Test suite (BATS framework)
│   ├── unit/             # Unit tests
│   ├── integration/      # Integration tests
│   ├── fixtures/         # Test data
│   └── helpers/          # Test utilities
├── mates/                # Built-in chatmate definitions
└── assets/               # Static assets
```
### Key Components

- **cmd/**: CLI interface using Cobra framework
- **internal/manager**: Core business logic for chatmate management
- **pkg/security**: Security validation and scanning
- **pkg/utils**: Cross-platform utilities and helpers
- **tests/**: Comprehensive test suite using BATS

## Development Workflow

### Making Changes

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Follow Go conventions
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**
   ```bash
   # Run all tests
   ./run-tests.sh

   # Run specific tests
   go test ./internal/manager/
   go test ./pkg/security/

   # Test CLI manually
   ./chatmate hire "Test Agent"
   ./chatmate list
   ```

4. **Check code quality**
   ```bash
   # Format code
   go fmt ./...

   # Vet code
   go vet ./...

   # Run security scan
   ./chatmate status --security-scan
   ```

### Adding New Commands

1. **Create command file**
   ```bash
   # Create new command in cmd/
   touch cmd/newcommand.go
   ```

2. **Implement command structure**
   ```go
   package cmd

   import (
       "github.com/spf13/cobra"
   )

   var newcommandCmd = &cobra.Command{
       Use:   "newcommand",
       Short: "Short description",
       Long:  `Detailed description...`,
       Run: func(cmd *cobra.Command, args []string) {
           // Implementation
       },
   }

   func init() {
       rootCmd.AddCommand(newcommandCmd)
   }
   ```

3. **Add tests**
   ```bash
   # Add unit tests
   touch tests/unit/newcommand.bats

   # Add integration tests if needed
   touch tests/integration/newcommand_integration.bats
   ```

### Adding Security Features

1. **Extend security validation**
   - Add new checks to `pkg/security/validation.go`
   - Include comprehensive error handling
   - Add corresponding tests

2. **Update security scanning**
   - Enhance the `SecurityScan` function
   - Add new vulnerability patterns
   - Update security scoring

## Testing

We use a multi-layered testing approach:

### Test Types

1. **Unit Tests** (Go's built-in testing)
   ```bash
   # Run Go unit tests
   go test ./...

   # With coverage
   go test -cover ./...

   # Specific package
   go test ./internal/manager/
   ```

2. **CLI Tests** (BATS framework)
   ```bash
   # Run all BATS tests
   ./run-tests.sh

   # Run specific test file
   ./tests/bats-core/bin/bats tests/unit/hire_script.bats

   # Run with verbose output
   ./tests/bats-core/bin/bats -t tests/unit/
   ```

3. **Integration Tests**
   ```bash
   # Full integration test suite
   ./tests/bats-core/bin/bats tests/integration/
   ```

### Writing Tests
#### BATS Test Example
```bash
#!/usr/bin/env bats

load '../helpers/test_helper'

@test "hire command creates chatmate file" {
    run ./chatmate hire "Test Agent"
    assert_success
    assert_output --partial "successfully hired"
}

@test "hire command with invalid input fails" {
    run ./chatmate hire ""
    assert_failure
    assert_output --partial "name is required"
}
```
#### Go Test Example
```go
func TestManager_HireChatmate(t *testing.T) {
    manager := &ChatmateManager{}

    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid name", "Test Agent", false},
        {"empty name", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := manager.HireChatmate(tt.input, "", "")
            if (err != nil) != tt.wantErr {
                t.Errorf("HireChatmate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Code Style

### Go Style Guidelines

1. **Follow standard Go conventions**
   - Use `gofmt` for formatting
   - Follow effective Go practices
   - Use meaningful variable names

2. **Documentation**
   - Add package documentation
   - Document all public functions
   - Include usage examples in GoDoc

3. **Error Handling**
   ```go
   // Good: Specific error messages
   if err != nil {
       return fmt.Errorf("failed to create chatmate %s: %w", name, err)
   }

   // Bad: Generic error handling
   if err != nil {
       return err
   }
   ```

4. **CLI Output**
   - Use emojis for visual appeal (✅ ❌ 🔍 🚀)
   - Provide clear, actionable messages
   - Include progress indicators for long operations

### Code Organization

1. **Package Structure**
   - Keep packages focused and cohesive
   - Minimize dependencies between packages
   - Use internal/ for private code

2. **Function Design**
   - Single responsibility principle
   - Clear input/output contracts
   - Comprehensive error handling


## Documentation

### Types of Documentation

1. **API Documentation** (GoDoc)
   - Package-level documentation
   - Function documentation with examples
   - Usage patterns and best practices

2. **User Documentation**
   - CLI help system (comprehensive)
   - User guide with examples
   - Installation instructions
   - Quick start guide

3. **Developer Documentation**
   - This contributing guide
   - Architecture decisions
   - API reference

### Writing Documentation

1. **CLI Help**
   ```go
   var cmd = &cobra.Command{
       Use:   "command",
       Short: "Brief description with emoji 🚀",
       Long: `Detailed description with:

   • Clear benefits
   • Usage examples
   • Common workflows
   • Troubleshooting tips`,
       Example: `  # Basic usage
     chatmate command arg1

     # Advanced usage
     chatmate command --flag value arg1`,
   }
   ```

2. **GoDoc Comments**
   ```go
   // Package manager provides core chatmate management functionality.
   //
   // This package handles chatmate lifecycle operations including hiring,
   // listing, and status checking with comprehensive validation and
   // security scanning.
   package manager

   // HireChatmate creates a new chatmate with the specified configuration.
   //
   // The function validates the input, creates the necessary directory
   // structure, and generates the chatmate file with proper permissions.
   //
   // Example:
   //   manager := &ChatmateManager{}
   //   err := manager.HireChatmate("Code Assistant", "helpful", "coding")
   //   if err != nil {
   //       log.Fatal(err)
   //   }

   func (cm *ChatmateManager) HireChatmate(name, personality, expertise string) error {
   ```

## Submitting Changes

### Pull Request Process

1. **Prepare your changes**
   ```bash
   # Update your branch
   git fetch upstream
   git rebase upstream/main

   # Run full test suite
   ./run-tests.sh

   # Check code quality
   go fmt ./...
   go vet ./...
   ```

2. **Create pull request**
   - Use clear, descriptive title
   - Explain the problem and solution
   - Include test results
   - Link any related issues

3. **PR Template**
   ```markdown
   ## Description
   Brief description of changes

   ## Changes Made
   - [ ] Feature/fix implemented
   - [ ] Tests added/updated
   - [ ] Documentation updated
   - [ ] Security considerations addressed

   ## Testing
   - [ ] Unit tests pass
   - [ ] Integration tests pass
   - [ ] Manual testing completed

   ## Security
   - [ ] No new security vulnerabilities
   - [ ] Input validation implemented
   - [ ] File permissions appropriate
   ```

### Code Review

- Be responsive to feedback
- Make requested changes promptly
- Ask questions if requirements are unclear
- Keep discussions professional and constructive

## Release Process

## Homebrew Tap Publishing (Maintainers)

After releasing a new version, update and publish the Homebrew tap:

1. Follow the release process above to tag and push a new version.
2. Update `homebrew-tap/Formula/chatmate.rb` with the new commit hash and version.
3. Test the formula locally:
   ```bash
   brew uninstall chatmate
   brew install --build-from-source ./homebrew-tap/Formula/chatmate.rb
   chatmate --help
   ```
4. Push the updated formula to the repository (and/or the separate tap repo if used).
5. See [docs/HOMEBREW_PUBLISHING.md](docs/HOMEBREW_PUBLISHING.md) for a detailed step-by-step guide.

Consider automating this process with a GitHub Action for future releases.

### Version Management

We use semantic versioning (SemVer):
- **Major** (X.0.0): Breaking changes
- **Minor** (0.X.0): New features, backwards compatible
- **Patch** (0.0.X): Bug fixes, backwards compatible

### Release Checklist

1. **Pre-release**
   - [ ] All tests passing
   - [ ] Documentation updated
   - [ ] Security scan clean
   - [ ] Performance benchmarks stable

2. **Release**
   - [ ] Version bumped
   - [ ] Changelog updated
   - [ ] Git tag created
   - [ ] Binaries built and tested

3. **Post-release**
   - [ ] Release notes published
   - [ ] Documentation deployed
   - [ ] Community notified

## Getting Help

### Resources

- **Documentation**: Check docs/ directory
- **Issues**: Search existing GitHub issues
- **Discussions**: Use GitHub Discussions for questions
- **Code Examples**: Look at tests/ for usage examples

### Communication

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and ideas
- **Pull Requests**: Code contributions and reviews

### Development Tips

1. **Debugging**
   ```bash
   # Verbose output
   ./chatmate --verbose hire "Test Agent"

   # Debug specific operations
   go run . --debug list
   ```

2. **Performance Testing**
   ```bash
   # Benchmark tests
   go test -bench=. ./...

   # Memory profiling
   go test -memprofile=mem.prof ./...
   ```

3. **Cross-platform Testing**
   ```bash
   # Test on different OS (if available)
   GOOS=linux go build
   GOOS=windows go build
   GOOS=darwin go build
   ```


## Thank You! 🙏
Thank you for contributing to Chatmate! Your contributions help make AI assistance more accessible and effective for developers worldwide.
For questions or support, please open an issue or start a discussion on GitHub. We're here to help and excited to see what you build! 🚀
