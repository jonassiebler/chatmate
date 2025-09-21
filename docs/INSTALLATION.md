# ChatMate Installation Guide 🚀

This guide covers installing ChatMate CLI and setting up chatmates for VS Code Copilot Chat.

## Prerequisites

Before installing ChatMate, ensure you have:

- **VS Code** installed and accessible from your PATH
- **VS Code Copilot Chat extension** enabled
- **Write permissions** to your VS Code user directory
- **Go 1.19+** (if building from source)

## Quick Installation

The fastest way to get started:

```bash
# Install or download the ChatMate CLI
# Then install all chatmates
chatmate hire
```

## Installation Verification

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

For more detailed usage information, see the [User Guide](USER_GUIDE.md).
