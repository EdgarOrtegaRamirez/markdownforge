// Package extract provides content extraction from Markdown documents.
package extract

import (
	"fmt"
	"regexp"
	"strings"
	
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
)

// Extractor extracts specific content from Markdown documents.
type Extractor struct {
	doc *parser.Document
}

// NewExtractor creates a new content extractor.
func NewExtractor(doc *parser.Document) *Extractor {
	return &Extractor{doc: doc}
}

// ExtractHeadings returns all headings at a specific level.
func (e *Extractor) ExtractHeadings(level int) []*parser.Node {
	var result []*parser.Node
	for _, h := range e.doc.Headings {
		if level == 0 || h.Level == level {
			result = append(result, h)
		}
	}
	return result
}

// ExtractSection returns the content of a section by heading text.
func (e *Extractor) ExtractSection(headingText string) (string, bool) {
	var lines []string
	inSection := false
	targetLevel := 0
	
	for _, block := range e.doc.Root.Children {
		if block.Type == parser.NodeHeading && strings.Contains(block.Content, headingText) {
			inSection = true
			targetLevel = block.Level
			lines = append(lines, fmt.Sprintf("%s %s", strings.Repeat("#", block.Level), block.Content))
			continue
		}
		
		if inSection {
			if block.Type == parser.NodeHeading && block.Level <= targetLevel {
				break
			}
			lines = append(lines, block.Content)
		}
	}
	
	if len(lines) == 0 {
		return "", false
	}
	return strings.Join(lines, "\n"), true
}

// ExtractCodeBlocks returns all code blocks, optionally filtered by language.
func (e *Extractor) ExtractCodeBlocks(language string) []string {
	var result []string
	for _, block := range e.doc.Blocks {
		if block.Type == parser.NodeFencedCode {
			if language == "" || block.Language == language {
				result = append(result, block.Content)
			}
		}
	}
	return result
}

// ExtractLinks returns all URLs from the document.
func (e *Extractor) ExtractLinks() []string {
	var result []string
	seen := make(map[string]bool)
	
	for _, link := range e.doc.Links {
		if link.URL != "" && !seen[link.URL] {
			seen[link.URL] = true
			result = append(result, link.URL)
		}
	}
	
	return result
}

// ExtractImages returns all image URLs from the document.
func (e *Extractor) ExtractImages() []string {
	var result []string
	seen := make(map[string]bool)
	
	for _, img := range e.doc.Images {
		if img.URL != "" && !seen[img.URL] {
			seen[img.URL] = true
			result = append(result, img.URL)
		}
	}
	
	return result
}

// ExtractByRegex extracts text matching a regex pattern.
func (e *Extractor) ExtractByRegex(pattern string) []string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}
	
	var result []string
	for _, line := range e.doc.Lines {
		matches := re.FindAllString(line, -1)
		result = append(result, matches...)
	}
	return result
}

// ExtractBlockquotes returns all blockquotes.
func (e *Extractor) ExtractBlockquotes() []string {
	var result []string
	for _, block := range e.doc.Blocks {
		if block.Type == parser.NodeBlockquote {
			result = append(result, block.Content)
		}
	}
	return result
}

// ExtractTables returns all tables as string representations.
func (e *Extractor) ExtractTables() []string {
	var result []string
	for _, block := range e.doc.Blocks {
		if block.Type == parser.NodeTable {
			result = append(result, renderTable(block))
		}
	}
	return result
}

// renderTable renders a table node as a string.
func renderTable(table *parser.Node) string {
	var sb strings.Builder
	
	for i, row := range table.Children {
		var cells []string
		for _, cell := range row.Children {
			cells = append(cells, cell.Content)
		}
		
		sb.WriteString("| ")
		sb.WriteString(strings.Join(cells, " | "))
		sb.WriteString(" |\n")
		
		if i == 0 {
			// Add separator after header
			seps := make([]string, len(cells))
			for j := range seps {
				seps[j] = "---"
			}
			sb.WriteString("| ")
			sb.WriteString(strings.Join(seps, " | "))
			sb.WriteString(" |\n")
		}
	}
	
	return sb.String()
}

// ExtractMetadata extracts YAML-like metadata from the document.
func (e *Extractor) ExtractMetadata() map[string]string {
	result := make(map[string]string)
	
	// Look for metadata at the beginning of the document
	end := 10
	if end > len(e.doc.Lines) {
		end = len(e.doc.Lines)
	}
	for _, line := range e.doc.Lines[:end] {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "---") {
			continue
		}
		
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" && value != "" {
				result[key] = value
			}
		}
	}
	
	return result
}
