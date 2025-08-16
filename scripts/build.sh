#!/bin/bash

# Build Script for Chatmate
# This script builds the chatmate binary with proper version information

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Configuration
BINARY_NAME="chatmate"
# Note: VERSION can be overridden by command line argument
DEFAULT_VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
VERSION=${VERSION:-$DEFAULT_VERSION}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION=$(go version | cut -d' ' -f3)

# Build flags for optimization
BUILD_FLAGS=(
    "-trimpath"
    "-ldflags=-s -w -X github.com/jonassiebler/chatmate/cmd.version=${VERSION} -X github.com/jonassiebler/chatmate/cmd.commit=${COMMIT} -X github.com/jonassiebler/chatmate/cmd.date=${DATE}"
)

# Platform configurations
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

# Functions
show_usage() {
    echo "Usage: $0 [options]"
    echo ""
    echo "Build chatmate binary with proper version information"
    echo ""
    echo "Options:"
    echo "  -v, --version VERSION    Set version (default: from git tag or 'dev')"
    echo "  -o, --output DIR         Output directory (default: current directory)"
    echo "  -p, --platforms          Build for all platforms"
    echo "  -r, --release            Build optimized release version"
    echo "  -t, --test               Run tests before building"
    echo "  -c, --clean              Clean before building"
    echo "  -h, --help               Show this help"
    echo ""
    echo "Examples:"
    echo "  $0                       # Build for current platform"
    echo "  $0 -r                    # Build optimized release"
    echo "  $0 -p                    # Build for all platforms"
    echo "  $0 -v v1.2.3             # Build with specific version"
    echo "  $0 -t -c -r              # Clean, test, and build release"
}

show_build_info() {
    echo ""
    log_info "🏗️  Build Configuration"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Binary:       $BINARY_NAME"
    echo "Version:      $VERSION"
    echo "Commit:       $COMMIT"
    echo "Date:         $DATE"
    echo "Go Version:   $GO_VERSION"
    echo ""
}

clean_build() {
    log_info "🧹 Cleaning build artifacts"
    rm -f ${BINARY_NAME}*
    rm -rf dist/
    go clean -cache
    log_success "Build artifacts cleaned"
}

run_tests() {
    log_info "🧪 Running tests"
    if [[ -f "./run-tests.sh" ]]; then
        ./run-tests.sh
    else
        go test ./...
    fi
    log_success "All tests passed"
}

build_binary() {
    local os=$1
    local arch=$2
    local output_name=$3
    
    local env_vars="GOOS=$os GOARCH=$arch"
    local output_path="$output_name"
    
    if [[ "$os" == "windows" ]]; then
        output_path="${output_path}.exe"
    fi
    
    log_info "Building for $os/$arch → $output_path"
    
    if ! env $env_vars go build "${BUILD_FLAGS[@]}" -o "$output_path" .; then
        log_error "Failed to build for $os/$arch"
        return 1
    fi
    
    # Show binary info
    if [[ "$os" == $(go env GOOS) ]] && [[ "$arch" == $(go env GOARCH) ]]; then
        local size=$(du -h "$output_path" | cut -f1)
        log_success "Built $output_path ($size)"
    else
        local size=$(du -h "$output_path" | cut -f1)
        log_success "Cross-compiled $output_path ($size)"
    fi
}

build_all_platforms() {
    log_info "🌍 Building for all platforms"
    mkdir -p dist
    
    for platform in "${PLATFORMS[@]}"; do
        IFS='/' read -r os arch <<< "$platform"
        local output_name="dist/${BINARY_NAME}-${os}-${arch}"
        build_binary "$os" "$arch" "$output_name"
    done
    
    log_success "Built binaries for all platforms in dist/"
    log_info "📁 Distribution directory contents:"
    ls -la dist/
}

build_current_platform() {
    local current_os=$(go env GOOS)
    local current_arch=$(go env GOARCH)
    build_binary "$current_os" "$current_arch" "$BINARY_NAME"
}

# Parse command line arguments
CLEAN=false
TEST=false
RELEASE=false
ALL_PLATFORMS=false
OUTPUT_DIR=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -p|--platforms)
            ALL_PLATFORMS=true
            shift
            ;;
        -r|--release)
            RELEASE=true
            shift
            ;;
        -t|--test)
            TEST=true
            shift
            ;;
        -c|--clean)
            CLEAN=true
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Main build process
main() {
    log_info "🚀 Starting Chatmate Build Process"
    
    # Check Go installation
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    # Show build information
    show_build_info
    
    # Clean if requested
    if [[ "$CLEAN" == "true" ]]; then
        clean_build
    fi
    
    # Run tests if requested
    if [[ "$TEST" == "true" ]]; then
        run_tests
    fi
    
    # Add release optimizations
    if [[ "$RELEASE" == "true" ]]; then
        log_info "🎯 Building optimized release version"
        BUILD_FLAGS+=("-tags=release")
    fi
    
    # Change to output directory if specified
    if [[ -n "$OUTPUT_DIR" ]]; then
        mkdir -p "$OUTPUT_DIR"
        cd "$OUTPUT_DIR"
    fi
    
    # Build binaries
    if [[ "$ALL_PLATFORMS" == "true" ]]; then
        build_all_platforms
    else
        build_current_platform
    fi
    
    echo ""
    log_success "🎉 Build completed successfully!"
    
    # Show final binary info for current platform
    if [[ "$ALL_PLATFORMS" != "true" ]]; then
        if [[ -f "$BINARY_NAME" ]]; then
            log_info "Testing binary version:"
            ./$BINARY_NAME version --quiet
        fi
    fi
}

# Run main function
main "$@"
