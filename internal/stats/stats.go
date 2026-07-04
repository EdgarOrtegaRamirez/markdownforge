// Package stats provides document statistics.
package stats

import (
	"fmt"
	"strings"
	"time"
	
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
)

// Stats holds document statistics.
type Stats struct {
	WordCount      int
	CharacterCount int
	LineCount      int
	ParagraphCount int
	HeadingCount   int
	CodeBlockCount int
	LinkCount      int
	ImageCount     int
	TableCount     int
	ListCount      int
	BlockquoteCount int
	ReadingTime    time.Duration
	ReadingTimeStr string
}

// Analyze computes statistics for a parsed document.
func Analyze(doc *parser.Document) *Stats {
	stats := &Stats{
		LineCount: len(doc.Lines),
	}
	
	// Count characters (including newlines)
	for _, line := range doc.Lines {
		stats.CharacterCount += len(line) + 1 // +1 for newline
	}
	
	// Count words
	fullText := strings.Join(doc.Lines, " ")
	words := strings.Fields(fullText)
	stats.WordCount = len(words)
	
	// Count blocks by type
	for _, block := range doc.Blocks {
		switch block.Type {
		case parser.NodeParagraph:
			stats.ParagraphCount++
		case parser.NodeFencedCode, parser.NodeCodeBlock:
			stats.CodeBlockCount++
		case parser.NodeTable:
			stats.TableCount++
		case parser.NodeList:
			stats.ListCount++
		case parser.NodeBlockquote:
			stats.BlockquoteCount++
		}
	}
	
	// Count headings, links, images
	stats.HeadingCount = len(doc.Headings)
	stats.LinkCount = len(doc.Links)
	stats.ImageCount = len(doc.Images)
	
	// Calculate reading time (average 200 words per minute)
	minutes := float64(stats.WordCount) / 200.0
	stats.ReadingTime = time.Duration(minutes * float64(time.Minute))
	stats.ReadingTimeStr = formatDuration(stats.ReadingTime)
	
	return stats
}

// formatDuration formats a duration for display.
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "< 1 min"
	}
	minutes := int(d.Minutes())
	if minutes == 1 {
		return "1 min"
	}
	return fmt.Sprintf("%d mins", minutes)
}

// RenderText renders stats as plain text.
func (s *Stats) RenderText() string {
	var sb strings.Builder
	sb.WriteString("Document Statistics\n")
	sb.WriteString("==================\n")
	sb.WriteString(fmt.Sprintf("Words:          %d\n", s.WordCount))
	sb.WriteString(fmt.Sprintf("Characters:     %d\n", s.CharacterCount))
	sb.WriteString(fmt.Sprintf("Lines:          %d\n", s.LineCount))
	sb.WriteString(fmt.Sprintf("Paragraphs:     %d\n", s.ParagraphCount))
	sb.WriteString(fmt.Sprintf("Headings:       %d\n", s.HeadingCount))
	sb.WriteString(fmt.Sprintf("Code Blocks:    %d\n", s.CodeBlockCount))
	sb.WriteString(fmt.Sprintf("Links:          %d\n", s.LinkCount))
	sb.WriteString(fmt.Sprintf("Images:         %d\n", s.ImageCount))
	sb.WriteString(fmt.Sprintf("Tables:         %d\n", s.TableCount))
	sb.WriteString(fmt.Sprintf("Lists:          %d\n", s.ListCount))
	sb.WriteString(fmt.Sprintf("Blockquotes:    %d\n", s.BlockquoteCount))
	sb.WriteString(fmt.Sprintf("Reading Time:   %s\n", s.ReadingTimeStr))
	return sb.String()
}

// RenderJSON renders stats as JSON.
func (s *Stats) RenderJSON() string {
	return fmt.Sprintf(`{
  "word_count": %d,
  "character_count": %d,
  "line_count": %d,
  "paragraph_count": %d,
  "heading_count": %d,
  "code_block_count": %d,
  "link_count": %d,
  "image_count": %d,
  "table_count": %d,
  "list_count": %d,
  "blockquote_count": %d,
  "reading_time": %q
}`, s.WordCount, s.CharacterCount, s.LineCount, s.ParagraphCount,
		s.HeadingCount, s.CodeBlockCount, s.LinkCount, s.ImageCount,
		s.TableCount, s.ListCount, s.BlockquoteCount, s.ReadingTimeStr)
}

// RenderMarkdown renders stats as a Markdown table.
func (s *Stats) RenderMarkdown() string {
	var sb strings.Builder
	sb.WriteString("| Metric | Value |\n")
	sb.WriteString("|--------|-------|\n")
	sb.WriteString(fmt.Sprintf("| Words | %d |\n", s.WordCount))
	sb.WriteString(fmt.Sprintf("| Characters | %d |\n", s.CharacterCount))
	sb.WriteString(fmt.Sprintf("| Lines | %d |\n", s.LineCount))
	sb.WriteString(fmt.Sprintf("| Paragraphs | %d |\n", s.ParagraphCount))
	sb.WriteString(fmt.Sprintf("| Headings | %d |\n", s.HeadingCount))
	sb.WriteString(fmt.Sprintf("| Code Blocks | %d |\n", s.CodeBlockCount))
	sb.WriteString(fmt.Sprintf("| Links | %d |\n", s.LinkCount))
	sb.WriteString(fmt.Sprintf("| Images | %d |\n", s.ImageCount))
	sb.WriteString(fmt.Sprintf("| Tables | %d |\n", s.TableCount))
	sb.WriteString(fmt.Sprintf("| Lists | %d |\n", s.ListCount))
	sb.WriteString(fmt.Sprintf("| Blockquotes | %d |\n", s.BlockquoteCount))
	sb.WriteString(fmt.Sprintf("| Reading Time | %s |\n", s.ReadingTimeStr))
	return sb.String()
}
