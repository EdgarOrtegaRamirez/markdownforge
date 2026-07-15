// Package lint provides Markdown linting capabilities.
package lint

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
)

// Severity represents the severity of a lint issue.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityError    Severity = "error"
	SeverityWarning  Severity = "warning"
	SeverityInfo     Severity = "info"
)

// Issue represents a lint issue found in the document.
type Issue struct {
	Line     int
	Column   int
	Severity Severity
	Rule     string
	Message  string
	Context  string
	Category string
}

// Rule represents a lint rule.
type Rule struct {
	Name        string
	Description string
	Severity    Severity
	Check       func(doc *parser.Document, rawLines []string, cfg LintConfig) []Issue
}

// LintConfig holds configuration for the linter.
type LintConfig struct {
	MaxLineLength     int    // 0 = unlimited
	RequireHeadingOrder bool // enforce sequential heading levels
	AllowEmptyAlt     bool   // allow images without alt text
}

// NewDefaultConfig returns default lint configuration.
func NewDefaultConfig() LintConfig {
	return LintConfig{
		MaxLineLength:       120,
		RequireHeadingOrder: true,
		AllowEmptyAlt:       false,
	}
}

// Linter checks Markdown documents for issues.
type Linter struct {
	rules   []Rule
	cfg     LintConfig
	rawLines []string
}

// New creates a linter with the given config.
func New(cfg LintConfig) *Linter {
	l := &Linter{cfg: cfg}
	l.AddDefaultRules()
	return l
}

// NewLinter creates a new linter with default config (for backward compat).
func NewLinter() *Linter {
	return New(NewDefaultConfig())
}

// SetRawLines sets the raw lines for line-based checks.
func (l *Linter) SetRawLines(lines []string) {
	l.rawLines = lines
}

// AddRule adds a custom lint rule.
func (l *Linter) AddRule(rule Rule) {
	l.rules = append(l.rules, rule)
}

// AddDefaultRules adds the default set of lint rules.
func (l *Linter) AddDefaultRules() {
	l.rules = append(l.rules,
	// Heading rules
		Rule{
			Name:        "heading-level",
			Description: "Headings should not skip levels (e.g., h1 to h3)",
			Severity:    SeverityWarning,
			Check:       checkHeadingLevels,
		},
		Rule{
			Name:        "heading-skip",
			Description: "Sequential heading levels without skips",
			Severity:    SeverityWarning,
			Check:       checkHeadingSkip,
		},
		// Link/image rules
		Rule{
			Name:        "empty-link-text",
			Description: "Links should not have empty text",
			Severity:    SeverityCritical,
			Check:       checkEmptyLinkText,
		},
		Rule{
			Name:        "no-empty-links",
			Description: "Links should not have empty text (parser-based)",
			Severity:    SeverityWarning,
			Check:       checkEmptyLinks,
		},
		Rule{
			Name:        "image-broken-syntax",
			Description: "Images must have valid syntax with closing parenthesis",
			Severity:    SeverityCritical,
			Check:       checkBrokenImageSyntax,
		},
		Rule{
			Name:        "empty-image-alt",
			Description: "Images should have descriptive alt text",
			Severity:    SeverityWarning,
			Check:       checkEmptyImages,
		},
		Rule{
			Name:        "alt-text-too-long",
			Description: "Image alt text should be under 150 characters",
			Severity:    SeverityWarning,
			Check:       checkAltTextLength,
		},
		// Spacing rules
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
			Name:        "line-too-long",
			Description: "Lines should not exceed maximum length",
			Severity:    SeverityWarning,
			Check:       checkLongLines,
		},
		// Structure rules
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
	)
}

// Lint runs all lint rules on a document.
func (l *Linter) Lint(doc *parser.Document) []Issue {
	var issues []Issue
	for _, rule := range l.rules {
		issues = append(issues, rule.Check(doc, l.rawLines, l.cfg)...)
	}
	return issues
}

// === Heading rules ===

