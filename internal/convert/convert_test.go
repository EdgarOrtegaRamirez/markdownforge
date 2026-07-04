package convert

import (
	"strings"
	"testing"
	
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
)

func TestToHTML(t *testing.T) {
	conv := NewConverter()
	
	html, err := conv.ToHTML("# Hello\n\nThis is a paragraph.")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if html == "" {
		t.Error("expected non-empty HTML output")
	}
}

func TestToHTMLFromDoc(t *testing.T) {
	source := `# Title

## Section

Content here.

` + "```go" + `
code
` + "```" + `

> Quote
`
	doc := parser.Parse(source)
	conv := NewConverter()
	
	html := conv.ToHTMLFromDoc(doc)
	
	if html == "" {
		t.Error("expected non-empty HTML output")
	}
	
	// Check for expected HTML elements
	checks := []string{
		"<article>",
		"<h1>",
		"<h2>",
		"<p>Content here.</p>",
		"<pre><code",
		"<blockquote>",
	}
	
	for _, check := range checks {
		if !containsString(html, check) {
			t.Errorf("expected %q in HTML output, got:\n%s", check, html)
		}
	}
}

func TestToPlainText(t *testing.T) {
	source := `# Title

Content here.

- Item 1
- Item 2
`
	doc := parser.Parse(source)
	conv := NewConverter()
	
	text := conv.ToPlainText(doc)
	
	if text == "" {
		t.Error("expected non-empty plain text output")
	}
	if !strings.Contains(text, "Title") {
		t.Errorf("expected 'Title' in plain text output, got:\n%s", text)
	}
}

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"<script>", "&lt;script&gt;"},
		{"a & b", "a &amp; b"},
		{`"quoted"`, "&quot;quoted&quot;"},
	}
	
	for _, tt := range tests {
		got := escapeHTML(tt.input)
		if got != tt.want {
			t.Errorf("escapeHTML(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
