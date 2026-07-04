package links

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"https://example.com", true},
		{"http://example.com/path", true},
		{"ftp://files.example.com", true},
		{"file:///path/to/file", true},
		{"relative/path", false},
		{"", false},
		{"://invalid", false},
	}

	for _, tt := range tests {
		err := ValidateURL(tt.input)
		if (err == nil) != tt.valid {
			t.Errorf("ValidateURL(%q): got err=%v, want valid=%v", tt.input, err, tt.valid)
		}
	}
}

func TestCheckInternal(t *testing.T) {
	// Create temp directory with a file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "README.md")
	os.WriteFile(testFile, []byte("# Test"), 0644)

	checker := NewChecker(tmpDir, 0, 1)

	// Test existing file
	status := &LinkStatus{}
	checker.checkInternal("README.md", status)
	if !status.Valid {
		t.Errorf("expected valid for existing file, got error: %s", status.Error)
	}

	// Test non-existing file
	status = &LinkStatus{}
	checker.checkInternal("nonexistent.md", status)
	if status.Valid {
		t.Error("expected invalid for non-existing file")
	}
}

func TestCheckInternalAnchor(t *testing.T) {
	tmpDir := t.TempDir()
	checker := NewChecker(tmpDir, 0, 1)

	status := &LinkStatus{}
	checker.checkInternal("#section", status)
	if !status.Valid {
		t.Errorf("expected valid for anchor-only link, got error: %s", status.Error)
	}
}

func TestSummary(t *testing.T) {
	results := []*LinkStatus{
		{URL: "https://valid.com", Valid: true},
		{URL: "https://broken.com", Valid: false, Error: "404 Not Found"},
		{URL: "https://also-valid.com", Valid: true},
		{URL: "https://also-broken.com", Valid: false, Error: "Connection refused"},
	}

	valid, invalid, errors := Summary(results)

	if valid != 2 {
		t.Errorf("expected 2 valid, got %d", valid)
	}
	if invalid != 2 {
		t.Errorf("expected 2 invalid, got %d", invalid)
	}
	if len(errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errors))
	}
}
