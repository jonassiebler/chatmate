#!/bin/bash

# Build optimization script for chatmate CLI
# This script creates different optimized builds with various flags

set -e

echo "🚀 Building optimized chatmate binaries..."

VERSION=${VERSION:-dev}
COMMIT=${COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo "none")}
BUILD_TIME=${BUILD_TIME:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}

LDFLAGS_BASE="-X main.version=$VERSION -X main.commit=$COMMIT -X main.buildTime=$BUILD_TIME"

# Create builds directory
mkdir -p builds

echo "📊 Current binary size:"
if [ -f chatmate ]; then
    ls -lh chatmate | awk '{print $5, $9}'
fi

echo ""
echo "🔧 Building with different optimization levels..."

# 1. Standard build (for comparison)
echo "Building standard..."
go build -o builds/chatmate-standard .
echo "Standard: $(ls -lh builds/chatmate-standard | awk '{print $5}')"

# 2. Stripped build (remove debug info)
echo "Building stripped..."
go build -ldflags="$LDFLAGS_BASE -s -w" -o builds/chatmate-stripped .
echo "Stripped: $(ls -lh builds/chatmate-stripped | awk '{print $5}')"

# 3. Optimized build (with trimpath)
echo "Building optimized..."
go build -ldflags="$LDFLAGS_BASE -s -w" -trimpath -o builds/chatmate-optimized .
echo "Optimized: $(ls -lh builds/chatmate-optimized | awk '{print $5}')"

# 4. Minimal build (all optimizations)
echo "Building minimal..."
CGO_ENABLED=0 go build -ldflags="$LDFLAGS_BASE -s -w -extldflags '-static'" -trimpath -a -installsuffix cgo -o builds/chatmate-minimal .
echo "Minimal: $(ls -lh builds/chatmate-minimal | awk '{print $5}')"

echo ""
echo "📈 Size comparison:"
echo "Standard:  $(ls -lh builds/chatmate-standard | awk '{print $5}')"
echo "Stripped:  $(ls -lh builds/chatmate-stripped | awk '{print $5}')"
echo "Optimized: $(ls -lh builds/chatmate-optimized | awk '{print $5}')"
echo "Minimal:   $(ls -lh builds/chatmate-minimal | awk '{print $5}')"

echo ""
echo "✅ Testing optimized builds..."

# Test each build
for build in standard stripped optimized minimal; do
    echo -n "Testing $build... "
    if ./builds/chatmate-$build --version >/dev/null 2>&1; then
        echo "✅"
    else
        echo "❌ Failed!"
    fi
done

echo ""
echo "🎯 Recommended build: chatmate-optimized (best size/functionality balance)"
echo "📦 Production build: chatmate-minimal (smallest size, static binary)"

# Copy optimized build as default
cp builds/chatmate-optimized chatmate
echo ""
echo "✅ Default binary updated with optimized build"
