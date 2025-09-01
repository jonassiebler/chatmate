# Homebrew Tap Publishing Guide for ChatMate

This document describes how to update and publish the Homebrew tap for the ChatMate CLI.

## Prerequisites
- You must have push access to the `jonassiebler/chatmate` and `jonassiebler/homebrew-tap` repositories.
- Homebrew and Go must be installed locally.

## Steps to Publish a New Release to Homebrew

1. **Release a new version of ChatMate:**
   - Run the release script to tag and build a new version:
     ```bash
     ./scripts/release.sh <version> --push
     ```
   - This will create a new tag and push it to GitHub.

2. **Update the Homebrew Formula:**
   - Edit `homebrew-tap/Formula/chatmate.rb`:
     - Update the `revision` field to the latest commit hash for the release.
     - Update the `version` field to the new version (e.g., `20250817113552`).
   - Commit and push the formula update:
     ```bash
     git add homebrew-tap/Formula/chatmate.rb
     git commit -m "chore: update Homebrew formula for v<version>"
     git push origin <branch>
     ```

3. **Test the Formula Locally:**
   - Uninstall any existing version:
     ```bash
     brew uninstall chatmate
     ```
   - Install from the updated formula:
     ```bash
     brew install --build-from-source ./homebrew-tap/Formula/chatmate.rb
     ```
   - Verify installation:
     ```bash
     chatmate --help
     ```

4. **Publish to the Homebrew Tap Repository:**
   - If you maintain a separate `homebrew-tap` repo, copy the updated formula there and push.
   - Otherwise, users can tap directly from this repo:
     ```bash
     brew tap jonassiebler/chatmate
     brew install chatmate
     ```

## Automation (Optional)
- Consider adding a GitHub Action to update the formula automatically on new releases.
- See the official Homebrew documentation for [creating and maintaining taps](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap).

## Troubleshooting
- If users report issues, ensure the formula revision and version match the latest release.
- Check build logs for Go or Homebrew errors.

---

For questions, open an issue or see the main README for support links.
