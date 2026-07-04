// Package links provides link checking for Markdown documents.
package links

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	
	"github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
)

// LinkStatus represents the status of a checked link.
type LinkStatus struct {
	URL      string
	Valid    bool
	StatusCode int
	Error    string
	Type     string // "external", "internal", "anchor"
}

// Checker checks links in Markdown documents.
type Checker struct {
	client     *http.Client
	cache      map[string]*LinkStatus
	mu         sync.RWMutex
	baseDir    string
	maxWorkers int
}

// NewChecker creates a new link checker.
func NewChecker(baseDir string, timeout time.Duration, maxWorkers int) *Checker {
	if maxWorkers <= 0 {
		maxWorkers = 10
	}
	return &Checker{
		client: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
		cache:      make(map[string]*LinkStatus),
		baseDir:    baseDir,
		maxWorkers: maxWorkers,
	}
}

// CheckAll checks all links in a document.
func (c *Checker) CheckAll(doc *parser.Document) []*LinkStatus {
	var links []string
	
	// Collect all links from the document
	for _, link := range doc.Links {
		if link.URL != "" {
			links = append(links, link.URL)
		}
	}
	for _, img := range doc.Images {
		if img.URL != "" {
			links = append(links, img.URL)
		}
	}
	
	// Deduplicate
	seen := make(map[string]bool)
	var unique []string
	for _, l := range links {
		if !seen[l] {
			seen[l] = true
			unique = append(unique, l)
		}
	}
	
	// Check links concurrently
	results := make([]*LinkStatus, len(unique))
	jobs := make(chan int, len(unique))
	
	var wg sync.WaitGroup
	for i := 0; i < c.maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				results[idx] = c.checkLink(unique[idx])
			}
		}()
	}
	
	for i := range unique {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	
	return results
}

// checkLink checks a single link.
func (c *Checker) checkLink(rawURL string) *LinkStatus {
	// Check cache
	c.mu.RLock()
	if cached, ok := c.cache[rawURL]; ok {
		c.mu.RUnlock()
		return cached
	}
	c.mu.RUnlock()
	
	status := &LinkStatus{URL: rawURL}
	
	// Determine link type
	if strings.HasPrefix(rawURL, "#") {
		status.Type = "anchor"
		status.Valid = true // Anchor links are always valid in context
	} else if strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") {
		status.Type = "external"
		c.checkExternal(rawURL, status)
	} else if strings.HasPrefix(rawURL, "mailto:") {
		status.Type = "external"
		status.Valid = true // Don't validate email links
	} else {
		status.Type = "internal"
		c.checkInternal(rawURL, status)
	}
	
	// Cache result
	c.mu.Lock()
	c.cache[rawURL] = status
	c.mu.Unlock()
	
	return status
}

// checkExternal checks an external HTTP link.
func (c *Checker) checkExternal(rawURL string, status *LinkStatus) {
	resp, err := c.client.Get(rawURL)
	if err != nil {
		status.Valid = false
		status.Error = err.Error()
		return
	}
	defer resp.Body.Close()
	
	status.StatusCode = resp.StatusCode
	status.Valid = resp.StatusCode >= 200 && resp.StatusCode < 400
	if !status.Valid {
		status.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
}

// checkInternal checks an internal file link.
func (c *Checker) checkInternal(rawURL string, status *LinkStatus) {
	// Split URL and anchor
	parts := strings.SplitN(rawURL, "#", 2)
	filePath := parts[0]
	
	if filePath == "" {
		// Anchor-only link
		status.Valid = true
		return
	}
	
	// Resolve path
 fullPath := filepath.Join(c.baseDir, filePath)
	
	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		status.Valid = false
		status.Error = "file not found"
		return
	}
	
	status.Valid = true
}

// Summary returns a summary of check results.
func Summary(results []*LinkStatus) (valid, invalid int, errors []string) {
	for _, r := range results {
		if r.Valid {
			valid++
		} else {
			invalid++
			errors = append(errors, fmt.Sprintf("%s: %s", r.URL, r.Error))
		}
	}
	return
}

// ValidateURL checks if a URL is valid.
func ValidateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme == "" {
		return fmt.Errorf("missing scheme")
	}
	if u.Host == "" && u.Scheme != "file" {
		return fmt.Errorf("missing host")
	}
	return nil
}
