#!/bin/bash

# ChatMate Man Page Installation Script
# This script installs the generated man pages to the system man directory

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
MAN_SOURCE_DIR="$PROJECT_ROOT/docs/man"
MAN_INSTALL_DIR="/usr/local/share/man/man1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}📖 ChatMate Man Page Installer${NC}"
echo ""

# Check if man pages exist
if [ ! -d "$MAN_SOURCE_DIR" ] || [ -z "$(ls -A "$MAN_SOURCE_DIR"/*.1 2>/dev/null)" ]; then
    echo -e "${RED}❌ Error: Man pages not found in $MAN_SOURCE_DIR${NC}"
    echo "Please generate man pages first:"
    echo "  cd $PROJECT_ROOT"
    echo "  go run scripts/generate-man-pages.go docs/man"
    exit 1
fi

# Check if we have permission to install
if [ ! -w "$(dirname "$MAN_INSTALL_DIR")" ] && [ "$EUID" -ne 0 ]; then
    echo -e "${YELLOW}⚠️  This script requires sudo to install man pages system-wide.${NC}"
    echo "Run with: sudo $0"
    echo ""
    echo "Or install locally to your user directory:"
    echo "  mkdir -p ~/.local/share/man/man1"
    echo "  cp $MAN_SOURCE_DIR/*.1 ~/.local/share/man/man1/"
    echo "  export MANPATH=~/.local/share/man:\$MANPATH"
    exit 1
fi

# Create install directory if it doesn't exist
echo -e "${BLUE}📁 Ensuring man directory exists...${NC}"
if [ ! -d "$MAN_INSTALL_DIR" ]; then
    mkdir -p "$MAN_INSTALL_DIR"
    echo -e "${GREEN}✅ Created $MAN_INSTALL_DIR${NC}"
fi

# Install man pages
echo -e "${BLUE}📖 Installing man pages...${NC}"
for manpage in "$MAN_SOURCE_DIR"/*.1; do
    if [ -f "$manpage" ]; then
        filename=$(basename "$manpage")
        echo "  Installing $filename"
        cp "$manpage" "$MAN_INSTALL_DIR/"
        chmod 644 "$MAN_INSTALL_DIR/$filename"
    fi
done

# Update man database (if mandb exists)
if command -v mandb >/dev/null 2>&1; then
    echo -e "${BLUE}🔄 Updating man database...${NC}"
    mandb -q 2>/dev/null || echo -e "${YELLOW}⚠️  Warning: Could not update man database${NC}"
else
    echo -e "${YELLOW}⚠️  mandb not found, skipping database update${NC}"
fi

echo ""
echo -e "${GREEN}✅ Man pages installed successfully!${NC}"
echo ""
echo -e "${BLUE}📖 Available man pages:${NC}"
for manpage in "$MAN_SOURCE_DIR"/*.1; do
    if [ -f "$manpage" ]; then
        filename=$(basename "$manpage" .1)
        echo "  man $filename"
    fi
done

echo ""
echo -e "${BLUE}🔍 Test installation:${NC}"
echo "  man chatmate"
echo "  man chatmate-hire"
echo "  man chatmate-list"
echo ""
echo -e "${GREEN}🎉 Installation complete!${NC}"
