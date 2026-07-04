# MarkdownForge

A comprehensive Markdown processing toolkit with CLI and library API.

## Features

- **TOC Generation** - Generate tables of contents from headings
- **Link Checking** - Validate all links (external HTTP + internal files)
- **Content Extraction** - Extract headings, code blocks, links, images, tables, blockquotes
- **Markdown Linting** - Check for common issues (heading levels, trailing whitespace, etc.)
- **HTML Conversion** - Convert Markdown to HTML or plain text
- **Badge Generation** - Generate shields.io-style badges
- **Document Statistics** - Word count, reading time, content analysis
- **Regex Extraction** - Extract content matching patterns

## Installation

```bash
go install github.com/EdgarOrtegaRamirez/markdownforge/cmd/markdownforge@latest
```

Or build from source:

```bash
git clone https://github.com/EdgarOrtegaRamirez/markdownforge
cd markdownforge
go build -o markdownforge .
```

## Quick Start

```bash
# Generate table of contents
markdownforge toc README.md

# Show document statistics
markdownforge stats README.md

# Lint a document
markdownforge lint README.md

# Check links
markdownforge links README.md

# Convert to HTML
markdownforge convert -f html README.md

# Extract code blocks
markdownforge extract -c go README.md

# Generate badges
markdownforge badge --name myproject --license MIT --language Go --tests passing
```

## Commands

### `toc` - Table of Contents

Generate a nested table of contents from document headings.

```bash
markdownforge toc document.md
```

Output:

```markdown
- [Title](#title)
  - [Section A](#section-a)
  - [Section B](#section-b)
    - [Subsection](#subsection)
```

### `stats` - Document Statistics

Show comprehensive statistics about a Markdown document.

```bash
markdownforge stats document.md
markdownforge stats -f json document.md  # JSON output
markdownforge stats -f markdown document.md  # Markdown table output
```

Output:

```
Document Statistics
==================
Words:          1234
Characters:     5678
Lines:          89
Paragraphs:     12
Headings:       8
Code Blocks:    3
Links:          5
Images:         2
Tables:         1
Lists:          4
Blockquotes:    2
Reading Time:   6 mins
```

### `lint` - Markdown Linting

Check for common Markdown issues.

```bash
markdownforge lint document.md
```

Rules:
- `heading-level` - Warns when heading levels skip (e.g., h1 to h3)
- `trailing-space` - Info about trailing whitespace
- `multiple-blanks` - Info about multiple consecutive blank lines
- `heading-punctuation` - Warns when headings end with punctuation
- `first-line-heading` - Info when document doesn't start with a heading
- `no-empty-links` - Warns about links with empty text

### `links` - Link Checking

Validate all links in a document.

```bash
markdownforge links document.md
markdownforge links --timeout 10s document.md
```

Output:

```
✓ https://example.com
✓ ./README.md
✗ https://broken-link.com: 404 Not Found

3 valid, 1 invalid
```

### `extract` - Content Extraction

Extract specific content from Markdown.

```bash
# Extract all links and images
markdownforge extract document.md

# Extract headings at level 2
markdownforge extract -l 2 document.md

# Extract code blocks by language
markdownforge extract -c python document.md

# Extract a specific section
markdownforge extract -s "Installation" document.md

# Extract by regex
markdownforge extract -r '\bTODO\b' document.md
```

### `convert` - Format Conversion

Convert Markdown to other formats.

```bash
markdownforge convert -f html document.md
markdownforge convert -f text document.md
```

### `badge` - Badge Generation

Generate shields.io-style badges.

```bash
markdownforge badge --name myproject --license MIT --language Go --tests passing
```

Output:

```markdown
![name](https://img.shields.io/badge/name-myproject-blue) ![license](https://img.shields.io/badge/license-MIT-green) ![language](https://img.shields.io/badge/language-Go-blue) ![tests](https://img.shields.io/badge/tests-passing-brightgreen)
```

## Library API

MarkdownForge can also be used as a Go library:

```go
package main

import (
    "fmt"
    "github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
    "github.com/EdgarOrtegaRamirez/markdownforge/internal/toc"
    "github.com/EdgarOrtegaRamirez/markdownforge/internal/stats"
    "github.com/EdgarOrtegaRamirez/markdownforge/internal/lint"
)

func main() {
    source := `# Title

## Section

Content here.
`
    
    // Parse document
    doc := parser.Parse(source)
    
    // Generate TOC
    t := toc.Generate(doc)
    fmt.Println(toc.RenderMarkdown(t, 0))
    
    // Get statistics
    s := stats.Analyze(doc)
    fmt.Println(s.RenderText())
    
    // Lint document
    l := lint.NewLinter()
    issues := l.Lint(doc)
    fmt.Printf("Found %d issues\n", len(issues))
}
```

## Architecture

```
markdownforge/
├── cmd/markdownforge/      # CLI entry point
│   └── main.go
├── internal/
│   ├── parser/             # Markdown AST parser
│   ├── toc/                # Table of contents generator
│   ├── links/              # Link checker
│   ├── extract/            # Content extractor
│   ├── lint/               # Markdown linter
│   ├── convert/            # HTML/text converter
│   ├── stats/              # Document statistics
│   └── badge/              # Badge generator
├── go.mod
├── go.sum
├── LICENSE
├── README.md
├── AGENTS.md
└── SECURITY.md
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `go test ./...`
6. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) for details.
