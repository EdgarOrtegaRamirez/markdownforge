package badge

import (
	"strings"
	"testing"
)

func TestBadgeURL(t *testing.T) {
	b := NewBadge("build", "passing", ColorBrightGreen)
	
	url := b.URL()
	if !strings.Contains(url, "img.shields.io") {
		t.Errorf("expected shields.io URL, got: %s", url)
	}
	if !strings.Contains(url, "build") {
		t.Errorf("expected 'build' in URL, got: %s", url)
	}
	if !strings.Contains(url, "passing") {
		t.Errorf("expected 'passing' in URL, got: %s", url)
	}
}

func TestBadgeMarkdown(t *testing.T) {
	b := NewBadge("status", "active", ColorGreen)
	
	md := b.Markdown()
	if !strings.HasPrefix(md, "![status](") {
		t.Errorf("expected markdown badge, got: %s", md)
	}
	if !strings.HasSuffix(md, ")") {
		t.Errorf("expected closing paren, got: %s", md)
	}
}

func TestBadgeHTML(t *testing.T) {
	b := NewBadge("version", "1.0", ColorBlue)
	
	html := b.HTML()
	if !strings.HasPrefix(html, "<img src=") {
		t.Errorf("expected HTML img tag, got: %s", html)
	}
}

func TestBadgeSet(t *testing.T) {
	bs := NewBadgeSet()
	bs.Add(NewBadge("a", "1", ColorRed))
	bs.Add(NewBadge("b", "2", ColorBlue))
	
	md := bs.Markdown()
	if !strings.Contains(md, "![a](") {
		t.Error("expected badge a in markdown")
	}
	if !strings.Contains(md, "![b](") {
		t.Error("expected badge b in markdown")
	}
	
	html := bs.HTML()
	if !strings.Contains(html, "<img") {
		t.Error("expected img tags in HTML")
	}
}

func TestGenerateProjectBadges(t *testing.T) {
	bs := GenerateProjectBadges("MyProject", "MIT", "Go")
	
	if len(bs.Badges) != 3 {
		t.Errorf("expected 3 badges, got %d", len(bs.Badges))
	}
}

func TestGenerateStatusBadges(t *testing.T) {
	bs := GenerateStatusBadges("passing", "passing", "85%")
	
	if len(bs.Badges) != 3 {
		t.Errorf("expected 3 badges, got %d", len(bs.Badges))
	}
}

func TestBadgeSpecialChars(t *testing.T) {
	b := NewBadge("my-badge", "value_with_under", ColorOrange)
	
	url := b.URL()
	// Should handle special characters
	if url == "" {
		t.Error("expected non-empty URL")
	}
}
