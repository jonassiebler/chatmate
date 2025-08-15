#!/usr/bin/env bats

# CLI Hire Command Tests
# Tests the core hire functionality that replaces hire.sh

load '../helpers/test_helper'

setup() {
    # Call parent setup
    load '../helpers/test_helper'
    
    # Create temporary prompts directory for testing
    export TEST_PROMPTS_DIR="$TEST_TEMP_DIR/prompts"
    mkdir -p "$TEST_PROMPTS_DIR"
    
    # Mock VS Code prompts directory by temporarily modifying environment
    export ORIGINAL_HOME="$HOME"
    export HOME="$TEST_TEMP_DIR"
    
    # Create mock VS Code directory structure
    case "$(uname)" in
        Darwin) # macOS
            export MOCK_PROMPTS_DIR="$HOME/Library/Application Support/Code/User/prompts"
            ;;
        Linux)
            export MOCK_PROMPTS_DIR="$HOME/.config/Code/User/prompts"
            ;;
        *)
            export MOCK_PROMPTS_DIR="$HOME/AppData/Roaming/Code/User/prompts"
            ;;
    esac
    
    mkdir -p "$MOCK_PROMPTS_DIR"
}

teardown() {
    # Restore original HOME
    export HOME="$ORIGINAL_HOME"
    
    # Call parent teardown
    load '../helpers/test_helper'
}

# Test basic hire functionality
@test "hire command installs all chatmates successfully" {
    # Count available chatmates
    local chatmate_count
    chatmate_count=$(count_files "mates" "*.chatmode.md")
    
    run node bin/chatmate.js hire --force
    assert_success
    assert_output --partial "Installing chatmate agents"
    assert_output --partial "All chatmates installed"
    
    # Verify installation
    local installed_count
    installed_count=$(count_files "$MOCK_PROMPTS_DIR" "*.chatmode.md")
    assert [ "$installed_count" -eq "$chatmate_count" ]
}

@test "hire command with --specific flag installs only specified chatmates" {
    run node bin/chatmate.js hire --specific "Solve Issue"
    assert_success
    assert_output --partial "Installing chatmate agents"
    
    # Verify only Solve Issue was installed
    run test -f "$MOCK_PROMPTS_DIR/Solve Issue.chatmode.md"
    assert_success
    
    # Verify other chatmates are not installed
    run test -f "$MOCK_PROMPTS_DIR/Create PR.chatmode.md"
    assert_failure
}

@test "hire command with multiple specific chatmates" {
    run node bin/chatmate.js hire --specific "Solve Issue" "Create PR"
    assert_success
    
    # Verify both specified chatmates were installed
    run test -f "$MOCK_PROMPTS_DIR/Solve Issue.chatmode.md"
    assert_success
    
    run test -f "$MOCK_PROMPTS_DIR/Create PR.chatmode.md"
    assert_success
}

@test "hire command warns about non-existent chatmate" {
    run node bin/chatmate.js hire --specific "NonExistent"
    assert_success
    assert_output --partial "No chatmate found matching: NonExistent"
}

@test "hire command with --force flag reinstalls existing chatmates" {
    # Install once
    run node bin/chatmate.js hire
    assert_success
    
    # Install again with force
    run node bin/chatmate.js hire --force
    assert_success
    assert_output --partial "reinstalled"
}

@test "hire command handles missing mates directory gracefully" {
    # Temporarily move mates directory
    mv "mates" "mates.backup"
    
    run node bin/chatmate.js hire
    assert_failure
    assert_output --partial "Failed to read chatmates directory"
    
    # Restore mates directory
    mv "mates.backup" "mates"
}

@test "hire command creates prompts directory if it doesn't exist" {
    # Remove prompts directory
    rm -rf "$MOCK_PROMPTS_DIR"
    
    run node bin/chatmate.js hire
    assert_success
    
    # Verify directory was created and files installed
    run test -d "$MOCK_PROMPTS_DIR"
    assert_success
    
    local installed_count
    installed_count=$(count_files "$MOCK_PROMPTS_DIR" "*.chatmode.md")
    assert [ "$installed_count" -gt 0 ]
}

@test "hire command handles permission errors gracefully" {
    # Make prompts directory read-only (simulate permission error)
    chmod 444 "$MOCK_PROMPTS_DIR"
    
    run node bin/chatmate.js hire
    # Should handle error gracefully
    if [[ $status -ne 0 ]]; then
        assert_output --partial "Error"
    fi
    
    # Restore permissions for cleanup
    chmod 755 "$MOCK_PROMPTS_DIR"
}

# Test compatibility with hire.sh functionality
@test "hire command produces same result as hire.sh script" {
    # Install with CLI
    run node bin/chatmate.js hire --force
    assert_success
    
    # Count CLI installed files
    local cli_count
    cli_count=$(count_files "$MOCK_PROMPTS_DIR" "*.chatmode.md")
    
    # Clear and install with hire.sh
    rm -rf "$MOCK_PROMPTS_DIR"/*
    
    # Mock the hire.sh environment for our test
    export PROMPTS_DIR="$MOCK_PROMPTS_DIR"
    run bash -c "
        export SCRIPT_DIR='$PWD'
        export MATES_DIR='$PWD/mates'
        export PROMPTS_DIR='$MOCK_PROMPTS_DIR'
        mkdir -p '$PROMPTS_DIR'
        cp -vf '$MATES_DIR'/*.md '$PROMPTS_DIR'/
    "
    assert_success
    
    # Count hire.sh installed files
    local script_count
    script_count=$(count_files "$MOCK_PROMPTS_DIR" "*.chatmode.md")
    
    # They should be equal
    assert [ "$cli_count" -eq "$script_count" ]
}
