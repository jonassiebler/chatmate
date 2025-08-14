# Contributing to chatmate 🤖

Thank you for your interest in contributing to chatmate! This guide will help you understand our development workflow and contribution standards.

## Quick Start

1. **Fork** the repository on GitHub
2. **Clone** your fork locally and add upstream remote:

   ```bash
   git clone https://github.com/your-username/chatmate.git
   cd chatmate
   git remote add upstream https://github.com/jonassiebler/chatmate.git
   ```

3. **Test the installation**: `./hire.sh`

## Development Workflow

### Branch Strategy

- **`main`** - Production branch (protected, stable releases)
- **`dev`** - Development integration branch (protected)  
- **`feature/*`** - Feature branches (your contributions)

### Contribution Process

1. **Create feature branch** from `dev`:

   ```bash
   git checkout dev
   git pull upstream dev
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following these guidelines:
   - **New chatmates**: Place in `mates/` directory with `.chatmode.md` extension
   - **Bug fixes**: Target the specific issue thoroughly
   - **Documentation**: Keep clear and concise

3. **Test thoroughly**:

   ```bash
   ./hire.sh  # Test installation
   # Restart VS Code and verify chatmates work
   ```

4. **Submit pull request** to `dev` branch using the PR template

## Chatmate Guidelines

### File Naming

- Use descriptive names: `Your Chatmate Name.chatmode.md`
- Follow existing patterns in the `mates/` directory

### Content Standards

- **Clear purpose**: Each chatmate should have a specific, well-defined role
- **Focused prompts**: Avoid overly broad or generic instructions
- **Tested functionality**: Ensure the chatmate works as expected in VS Code
- **Unique value**: Don't duplicate existing functionality

### Quality Checklist

- [ ] File follows naming convention
- [ ] Chatmate has clear, specific purpose
- [ ] Prompt is well-structured and focused
- [ ] Tested in VS Code Copilot Chat
- [ ] Installation script works correctly
- [ ] No conflicts with existing chatmates

## Review Process

- All changes require pull requests to protected branches
- Admin review is mandatory before merging
- Automated tests must pass (structure validation, security scans, etc.)
- Follow the pull request template for complete submissions

## Code of Conduct

- **Be respectful** and inclusive in all interactions
- **Test thoroughly** before submitting changes  
- **Follow guidelines** outlined in this document
- **Respond promptly** to review feedback
- **Keep PRs focused** - one feature/fix per pull request

## Issue Reporting

When reporting bugs or requesting features:

1. **Search existing issues** first to avoid duplicates
2. **Provide detailed information**:
   - Clear description of the issue/request
   - Steps to reproduce (for bugs)
   - Environment details (OS, VS Code version, etc.)
   - Expected vs actual behavior
3. **Be patient** for responses from maintainers

## Repository Structure

```text
chatmate/
├── .github/               # GitHub workflows and templates
├── mates/                # Chatmate definitions (.chatmode.md files)
├── CONTRIBUTING.md       # This file
├── LICENSE              # MIT license
├── README.md            # Main documentation
└── hire.sh              # Installation script
```

## Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/jonassiebler/chatmate/issues)
- 💡 **Feature Requests**: [GitHub Issues](https://github.com/jonassiebler/chatmate/issues)
- ❓ **Questions**: [GitHub Issues](https://github.com/jonassiebler/chatmate/issues) or [GitHub Discussions](https://github.com/jonassiebler/chatmate/discussions)

## Maintainer Response Times

We strive to:

- **Acknowledge issues** within 48 hours
- **Review PRs** within 1 week  
- **Respond to discussions** within 72 hours

---

Thank you for contributing to chatmate! Together, we're building the best collection of AI assistants for developers. 🚀
