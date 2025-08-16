# ChatMate User Guide 📖

Welcome to the comprehensive ChatMate user guide! This document will help you master the ChatMate CLI and get the most out of your specialized AI agents for VS Code Copilot Chat.

## Table of Contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Getting Started](#getting-started)
4. [Command Reference](#command-reference)
5. [Chatmate Catalog](#chatmate-catalog)
6. [Usage Patterns](#usage-patterns)
7. [Troubleshooting](#troubleshooting)
8. [Best Practices](#best-practices)
9. [Advanced Usage](#advanced-usage)
10. [FAQ](#faq)

## Introduction

### What is ChatMate?

ChatMate is a CLI tool that manages specialized AI agents (called "chatmates") for VS Code Copilot Chat. Each chatmate is a carefully crafted prompt designed to excel at specific development tasks, transforming your Copilot Chat experience into a powerhouse of specialized expertise.

### Why Use ChatMate?

- **Specialized Expertise**: Each chatmate brings focused knowledge to specific tasks
- **Consistent Quality**: Well-tested prompts ensure reliable, high-quality responses
- **Easy Management**: Simple CLI commands to install, update, and manage chatmates
- **Seamless Integration**: Works natively with VS Code Copilot Chat
- **Open Source**: Community-driven development with transparent, extensible agents

## Installation

### Prerequisites

Before installing ChatMate, ensure you have:

- **VS Code** installed and accessible from your PATH
- **VS Code Copilot Chat extension** enabled
- **Write permissions** to your VS Code user directory
- **Go 1.19+** (if building from source)

### Quick Installation

The fastest way to get started:

```bash
# Install or download the ChatMate CLI
# Then install all chatmates
chatmate hire
```

### Installation Verification

After installation, verify everything is working:

```bash
# Check system status
chatmate status

# List installed chatmates
chatmate list --installed

# View configuration
chatmate config
```

You should see output indicating VS Code is detected and chatmates are installed.

## Getting Started

### Your First Steps

1. **Install ChatMate and all chatmates:**
   ```bash
   chatmate hire
   ```

2. **Restart VS Code** to load the new chatmates

3. **Open Copilot Chat** in VS Code (`Ctrl/Cmd+Shift+P` → "Chat: Open Chat")

4. **Try a chatmate:**
   - Type `@Solve Issue` to get help with debugging
   - Type `@Code Review` for code analysis
   - Type `@Testing` for test generation

### Basic Workflow

The typical ChatMate workflow:

```bash
# 1. Check status and available chatmates
chatmate status
chatmate list

# 2. Install specific chatmates you need
chatmate hire "Solve Issue" "Testing" "Code Review"

# 3. Use in VS Code Copilot Chat
# Open VS Code, use Copilot Chat with @ChatmateName

# 4. Update or manage chatmates as needed
chatmate hire --force  # Update all
chatmate uninstall "Unused Chatmate"  # Remove specific ones
```

## Command Reference

### `chatmate hire`

Install chatmate agents to your VS Code setup.

**Syntax:**
```bash
chatmate hire [chatmate names...] [flags]
```

**Options:**
- `--force, -f`: Force reinstall existing chatmates
- `--specific, -s`: Install specific chatmates by name (alternative to args)
- `--help`: Show help for the hire command

**Examples:**
```bash
# Install all available chatmates (recommended for first-time users)
chatmate hire

# Install specific chatmates by name
chatmate hire "Solve Issue" "Code Review" "Testing"

# Force reinstall all chatmates (useful after updates)
chatmate hire --force

# Force reinstall specific chatmates
chatmate hire --force "Solve Issue" "Testing"

# Using the --specific flag (alternative syntax)
chatmate hire --specific "Code Review" --specific "Documentation"
```

**What it does:**
1. Validates VS Code installation and prompts directory
2. Copies chatmate files to VS Code user prompts directory
3. Handles existing files with smart overwrite logic
4. Reports installation status and any conflicts

### `chatmate list`

Display information about available and installed chatmate agents.

**Syntax:**
```bash
chatmate list [flags]
```

**Options:**
- `--available, -a`: Show only available chatmates
- `--installed, -i`: Show only installed chatmates
- `--help`: Show help for the list command

**Examples:**
```bash
# List all chatmates with installation status (default)
chatmate list

# Show only available chatmates (not yet installed)
chatmate list --available

# Show only installed chatmates
chatmate list --installed

# Combine with grep for filtering
chatmate list --available | grep "Testing"  # Find testing-related chatmates
```

**Output format:**
- ✅ **Installed chatmates**: Green checkmark with "installed" status
- ❌ **Available chatmates**: Red X with "not installed" status
- 📊 **Summary**: Count of installed vs available chatmates

### `chatmate status`

Show comprehensive ChatMate installation status and system information.

**Syntax:**
```bash
chatmate status
```

**Examples:**
```bash
# Show complete ChatMate installation status
chatmate status

# Save status info for support requests
chatmate status > chatmate-status.txt

# Common troubleshooting workflow
chatmate status          # Check system health
chatmate list           # Verify chatmate availability
chatmate hire --force   # Force reinstall if needed
```

**Information provided:**
- VS Code installation detection and path
- ChatMate prompts directory location and permissions
- Count of installed vs available chatmates
- System platform and environment details
- Integration health status

### `chatmate uninstall`

Remove chatmate agents from your VS Code setup.

**Syntax:**
```bash
chatmate uninstall [chatmate names...] [flags]
```

**Options:**
- `--all`: Uninstall all chatmates
- `--help`: Show help for the uninstall command

**Examples:**
```bash
# Uninstall a specific chatmate
chatmate uninstall "Solve Issue"

# Uninstall multiple chatmates at once
chatmate uninstall "Create PR" "Merge PR" "Review PR"

# Uninstall all chatmates (nuclear option)
chatmate uninstall --all

# Common workflow: check what's installed, then remove unused ones
chatmate list --installed
chatmate uninstall "Documentation" "Optimize Issues"
```

**What happens:**
- Chatmate files are removed from VS Code prompts directory
- Existing chat history and conversations are preserved
- You can always reinstall chatmates later with `chatmate hire`

### `chatmate config`

Display detailed ChatMate configuration information.

**Syntax:**
```bash
chatmate config [flags]
```

**Options:**
- `--show, -s`: Show configuration details (default: true)
- `--help`: Show help for the config command

**Examples:**
```bash
# Show complete configuration information
chatmate config

# Save configuration for support requests
chatmate config > chatmate-config.txt

# Common troubleshooting workflow
chatmate config    # Check paths and configuration
chatmate status    # Verify system integration
chatmate list      # Test chatmate discovery
```

**Configuration details:**
- ChatMate installation directory and embedded resources
- VS Code user directory and prompts path
- Platform-specific paths and conventions
- Environment variables and system settings
- File permissions and accessibility information

### Global Options

All commands support these global options:

- `--verbose, -v`: Enable verbose output for debugging
- `--help, -h`: Show help information
- `--version`: Show version information

## Chatmate Catalog

### Development & Debugging

#### **Solve Issue** 🐛
- **Purpose**: Systematic debugging and problem resolution
- **Best for**: Bug analysis, troubleshooting, root cause investigation
- **Usage**: `@Solve Issue` followed by your problem description
- **Example**: "@Solve Issue My React component isn't rendering properly"

#### **Code Review** 👁️
- **Purpose**: Expert code analysis and improvement suggestions
- **Best for**: Code quality assessment, security reviews, best practices
- **Usage**: `@Code Review` with your code
- **Example**: "@Code Review Please review this authentication function"

#### **Testing** 🧪
- **Purpose**: Comprehensive test generation and debugging
- **Best for**: Unit tests, integration tests, test strategy
- **Usage**: `@Testing` with your code or requirements
- **Example**: "@Testing Generate unit tests for this service class"

### Project Management

#### **Create PR** 📝
- **Purpose**: Generate comprehensive pull requests
- **Best for**: PR creation, description writing, change documentation
- **Usage**: `@Create PR` with your changes
- **Example**: "@Create PR Create a PR for the new authentication feature"

#### **Review PR** 🔍
- **Purpose**: Analyze pull requests for quality and best practices
- **Best for**: PR reviews, change assessment, feedback generation
- **Usage**: `@Review PR` with PR details
- **Example**: "@Review PR Review this authentication enhancement PR"

#### **Create Issue** 🎯
- **Purpose**: Create well-structured GitHub issues
- **Best for**: Bug reports, feature requests, task documentation
- **Usage**: `@Create Issue` with issue details
- **Example**: "@Create Issue User authentication sometimes fails on mobile"

#### **Optimize Issues** ⚡
- **Purpose**: Improve existing issues for clarity and actionability
- **Best for**: Issue refinement, requirement clarification
- **Usage**: `@Optimize Issues` with existing issue content
- **Example**: "@Optimize Issues Make this bug report more actionable"

#### **Merge PR** 🔀
- **Purpose**: Handle pull request merging and cleanup
- **Best for**: Merge decisions, conflict resolution, post-merge tasks
- **Usage**: `@Merge PR` with merge context
- **Example**: "@Merge PR Help me merge this complex feature branch"

### Documentation & Content

#### **Documentation** 📚
- **Purpose**: Technical writing and API documentation
- **Best for**: README files, API docs, user guides, inline comments
- **Usage**: `@Documentation` with your code or topic
- **Example**: "@Documentation Write API documentation for this REST endpoint"

### Specialized Assistants

#### **Code Claude Sonnet 4** 🎯
- **Purpose**: Advanced coding assistant optimized for complex development tasks
- **Best for**: Architecture decisions, complex algorithms, system design
- **Usage**: `@Code Claude Sonnet 4` with development questions
- **Example**: "@Code Claude Sonnet 4 Design a scalable microservice architecture"

#### **Code GPT-4.1** 💡
- **Purpose**: Versatile coding assistant for general development needs
- **Best for**: Code generation, explanations, quick solutions
- **Usage**: `@Code GPT-4.1` with coding questions
- **Example**: "@Code GPT-4.1 Explain this complex algorithm"

### Meta & Customization

#### **Create Chatmode** 🛠️
- **Purpose**: Help create new custom Copilot Chat agents
- **Best for**: Building specialized chatmates, prompt engineering
- **Usage**: `@Create Chatmode` with requirements for new agent
- **Example**: "@Create Chatmode Create a chatmate specialized in database optimization"

## Usage Patterns

### For Individual Developers

**Daily Development Workflow:**
```bash
# Morning setup - check status and install needed chatmates
chatmate status
chatmate list --installed

# Work on specific task - use appropriate chatmates in VS Code
# @Solve Issue for debugging
# @Code Review for quality checks
# @Testing for test creation

# End of day - clean up unused chatmates
chatmate list --installed
chatmate uninstall "Rarely Used Chatmate"
```

### For Team Leads

**Team Management:**
```bash
# Install comprehensive set for team leadership
chatmate hire "Code Review" "Create PR" "Review PR" "Create Issue"

# Weekly team support workflow:
# Use @Review PR for pull request reviews
# Use @Create Issue for task assignment
# Use @Code Review for architecture decisions
```

### For Project Managers

**Project Oversight:**
```bash
# Install project management focused chatmates
chatmate hire "Create Issue" "Optimize Issues" "Review PR" "Merge PR"

# Daily project management:
# Use @Create Issue for requirements documentation
# Use @Review PR for understanding technical changes
# Use @Optimize Issues for clarifying requirements
```

### For QA Engineers

**Quality Assurance:**
```bash
# Install testing and quality focused chatmates
chatmate hire "Testing" "Code Review" "Solve Issue"

# QA workflow:
# Use @Testing for test case generation
# Use @Code Review for code quality assessment
# Use @Solve Issue for bug investigation
```

## Troubleshooting

### Common Issues

#### "VS Code not detected"

**Problem**: ChatMate can't find your VS Code installation.

**Solutions:**
1. Ensure VS Code is in your PATH:
   ```bash
   code --version  # Should show VS Code version
   ```

2. Check VS Code installation location:
   - **macOS**: `/Applications/Visual Studio Code.app`
   - **Windows**: `C:\Users\{username}\AppData\Local\Programs\Microsoft VS Code`
   - **Linux**: `/usr/share/code` or `/opt/visual-studio-code`

3. Manually set VS Code path in your environment:
   ```bash
   export PATH="/Applications/Visual Studio Code.app/Contents/Resources/app/bin:$PATH"
   ```

#### "Permission denied" errors

**Problem**: ChatMate can't write to the prompts directory.

**Solutions:**
1. Check directory permissions:
   ```bash
   chatmate config  # Shows prompts directory path
   ls -la "$(dirname "$(chatmate config | grep 'Prompts directory')")"
   ```

2. Fix permissions:
   ```bash
   # macOS/Linux
   chmod 755 ~/Library/Application\ Support/Code/User/prompts/

   # Or create the directory if it doesn't exist
   mkdir -p ~/Library/Application\ Support/Code/User/prompts/
   ```

3. Run with appropriate permissions:
   ```bash
   sudo chatmate hire  # Only if necessary
   ```

#### "Chatmates not appearing in Copilot Chat"

**Problem**: Installed chatmates don't show up in VS Code.

**Solutions:**
1. **Restart VS Code completely**
   - Close all VS Code windows
   - Reopen VS Code
   - Try using `@` in Copilot Chat to see available chatmates

2. Verify installation:
   ```bash
   chatmate status
   chatmate list --installed
   ```

3. Check VS Code Copilot Chat extension:
   - Ensure Copilot Chat extension is enabled
   - Check that you're signed into GitHub Copilot
   - Try reloading VS Code window: `Ctrl/Cmd+Shift+P` → "Developer: Reload Window"

4. Manual verification:
   ```bash
   # Check if files exist in prompts directory
   ls -la ~/Library/Application\ Support/Code/User/prompts/*.chatmode.md
   ```

#### "Command not found: chatmate"

**Problem**: The ChatMate CLI isn't accessible from your shell.

**Solutions:**
1. Ensure ChatMate is installed and in your PATH:
   ```bash
   which chatmate  # Should show the path to chatmate
   ```

2. Add to PATH if necessary:
   ```bash
   # Add to your shell profile (.bashrc, .zshrc, etc.)
   export PATH="/path/to/chatmate/directory:$PATH"
   ```

3. Use full path temporarily:
   ```bash
   /full/path/to/chatmate hire
   ```

### Diagnostic Commands

When troubleshooting, run these commands to gather information:

```bash
# System information
chatmate status
chatmate config

# Installation verification
chatmate list --installed
ls -la ~/Library/Application\ Support/Code/User/prompts/

# VS Code verification
code --version
which code

# Environment check
echo $PATH
env | grep -i code
```

### Getting Help

If you're still experiencing issues:

1. **Create a diagnostic report:**
   ```bash
   {
     echo "=== ChatMate Status ==="
     chatmate status
     echo -e "\n=== ChatMate Config ==="
     chatmate config
     echo -e "\n=== Installed Chatmates ==="
     chatmate list --installed
     echo -e "\n=== Environment ==="
     echo "OS: $(uname -s)"
     echo "Shell: $SHELL"
     echo "VS Code: $(code --version | head -1)"
   } > chatmate-diagnostic.txt
   ```

2. **Open an issue**: Include your diagnostic report in a [GitHub issue](https://github.com/jonassiebler/chatmate/issues)

3. **Check existing issues**: Search for similar problems in the issue tracker

## Best Practices

### Chatmate Management

#### **Start Simple**
- Begin with `chatmate hire` to install all chatmates
- Use `chatmate list` to familiarize yourself with available agents
- Remove chatmates you don't use with `chatmate uninstall`

#### **Keep Organized**
```bash
# Weekly maintenance routine
chatmate status                    # Check system health
chatmate list --installed         # Review installed chatmates
chatmate uninstall "Unused One"   # Remove chatmates you don't use
chatmate hire --force             # Update existing chatmates
```

#### **Team Synchronization**
- Use `chatmate list > team-chatmates.txt` to share team configurations
- Document which chatmates are recommended for specific roles
- Keep chatmate selections consistent across team members

### VS Code Integration

#### **Effective Usage in Copilot Chat**
- **Be specific**: "@Solve Issue The login form validation isn't working correctly"
- **Provide context**: Include relevant code, error messages, or requirements
- **Use appropriate chatmate**: Match the task to the specialized agent
- **Follow up**: Ask clarifying questions to get better responses

#### **Workflow Integration**
- Use `@Code Review` before committing changes
- Use `@Testing` when implementing new features
- Use `@Documentation` when creating public APIs
- Use `@Solve Issue` when debugging problems

### Performance Optimization

#### **Selective Installation**
```bash
# Instead of installing everything
chatmate hire

# Install only what you need
chatmate hire "Code Review" "Testing" "Solve Issue"
```

#### **Regular Maintenance**
```bash
# Monthly cleanup routine
chatmate list --installed          # Review current chatmates
chatmate uninstall "Unused Ones"   # Remove unused chatmates
chatmate hire --force             # Update remaining chatmates
```

### Security Considerations

#### **Verify Installations**
```bash
# Always verify after installation
chatmate status
chatmate list --installed

# Check file integrity
ls -la ~/Library/Application\ Support/Code/User/prompts/
```

#### **Keep Updated**
```bash
# Regular update routine
chatmate hire --force  # Update all chatmates
```

#### **Review Chatmate Content**
- Chatmates are open source - review their content
- Understand what each chatmate does before using
- Report any security concerns to the project maintainers

## Advanced Usage

### Custom Chatmate Development

#### **Using Create Chatmode**
1. Use the `@Create Chatmode` agent in VS Code:
   ```
   @Create Chatmode I need a chatmate specialized in database optimization
   ```

2. Follow the guidance to create your custom chatmate

3. Save as `.chatmode.md` file in the `mates/` directory

4. Install with `chatmate hire`

#### **Manual Chatmate Creation**
Create a file like `Custom Agent.chatmode.md`:

```markdown
---
name: Custom Agent
description: Description of your custom agent
author: Your Name
version: 1.0.0
tags: [custom, specialized]
---

# Custom Agent Instructions

[Your specialized prompt content here]
```

### Automation and Scripting

#### **Automated Setup Scripts**
```bash
#!/bin/bash
# team-setup.sh - Automated team chatmate setup

echo "Setting up ChatMate for development team..."

# Install core development chatmates
chatmate hire "Code Review" "Testing" "Solve Issue" "Documentation"

# Verify installation
if chatmate status | grep -q "✅"; then
    echo "✅ ChatMate setup complete!"
    echo "Restart VS Code to use your new chatmates"
else
    echo "❌ Setup failed. Check chatmate status"
    exit 1
fi
```

#### **CI/CD Integration**
```yaml
# .github/workflows/chatmate-update.yml
name: Update ChatMates
on:
  schedule:
    - cron: '0 9 * * 1'  # Weekly on Mondays

jobs:
  update-chatmates:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Update ChatMates
        run: |
          chatmate hire --force
          chatmate status
```

### Environment-Specific Configurations

#### **Development Environment**
```bash
# Install comprehensive development set
chatmate hire "Solve Issue" "Code Review" "Testing" "Documentation"
```

#### **Production Support**
```bash
# Install focused production support set
chatmate hire "Solve Issue" "Code Review" "Documentation"
```

#### **Project Management**
```bash
# Install project management focused set
chatmate hire "Create Issue" "Review PR" "Create PR" "Optimize Issues"
```

## FAQ

### General Questions

**Q: What's the difference between chatmates and regular Copilot Chat?**
A: Chatmates are specialized agents with focused expertise. Regular Copilot Chat is general-purpose, while chatmates like "@Solve Issue" are specifically trained for debugging, "@Code Review" for code analysis, etc.

**Q: Do chatmates work offline?**
A: No, chatmates require VS Code Copilot Chat, which needs an internet connection. However, once installed, the chatmate prompts are stored locally.

**Q: Can I use multiple chatmates in the same conversation?**
A: Yes! You can switch between chatmates in the same VS Code Copilot Chat session by using different `@` mentions.

### Installation Questions

**Q: Do I need to reinstall chatmates after VS Code updates?**
A: Usually no, but if you experience issues after VS Code updates, try `chatmate hire --force` to refresh the installation.

**Q: Can I install chatmates on multiple machines?**
A: Yes, run `chatmate hire` on each machine where you use VS Code.

**Q: What happens if I reinstall VS Code?**
A: You'll need to run `chatmate hire` again, as VS Code reinstallation typically clears the user prompts directory.

### Usage Questions

**Q: How do I know which chatmate to use for my task?**
A: Use `chatmate list` to see descriptions, or refer to the [Chatmate Catalog](#chatmate-catalog) section for detailed guidance.

**Q: Can I customize existing chatmates?**
A: You can create custom versions using `@Create Chatmode` or by manually editing the `.chatmode.md` files in your prompts directory.

**Q: Why aren't my chatmates appearing in VS Code?**
A: The most common solution is restarting VS Code after installation. Also verify with `chatmate status` that everything is installed correctly.

### Technical Questions

**Q: Where are chatmates stored?**
A: In your VS Code user prompts directory:
- **macOS**: `~/Library/Application Support/Code/User/prompts`
- **Linux**: `~/.config/Code/User/prompts`
- **Windows**: `%APPDATA%/Code/User/prompts`

**Q: How much disk space do chatmates use?**
A: Very little - typically under 1MB total for all chatmates, as they're just text-based prompt files.

**Q: Can I backup my chatmate configuration?**
A: Yes, backup your VS Code prompts directory or use `chatmate list --installed > my-chatmates.txt` to record your setup.

### Troubleshooting Questions

**Q: What should I do if `chatmate status` shows errors?**
A: Check the [Troubleshooting](#troubleshooting) section for specific error solutions, or run the diagnostic commands to gather information.

**Q: Can I run ChatMate without admin/sudo permissions?**
A: Yes, ChatMate should work with standard user permissions. If you get permission errors, check the troubleshooting section.

**Q: What if VS Code is installed in a non-standard location?**
A: ChatMate should auto-detect most installations. If not, ensure VS Code is in your PATH or check the troubleshooting section.

---

## Need More Help?

- 📖 **Command Help**: Run `chatmate --help` or `chatmate [command] --help`
- 🐛 **Report Issues**: [GitHub Issues](https://github.com/jonassiebler/chatmate/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/jonassiebler/chatmate/discussions)
- 📚 **Documentation**: This guide and the README
- 🔧 **Man Pages**: `man chatmate` (if installed)

Happy chatmate-ing! 🤖✨
