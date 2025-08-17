# ChatMate 🤖

Open source collection of specialized AI agents for VS Code Copilot Chat. Each chatmate is a carefully crafted prompt designed to excel at specific development tasks.

[![Go](https://github.com/jonassiebler/chatmate/actions/workflows/go.yml/badge.svg)](https://github.com/jonassiebler/chatmate/actions/workflows/go.yml)
[![Security](https://github.com/jonassiebler/chatmate/actions/workflows/security.yml/badge.svg)](https://github.com/jonassiebler/chatmate/actions/workflows/security.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jonassiebler/chatmate)](https://goreportcard.com/report/github.com/jonassiebler/chatmate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> 🎉 **Recent Update**: ChatMate has been completely rewritten in Go with native testing, enhanced CLI interface, comprehensive security validation, and cross-platform support. All features are backward compatible with improved performance and reliability.

## 🚀 Quick Start

### Install the CLI Tool


```bash
# Download or build the chatmate CLI
# Then install all chatmates
chatmate hire
```

**That's it!** Restart VS Code and start using your new chatmates with `@ChatmateName` in Copilot Chat.

### Verify Installation


```bash
chatmate status    # Check installation health
chatmate list      # See installed chatmates
```

## 🤖 What are Chatmates?

Chatmates are specialized AI agents that bring focused expertise to your VS Code Copilot Chat:

| Chatmate | Purpose | Use Case |
|----------|---------|----------|
| **Solve Issue** 🐛 | Systematic debugging | "@Solve Issue My React component won't render" |
| **Code Review** 👁️ | Expert code analysis | "@Code Review Check this authentication logic" |
| **Testing** 🧪 | Test generation & debugging | "@Testing Generate unit tests for this service" |
| **Create PR** 📝 | Pull request creation | "@Create PR Create PR for the new auth feature" |
| **Documentation** 📚 | Technical writing | "@Documentation Write API docs for this endpoint" |
| **Create Issue** 🎯 | GitHub issue creation | "@Create Issue User login fails on mobile" |
| ... and more! | | Run `chatmate list` for all available agents |

## 📋 Command Reference

### Essential Commands


```bash
# Install all chatmates (recommended)
chatmate hire

# Install specific chatmates
chatmate hire "Solve Issue" "Code Review" "Testing"

# Force reinstall/update all chatmates
chatmate hire --force

# List available and installed chatmates
chatmate list

# Show installation status and system info
chatmate status

# View configuration and paths
chatmate config

# Remove specific chatmates
chatmate uninstall "Unused Chatmate"

# Get help
chatmate --help
chatmate hire --help
```

### Example Workflows


```bash
# First-time setup
chatmate hire
# → Restart VS Code → Use @ChatmateName in Copilot Chat

# Daily development
chatmate status          # Check health
# → Use appropriate chatmates for tasks

# Team coordination
chatmate list --installed > team-setup.txt
# → Share configuration with team

# Maintenance
chatmate hire --force    # Update all chatmates
chatmate uninstall "Rarely Used"  # Clean up
```

## 📖 Documentation

- **[Complete User Guide](docs/USER_GUIDE.md)** - Comprehensive documentation with examples, troubleshooting, and best practices
- **[Man Pages](docs/man/)** - Unix man pages for all commands (`man chatmate`)
- **Command Help** - Run `chatmate --help` or `chatmate [command] --help`

## 🛠 Requirements

- **VS Code** with GitHub Copilot Chat extension
- **Write permissions** to VS Code user directory
- **Supported platforms**: macOS, Linux, Windows

## 🔧 Installation Methods

git clone https://github.com/jonassiebler/chatmate.git
go build -o chatmate .

### Homebrew Tap Installation (Recommended)

Install ChatMate easily via Homebrew:

```bash
# Add the ChatMate tap
brew tap jonassiebler/chatmate

# Install the CLI
brew install chatmate

# (Optional) Update to latest version
brew upgrade chatmate

# Install all chatmates
chatmate hire
```

**That's it!** Restart VS Code and start using your new chatmates with `@ChatmateName` in Copilot Chat.

---

### CLI Installation (Manual)

Build and install the ChatMate CLI tool manually:

```bash
# Clone the repository
git clone https://github.com/jonassiebler/chatmate.git
cd chatmate

# Build the CLI
go build -o chatmate .

# Install chatmates
./chatmate hire
```

### Legacy Script Installation

For users who prefer the original method:


```bash
git clone https://github.com/jonassiebler/chatmate.git
cd chatmate
./hire.sh
```

### Manual Installation

Install specific chatmates manually:

1. **Find your VS Code prompts directory:**
   - **macOS:** `~/Library/Application Support/Code/User/prompts`
   - **Linux:** `~/.config/Code/User/prompts`
   - **Windows:** `%APPDATA%/Code/User/prompts`

2. **Copy chatmate files:**

   ```bash
   cp mates/"Solve Issue.chatmode.md" "/path/to/prompts/"
   ```

## 🚨 Troubleshooting

### Common Issues

- **Chatmates not appearing?** → Restart VS Code completely
- **Permission denied?** → Check directory permissions with `chatmate config`
- **VS Code not detected?** → Ensure VS Code is in your PATH (`code --version`)

### Getting Help


```bash
# Diagnostic commands
chatmate status    # System health check
chatmate config    # Configuration details
chatmate --help    # Command help
```

For detailed troubleshooting, see the [User Guide](docs/USER_GUIDE.md#troubleshooting).

## 🧪 Development & Testing

ChatMate uses Go's native testing framework with Testify for comprehensive test coverage:


```bash
# Run all tests with coverage
./run-tests.sh

# Run only unit tests
./run-tests.sh --unit

# Run only integration tests
./run-tests.sh --integration

# Run with benchmarks
./run-tests.sh --benchmark

# Generate coverage report
./run-tests.sh --coverage

# Run quality checks
./run-tests.sh --quality

# Run security checks
./run-tests.sh --security

# Quick Go test commands
go test -v ./...                    # Run all tests
go test -cover ./...               # Run with coverage
go test -bench=. ./...             # Run benchmarks
```


### Test Structure

- **Unit Tests**: Located alongside source code as `*_test.go` files
- **Integration Tests**: End-to-end functionality testing
- **Test Helpers**: `internal/testing/helpers/` - shared testing utilities
- **Test Fixtures**: `internal/testing/fixtures/` - sample test data

## 🤝 Contributing

We welcome contributions! Whether it's new chatmates, bug fixes, or documentation improvements.

### Quick Contributing Guide

1. Fork the repository
2. Create a feature branch from `dev`
3. Make your changes and test thoroughly
4. Submit a pull request to `dev` branch

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

### Creating Custom Chatmates

1. Use the **@Create Chatmode** agent to design new chatmates
2. Add your `.chatmode.md` files to the `mates/` directory
3. Install with `chatmate hire`
4. Share with the community via pull request

## 📊 Quality Assurance

ChatMate maintains high quality through:

- **Automated testing** on multiple platforms
- **Security scanning** for vulnerabilities
- **Code quality checks** with linting and static analysis
- **Performance optimization** for fast, efficient operation
- **Comprehensive documentation** and user guides

## 🏗 Repository Structure


```text
chatmate/
├── cmd/                # CLI commands and interfaces
├── internal/           # Core application logic
│   ├── assets/         # Embedded resources
│   ├── manager/        # Chatmate management
│   └── testing/        # Test helpers and fixtures
├── pkg/               # Public packages (utils, security)
├── mates/             # Chatmate agent definitions
├── docs/              # Documentation and man pages
├── scripts/           # Build and utility scripts
├── *_test.go          # Go test files
├── run-tests.sh       # Test runner script
├── TEST_CONFIG.md     # Testing documentation
└── README.md          # This file
```

## 📋 Platform Support

- ✅ **macOS** (Intel & Apple Silicon)
- ✅ **Linux** (x64, ARM64)
- ✅ **Windows** (x64, Git Bash/WSL)

All platforms tested in CI/CD pipeline with comprehensive integration tests.

## 🔐 Security

ChatMate prioritizes security with:

- **Input validation** for all file operations
- **Path traversal protection** for safe file handling
- **Vulnerability scanning** in CI/CD pipeline
- **Code signing** for release binaries (planned)
- **Regular dependency updates** for security patches

Security Grade: **A** ✅ ([View Security Report](security-reports/))

## 📄 License

MIT License - Feel free to use these chatmates in your own projects!

## 🆘 Support & Community

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/jonassiebler/chatmate/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/jonassiebler/chatmate/discussions)
- � **Documentation**: [User Guide](docs/USER_GUIDE.md) | `man chatmate`
- 🤝 **Contributing**: [CONTRIBUTING.md](CONTRIBUTING.md)

---

**Ready to supercharge your VS Code experience?** 🚀


```bash
chatmate hire
```

Then restart VS Code and start using `@ChatmateName` in Copilot Chat!
