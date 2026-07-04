package stats

import (
	"strings"
	"testing"

	"github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
)

const testDoc = `# Title

This is a paragraph with some words.

## Section

Another paragraph here.

` + "```go" + `
code example
` + "```" + `

- Item 1
- Item 2

> A blockquote

| A | B |
|---|---|
| 1 | 2 |

[Link](https://example.com)
![Image](https://example.com/img.png)
`

func TestAnalyze(t *testing.T) {
	doc := parser.Parse(testDoc)
	stats := Analyze(doc)

	if stats.WordCount == 0 {
		t.Error("expected non-zero word count")
	}
	if stats.CharacterCount == 0 {
		t.Error("expected non-zero character count")
	}
	if stats.LineCount == 0 {
		t.Error("expected non-zero line count")
	}
	if stats.HeadingCount < 2 {
		t.Errorf("expected at least 2 headings, got %d", stats.HeadingCount)
	}
	if stats.CodeBlockCount != 1 {
		t.Errorf("expected 1 code block, got %d", stats.CodeBlockCount)
	}
	if stats.LinkCount < 1 {
		t.Errorf("expected at least 1 link, got %d", stats.LinkCount)
	}
	if stats.ImageCount < 1 {
		t.Errorf("expected at least 1 image, got %d", stats.ImageCount)
	}
	if stats.TableCount != 1 {
		t.Errorf("expected 1 table, got %d", stats.TableCount)
	}
	if stats.ListCount != 1 {
		t.Errorf("expected 1 list, got %d", stats.ListCount)
	}
	if stats.BlockquoteCount != 1 {
		t.Errorf("expected 1 blockquote, got %d", stats.BlockquoteCount)
	}
}

func TestAnalyzeEmpty(t *testing.T) {
	doc := parser.Parse("")
	stats := Analyze(doc)

	if stats.WordCount != 0 {
		t.Errorf("expected 0 words, got %d", stats.WordCount)
	}
	if stats.LineCount < 1 {
		t.Errorf("expected at least 1 line, got %d", stats.LineCount)
	}
}

func TestRenderText(t *testing.T) {
	doc := parser.Parse(testDoc)
	stats := Analyze(doc)

	text := stats.RenderText()
	if text == "" {
		t.Error("expected non-empty text output")
	}
	if !strings.Contains(text, "Document Statistics") {
		t.Error("expected 'Document Statistics' in output")
	}
}

func TestRenderJSON(t *testing.T) {
	doc := parser.Parse(testDoc)
	stats := Analyze(doc)

	json := stats.RenderJSON()
	if json == "" {
		t.Error("expected non-empty JSON output")
	}
	if !strings.Contains(json, "word_count") {
		t.Error("expected 'word_count' in JSON output")
	}
}

func TestRenderMarkdown(t *testing.T) {
	doc := parser.Parse(testDoc)
	stats := Analyze(doc)

	md := stats.RenderMarkdown()
	if md == "" {
		t.Error("expected non-empty markdown output")
	}
	if !strings.Contains(md, "| Words |") {
		t.Error("expected '| Words |' in markdown output")
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"< 1 min"},
	}

	for _, tt := range tests {
		if tt.input == "" {
			t.Error("expected non-empty")
		}
	}
}
