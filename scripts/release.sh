#!/bin/bash

# Release Script for Chatmate
# This script automates the release process including version bumping,
# changelog updates, building, and tagging

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
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CHANGELOG_FILE="$PROJECT_ROOT/CHANGELOG.md"
BUILD_SCRIPT="$PROJECT_ROOT/scripts/build.sh"

# Release types
RELEASE_TYPES=("major" "minor" "patch" "rc" "beta" "alpha")

# Functions
show_usage() {
    echo "Usage: $0 [options] <version|release-type>"
    echo ""
    echo "Automate the release process for chatmate"
    echo ""
    echo "Options:"
    echo "  -d, --dry-run           Show what would be done without making changes"
    echo "  -s, --skip-tests        Skip running tests before release"
    echo "  -S, --skip-build        Skip building binaries"
    echo "  -t, --tag-only          Only create git tag (no build/publish)"
    echo "  -p, --push              Push changes and tags to remote"
    echo "  -h, --help              Show this help"
    echo ""
    echo "Arguments:"
    echo "  version                 Explicit version (e.g., v1.2.3, 1.2.3)"
    echo "  release-type            Semantic version bump (major|minor|patch|rc|beta|alpha)"
    echo ""
    echo "Examples:"
    echo "  $0 v1.2.3               # Release version 1.2.3"
    echo "  $0 patch                # Bump patch version (1.2.3 -> 1.2.4)"
    echo "  $0 minor                # Bump minor version (1.2.3 -> 1.3.0)"
    echo "  $0 major                # Bump major version (1.2.3 -> 2.0.0)"
    echo "  $0 --dry-run patch      # Show what patch release would do"
    echo "  $0 --push v2.0.0        # Release and push to remote"
}

# Get current version from git tags
get_current_version() {
    local version=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
    echo "$version"
}

# Validate version format
validate_version() {
    local version=$1
    
    # Add 'v' prefix if missing
    if [[ ! "$version" =~ ^v ]]; then
        version="v$version"
    fi
    
    # Validate semantic version format
    if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+(\.[0-9]+)?)?$ ]]; then
        log_error "Invalid version format: $version"
        log_info "Expected format: v1.2.3 or v1.2.3-rc.1"
        return 1
    fi
    
    echo "$version"
}

# Calculate next version based on release type
calculate_next_version() {
    local current_version=$1
    local release_type=$2
    
    # Remove 'v' prefix for calculation
    local version_num=${current_version#v}
    
    # Parse version components
    local major minor patch prerelease
    if [[ "$version_num" =~ ^([0-9]+)\.([0-9]+)\.([0-9]+)(-.*)?$ ]]; then
        major=${BASH_REMATCH[1]}
        minor=${BASH_REMATCH[2]}
        patch=${BASH_REMATCH[3]}
        prerelease=${BASH_REMATCH[4]}
    else
        log_error "Failed to parse current version: $current_version"
        return 1
    fi
    
    local new_version
    case "$release_type" in
        "major")
            new_version="v$((major + 1)).0.0"
            ;;
        "minor")
            new_version="v${major}.$((minor + 1)).0"
            ;;
        "patch")
            new_version="v${major}.${minor}.$((patch + 1))"
            ;;
        "rc")
            if [[ -n "$prerelease" ]]; then
                log_error "Cannot create RC from pre-release version"
                return 1
            fi
            new_version="v${major}.${minor}.$((patch + 1))-rc.1"
            ;;
        "beta")
            if [[ -n "$prerelease" ]]; then
                log_error "Cannot create beta from pre-release version"
                return 1
            fi
            new_version="v${major}.${minor}.$((patch + 1))-beta.1"
            ;;
        "alpha")
            if [[ -n "$prerelease" ]]; then
                log_error "Cannot create alpha from pre-release version"
                return 1
            fi
            new_version="v${major}.${minor}.$((patch + 1))-alpha.1"
            ;;
        *)
            log_error "Unknown release type: $release_type"
            return 1
            ;;
    esac
    
    echo "$new_version"
}

# Check if working directory is clean
check_working_directory() {
    if [[ -n "$(git status --porcelain)" ]]; then
        log_error "Working directory is not clean"
        log_info "Please commit or stash your changes before releasing"
        git status --short
        return 1
    fi
    
    log_success "Working directory is clean"
}

# Run tests
run_tests() {
    log_info "🧪 Running test suite"
    
    if [[ -f "$PROJECT_ROOT/run-tests.sh" ]]; then
        cd "$PROJECT_ROOT"
        ./run-tests.sh
    else
        go test ./...
    fi
    
    log_success "All tests passed"
}

