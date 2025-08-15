#!/usr/bin/env bats

# Test configuration and setup for Bats testing framework
# This file provides shared setup and helper functions for all tests

# Load bats libraries using relative paths
load '../bats-support/load'
load '../bats-assert/load'

# Test environment setup
setup() {
    # Create temporary directory for tests
    export TEST_TEMP_DIR="$(mktemp -d)"
    export ORIGINAL_PWD="$PWD"
    
    # Set up test environment variables
    export CHATMATE_TEST_MODE=true
    export TEST_MATES_DIR="$PWD/mates"
    export TEST_SCRIPT_DIR="$PWD"
    
    # Mock common system paths for testing
    export MOCK_HOME="$TEST_TEMP_DIR/home"
    export MOCK_APPDATA="$TEST_TEMP_DIR/appdata"
    
    # Create mock directories
    mkdir -p "$MOCK_HOME/Library/Application Support/Code/User/prompts"
    mkdir -p "$MOCK_HOME/.config/Code/User/prompts"
    mkdir -p "$MOCK_APPDATA/Code/User/prompts"
}

# Test environment cleanup
teardown() {
    # Clean up temporary directory
    if [[ -n "$TEST_TEMP_DIR" && -d "$TEST_TEMP_DIR" ]]; then
        rm -rf "$TEST_TEMP_DIR"
    fi
    
    # Return to original directory
    cd "$ORIGINAL_PWD" || exit 1
}

# Helper function to count files in directory
count_files() {
    local dir="$1"
    local pattern="${2:-*}"
    
    if [[ ! -d "$dir" ]]; then
        echo "0"
        return
    fi
    
    find "$dir" -maxdepth 1 -name "$pattern" -type f | wc -l | tr -d ' '
}

# Helper function to check if file contains specific content
file_contains() {
    local file="$1"
    local content="$2"
    
    if [[ ! -f "$file" ]]; then
        return 1
    fi
    
    grep -q "$content" "$file"
}

# Helper function to validate chatmate file format
is_valid_chatmate() {
    local file="$1"
    
    # Check file extension
    if [[ ! "$file" =~ \.chatmode\.md$ ]]; then
        return 1
    fi
    
    # Check file exists and is not empty
    if [[ ! -s "$file" ]]; then
        return 1
    fi
    
    # Check for basic content structure (headers)
    if ! grep -q "^#" "$file"; then
        return 1
    fi
    
    return 0
}

# Helper function to simulate different OS environments
simulate_os() {
    local os_type="$1"
    
    case "$os_type" in
        "macos"|"darwin")
            export OSTYPE="darwin20.0"
            export HOME="$MOCK_HOME"
            ;;
        "linux")
            export OSTYPE="linux-gnu"
            export HOME="$MOCK_HOME"
            ;;
        "windows")
            export OSTYPE="msys"
            export APPDATA="$MOCK_APPDATA"
            ;;
        *)
            export OSTYPE="unknown"
            ;;
    esac
}

# Helper function to create test chatmate files
create_test_chatmate() {
    local name="$1"
    local content="${2:-# Test Chatmate\n\nThis is a test chatmate file."}"
    local file="$TEST_MATES_DIR/$name.chatmode.md"
    
    mkdir -p "$TEST_MATES_DIR"
    echo -e "$content" > "$file"
    echo "$file"
}

# Helper function to verify installation results
verify_installation() {
    local prompts_dir="$1"
    local expected_count="$2"
    
    # Check directory exists
    if [[ ! -d "$prompts_dir" ]]; then
        return 1
    fi
    
    # Check file count
    local actual_count
    actual_count=$(count_files "$prompts_dir" "*.md")
    
    if [[ "$actual_count" != "$expected_count" ]]; then
        return 1
    fi
    
    return 0
}
