// Package badge generates shields.io-style badges.
package badge

import (
	"fmt"
	"strings"
)

// BadgeColor represents a badge color.
type BadgeColor string

const (
	ColorBrightGreen BadgeColor = "brightgreen"
	ColorGreen       BadgeColor = "green"
	ColorYellow      BadgeColor = "yellow"
	ColorYellowGreen BadgeColor = "yellowgreen"
	ColorOrange      BadgeColor = "orange"
	ColorRed         BadgeColor = "red"
	ColorBlue        BadgeColor = "blue"
	ColorLightGrey   BadgeColor = "lightgrey"
	ColorSlateBlue   BadgeColor = "slateblue"
	ColorF5A623      BadgeColor = "F5A623"
)

// Badge represents a shields.io badge.
type Badge struct {
	Label string
	Value string
	Color BadgeColor
}

// NewBadge creates a new badge.
func NewBadge(label, value string, color BadgeColor) *Badge {
	return &Badge{
		Label: label,
		Value: value,
		Color: color,
	}
}

// URL returns the shields.io badge URL.
func (b *Badge) URL() string {
	label := strings.ReplaceAll(b.Label, "-", "--")
	label = strings.ReplaceAll(label, "_", "__")
	value := strings.ReplaceAll(b.Value, "-", "--")
	value = strings.ReplaceAll(value, "_", "__")
	
	return fmt.Sprintf("https://img.shields.io/badge/%s-%s-%s", 
		urlEncode(label), urlEncode(value), b.Color)
}

// Markdown returns the badge as Markdown.
func (b *Badge) Markdown() string {
	return fmt.Sprintf("![%s](%s)", b.Label, b.URL())
}

// HTML returns the badge as an HTML img tag.
func (b *Badge) HTML() string {
	return fmt.Sprintf(`<img src="%s" alt="%s">`, b.URL(), b.Label)
}

// BadgeSet represents a collection of badges.
type BadgeSet struct {
	Badges []*Badge
}

// NewBadgeSet creates a new badge set.
func NewBadgeSet() *BadgeSet {
	return &BadgeSet{}
}

// Add adds a badge to the set.
func (bs *BadgeSet) Add(badge *Badge) {
	bs.Badges = append(bs.Badges, badge)
}

// Markdown returns all badges as Markdown.
func (bs *BadgeSet) Markdown() string {
	var parts []string
	for _, b := range bs.Badges {
		parts = append(parts, b.Markdown())
	}
	return strings.Join(parts, " ")
}

// HTML returns all badges as HTML.
func (bs *BadgeSet) HTML() string {
	var parts []string
	for _, b := range bs.Badges {
		parts = append(parts, b.HTML())
	}
	return strings.Join(parts, "\n")
}

// urlEncode encodes a string for use in URLs.
func urlEncode(s string) string {
	// Simple percent encoding for shields.io
	s = strings.ReplaceAll(s, "-", "--")
	return s
}

// GenerateProjectBadges creates common project badges.
func GenerateProjectBadges(name, license, language string) *BadgeSet {
	bs := NewBadgeSet()
	
	if name != "" {
		bs.Add(NewBadge("name", name, ColorBlue))
	}
	if license != "" {
		bs.Add(NewBadge("license", license, ColorGreen))
	}
	if language != "" {
		bs.Add(NewBadge("language", language, ColorBlue))
	}
	
	return bs
}

// GenerateStatusBadges creates status/quality badges.
func GenerateStatusBadges(tests, lint, coverage string) *BadgeSet {
	bs := NewBadgeSet()
	
	if tests != "" {
		color := ColorBrightGreen
		if tests == "failing" {
			color = ColorRed
		}
		bs.Add(NewBadge("tests", tests, color))
	}
	if lint != "" {
		color := ColorBrightGreen
		if lint == "failing" {
			color = ColorRed
		}
		bs.Add(NewBadge("lint", lint, color))
	}
	if coverage != "" {
		bs.Add(NewBadge("coverage", coverage, ColorGreen))
	}
	
	return bs
}
