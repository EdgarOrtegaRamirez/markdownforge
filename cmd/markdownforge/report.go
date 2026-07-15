package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/EdgarOrtegaRamirez/markdownforge/internal/lint"
)

// printJSONReport outputs the lint report as JSON.
func printJSONReport(issues []lint.Issue, score float64, grade string, totalLines int) {
	type IssueOut struct {
		Rule     string `json:"rule"`
		Severity string `json:"severity"`
		Message  string `json:"message"`
		Line     int    `json:"line"`
		Category string `json:"category"`
	}

	type StatsOut struct {
		TotalIssues int `json:"total_issues"`
		Critical    int `json:"critical"`
		Warning     int `json:"warning"`
		Info        int `json:"info"`
	}

	out := struct {
		TotalLines int        `json:"total_lines"`
		Issues     []IssueOut `json:"issues"`
		Stats      StatsOut   `json:"stats"`
		Score      float64    `json:"score"`
		Grade      string     `json:"grade"`
	}{
		TotalLines: totalLines,
		Stats: StatsOut{
			TotalIssues: len(issues),
		},
		Score: score,
		Grade: grade,
	}

	for _, issue := range issues {
		out.Issues = append(out.Issues, IssueOut{
			Rule:     issue.Rule,
			Severity: string(issue.Severity),
			Message:  issue.Message,
			Line:     issue.Line,
			Category: issue.Category,
		})
		switch issue.Severity {
		case lint.SeverityCritical, lint.SeverityError:
			out.Stats.Critical++
		case lint.SeverityWarning:
			out.Stats.Warning++
		case lint.SeverityInfo:
			out.Stats.Info++
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(out)
}

// printGitHubReport outputs issues in GitHub Actions annotation format.
func printGitHubReport(issues []lint.Issue) {
	for _, issue := range issues {
		switch issue.Severity {
		case lint.SeverityCritical, lint.SeverityError:
			fmt.Fprintf(os.Stderr, "::error file=%s,line=%d,col=%d::[%s] %s\n",
				"<stdin>", issue.Line, issue.Column, issue.Rule, issue.Message)
		case lint.SeverityWarning:
			fmt.Fprintf(os.Stderr, "::warning file=%s,line=%d,col=%d::[%s] %s\n",
				"<stdin>", issue.Line, issue.Column, issue.Rule, issue.Message)
		case lint.SeverityInfo:
			fmt.Fprintf(os.Stderr, "::notice file=%s,line=%d,col=%d::[%s] %s\n",
				"<stdin>", issue.Line, issue.Column, issue.Rule, issue.Message)
		}
	}
}
