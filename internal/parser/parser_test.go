package parser

import (
	"testing"
)

func TestParseHeadings(t *testing.T) {
	source := `# Title
## Subtitle
### Section
#### Subsection
##### Deep
###### Deepest
`
	doc := Parse(source)

	if len(doc.Headings) != 6 {
		t.Fatalf("expected 6 headings, got %d", len(doc.Headings))
	}

	expected := []struct {
		level int
		text  string
	}{
		{1, "Title"},
		{2, "Subtitle"},
		{3, "Section"},
		{4, "Subsection"},
		{5, "Deep"},
		{6, "Deepest"},
	}

	for i, exp := range expected {
		if doc.Headings[i].Level != exp.level {
			t.Errorf("heading %d: expected level %d, got %d", i, exp.level, doc.Headings[i].Level)
		}
		if doc.Headings[i].Content != exp.text {
			t.Errorf("heading %d: expected %q, got %q", i, exp.text, doc.Headings[i].Content)
		}
	}
}

func TestParseParagraphs(t *testing.T) {
	source := `First paragraph.

Second paragraph.

Third paragraph.
`
	doc := Parse(source)

	paragraphs := 0
	for _, block := range doc.Blocks {
		if block.Type == NodeParagraph {
			paragraphs++
		}
	}

	if paragraphs != 3 {
		t.Errorf("expected 3 paragraphs, got %d", paragraphs)
	}
}

func TestParseFencedCode(t *testing.T) {
	source := "```go\npackage main\n\nfunc main() {}\n```"
	doc := Parse(source)

	if len(doc.Blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(doc.Blocks))
	}

	block := doc.Blocks[0]
	if block.Type != NodeFencedCode {
		t.Errorf("expected fenced code, got %s", block.Type)
	}
	if block.Language != "go" {
		t.Errorf("expected language 'go', got %q", block.Language)
	}
	if block.Content != "package main\n\nfunc main() {}" {
		t.Errorf("unexpected content: %q", block.Content)
	}
}

func TestParseIndentedCode(t *testing.T) {
	source := "    code line 1\n    code line 2\n    code line 3"
	doc := Parse(source)

	if len(doc.Blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(doc.Blocks))
	}

	block := doc.Blocks[0]
	if block.Type != NodeCodeBlock {
		t.Errorf("expected code block, got %s", block.Type)
	}
	if block.Content != "code line 1\ncode line 2\ncode line 3" {
		t.Errorf("unexpected content: %q", block.Content)
	}
}

func TestParseBlockquote(t *testing.T) {
	source := "> This is a quote\n> With multiple lines\n> And more"
	doc := Parse(source)

	if len(doc.Blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(doc.Blocks))
	}

	block := doc.Blocks[0]
	if block.Type != NodeBlockquote {
		t.Errorf("expected blockquote, got %s", block.Type)
	}
}

func TestParseTable(t *testing.T) {
	source := `| Name | Age |
|------|-----|
| Alice | 30 |
| Bob | 25 |`
	doc := Parse(source)

	if len(doc.Blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(doc.Blocks))
	}

	block := doc.Blocks[0]
	if block.Type != NodeTable {
		t.Errorf("expected table, got %s", block.Type)
	}
	if len(block.Children) != 3 {
		t.Errorf("expected 3 rows, got %d", len(block.Children))
	}
}

func TestParseHorizontalRule(t *testing.T) {
	source := `---

***

___`
	doc := Parse(source)

	rules := 0
	for _, block := range doc.Blocks {
		if block.Type == NodeHorizontalRule {
			rules++
		}
	}

	if rules != 3 {
		t.Errorf("expected 3 horizontal rules, got %d", rules)
	}
}

func TestParseList(t *testing.T) {
	source := `- Item 1
- Item 2
- Item 3`
	doc := Parse(source)

	if len(doc.Blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(doc.Blocks))
	}

	block := doc.Blocks[0]
	if block.Type != NodeList {
		t.Errorf("expected list, got %s", block.Type)
	}
	if len(block.Children) != 3 {
		t.Errorf("expected 3 items, got %d", len(block.Children))
	}
}

func TestIsHorizontalRule(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"---", true},
		{"***", true},
		{"___", true},
		{"- - -", true},
		{"* * *", true},
		{"_ _ _", true},
		{"--", false},
		{"**", false},
		{"-- some text --", false},
		{"not a rule", false},
	}

	for _, tt := range tests {
		got := isHorizontalRule(tt.input)
		if got != tt.want {
			t.Errorf("isHorizontalRule(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestIsListStart(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"- Item", true},
		{"* Item", true},
		{"+ Item", true},
		{"1. First", true},
		{"10. Tenth", true},
		{"not a list", false},
		{"-no space", false},
	}

	for _, tt := range tests {
		got := isListStart(tt.input)
		if got != tt.want {
			t.Errorf("isListStart(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseEmpty(t *testing.T) {
	doc := Parse("")
	if len(doc.Headings) != 0 {
		t.Errorf("expected 0 headings, got %d", len(doc.Headings))
	}
	if len(doc.Blocks) != 0 {
		t.Errorf("expected 0 blocks, got %d", len(doc.Blocks))
	}
}

func TestParseMixedContent(t *testing.T) {
	source := `# Title

This is a paragraph.

` + "```" + `
code here
` + "```" + `

- Item 1
- Item 2

> A quote
`
	doc := Parse(source)

	if len(doc.Headings) != 1 {
		t.Errorf("expected 1 heading, got %d", len(doc.Headings))
	}
	if doc.Headings[0].Content != "Title" {
		t.Errorf("expected heading 'Title', got %q", doc.Headings[0].Content)
	}
}

func TestNodeTypeString(t *testing.T) {
	tests := []struct {
		nodeType NodeType
		expected string
	}{
		{NodeDocument, "document"},
		{NodeHeading, "heading"},
		{NodeParagraph, "paragraph"},
		{NodeCodeBlock, "code_block"},
		{NodeFencedCode, "fenced_code"},
		{NodeBlockquote, "blockquote"},
		{NodeList, "list"},
		{NodeListItem, "list_item"},
		{NodeTable, "table"},
		{NodeTableRow, "table_row"},
		{NodeTableCell, "table_cell"},
		{NodeHorizontalRule, "horizontal_rule"},
		{NodeHTML, "html"},
		{NodeLink, "link"},
		{NodeImage, "image"},
		{NodeStrong, "strong"},
		{NodeEmphasis, "emphasis"},
		{NodeCode, "code"},
		{NodeText, "text"},
	}

	for _, tt := range tests {
		got := tt.nodeType.String()
		if got != tt.expected {
			t.Errorf("NodeType(%d).String() = %q, want %q", int(tt.nodeType), got, tt.expected)
		}
	}
}
