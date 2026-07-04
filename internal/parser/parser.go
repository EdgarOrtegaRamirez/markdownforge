// Package parser implements Markdown document parsing into an AST.
package parser

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// NodeType represents the type of a Markdown node.
type NodeType int

const (
	NodeDocument NodeType = iota
	NodeHeading
	NodeParagraph
	NodeCodeBlock
	NodeFencedCode
	NodeBlockquote
	NodeList
	NodeListItem
	NodeTable
	NodeTableRow
	NodeTableCell
	NodeHorizontalRule
	NodeHTML
	NodeLink
	NodeImage
	NodeStrong
	NodeEmphasis
	NodeCode
	NodeText
	NodeSoftBreak
	NodeHardBreak
)

// Node represents a node in the Markdown AST.
type Node struct {
	Type     NodeType
	Level    int      // For headings (1-6)
	Content  string   // Text content
	Language string   // For fenced code blocks
	URL      string   // For links and images
	Title    string   // For links and images
	Alt      string   // For images
	Children []*Node  // Child nodes
	Parent   *Node    // Parent node
	Line     int      // Line number in source
	Attrs    map[string]string // For tables
}

// Document represents a parsed Markdown document.
type Document struct {
	Root     *Node
	Headings []*Node
	Links    []*Node
	Images   []*Node
	Blocks   []*Node
	Lines    []string
}

// Parse parses a Markdown string into a Document AST.
func Parse(source string) *Document {
	doc := &Document{
		Root:  &Node{Type: NodeDocument},
		Lines: strings.Split(source, "\n"),
	}
	
	parser := &blockParser{
		source: source,
		lines:  doc.Lines,
		doc:    doc,
		pos:    0,
	}
	
	parser.parseBlocks(doc.Root)
	
	// Collect nodes by type
	collectNodes(doc.Root, doc)
	
	return doc
}

// collectNodes traverses the AST and collects nodes by type.
func collectNodes(root *Node, doc *Document) {
	for _, child := range root.Children {
		switch child.Type {
		case NodeHeading:
			doc.Headings = append(doc.Headings, child)
			// Parse inline elements in headings
			parseInline(child, doc)
		case NodeLink:
			doc.Links = append(doc.Links, child)
		case NodeImage:
			doc.Images = append(doc.Images, child)
		case NodeParagraph, NodeBlockquote:
			// Parse inline elements in paragraphs and blockquotes
			parseInline(child, doc)
			doc.Blocks = append(doc.Blocks, child)
		case NodeCodeBlock, NodeFencedCode, NodeList, NodeTable, NodeHorizontalRule, NodeHTML:
			doc.Blocks = append(doc.Blocks, child)
		}
		collectNodes(child, doc)
	}
}

// parseInline extracts inline elements (links, images) from text content.
func parseInline(node *Node, doc *Document) {
	content := node.Content
	
	// Extract images: ![alt](url "title")
	imgRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)\s]+)(?:\s+"([^"]*)")?\)`)
	for _, match := range imgRegex.FindAllStringSubmatch(content, -1) {
		img := &Node{
			Type:    NodeImage,
			Content: match[1],
			URL:     match[2],
			Alt:     match[1],
			Line:    node.Line,
		}
		if len(match) > 3 {
			img.Title = match[3]
		}
		doc.Images = append(doc.Images, img)
		node.Children = append(node.Children, img)
	}
	
	// Extract links: [text](url "title")
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)\s]+)(?:\s+"([^"]*)")?\)`)
	for _, match := range linkRegex.FindAllStringSubmatch(content, -1) {
		link := &Node{
			Type:    NodeLink,
			Content: match[1],
			URL:     match[2],
			Line:    node.Line,
		}
		if len(match) > 3 {
			link.Title = match[3]
		}
		doc.Links = append(doc.Links, link)
		node.Children = append(node.Children, link)
	}
}

// blockParser parses block-level Markdown elements.
type blockParser struct {
	source string
	lines  []string
	doc    *Document
	pos    int
}

