# chatmate 🤖

Open source collection of specialized AI agents for VS Code Copilot Chat. Each chatmate is a carefully crafted prompt designed to excel at specific development tasks.

## What are chatmates?

Chatmates are custom Copilot Chatmodes, configured with specialized agentic skills, stored as markdown files. When installed, they become available as specialized assistants in VS Code, each with their own expertise:

- **Solve Issue** - Automatically analyzes and resolves GitHub issues
- **Code Reviewer** - Provides thorough code reviews with security and quality insights  
- **Create PR** - Generates comprehensive pull requests with detailed descriptions
- **Review PR** - Analyzes pull requests for quality, security, and best practices
- **Create Issue** - Creates well-structured GitHub issues with proper templates
- **Optimize Issues** - Improves existing issues for clarity and actionability
- **Create Chatmode** - Helps you build new custom Copilot Chat agents
- **Code Claude Sonnet 4** & **Code GPT-4.1** - Specialized coding assistants

## Quick Start

1. **Clone the repository:**

   ```bash
   git clone https://github.com/jonassiebler/chatmate.git
   cd chatmate
   ```

2. **Install all chatmates:**

   ```bash
   ./hire.sh
   ```

3. **Restart VS Code** to use your new chatmates

## Repository Structure

```text
chatmate/
├── mates/          # Chatmate markdown files (the AI agents)
├── hire.sh         # Installation script
└── README.md       # This file
```

## Installation

### hire.sh

The installation script that copies all chatmates to your VS Code prompts folder.

```bash
./hire.sh
```

## Manual Installation

If you prefer to install chatmates manually or selectively:

1. **Find your VS Code prompts directory:**
   - **macOS:** `~/Library/Application Support/Code/User/prompts`
   - **Linux:** `~/.config/Code/User/prompts`
   - **Windows:** `%APPDATA%/Code/User/prompts`

2. **Copy specific chatmates:**

   ```bash
   cp mates/"Solve Issue.chatmode.md" "/path/to/prompts/folder/"
   ```

## Uninstalling

To remove chatmates, simply delete them from your VS Code prompts directory:

```bash
rm ~/Library/Application\ Support/Code/User/prompts/*.chatmode.md
```

## Creating Custom Chatmates

1. Use the **Create Chatmode** agent to help build new chatmates
2. Add your custom `.md` files to the `mates/` directory
3. Run `./hire.sh` to install them

## Contributing

We welcome contributions from the community! Whether you're fixing bugs, adding new chatmates, or improving documentation, your help makes chatmate better for everyone.

### Development Workflow

We follow a structured workflow to maintain code quality:

- **`main`** - Production branch (stable releases)
- **`dev`** - Development branch (integration and testing)
- **`feature/*`** - Feature branches (your contributions)

### Quick Contribution Guide

1. **Fork** the repository
2. **Create** a feature branch from `dev`
3. **Make** your changes following our guidelines
4. **Test** thoroughly using `./hire.sh`
5. **Submit** a pull request to `dev` branch

### Contribution Types

- 🐛 **Bug fixes** for existing chatmates
- ✨ **New chatmate** agents
- 📖 **Documentation** improvements
- 🔧 **Installation script** enhancements

For detailed guidelines, see [CONTRIBUTING.md](CONTRIBUTING.md)

### Code of Conduct

- Be respectful and inclusive
- Follow our coding standards
- Test your changes thoroughly
- Provide clear documentation

## Branch Protection & Review Process

### Protected Branches

- **`main`** and **`dev`** branches are protected
- All changes require pull requests
- Admin review is mandatory before merging
- Direct commits restricted to admin group

### Review Requirements

- ✅ Automated tests must pass
- ✅ Admin approval required
- ✅ No merge conflicts
- ✅ Follow contribution guidelines

This ensures high-quality, stable releases for all users.

## Requirements

- VS Code with GitHub Copilot Chat extension
- Git (for cloning the repository)
- Bash shell (for installation scripts)

## Supported Platforms

- ✅ macOS
- ✅ Linux  
- ✅ Windows (Git Bash/WSL)

## Quality Assurance

### Automated Testing

Every contribution is automatically validated through our CI/CD pipeline:

- 🔍 **Structure validation** - Repository and file organization
- 🤖 **Chatmate validation** - Syntax and naming conventions
- 🧪 **Installation testing** - Cross-platform compatibility
- 📝 **Documentation linting** - Markdown formatting and links
- 🔒 **Security scanning** - Sensitive data detection

### Continuous Integration

- **GitHub Actions** automatically test all pull requests
- **Multi-platform testing** on Ubuntu and macOS
- **Automated releases** with comprehensive changelogs
- **Link checking** ensures documentation stays current

## License

MIT License - Feel free to use these chatmates in your own projects!

## Support

- 🐛 **Issues:** [GitHub Issues](https://github.com/jonassiebler/chatmate/issues)
- 💡 **Feature Requests:** [GitHub Discussions](https://github.com/jonassiebler/chatmate/discussions)
- 📧 **Contact:** Create an issue for questions

---

Happy coding with your new chatmates!
