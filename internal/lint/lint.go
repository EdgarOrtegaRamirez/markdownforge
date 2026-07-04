// Package lint provides Markdown linting capabilities.
package lint

import (
	"fmt"
	"strings"
	
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
)

// Severity represents the severity of a lint issue.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Issue represents a lint issue found in the document.
type Issue struct {
	Line     int
	Column   int
	Severity Severity
	Rule     string
	Message  string
}

// Rule represents a lint rule.
type Rule struct {
	Name        string
	Description string
	Severity    Severity
	Check       func(doc *parser.Document) []Issue
}

// Linter checks Markdown documents for issues.
type Linter struct {
	rules []Rule
}

// NewLinter creates a new linter with default rules.
func NewLinter() *Linter {
	l := &Linter{}
	l.AddDefaultRules()
	return l
}

// AddRule adds a custom lint rule.
func (l *Linter) AddRule(rule Rule) {
	l.rules = append(l.rules, rule)
}

// AddDefaultRules adds the default set of lint rules.
func (l *Linter) AddDefaultRules() {
	l.rules = append(l.rules,
		Rule{
			Name:        "heading-level",
			Description: "Headings should not skip levels (e.g., h1 to h3)",
			Severity:    SeverityWarning,
			Check:       checkHeadingLevels,
		},
		Rule{
			Name:        "trailing-space",
			Description: "Lines should not have trailing whitespace",
			Severity:    SeverityInfo,
			Check:       checkTrailingSpace,
		},
		Rule{
			Name:        "multiple-blanks",
			Description: "Multiple consecutive blank lines should be reduced to one",
			Severity:    SeverityInfo,
			Check:       checkMultipleBlanks,
		},
		Rule{
			Name:        "heading-punctuation",
			Description: "Headings should not end with punctuation",
			Severity:    SeverityWarning,
			Check:       checkHeadingPunctuation,
		},
		Rule{
			Name:        "first-line-heading",
			Description: "Document should start with a heading",
			Severity:    SeverityInfo,
			Check:       checkFirstLineHeading,
		},
		Rule{
			Name:        "no-empty-links",
			Description: "Links should not have empty text",
			Severity:    SeverityWarning,
			Check:       checkEmptyLinks,
		},
	)
}

// Lint runs all lint rules on a document.
func (l *Linter) Lint(doc *parser.Document) []Issue {
	var issues []Issue
	for _, rule := range l.rules {
		issues = append(issues, rule.Check(doc)...)
	}
	return issues
}

// checkHeadingLevels checks for skipped heading levels.
func checkHeadingLevels(doc *parser.Document) []Issue {
	var issues []Issue
	var prevLevel int
	
	for _, h := range doc.Headings {
		if prevLevel > 0 && h.Level > prevLevel+1 {
			issues = append(issues, Issue{
				Line:     h.Line,
				Severity: SeverityWarning,
				Rule:     "heading-level",
				Message:  fmt.Sprintf("Heading level jumps from h%d to h%d", prevLevel, h.Level),
			})
		}
		prevLevel = h.Level
	}
	
	return issues
}

// checkTrailingSpace checks for trailing whitespace.
func checkTrailingSpace(doc *parser.Document) []Issue {
	var issues []Issue
	for i, line := range doc.Lines {
		if strings.TrimRight(line, " \t") != line && strings.TrimSpace(line) != "" {
			issues = append(issues, Issue{
				Line:     i + 1,
				Severity: SeverityInfo,
				Rule:     "trailing-space",
				Message:  "Line has trailing whitespace",
			})
		}
	}
	return issues
}

// checkMultipleBlanks checks for multiple consecutive blank lines.
func checkMultipleBlanks(doc *parser.Document) []Issue {
	var issues []Issue
	blankCount := 0
	
	for i, line := range doc.Lines {
		if strings.TrimSpace(line) == "" {
			blankCount++
			if blankCount > 1 {
				issues = append(issues, Issue{
					Line:     i + 1,
					Severity: SeverityInfo,
					Rule:     "multiple-blanks",
					Message:  "Multiple consecutive blank lines",
				})
			}
		} else {
			blankCount = 0
		}
	}
	
	return issues
}

// checkHeadingPunctuation checks if headings end with punctuation.
func checkHeadingPunctuation(doc *parser.Document) []Issue {
	var issues []Issue
	punctuation := []string{".", ",", ";", ":", "!", "?"}
	
	for _, h := range doc.Headings {
		for _, p := range punctuation {
			if strings.HasSuffix(h.Content, p) {
				issues = append(issues, Issue{
					Line:     h.Line,
					Severity: SeverityWarning,
					Rule:     "heading-punctuation",
					Message:  fmt.Sprintf("Heading ends with '%s'", p),
				})
				break
			}
		}
	}
	
	return issues
}

// checkFirstLineHeading checks if the document starts with a heading.
func checkFirstLineHeading(doc *parser.Document) []Issue {
	var issues []Issue
	if len(doc.Headings) == 0 || doc.Headings[0].Line != 1 {
		issues = append(issues, Issue{
			Line:     1,
			Severity: SeverityInfo,
			Rule:     "first-line-heading",
			Message:  "Document does not start with a heading",
		})
	}
	return issues
}

// checkEmptyLinks checks for links with empty text.
func checkEmptyLinks(doc *parser.Document) []Issue {
	var issues []Issue
	for _, link := range doc.Links {
		if link.Content == "" {
			issues = append(issues, Issue{
				Line:     link.Line,
				Severity: SeverityWarning,
				Rule:     "no-empty-links",
				Message:  "Link has empty text",
			})
		}
	}
	return issues
}

// Summary returns a summary of lint issues.
func Summary(issues []Issue) (errors, warnings, infos int) {
	for _, issue := range issues {
		switch issue.Severity {
		case SeverityError:
			errors++
		case SeverityWarning:
			warnings++
		case SeverityInfo:
			infos++
		}
	}
	return
}
