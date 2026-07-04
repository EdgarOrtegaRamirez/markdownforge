// Package toc generates tables of contents from parsed Markdown.
package toc

import (
	"fmt"
	"strings"
	
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
)

// TOCEntry represents an entry in the table of contents.
type TOCEntry struct {
	Level  int
	Text   string
	Slug   string
	Children []*TOCEntry
}

// Generate creates a table of contents from a parsed document.
func Generate(doc *parser.Document) *TOCEntry {
	root := &TOCEntry{
		Level: 0,
		Text:  "Root",
	}
	
	stack := []*TOCEntry{root}
	
	for _, heading := range doc.Headings {
		entry := &TOCEntry{
			Level: heading.Level,
			Text:  heading.Content,
			Slug:  generateSlug(heading.Content),
		}
		
		// Find the right parent
		for len(stack) > 1 && stack[len(stack)-1].Level >= entry.Level {
			stack = stack[:len(stack)-1]
		}
		
		parent := stack[len(stack)-1]
		parent.Children = append(parent.Children, entry)
		stack = append(stack, entry)
	}
	
	return root
}

// generateSlug creates a URL-friendly slug from text.
func generateSlug(text string) string {
	slug := strings.ToLower(text)
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == ' ' {
			return r
		}
		return -1
	}, slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	slug = strings.Trim(slug, "-")
	return slug
}

// RenderMarkdown renders the TOC as a Markdown list.
func RenderMarkdown(toc *TOCEntry, depth int) string {
	var sb strings.Builder
	if depth > 0 {
		sb.WriteString("\n")
	}
	
	for _, entry := range toc.Children {
		indent := strings.Repeat("  ", entry.Level-1)
		slug := entry.Slug
		sb.WriteString(fmt.Sprintf("%s- [%s](#%s)\n", indent, entry.Text, slug))
		
		if len(entry.Children) > 0 {
			sb.WriteString(RenderMarkdown(entry, depth+1))
		}
	}
	
	return sb.String()
}

// RenderHTML renders the TOC as an HTML unordered list.
func RenderHTML(toc *TOCEntry) string {
	var sb strings.Builder
	sb.WriteString("<nav class=\"toc\">\n")
	renderHTMLList(toc, &sb, 1)
	sb.WriteString("</nav>\n")
	return sb.String()
}

func renderHTMLList(entry *TOCEntry, sb *strings.Builder, depth int) {
	if len(entry.Children) == 0 {
		return
	}
	
	indent := strings.Repeat("  ", depth)
	sb.WriteString(fmt.Sprintf("%s<ul>\n", indent))
	
	for _, child := range entry.Children {
		sb.WriteString(fmt.Sprintf("%s  <li><a href=\"#%s\">%s</a></li>\n", indent, child.Slug, child.Text))
		if len(child.Children) > 0 {
			renderHTMLList(child, sb, depth+1)
		}
	}
	
	sb.WriteString(fmt.Sprintf("%s</ul>\n", indent))
}

// RenderText renders the TOC as plain text with indentation.
func RenderText(toc *TOCEntry) string {
	var sb strings.Builder
	renderTextEntry(toc, &sb, 0)
	return sb.String()
}

func renderTextEntry(entry *TOCEntry, sb *strings.Builder, depth int) {
	for _, child := range entry.Children {
		indent := strings.Repeat("  ", depth)
		sb.WriteString(fmt.Sprintf("%s%s\n", indent, child.Text))
		if len(child.Children) > 0 {
			renderTextEntry(child, sb, depth+1)
		}
	}
}
