# ChatMate Troubleshooting Guide 🔧

This guide helps you resolve common ChatMate installation and usage issues.

## Common Issues

### "VS Code not detected"

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

### "Permission denied" errors

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

### "Chatmates not appearing in Copilot Chat"

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

### "Command not found: chatmate"

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

## Diagnostic Commands

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

## Getting Help

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
A: Use `chatmate list` to see descriptions, or refer to the User Guide for detailed guidance.

**Q: Can I customize existing chatmates?**
A: You can create custom versions or manually edit the `.chatmode.md` files in your prompts directory.

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

---

## Need More Help?

- 📖 **Command Help**: Run `chatmate --help` or `chatmate [command] --help`
- 🐛 **Report Issues**: [GitHub Issues](https://github.com/jonassiebler/chatmate/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/jonassiebler/chatmate/discussions)
- 📚 **Documentation**: User Guide and README
- 🔧 **Man Pages**: `man chatmate` (if installed)
