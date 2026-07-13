// Package main provides the CLI entry point for MarkdownForge.
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/EdgarOrtegaRamirez/markdownforge/internal/badge"
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/convert"
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/extract"
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/links"
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/lint"
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/spellcheck"
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/stats"
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/toc"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "markdownforge",
		Short: "A comprehensive Markdown processing toolkit",
		Long:  `MarkdownForge provides tools for processing, analyzing, and converting Markdown documents.`,
	}

	// Add commands
	rootCmd.AddCommand(
		newTOCCmd(),
		newStatsCmd(),
		newLintCmd(),
		newLinksCmd(),
		newExtractCmd(),
		newConvertCmd(),
		newBadgeCmd(),
		newSpellcheckCmd(),
		newVersionCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func readInput(args []string) (string, error) {
	if len(args) > 0 {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return "", fmt.Errorf("reading file: %w", err)
		}
		return string(data), nil
	}

	// Read from stdin
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeCharDevice) == 0 {
		var sb strings.Builder
		buf := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buf)
			if n > 0 {
				sb.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
		return sb.String(), nil
	}

	return "", fmt.Errorf("no input provided")
}

func newTOCCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "toc [file]",
		Short: "Generate table of contents",
		RunE: func(cmd *cobra.Command, args []string) error {
			source, err := readInput(args)
			if err != nil {
				return err
			}

			doc := parser.Parse(source)
			t := toc.Generate(doc)
			output := toc.RenderMarkdown(t, 0)
			fmt.Print(output)
			return nil
		},
	}
	return cmd
}

func newStatsCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "stats [file]",
		Short: "Show document statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			source, err := readInput(args)
			if err != nil {
				return err
			}

			doc := parser.Parse(source)
			s := stats.Analyze(doc)

			switch format {
			case "json":
				fmt.Print(s.RenderJSON())
			case "markdown":
				fmt.Print(s.RenderMarkdown())
			default:
				fmt.Print(s.RenderText())
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json, markdown")
	return cmd
}

func newLintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint [file]",
		Short: "Lint Markdown document",
		RunE: func(cmd *cobra.Command, args []string) error {
			source, err := readInput(args)
			if err != nil {
				return err
			}

			doc := parser.Parse(source)
			l := lint.NewLinter()
			issues := l.Lint(doc)

			if len(issues) == 0 {
				color.Green("✓ No issues found")
				return nil
			}

			errors, warnings, infos := lint.Summary(issues)

			for _, issue := range issues {
				var prefix string
				switch issue.Severity {
				case lint.SeverityError:
					prefix = color.New(color.FgRed).Sprint("error")
				case lint.SeverityWarning:
					prefix = color.New(color.FgYellow).Sprint("warning")
				case lint.SeverityInfo:
					prefix = color.New(color.FgCyan).Sprint("info")
				}
				fmt.Printf("%d:%d %s [%s] %s\n", issue.Line, issue.Column, prefix, issue.Rule, issue.Message)
			}

			fmt.Printf("\n%d errors, %d warnings, %d info\n", errors, warnings, infos)
			return nil
		},
	}
	return cmd
}

func newLinksCmd() *cobra.Command {
	var timeout time.Duration
	cmd := &cobra.Command{
		Use:   "links [file]",
		Short: "Check links in Markdown document",
		RunE: func(cmd *cobra.Command, args []string) error {
			source, err := readInput(args)
			if err != nil {
				return err
			}

			doc := parser.Parse(source)
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			checker := links.NewChecker(dir, timeout, 10)
			results := checker.CheckAll(doc)

			valid, invalid, errors := links.Summary(results)

			for _, r := range results {
				if r.Valid {
					color.Green("✓ %s", r.URL)
				} else {
					color.Red("✗ %s: %s", r.URL, r.Error)
				}
			}

			fmt.Printf("\n%d valid, %d invalid\n", valid, invalid)
			if len(errors) > 0 {
				fmt.Println("\nErrors:")
				for _, e := range errors {
					fmt.Printf("  - %s\n", e)
				}
			}
			return nil
		},
	}
	cmd.Flags().DurationVar(&timeout, "timeout", 5*time.Second, "Request timeout")
	return cmd
}

