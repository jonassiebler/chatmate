#!/usr/bin/env bats

# CLI Commands Tests
# Tests list, uninstall, status, and config commands

load '../helpers/test_helper'

setup() {
    # Call parent setup
    load '../helpers/test_helper'
    
    # Create temporary prompts directory for testing
    export TEST_PROMPTS_DIR="$TEST_TEMP_DIR/prompts"
    mkdir -p "$TEST_PROMPTS_DIR"
    
    # Mock VS Code prompts directory
    export ORIGINAL_HOME="$HOME"
    export HOME="$TEST_TEMP_DIR"
    
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
    
    # Pre-install some chatmates for testing
    node bin/chatmate.js hire --specific "Solve Issue" "Create PR" >/dev/null 2>&1
}

teardown() {
    export HOME="$ORIGINAL_HOME"
    load '../helpers/test_helper'
}

# Test list command
@test "list command shows available and installed chatmates" {
    run node bin/chatmate.js list
    assert_success
    assert_output --partial "Available Chatmates:"
    assert_output --partial "Installed Chatmates:"
    assert_output --partial "Summary:"
}

@test "list command with --available flag shows only available" {
    run node bin/chatmate.js list --available
    assert_success
    assert_output --partial "Available Chatmates:"
    # Should not show installed section
    refute_output --partial "Installed Chatmates:"
}

@test "list command with --installed flag shows only installed" {
    run node bin/chatmate.js list --installed
    assert_success
    assert_output --partial "Installed Chatmates:"
    # Should not show available section
    refute_output --partial "Available Chatmates:"
}

@test "list command shows correct installation status" {
    run node bin/chatmate.js list
    assert_success
    # Should show our pre-installed chatmates
    assert_output --partial "Solve Issue"
    assert_output --partial "Create PR"
    assert_output --partial "installed"
}

# Test uninstall command
@test "uninstall command removes specific chatmate" {
    # Verify chatmate is installed
    run test -f "$MOCK_PROMPTS_DIR/Solve Issue.chatmode.md"
    assert_success
    
    # Uninstall it
    run node bin/chatmate.js uninstall "Solve Issue"
    assert_success
    assert_output --partial "Uninstalling chatmate agents"
    assert_output --partial "Solve Issue.chatmode.md (uninstalled)"
    
    # Verify it's removed
    run test -f "$MOCK_PROMPTS_DIR/Solve Issue.chatmode.md"
    assert_failure
}

@test "uninstall command with --all flag removes all chatmates" {
    # Verify chatmates are installed
    run test -f "$MOCK_PROMPTS_DIR/Solve Issue.chatmode.md"
    assert_success
    run test -f "$MOCK_PROMPTS_DIR/Create PR.chatmode.md"
    assert_success
    
    # Uninstall all
    run node bin/chatmate.js uninstall --all
    assert_success
    assert_output --partial "Uninstalling chatmate agents"
    
    # Verify all are removed
    local remaining_count
    remaining_count=$(count_files "$MOCK_PROMPTS_DIR" "*.chatmode.md")
    assert [ "$remaining_count" -eq 0 ]
}

@test "uninstall command warns about non-existent chatmate" {
    run node bin/chatmate.js uninstall "NonExistent"
    assert_success
    assert_output --partial "No installed chatmate found matching: NonExistent"
}

@test "uninstall command handles multiple chatmates" {
    run node bin/chatmate.js uninstall "Solve Issue" "Create PR"
    assert_success
    assert_output --partial "Solve Issue.chatmode.md (uninstalled)"
    assert_output --partial "Create PR.chatmode.md (uninstalled)"
}

# Test status command
@test "status command shows VS Code and prompts directory status" {
    run node bin/chatmate.js status
    assert_success
    assert_output --partial "ChatMate Installation Status"
    assert_output --partial "VS Code:"
    assert_output --partial "Prompts Directory:"
    assert_output --partial "Chatmate Statistics:"
}

@test "status command shows correct chatmate counts" {
    run node bin/chatmate.js status
    assert_success
    assert_output --partial "Available:"
    assert_output --partial "Installed:"
    # Should show that we have some installed chatmates
    assert_output --partial "chatmates"
}

@test "status command detects prompts directory" {
    run node bin/chatmate.js status
    assert_success
    assert_output --partial "Prompts directory exists"
    assert_output --partial "$MOCK_PROMPTS_DIR"
}

# Test config command
@test "config command with --show displays configuration" {
    run node bin/chatmate.js config --show
    assert_success
    assert_output --partial "ChatMate Configuration:"
    assert_output --partial "Mates Directory:"
    assert_output --partial "Prompts Directory:"
    assert_output --partial "Platform:"
    assert_output --partial "Node Version:"
}

@test "config command without flags shows help" {
    run node bin/chatmate.js config
    assert_success
    assert_output --partial "Configuration Management:"
    assert_output --partial "--show"
    assert_output --partial "--reset"
}

@test "config command with --reset shows reset message" {
    run node bin/chatmate.js config --reset
    assert_success
    assert_output --partial "Configuration reset is not implemented yet"
}

# Test error handling
@test "commands handle missing VS Code gracefully" {
    # This test verifies the CLI doesn't crash when VS Code isn't found
    run node bin/chatmate.js status
    # Should succeed regardless of VS Code installation
    assert_success
}

@test "commands work with relative paths" {
    # Change to a different directory and test
    cd "$TEST_TEMP_DIR" || exit 1
    
    run node "$ORIGINAL_PWD/bin/chatmate.js" --version
    assert_success
    assert_output "1.0.0"
}

# Test cross-platform compatibility
@test "CLI detects correct prompts directory for platform" {
    run node bin/chatmate.js config --show
    assert_success
    
    case "$(uname)" in
        Darwin)
            assert_output --partial "Library/Application Support/Code/User/prompts"
            ;;
        Linux)
            assert_output --partial ".config/Code/User/prompts"
            ;;
    esac
}