func checkHeadingLevels(doc *parser.Document, _ []string, _ LintConfig) []Issue {
	var issues []Issue
	var prevLevel int

	for _, h := range doc.Headings {
		if prevLevel > 0 && h.Level > prevLevel+1 {
			issues = append(issues, Issue{
				Line:     h.Line,
				Severity: SeverityWarning,
				Rule:     "heading-level",
				Message:  fmt.Sprintf("Heading level jumps from h%d to h%d", prevLevel, h.Level),
				Category: "heading",
			})
		}
		prevLevel = h.Level
	}
	return issues
}

func checkHeadingSkip(doc *parser.Document, rawLines []string, cfg LintConfig) []Issue {
	if !cfg.RequireHeadingOrder {
		return nil
	}
	var issues []Issue
	var headingStack []int

	for i, line := range rawLines {
		trimmed := strings.TrimSpace(line)
		// Skip code blocks
		if strings.HasPrefix(trimmed, "```") {
			continue
		}
		headingMatch := regexp.MustCompile(`^(#{1,6})\s+(.*)`).FindStringSubmatch(trimmed)
		if headingMatch != nil {
			level := len(headingMatch[1])
			if len(headingStack) > 0 {
				lastLevel := headingStack[len(headingStack)-1]
				if level > lastLevel+1 {
					issues = append(issues, Issue{
						Line:     i + 1,
						Severity: SeverityWarning,
						Rule:     "heading-skip",
						Message:  fmt.Sprintf("Heading level jumps from H%d to H%d — consider adding H%d", lastLevel, level, lastLevel+1),
						Category: "heading",
					})
				}
			}
			headingStack = append(headingStack, level)
		}
	}
	return issues
}

// === Link/Image rules ===

func checkEmptyLinkText(doc *parser.Document, rawLines []string, cfg LintConfig) []Issue {
	var issues []Issue
	for i, line := range rawLines {
		trimmed := strings.TrimSpace(line)
		// Skip code blocks
		if strings.HasPrefix(trimmed, "```") {
			continue
		}
		// Check for empty link text: []()
		matches := regexp.MustCompile(`\[([^\]]*)\]\(`).FindAllStringSubmatch(trimmed, -1)
		for _, m := range matches {
			if strings.TrimSpace(m[1]) == "" {
				issues = append(issues, Issue{
					Line:     i + 1,
					Severity: SeverityCritical,
					Rule:     "empty-link-text",
					Message:  "Link has empty text — add descriptive link text",
					Category: "link",
				})
			}
		}
	}
	return issues
}

func checkEmptyLinks(doc *parser.Document, _ []string, _ LintConfig) []Issue {
	var issues []Issue
	for _, link := range doc.Links {
		if link.Content == "" {
			issues = append(issues, Issue{
				Line:     link.Line,
				Severity: SeverityWarning,
				Rule:     "no-empty-links",
				Message:  "Link has empty text",
				Category: "link",
			})
		}
	}
	return issues
}

