# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| Latest  | :white_check_mark: |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue, please follow these steps:

1. **DO NOT** create a public GitHub issue
2. Email security details to the maintainers
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if available)

## Security Best Practices

When using ERD:

1. **Diagram Content**: ERD diagrams may reveal your domain model structure. Ensure you don't expose sensitive entity relationships in public diagrams.

2. **Validation**: Always validate diagram structures before rendering to prevent potential injection attacks in diagram output.

3. **File Generation**: If generating diagram files programmatically, ensure proper file path validation to prevent path traversal attacks.

## Security Features

ERD is designed with security in mind:

- Zero external dependencies (reduces supply chain risks)
- No network operations
- No file system operations (output is returned as strings)
- Pure data structure manipulation
- Thread-safe operations

## Acknowledgments

We appreciate responsible disclosure of security vulnerabilities.