func newExtractCmd() *cobra.Command {
	var (
		level   int
		lang    string
		section string
		regex   string
	)
	cmd := &cobra.Command{
		Use:   "extract [file]",
		Short: "Extract content from Markdown",
		RunE: func(cmd *cobra.Command, args []string) error {
			source, err := readInput(args)
			if err != nil {
				return err
			}

			doc := parser.Parse(source)
			e := extract.NewExtractor(doc)

			if section != "" {
				content, found := e.ExtractSection(section)
				if !found {
					return fmt.Errorf("section %q not found", section)
				}
				fmt.Print(content)
				return nil
			}

			if lang != "" {
				blocks := e.ExtractCodeBlocks(lang)
				for _, block := range blocks {
					fmt.Println(block)
					fmt.Println("---")
				}
				return nil
			}

			if level > 0 {
				headings := e.ExtractHeadings(level)
				for _, h := range headings {
					fmt.Printf("%s %s\n", strings.Repeat("#", h.Level), h.Content)
				}
				return nil
			}

			if regex != "" {
				results := e.ExtractByRegex(regex)
				for _, r := range results {
					fmt.Println(r)
				}
				return nil
			}

			// Default: extract all
			fmt.Println("=== Links ===")
			for _, link := range e.ExtractLinks() {
				fmt.Println(link)
			}
			fmt.Println("\n=== Images ===")
			for _, img := range e.ExtractImages() {
				fmt.Println(img)
			}
			return nil
		},
	}
	cmd.Flags().IntVarP(&level, "level", "l", 0, "Extract headings at specific level")
	cmd.Flags().StringVarP(&lang, "code", "c", "", "Extract code blocks by language")
	cmd.Flags().StringVarP(&section, "section", "s", "", "Extract section by heading text")
	cmd.Flags().StringVarP(&regex, "regex", "r", "", "Extract by regex pattern")
	return cmd
}

func newConvertCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "convert [file]",
		Short: "Convert Markdown to other formats",
		RunE: func(cmd *cobra.Command, args []string) error {
			source, err := readInput(args)
			if err != nil {
				return err
			}

			doc := parser.Parse(source)
			c := convert.NewConverter()

			switch format {
			case "html":
				html, err := c.ToHTML(source)
				if err != nil {
					return err
				}
				fmt.Print(html)
			case "text":
				fmt.Print(c.ToPlainText(doc))
			default:
				return fmt.Errorf("unsupported format: %s", format)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "html", "Output format: html, text")
	return cmd
}

func newBadgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "badge",
		Short: "Generate shields.io badges",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			license, _ := cmd.Flags().GetString("license")
			lang, _ := cmd.Flags().GetString("language")
			tests, _ := cmd.Flags().GetString("tests")

			bs := badge.GenerateProjectBadges(name, license, lang)
			if tests != "" {
				sb := badge.GenerateStatusBadges(tests, "", "")
				for _, b := range sb.Badges {
					bs.Add(b)
				}
			}

			fmt.Println(bs.Markdown())
			return nil
		},
	}
	cmd.Flags().String("name", "", "Project name")
	cmd.Flags().String("license", "", "License type")
	cmd.Flags().String("language", "", "Programming language")
	cmd.Flags().String("tests", "", "Test status (passing/failing)")
	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("markdownforge %s (commit: %s, built: %s)\n", version, commit, buildDate)
		},
	}
}

func newSpellcheckCmd() *cobra.Command {
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "spellcheck [file]",
		Short: "Check spelling in Markdown document",
		RunE: func(cmd *cobra.Command, args []string) error {
			source, err := readInput(args)
			if err != nil {
				return err
			}

			doc := parser.Parse(source)
			checker := spellcheck.NewChecker()
			issues := checker.Check(doc)

			if len(issues) == 0 {
				color.Green("✓ No spelling issues found")
				return nil
			}

			if jsonOutput {
				fmt.Println(formatJSON(issues))
				return nil
			}

			for _, issue := range issues {
				fmt.Printf("%d: %s → %s\n", issue.Line, issue.Word, issue.Suggestion)
			}

			fmt.Printf("\n%d possible misspelling(s) found\n", len(issues))
			return nil
		},
	}
	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output as JSON")
	return cmd
}

func formatJSON(issues []spellcheck.Issue) string {
	type jsonIssue struct {
		Line       int    `json:"line"`
		Word       string `json:"word"`
		Suggestion string `json:"suggestion"`
		Message    string `json:"message"`
	}
	out := "["
	for i, issue := range issues {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf(`{"line":%d,"word":"%s","suggestion":"%s","message":"%s"}`,
			issue.Line, issue.Word, issue.Suggestion, issue.Message)
	}
	out += "]"
	return out
}
