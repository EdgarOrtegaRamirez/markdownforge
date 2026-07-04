// Package convert provides Markdown to HTML conversion.
package convert

import (
	"bytes"
	"fmt"
	"strings"
	
	"github.com/yuin/goldmark"
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
)

// Converter converts Markdown to HTML.
type Converter struct {
	md goldmark.Markdown
}

// NewConverter creates a new Markdown to HTML converter.
func NewConverter() *Converter {
	return &Converter{
		md: goldmark.New(),
	}
}

// ToHTML converts Markdown source to HTML.
func (c *Converter) ToHTML(source string) (string, error) {
	var buf bytes.Buffer
	if err := c.md.Convert([]byte(source), &buf); err != nil {
		return "", fmt.Errorf("converting to HTML: %w", err)
	}
	return buf.String(), nil
}

// ToHTMLFromDoc converts a parsed document to HTML.
func (c *Converter) ToHTMLFromDoc(doc *parser.Document) string {
	var sb strings.Builder
	sb.WriteString("<article>\n")
	
	// Use root children to iterate all blocks including headings
	for _, block := range doc.Root.Children {
		switch block.Type {
		case parser.NodeHeading:
			sb.WriteString(fmt.Sprintf("<h%d>%s</h%d>\n", block.Level, block.Content, block.Level))
		case parser.NodeParagraph:
			sb.WriteString(fmt.Sprintf("<p>%s</p>\n", block.Content))
		case parser.NodeFencedCode:
			if block.Language != "" {
				sb.WriteString(fmt.Sprintf("<pre><code class=\"language-%s\">%s</code></pre>\n", block.Language, escapeHTML(block.Content)))
			} else {
				sb.WriteString(fmt.Sprintf("<pre><code>%s</code></pre>\n", escapeHTML(block.Content)))
			}
		case parser.NodeCodeBlock:
			sb.WriteString(fmt.Sprintf("<pre><code>%s</code></pre>\n", escapeHTML(block.Content)))
		case parser.NodeBlockquote:
			sb.WriteString(fmt.Sprintf("<blockquote><p>%s</p></blockquote>\n", block.Content))
		case parser.NodeHorizontalRule:
			sb.WriteString("<hr>\n")
		case parser.NodeList:
			sb.WriteString("<ul>\n")
			for _, item := range block.Children {
				sb.WriteString(fmt.Sprintf("  <li>%s</li>\n", item.Content))
			}
			sb.WriteString("</ul>\n")
		case parser.NodeTable:
			sb.WriteString(renderTableHTML(block))
		}
	}
	
	sb.WriteString("</article>\n")
	return sb.String()
}

// ToPlainText converts Markdown to plain text.
func (c *Converter) ToPlainText(doc *parser.Document) string {
	var sb strings.Builder
	
	for _, block := range doc.Root.Children {
		switch block.Type {
		case parser.NodeHeading:
			sb.WriteString(strings.Repeat("#", block.Level) + " " + block.Content + "\n\n")
		case parser.NodeParagraph:
			sb.WriteString(block.Content + "\n\n")
		case parser.NodeFencedCode, parser.NodeCodeBlock:
			sb.WriteString(block.Content + "\n\n")
		case parser.NodeBlockquote:
			lines := strings.Split(block.Content, "\n")
			for _, line := range lines {
				sb.WriteString("> " + line + "\n")
			}
			sb.WriteString("\n")
		case parser.NodeHorizontalRule:
			sb.WriteString("---\n\n")
		case parser.NodeList:
			for _, item := range block.Children {
				sb.WriteString("- " + item.Content + "\n")
			}
			sb.WriteString("\n")
		}
	}
	
	return sb.String()
}

// renderTableHTML renders a table node as HTML.
func renderTableHTML(table *parser.Node) string {
	var sb strings.Builder
	sb.WriteString("<table>\n")
	
	for i, row := range table.Children {
		if i == 0 {
			sb.WriteString("<thead>\n")
		} else if i == 1 {
			sb.WriteString("<tbody>\n")
		}
		
		sb.WriteString("  <tr>\n")
		for _, cell := range row.Children {
			tag := "td"
			if i == 0 {
				tag = "th"
			}
			sb.WriteString(fmt.Sprintf("    <%s>%s</%s>\n", tag, cell.Content, tag))
		}
		sb.WriteString("  </tr>\n")
	}
	
	if len(table.Children) > 0 {
		sb.WriteString("</tbody>\n")
	}
	sb.WriteString("</table>\n")
	
	return sb.String()
}

// escapeHTML escapes HTML special characters.
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
