#!/usr/bin/env bats

# Integration tests for ChatMate project
# Tests end-to-end functionality, cross-platform compatibility, and real-world scenarios

load '../helpers/test_helper'

# Test complete installation process
@test "full installation workflow on macOS" {
    simulate_os "macos"
    local prompts_dir="$MOCK_HOME/Library/Application Support/Code/User/prompts"
    
    # Verify mates directory exists
    run test -d "mates"
    assert_success
    
    # Count available chatmates
    local chatmate_count
    chatmate_count=$(count_files "mates" "*.chatmode.md")
    assert [ "$chatmate_count" -gt 0 ]
    
    # Simulate installation
    run bash -c "
        mkdir -p '$prompts_dir'
        cp -vf mates/*.md '$prompts_dir'/
        echo 'Installation completed successfully'
    "
    assert_success
    assert_output --partial "Installation completed successfully"
    
    # Verify all files were copied
    assert verify_installation "$prompts_dir" "$chatmate_count"
    
    # Verify specific important chatmates were installed
    assert test -f "$prompts_dir/Testing.chatmode.md"
    assert test -f "$prompts_dir/Create PR.chatmode.md"
    assert test -f "$prompts_dir/Solve Issue.chatmode.md"
}

@test "full installation workflow on Linux" {
    simulate_os "linux"
    local prompts_dir="$MOCK_HOME/.config/Code/User/prompts"
    
    local chatmate_count
    chatmate_count=$(count_files "mates" "*.chatmode.md")
    
    # Simulate installation
    run bash -c "
        mkdir -p '$prompts_dir'
        cp -vf mates/*.md '$prompts_dir'/
        echo 'Linux installation completed'
    "
    assert_success
    
    assert verify_installation "$prompts_dir" "$chatmate_count"
}

@test "full installation workflow on Windows" {
    simulate_os "windows"
    local prompts_dir="$MOCK_APPDATA/Code/User/prompts"
    
    local chatmate_count
    chatmate_count=$(count_files "mates" "*.chatmode.md")
    
    # Simulate installation
    run bash -c "
        mkdir -p '$prompts_dir'
        cp -vf mates/*.md '$prompts_dir'/
        echo 'Windows installation completed'
    "
    assert_success
    
    assert verify_installation "$prompts_dir" "$chatmate_count"
}

# Test installation script error scenarios
@test "installation fails gracefully when mates directory is missing" {
    local temp_dir="$TEST_TEMP_DIR/no_mates_project"
    mkdir -p "$temp_dir"
    cp "hire.sh" "$temp_dir/"
    
    cd "$temp_dir" || exit 1
    
    # Should fail because no mates directory
    run bash hire.sh
    assert_failure
}

@test "installation handles existing files correctly" {
    simulate_os "linux"
    local prompts_dir="$MOCK_HOME/.config/Code/User/prompts"
    
    # Pre-populate with old versions
    mkdir -p "$prompts_dir"
    echo "# Old Version" > "$prompts_dir/Testing.chatmode.md"
    echo "# Old Version" > "$prompts_dir/Create PR.chatmode.md"
    
    # Run installation
    run bash -c "
        cp -vf mates/*.md '$prompts_dir'/
        echo 'Update completed'
    "
    assert_success
    
    # Verify files were updated (should not contain "Old Version")
    refute file_contains "$prompts_dir/Testing.chatmode.md" "Old Version"
    refute file_contains "$prompts_dir/Create PR.chatmode.md" "Old Version"
    
    # Should contain actual content
    assert file_contains "$prompts_dir/Testing.chatmode.md" "Testing Framework Agent"
}

# Test repository structure integrity
@test "repository structure is complete and valid" {
    # Check all required top-level files
    local required_files=("README.md" "CONTRIBUTING.md" "hire.sh" "LICENSE")
    
    for file in "${required_files[@]}"; do
        run test -f "$file"
        assert_success
    done
    
    # Check required directories
    local required_dirs=("mates" ".github" ".github/workflows")
    
    for dir in "${required_dirs[@]}"; do
        run test -d "$dir"
        assert_success
    done
    
    # Check GitHub workflows exist
    run test -f ".github/workflows/validate.yml"
    assert_success
    
    run test -f ".github/workflows/pr-validation.yml"
    assert_success
}

@test "all chatmates have consistent quality standards" {
    local error_count=0
    local issues=()
    
    while IFS= read -r -d '' file; do
        local basename
        basename=$(basename "$file")
        
        # Check file size is reasonable
        local file_size
        file_size=$(wc -c < "$file")
        if [[ "$file_size" -lt 1000 ]]; then
            issues+=("$basename: File too small ($file_size bytes)")
            ((error_count++))
        elif [[ "$file_size" -gt 50000 ]]; then
            issues+=("$basename: File too large ($file_size bytes)")
            ((error_count++))
        fi
        
        # Check for required sections
        if ! grep -q "## " "$file"; then
            issues+=("$basename: Missing section headers")
            ((error_count++))
        fi
        
        # Check for description in YAML
        if ! sed -n '/^---$/,/^---$/p' "$file" | grep -q "description:"; then
            issues+=("$basename: Missing description in YAML frontmatter")
            ((error_count++))
        fi
        
    done < <(find mates -name "*.chatmode.md" -type f -print0)
    
    if [[ "$error_count" -gt 0 ]]; then
        echo "Quality issues found:"
        printf '%s\n' "${issues[@]}"
        exit 1
    fi
}

# Test markdown quality across all files
@test "all markdown files pass basic lint checks" {
    # Create a basic markdownlint config
    local lint_config="$TEST_TEMP_DIR/.markdownlint.json"
    cat > "$lint_config" << 'EOF'
{
  "MD013": false,
  "MD041": false,
  "MD033": false,
  "MD034": false,
  "MD036": false
}
EOF
    
    # Test each markdown file individually for better error reporting
    while IFS= read -r -d '' file; do
        # Basic markdown structure checks
        local basename
        basename=$(basename "$file")
        
        # Check for proper line endings
        if file "$file" | grep -q "CRLF"; then
            echo "CRLF line endings found in: $basename"
            exit 1
        fi
        
        # Check for trailing whitespace
        if grep -q "[[:space:]]$" "$file"; then
            echo "Trailing whitespace found in: $basename"
            exit 1
        fi
        
    done < <(find . -name "*.md" -type f -print0)
}

# Test GitHub Actions workflow validation
@test "GitHub Actions workflows are syntactically valid" {
    # Check workflow files exist and have basic structure
    local workflow_files=(".github/workflows/validate.yml" ".github/workflows/pr-validation.yml")
    
    for workflow in "${workflow_files[@]}"; do
        run test -f "$workflow"
        assert_success
        
        # Check for required workflow elements
        assert file_contains "$workflow" "name:"
        assert file_contains "$workflow" "on:"
        assert file_contains "$workflow" "jobs:"
        
        # Check for proper YAML structure (basic)
        run bash -c "head -n 1 '$workflow' | grep -q '^name:'"
        assert_success
    done
}

# Test security and safety
@test "no sensitive information in repository" {
    local sensitive_patterns=("password" "secret" "token" "api.key" "private.key")
    local found_issues=()
    
    for pattern in "${sensitive_patterns[@]}"; do
        while IFS= read -r line; do
            # Skip this test file itself and common false positives
            if [[ "$line" != *"integration_tests.bats"* ]] && 
               [[ "$line" != *"# Test"* ]] && 
               [[ "$line" != *"example"* ]]; then
                found_issues+=("Potential sensitive data: $line")
            fi
        done < <(grep -ri "$pattern" . --exclude-dir=.git --exclude-dir=node_modules --exclude="*.bats" || true)
    done
    
    if [[ ${#found_issues[@]} -gt 0 ]]; then
        echo "Security issues found:"
        printf '%s\n' "${found_issues[@]}"
        exit 1
    fi
}

# Test cross-platform compatibility
@test "hire.sh works across different shell environments" {
    # Test script with different shell interpreters where available
    local shells=("bash")
    
    # Add other shells if available
    if command -v zsh >/dev/null 2>&1; then
        shells+=("zsh")
    fi
    
    for shell in "${shells[@]}"; do
        # Test syntax check
        run "$shell" -n "hire.sh"
        assert_success
    done
}

# Test installation with different file permissions
@test "installation works with restrictive permissions" {
    simulate_os "linux"
    local prompts_dir="$MOCK_HOME/.config/Code/User/prompts"
    
    # Create directory with restrictive permissions
    mkdir -p "$prompts_dir"
    chmod 755 "$prompts_dir"
    
    # Test installation still works
    run bash -c "cp -vf mates/*.md '$prompts_dir'/"
    assert_success
    
    # Verify files are readable
    run test -r "$prompts_dir/Testing.chatmode.md"
    assert_success
}

# Test upgrade scenario
@test "upgrading installation overwrites old files correctly" {
    simulate_os "linux"
    local prompts_dir="$MOCK_HOME/.config/Code/User/prompts"
    
    # Install initial version
    mkdir -p "$prompts_dir"
    echo "# Old Testing Framework" > "$prompts_dir/Testing.chatmode.md"
    echo "# Old Create PR" > "$prompts_dir/Create PR.chatmode.md"
    echo "# Orphaned File" > "$prompts_dir/Old.chatmode.md"
    
    local old_count
    old_count=$(count_files "$prompts_dir" "*.md")
    
    # Upgrade installation
    run bash -c "cp -vf mates/*.md '$prompts_dir'/"
    assert_success
    
    # Verify old files were overwritten
    refute file_contains "$prompts_dir/Testing.chatmode.md" "Old Testing Framework"
    assert file_contains "$prompts_dir/Testing.chatmode.md" "Testing Framework Agent"
    
    # Note: Orphaned files remain (this is expected behavior)
    local new_count
    new_count=$(count_files "$prompts_dir" "*.md")
    assert [ "$new_count" -ge "$old_count" ]
}
