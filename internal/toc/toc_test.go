package toc

import (
	"strings"
	"testing"
	
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
)

func TestGenerateTOC(t *testing.T) {
	source := `# Title
## First Section
Some content
### Subsection
## Second Section
### Another Sub
#### Deep Section
`
	doc := parser.Parse(source)
	toc := Generate(doc)
	
	// The h1 is the root, h2s are children
	if len(toc.Children) != 1 {
		t.Fatalf("expected 1 top-level entry (h1), got %d", len(toc.Children))
	}
	
	if toc.Children[0].Text != "Title" {
		t.Errorf("expected 'Title', got %q", toc.Children[0].Text)
	}
	if toc.Children[0].Level != 1 {
		t.Errorf("expected level 1, got %d", toc.Children[0].Level)
	}
	
	// Title should have children (h2s)
	if len(toc.Children[0].Children) != 2 {
		t.Fatalf("expected 2 h2 entries under Title, got %d", len(toc.Children[0].Children))
	}
	
	first := toc.Children[0].Children[0]
	if first.Text != "First Section" {
		t.Errorf("expected 'First Section', got %q", first.Text)
	}
	
	// First Section should have a child (h3)
	if len(first.Children) != 1 {
		t.Fatalf("expected 1 h3 entry under First Section, got %d", len(first.Children))
	}
	if first.Children[0].Text != "Subsection" {
		t.Errorf("expected 'Subsection', got %q", first.Children[0].Text)
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello World", "hello-world"},
		{"Title with Numbers 123", "title-with-numbers-123"},
		{"Special!@#$Chars", "specialchars"},
		{"  Extra  Spaces  ", "extra-spaces"},
		{"--leading--", "leading"},
	}
	
	for _, tt := range tests {
		got := generateSlug(tt.input)
		if got != tt.want {
			t.Errorf("generateSlug(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestRenderMarkdown(t *testing.T) {
	source := `# Title
## Section A
## Section B
### Sub B1
`
	doc := parser.Parse(source)
	toc := Generate(doc)
	md := RenderMarkdown(toc, 0)
	
	if len(md) == 0 {
		t.Error("expected non-empty markdown output")
	}
	
	// Check that it contains the expected links
	if !strings.Contains(md, "[Title](#title)") {
		t.Errorf("expected '[Title](#title)' in output, got:\n%s", md)
	}
	if !strings.Contains(md, "[Section A](#section-a)") {
		t.Errorf("expected '[Section A](#section-a)' in output, got:\n%s", md)
	}
	if !strings.Contains(md, "[Section B](#section-b)") {
		t.Errorf("expected '[Section B](#section-b)' in output, got:\n%s", md)
	}
}

func TestRenderText(t *testing.T) {
	source := `# Title
## Section A
## Section B
`
	doc := parser.Parse(source)
	toc := Generate(doc)
	text := RenderText(toc)
	
	if !strings.Contains(text, "Title") {
		t.Errorf("expected 'Title' in output, got:\n%s", text)
	}
	if !strings.Contains(text, "Section A") {
		t.Errorf("expected 'Section A' in output, got:\n%s", text)
	}
}

func TestRenderHTML(t *testing.T) {
	source := `# Title
## Section A
`
	doc := parser.Parse(source)
	toc := Generate(doc)
	html := RenderHTML(toc)
	
	if !strings.Contains(html, "<nav class=\"toc\">") {
		t.Errorf("expected '<nav class=\"toc\">' in output, got:\n%s", html)
	}
	if !strings.Contains(html, "<a href=\"#title\">Title</a>") {
		t.Errorf("expected link to title in output, got:\n%s", html)
	}
}

func TestEmptyDocument(t *testing.T) {
	doc := parser.Parse("")
	toc := Generate(doc)
	
	if len(toc.Children) != 0 {
		t.Errorf("expected 0 entries for empty document, got %d", len(toc.Children))
	}
}
