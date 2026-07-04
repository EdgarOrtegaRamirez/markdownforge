package lint

import (
	"testing"

	"github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
)

func TestLintHeadingLevels(t *testing.T) {
	source := `# Title
### Skipped Level
`
	doc := parser.Parse(source)
	linter := NewLinter()
	issues := linter.Lint(doc)

	found := false
	for _, issue := range issues {
		if issue.Rule == "heading-level" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected heading-level issue")
	}
}

func TestLintTrailingSpace(t *testing.T) {
	source := "Line with trailing   \nGood line\n"
	doc := parser.Parse(source)
	linter := NewLinter()
	issues := linter.Lint(doc)

	found := false
	for _, issue := range issues {
		if issue.Rule == "trailing-space" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected trailing-space issue")
	}
}

func TestLintMultipleBlanks(t *testing.T) {
	source := "Line 1\n\n\n\nLine 2\n"
	doc := parser.Parse(source)
	linter := NewLinter()
	issues := linter.Lint(doc)

	found := false
	for _, issue := range issues {
		if issue.Rule == "multiple-blanks" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected multiple-blanks issue")
	}
}

func TestLintHeadingPunctuation(t *testing.T) {
	source := "# Title.\n## Question?\n"
	doc := parser.Parse(source)
	linter := NewLinter()
	issues := linter.Lint(doc)

	count := 0
	for _, issue := range issues {
		if issue.Rule == "heading-punctuation" {
			count++
		}
	}
	if count != 2 {
		t.Errorf("expected 2 heading-punctuation issues, got %d", count)
	}
}

func TestLintNoIssues(t *testing.T) {
	source := `# Title

Good content here.

## Section

More content.
`
	doc := parser.Parse(source)
	linter := NewLinter()
	issues := linter.Lint(doc)

	// Should only have info-level issues, no warnings
	for _, issue := range issues {
		if issue.Severity == SeverityWarning {
			t.Errorf("unexpected warning: %s", issue.Message)
		}
	}
}

func TestSummary(t *testing.T) {
	issues := []Issue{
		{Severity: SeverityError},
		{Severity: SeverityWarning},
		{Severity: SeverityWarning},
		{Severity: SeverityInfo},
	}

	errors, warnings, infos := Summary(issues)
	if errors != 1 {
		t.Errorf("expected 1 error, got %d", errors)
	}
	if warnings != 2 {
		t.Errorf("expected 2 warnings, got %d", warnings)
	}
	if infos != 1 {
		t.Errorf("expected 1 info, got %d", infos)
	}
}

func TestCustomRule(t *testing.T) {
	source := `# Title
`
	doc := parser.Parse(source)
	linter := NewLinter()

	linter.AddRule(Rule{
		Name:        "test-rule",
		Description: "Test rule",
		Severity:    SeverityInfo,
		Check: func(doc *parser.Document) []Issue {
			return []Issue{{Message: "custom issue", Rule: "test-rule"}}
		},
	})

	issues := linter.Lint(doc)
	found := false
	for _, issue := range issues {
		if issue.Rule == "test-rule" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected custom rule issue")
	}
}