// parseBlocks parses block-level elements.
func (p *blockParser) parseBlocks(parent *Node) {
	for p.pos < len(p.lines) {
		line := p.lines[p.pos]
		trimmed := strings.TrimSpace(line)
		
		// Empty line
		if trimmed == "" {
			p.pos++
			continue
		}
		
		// Heading
		if level, text := parseHeading(line); level > 0 {
			node := &Node{
				Type:    NodeHeading,
				Level:   level,
				Content: strings.TrimSpace(text),
				Line:    p.pos + 1,
			}
			parent.Children = append(parent.Children, node)
			p.pos++
			continue
		}
		
		// Horizontal rule
		if isHorizontalRule(trimmed) {
			node := &Node{
				Type: NodeHorizontalRule,
				Line: p.pos + 1,
			}
			parent.Children = append(parent.Children, node)
			p.pos++
			continue
		}
		
		// Fenced code block
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			node := p.parseFencedCode()
			if node != nil {
				parent.Children = append(parent.Children, node)
			}
			continue
		}
		
		// Indented code block
		if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
			node := p.parseIndentedCode()
			if node != nil {
				parent.Children = append(parent.Children, node)
			}
			continue
		}
		
		// Blockquote
		if strings.HasPrefix(trimmed, ">") {
			node := p.parseBlockquote()
			if node != nil {
				parent.Children = append(parent.Children, node)
			}
			continue
		}
		
		// Table
		if isTableStart(p.lines, p.pos) {
			node := p.parseTable()
			if node != nil {
				parent.Children = append(parent.Children, node)
			}
			continue
		}
		
		// List
		if isListStart(trimmed) {
			node := p.parseList()
			if node != nil {
				parent.Children = append(parent.Children, node)
			}
			continue
		}
		
		// HTML block
		if strings.HasPrefix(trimmed, "<") && (strings.HasSuffix(trimmed, ">") || 
			strings.Contains(trimmed, "</")) {
			node := &Node{
				Type:    NodeHTML,
				Content: line,
				Line:    p.pos + 1,
			}
			parent.Children = append(parent.Children, node)
			p.pos++
			continue
		}
		
		// Paragraph
		node := p.parseParagraph()
		if node != nil {
			parent.Children = append(parent.Children, node)
		}
	}
}

// parseHeading parses a heading line. Returns level and text.
func parseHeading(line string) (int, string) {
	trimmed := strings.TrimSpace(line)
	level := 0
	for i, ch := range trimmed {
		if ch == '#' {
			level++
		} else {
			if level > 0 && level <= 6 && (i == level) && (i >= len(trimmed) || trimmed[i] == ' ' || trimmed[i] == '\t') {
				text := strings.TrimSpace(trimmed[i:])
				return level, text
			}
			return 0, ""
		}
	}
	return 0, ""
}

// isHorizontalRule checks if a line is a horizontal rule.
func isHorizontalRule(line string) bool {
	count := 0
	var char rune
	for _, ch := range line {
		if ch == '-' || ch == '*' || ch == '_' {
			if char == 0 {
				char = ch
			}
			if ch != char {
				return false
			}
			count++
		} else if ch != ' ' && ch != '\t' {
			return false
		}
	}
	return count >= 3
}

// parseFencedCode parses a fenced code block.
func (p *blockParser) parseFencedCode() *Node {
	line := strings.TrimSpace(p.lines[p.pos])
	fence := ""
	if strings.HasPrefix(line, "```") {
		fence = "```"
	} else if strings.HasPrefix(line, "~~~") {
		fence = "~~~"
	}
	lang := strings.TrimSpace(line[len(fence):])
	
	node := &Node{
		Type:     NodeFencedCode,
		Language: lang,
		Line:     p.pos + 1,
	}
	
	p.pos++
	var content []string
	for p.pos < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.pos])
		if strings.HasPrefix(line, fence) && len(strings.Trim(line, fence[0:1])) == 0 {
			p.pos++
			break
		}
		content = append(content, p.lines[p.pos])
		p.pos++
	}
	
	node.Content = strings.Join(content, "\n")
	return node
}

// parseIndentedCode parses an indented code block.
func (p *blockParser) parseIndentedCode() *Node {
	node := &Node{
		Type: NodeCodeBlock,
		Line: p.pos + 1,
	}
	
	var content []string
	for p.pos < len(p.lines) {
		line := p.lines[p.pos]
		if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
			content = append(content, strings.TrimPrefix(strings.TrimPrefix(line, "    "), "\t"))
			p.pos++
		} else {
			break
		}
	}
	
	node.Content = strings.Join(content, "\n")
	return node
}

// parseBlockquote parses a blockquote.
func (p *blockParser) parseBlockquote() *Node {
	node := &Node{
		Type: NodeBlockquote,
		Line: p.pos + 1,
	}
	
	var content []string
	for p.pos < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.pos])
		if strings.HasPrefix(line, ">") {
			content = append(content, strings.TrimPrefix(line, ">"))
			p.pos++
		} else if line == "" {
			p.pos++
		} else {
			break
		}
	}
	
	node.Content = strings.Join(content, "\n")
	return node
}

// isTableStart checks if a line starts a table.
func isTableStart(lines []string, pos int) bool {
	if pos >= len(lines) {
		return false
	}
	line := strings.TrimSpace(lines[pos])
	if !strings.Contains(line, "|") {
		return false
	}
	if pos+1 < len(lines) {
		next := strings.TrimSpace(lines[pos+1])
		if strings.Contains(next, "---") && strings.Contains(next, "|") {
			return true
		}
	}
	return false
}

