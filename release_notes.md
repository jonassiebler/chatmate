🎉 **First Official Release of ChatMate!**

ChatMate is now available as a fully-featured CLI tool for managing VS Code Copilot Chat agents.

## 🚀 What's New
- Complete rewrite in Go with comprehensive testing framework
- Native cross-platform support (macOS, Linux, Windows)  
- Homebrew tap integration for easy installation
- Enhanced security validation and path protection
- Performance optimizations and better error handling
- Rich CLI interface with detailed help and configuration

## 📦 Installation

### Homebrew (Recommended)
```bash
brew tap jonassiebler/chatmate https://github.com/jonassiebler/chatmate.git
brew install chatmate
chatmate hire
```

### Manual Installation
Download the appropriate binary for your platform from the assets below, make it executable, and run `chatmate hire`.

## 🤖 Available Chatmates
- **Solve Issue** - Systematic debugging and problem resolution  
- **Code Review** - Expert code analysis and feedback
- **Testing** - Test generation and debugging assistance
- **Create PR** - Pull request creation and optimization
- **Create Issue** - GitHub issue creation with proper formatting
- And many more! Run `chatmate list` to see all available agents.

## 🔄 Upgrading from Previous Version
If you were using the legacy script installation, simply run `chatmate hire` to migrate to the new system seamlessly.

## 📋 Requirements
- VS Code with GitHub Copilot Chat extension
- Write permissions to VS Code user directory

## 🆘 Support
- 📖 [User Guide](https://github.com/jonassiebler/chatmate/blob/main/docs/USER_GUIDE.md)
- 🐛 [Report Issues](https://github.com/jonassiebler/chatmate/issues)
- 💡 [Discussions](https://github.com/jonassiebler/chatmate/discussions)

**Ready to supercharge your VS Code experience?** Install ChatMate and start using `@ChatmateName` in Copilot Chat!