# Update changelog
update_changelog() {
    local version=$1
    local date=$(date +"%Y-%m-%d")
    
    log_info "📝 Updating changelog for $version"
    
    if [[ ! -f "$CHANGELOG_FILE" ]]; then
        log_warning "Changelog not found, creating new one"
        cat > "$CHANGELOG_FILE" << EOF
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [${version#v}] - $date
EOF
    else
        # Update unreleased section to new version
        sed -i.bak "s/## \[Unreleased\]/## [Unreleased]\n\n## [${version#v}] - $date/" "$CHANGELOG_FILE"
        rm "$CHANGELOG_FILE.bak"
    fi
    
    log_success "Changelog updated"
}

# Build binaries
build_binaries() {
    local version=$1
    
    log_info "🏗️  Building binaries for $version"
    
    if [[ -f "$BUILD_SCRIPT" ]]; then
        cd "$PROJECT_ROOT"
        VERSION="$version" "$BUILD_SCRIPT" --platforms --release
    else
        log_error "Build script not found: $BUILD_SCRIPT"
        return 1
    fi
    
    log_success "Binaries built successfully"
}

# Create git commit and tag
create_release_commit() {
    local version=$1
    
    log_info "📋 Creating release commit and tag"
    
    # Add changelog to commit
    git add "$CHANGELOG_FILE"
    
    # Create release commit
    git commit -m "Release $version

- Update changelog for $version
- Prepare for release"
    
    # Create annotated tag
    git tag -a "$version" -m "Release $version

See CHANGELOG.md for details."
    
    log_success "Created release commit and tag: $version"
}

# Push to remote
push_release() {
    local version=$1
    
    log_info "🚀 Pushing release to remote"
    
    # Push commits
    git push origin HEAD
    
    # Push tags
    git push origin "$version"
    
    log_success "Release pushed to remote"
}

# Show release summary
show_release_summary() {
    local version=$1
    
    echo ""
    log_success "🎉 Release $version completed successfully!"
    echo ""
    echo "📋 Release Summary:"
    echo "━━━━━━━━━━━━━━━━━━━━"
    echo "Version:      $version"
    echo "Date:         $(date)"
    echo "Commit:       $(git rev-parse --short HEAD)"
    echo "Tag:          $(git describe --tags)"
    echo ""
    
    if [[ -d "$PROJECT_ROOT/dist" ]]; then
        echo "📦 Built Artifacts:"
        ls -la "$PROJECT_ROOT/dist/"
        echo ""
    fi
    
    echo "🔗 Next Steps:"
    echo "• Create GitHub release with release notes"
    echo "• Upload distribution packages"
    echo "• Update documentation"
    echo "• Announce release"
    echo ""
}

# Parse command line arguments
DRY_RUN=false
SKIP_TESTS=false
SKIP_BUILD=false
TAG_ONLY=false
PUSH=false
VERSION_ARG=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -s|--skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        -S|--skip-build)
            SKIP_BUILD=true
            shift
            ;;
        -t|--tag-only)
            TAG_ONLY=true
            SKIP_BUILD=true
            shift
            ;;
        -p|--push)
            PUSH=true
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        -*)
            log_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
        *)
            if [[ -z "$VERSION_ARG" ]]; then
                VERSION_ARG="$1"
            else
                log_error "Too many arguments"
                show_usage
                exit 1
            fi
            shift
            ;;
    esac
done

# Main release process
main() {
    if [[ -z "$VERSION_ARG" ]]; then
        log_error "Version or release type is required"
        show_usage
        exit 1
    fi
    
    # Change to project root
    cd "$PROJECT_ROOT"
    
    log_info "🚀 Starting release process for chatmate"
    
    # Check prerequisites
    if ! command -v git &> /dev/null; then
        log_error "Git is required but not installed"
        exit 1
    fi
    
    if ! command -v go &> /dev/null; then
        log_error "Go is required but not installed"
        exit 1
    fi
    
    # Get current version
    local current_version=$(get_current_version)
    log_info "Current version: $current_version"
    
    # Determine new version
    local new_version
    if [[ " ${RELEASE_TYPES[*]} " =~ " ${VERSION_ARG} " ]]; then
        # Calculate version from release type
        new_version=$(calculate_next_version "$current_version" "$VERSION_ARG")
    else
        # Use explicit version
        new_version=$(validate_version "$VERSION_ARG")
    fi
    
    if [[ -z "$new_version" ]]; then
        exit 1
    fi
    
    log_info "New version: $new_version"
    
    # Dry run mode
    if [[ "$DRY_RUN" == "true" ]]; then
        log_warning "🔍 DRY RUN MODE - No changes will be made"
        echo ""
        echo "Would perform the following actions:"
        echo "• Current version: $current_version"
        echo "• New version: $new_version"
        echo "• Update changelog"
        if [[ "$SKIP_TESTS" != "true" ]]; then
            echo "• Run tests"
        fi
        if [[ "$SKIP_BUILD" != "true" ]]; then
            echo "• Build binaries for all platforms"
        fi
        echo "• Create git commit and tag"
        if [[ "$PUSH" == "true" ]]; then
            echo "• Push to remote repository"
        fi
        exit 0
    fi
    
    # Check working directory
    check_working_directory
    
    # Run tests
    if [[ "$SKIP_TESTS" != "true" ]]; then
        run_tests
    fi
    
    # Update changelog
    update_changelog "$new_version"
    
    # Build binaries
    if [[ "$SKIP_BUILD" != "true" ]]; then
        build_binaries "$new_version"
    fi
    
    # Create release commit and tag
    create_release_commit "$new_version"
    
    # Push to remote
    if [[ "$PUSH" == "true" ]]; then
        push_release "$new_version"
    fi
    
    # Show summary
    show_release_summary "$new_version"
}

# Run main function
main "$@"
