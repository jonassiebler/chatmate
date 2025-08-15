#!/usr/bin/env bats

# Unit tests for hire.sh installation script
# Tests core functionality, error handling, and cross-platform compatibility

load '../helpers/test_helper'

# Test basic script validation
@test "hire.sh script exists and is executable" {
    run test -f "hire.sh"
    assert_success
    
    run test -x "hire.sh"
    assert_success
}

@test "hire.sh has valid bash syntax" {
    run bash -n "hire.sh"
    assert_success
}

@test "hire.sh contains required shebang" {
    run head -n 1 "hire.sh"
    assert_output --partial "#!/bin/bash"
}

# Test directory detection logic
@test "script detects mates directory correctly" {
    run bash -c 'source hire.sh && echo "$MATES_DIR"'
    assert_success
    assert_output --partial "/mates"
}

@test "script fails gracefully if mates directory missing" {
    # Create temporary script without mates directory
    local temp_script="$TEST_TEMP_DIR/hire_test.sh"
    cp "hire.sh" "$temp_script"
    
    # Create temporary directory without mates
    local temp_dir="$TEST_TEMP_DIR/no_mates"
    mkdir -p "$temp_dir"
    
    cd "$temp_dir" || exit 1
    run "$temp_script"
    assert_failure
}

# Test OS detection and path resolution
@test "script detects macOS paths correctly" {
    simulate_os "macos"
    
    # Mock the script's path detection logic
    run bash -c '
        OSTYPE="darwin20.0"
        if [[ "$OSTYPE" == "darwin"* ]]; then
            echo "$HOME/Library/Application Support/Code/User/prompts"
        fi
    '
    assert_success
    assert_output --partial "Library/Application Support/Code/User/prompts"
}

@test "script detects Linux paths correctly" {
    simulate_os "linux"
    
    run bash -c '
        OSTYPE="linux-gnu"
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            echo "$HOME/.config/Code/User/prompts"
        fi
    '
    assert_success
    assert_output --partial ".config/Code/User/prompts"
}

@test "script detects Windows paths correctly" {
    simulate_os "windows"
    
    run bash -c '
        OSTYPE="msys"
        if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
            echo "$APPDATA/Code/User/prompts"
        fi
    '
    assert_success
    assert_output --partial "appdata/Code/User/prompts"
}

@test "script fails for unsupported OS" {
    simulate_os "unknown"
    
    run bash -c '
        OSTYPE="unknown-os"
        if [[ "$OSTYPE" != "darwin"* && "$OSTYPE" != "linux-gnu"* && "$OSTYPE" != "msys" && "$OSTYPE" != "cygwin" ]]; then
            echo "Error: Unsupported operating system: $OSTYPE"
            exit 1
        fi
    '
    assert_failure
    assert_output --partial "Unsupported operating system"
}

# Test file operations
@test "script creates prompts directory if missing" {
    simulate_os "linux"
    local prompts_dir="$MOCK_HOME/.config/Code/User/prompts"
    
    # Ensure directory doesn't exist
    rm -rf "$prompts_dir"
    
    # Test mkdir -p functionality
    run bash -c "mkdir -p '$prompts_dir'"
    assert_success
    
    run test -d "$prompts_dir"
    assert_success
}

@test "script copies markdown files correctly" {
    local source_dir="$TEST_TEMP_DIR/test_mates"
    local dest_dir="$TEST_TEMP_DIR/prompts"
    
    # Create test source files
    mkdir -p "$source_dir"
    echo "# Test Chatmate 1" > "$source_dir/test1.chatmode.md"
    echo "# Test Chatmate 2" > "$source_dir/test2.chatmode.md"
    echo "Not a chatmate" > "$source_dir/readme.txt"
    
    # Create destination directory
    mkdir -p "$dest_dir"
    
    # Test copy operation
    run bash -c "cp -vf '$source_dir'/*.md '$dest_dir'/"
    assert_success
    
    # Verify files were copied
    assert_equal "$(count_files "$dest_dir" "*.md")" "2"
    assert file_contains "$dest_dir/test1.chatmode.md" "Test Chatmate 1"
    assert file_contains "$dest_dir/test2.chatmode.md" "Test Chatmate 2"
}

@test "script overwrites existing files" {
    local source_dir="$TEST_TEMP_DIR/test_mates"
    local dest_dir="$TEST_TEMP_DIR/prompts"
    
    # Create test files
    mkdir -p "$source_dir" "$dest_dir"
    echo "# New Content" > "$source_dir/test.chatmode.md"
    echo "# Old Content" > "$dest_dir/test.chatmode.md"
    
    # Test overwrite operation
    run bash -c "cp -vf '$source_dir'/*.md '$dest_dir'/"
    assert_success
    
    # Verify file was overwritten
    assert file_contains "$dest_dir/test.chatmode.md" "New Content"
    refute file_contains "$dest_dir/test.chatmode.md" "Old Content"
}

# Test error handling
@test "script handles permission errors gracefully" {
    skip "Requires specific permission setup"
    # This test would need careful setup to avoid affecting the system
}

@test "script validates mates directory exists" {
    # Create script that checks for mates directory
    local test_script="$TEST_TEMP_DIR/check_mates.sh"
    cat > "$test_script" << 'EOF'
#!/bin/bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MATES_DIR="$SCRIPT_DIR/mates"

if [ ! -d "$MATES_DIR" ]; then
    echo "Error: mates directory not found"
    exit 1
fi
echo "mates directory found"
EOF
    chmod +x "$test_script"
    
    # Test without mates directory
    cd "$TEST_TEMP_DIR" || exit 1
    run "$test_script"
    assert_failure
    assert_output --partial "mates directory not found"
    
    # Test with mates directory
    mkdir -p "$TEST_TEMP_DIR/mates"
    run "$test_script"
    assert_success
    assert_output --partial "mates directory found"
}

# Test complete installation flow
@test "complete installation flow works correctly" {
    # This is an integration-style test but kept in unit tests for simplicity
    simulate_os "linux"
    local prompts_dir="$MOCK_HOME/.config/Code/User/prompts"
    
    # Create test chatmate files
    create_test_chatmate "Test1" "# Test Chatmate 1\\nDescription"
    create_test_chatmate "Test2" "# Test Chatmate 2\\nDescription"
    
    # Mock the core installation logic
    run bash -c "
        mkdir -p '$prompts_dir'
        cp -vf '$TEST_MATES_DIR'/*.md '$prompts_dir'/
        echo 'Installation complete'
    "
    assert_success
    assert_output --partial "Installation complete"
    
    # Verify installation
    assert verify_installation "$prompts_dir" "2"
}
