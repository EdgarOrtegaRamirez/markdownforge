# Security Policy

## Overview

MarkdownForge processes local Markdown files and makes HTTP requests to validate links. This document outlines the security considerations for using this tool.

## Network Usage

When using the `links` command, MarkdownForge makes HTTP HEAD requests to verify external URLs. These requests:

- Do not send any sensitive data
- Use HEAD requests only (no data transfer)
- Respect timeout settings (default 10 seconds)
- Do not store or log response bodies

### Link Checking Security

- Only HEAD requests are made (no data transmission)
- Request timeout is configurable (default 10 seconds)
- No cookies or authentication headers are sent
- No user-agent modification

## Input Validation

All user inputs are validated:

- File paths are checked for existence before processing
- URLs are parsed and validated before requests
- Pattern matching uses standard Go regex (RE2)
- File operations use safe path handling

## Dependencies

MarkdownForge uses minimal dependencies:

- `github.com/spf13/cobra` - CLI framework (widely used, actively maintained)
- `github.com/fatih/color` - Terminal colors (minimal)
- `github.com/yuin/goldmark` - Markdown parser (pure Go, no native code)

## Vulnerability Reporting

If you discover a security vulnerability, please report it responsibly:

1. Do not open a public GitHub issue
2. Email [security@example.com] (replace with actual)
3. Include detailed reproduction steps
4. Allow reasonable time for response before disclosure

## Best Practices

When using MarkdownForge:

- Only process Markdown files you trust
- Use the `--timeout` flag for link checking on untrusted content
- Review extracted content before using it
- Do not pass user-controlled input as command-line arguments to system commands

## Audit Status

- Last security review: 2026-07-04
- No known vulnerabilities
- No use of unsafe Go packages
