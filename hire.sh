#!/bin/bash

# hire.sh - Install chatmate agents to VS Code Copilot Chat
# This script copies chatmate markdown files to the local prompts folder

set -e

# Get script directory and set paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MATES_DIR="$SCRIPT_DIR/mates"

# Find VS Code Copilot Chat prompts directory based on OS
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    PROMPTS_DIR="$HOME/Library/Application Support/Code/User/prompts"
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    PROMPTS_DIR="$HOME/.config/Code/User/prompts"
elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
    # Windows (Git Bash/Cygwin)
    PROMPTS_DIR="$APPDATA/Code/User/prompts"
else
    echo "Error: Unsupported operating system: $OSTYPE"
    exit 1
fi

# Create prompts directory if it doesn't exist
mkdir -p "$PROMPTS_DIR"

# Copy all markdown files from mates to prompts, always overwriting
echo "Installing chatmates to: $PROMPTS_DIR"
cp -vf "$MATES_DIR"/*.md "$PROMPTS_DIR"/

echo "✅ All chatmates installed! Restart VS Code to use them."
