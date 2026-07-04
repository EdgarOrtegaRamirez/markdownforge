package extract

import (
	"strings"
	"testing"

	"github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
)

const testMarkdown = `# Title

## Section A

Content in section A.

## Section B

More content here.

### Subsection B1

Sub content.

` + "```go" + `
package main

func main() {}
` + "```" + `

` + "```python" + `
print("hello")
` + "```" + `

> This is a blockquote.

| Name | Age |
|------|-----|
| Alice | 30 |
| Bob | 25 |

[Link](https://example.com)
![Image](https://example.com/img.png)
`

func TestExtractHeadings(t *testing.T) {
	doc := parser.Parse(testMarkdown)
	ext := NewExtractor(doc)

	// All headings
	all := ext.ExtractHeadings(0)
	if len(all) < 4 {
		t.Errorf("expected at least 4 headings, got %d", len(all))
	}

	// Level 2 only
	h2s := ext.ExtractHeadings(2)
	if len(h2s) < 2 {
		t.Errorf("expected at least 2 h2 headings, got %d", len(h2s))
	}
}

func TestExtractSection(t *testing.T) {
	doc := parser.Parse(testMarkdown)
	ext := NewExtractor(doc)

	content, found := ext.ExtractSection("Section B")
	if !found {
		t.Fatal("expected to find Section B")
	}
	if !strings.Contains(content, "More content here") {
		t.Errorf("expected section content, got: %s", content)
	}
}

func TestExtractCodeBlocks(t *testing.T) {
	doc := parser.Parse(testMarkdown)
	ext := NewExtractor(doc)

	// All code blocks
	all := ext.ExtractCodeBlocks("")
	if len(all) != 2 {
		t.Errorf("expected 2 code blocks, got %d", len(all))
	}

	// Python only
	python := ext.ExtractCodeBlocks("python")
	if len(python) != 1 {
		t.Errorf("expected 1 python code block, got %d", len(python))
	}
}

func TestExtractLinks(t *testing.T) {
	doc := parser.Parse(testMarkdown)
	ext := NewExtractor(doc)

	links := ext.ExtractLinks()
	if len(links) < 1 {
		t.Errorf("expected at least 1 link, got %d", len(links))
	}
}

func TestExtractImages(t *testing.T) {
	doc := parser.Parse(testMarkdown)
	ext := NewExtractor(doc)

	images := ext.ExtractImages()
	if len(images) < 1 {
		t.Errorf("expected at least 1 image, got %d", len(images))
	}
}

func TestExtractBlockquotes(t *testing.T) {
	doc := parser.Parse(testMarkdown)
	ext := NewExtractor(doc)

	quotes := ext.ExtractBlockquotes()
	if len(quotes) != 1 {
		t.Errorf("expected 1 blockquote, got %d", len(quotes))
	}
}

func TestExtractTables(t *testing.T) {
	doc := parser.Parse(testMarkdown)
	ext := NewExtractor(doc)

	tables := ext.ExtractTables()
	if len(tables) != 1 {
		t.Errorf("expected 1 table, got %d", len(tables))
	}
}

func TestExtractByRegex(t *testing.T) {
	doc := parser.Parse(testMarkdown)
	ext := NewExtractor(doc)

	// Find all words starting with 'S'
	results := ext.ExtractByRegex(`\bS\w+`)
	if len(results) == 0 {
		t.Error("expected some matches")
	}
}

func TestExtractMetadata(t *testing.T) {
	source := `---
title: My Document
author: Test
---

# Content
`
	doc := parser.Parse(source)
	ext := NewExtractor(doc)

	meta := ext.ExtractMetadata()
	if meta["title"] != "My Document" {
		t.Errorf("expected title 'My Document', got %q", meta["title"])
	}
}
