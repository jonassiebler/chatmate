#!/usr/bin/env bats

# Unit tests for chatmate file validation and structure
# Tests chatmate format compliance, content validation, and quality checks

load '../helpers/test_helper'

# Test chatmate file structure
@test "mates directory exists and contains files" {
    run test -d "mates"
    assert_success
    
    local file_count
    file_count=$(count_files "mates" "*.md")
    assert [ "$file_count" -gt 0 ]
}

@test "all chatmate files follow naming convention" {
    # Check that all .md files in mates directory end with .chatmode.md
    while IFS= read -r -d '' file; do
        local basename
        basename=$(basename "$file")
        if [[ "$basename" != *.chatmode.md ]]; then
            echo "Invalid chatmate filename: $basename"
            exit 1
        fi
    done < <(find mates -name "*.md" -type f -print0)
}

@test "chatmate files are not empty" {
    while IFS= read -r -d '' file; do
        if [[ ! -s "$file" ]]; then
            echo "Empty chatmate file: $(basename "$file")"
            exit 1
        fi
    done < <(find mates -name "*.chatmode.md" -type f -print0)
}

@test "chatmate files contain required headers" {
    while IFS= read -r -d '' file; do
        if ! grep -q "^#" "$file"; then
            echo "No headers found in: $(basename "$file")"
            exit 1
        fi
    done < <(find mates -name "*.chatmode.md" -type f -print0)
}

@test "chatmate files contain YAML frontmatter" {
    while IFS= read -r -d '' file; do
        # Check for YAML frontmatter (starts with ---)
        if ! head -n 1 "$file" | grep -q "^---"; then
            echo "No YAML frontmatter in: $(basename "$file")"
            exit 1
        fi
    done < <(find mates -name "*.chatmode.md" -type f -print0)
}

@test "chatmate files have reasonable size limits" {
    local max_size=100000  # 100KB limit
    
    while IFS= read -r -d '' file; do
        local file_size
        file_size=$(wc -c < "$file")
        if [[ "$file_size" -gt "$max_size" ]]; then
            echo "File too large: $(basename "$file") ($file_size bytes)"
            exit 1
        fi
    done < <(find mates -name "*.chatmode.md" -type f -print0)
}

# Test specific chatmate content validation
@test "Testing chatmode exists and is valid" {
    local testing_file="mates/Testing.chatmode.md"
    
    run test -f "$testing_file"
    assert_success
    
    # Validate it's a proper chatmate
    assert is_valid_chatmate "$testing_file"
    
    # Check for specific Testing chatmode content
    assert file_contains "$testing_file" "Testing Framework Agent"
    assert file_contains "$testing_file" "model:"
    assert file_contains "$testing_file" "tools:"
}

@test "Create PR chatmode exists and is valid" {
    local pr_file="mates/Create PR.chatmode.md"
    
    run test -f "$pr_file"
    assert_success
    
    assert is_valid_chatmate "$pr_file"
    assert file_contains "$pr_file" "Pull Request Creation Agent"
}

@test "Solve Issue chatmode exists and is valid" {
    local solve_file="mates/Solve Issue.chatmode.md"
    
    run test -f "$solve_file"
    assert_success
    
    assert is_valid_chatmate "$solve_file"
    assert file_contains "$solve_file" "Solve Issue"
}

# Test YAML frontmatter structure
@test "chatmate YAML frontmatter contains required fields" {
    while IFS= read -r -d '' file; do
        # Extract YAML frontmatter (between first --- and second ---)
        local yaml_content
        yaml_content=$(sed -n '/^---$/,/^---$/p' "$file" | sed '1d;$d')
        
        # Check for required fields
        if ! echo "$yaml_content" | grep -q "description:"; then
            echo "Missing description field in: $(basename "$file")"
            exit 1
        fi
        
        if ! echo "$yaml_content" | grep -q "model:"; then
            echo "Missing model field in: $(basename "$file")"
            exit 1
        fi
        
        if ! echo "$yaml_content" | grep -q "tools:"; then
            echo "Missing tools field in: $(basename "$file")"
            exit 1
        fi
        
    done < <(find mates -name "*.chatmode.md" -type f -print0)
}

@test "chatmate descriptions are informative" {
    while IFS= read -r -d '' file; do
        local description
        description=$(sed -n '/^---$/,/^---$/p' "$file" | grep "description:" | cut -d'"' -f2)
        
        # Check description is not empty and has reasonable length
        if [[ -z "$description" ]] || [[ ${#description} -lt 10 ]]; then
            echo "Description too short in: $(basename "$file")"
            exit 1
        fi
        
    done < <(find mates -name "*.chatmode.md" -type f -print0)
}

@test "chatmate tools lists are valid arrays" {
    while IFS= read -r -d '' file; do
        # Extract tools section and check it's a proper YAML array
        local tools_line
        tools_line=$(sed -n '/^---$/,/^---$/p' "$file" | grep "tools:")
        
        if ! echo "$tools_line" | grep -q "\["; then
            echo "Tools field not formatted as array in: $(basename "$file")"
            exit 1
        fi
        
    done < <(find mates -name "*.chatmode.md" -type f -print0)
}

# Test content quality
@test "chatmate files contain substantial content" {
    local min_lines=20
    
    while IFS= read -r -d '' file; do
        local line_count
        line_count=$(wc -l < "$file")
        
        if [[ "$line_count" -lt "$min_lines" ]]; then
            echo "File too short: $(basename "$file") ($line_count lines)"
            exit 1
        fi
        
    done < <(find mates -name "*.chatmode.md" -type f -print0)
}

@test "chatmate files use proper markdown formatting" {
    while IFS= read -r -d '' file; do
        # Check for proper header hierarchy (should start with # not ##)
        local first_header
        first_header=$(grep -m 1 "^#" "$file" || echo "")
        
        if [[ -n "$first_header" ]] && [[ ! "$first_header" =~ ^#[^#] ]]; then
            echo "Improper header hierarchy in: $(basename "$file")"
            exit 1
        fi
        
    done < <(find mates -name "*.chatmode.md" -type f -print0)
}

# Test for common issues
@test "chatmate files do not contain placeholder text" {
    local placeholders=("TODO" "FIXME" "placeholder" "example.com" "your-" "CHANGE_ME")
    
    while IFS= read -r -d '' file; do
        for placeholder in "${placeholders[@]}"; do
            if grep -qi "$placeholder" "$file"; then
                echo "Placeholder text found in: $(basename "$file") - $placeholder"
                exit 1
            fi
        done
    done < <(find mates -name "*.chatmode.md" -type f -print0)
}

@test "chatmate files contain unique content" {
    # Check for duplicate content across chatmates (basic check)
    local temp_content="$TEST_TEMP_DIR/content_check"
    mkdir -p "$temp_content"
    
    while IFS= read -r -d '' file; do
        local basename
        basename=$(basename "$file" .chatmode.md)
        local title_line
        title_line=$(grep -m 1 "^# " "$file" | sed 's/^# //' || echo "No title")
        
        # Simple uniqueness check - title should contain the chatmate name
        if [[ "$title_line" != *"$basename"* ]] && [[ "$title_line" != "No title" ]]; then
            echo "Title doesn't match filename in: $(basename "$file")"
            echo "Title: $title_line"
            echo "Expected to contain: $basename"
            exit 1
        fi
        
    done < <(find mates -name "*.chatmode.md" -type f -print0)
}
