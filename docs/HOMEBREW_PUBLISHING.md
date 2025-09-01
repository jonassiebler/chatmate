# Homebrew Tap Publishing Guide for ChatMate

This document describes how to update and publish the Homebrew tap for the ChatMate CLI.

## Prerequisites
- You must have push access to the `jonassiebler/chatmate` repository
- Homebrew and Go must be installed locally

## Current Setup
ChatMate uses a **single repository approach** with the Formula located at `Formula/chatmate.rb` in the main repository. Users install via:

```bash
brew tap jonassiebler/chatmate https://github.com/jonassiebler/chatmate.git
brew install chatmate
```

## Steps to Publish a New Release to Homebrew

### 1. Create a New Release
```bash
# Create and push a new GitHub release with binaries
gh release create v1.x.x --generate-notes --title "ChatMate v1.x.x"
```

### 2. Update the Homebrew Formula
Edit `Formula/chatmate.rb` and update:
- `url` field to point to the new release tarball
- `sha256` field with the new tarball's checksum

Example:
```ruby
class Chatmate < Formula
  desc "CLI tool for managing AI-powered VS Code Copilot Chat agents"
  homepage "https://github.com/jonassiebler/chatmate"
  url "https://github.com/jonassiebler/chatmate/archive/refs/tags/v1.x.x.tar.gz"
  sha256 "your-new-sha256-hash-here"
  license "MIT"
  # ... rest of formula
end
```

### 3. Calculate SHA256 for New Release
```bash
# Download the release tarball and calculate its SHA256
curl -L https://github.com/jonassiebler/chatmate/archive/refs/tags/v1.x.x.tar.gz | sha256sum
```

### 4. Test the Formula Locally
```bash
# Test installation from local formula
brew install --build-from-source Formula/chatmate.rb

# Verify installation
chatmate --version
```

### 5. Commit and Push Formula Updates
```bash
git add Formula/chatmate.rb
git commit -m "chore: update Homebrew formula for v1.x.x"
git push origin main
```

## Testing User Installation
After publishing, users can install ChatMate via:

```bash
# Add the tap with full URL (required for single repository approach)
brew tap jonassiebler/chatmate https://github.com/jonassiebler/chatmate.git

# Install ChatMate (first time)
brew install chatmate

# Upgrade ChatMate (for existing installations)
brew update && brew upgrade chatmate

# Verify version
chatmate --version
```

**Note:** Users with existing installations can simply run `brew upgrade chatmate` to get the latest version. No need to uninstall/reinstall.

## Troubleshooting
- **"Repository not found" error**: This happens when using the short form `brew tap jonassiebler/chatmate`. Use the full URL form instead.
- **Formula not found**: Ensure `Formula/chatmate.rb` exists in the repository root
- **Build failures**: Check that the SHA256 matches the actual tarball and the Go build process works
- **Version mismatches**: Ensure the GitHub release tag matches the URL in the formula

## Future Considerations
If ChatMate gains significant popularity, consider:
1. Submitting to [homebrew-core](https://github.com/Homebrew/homebrew-core) for better discoverability
2. Setting up automated formula updates via GitHub Actions
3. Using [goreleaser](https://goreleaser.com/) for automated releases and homebrew integration

## Automation (Optional)
- Consider adding a GitHub Action to update the formula automatically on new releases.
- See the official Homebrew documentation for [creating and maintaining taps](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap).

---

For questions, open an issue or see the main README for support links.