// parseTable parses a table.
func (p *blockParser) parseTable() *Node {
	node := &Node{
		Type: NodeTable,
		Line: p.pos + 1,
	}
	
	// Parse header
	headerLine := strings.TrimSpace(p.lines[p.pos])
	cells := splitTableRow(headerLine)
	headerRow := &Node{Type: NodeTableRow}
	for _, cell := range cells {
		headerRow.Children = append(headerRow.Children, &Node{
			Type:    NodeTableCell,
			Content: strings.TrimSpace(cell),
		})
	}
	node.Children = append(node.Children, headerRow)
	p.pos++
	
	// Skip separator
	if p.pos < len(p.lines) {
		p.pos++
	}
	
	// Parse body rows
	for p.pos < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.pos])
		if !strings.Contains(line, "|") || line == "" {
			break
		}
		
		row := &Node{Type: NodeTableRow}
		cells := splitTableRow(line)
		for _, cell := range cells {
			row.Children = append(row.Children, &Node{
				Type:    NodeTableCell,
				Content: strings.TrimSpace(cell),
			})
		}
		node.Children = append(node.Children, row)
		p.pos++
	}
	
	return node
}

// splitTableRow splits a table row into cells.
func splitTableRow(line string) []string {
	line = strings.Trim(line, "| ")
	return strings.Split(line, "|")
}

// isListStart checks if a line starts a list.
func isListStart(line string) bool {
	if len(line) < 2 {
		return false
	}
	if (line[0] == '-' || line[0] == '*' || line[0] == '+') && (len(line) == 1 || line[1] == ' ' || line[1] == '\t') {
		return true
	}
	if unicode.IsDigit(rune(line[0])) {
		for i, ch := range line {
			if ch == '.' && i > 0 {
				return true
			}
			if !unicode.IsDigit(ch) && ch != '.' {
				return false
			}
		}
	}
	return false
}

// parseList parses a list.
func (p *blockParser) parseList() *Node {
	node := &Node{
		Type: NodeList,
		Line: p.pos + 1,
	}
	
	for p.pos < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.pos])
		if !isListStart(line) && line != "" {
			break
		}
		if line == "" {
			p.pos++
			continue
		}
		
		item := &Node{
			Type: NodeListItem,
			Line: p.pos + 1,
		}
		
		// Extract list marker and text
		if line[0] == '-' || line[0] == '*' || line[0] == '+' {
			item.Content = strings.TrimSpace(line[2:])
		} else {
			// Ordered list
			for i, ch := range line {
				if ch == '.' {
					item.Content = strings.TrimSpace(line[i+1:])
					break
				}
			}
		}
		
		node.Children = append(node.Children, item)
		p.pos++
	}
	
	return node
}

// parseParagraph parses a paragraph.
func (p *blockParser) parseParagraph() *Node {
	node := &Node{
		Type: NodeParagraph,
		Line: p.pos + 1,
	}
	
	var lines []string
	for p.pos < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.pos])
		if line == "" || isHorizontalRule(line) {
			break
		}
		if level, _ := parseHeading(p.lines[p.pos]); level > 0 {
			break
		}
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			break
		}
		if strings.HasPrefix(line, ">") {
			break
		}
		if isTableStart(p.lines, p.pos) {
			break
		}
		if isListStart(line) {
			break
		}
		lines = append(lines, line)
		p.pos++
	}
	
	node.Content = strings.Join(lines, " ")
	return node
}

// String returns a string representation of the node type.
func (t NodeType) String() string {
	switch t {
	case NodeDocument:
		return "document"
	case NodeHeading:
		return "heading"
	case NodeParagraph:
		return "paragraph"
	case NodeCodeBlock:
		return "code_block"
	case NodeFencedCode:
		return "fenced_code"
	case NodeBlockquote:
		return "blockquote"
	case NodeList:
		return "list"
	case NodeListItem:
		return "list_item"
	case NodeTable:
		return "table"
	case NodeTableRow:
		return "table_row"
	case NodeTableCell:
		return "table_cell"
	case NodeHorizontalRule:
		return "horizontal_rule"
	case NodeHTML:
		return "html"
	case NodeLink:
		return "link"
	case NodeImage:
		return "image"
	case NodeStrong:
		return "strong"
	case NodeEmphasis:
		return "emphasis"
	case NodeCode:
		return "code"
	case NodeText:
		return "text"
	case NodeSoftBreak:
		return "soft_break"
	case NodeHardBreak:
		return "hard_break"
	default:
		return fmt.Sprintf("unknown(%d)", int(t))
	}
}
