#!/usr/bin/env bats

# CLI Basic Functionality Tests
# Tests core CLI commands, help system, and basic functionality

load '../helpers/test_helper'

# Test CLI binary existence and setup
@test "chatmate CLI binary exists and is executable" {
    run test -f "bin/chatmate.js"
    assert_success
    
    run test -x "bin/chatmate.js"
    assert_success
}

@test "CLI has valid Node.js shebang" {
    run head -n 1 "bin/chatmate.js"
    assert_output --partial "#!/usr/bin/env node"
}

@test "CLI displays help when run without arguments" {
    run node bin/chatmate.js
    assert_success
    assert_output --partial "Usage: chatmate"
    assert_output --partial "Commands:"
}

@test "CLI displays version with --version flag" {
    run node bin/chatmate.js --version
    assert_success
    assert_output "1.0.0"
}

@test "CLI displays help with --help flag" {
    run node bin/chatmate.js --help
    assert_success
    assert_output --partial "Usage: chatmate"
    assert_output --partial "hire"
    assert_output --partial "list"
    assert_output --partial "uninstall"
    assert_output --partial "status"
    assert_output --partial "config"
}

# Test main commands exist
@test "hire command shows help" {
    run node bin/chatmate.js hire --help
    assert_success
    assert_output --partial "Install all chatmate agents"
    assert_output --partial "--specific"
    assert_output --partial "--force"
}

@test "list command shows help" {
    run node bin/chatmate.js list --help
    assert_success
    assert_output --partial "List available and installed chatmate agents"
    assert_output --partial "--available"
    assert_output --partial "--installed"
}

@test "uninstall command shows help" {
    run node bin/chatmate.js uninstall --help
    assert_success
    assert_output --partial "Uninstall specific chatmate agents"
    assert_output --partial "--all"
}

@test "status command shows help" {
    run node bin/chatmate.js status --help
    assert_success
    assert_output --partial "Show ChatMate and VS Code installation status"
}

@test "config command shows help" {
    run node bin/chatmate.js config --help
    assert_success
    assert_output --partial "Manage ChatMate configuration settings"
    assert_output --partial "--show"
    assert_output --partial "--reset"
}

# Test invalid commands
@test "CLI shows error for invalid command" {
    run node bin/chatmate.js invalid-command
    assert_failure
    assert_output --partial "error: unknown command 'invalid-command'"
}

# Test CLI dependencies
@test "CLI can import required dependencies" {
    run node -e "
        require('./bin/chatmate.js');
        console.log('Dependencies loaded successfully');
    " 2>/dev/null
    assert_success
}

@test "CLI lib modules can be required" {
    run node -e "
        const { ChatMateManager } = require('./lib/chatmate-manager');
        console.log('ChatMateManager loaded');
    "
    assert_success
    assert_output "ChatMateManager loaded"
}
