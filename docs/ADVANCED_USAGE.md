# ChatMate Advanced Usage Guide 🚀

This guide covers advanced ChatMate features, customization, and automation for power users.

## Usage Patterns

### For Individual Developers

**Daily Development Workflow:**
```bash
# Morning setup - check status and install needed chatmates
chatmate status
chatmate list --installed

# Work on specific task - use appropriate chatmates in VS Code
# @Solve Issue for debugging
# @Code Review for quality checks
# @Testing for test creation

# End of day - clean up unused chatmates
chatmate list --installed
chatmate uninstall "Rarely Used Chatmate"
```

### For Team Leads

**Team Management:**
```bash
# Install comprehensive set for team leadership
chatmate hire "Code Review" "Create PR" "Review PR" "Create Issue"

# Weekly team support workflow:
# Use @Review PR for pull request reviews
# Use @Create Issue for task assignment
# Use @Code Review for architecture decisions
```

### For Project Managers

**Project Oversight:**
```bash
# Install project management focused chatmates
chatmate hire "Create Issue" "Optimize Issues" "Review PR" "Merge PR"

# Daily project management:
# Use @Create Issue for requirements documentation
# Use @Review PR for understanding technical changes
# Use @Optimize Issues for clarifying requirements
```

### For QA Engineers

**Quality Assurance:**
```bash
# Install testing and quality focused chatmates
chatmate hire "Testing" "Code Review" "Solve Issue"

# QA workflow:
# Use @Testing for test case generation
# Use @Code Review for code quality assessment
# Use @Solve Issue for bug investigation
```

## Best Practices

### Chatmate Management

#### **Start Simple**
- Begin with `chatmate hire` to install all chatmates
- Use `chatmate list` to familiarize yourself with available agents
- Remove chatmates you don't use with `chatmate uninstall`

#### **Keep Organized**
```bash
# Weekly maintenance routine
chatmate status                    # Check system health
chatmate list --installed         # Review installed chatmates
chatmate uninstall "Unused One"   # Remove chatmates you don't use
chatmate hire --force             # Update existing chatmates
```

#### **Team Synchronization**
- Use `chatmate list > team-chatmates.txt` to share team configurations
- Document which chatmates are recommended for specific roles
- Keep chatmate selections consistent across team members

### VS Code Integration

#### **Effective Usage in Copilot Chat**
- **Be specific**: "@Solve Issue The login form validation isn't working correctly"
- **Provide context**: Include relevant code, error messages, or requirements
- **Use appropriate chatmate**: Match the task to the specialized agent
- **Follow up**: Ask clarifying questions to get better responses

#### **Workflow Integration**
- Use `@Code Review` before committing changes
- Use `@Testing` when implementing new features
- Use `@Documentation` when creating public APIs
- Use `@Solve Issue` when debugging problems

### Performance Optimization

#### **Selective Installation**
```bash
# Instead of installing everything
chatmate hire

# Install only what you need
chatmate hire "Code Review" "Testing" "Solve Issue"
```

#### **Regular Maintenance**
```bash
# Monthly cleanup routine
chatmate list --installed          # Review current chatmates
chatmate uninstall "Unused Ones"   # Remove unused chatmates
chatmate hire --force             # Update remaining chatmates
```

### Security Considerations

#### **Verify Installations**
```bash
# Always verify after installation
chatmate status
chatmate list --installed

# Check file integrity
ls -la ~/Library/Application\ Support/Code/User/prompts/
```

#### **Keep Updated**
```bash
# Regular update routine
chatmate hire --force  # Update all chatmates
```

#### **Review Chatmate Content**
- Chatmates are open source - review their content
- Understand what each chatmate does before using
- Report any security concerns to the project maintainers

## Custom Chatmate Development

### Using Create Chatmode

1. Use the `@Create Chatmode` agent in VS Code:
   ```
   @Create Chatmode I need a chatmate specialized in database optimization
   ```

2. Follow the guidance to create your custom chatmate

3. Save as `.chatmode.md` file in the `mates/` directory

4. Install with `chatmate hire`

### Manual Chatmate Creation

Create a file like `Custom Agent.chatmode.md`:

```markdown
---
name: Custom Agent
description: Description of your custom agent
author: Your Name
version: 1.0.0
tags: [custom, specialized]
---

# Custom Agent Instructions

[Your specialized prompt content here]
```

## Automation and Scripting

### Automated Setup Scripts

```bash
#!/bin/bash
# team-setup.sh - Automated team chatmate setup

echo "Setting up ChatMate for development team..."

# Install core development chatmates
chatmate hire "Code Review" "Testing" "Solve Issue" "Documentation"

# Verify installation
if chatmate status | grep -q "✅"; then
    echo "✅ ChatMate setup complete!"
    echo "Restart VS Code to use your new chatmates"
else
    echo "❌ Setup failed. Check chatmate status"
    exit 1
fi
```

### CI/CD Integration

```yaml
# .github/workflows/chatmate-update.yml
name: Update ChatMates
on:
  schedule:
    - cron: '0 9 * * 1'  # Weekly on Mondays

jobs:
  update-chatmates:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Update ChatMates
        run: |
          chatmate hire --force
          chatmate status
```

## Environment-Specific Configurations

### Development Environment
```bash
# Install comprehensive development set
chatmate hire "Solve Issue" "Code Review" "Testing" "Documentation"
```

### Production Support
```bash
# Install focused production support set
chatmate hire "Solve Issue" "Code Review" "Documentation"
```

### Project Management
```bash
# Install project management focused set
chatmate hire "Create Issue" "Review PR" "Create PR" "Optimize Issues"
```

## Integration with Other Tools

### Git Hooks Integration
```bash
# .git/hooks/pre-commit
#!/bin/bash
# Use ChatMate for pre-commit code review
if command -v chatmate &> /dev/null; then
    echo "Running ChatMate code review..."
    # Add your custom integration here
fi
```

### IDE Integration
- Set up custom shortcuts for common chatmate commands
- Integrate with your IDE's task runner
- Use chatmate status checks in build processes

### Team Workflow Integration
- Include chatmate setup in onboarding documentation
- Use chatmates in code review checklists
- Integrate with team communication tools

---

For more information, see the [User Guide](USER_GUIDE.md) and [Troubleshooting Guide](TROUBLESHOOTING.md).
