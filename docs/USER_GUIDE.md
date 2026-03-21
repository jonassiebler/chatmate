# ChatMate User Guide 📖

Welcome to ChatMate! This guide covers the commands and chatmates available for VS Code Copilot Chat.

## Quick Start

For installation and setup, see the [Installation Guide](INSTALLATION.md).
For troubleshooting, see the [Troubleshooting Guide](TROUBLESHOOTING.md).
For advanced usage, see the [Advanced Usage Guide](ADVANCED_USAGE.md).

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

### `chatmate legacy`

Interact with legacy chatmates that are no longer part of the default hire workflow.

**Subcommands:**
- `chatmate legacy list` — Show which legacy chatmates exist and are installed
- `chatmate legacy hire <names...>` — Manually install specific legacy chatmates (supports `--force`)
- `chatmate legacy fire` — Interactively remove installed legacy chatmates with per-mate confirmation

**Examples:**
```bash
# Inspect available legacy chatmates
chatmate legacy list

# Install a specific legacy chatmate
chatmate legacy hire "Solve Issue"

# Force reinstall legacy prompts after manual edits
chatmate legacy hire --force "Solve Issue"

# Clean up installed legacy chatmates with prompts per mate
chatmate legacy fire
```

Legacy chatmates stay available for specialized workflows, but they do not appear in standard `chatmate hire` or `chatmate list` output.

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

> **💡 Optimized Design**: All chatmates feature streamlined, language-agnostic instructions with 3-Domain Safety Paradigm for Implementation-Testing-Documentation validation, ensuring reliable and efficient development workflows.

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

#### **Review Repo** 🔍
- **Purpose**: Comprehensive repository analysis and health assessment
- **Best for**: Code quality evaluation, architectural review, improvement recommendations
- **Usage**: `@Review Repo` for repository-wide analysis
- **Example**: "@Review Repo Analyze this repository for technical debt and improvement opportunities"

---

## Additional Resources

- 📖 **Installation**: [Installation Guide](INSTALLATION.md)
- 🔧 **Troubleshooting**: [Troubleshooting Guide](TROUBLESHOOTING.md)
- 🚀 **Advanced Usage**: [Advanced Usage Guide](ADVANCED_USAGE.md)
- 📖 **Command Help**: Run `chatmate --help` or `chatmate [command] --help`
- 🐛 **Report Issues**: [GitHub Issues](https://github.com/jonassiebler/chatmate/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/jonassiebler/chatmate/discussions)

Happy chatmate-ing! 🤖✨
