## Development Guidelines

### Code Style

#### Go Style Guidelines

1. **Follow standard Go conventions**
   - Use `gofmt` for formatting
   - Follow effective Go practices
   - Use meaningful variable names

2. **Documentation**
   - Add package documentation
   - Document all public functions
   - Include usage examples in GoDoc

3. **Error Handling**
   ```go
   // Good: Specific error messages
   if err != nil {
       return fmt.Errorf("failed to create chatmate %s: %w", name, err)
   }

   // Bad: Generic error handling
   if err != nil {
       return err
   }
   ```

4. **CLI Output**
   - Use emojis for visual appeal (✅ ❌ 🔍 🚀)
   - Provide clear, actionable messages
   - Include progress indicators for long operations

### Code Organization

1. **Package Structure**
   - Keep packages focused and cohesive
   - Minimize dependencies between packages
   - Use internal/ for private code

2. **Function Design**
   - Single responsibility principle
   - Clear input/output contracts
   - Comprehensive error handling

### Testing

We use a multi-layered testing approach:

#### Test Types

1. **Unit Tests** (Go's built-in testing)
   ```bash
   go test ./...
   go test -cover ./...
   go test ./internal/manager/
   ```
2. **CLI Tests** (BATS framework)
   ```bash
   ./run-tests.sh
   ./tests/bats-core/bin/bats tests/unit/hire_script.bats
   ./tests/bats-core/bin/bats -t tests/unit/
   ```
3. **Integration Tests**
   ```bash
   ./tests/bats-core/bin/bats tests/integration/
   ```

#### Writing Tests
##### BATS Test Example
```bash
#!/usr/bin/env bats
load '../helpers/test_helper'
@test "hire command creates chatmate file" {
    run ./chatmate hire "Test Agent"
    assert_success
    assert_output --partial "successfully hired"
}
@test "hire command with invalid input fails" {
    run ./chatmate hire ""
    assert_failure
    assert_output --partial "name is required"
}
```
##### Go Test Example
```go
func TestManager_HireChatmate(t *testing.T) {
    manager := &ChatmateManager{}
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid name", "Test Agent", false},
        {"empty name", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := manager.HireChatmate(tt.input, "", "")
            if (err != nil) != tt.wantErr {
                t.Errorf("HireChatmate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Adding New Commands

1. **Create command file**
   ```bash
   touch cmd/newcommand.go
   ```
2. **Implement command structure**
   ```go
   package cmd
   import (
       "github.com/spf13/cobra"
   )
   var newcommandCmd = &cobra.Command{
       Use:   "newcommand",
       Short: "Short description",
       Long:  `Detailed description...`,
       Run: func(cmd *cobra.Command, args []string) {
           // Implementation
       },
   }
   func init() {
       rootCmd.AddCommand(newcommandCmd)
   }
   ```
3. **Add tests**
   ```bash
   touch tests/unit/newcommand.bats
   touch tests/integration/newcommand_integration.bats
   ```

### Adding Security Features

1. **Extend security validation**
   - Add new checks to `pkg/security/validation.go`
   - Include comprehensive error handling
   - Add corresponding tests
2. **Update security scanning**
   - Enhance the `SecurityScan` function
   - Add new vulnerability patterns
   - Update security scoring