---
description: 'Invalid Test Chatmate'
# Missing author field (required)
model: 'Claude Sonnet 4'
tools: ['codebase']
---

# Invalid Chatmate - Missing Required Fields

This chatmate has YAML frontmatter but is missing required fields like 'author'.

## Purpose

Test validation of incomplete YAML frontmatter.

## Expected Behavior

- Should fail validation due to missing required author field
- Used to test YAML field validation
- Should be rejected during installation process

## Testing Notes

This file intentionally omits the required 'author' field to test validation error handling.