func checkBrokenImageSyntax(doc *parser.Document, rawLines []string, cfg LintConfig) []Issue {
	var issues []Issue
	for i, line := range rawLines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			continue
		}
		// Check for broken image: has ![...]( but no closing )
		if regexp.MustCompile(`!\[[^\]]*\]\(`).MatchString(trimmed) {
			if !regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)$`).MatchString(trimmed) {
				issues = append(issues, Issue{
					Line:     i + 1,
					Severity: SeverityCritical,
					Rule:     "image-broken-syntax",
					Message:  "Image syntax is broken — missing closing parenthesis",
					Category: "image",
				})
			}
		}
	}
	return issues
}

func checkEmptyImages(doc *parser.Document, rawLines []string, cfg LintConfig) []Issue {
	if cfg.AllowEmptyAlt {
		return nil
	}
	var issues []Issue
	for i, line := range rawLines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			continue
		}
		if regexp.MustCompile(`!\[\]\([^)]+\)`).MatchString(trimmed) {
			issues = append(issues, Issue{
				Line:     i + 1,
				Severity: SeverityWarning,
				Rule:     "empty-image-alt",
				Message:  "Image has empty alt text — consider adding descriptive alt text or using decorative marker",
				Category: "image",
			})
		}
	}
	return issues
}

func checkAltTextLength(doc *parser.Document, rawLines []string, cfg LintConfig) []Issue {
	if cfg.AllowEmptyAlt {
		return nil
	}
	var issues []Issue
	for i, line := range rawLines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			continue
		}
		m := regexp.MustCompile(`!\[([^\]]*)\]\(`).FindStringSubmatch(trimmed)
		if m != nil {
			alt := m[1]
			if len(alt) > 150 {
				issues = append(issues, Issue{
					Line:     i + 1,
					Severity: SeverityWarning,
					Rule:     "alt-text-too-long",
					Message:  fmt.Sprintf("Image alt text is %d characters — consider shortening or moving to caption", len(alt)),
					Category: "image",
				})
			}
		}
	}
	return issues
}

// === Spacing rules ===

func checkTrailingSpace(doc *parser.Document, rawLines []string, cfg LintConfig) []Issue {
	var issues []Issue
	for i, line := range rawLines {
		if strings.TrimRight(line, " \t") != line && strings.TrimSpace(line) != "" {
			issues = append(issues, Issue{
				Line:     i + 1,
				Severity: SeverityInfo,
				Rule:     "trailing-space",
				Message:  "Line has trailing whitespace",
				Category: "spacing",
			})
		}
	}
	return issues
}

func checkMultipleBlanks(doc *parser.Document, _ []string, _ LintConfig) []Issue {
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
					Category: "spacing",
				})
			}
		} else {
			blankCount = 0
		}
	}
	return issues
}

func checkLongLines(doc *parser.Document, rawLines []string, cfg LintConfig) []Issue {
	if cfg.MaxLineLength == 0 {
		return nil
	}
	var issues []Issue
	for i, line := range rawLines {
		if len(line) > cfg.MaxLineLength {
			issues = append(issues, Issue{
				Line:     i + 1,
				Severity: SeverityWarning,
				Rule:     "line-too-long",
				Message:  fmt.Sprintf("Line is %d characters (max %d)", len(line), cfg.MaxLineLength),
				Category: "spacing",
			})
		}
	}
	return issues
}

// === Structure rules ===

func checkHeadingPunctuation(doc *parser.Document, _ []string, _ LintConfig) []Issue {
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
					Category: "heading",
				})
				break
			}
		}
	}
	return issues
}

func checkFirstLineHeading(doc *parser.Document, _ []string, _ LintConfig) []Issue {
	var issues []Issue
	if len(doc.Headings) == 0 || doc.Headings[0].Line != 1 {
		issues = append(issues, Issue{
			Line:     1,
			Severity: SeverityInfo,
			Rule:     "first-line-heading",
			Message:  "Document does not start with a heading",
			Category: "structure",
		})
	}
	return issues
}

// === Scoring ===

// CalculateScore computes a 0-100 quality score based on issues.
func CalculateScore(issues []Issue) float64 {
	if len(issues) == 0 {
		return 100.0
	}

	criticalPenalty := 15.0
	warningPenalty := 5.0
	infoPenalty := 1.0

	penalty := 0.0
	for _, issue := range issues {
		switch issue.Severity {
		case SeverityCritical, SeverityError:
			penalty += criticalPenalty
		case SeverityWarning:
			penalty += warningPenalty
		case SeverityInfo:
			penalty += infoPenalty
		}
	}

	score := 100.0 - penalty
	if score < 0 {
		score = 0
	}
	return float64(int(score*10)) / 10
}

// ScoreToGrade converts a score to a letter grade.
func ScoreToGrade(score float64) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}

// GradeSymbol returns a visual symbol for the grade.
func GradeSymbol(grade string) string {
	symbols := map[string]string{
		"A": "🌟",
		"B": "✅",
		"C": "👍",
		"D": "⚠️",
		"F": "❌",
	}
	if s, ok := symbols[grade]; ok {
		return s + " " + grade
	}
	return grade
}

// Summary returns a summary of lint issues.
func Summary(issues []Issue) (errors, warnings, infos int) {
	for _, issue := range issues {
		switch issue.Severity {
		case SeverityCritical, SeverityError:
			errors++
		case SeverityWarning:
			warnings++
		case SeverityInfo:
			infos++
		}
	}
	return
}
