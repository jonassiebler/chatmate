# ChatMate Installation Guide 🚀

Quick installation guide for getting ChatMate up and running with VS Code Copilot Chat.

## Prerequisites ✅

Before installing ChatMate, ensure you have:

- ✅ **VS Code** installed and accessible
- ✅ **GitHub Copilot Chat extension** enabled in VS Code
- ✅ **Active GitHub Copilot subscription**
- ✅ **Write permissions** to your VS Code user directory

## Quick Installation (Recommended)

### 1. Get ChatMate CLI

#### Option A: Download Release Binary

```bash
# Download the latest release for your platform
# Extract and move to your PATH
```

#### Option B: Build from Source

```bash
git clone https://github.com/jonassiebler/chatmate.git
cd chatmate
go build -o chatmate .
```

### 2. Install All Chatmates

```bash
# Install all available chatmates (recommended for new users)
./chatmate hire
```

### 3. Restart VS Code

- Close all VS Code windows
- Reopen VS Code
- Open Copilot Chat (`Ctrl/Cmd+Shift+P` → "Chat: Open Chat")

### 4. Test Installation

In VS Code Copilot Chat, try:
```
@Solve Issue Hello! Can you help me debug code?
```

If you see a response from the Solve Issue chatmate, you're all set! 🎉

## Alternative Installation Methods

### Legacy Script Installation

```bash
git clone https://github.com/jonassiebler/chatmate.git
cd chatmate
chmod +x hire.sh
./hire.sh
```

### Manual Installation

1. **Find your VS Code prompts directory:**

   **macOS:**
   ```bash
   ~/Library/Application Support/Code/User/prompts
   ```

   **Linux:**
   ```bash
   ~/.config/Code/User/prompts
   ```

   **Windows:**
   ```bash
   %APPDATA%\Code\User\prompts
   ```

2. **Copy chatmate files:**
   ```bash
   cp mates/*.chatmode.md "/path/to/your/prompts/directory/"
   ```

## Verification Commands

After installation, verify everything is working:

```bash
# Check system status
./chatmate status

# List installed chatmates
./chatmate list --installed

# View configuration
./chatmate config
```

Expected output should show:
- ✅ VS Code detected
- ✅ Prompts directory exists
- ✅ Chatmates installed

## Troubleshooting Quick Fixes

### "VS Code not detected"
```bash
# Ensure VS Code is in your PATH
code --version

# If not found, add VS Code to PATH (macOS example)
export PATH="/Applications/Visual Studio Code.app/Contents/Resources/app/bin:$PATH"
```

### "Permission denied"
```bash
# Check prompts directory permissions
ls -la ~/Library/Application\ Support/Code/User/

# Create directory if missing
mkdir -p ~/Library/Application\ Support/Code/User/prompts

# Fix permissions if needed
chmod 755 ~/Library/Application\ Support/Code/User/prompts
```

### "Chatmates not appearing in VS Code"
1. **Completely restart VS Code** (close all windows, reopen)
2. Verify Copilot Chat extension is enabled
3. Check you're signed into GitHub Copilot
4. Try `Ctrl/Cmd+Shift+P` → "Developer: Reload Window"

## What Gets Installed?

ChatMate installs these specialized AI agents to your VS Code:

| Agent | Purpose | Example Usage |
|-------|---------|---------------|
| Solve Issue | Debugging & problem solving | `@Solve Issue` |
| Code Review | Code analysis & quality | `@Code Review` |
| Testing | Test generation & debugging | `@Testing` |
| Create PR | Pull request creation | `@Create PR` |
| Documentation | Technical writing | `@Documentation` |
| Create Issue | GitHub issue creation | `@Create Issue` |

## Next Steps

1. **Explore chatmates**: Run `./chatmate list` to see all available agents
2. **Read the User Guide**: [docs/USER_GUIDE.md](USER_GUIDE.md) for comprehensive usage
3. **Customize your setup**: Install only the chatmates you need
4. **Stay updated**: Use `./chatmate hire --force` to update chatmates

## Need Help?

- 📖 **Full Documentation**: [User Guide](USER_GUIDE.md)
- 🐛 **Issues**: [GitHub Issues](https://github.com/jonassiebler/chatmate/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/jonassiebler/chatmate/discussions)
- 🔧 **Command Help**: `./chatmate --help`

---

**Installation complete?** Start using your chatmates:
```
@Solve Issue How do I optimize this React component?
@Code Review Please review this authentication function
@Testing Generate unit tests for this service class
```

Happy coding! 🤖✨
