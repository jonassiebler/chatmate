## Security Guidelines

Security is a top priority for Chatmate:

### Security Guidelines

1. **Input Validation**
   - Validate all user inputs
   - Sanitize file paths
   - Check for malicious content

2. **File Operations**
   - Use safe file permissions (0644 for files, 0755 for directories)
   - Validate file paths to prevent directory traversal
   - Handle symlinks carefully

3. **External Dependencies**
   - Keep dependencies minimal and up-to-date
   - Review security advisories
   - Use `go mod audit` for vulnerability scanning

### Security Testing

```bash
# Run security scan
./chatmate status --security-scan
# Check for vulnerabilities
go list -json -deps | nancy sleuth
# Static analysis
go vet ./...
staticcheck ./...
```