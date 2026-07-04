# AGENTS.md

## Overview

MarkdownForge is a comprehensive Markdown processing toolkit written in Go. It provides tools for parsing, analyzing, linting, and converting Markdown documents.

## Architecture

### Core Modules

1. **parser/** - Markdown AST parser with block and inline element extraction
2. **toc/** - Table of contents generator with nested hierarchy
3. **links/** - Link checker with HTTP and file validation
4. **extract/** - Content extractor for headings, code blocks, links, images
5. **lint/** - Markdown linter with customizable rules
6. **convert/** - HTML and plain text conversion using goldmark
7. **stats/** - Document statistics (word count, reading time, etc.)
8. **badge/** - Shields.io badge generator

## Development

### Building

```bash
go build -o markdownforge .
```

### Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/parser/...
go test ./internal/toc/...
```

### Adding New Lint Rules

1. Add a check function in `internal/lint/lint.go`
2. Register it in `AddDefaultRules()`
3. Add tests in `internal/lint/lint_test.go`

### Adding New Commands

1. Create a new command function in `cmd/markdownforge/main.go`
2. Follow the pattern of existing commands
3. Add the command to the root command

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/fatih/color` - Terminal colors
- `github.com/yuin/goldmark` - Markdown to HTML conversion

## Testing Guidelines

- Test all public functions
- Test edge cases (empty input, invalid input)
- Test error paths
- Use table-driven tests where appropriate
- Aim for >80% coverage
